package lib_test

import (
	"bytes"
	"testing"

	"github.com/GuustTaillieu/idiomatic-go/internal/lib"
)

func BenchmarkUnpooledBuffer(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		buf.WriteString("Hello, World!")
		_ = buf.Bytes()
	}
}

func BenchmarkPooledBuffer(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	pool := lib.NewBufferPool()

	for b.Loop() {
		buf := pool.Get()
		buf.WriteString("Hello, World!")
		_ = buf.Bytes()
		pool.Put(buf)
	}
}
