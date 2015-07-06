package puny

import "testing"

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Decode("xn----7sbab8abeduuih7bb2byd6cycj")
	}
}

func BenchmarkDecodeHostname(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DecodeHostname("xn----7sbab8abeduuih7bb2byd6cycj.xn--p1ai")
	}
}
