# Some Example Tests

This suite-level documentation is copied to the top of the generated
documentation file, when documentation is generated.
  
  * The first example thing,
  * The second example thing,
  * The third and final example thing.

Read more about [HUnit](https://github.com/bww/hunit).

## Contents

* [GET /example/:entity_name](#get-exampleentity_name)
* [GET https://raw.githubusercontent.com/bww/hunit/master/example/entity.json](#get-httpsrawgithubusercontentcombwwhunitmasterexampleentityjson)

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
Date: Sat, 13 May 2017 22:03:46 GMT
Via: 1.1 varnish
X-Cache: HIT
X-Cache-Hits: 1
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
X-Xss-Protection: 1; mode=block
Etag: "97a3ed924dddfaaaa2c566cf556e9fa1379cbd02"
Connection: keep-alive
Expires: Sat, 13 May 2017 22:08:46 GMT
X-Geo-Block-List: 
X-Github-Request-Id: 5F42:3F31:424ED5F:44E4DD7:591782C2
Accept-Ranges: bytes
Cache-Control: max-age=300
X-Served-By: cache-sjc3650-SJC
X-Timer: S1494713027.724161,VS0,VE0
Strict-Transport-Security: max-age=31536000
X-Content-Type-Options: nosniff
Content-Type: text/plain; charset=utf-8
Source-Age: 0
X-Fastly-Request-Id: b110ed48cef06a65d6646e65d96f5eb21a3d122b
X-Frame-Options: deny
Vary: Authorization,Accept-Encoding
Access-Control-Allow-Origin: *

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
X-Content-Type-Options: nosniff
Via: 1.1 varnish
X-Cache: MISS
X-Timer: S1494713027.753970,VS0,VE95
X-Fastly-Request-Id: 908ef5ffbd99bf759ad5b7a694cf553a9cbd3a38
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
X-Xss-Protection: 1; mode=block
Etag: "f69bdd0be1577adf16208754ae8f62cd6e7fcb1a"
Cache-Control: max-age=300
Accept-Ranges: bytes
Date: Sat, 13 May 2017 22:03:46 GMT
Connection: keep-alive
X-Served-By: cache-sjc3650-SJC
X-Cache-Hits: 0
Content-Type: text/plain; charset=utf-8
X-Geo-Block-List: 
X-Github-Request-Id: B85C:0220:4A66572:4D6A5DC:591782C1
Source-Age: 0
Strict-Transport-Security: max-age=31536000
X-Frame-Options: deny
Vary: Authorization,Accept-Encoding
Access-Control-Allow-Origin: *
Expires: Sat, 13 May 2017 22:08:46 GMT

{
  "/": false,
  "a": 123,
  "z": "Hello, this is the value"
}
```


