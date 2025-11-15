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