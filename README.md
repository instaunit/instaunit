# Instaunit

### Instaunit runs your tests against HTTP and Websocket APIs.

It makes managing your tests simple and declarative. And since tests and documentation are naturally maintained in parallel, Instaunit combines these two highly-related tasks into one: descriptions can be added to your tests to automatically generate documentation of your endpoints.

## Describing Tests

Tests are described by a simple YAML-based document format. Just describe your request, the response you expect, and that's basically it. Instaunit also supports a bunch of advanced features if you need them:

* **Compare entities semantically** to ignore insignificant differences like whitespace and map key order,
* Reference the output of previously-run tests to **chain related tests together**,
* Easily build **mock services** so your tests don't need any external dependencies,
* Use **convenient, built-in functions** to generate input and randomize your tests.

Tests and documentation are maintained together. To generate documentation, simply add descriptions to a representative test case for your endpoint. You can pick and choose which tests generate documentation.

Refer to the full [`test.yml`](https://github.com/instaunit/instaunit/blob/master/example/test.yml) file for a more complete illustration of test cases.

```yaml
tests:
  -
    doc: A description of this endpoint.
  	
    request:
      method: GET
      url: https://raw.githubusercontent.com/instaunit/instaunit/master/example/test.txt
    
    response:
      status: 200
      entity: Here's a simple response from the server.
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
             -Here's a simple response from the server.
             +Heres a simple response from the server.
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
