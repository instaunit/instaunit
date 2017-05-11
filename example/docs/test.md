# Some Example Tests

This suite-level documentation is copied to the top of the generated
documentation file, when documentation is generated.
  
  * The first example thing,
  * The second example thing,
  * The third and final example thing.

Read more about [HUnit](https://github.com/bww/hunit).

## GET /example/:entity_name
Fetch the *entity text* from Github.

The entity text is just some example text used to illustrate how literal
entites can be compared in a test case using [*HUnit*](https://github.com/bww/hunit).
### Example request

    GET /bww/hunit/master/example/entity.txt HTTP/1.1
    Host: raw.githubusercontent.com
    Origin: localhost
    
    
### Example response

    HTTP/1.1 200 OK
    X-Frame-Options: deny
    Accept-Ranges: bytes
    Access-Control-Allow-Origin: *
    Source-Age: 203
    Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
    Strict-Transport-Security: max-age=31536000
    Cache-Control: max-age=300
    X-Geo-Block-List: 
    Date: Thu, 11 May 2017 00:28:11 GMT
    X-Cache: HIT
    Vary: Authorization,Accept-Encoding
    X-Content-Type-Options: nosniff
    X-Xss-Protection: 1; mode=block
    Connection: keep-alive
    X-Cache-Hits: 2
    X-Fastly-Request-Id: f9b2f4a9f48a6f5d3e4332fdf9c36eaceb542656
    Expires: Thu, 11 May 2017 00:33:11 GMT
    Etag: "9b5817fe39f818381e7cc03a1105ee2686f2831a"
    Content-Type: text/plain; charset=utf-8
    Via: 1.1 varnish
    X-Served-By: cache-jfk8120-JFK
    X-Timer: S1494462491.095142,VS0,VE0
    X-Github-Request-Id: A618:34B9:1ED6E34:20038A1:5913AF50
    
    Here's a simple
    response from the
    server.
    


