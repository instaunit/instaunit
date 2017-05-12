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
Ok: FOOBAR!!
Origin: localhost

```
### Example response

```http
HTTP/1.1 200 OK
Via: 1.1 varnish
X-Cache: HIT
X-Cache-Hits: 1
Vary: Authorization,Accept-Encoding
X-Xss-Protection: 1; mode=block
Accept-Ranges: bytes
Date: Fri, 12 May 2017 00:21:22 GMT
Access-Control-Allow-Origin: *
Cache-Control: max-age=300
Connection: keep-alive
Source-Age: 0
X-Served-By: cache-jfk8122-JFK
X-Timer: S1494548482.259046,VS0,VE0
Expires: Fri, 12 May 2017 00:26:22 GMT
X-Frame-Options: deny
Content-Type: text/plain; charset=utf-8
X-Geo-Block-List: 
X-Github-Request-Id: E92C:34B7:CB35B6:D3A157:59150001
X-Fastly-Request-Id: a8562db21548a9d6e768712376b357b3e489c92a
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
Strict-Transport-Security: max-age=31536000
X-Content-Type-Options: nosniff
Etag: "9b5817fe39f818381e7cc03a1105ee2686f2831a"

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
Ok: FOOBAR!!
Origin: localhost

```
### Example response

```http
HTTP/1.1 200 OK
X-Geo-Block-List: 
Accept-Ranges: bytes
Via: 1.1 varnish
X-Frame-Options: deny
Etag: "bd5fdee5109d273e7d9d396848c196eb43ab9f77"
Cache-Control: max-age=300
Date: Fri, 12 May 2017 00:21:22 GMT
X-Cache: MISS
Expires: Fri, 12 May 2017 00:26:22 GMT
Connection: keep-alive
X-Served-By: cache-jfk8122-JFK
X-Cache-Hits: 0
X-Timer: S1494548482.284359,VS0,VE31
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
X-Content-Type-Options: nosniff
X-Xss-Protection: 1; mode=block
Vary: Authorization,Accept-Encoding
Source-Age: 0
Access-Control-Allow-Origin: *
X-Fastly-Request-Id: 97220a922aa5fb8ed6fb98e5aeea868141fc54ea
Strict-Transport-Security: max-age=31536000
Content-Type: text/plain; charset=utf-8
X-Github-Request-Id: B8BE:02E8:29F8034:2BCCA0F:59150002

{
  "z": "Hello, this is the value",
  "a": 123,
  "/": false
}
```


