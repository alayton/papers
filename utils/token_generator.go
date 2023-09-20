package utils

/*
Adapted from https://github.com/kjk/betterguid
Increased random bits from 72 bits/12 characters to 144 bits/24 characters
Always regenerate random bits instead of incrementing on the same timestamp

The MIT License (MIT)

Copyright (c) 2015 Krzysztof Kowalczyk

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import (
	"math/rand"
	"strings"
	"time"
)

const (
	// Modeled after base64 web-safe chars, but ordered by ASCII.
	pushChars        = "-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
	tokenTimeBytes   = 8
	tokenRandomBytes = 32
	tokenShortBytes  = 8
	tokenLength      = tokenTimeBytes + tokenRandomBytes
	tokenShortLength = tokenTimeBytes + tokenShortBytes
)

var (
	rnd *rand.Rand
)

func init() {
	// seed to get randomness
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Creates a unique token suitable for confirmation and recovery
func NewToken() string {
	var token [tokenTimeBytes + tokenRandomBytes]byte

	// put current time at the beginning
	timeMs := time.Now().UTC().UnixNano() / 1e6
	for i := tokenTimeBytes - 1; i >= 0; i-- {
		n := int(timeMs % 64)
		token[i] = pushChars[n]
		timeMs = timeMs / 64
	}

	for i := tokenTimeBytes; i < tokenLength; i++ {
		token[i] = pushChars[rnd.Intn(64)]
	}

	return string(token[:])
}

// Creates a shorter token for identifying access and refresh token chains
func NewChain() string {
	var token [tokenTimeBytes + tokenShortBytes]byte

	// put current time at the beginning
	timeMs := time.Now().UTC().UnixNano() / 1e6
	for i := tokenTimeBytes - 1; i >= 0; i-- {
		n := int(timeMs % 64)
		token[i] = pushChars[n]
		timeMs = timeMs / 64
	}

	for i := tokenTimeBytes; i < tokenShortLength; i++ {
		token[i] = pushChars[rnd.Intn(64)]
	}

	return string(token[:])
}

// Returns the time a token was created at
func TokenTime(token string) time.Time {
	var t = []byte(token)

	ms := 0
	for i := 0; i < tokenTimeBytes; i++ {
		v := strings.IndexByte(pushChars, t[i])
		ms *= 64
		ms += v
	}

	return time.UnixMilli(int64(ms))
}
