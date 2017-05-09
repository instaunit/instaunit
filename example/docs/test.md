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
    Strict-Transport-Security: max-age=31536000
    X-Content-Type-Options: nosniff
    X-Frame-Options: deny
    Accept-Ranges: bytes
    X-Served-By: cache-jfk8142-JFK
    Date: Tue, 09 May 2017 00:05:55 GMT
    X-Timer: S1494288355.409124,VS0,VE0
    Vary: Authorization,Accept-Encoding
    X-Fastly-Request-Id: 12b63cbbb7f75f39b08ec05d3bc27d6e907fa066
    Expires: Tue, 09 May 2017 00:10:55 GMT
    Etag: "9b5817fe39f818381e7cc03a1105ee2686f2831a"
    Content-Type: text/plain; charset=utf-8
    X-Cache: HIT
    Access-Control-Allow-Origin: *
    Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
    X-Xss-Protection: 1; mode=block
    Cache-Control: max-age=300
    X-Geo-Block-List: 
    X-Github-Request-Id: 2CCC:34B4:24397DC:2596EAA:591107A5
    Via: 1.1 varnish
    Connection: keep-alive
    X-Cache-Hits: 2
    Source-Age: 62
    
    Here's a simple
    response from the
    server.
    


