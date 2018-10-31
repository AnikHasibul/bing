# bing
--
    import "github.com/anikhasibul/bing"

package bing provides web search functionalities by scraping bing search engine.

## Usage

#### func  Search

```go
func Search(query string, blacklist []string) (results []string, err error)
```
Search searches a query on bing. NOTE: results may not be empty even if the
error is not nil. Because, an error can be occured at 1000th page. So you can
still be able to get the data of previous 999 pages. ;)

#### func  SearchWithTimeout

```go
func SearchWithTimeout(query string, blacklist []string, duration time.Duration) (results []string, err error)
```
SearchWithTimeout searches a query on bing with timeout, usefull for the long
result queries. NOTE: results may not be empty even if the error is not nil.
Because, an error can be occured at 1000th page. So you can still be able to get
the data of previous 999 pages. ;)

#### func  Unique

```go
func Unique(sites []string) []string
```
Unique removes duplicate values from the given list.

#### func  HostOnly

```go
func HostOnly(list []string) []string
```
HostOnly returns only valid hostname from a result list.

