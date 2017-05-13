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
Origin: localhost
Host: api.example.com

```
### Example response

```http
HTTP/1.1 200 OK
X-Cache-Hits: 2
Vary: Authorization,Accept-Encoding
Content-Type: text/plain; charset=utf-8
Cache-Control: max-age=300
X-Geo-Block-List: 
Accept-Ranges: bytes
Expires: Sat, 13 May 2017 23:05:34 GMT
X-Content-Type-Options: nosniff
Date: Sat, 13 May 2017 23:00:34 GMT
Connection: keep-alive
X-Timer: S1494716434.275117,VS0,VE0
X-Fastly-Request-Id: a42be1ac55eeb5c2d07d12d6bb54273e4b043ad7
Source-Age: 37
X-Xss-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000
X-Frame-Options: deny
Etag: "97a3ed924dddfaaaa2c566cf556e9fa1379cbd02"
Via: 1.1 varnish
Access-Control-Allow-Origin: *
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'

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
Host: api.example.com

```
### Example response

```http
HTTP/1.1 200 OK
X-Fastly-Request-Id: 4195f4962c3b346a00767e9f31070f10d9184dae
Expires: Sat, 13 May 2017 23:05:34 GMT
Date: Sat, 13 May 2017 23:00:34 GMT
Via: 1.1 varnish
Vary: Authorization,Accept-Encoding
Accept-Ranges: bytes
X-Timer: S1494716434.314311,VS0,VE1
Access-Control-Allow-Origin: *
X-Xss-Protection: 1; mode=block
Etag: "f69bdd0be1577adf16208754ae8f62cd6e7fcb1a"
Content-Type: text/plain; charset=utf-8
X-Geo-Block-List: 
Connection: keep-alive
Source-Age: 37
Strict-Transport-Security: max-age=31536000
X-Frame-Options: deny
X-Cache-Hits: 1
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
X-Content-Type-Options: nosniff
Cache-Control: max-age=300

{
  "/": false,
  "a": 123,
  "z": "Hello, this is the value"
}
```


