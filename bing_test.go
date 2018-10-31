package bing

import (
	"fmt"
	"testing"
	"time"
)

func TestSearch(t *testing.T) {
	res, err := SearchWithTimeout("hello", []string{"wiki"}, time.Duration(time.Minute))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}
