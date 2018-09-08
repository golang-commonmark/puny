// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
