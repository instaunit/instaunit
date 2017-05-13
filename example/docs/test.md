# Some Example Tests

This suite-level documentation is copied to the top of the generated
documentation file, when documentation is generated.
  
  * The first example thing,
  * The second example thing,
  * The third and final example thing.

Read more about [HUnit](https://github.com/bww/hunit).

## Contents

* [GET /example/:entity_name](#get-exampleentityname)
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
X-Xss-Protection: 1; mode=block
Accept-Ranges: bytes
Date: Sat, 13 May 2017 21:47:54 GMT
Via: 1.1 varnish
X-Fastly-Request-Id: 687a9fe1ef6bb1433ed197a39dc32ce9b4a51872
X-Geo-Block-List: 
X-Github-Request-Id: 23C4:3F31:4241A99:44D7165:59177F09
X-Cache: HIT
X-Timer: S1494712074.244150,VS0,VE0
Vary: Authorization,Accept-Encoding
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
X-Content-Type-Options: nosniff
Etag: "97a3ed924dddfaaaa2c566cf556e9fa1379cbd02"
Cache-Control: max-age=300
X-Served-By: cache-sjc3151-SJC
X-Cache-Hits: 1
Strict-Transport-Security: max-age=31536000
X-Frame-Options: deny
Content-Type: text/plain; charset=utf-8
Connection: keep-alive
Access-Control-Allow-Origin: *
Expires: Sat, 13 May 2017 21:52:54 GMT
Source-Age: 0

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
Accept-Ranges: bytes
Connection: keep-alive
X-Fastly-Request-Id: b20480fde11abae2a39f46668b2715c76236ce49
X-Xss-Protection: 1; mode=block
Etag: "f69bdd0be1577adf16208754ae8f62cd6e7fcb1a"
Content-Type: text/plain; charset=utf-8
Date: Sat, 13 May 2017 21:47:54 GMT
X-Served-By: cache-sjc3151-SJC
X-Cache: MISS
X-Timer: S1494712074.281772,VS0,VE92
Vary: Authorization,Accept-Encoding
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
Strict-Transport-Security: max-age=31536000
X-Frame-Options: deny
Access-Control-Allow-Origin: *
X-Github-Request-Id: 6C04:0218:6990CA:6EF727:59177F0A
Expires: Sat, 13 May 2017 21:52:54 GMT
X-Content-Type-Options: nosniff
Cache-Control: max-age=300
X-Geo-Block-List: 
Via: 1.1 varnish
X-Cache-Hits: 0
Source-Age: 0

{
  "/": false,
  "a": 123,
  "z": "Hello, this is the value"
}
```


