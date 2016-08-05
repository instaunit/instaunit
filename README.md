# HUnit

A simple testing tool for HTTP APIs.

	$ hunit test.yml
    ====> test.yml
    ----> GET https://raw.githubusercontent.com/bww/hunit/master/example/test.txt
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


### test.yml

	- 
      request:
        method: GET
        url: https://raw.githubusercontent.com/bww/hunit/master/example/test.txt
        headers:
          Origin: localhost
      
      response:
        status: 200
        headers:
          Content-Type: text/plain; charset=utf-8
        entity: |+
          Heres a simple
          response from the
          server.      

