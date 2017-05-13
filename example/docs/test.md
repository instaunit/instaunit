# Some Example Tests

This suite-level documentation is copied to the top of the generated
documentation file, when documentation is generated.
  
  * The first example thing,
  * The second example thing,
  * The third and final example thing.

Read more about [HUnit](https://github.com/bww/hunit).

## Contents

* [GET https://raw.githubusercontent.com/bww/hunit/master/example/entity.json](#get-httpsrawgithubusercontentcombwwhunitmasterexampleentityjson)
* [GET /example/:entity_name](#get-exampleentity_name)

## GET /example/:entity_name

Fetch the *entity text* from Github.

The entity text is just some example text used to illustrate how literal
entites can be compared in a test case using [*HUnit*](https://github.com/bww/hunit).

### Example request

```http
GET /bww/hunit/master/example/entity.txt HTTP/1.1
Origin: localhost

```
### Example response

```http
HTTP/1.1 200 OK
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
Cache-Control: max-age=300
X-Geo-Block-List: 
X-Github-Request-Id: 36F2:3F2F:2BF3DCB:2DB2CF2:59178FEC
Accept-Ranges: bytes
Date: Sat, 13 May 2017 23:02:41 GMT
Via: 1.1 varnish
X-Cache: HIT
Strict-Transport-Security: max-age=31536000
Etag: "97a3ed924dddfaaaa2c566cf556e9fa1379cbd02"
Connection: keep-alive
X-Served-By: cache-sjc3127-SJC
X-Frame-Options: deny
X-Xss-Protection: 1; mode=block
Access-Control-Allow-Origin: *
X-Fastly-Request-Id: 4e8796ff4b2d6b962679b629d4a83ca8c6d5f1e8
Expires: Sat, 13 May 2017 23:07:41 GMT
X-Content-Type-Options: nosniff
X-Cache-Hits: 2
X-Timer: S1494716561.261468,VS0,VE0
Vary: Authorization,Accept-Encoding
Source-Age: 164

Here's a simple
response from the
server.

```


## GET https://raw.githubusercontent.com/bww/hunit/master/example/entity.json

An entity

### Example request

```http
GET /bww/hunit/master/example/entity.json HTTP/1.1
Origin: localhost

```
### Example response

```http
HTTP/1.1 200 OK
X-Xss-Protection: 1; mode=block
X-Github-Request-Id: 4A58:0220:4A9AAE6:4DA1284:59178FEC
Accept-Ranges: bytes
Access-Control-Allow-Origin: *
Source-Age: 164
X-Frame-Options: deny
X-Fastly-Request-Id: dd224dfab41576b9f09107a4369dfb4ef132dd78
Cache-Control: max-age=300
Etag: "f69bdd0be1577adf16208754ae8f62cd6e7fcb1a"
X-Geo-Block-List: 
Date: Sat, 13 May 2017 23:02:41 GMT
Via: 1.1 varnish
X-Served-By: cache-sjc3127-SJC
X-Cache: HIT
Vary: Authorization,Accept-Encoding
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
Expires: Sat, 13 May 2017 23:07:41 GMT
X-Content-Type-Options: nosniff
Connection: keep-alive
X-Cache-Hits: 1
X-Timer: S1494716561.297269,VS0,VE1
Strict-Transport-Security: max-age=31536000

{
  "/": false,
  "a": 123,
  "z": "Hello, this is the value"
}
```


