package speeds_test

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/requester/rio/speeds"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	r := speeds.NewRateLimit(100)
	fmt.Println("adding 101...")
	r.Add(101)
	fmt.Println("adding 10...")
	r.Add(10)
	fmt.Println("adding 11...")
	r.Add(11)
	fmt.Println("adding 12...")
	r.Add(12)
	fmt.Println("adding 13...")
	r.Add(13)
	fmt.Println("adding 22...")
	r.Add(22)
	fmt.Println("adding 35...")
	r.Add(35)
	fmt.Println("adding 25...")
	r.Add(25)
	fmt.Println("adding 11...")
	r.Add(11)

	r.Stop()
	time.Sleep(10e9)
}
