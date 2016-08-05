package main

import (
  "os"
  "fmt"
  "flag"
  "path"
  "hunit"
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
  fTrimEntity   := cmdline.Bool     ("entity:trim",     true,                                                         "Trim trailing whitespace from entities.")
  fDebug        := cmdline.Bool     ("debug",           false,                                                        "Enable debugging mode.")
  fVerbose      := cmdline.Bool     ("verbose",         false,                                                        "Enable verbose debugging mode.")
  cmdline.Parse(os.Args[1:])
  
  DEBUG = *fDebug
  DEBUG_VERBOSE = *fVerbose
  
  var options hunit.Options
  if *fTrimEntity {
    options |= hunit.OptionEntityTrimTrailingWhitespace
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
