// package bing provides web search functionalities by scraping bing search engine.
package bing

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

/*
Search searches a query on bing.
NOTE: results may not be empty even if the error is not nil.
Because, an error can be occured at 1000th page.
So you can still be able to get the data of previous 999 pages.
;)
*/
func Search(query string, blacklist []string) (results []string, err error) {

	blacklist = append(
		blacklist,
		"go.microsoft.com",
	)
	return searchwithpagelimit2(
		query,
		blacklist,
		time.Duration(-1),
	)
}

/*
SearchWithTimeout searches a query on bing with timeout, usefull for the long result queries.
NOTE: results may not be empty even if the error is not nil.
Because, an error can be occured at 1000th page.
So you can still be able to get the data of previous 999 pages.
;)
*/
func SearchWithTimeout(query string, blacklist []string, duration time.Duration) (results []string, err error) {

	blacklist = append(
		blacklist,
		"go.microsoft.com",
	)
	return search(
		query,
		blacklist,
		duration,
	)
}

// Unique removes duplicate values from the given list.
func Unique(sites []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range sites {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// HostOnly returns only valid hostname from a result list.
func HostOnly(list []string) []string {
	var s []string
	for _, uri := range list {
		v, err := url.ParseRequestURI(
			uri,
		)
		if err != nil {
			continue
		}
		s = append(s, "http://"+v.Host)
		// list only domains (not the full url)
		//[BUG only http:// listing available, https:// are not available ]
	}
	// decline the duplicate hosts
	// because a result list can contain many url from a single host
	return Unique(s)
}

// search searches the dorks on bing
func search(searchStr string, blacklist []string, duration time.Duration) (results []string, err error) {
	var timeout uint32
	if searchStr == "" {
		return results, errors.New(
			"bing: empty string given",
		)
	}
	if duration > 0 {
		go func() {
			<-time.After(duration)
			atomic.AddUint32(&timeout, 1)
		}()
	}
	// pre: previous page's results
	// nex: next(current) page's results
	// sey: final results to add to collection
	var pre, nex, sey []string

	/* Searching starts here */
	var page int
	nex = []string{""}

	// Search until we reach the last page
	for {
		if atomic.LoadUint32(&timeout) != uint32(0) {
			return results, err
		}
		//Bing page indicator
		page += 10
		client := &http.Client{
			Timeout: time.Duration(10 * time.Second),
		}
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf(
				"https://www.bing.com/search?q=%s&first=%d",
				url.QueryEscape(searchStr),
				page,
			),
			nil,
		)
		if err != nil {
			return results, err
		}
		// Spoof a old browser user agent
		// for getting a simple html page
		req.Header.Set("User-Agent",
			"Nokia2700c/10.0.011 (SymbianOS/9.4; U; Series60/5.0 Opera/5.0; Profile/MIDP-2.1 Configuration/CLDC-1.1 ) AppleWebKit/525 (KHTML, like Gecko) Safari/525 3gpp-gba",
		)

		resp, err := client.Do(req)
		if err != nil {
			return results, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(
			resp.Body,
		)
		if err != nil {
			return results, err
		}
		/* HTTP client ends here! we got the response in body */
		anchors := strings.Split(
			string(body),
			`<a _ctf="rdr_T" href="http`,
		)

		if len(anchors) < 10 || eq(pre, nex) {
			// TODO return results, the result
		}
		pre = nex
		nex = []string{}
		// lets extract the links from <a> tag
		for i := range anchors {
			if i == 0 {
				continue
			}
			href := strings.Split(
				anchors[i],
				`"`,
			)
			if kickBlacklist(
				href[0],
				blacklist,
			) {
				// ignore denied links or sites
				continue
			}
			nex = append(nex, href[0])
			// append them with previous results
			results = append(
				results,
				"http"+href[0],
			)
			sey = append(
				sey,
				"http"+href[0],
			)
		}
		resp.Body.Close()
	}
	return results, nil
}

// search searches the dorks on bing
func searchwithpagelimit2(searchStr string, blacklist []string, duration time.Duration) (results []string, err error) {
	var timeout uint32
	if searchStr == "" {
		return results, errors.New(
			"bing: empty string given",
		)
	}
	if duration > 0 {
		go func() {
			<-time.After(duration)
			atomic.AddUint32(&timeout, 1)
		}()
	}
	// pre: previous page's results
	// nex: next(current) page's results
	// sey: final results to add to collection
	var pre, nex, sey []string

	/* Searching starts here */
	var page int
	nex = []string{""}

	// Search until we reach the last page
	for i := 0; i<2; i ++ {
		if atomic.LoadUint32(&timeout) != uint32(0) {
			return results, err
		}
		//Bing page indicator
		page += 10
		client := &http.Client{
			Timeout: time.Duration(10 * time.Second),
		}
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf(
				"https://www.bing.com/search?q=%s&first=%d",
				url.QueryEscape(searchStr),
				page,
			),
			nil,
		)
		if err != nil {
			return results, err
		}
		// Spoof a old browser user agent
		// for getting a simple html page
		req.Header.Set("User-Agent",
			"Nokia2700c/10.0.011 (SymbianOS/9.4; U; Series60/5.0 Opera/5.0; Profile/MIDP-2.1 Configuration/CLDC-1.1 ) AppleWebKit/525 (KHTML, like Gecko) Safari/525 3gpp-gba",
		)

		resp, err := client.Do(req)
		if err != nil {
			return results, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(
			resp.Body,
		)
		if err != nil {
			return results, err
		}
		/* HTTP client ends here! we got the response in body */
		anchors := strings.Split(
			string(body),
			`<a _ctf="rdr_T" href="http`,
		)

		if len(anchors) < 10 || eq(pre, nex) {
			// TODO return results, the result
		}
		pre = nex
		nex = []string{}
		// lets extract the links from <a> tag
		for i := range anchors {
			if i == 0 {
				continue
			}
			href := strings.Split(
				anchors[i],
				`"`,
			)
			if kickBlacklist(
				href[0],
				blacklist,
			) {
				// ignore denied links or sites
				continue
			}
			nex = append(nex, href[0])
			// append them with previous results
			results = append(
				results,
				"http"+href[0],
			)
			sey = append(
				sey,
				"http"+href[0],
			)
		}
		resp.Body.Close()
	}
	return results, nil
}

// kickBlacklist returns the true if the site/link matches any of Deny params
func kickBlacklist(str string, blck []string) bool {
	for i := range blck {
		if blck[i] != "" {
			if strings.Contains(str, blck[i]) {
				return true
			}
		}
	}
	return false
}

/* checking the last page (a bit tricky :p )*/
func eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
