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
X-Frame-Options: deny
X-Geo-Block-List: 
Via: 1.1 varnish
Strict-Transport-Security: max-age=31536000
X-Xss-Protection: 1; mode=block
Etag: "9b5817fe39f818381e7cc03a1105ee2686f2831a"
Content-Type: text/plain; charset=utf-8
Cache-Control: max-age=300
Date: Fri, 12 May 2017 15:56:17 GMT
X-Fastly-Request-Id: b6e6c1473e413ee3f75f3b0fc58728e6811c04a5
Expires: Fri, 12 May 2017 16:01:17 GMT
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
X-Content-Type-Options: nosniff
Connection: keep-alive
X-Served-By: cache-jfk8137-JFK
Vary: Authorization,Accept-Encoding
Access-Control-Allow-Origin: *
Source-Age: 262
X-Github-Request-Id: 8B76:2606:D11AE4:DA5611:5915DA1A
Accept-Ranges: bytes
X-Cache: HIT
X-Cache-Hits: 2
X-Timer: S1494604577.063598,VS0,VE0

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
Cache-Control: max-age=300
Accept-Ranges: bytes
X-Served-By: cache-jfk8137-JFK
X-Cache: HIT
X-Xss-Protection: 1; mode=block
Date: Fri, 12 May 2017 15:56:17 GMT
X-Timer: S1494604577.094940,VS0,VE1
X-Fastly-Request-Id: 6c8aa8b4234ce45f7264f8e7e6bcda7e2c8c4c90
Source-Age: 262
X-Geo-Block-List: 
X-Frame-Options: deny
Etag: "bd5fdee5109d273e7d9d396848c196eb43ab9f77"
X-Github-Request-Id: 79CC:02E6:19EF3A6:1B17F29:5915DA1B
Via: 1.1 varnish
Vary: Authorization,Accept-Encoding
Access-Control-Allow-Origin: *
Strict-Transport-Security: max-age=31536000
Expires: Fri, 12 May 2017 16:01:17 GMT
X-Content-Type-Options: nosniff
Content-Type: text/plain; charset=utf-8
Connection: keep-alive
X-Cache-Hits: 1
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'

{
  "z": "Hello, this is the value",
  "a": 123,
  "/": false
}
```


