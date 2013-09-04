// Copyright 2013 Eric Myhre
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package iox

import (
	"bytes"
	"io"
	"strings"
)

func ReaderFromInterface(x interface{}) io.Reader {
	switch y := x.(type) {
	case string:
		return ReaderFromString(y)
	case []byte:
		return ReaderFromByteSlice(y)
	case io.Reader:
		return y
	case bytes.Buffer:
		return &y
	case <-chan string:
		return ReaderFromChanString(y)
	case chan string:
		return ReaderFromChanString(y)
	case <-chan []byte:
		return ReaderFromChanByteSlice(y)
	case chan []byte:
		return ReaderFromChanByteSlice(y)
	default:
		return nil
	}
}

func ReaderFromString(str string) io.Reader {
	return strings.NewReader(str)
}

func ReaderFromByteSlice(bats []byte) io.Reader {
	return bytes.NewReader(bats)
}

func ReaderFromChanString(ch <-chan string) io.Reader {
	return &readerChanString{ch: ch}
}

type readerChanString struct {
	ch  <-chan string
	buf []byte
}

func (r *readerChanString) Read(p []byte) (n int, err error) {
	w := 0
	if len(r.buf) == 0 {
		// skip
	} else if len(p) >= len(r.buf) {
		// copy whole buffer out
		w = copy(p, r.buf)
		r.buf = r.buf[0:0]
	} else {
		// not room for the whole buffer; copy what there's room for, shift buf, return.
		w = copy(p, r.buf[:len(p)])
		r.buf = r.buf[len(p):0]
		return w, nil
	}

	str, open := <-r.ch
	bytes := []byte(str)
	w2 := copy(p, bytes)
	r.buf = bytes[w2:]

	if open || len(r.buf) > 0 {
		return w + w2, nil
	} else {
		return w + w2, io.EOF
	}
}

func ReaderFromChanByteSlice(ch <-chan []byte) io.Reader {
	return &readerChanByteSlice{ch: ch}
}

type readerChanByteSlice struct {
	ch  <-chan []byte
	buf []byte
}

func (r *readerChanByteSlice) Read(p []byte) (n int, err error) {
	w := 0
	if len(r.buf) == 0 {
		// skip
	} else if len(p) >= len(r.buf) {
		// copy whole buffer out
		w = copy(p, r.buf)
		r.buf = r.buf[0:0]
	} else {
		// not room for the whole buffer; copy what there's room for, shift buf, return.
		w = copy(p, r.buf[:len(p)])
		r.buf = r.buf[len(p):0]
		return w, nil
	}

	bytes, open := <-r.ch
	w2 := copy(p, bytes)
	r.buf = bytes[w2:]

	if open || len(r.buf) > 0 {
		return w + w2, nil
	} else {
		return w + w2, io.EOF
	}
}
