package main

import (
	"fmt"
	"strconv"

	lru "github.com/hashicorp/golang-lru/v2"
)

func main() {
	cache, err := lru.New[string, any](100)
	if err != nil {
		panic(err)
	}

	nums := 100
	for ix := 0; ix < nums; ix++ {
		cache.Add("k"+strconv.Itoa(ix), "v"+strconv.Itoa(ix))
		cache.Get("k" + strconv.Itoa(ix+1))
	}
	fmt.Println(cache.Len())
	cache.Add("k"+strconv.Itoa(101), 2)
	fmt.Println(cache.Get("k" + strconv.Itoa(1)))
}
