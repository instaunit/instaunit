package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/instaunit/instaunit/hunit"
	"github.com/instaunit/instaunit/hunit/cache"
	"github.com/instaunit/instaunit/hunit/doc"
	"github.com/instaunit/instaunit/hunit/exec"
	"github.com/instaunit/instaunit/hunit/net/await"
	"github.com/instaunit/instaunit/hunit/report"
	"github.com/instaunit/instaunit/hunit/service"
	"github.com/instaunit/instaunit/hunit/service/backend/rest"
	"github.com/instaunit/instaunit/hunit/syncio"
	"github.com/instaunit/instaunit/hunit/test"
	"github.com/instaunit/instaunit/hunit/text"

	doc_emit "github.com/instaunit/instaunit/hunit/doc/emit"
	report_emit "github.com/instaunit/instaunit/hunit/report/emit"

	"github.com/bww/go-util/debug"
	"github.com/fatih/color"
)

const (
	stdinPath = "-"
	cacheBase = "./.instaunit/cache"
)

var ( // set at compile time via the linker
	version = "v0.0.0"
	githash = "000000"
)

var (
	colorErr   = []color.Attribute{color.FgYellow}
	colorSuite = []color.Attribute{color.Bold}
)

var syncStdout = syncio.NewWriter(os.Stdout)

// You know what it does
func main() {
	os.Exit(app())
}

