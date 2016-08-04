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
  
  cmdline     := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
  fBaseURL    := cmdline.String   ("base-url",    coalesce(os.Getenv("HUNIT_BASE_URL"), "http://localhost/"),   "The base URL for requests.")
  fDebug      := cmdline.Bool     ("debug",       false,                                                        "Enable debugging mode.")
  fVerbose    := cmdline.Bool     ("verbose",     false,                                                        "Enable verbose debugging mode.")
  cmdline.Parse(os.Args[1:])
  
  DEBUG = *fDebug
  DEBUG_VERBOSE = *fVerbose
  
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
    
    results, err := suite.Run(hunit.Context{BaseURL: *fBaseURL, Debug: DEBUG})
    if err != nil {
      fmt.Printf("Could not run test suite: %v\n", err)
      continue
    }
    
    for _, r := range results {
      if !r.Success {
        success = false
      }
      if r.Errors != nil {
        for _, e := range r.Errors {
          fmt.Printf("%v: %v\n", base, e)
        }
      }
    }
    
  }
  
  if !success {
    os.Exit(1)
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
