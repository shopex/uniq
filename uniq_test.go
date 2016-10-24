package main

import (
	"github.com/bradfitz/gomemcache/memcache"
	"strconv"
	"testing"
)

func TestOpen(t *testing.T) {
	var err error
	uniq, err = Open("data", 1<<12, 3)
	if err != nil {
		t.Error("open error", err)
	}
}

func TestWriteCore(t *testing.T) {
	s := uniq.TestAndAdd([]byte("foo"), 0)
	t.Log(s)
}

func BenchmarkWriteCore(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uniq.TestAndAdd([]byte("foo"), 0)
	}
}

func BenchmarkWithProto(b *testing.B) {
	mc := memcache.New("127.0.0.1:6532")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Get(strconv.Itoa(i) + "foo")
	}
}
