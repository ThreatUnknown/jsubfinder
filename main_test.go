package main

import (
	"testing"
	)

var result int

func BenchmarkMain(b *testing.B) {
	for x := 0; x < b.N; x++ {
		//run my awesome test method
		main()
		//fmt.Printf("Y = %d\n", y)
	}
}
