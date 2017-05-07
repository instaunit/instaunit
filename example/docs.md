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
    Date: Sun, 07 May 2017 22:04:32 GMT
    X-Cache: HIT
    X-Cache-Hits: 1
    Strict-Transport-Security: max-age=31536000
    Etag: "9b5817fe39f818381e7cc03a1105ee2686f2831a"
    X-Github-Request-Id: C2AC:34B9:1350B7E:140E7B4:590F99F0
    X-Served-By: cache-jfk8147-JFK
    Vary: Authorization,Accept-Encoding
    Access-Control-Allow-Origin: *
    X-Content-Type-Options: nosniff
    X-Frame-Options: deny
    Cache-Control: max-age=300
    Accept-Ranges: bytes
    Via: 1.1 varnish
    Connection: keep-alive
    X-Timer: S1494194673.585176,VS0,VE0
    Source-Age: 0
    X-Xss-Protection: 1; mode=block
    Content-Type: text/plain; charset=utf-8
    Expires: Sun, 07 May 2017 22:09:32 GMT
    Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
    X-Geo-Block-List: 
    X-Fastly-Request-Id: 5ca28d71c8cf88dc2253655d1f9739f0f685ae5b
    
    Here's a simple
    response from the
    server.
    


