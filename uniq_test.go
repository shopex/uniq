package main

import (
	"testing"
	"fmt"
)

func TestOpen(t *testing.T){
    var err error
    uniq, err = Open("foo", 1 << 12, 3)
    if err != nil {
        t.Error("open error", err)
    }
}

func TestWrite(t *testing.T){
	uniq.Write([]byte("foo"), 0)
	fmt.Println(u.Test([]byte("foo"), 10))
}