// You know what it does
func app() int {
	var tests, skipped, failures, errno int
	var headerSpecs, serviceSpecs flagList

	cmdline := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var (
		fBaseURL         = cmdline.String("base-url", coalesce(os.Getenv("HUNIT_BASE_URL"), "http://localhost/"), "The base URL for requests. Overrides: $HUNIT_BASE_URL.")
		fWait            = cmdline.Duration("wait", strToDuration(os.Getenv("HUNIT_WAIT_ON_START")), "Wait an interval before test suites are run to allow services to settle. Overrides: $HUNIT_WAIT_ON_START.")
		fExpandVars      = cmdline.Bool("expand", strToBool(os.Getenv("HUNIT_EXPAND_VARS"), true), "Expand variables in test cases. Overrides: $HUNIT_EXPAND_VARS.")
		fTrimEntity      = cmdline.Bool("entity:trim", strToBool(os.Getenv("HUNIT_TRIM_ENTITY"), true), "Trim trailing whitespace from entities. Overrides: $HUNIT_TRIM_ENTITY.")
		fDumpRequest     = cmdline.Bool("dump:request", strToBool(os.Getenv("HUNIT_DUMP_REQUESTS")), "Dump requests to standard output as they are processed. Overrides: $HUNIT_DUMP_REQUESTS.")
		fDumpResponse    = cmdline.Bool("dump:response", strToBool(os.Getenv("HUNIT_DUMP_RESPONSES")), "Dump responses to standard output as they are processed. Overrides: $HUNIT_DUMP_RESPONSES.")
		fGendoc          = cmdline.Bool("gendoc", strToBool(os.Getenv("HUNIT_GENDOC")), "Generate documentation. Overrides: $HUNIT_GENDOC.")
		fDocpath         = cmdline.String("doc:output", coalesce(os.Getenv("HUNIT_DOC_OUTPUT"), "./docs"), "The directory in which generated documentation should be written. Overrides: $HUNIT_DOC_OUTPUT.")
		fDoctype         = cmdline.String("doc:type", coalesce(os.Getenv("HUNIT_DOC_TYPE"), "markdown"), "The format to generate documentation in. Overrides: $HUNIT_DOC_TYPE.")
		fDocInclHTTP     = cmdline.Bool("doc:include-http", strToBool(os.Getenv("HUNIT_DOC_INCLUDE_HTTP")), "Include HTTP in request and response examples (as opposed to just routes and entities). Overrides: $HUNIT_DOC_INCLUDE_HTTP.")
		fDocFormatEntity = cmdline.Bool("doc:format-entities", strToBool(os.Getenv("HUNIT_DOC_FORMAT_ENTITIES")), "Pretty-print supported request and response entities in documentation output. Overrides: $HUNIT_DOC_FORMAT_ENTITIES.")
		fReport          = cmdline.Bool("report", strToBool(os.Getenv("HUNIT_REPORT")), "Generate a report. Overrides: $HUNIT_REPORT.")
		fReportPath      = cmdline.String("report:output", coalesce(os.Getenv("HUNIT_REPORT_OUTPUT"), "./reports"), "The directory in which generated reports should be written. Overrides: $HUNIT_REPORT_OUTPUT.")
		fReportType      = cmdline.String("report:type", coalesce(os.Getenv("HUNIT_REPORT_TYPE"), "junit"), "The format to generate reports in. Overrides: $HUNIT_REPORT_TYPE.")
		fCache           = cmdline.Bool("cache", strToBool(os.Getenv("HUNIT_CACHE_RESULTS")), "Cache results. When enabled, test suites run against a managed service will cache results if neither the service binary nor the test suite has changed. Overrides: $HUNIT_CACHE_RESULTS.")
		fIOGracePeriod   = cmdline.Duration("net:grace-period", strToDuration(os.Getenv("HUNIT_NET_IO_GRACE_PERIOD")), "The grace period to wait for long-running I/O to complete before shutting down websocket/persistent connections. Overrides: $HUNIT_NET_IO_GRACE_PERIOD.")
		fExec            = cmdline.String("exec", os.Getenv("HUNIT_EXEC_COMMAND"), "The command to execute before running tests, usually the program that is being tested. This process will be interrupted after tests have completed. Overrides: $HUNIT_EXEC_COMMAND.")
		fExecLog         = cmdline.String("exec:log", os.Getenv("HUNIT_EXEC_LOG"), "The path to log command output to. If omitted, output is redirected to standard output. Overrides: $HUNIT_EXEC_LOG.")
		fDebug           = cmdline.Bool("debug", strToBool(os.Getenv("HUNIT_DEBUG")), "Enable debugging mode. Overrides: $HUNIT_DEBUG.")
		fColor           = cmdline.Bool("color", strToBool(coalesce(os.Getenv("HUNIT_COLOR_OUTPUT"), "true")), "Colorize output when it's to a terminal. Overrides: $HUNIT_COLOR_OUTPUT.")
		fVerbose         = cmdline.Bool("verbose", strToBool(os.Getenv("HUNIT_VERBOSE")), "Be more verbose. Overrides: $HUNIT_VERBOSE and $VERBOSE.")
		fVersion         = cmdline.Bool("version", false, "Display the version and exit.")
	)
	cmdline.Var(&headerSpecs, "header", "Define a header to be set for every request, specified as 'Header-Name: <value>'. Provide -header repeatedly to set many headers.")
	cmdline.Var(&serviceSpecs, "service", "Define a mock service, specified as '[host]:<port>=endpoints.yml'. The service is available while tests are running.")
	cmdline.Parse(os.Args[1:])

	if *fVersion {
		fmt.Println(formatVersion())
		return 0
	}

	debug.DEBUG = debug.DEBUG || *fDebug
	debug.VERBOSE = debug.VERBOSE || *fVerbose
	color.NoColor = !*fColor

	var options test.Options
	if *fTrimEntity {
		options |= test.OptionEntityTrimTrailingWhitespace
	}
	if *fExpandVars {
		options |= test.OptionInterpolateVariables
	}
	if *fDumpRequest {
		options |= test.OptionDisplayRequests
	}
	if *fDumpResponse {
		options |= test.OptionDisplayResponses
	}
	if debug.VERBOSE {
		options |= test.OptionDisplayRequests | test.OptionDisplayResponses
	}

	options |= test.OptionDisplayRequestsOnFailure
	options |= test.OptionDisplayResponsesOnFailure

	var config test.Config
	if *fDocInclHTTP {
		config.Doc.IncludeHTTP = true
	}
	if *fDocFormatEntity {
		config.Doc.FormatEntities = true
	}
	if *fIOGracePeriod > 0 {
		config.Net.StreamIOGracePeriod = *fIOGracePeriod
	}

	var globalHeaders map[string]string
	if headerSpecs != nil && len(headerSpecs) > 0 {
		globalHeaders = make(map[string]string)
		for _, e := range headerSpecs {
			x := strings.Index(e, ":")
			if x < 1 {
				color.New(colorErr...).Printf("* * * Invalid header: %v\n", e)
				return 1
			}
			globalHeaders[strings.TrimSpace(e[:x])] = strings.TrimSpace(e[x+1:])
		}
	}

	var doctype doc_emit.Doctype
	var docname map[string]int
	if *fGendoc {
		var err error
		doctype, err = doc_emit.ParseDoctype(*fDoctype)
		if err != nil {
			color.New(colorErr...).Printf("* * * Invalid documentation type: %v\n", err)
			return 1
		}
		err = os.MkdirAll(*fDocpath, 0755)
		if err != nil {
			color.New(colorErr...).Printf("* * * Could not create documentation base: %v\n", err)
			return 1
		}
		docname = make(map[string]int)
	}

	var reports []report.Generator
	if *fReport {
		err := os.MkdirAll(*fReportPath, 0755)
		if err != nil {
			color.New(colorErr...).Printf("* * * Could not create documentation base: %v\n", err)
			return 1
		}

		rtype, err := report_emit.ParseDoctype(*fReportType)
		if err != nil {
			color.New(colorErr...).Printf("* * * Invalid report type: %v\n", err)
			return 1
		}

		out, err := os.OpenFile(path.Join(*fReportPath, rtype.String()+rtype.Ext()), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			color.New(colorErr...).Printf("* * * Could not open report output: %v\n", err)
			return 1
		}

		gen, err := report.New(rtype, out, fmt.Sprint(time.Now().Unix()))
		if err != nil {
			color.New(colorErr...).Printf("* * * Could create report generator: %v\n", err)
			return 1
		}

		err = gen.Init()
		if err != nil {
			color.New(colorErr...).Printf("* * * Could initialize report generator: %v\n", err)
			return 1
		}

		reports = []report.Generator{gen} // just one for now
	}

	services := 0
	for _, e := range serviceSpecs {
		conf, err := service.ParseConfig(e)
		if err != nil {
			color.New(colorErr...).Printf("* * * Could not create mock service: %v\n", err)
			return 1
		}
		svc, err := rest.New(conf) // only REST is supported for now...
		if err != nil {
			color.New(colorErr...).Printf("* * * Could not create mock service: %v\n", err)
			return 1
		}
		err = svc.Start()
		if err != nil {
			color.New(colorErr...).Printf("* * * Could not start mock service: %v\n", err)
			return 1
		}
		if debug.VERBOSE {
			fmt.Println()
		}
		defer func(s service.Service, c service.Config) {
			c.Resource.Close()
			s.Stop()
		}(svc, conf)
		fmt.Printf("----> Service %v (%v)\n", conf.Addr, conf.Path)
		services++
	}

	var done <-chan struct{}
	if *fExec != "" {
		var proc *exec.Process
		var err error
		proc, done, err = execCommandAsync(exec.NewCommand(*fExec, *fExec), *fExecLog)
		if err != nil {
			color.New(colorErr...).Printf("* * * %v\n", err)
			return 1
		}
		defer proc.Kill()
	}

	// give services and processes a second to settle
	if services > 0 {
		<-time.After(time.Second / 4)
	}
	if *fWait > 0 {
		fmt.Println("----> Waiting", *fWait)
		<-time.After(*fWait)
	}

	// setup caching
	var rcache, wcache *cache.Cache
	var cachePath string
	if *fCache && *fExec != "" {
		var err error
		sum, err := cache.Checksum(*fExec)
		if err != nil {
			color.New(colorErr...).Println("* * *", err)
			return 1
		}

		cachePath = path.Join(cacheBase, *fExec+".cache")
		wcache = &cache.Cache{
			Version: formatVersion(),
			Created: time.Now(),
			Binary:  sum,
		}

		c, err := cache.Read(cachePath)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("----> No results cache available:", cachePath)
		} else if err != nil {
			color.New(colorErr...).Println("* * *", err)
			return 1
		} else if b := c.Binary; b != nil && b.Path == *fExec && b.Checksum == sum.Checksum {
			fmt.Println("----> Using results cache:", cachePath)
			rcache = c
		}
	}

	var proc *exec.Process
	success := true
	start := time.Now()

