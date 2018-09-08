// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package puny

import "testing"

type testCase struct {
	desc string
	in   string
	want string
	err  error
}

func TestDecode(t *testing.T) {
	testCases := []testCase{
		{
			in: "", want: "",
		},
		{
			desc: "a single basic code point",
			in:   "Bach-",
			want: "Bach",
		},
		{
			desc: "a single non-ASCII character",
			in:   "tda",
			want: "ü",
		},
		{
			desc: "multiple non-ASCII characters",
			in:   "4can8av2009b",
			want: "üëäö♥",
		},
		{
			desc: "mix of ASCII and non-ASCII characters",
			in:   "bcher-kva",
			want: "bücher",
		},
		{
			desc: "long string with both ASCII and non-ASCII characters",
			in:   "Willst du die Blthe des frhen, die Frchte des spteren Jahres-x9e96lkal",
			want: "Willst du die Blüthe des frühen, die Früchte des späteren Jahres",
		},
		{
			desc: "Arabic (Egyptian)",
			in:   "egbpdaj6bu4bxfgehfvwxn",
			want: "ليهمابتكلموشعربي؟",
		},
		{
			desc: "Chinese (simplified)",
			in:   "ihqwcrb4cv8a8dqg056pqjye",
			want: "他们为什么不说中文",
		},
		{
			desc: "Chinese (traditional)",
			in:   "ihqwctvzc91f659drss3x8bo0yb",
			want: "他們爲什麽不說中文",
		},
		{
			desc: "Czech",
			in:   "Proprostnemluvesky-uyb24dma41a",
			want: "Pročprostěnemluvíčesky",
		},
		{
			desc: "Hebrew",
			in:   "4dbcagdahymbxekheh6e0a7fei0b",
			want: "למההםפשוטלאמדבריםעברית",
		},
		{
			desc: "Hindi (Devanagari)",
			in:   "i1baa7eci9glrd9b2ae1bj0hfcgg6iyaf8o0a1dig0cd",
			want: "यहलोगहिन्दीक्योंनहींबोलसकतेहैं",
		},
		{
			desc: "Japanese (kanji and hiragana)",
			in:   "n8jok5ay5dzabd5bym9f0cm5685rrjetr6pdxa",
			want: "なぜみんな日本語を話してくれないのか",
		},
		{
			desc: "Korean (Hangul syllables)",
			in:   "989aomsvi5e83db1d2a355cv1e0vak1dwrv93d5xbh15a0dt30a5jpsd879ccm6fea98c",
			want: "세계의모든사람들이한국어를이해한다면얼마나좋을까",
		},
		{
			desc: "Russian (Cyrillic)",
			in:   "b1abfaaepdrnnbgefbadotcwatmq2g4l",
			want: "почемужеонинеговорятпорусски",
		},
		{
			desc: "Spanish",
			in:   "PorqunopuedensimplementehablarenEspaol-fmd56a",
			want: "PorquénopuedensimplementehablarenEspañol",
		},
		{
			desc: "Vietnamese",
			in:   "TisaohkhngthchnitingVit-kjcr8268qyxafd2f1b9g",
			want: "TạisaohọkhôngthểchỉnóitiếngViệt",
		},
		{
			in:   "3B-ww4c5e180e575a65lsy2b",
			want: "3年B組金八先生",
		},
		{
			in:   "-with-SUPER-MONKEYS-pc58ag80a8qai00g7n9n",
			want: "安室奈美恵-with-SUPER-MONKEYS",
		},
		{
			in:   "Hello-Another-Way--fc4qua05auwb3674vfr0b",
			want: "Hello-Another-Way-それぞれの場所",
		},
		{
			in:   "2-u9tlzr9756bt3uc0v",
			want: "ひとつ屋根の下2",
		},
		{
			in:   "MajiKoi5-783gue6qz075azm5e",
			want: "MajiでKoiする5秒前",
		},
		{
			in:   "de-jg4avhby1noc0d",
			want: "パフィーdeルンバ",
		},
		{
			in:   "d9juau41awczczp",
			want: "そのスピードで",
		},
		{
			desc: "ASCII string that breaks the existing rules for host-name labels",
			in:   "-> $1.00 <--",
			want: "-> $1.00 <-",
		},
		{
			desc: "Punycode in uppercase",
			in:   "bcher-KVA",
			want: "bücher",
		},
	}
	for _, tc := range testCases {
		got, _ := Decode(tc.in)
		if got != tc.want {
			t.Errorf("Decode(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestDecodeError(t *testing.T) {
	testCases := []testCase{
		{
			in:  `\%&`,
			err: ErrOverflow,
		},
		{
			in:  "0000000d000000w00000xb",
			err: ErrOverflow,
		},
		{
			in:  "ква-kva",
			err: ErrNotBasic,
		},
		{
			in:  "UB4",
			err: ErrInvalidInput,
		},
	}
	for _, tc := range testCases {
		_, err := Decode(tc.in)
		if err != tc.err {
			t.Errorf("Decode(%q): expected error %v, got %v", tc.in, tc.err, err)
		}
	}
}

func TestDecodeHostname(t *testing.T) {
	testCases := []testCase{
		{
			in: "", want: "",
		},
		{
			in:   "xn--maana-pta.com",
			want: "mañana.com",
		},
		{
			in:   "example.com.",
			want: "example.com.",
		},
		{
			in:   "xn--bcher-kva.com",
			want: "bücher.com",
		},
		{
			in:   "xn--caf-dma.com",
			want: "café.com",
		},
		{
			in:   "xn----dqo34k.com",
			want: "☃-⌘.com",
		},
		{
			in:   "xn----dqo34kn65z.com",
			want: "퐀☃-⌘.com",
		},
		{
			desc: "Emoji",
			in:   "xn--ls8h.la",
			want: "💩.la",
		},
		{
			desc: "Using U+002E as separator",
			in:   "xn--maana-pta.com",
			want: "mañana.com",
		},
		{
			desc: "Using U+3002 as separator",
			in:   "xn--maana-pta。com",
			want: "mañana.com",
		},
		{
			desc: "Using U+FF0E as separator",
			in:   "xn--maana-pta．com",
			want: "mañana.com",
		},
		{
			desc: "Using U+FF61 as separator",
			in:   "xn--maana-pta｡com",
			want: "mañana.com",
		},
		{
			desc: "Cyrillic hostname",
			in:   "xn----7sbhfccwe7ahoby0si.xn--p1ai",
			want: "единая-мордовия.рф",
		},
		{
			in:   "xn--UB4",
			want: "",
		},
		{
			in:   "xn--UB4.xn--UB4",
			want: "",
		},
	}
	for _, tc := range testCases {
		got, _ := DecodeHostname(tc.in)
		if got != tc.want {
			t.Errorf("DecodeHostname(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestDecodeEmail(t *testing.T) {
	testCases := []testCase{
		{
			in: "", want: "",
		},
		{
			in: "admin@xn--UB4", want: "",
		},
		{
			in:   "ЛюбимыйРуководитель@xn----7sbhfccwe7ahoby0si.xn--p1ai",
			want: "ЛюбимыйРуководитель@единая-мордовия.рф",
		},
	}
	for _, tc := range testCases {
		got, _ := DecodeEmail(tc.in)
		if got != tc.want {
			t.Errorf("DecodeEmail(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
