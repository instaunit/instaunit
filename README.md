# Instaunit tests your Web APIs

Instaunit makes writing tests for REST and Websocket services simple and declarative. You can use Instaunit locally for development and on your CI infrastructure to run automated tests against your Web services.

Since tests and documentation are naturally maintained in parallel, Instaunit can also combine these two highly-related tasks into one: add optional usage descriptions to your tests and Instaunit can generate documentation for your endpoints, complete with examples.

## Describing Tests

Tests are described by a simple YAML-based document format. Just describe your request, the response you expect, and that's basically it. Here's a very simple test suite containing a single test case:

```yaml
tests:
  -
    request: # The request to perform
      method: GET
      url: https://www.example.com/status
    response: # The response we expect to get
      status: 200
      entity: {"status": "Ok"}
```

Instaunit supports a lot of other request properties – from headers to authorization to parameter encoding – as well as a bunch of advanced features if you need them:

* **Compare entities semantically** to ignore insignificant differences like whitespace and map key order,
* Reference the output of previously-run tests to **chain related tests together**,
* **Wait for dependency services** to become available before a test suite starts running,
* Declare **mock services** so your tests avoid external dependencies without making code-level changes,
* Evaluate **expressions and built-in functions** to generate input and randomize your tests.

## Documenting Tests

Tests and documentation are naturally maintained together. To generate documentation, simply add a description to a representative test case for your endpoint. You can pick and choose which tests generate documentation.

Instaunit supports a variety of properties you can include to document a request but all you need to get started is `doc`:

```yaml
tests:
  -
    doc: |
      About this endpoint. Use _Markdown_ to add formatting!
    request:
      method: GET
      url: https://www.example.com/status
    response:
      status: 200
      entity: {"status": "Ok"}
```

## Running Tests

Tests can be run by pointing `instaunit` to a test suite document (or many of them). 

Try running `instaunit -h` to view the options it supports.

```
$ instaunit test.yml
  ====> test.yml
  ----> GET https://www.example.com/status
        #1 Unexpected status code:
             expected: (int) 200
               actual: (int) 500

        #2 Entities do not match:
             --- Expected
             +++ Actual
             @@ -1,2 +1,2 @@
             -{"status": "Ok"}
             +{"status": "Not great"}
```
