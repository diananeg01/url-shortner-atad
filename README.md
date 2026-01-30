# URL SHORTENER

## First checkpoint (16.11.2025)

* wrote a few lines of go code to get to familiar with the syntax
* tested `gomponents` and `go-echarts`, played a bit with them
* created database schema and connection
* for custom given sequence (no rules set yet), generate short url
  * in the page, the render url is just an anchor that uses the original url
  * redirect can be tested with `localhost:8080/redirect/<given sequence>`
  * url is stored in the database (no rules on unicity yet)
* no login yet
* no error handling yet

## Final submission checklist

| User Stories                                                           | Done |
|------------------------------------------------------------------------|------|
| I can customize my short link                                          | [x]  |
| I can view statistics: total clicks, unique visitors                   | [x]  |
| I can generate a QR code for my short link                             | [x]  |
| I can set an expiration date for my links                              | [x]  |


| Technical Requirements                                                 | Done |
|------------------------------------------------------------------------|------|
| RESTful API with endpoints for shortening, redirecting, and analytics  | [x]  |
| Short code generation (6-8 characters, alphanumeric)                   | [x]  |
| Collision detection and handling                                       | [ ]  |
| Rate limiting (e.g., 10 requests/minute per IP)                        | [ ]  |
| Web dashboard showing all links and statistics                         | [x]  |
| Database for URLs and click events                                     | [x]  |
