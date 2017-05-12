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

```http
GET /bww/hunit/master/example/entity.txt HTTP/1.1
Host: raw.githubusercontent.com
Origin: localhost

```
### Example response

```http
HTTP/1.1 200 OK
Source-Age: 156
X-Geo-Block-List: 
X-Github-Request-Id: 964A:34B9:22B2082:2408244:5915044E
X-Served-By: cache-jfk8132-JFK
X-Cache: HIT
X-Cache-Hits: 2
Expires: Fri, 12 May 2017 00:47:18 GMT
Strict-Transport-Security: max-age=31536000
Content-Type: text/plain; charset=utf-8
Accept-Ranges: bytes
Connection: keep-alive
Vary: Authorization,Accept-Encoding
Via: 1.1 varnish
X-Timer: S1494549738.408763,VS0,VE0
Access-Control-Allow-Origin: *
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
X-Content-Type-Options: nosniff
X-Frame-Options: deny
X-Xss-Protection: 1; mode=block
Cache-Control: max-age=300
X-Fastly-Request-Id: c6e92114a36fbce0d40294c83bc5791f7ca0414c
Etag: "9b5817fe39f818381e7cc03a1105ee2686f2831a"
Date: Fri, 12 May 2017 00:42:18 GMT

Here's a simple
response from the
server.
```


## GET https://raw.githubusercontent.com/bww/hunit/master/example/entity.json

An entity

### Example request

```http
GET /bww/hunit/master/example/entity.json HTTP/1.1
Host: raw.githubusercontent.com
Origin: localhost

```
### Example response

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Connection: keep-alive
Accept-Ranges: bytes
Vary: Authorization,Accept-Encoding
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
X-Content-Type-Options: nosniff
X-Xss-Protection: 1; mode=block
Cache-Control: max-age=300
X-Served-By: cache-jfk8132-JFK
X-Timer: S1494549738.449243,VS0,VE1
Access-Control-Allow-Origin: *
X-Fastly-Request-Id: a622bc9b89d090b0d45d9cb0c752e5931b4d0d09
X-Frame-Options: deny
Etag: "bd5fdee5109d273e7d9d396848c196eb43ab9f77"
X-Geo-Block-List: 
Via: 1.1 varnish
Expires: Fri, 12 May 2017 00:47:18 GMT
X-Cache-Hits: 1
Source-Age: 156
Strict-Transport-Security: max-age=31536000
X-Github-Request-Id: 6AAE:02E6:183CC7C:1952513:5915044E
Date: Fri, 12 May 2017 00:42:18 GMT
X-Cache: HIT

{
  "z": "Hello, this is the value",
  "a": 123,
  "/": false
}
```


