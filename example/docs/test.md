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
Source-Age: 266
X-Frame-Options: deny
X-Geo-Block-List: 
Date: Fri, 12 May 2017 16:05:43 GMT
X-Served-By: cache-jfk8129-JFK
X-Fastly-Request-Id: 8f427c0a1db8a9675797009fcb58a746a70a4c47
X-Timer: S1494605144.748150,VS0,VE0
Vary: Authorization,Accept-Encoding
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
Strict-Transport-Security: max-age=31536000
X-Content-Type-Options: nosniff
X-Xss-Protection: 1; mode=block
Accept-Ranges: bytes
X-Cache-Hits: 2
Expires: Fri, 12 May 2017 16:10:43 GMT
X-Github-Request-Id: 09C8:34B4:3A72D9D:3CA91F4:5915DC4D
Via: 1.1 varnish
X-Cache: HIT
Access-Control-Allow-Origin: *
Etag: "97a3ed924dddfaaaa2c566cf556e9fa1379cbd02"
Content-Type: text/plain; charset=utf-8
Cache-Control: max-age=300
Connection: keep-alive

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
Date: Fri, 12 May 2017 16:05:43 GMT
Via: 1.1 varnish
Connection: keep-alive
X-Cache: HIT
Access-Control-Allow-Origin: *
Strict-Transport-Security: max-age=31536000
X-Content-Type-Options: nosniff
X-Frame-Options: deny
X-Github-Request-Id: 6520:02E6:19F456C:1B1D3A5:5915DC4D
Accept-Ranges: bytes
X-Fastly-Request-Id: 78d78f39e382f71fd54080391efa4ef032625af5
Source-Age: 266
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
X-Xss-Protection: 1; mode=block
Cache-Control: max-age=300
X-Cache-Hits: 1
X-Timer: S1494605144.778335,VS0,VE1
Vary: Authorization,Accept-Encoding
Content-Type: text/plain; charset=utf-8
X-Geo-Block-List: 
Expires: Fri, 12 May 2017 16:10:43 GMT
Etag: "f69bdd0be1577adf16208754ae8f62cd6e7fcb1a"
X-Served-By: cache-jfk8129-JFK

{
  "/": false,
  "a": 123,
  "z": "Hello, this is the value"
}
```


