# Some Example Tests

This suite-level documentation is copied to the top of the generated
documentation file, when documentation is generated.
  
  * The first example thing,
  * The second example thing,
  * The third and final example thing.

Read more about [Instaunit](https://github.com/instaunit/instaunit).

## Contents

* [GET /example/:entity_name](#get-exampleentity_name)
* [GET https://raw.githubusercontent.com/instaunit/instaunit/master/example/entity.json](#get-httpsrawgithubusercontentcominstaunitinstaunitmasterexampleentityjson)

## GET /example/:entity_name

Fetch the *entity text* from Github.

The entity text is just some example text used to illustrate how literal
entites can be compared in a test case using [*Instaunit*](https://github.com/instaunit/instaunit).

### Example request

```http
GET /instaunit/instaunit/master/example/entity.txt HTTP/1.1
Host: raw.githubusercontent.com
Origin: localhost

```
### Example response

```http
HTTP/1.1 200 OK
X-Github-Request-Id: 8980:154BF:A25F28:A82111:5951C458
Date: Tue, 27 Jun 2017 02:38:15 GMT
Connection: keep-alive
Vary: Authorization,Accept-Encoding
X-Fastly-Request-Id: f6a1318d5de7080b4f3a8f36385bd541c5c71bf9
Strict-Transport-Security: max-age=31536000
X-Frame-Options: deny
X-Xss-Protection: 1; mode=block
Via: 1.1 varnish
X-Served-By: cache-jfk8129-JFK
X-Cache: HIT
Expires: Tue, 27 Jun 2017 02:43:15 GMT
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
Source-Age: 192
Cache-Control: max-age=300
X-Geo-Block-List: 
X-Timer: S1498531095.488822,VS0,VE0
Access-Control-Allow-Origin: *
X-Content-Type-Options: nosniff
Content-Type: text/plain; charset=utf-8
Accept-Ranges: bytes
X-Cache-Hits: 4
Etag: "97a3ed924dddfaaaa2c566cf556e9fa1379cbd02"

Here's a simple
response from the
server.

```


## GET https://raw.githubusercontent.com/instaunit/instaunit/master/example/entity.json

An entity

### Example request

```http
GET /instaunit/instaunit/master/example/entity.json HTTP/1.1
Host: raw.githubusercontent.com
Origin: localhost

```
### Example response

```http
HTTP/1.1 200 OK
Connection: keep-alive
Access-Control-Allow-Origin: *
X-Xss-Protection: 1; mode=block
Etag: "f69bdd0be1577adf16208754ae8f62cd6e7fcb1a"
X-Geo-Block-List: 
X-Served-By: cache-jfk8129-JFK
X-Cache-Hits: 1
X-Timer: S1498531096.523180,VS0,VE1
Vary: Authorization,Accept-Encoding
X-Fastly-Request-Id: 4e625d9fd33216932c7428865e10619909b58bc0
Strict-Transport-Security: max-age=31536000
X-Content-Type-Options: nosniff
X-Github-Request-Id: 887A:18B9:23DA906:2521DE6:5951C457
Expires: Tue, 27 Jun 2017 02:43:15 GMT
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
X-Frame-Options: deny
Accept-Ranges: bytes
Date: Tue, 27 Jun 2017 02:38:15 GMT
Via: 1.1 varnish
X-Cache: HIT
Content-Type: text/plain; charset=utf-8
Cache-Control: max-age=300
Source-Age: 192

{
  "/": false,
  "a": 123,
  "z": "Hello, this is the value"
}
```


