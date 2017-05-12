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
Etag: "9b5817fe39f818381e7cc03a1105ee2686f2831a"
Content-Type: text/plain; charset=utf-8
Accept-Ranges: bytes
Via: 1.1 varnish
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
X-Content-Type-Options: nosniff
X-Xss-Protection: 1; mode=block
Source-Age: 4
X-Cache-Hits: 4
X-Fastly-Request-Id: 217c6140bd53affd2ecff764c5fc32b072475b29
Expires: Fri, 12 May 2017 00:39:40 GMT
Date: Fri, 12 May 2017 00:34:40 GMT
X-Served-By: cache-jfk8145-JFK
X-Cache: HIT
X-Github-Request-Id: 3938:34B4:36DF476:38F20DA:5915031C
Connection: keep-alive
X-Timer: S1494549281.938835,VS0,VE0
Strict-Transport-Security: max-age=31536000
X-Frame-Options: deny
Cache-Control: max-age=300
Access-Control-Allow-Origin: *
X-Geo-Block-List: 
Vary: Authorization,Accept-Encoding

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
X-Frame-Options: deny
X-Github-Request-Id: 243E:02E4:99986D:A11207:5915031C
Date: Fri, 12 May 2017 00:34:40 GMT
X-Served-By: cache-jfk8145-JFK
X-Cache: HIT
X-Fastly-Request-Id: 5dc7c777223c5282d6c8ac710dbf699b676865d8
Source-Age: 4
Strict-Transport-Security: max-age=31536000
Etag: "bd5fdee5109d273e7d9d396848c196eb43ab9f77"
Vary: Authorization,Accept-Encoding
X-Xss-Protection: 1; mode=block
Accept-Ranges: bytes
Connection: keep-alive
X-Cache-Hits: 1
X-Timer: S1494549281.965496,VS0,VE0
Access-Control-Allow-Origin: *
X-Geo-Block-List: 
X-Content-Type-Options: nosniff
Content-Type: text/plain; charset=utf-8
Cache-Control: max-age=300
Via: 1.1 varnish
Expires: Fri, 12 May 2017 00:39:40 GMT
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'

{
  "z": "Hello, this is the value",
  "a": 123,
  "/": false
}
```


