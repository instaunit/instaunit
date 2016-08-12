package main

import (
  "os"
  "fmt"
  "flag"
  "path"
  "hunit"
  "strings"
)

var DEBUG bool
var DEBUG_VERBOSE bool

/**
 * You know what it does
 */
func main() {
  var tests, failures int
  
  cmdline       := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
  fBaseURL      := cmdline.String   ("base-url",        coalesce(os.Getenv("HUNIT_BASE_URL"), "http://localhost/"),   "The base URL for requests.")
  fExpandVars   := cmdline.Bool     ("expand",          strToBool(os.Getenv("HUNIT_EXPAND_VARS"), true),              "Expand variables in test cases.")
  fTrimEntity   := cmdline.Bool     ("entity:trim",     strToBool(os.Getenv("HUNIT_TRIM_ENTITY"), true),              "Trim trailing whitespace from entities.")
  fDumpRequest  := cmdline.Bool     ("dump:request",    strToBool(os.Getenv("HUNIT_DUMP_REQUESTS")),                  "Dump requests to standard output as they are processed.")
  fDumpResponse := cmdline.Bool     ("dump:response",   strToBool(os.Getenv("HUNIT_DUMP_RESPONSES")),                 "Dump responses to standard output as they are processed.")
  fDebug        := cmdline.Bool     ("debug",           strToBool(os.Getenv("HUNIT_DEBUG")),                          "Enable debugging mode.")
  fVerbose      := cmdline.Bool     ("verbose",         strToBool(os.Getenv("HUNIT_VERBOSE")),                        "Be more verbose.")
  cmdline.Parse(os.Args[1:])
  
  DEBUG = *fDebug
  DEBUG_VERBOSE = *fVerbose
  
  var options hunit.Options
  if *fTrimEntity {
    if DEBUG_VERBOSE { fmt.Println("Enabled: Trim entity trailing whitespace") }
    options |= hunit.OptionEntityTrimTrailingWhitespace
  }
  if *fExpandVars {
    if DEBUG_VERBOSE { fmt.Println("Enabled: Expand variables in test cases") }
    options |= hunit.OptionInterpolateVariables
  }
  if *fDumpRequest {
    if DEBUG_VERBOSE { fmt.Println("Enabled: Dump requests") }
    options |= hunit.OptionDisplayRequests
  }
  if *fDumpResponse {
    if DEBUG_VERBOSE { fmt.Println("Enabled: Dump responses") }
    options |= hunit.OptionDisplayResponses
  }
  if DEBUG && DEBUG_VERBOSE {
    fmt.Println("Enabled: Verbose output")
    options |= hunit.OptionDisplayRequests | hunit.OptionDisplayResponses
  }
  
  success := true
  for _, e := range cmdline.Args() {
    base := path.Base(e)
    if DEBUG {
      fmt.Printf("====> %v\n", base)
    }
    
    suite, err := hunit.LoadSuiteFromFile(e)
    if err != nil {
      fmt.Printf("Could not load test suite: %v\n", err)
      continue
    }
    
    results, err := suite.Run(hunit.Context{BaseURL: *fBaseURL, Options: options, Debug: DEBUG})
    if err != nil {
      fmt.Printf("Could not run test suite: %v\n", err)
      continue
    }
    
    if (options & (hunit.OptionDisplayRequests | hunit.OptionDisplayResponses)) != 0 {
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
          fmt.Printf("      #%d %v\n", count, e)
        }
      }
    }
    
  }
  
  fmt.Println()
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
