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
    Vary: Authorization,Accept-Encoding
    Access-Control-Allow-Origin: *
    X-Fastly-Request-Id: c8bd0febd27956d55e5e7616e9cf2af51c69e9f6
    Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
    X-Xss-Protection: 1; mode=block
    X-Geo-Block-List: 
    Connection: keep-alive
    Strict-Transport-Security: max-age=31536000
    Via: 1.1 varnish
    X-Served-By: cache-jfk8120-JFK
    X-Cache: HIT
    Date: Thu, 11 May 2017 02:01:20 GMT
    Source-Age: 44
    X-Content-Type-Options: nosniff
    X-Frame-Options: deny
    Content-Type: text/plain; charset=utf-8
    X-Github-Request-Id: 095A:34B4:30D6A2B:32AAD17:5913C5C4
    Expires: Thu, 11 May 2017 02:06:20 GMT
    Etag: "9b5817fe39f818381e7cc03a1105ee2686f2831a"
    Cache-Control: max-age=300
    Accept-Ranges: bytes
    X-Cache-Hits: 2
    X-Timer: S1494468080.319668,VS0,VE0
    
    Here's a simple
    response from the
    server.
    


## GET https://raw.githubusercontent.com/bww/hunit/master/example/entity.json
Testing json request

### Example request

    GET /bww/hunit/master/example/entity.json HTTP/1.1
    Host: raw.githubusercontent.com
    Content-Type: application/json
    Origin: localhost
    
    
### Example response

    HTTP/1.1 200 OK
    Source-Age: 44
    X-Xss-Protection: 1; mode=block
    X-Github-Request-Id: 56AC:02E6:1586082:16789E5:5913C5C4
    Via: 1.1 varnish
    Access-Control-Allow-Origin: *
    X-Fastly-Request-Id: 6c4a166e4a7cda7becd2bf16e51855b5279f59a3
    X-Cache: HIT
    Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
    X-Content-Type-Options: nosniff
    X-Frame-Options: deny
    Etag: "bd5fdee5109d273e7d9d396848c196eb43ab9f77"
    Date: Thu, 11 May 2017 02:01:20 GMT
    X-Served-By: cache-jfk8120-JFK
    X-Timer: S1494468080.367789,VS0,VE1
    Accept-Ranges: bytes
    Connection: keep-alive
    X-Cache-Hits: 1
    Strict-Transport-Security: max-age=31536000
    Content-Type: text/plain; charset=utf-8
    Cache-Control: max-age=300
    X-Geo-Block-List: 
    Vary: Authorization,Accept-Encoding
    Expires: Thu, 11 May 2017 02:06:20 GMT
    
    {
      "z": "Hello, this is the value",
      "a": 123,
      "/": false
    }
    


