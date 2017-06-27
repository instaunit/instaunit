# Instaunit

Instaunit is a command-line tool that runs tests against HTTP APIs. It makes managing your tests simple and declarative without skimping on the features.

Since tests and documentation are naturally maintained in parallel, Instaunit combines these two highly-related tasks into one. Descriptions can be added to your tests to automatically generate documentation of your endpoints.

## Describing Tests

Tests are described by a simple YAML-based document format. Tests and documentation are managed together.

Just describe your request, the response you expect, and you're done! Instaunit also supports plenty of advanced features if you need them:

* Compare entities *semantically* to ignore insignificant differences like whitespace and map key order,
* Reference the output of previously-run tests from subsequent, dependent tests,
* Mock services can produce test output to be consumed by your tests via HTTP,
* Handy functions can generate test input to randomize your tests.

Refer to the full [`test.yml`](https://github.com/instaunit/instaunit/blob/master/example/test.yml) file for a more complete illustration of test cases.

```yaml
- 
    doc: |
      Fetch a document from our [Github repo](github.com/instaunit/instaunit).
    
    request:
      method: GET
      url: https://raw.githubusercontent.com/instaunit/instaunit/master/example/test.txt
      headers:
        Origin: localhost
    
    response:
      status: 200
      headers:
        Content-Type: text/plain; charset=utf-8
      entity: |
        Heres a simple
        response from the
        server.      
```

## Running Tests

Tests can be run by pointing `instaunit` to a test suite document (or many of them). 

Try running `instaunit -h` to view the options it supports.

```
$ instaunit test.yml
  ====> test.yml
  ----> GET https://raw.githubusercontent.com/instaunit/instaunit/master/example/test.txt
        #1 Unexpected status code:
             expected: (int) 404
               actual: (int) 200

        #2 Entities do not match:
             --- Expected
             +++ Actual
             @@ -1,2 +1,2 @@
             -Heres a simple
             +Here's a simple
              response from the
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