suites:
	for _, e := range cmdline.Args() {
		if proc != nil {
			if proc.Running() {
				if l := proc.Linger(); l > 0 {
					color.New(colorSuite...).Printf("----> Waiting %v for process to complete...\n", l)
				}
				proc.Kill()
			}
			proc = nil
		}

		var suite *test.Suite
		var base string
		var err error
		cdup := config // copy global configs and update them

		var reader io.Reader
		if e == stdinPath {
			base = "(stdin)"
			color.New(colorSuite...).Printf("====> %s", base)
			reader = os.Stdin
		} else {
			base = path.Base(e)
			color.New(colorSuite...).Printf("====> %s", base)
			f, err := os.Open(e)
			if err != nil {
				color.New(colorErr...).Println("\n* * * Could not load test suite:", err)
				errno++
				break
			}
			defer f.Close()
			reader = f
		}

		suite, err = test.LoadSuiteFromReader(&cdup, reader)
		if err != nil {
			color.New(colorErr...).Println("\n* * * Could not load test suite:", err)
			errno++
			break
		}

		if suite.Title != "" {
			color.New(colorSuite...).Printf(" (%v)", suite.Title)
		}

		var sum *cache.Resource
		if (rcache != nil || wcache != nil) && e != stdinPath {
			sum, err = cache.Checksum(e)
			if err != nil {
				color.New(colorErr...).Println("* * * Could not load suite checksum:", err)
				errno++
				break
			}
			color.New(colorSuite...).Printf(" (cache: %s)\n", sum.Checksum)
		} else {
			fmt.Println()
		}

		if rcache != nil && e != stdinPath {
			cached := rcache.Suite(sum.Checksum)
			if cached != nil {
				fmt.Println("----> Reporting cached results from:", rcache.Created)
				results := rcache.ResultsForSuite(sum)
				success = reportResults(options, true, results, &tests, &failures, &skipped) && success
				wcache.AddSuite(cached, results)
				continue
			}
		}

		var gendoc []doc.Generator
		if *fGendoc {
			base := disambigFile(base, doctype.Ext(), docname)
			out, err := os.OpenFile(path.Join(*fDocpath, base), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				color.New(colorErr...).Println("* * * Could not open documentation output:", err)
				return 1
			}

			gen, err := doc.New(doctype, out)
			if err != nil {
				color.New(colorErr...).Println("* * * Could create documentation generator:", err)
				return 1
			}

			gendoc = []doc.Generator{gen} // just one for now
		}

		if len(suite.Setup) > 0 {
			if execCommands(suite.Setup) != nil {
				continue suites
			}
		}

		if suite.Exec != nil {
			cmd := suite.Exec
			cmd.Environment = exec.Environ(cmd.Environment)
			proc, _, err = execCommandAsync(*cmd, *fExecLog) // ignore done on per-suite tests
			if err != nil {
				color.New(colorErr...).Printf("* * * %v\n", err)
				continue suites
			}
		}

		if deps := suite.Deps; deps != nil {
			var deadline string
			if deps.Timeout == 0 {
				deadline = "forever"
			} else {
				deadline = fmt.Sprint(deps.Timeout)
			}
			if l := len(deps.Resources); l > 0 {
				if l == 1 {
					color.New(colorSuite...).Printf("----> Waiting %s for one dependency...\n", deadline)
				} else {
					color.New(colorSuite...).Printf("----> Waiting %s for %d dependencies...\n", deadline, l)
				}
				err := await.Await(context.Background(), deps.Resources, deps.Timeout)
				if err != nil {
					color.New(colorErr...).Printf("* * * Error waiting for dependencies: %v\n", err)
					errno++
					continue
				}
			}
		}

		startSuite := time.Now()
		results, err := hunit.RunSuite(suite, hunit.Context{BaseURL: *fBaseURL, Options: options, Headers: globalHeaders, Debug: debug.DEBUG, Gendoc: gendoc, Config: cdup})
		if err != nil {
			color.New(colorErr...).Printf("* * * Could not run test suite: %v\n", err)
			errno++
		}
		suiteDuration := time.Since(startSuite)

		if (options & (test.OptionDisplayRequests | test.OptionDisplayResponses)) != 0 {
			if len(results) > 0 {
				fmt.Println()
			}
		}

		for _, e := range reports {
			err := e.Suite(cdup, suite, &report_emit.Results{results, suiteDuration})
			if err != nil {
				color.New(colorErr...).Printf("* * * Could not emit report: %v\n", err)
			}
		}

		for _, e := range gendoc {
			err := e.Close()
			if err != nil {
				color.New(colorErr...).Printf("* * * Could not finalize documentation writer: %v\n", err)
			}
		}

		success = reportResults(options, false, results, &tests, &failures, &skipped) && success
		if wcache != nil && sum != nil {
			wcache.AddSuite(sum, results)
		}

		if len(suite.Teardown) > 0 {
			if execCommands(suite.Teardown) != nil {
				continue suites
			}
		}
	}

	for _, e := range reports {
		err := e.Finalize()
		if err != nil {
			color.New(colorErr...).Printf("* * * Could not finalize report writer: %v\n", err)
		}
	}

	if tests < 1 && errno < 1 && services > 0 {
		if done != nil {
			color.New(colorSuite...).Println("====> No tests; running services until process exits...")
			<-done
		} else {
			color.New(colorSuite...).Println("====> No tests; running services until we're interrupted...")
			<-make(chan struct{})
		}
	}

	if proc != nil {
		if proc.Running() {
			if l := proc.Linger(); l > 0 {
				color.New(colorSuite...).Printf("----> Waiting %v for process to complete...\n", l)
			}
			proc.Kill()
		}
		proc = nil
	}

	duration := time.Since(start)
	fmt.Println()

	if wcache != nil {
		err := cache.Write(cachePath, wcache)
		if err != nil {
			color.New(colorErr...).Printf("* * * Could not write cache: %v\n\n", err)
		}
	}

	if errno > 0 {
		color.New(color.BgHiRed, color.Bold, color.FgBlack).Printf(" ERRORS! ")
		fmt.Printf(" %d %s could not be run due to errors.\n\n", errno, plural(errno, "test", "tests"))
		return 1
	}

	fmt.Printf("Finished in %v.\n\n", duration)

	if !success {
		color.New(color.BgHiRed, color.Bold, color.FgBlack).Printf(" FAIL! ")
		fmt.Printf(" %d of %d tests failed (%d implicit).\n", failures, tests, skipped)
		return 1
	}

	color.New(color.BgHiGreen, color.Bold, color.FgBlack).Printf(" PASS! ")
	if tests == 0 {
		fmt.Printf(" Hmm, nothing to do, really...\n")
	} else if tests == 1 {
		fmt.Printf(" The test passed.\n")
	} else {
		fmt.Printf(" All %d tests passed.\n", tests)
	}
	return 0
}

