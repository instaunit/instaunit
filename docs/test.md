## GET /example/:entity_name
Fetch the *entity text* from Github.

The entity text is just some example text used to illustrate how literal
entites can be compared in a test case using [*HUnit*](https://github.com/shangsunset/hunit).

### Example request

    GET /shangsunset/hunit/master/example/entity.txt HTTP/1.1
    Host: raw.githubusercontent.com
    Origin: localhost
    
    
### Example response

    HTTP/1.1 200 OK
    X-Frame-Options: deny
    Date: Wed, 10 May 2017 04:17:24 GMT
    Via: 1.1 varnish
    X-Cache-Hits: 2
    X-Timer: S1494389844.391998,VS0,VE0
    Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'
    X-Content-Type-Options: nosniff
    X-Geo-Block-List: 
    X-Github-Request-Id: 45B4:4018:23EB840:255D67E:5912942D
    Connection: keep-alive
    X-Cache: HIT
    Access-Control-Allow-Origin: *
    Expires: Wed, 10 May 2017 04:22:24 GMT
    Strict-Transport-Security: max-age=31536000
    Cache-Control: max-age=300
    Accept-Ranges: bytes
    X-Served-By: cache-jfk8131-JFK
    Vary: Authorization,Accept-Encoding
    X-Fastly-Request-Id: d4b2112511585ec1c41ad87627a6bff1dce98027
    X-Xss-Protection: 1; mode=block
    Content-Type: text/plain; charset=utf-8
    Source-Age: 39
    Etag: "adf1355d43286ca52615dbef27463748d1178d24"
    
    
    {
      "results": [
        {
          "id": "1d86dd70-1a2a-48a3-b3e0-5d125ce64872",
          "created_at": "2017-05-08T19:23:39.802Z",
          "updated_at": "2017-05-08T19:23:39.802Z",
          "resume_path": "path/to/resume",
          "rating": null
        },
        {
          "id": "8f465a1f-09b5-437e-8014-6f856c515555",
          "created_at": "2017-05-08T19:23:39.844Z",
          "updated_at": "2017-05-08T19:23:39.844Z",
          "resume_path": "path/to/resume",
          "candidate": {
            "id": "5910c5bb55840926b91a930d",
            "slugs": [
              "joelknlsadnfkjsdfn-blowsdl-29"
            ],
            "email": "hello@example.com",
            "summary": "Hello, I'm a candidate.",
            "name": {
              "first": "Joelknlsadnfkjsdfn",
              "last": "Blowsdl"
            }
          },
          "rating": 4
        },
        {
          "id": "34ce0ff3-c3d9-4709-8644-3354bb333c64",
          "created_at": "2017-05-08T19:23:39.877Z",
          "updated_at": "2017-05-08T19:23:39.877Z",
          "resume_path": "other/path/to/resume",
          "candidate": {
            "id": "5910c5bb55840926b91a930d",
            "slugs": [
              "joelknlsadnfkjsdfn-blowsdl-29"
            ],
            "email": "hello@example.com",
            "summary": "Hello, I'm a candidate.",
            "name": {
              "first": "Joelknlsadnfkjsdfn",
              "last": "Blowsdl"
            }
          },
          "rating": null
        }
      ],
      "page": {
        "this_page": 0,
        "page_length": 25
      },
      "meta": {
        "total_results": 3
      }
    }
    
    


