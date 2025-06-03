# Instaunit tests your HTTP APIs

Instaunit is a tool that lets you write integration tests for REST and Websocket services declaratively and run them automatically. You can use Instaunit locally for development and on your CI infrastructure as part of your integration tests.

**Instaunit manages the repetitive details** of setting up and tearing down, executing requests and evaluating responses, which allows you to get more value out of your tests faster.

Since tests and documentation are naturally maintained in parallel, Instaunit can combine these two highly-related tasks into one: add optional usage descriptions to your tests and Instaunit can generate documentation for your endpoints, complete with examples.

# Getting Started

Get up and running quickly with our [**Getting Started Tutorial**](https://github.com/instaunit/instaunit/wiki/Getting-Started).

## Installing Instaunit

MacOS users can install the latest version of Instaunit via Homebrew by running the following command.

```
brew install bww/stable/instaunit
```

Alternativelly, you can install Instaunit by:

* [Downloading a binary release](https://github.com/instaunit/instaunit/releases),
* Cloning this repo and building from source via: `make` or `make install`.

# Writing Tests

Tests are described by a YAML-based document format. Describe your request, the response you expect, and it's that simple. When you need it, reach for Instaunit's many advanced features. Here's a very simple test suite containing a single test case:

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
* **Run matrixes of related tests** to simplify maintaining tests for repeating endpoint patterns,
* **Wait for dependency services** to become available before a test suite starts running,
* Declare [**mock services**](https://github.com/instaunit/instaunit/wiki/Mock-Services) so your tests avoid external dependencies without making code-level changes,
* Evaluate **expressions and built-in functions** to generate input and randomize your tests.
* **Manage the process you're testing** by starting it before tests are run and stopping it after tests have completed.

# Running Tests

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
             -{"status": "Ok"}
             +{"status": "Not great"}
```

## Producing Reports

Reports describing the results of your tests can be produced in addition to the information logged to standard output. Currently the [JUnit](https://junit.org/junit5/) report format is supported.

Try running `instaunit -report <test_suite>` to generate a report.

You can optionally specify the report output directory with the `-report:output` flag or the generated report format with `-report:type`. The JUnit report format is the default (and currently, only) output format.

# Documenting Tests

Tests and documentation are naturally maintained together: when an endpoint is added or changed you must update your tests as well as the documentation that describes it. To generate documentation, simply add a description to a representative test case for your endpoint. You can pick and choose which tests generate documentation.

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
