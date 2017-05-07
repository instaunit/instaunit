package env

import (
  "time"
  "net/http"
  "io/ioutil"
)

/**
 * Shared client
 */
var httpClient = &http.Client{Timeout:time.Second * 5}

/**
 * Fetch AWS metadata
 */
func awsGet(u string) (string, error) {
  
  rsp, err := httpClient.Get(u)
  if err != nil {
    return "", err
  }
  if rsp.Body != nil {
    defer rsp.Body.Close()
  }
  
  data, err := ioutil.ReadAll(rsp.Body)
  if err != nil{
    return "", err
  }
  
  return string(data), nil
}

/**
 * Determine our hostname
 */
func awsLocalHostname() (string, error) {
  return awsGet("http://169.254.169.254/latest/meta-data/local-hostname")
}

/**
 * Determine our local IP
 */
func awsLocalIPv4() (string, error) {
  return awsGet("http://169.254.169.254/latest/meta-data/local-ipv4")
}

/**
 * Determine our public IP
 */
func awsPublicIPv4() (string, error) {
  return awsGet("http://169.254.169.254/latest/meta-data/public-ipv4")
}
