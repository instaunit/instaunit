package main

import (
  "os"
  "fmt"
  "flag"
  "path"
  "strings"
  "hunit"
  "hunit/test"
  "hunit/text"
  "hunit/doc"
  "hunit/doc/emit"
)

var DEBUG bool
var VERBOSE bool

const ws = " \n\r\t\v"

/**
 * You know what it does
 */
func main() {
  var tests, failures, errors int
  var headerSpecs flagList
  
  cmdline       := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
  fBaseURL      := cmdline.String   ("base-url",          coalesce(os.Getenv("HUNIT_BASE_URL"), "http://localhost/"),   "The base URL for requests.")
  fExpandVars   := cmdline.Bool     ("expand",            strToBool(os.Getenv("HUNIT_EXPAND_VARS"), true),              "Expand variables in test cases.")
  fTrimEntity   := cmdline.Bool     ("entity:trim",       strToBool(os.Getenv("HUNIT_TRIM_ENTITY"), true),              "Trim trailing whitespace from entities.")
  fDumpRequest  := cmdline.Bool     ("dump:request",      strToBool(os.Getenv("HUNIT_DUMP_REQUESTS")),                  "Dump requests to standard output as they are processed.")
  fDumpResponse := cmdline.Bool     ("dump:response",     strToBool(os.Getenv("HUNIT_DUMP_RESPONSES")),                 "Dump responses to standard output as they are processed.")
  fGendoc       := cmdline.Bool     ("gendoc",            strToBool(os.Getenv("HUNIT_GENDOC")),                         "Generate documentation.")
  fDocpath      := cmdline.String   ("doc:output",        coalesce(os.Getenv("HUNIT_DOC_OUTPUT"), "./docs"),            "The directory in which generated documentation should be written.")
  fDoctype      := cmdline.String   ("doc:type",          coalesce(os.Getenv("HUNIT_DOC_TYPE"), "markdown"),            "The format to generate documentation in.")
  fDocInclHTTP  := cmdline.Bool     ("doc:include-http",  strToBool(os.Getenv("HUNIT_DOC_INCLUDE_HTTP")),               "Include HTTP in request and response examples (as opposed to just routes and entities).")
  fDebug        := cmdline.Bool     ("debug",             strToBool(os.Getenv("HUNIT_DEBUG")),                          "Enable debugging mode.")
  fVerbose      := cmdline.Bool     ("verbose",           strToBool(os.Getenv("HUNIT_VERBOSE")),                        "Be more verbose.")
  cmdline.Var    (&headerSpecs,      "header",                                                                          "Define a header to be set for every request, specified as 'Header-Name: <value>'. Provide -header repeatedly to set many headers.")
  cmdline.Parse(os.Args[1:])
  
  DEBUG = *fDebug
  VERBOSE = *fVerbose
  
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
  if VERBOSE {
    options |= test.OptionDisplayRequests | test.OptionDisplayResponses
  }
  
  var config test.Config
  if *fDocInclHTTP {
    config.Doc.IncludeHTTP = true
  }
  
  var globalHeaders map[string]string
  if headerSpecs != nil && len(headerSpecs) > 0 {
    globalHeaders = make(map[string]string)
    for _, e := range headerSpecs {
      x := strings.Index(e, ":")
      if x < 1 {
        fmt.Printf("* * * Invalid header: %v\n", e)
        return
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
      fmt.Printf("* * * Invalid documentation type: %v\n", err)
      return
    }
    err = os.MkdirAll(*fDocpath, 0755)
    if err != nil {
      fmt.Printf("* * * Could not create documentation base: %v\n", err)
      return
    }
    docname = make(map[string]int)
  }
  
  success := true
  for _, e := range cmdline.Args() {
    base := path.Base(e)
    fmt.Printf("====> %v", base)
    
    cdup := config // copy global configs and update them
    suite, err := test.LoadSuiteFromFile(&cdup, e)
    if err != nil {
      fmt.Println()
      fmt.Printf("* * * Could not load test suite: %v\n", err)
      errors++
      continue
    }
    
    if suite.Title != "" {
      fmt.Printf(" (%v)\n", suite.Title)
    }else{
      fmt.Println()
    }
    
    var gendoc []doc.Generator
    if *fGendoc {
      ext := path.Ext(base)
      stem := base[:len(base) - len(ext)]
      
      n, ok := docname[stem]
      if ok && n > 0 {
        stem = fmt.Sprintf("%v-%d", stem, n)
      }
      docname[stem] = n + 1
      
      out, err := os.OpenFile(path.Join(*fDocpath, stem + doctype.Ext()), os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
      if err != nil {
        fmt.Printf("* * * Could not open documentation output: %v\n", err)
        return
      }
      
      gen, err := doc.New(doctype, out)
      if err != nil {
        fmt.Printf("* * * Could create documentation generator: %v\n", err)
        return
      }
      
      gendoc = []doc.Generator{gen} // just one for now
    }
    
    results, err := hunit.RunSuite(suite, hunit.Context{BaseURL: *fBaseURL, Options: options, Headers: globalHeaders, Debug: DEBUG, Gendoc: gendoc, Config: cdup})
    if err != nil {
      fmt.Printf("* * * Could not run test suite: %v\n", err)
      errors++
    }
    
    if (options & (test.OptionDisplayRequests | test.OptionDisplayResponses)) != 0 {
      if len(results) > 0 {
        fmt.Println()
      }
    }
    
    var count int
    for _, r := range results {
      fmt.Printf("----> %v", r.Name)
      tests++
      if !r.Success {
        success = false
        failures++
      }
      if r.Errors != nil {
        for _, e := range r.Errors {
          count++
          fmt.Println(text.IndentWithOptions(fmt.Sprintf("      #%d %v\n", count, e), "           ", 0))
        }
      }
    }
    
  }
  
  fmt.Println()
  if errors > 0 {
    fmt.Printf("ERRORS! %d %s could not be run due to errors.\n", errors, plural(errors, "test", "tests"))
    os.Exit(1)
  }
  if !success {
    fmt.Printf("FAILURES! %d of %d tests failed.\n", failures, tests)
    os.Exit(1)
  }
  if tests == 1 {
    fmt.Printf("SUCCESS! The test passed.\n")
  }else{
    fmt.Printf("SUCCESS! All %d tests passed.\n", tests)
  }
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