func reportResults(options test.Options, cached bool, results []*hunit.Result, tests, failures, skipped *int) bool {
	var count int
	var success bool
	var prefix string
	if cached {
		prefix = "(cached) "
	}
	for _, r := range results {
		*tests++
		if !r.Success {
			success = false
			*failures++
		}
		if r.Skipped {
			color.New(color.FgYellow).Printf("----> %s%v", prefix, r.Name)
			*skipped++
			continue
		}
		if r.Success {
			color.New(color.FgCyan).Printf("----> %s%v", prefix, r.Name)
		} else {
			color.New(color.FgRed).Printf("----> %s%v", prefix, r.Name)
		}
		if r.Errors != nil {
			for _, e := range r.Errors {
				count++
				fmt.Println(text.IndentWithOptions(fmt.Sprintf("        #%d %s", count, e), "             ", 0))
				fmt.Println()
			}
		}
		preq := len(r.Reqdata) > 0 && ((options&test.OptionDisplayRequests) == test.OptionDisplayRequests || (!r.Success && (options&test.OptionDisplayRequestsOnFailure) == test.OptionDisplayRequestsOnFailure))
		prsp := len(r.Rspdata) > 0 && ((options&test.OptionDisplayResponses) == test.OptionDisplayResponses || (!r.Success && (options&test.OptionDisplayResponsesOnFailure) == test.OptionDisplayResponsesOnFailure))
		if preq {
			fmt.Println(text.Indent(string(r.Reqdata), "      > "))
		}
		if preq && prsp {
			fmt.Println("      * ")
		}
		if prsp {
			fmt.Println(text.Indent(string(r.Rspdata), "      < "))
		}
		if preq || prsp {
			fmt.Println()
		}
	}
	return success
}

