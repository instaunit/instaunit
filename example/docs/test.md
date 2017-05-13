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
Accept-Ranges: bytes
Via: 1.1 varnish
X-Cache: HIT
X-Geo-Block-List: 
Connection: keep-alive
X-Timer: S1494712410.609448,VS0,VE0
Vary: Authorization,Accept-Encoding
X-Fastly-Request-Id: 5f4c13598bd1c6959fb1ef73597824e7fcd88bc7
Source-Age: 0
Etag: "97a3ed924dddfaaaa2c566cf556e9fa1379cbd02"
Date: Sat, 13 May 2017 21:53:29 GMT
X-Content-Type-Options: nosniff
Content-Type: text/plain; charset=utf-8
X-Served-By: cache-sjc3642-SJC
X-Cache-Hits: 1
Access-Control-Allow-Origin: *
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
Strict-Transport-Security: max-age=31536000
Cache-Control: max-age=300
X-Github-Request-Id: A372:3F31:4246163:44DBB97:59178058
Expires: Sat, 13 May 2017 21:58:29 GMT
X-Frame-Options: deny
X-Xss-Protection: 1; mode=block

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
X-Cache: MISS
X-Cache-Hits: 0
Access-Control-Allow-Origin: *
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
Strict-Transport-Security: max-age=31536000
Etag: "f69bdd0be1577adf16208754ae8f62cd6e7fcb1a"
Accept-Ranges: bytes
Date: Sat, 13 May 2017 21:53:29 GMT
X-Timer: S1494712410.639929,VS0,VE96
X-Content-Type-Options: nosniff
Cache-Control: max-age=300
X-Github-Request-Id: 1448:0220:4A5C770:4D60120:59178058
Connection: keep-alive
X-Served-By: cache-sjc3642-SJC
X-Xss-Protection: 1; mode=block
Via: 1.1 varnish
Expires: Sat, 13 May 2017 21:58:29 GMT
Source-Age: 0
X-Frame-Options: deny
Content-Type: text/plain; charset=utf-8
X-Geo-Block-List: 
Vary: Authorization,Accept-Encoding
X-Fastly-Request-Id: 20242363a93ef617565fac8d7ec029ed3553a915

{
  "/": false,
  "a": 123,
  "z": "Hello, this is the value"
}
```


