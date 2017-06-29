package main

import (
  "os"
  "io"
  "fmt"
  "flag"
  "path"
  "time"
  "strings"
  "hunit"
  "hunit/test"
  "hunit/text"
  "hunit/doc"
  "hunit/doc/emit"
  "hunit/service"
  "hunit/service/backend/rest"
)

import (
  "github.com/fatih/color"
  "github.com/bww/go-util/debug"
)

var (
  colorErr    = []color.Attribute{color.FgYellow}
  colorSuite  = []color.Attribute{color.Bold}
  colorDetail = []color.Attribute{color.Faint}
)

/**
 * You know what it does
 */
func main() {
  os.Exit(app())
}

/**
 * You know what it does
 */
func app() int {
  var tests, failures, errors int
  var headerSpecs, serviceSpecs flagList
  
  cmdline           := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
  fBaseURL          := cmdline.String   ("base-url",              coalesce(os.Getenv("HUNIT_BASE_URL"), "http://localhost/"),   "The base URL for requests.")
  fExpandVars       := cmdline.Bool     ("expand",                strToBool(os.Getenv("HUNIT_EXPAND_VARS"), true),              "Expand variables in test cases.")
  fTrimEntity       := cmdline.Bool     ("entity:trim",           strToBool(os.Getenv("HUNIT_TRIM_ENTITY"), true),              "Trim trailing whitespace from entities.")
  fDumpRequest      := cmdline.Bool     ("dump:request",          strToBool(os.Getenv("HUNIT_DUMP_REQUESTS")),                  "Dump requests to standard output as they are processed.")
  fDumpResponse     := cmdline.Bool     ("dump:response",         strToBool(os.Getenv("HUNIT_DUMP_RESPONSES")),                 "Dump responses to standard output as they are processed.")
  fGendoc           := cmdline.Bool     ("gendoc",                strToBool(os.Getenv("HUNIT_GENDOC")),                         "Generate documentation.")
  fDocpath          := cmdline.String   ("doc:output",            coalesce(os.Getenv("HUNIT_DOC_OUTPUT"), "./docs"),            "The directory in which generated documentation should be written.")
  fDoctype          := cmdline.String   ("doc:type",              coalesce(os.Getenv("HUNIT_DOC_TYPE"), "markdown"),            "The format to generate documentation in.")
  fDocInclHTTP      := cmdline.Bool     ("doc:include-http",      strToBool(os.Getenv("HUNIT_DOC_INCLUDE_HTTP")),               "Include HTTP in request and response examples (as opposed to just routes and entities).")
  fDocFormatEntity  := cmdline.Bool     ("doc:format-entities",   strToBool(os.Getenv("HUNIT_DOC_FORMAT_ENTITIES")),            "Pretty-print supported request and response entities in documentation output.")
  fDebug            := cmdline.Bool     ("debug",                 strToBool(os.Getenv("HUNIT_DEBUG")),                          "Enable debugging mode.")
  fColor            := cmdline.Bool     ("color",                 strToBool(coalesce(os.Getenv("HUNIT_COLOR_OUTPUT"), "true")), "Colorize output when it's to a terminal.")
  fVerbose          := cmdline.Bool     ("verbose",               strToBool(os.Getenv("HUNIT_VERBOSE")),                        "Be more verbose.")
  cmdline.Var        (&headerSpecs,      "header",                                                                              "Define a header to be set for every request, specified as 'Header-Name: <value>'. Provide -header repeatedly to set many headers.")
  cmdline.Var        (&serviceSpecs,     "service",                                                                             "Define a mock service, specified as '[host]:<port>=endpoints.yml'. The service is available while tests are running.")
  cmdline.Parse(os.Args[1:])
  
  debug.DEBUG = *fDebug
  debug.VERBOSE = *fVerbose
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
  
  var doctype emit.Doctype
  var docname map[string]int
  if *fGendoc {
    var err error
    doctype, err = emit.ParseDoctype(*fDoctype)
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
  
  services := 0
  for _, e := range serviceSpecs {
    conf, err := service.ParseConfig(e)
    if err != nil {
      color.New(colorErr...).Printf("* * * Could not create mock service: %v\n", err)
      return 1
    }
    svc, err := rest.New(conf)
    if err != nil {
      color.New(colorErr...).Printf("* * * Could not create mock service: %v\n", err)
      return 1
    }
    err = svc.StartService()
    if err != nil {
      color.New(colorErr...).Printf("* * * Could not start mock service: %v\n", err)
      return 1
    }
    defer func(s service.Service, c service.Config){
      c.Resource.Close()
      s.StopService()
    }(svc, conf)
    fmt.Printf("----> Service %v (%v)\n", conf.Addr, conf.Path)
    services++
  }
  
  // give service a second to settle
  if services > 0 {
    <-time.After(time.Second / 2)
  }
  
  success := true
  start := time.Now()
  for _, e := range cmdline.Args() {
    base := path.Base(e)
    color.New(colorSuite...).Printf("====> %v", base)
    
    cdup := config // copy global configs and update them
    suite, err := test.LoadSuiteFromFile(&cdup, e)
    if err != nil {
      fmt.Println()
      color.New(colorErr...).Printf("* * * Could not load test suite: %v\n", err)
      errors++
      continue
    }
    
    if suite.Title != "" {
      color.New(colorSuite...).Printf(" (%v)\n", suite.Title)
    }else{
      fmt.Println()
    }
    
    var out io.WriteCloser
    var gendoc []doc.Generator
    if *fGendoc {
      var err error
      ext := path.Ext(base)
      stem := base[:len(base) - len(ext)]
      
      n, ok := docname[stem]
      if ok && n > 0 {
        stem = fmt.Sprintf("%v-%d", stem, n)
      }
      docname[stem] = n + 1
      
      out, err = os.OpenFile(path.Join(*fDocpath, stem + doctype.Ext()), os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
      if err != nil {
        color.New(colorErr...).Printf("* * * Could not open documentation output: %v\n", err)
        return 1
      }
      
      gen, err := doc.New(doctype, out)
      if err != nil {
        color.New(colorErr...).Printf("* * * Could create documentation generator: %v\n", err)
        return 1
      }
      
      gendoc = []doc.Generator{gen} // just one for now
    }
    
    results, err := hunit.RunSuite(suite, hunit.Context{BaseURL: *fBaseURL, Options: options, Headers: globalHeaders, Debug: debug.DEBUG, Gendoc: gendoc, Config: cdup})
    if err != nil {
      color.New(colorErr...).Printf("* * * Could not run test suite: %v\n", err)
      errors++
    }
    
    if (options & (test.OptionDisplayRequests | test.OptionDisplayResponses)) != 0 {
      if len(results) > 0 {
        fmt.Println()
      }
    }
    
    if out != nil {
      err := out.Close()
      if err != nil {
        color.New(colorErr...).Printf("* * * Could not close documentation writer: %v\n", err)
      }
    }
    
    var count int
    for _, r := range results {
      if r.Success {
        color.New(color.FgCyan).Printf("----> %v", r.Name)
      }else{
        color.New(color.FgRed).Printf("----> %v", r.Name)
      }
      tests++
      if !r.Success {
        success = false
        failures++
      }
      if r.Errors != nil {
        for _, e := range r.Errors {
          count++
          fmt.Println(text.IndentWithOptions(fmt.Sprintf("        #%d %v", count, e), "             ", 0))
          fmt.Println()
        }
      }
      preq := (options & test.OptionDisplayRequests) == test.OptionDisplayRequests || (!r.Success && (options & test.OptionDisplayRequestsOnFailure) == test.OptionDisplayRequestsOnFailure)
      prsp := (options & test.OptionDisplayResponses) == test.OptionDisplayResponses || (!r.Success && (options & test.OptionDisplayResponsesOnFailure) == test.OptionDisplayResponsesOnFailure)
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
        fmt.Println("\n")
      }
    }
    
  }
  
  if tests < 1 && services > 0 {
    color.New(colorSuite...).Println("====> No tests; running services until we're interrupted...")
    <- make(chan struct{})
  }
  
  duration := time.Since(start)
  fmt.Println()
  
  if errors > 0 {
    color.New(color.BgHiRed, color.Bold, color.FgBlack).Printf(" ERRORS! ")
    fmt.Printf(" %d %s could not be run due to errors.\n", errors, plural(errors, "test", "tests"))
    return 1
  }
  
  fmt.Printf("Finished in %v.\n\n", duration)
  
  if !success {
    color.New(color.BgHiRed, color.Bold, color.FgBlack).Printf(" FAIL! ")
    fmt.Printf(" %d of %d tests failed.\n", failures, tests)
    return 1
  }
  
  color.New(color.BgHiGreen, color.Bold, color.FgBlack).Printf(" PASS! ")
  if tests == 0 {
    fmt.Printf(" Hmm, nothing to do, really...\n")
  }else if tests == 1 {
    fmt.Printf(" The test passed.\n")
  }else{
    fmt.Printf(" All %d tests passed.\n", tests)
  }
  return 0
}

/**
 * Flag string list
 */
type flagList []string

/**
 * Set a flag
 */
func (s *flagList) Set(v string) error {
  *s = append(*s, v)
  return nil
}

/**
 * Describe
 */
func (s *flagList) String() string {
  return fmt.Sprintf("%+v", *s)
}

/**
 * Pluralize
 */
func plural(v int, s, p string) string {
  if v == 1 {
    return s
  }else{
    return p
  }
}

/**
 * Return the first non-empty string from those provided
 */
func coalesce(v... string) string {
  for _, e := range v {
    if e != "" {
      return e
    }
  }
  return ""
}

/**
 * String to bool
 */
func strToBool(s string, d ...bool) bool {
  if s == "" {
    if len(d) > 0 {
      return d[0]
    }else{
      return false
    }
  }
  return strings.EqualFold(s, "t") || strings.EqualFold(s, "true") || strings.EqualFold(s, "y") || strings.EqualFold(s, "yes")
}
