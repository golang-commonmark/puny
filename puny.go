// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package puny provides functions for decoding punycode-encoded strings.
package puny

import (
	"errors"
	"strings"
)

const (
	maxInt32      int32 = 2147483647
	base          int32 = 36
	tMin          int32 = 1
	baseMinusTMin       = base - tMin
	tMax          int32 = 26
	skew          int32 = 38
	damp          int32 = 700
	initialBias   int32 = 72
	initialN      int32 = 128

	acePrefix = "xn--"
	delimiter = '-'
)

var (
	ErrOverflow     = errors.New("overflow: input needs wider integers to process")
	ErrNotBasic     = errors.New("illegal input >= 0x80 (not a basic code point)")
	ErrInvalidInput = errors.New("invalid input")

	ErrNoAtSignInEmail = errors.New("no at sign in the e-mail address")
)

func adapt(delta, numPoints int32, firstTime bool) int32 {
	if firstTime {
		delta /= damp
	} else {
		delta /= 2
	}
	delta += delta / numPoints
	k := int32(0)
	for delta > baseMinusTMin*tMax/2 {
		delta = delta / baseMinusTMin
		k += base
	}
	return k + (baseMinusTMin+1)*delta/(delta+skew)
}

func basicToDigit(b byte) int32 {
	switch {
	case b >= '0' && b <= '9':
		return int32(b - 22)
	case b >= 'A' && b <= 'Z':
		return int32(b - 'A')
	case b >= 'a' && b <= 'z':
		return int32(b - 'a')
	}
	return base
}

func lastIndex(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}

// Decode decodes a punycode-encoded string.
func Decode(s string) (string, error) {
	basic := lastIndex(s, delimiter)
	output := make([]rune, 0, len(s))
	for i := 0; i < basic; i++ {
		b := s[i]
		if b >= 0x80 {
			return "", ErrNotBasic
		}
		output = append(output, rune(b))
	}

	i, n, bias, pos := int32(0), initialN, initialBias, basic+1

	for pos < len(s) {
		oldi, w, k := i, int32(1), base
		for {
			digit := basicToDigit(s[pos])
			pos++

			if digit >= base || digit > (maxInt32-i)/w {
				return "", ErrOverflow
			}

			i += digit * w

			t := k - bias
			if t < tMin {
				t = tMin
			} else if t > tMax {
				t = tMax
			}

			if digit < t {
				break
			}

			if pos == len(s) {
				return "", ErrInvalidInput
			}

			baseMinusT := base - t
			if w > maxInt32/baseMinusT {
				return "", ErrOverflow
			}

			w *= baseMinusT
			k += base
		}

		out := int32(len(output) + 1)
		bias = adapt(i-oldi, out, oldi == 0)

		if i/out > maxInt32-n {
			return "", ErrOverflow
		}

		n += i / out
		i %= out

		output = append(output, 0)
		copy(output[i+1:], output[i:])
		output[i] = rune(n)

		i++
	}

	return string(output), nil
}

func decode(s string) (string, error) {
	if !strings.HasPrefix(s, acePrefix) {
		return s, nil
	}

	uni, err := Decode(s[len(acePrefix):])
	if err != nil {
		return "", err
	}

	return strings.ToLower(uni), nil
}

func issep(r rune) bool {
	return r == '.' || r == '。' || r == '．' || r == '｡'
}

// DecodeHostname decodes a punycode-encoded host name.
func DecodeHostname(s string) (_ string, err error) {
	if !strings.Contains(s, acePrefix) {
		return s, nil
	}

	a := make([]string, 0, 1)
	start := 0
	var uni string
	for i, r := range s {
		if issep(r) {
			uni, err = decode(s[start:i])
			if err != nil {
				return
			}
			a = append(a, uni)
			start = i + 1
			if r != '.' {
				start += 2
			}
		}
	}

	uni, err = decode(s[start:])
	if err != nil {
		return
	}
	a = append(a, uni)

	return strings.Join(a, "."), nil
}

// DecodeEmail decodes a punycode-encoded e-mail address.
func DecodeEmail(s string) (string, error) {
	i := strings.IndexByte(s, '@')
	if i < 0 {
		return "", ErrNoAtSignInEmail
	}

	host, err := DecodeHostname(s[i+1:])
	if err != nil {
		return "", err
	}

	return s[:i+1] + host, nil
}