// Execute a set of commands in sequence, allowing each to terminate before
// the next is executed.
func execCommands(cmds []*exec.Command) error {
	for i, e := range cmds {
		if i > 0 && debug.VERBOSE {
			fmt.Println()
		}

		if e.Command == "" {
			color.New(colorErr...).Printf("* * * Setup command #%d is empty (did you set 'run'?)", i+1)
			return fmt.Errorf("Empty command")
		}

		if e.Display != "" {
			fmt.Printf("----> %v ", e.Display)
		} else {
			fmt.Printf("----> $ %v ", e.Command)
		}

		out, err := e.Exec()
		if err != nil {
			fmt.Println()
			color.New(colorErr...).Printf("* * * Setup command #%d failed: %v\n", i+1, err)
			if len(out) > 0 {
				fmt.Println(text.Indent(string(out), "      < "))
			}
			return err
		}

		color.New(color.Bold, color.FgHiGreen).Println("OK")
		if debug.VERBOSE && len(out) > 0 {
			fmt.Println(text.Indent(string(out), "      < "))
		}
	}
	return nil
}

// Execute a single command and do not wait for it to terminate
func execCommandAsync(cmd exec.Command, logs string) (*exec.Process, <-chan struct{}, error) {
	if cmd.Command == "" {
		return nil, nil, fmt.Errorf("Empty command (did you set 'run'?)")
	}

	var wout, werr io.WriteCloser
	if logs != "" {
		var err error
		out, err := os.OpenFile(logs, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("Could not open exec log: %v", err)
		}
		wout, werr = out, out
	} else {
		wout = exec.NewPrefixWriter(syncStdout, "      ◇ ")
		werr = exec.NewPrefixWriter(syncStdout, color.New(color.FgRed).Sprint("      ◆ "))
	}

	proc, err := cmd.Start(wout, werr)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not exec process: %v", err)
	}

	color.New(colorSuite...).Printf("----> $ %v\n", proc)

	done := make(chan struct{})
	go func() {
		state := proc.Monitor()
		color.New(colorSuite...).Printf("----> * %v (pid %d; %s)\n", proc, state.Pid(), state)
		close(done)
	}()

	if cmd.Wait > 0 {
		color.New(colorSuite...).Printf("----> Waiting %v for process to settle...\n", cmd.Wait)
		<-time.After(cmd.Wait)
	}
	return proc, done, nil
}

func formatVersion() string {
	if version == githash {
		return version
	} else {
		return fmt.Sprintf("%s (%s)", version, githash)
	}
}
