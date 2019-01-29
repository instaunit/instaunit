# Instaunit tests your service APIs

Instaunit makes building tests for REST and Websocket services simple and declarative. And since tests and documentation are naturally maintained in parallel, Instaunit can combine these two highly-related tasks into one: descriptions can be optionally added to your tests to automatically generate documentation of your endpoints.

## Describing Tests

Tests are described by a simple YAML-based document format. Just describe your request, the response you expect, and that's basically it. Here's a very simple test suite:

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

Instaunit also supports a bunch of advanced features if you need them:

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

## Generating Documentation

Documentation can be generated for the endpoints in your tests by passing the `-gendoc` flag. When set, tests are run and their input and output are automatically incorporated into documentation as examples.

Currently, documentation can be generated as [Markdown](https://en.wikipedia.org/wiki/Markdown) (which, of course, can easily be converted to HTML). We'd be [interested to hear](https://github.com/instaunit/instaunit/issues) about other documentation formats that may be worth supporting.

Refer to [`example/docs`](https://github.com/instaunit/instaunit/blob/master/example/docs) for an example of the docs generated for the example tests.

```
$ instaunit -gendoc test.yml
====> test.yml
# ...
```
