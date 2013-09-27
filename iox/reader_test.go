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
	"github.com/coocood/assrt"
	"io"
	"testing"
)

func TestReaderFromChanString(t *testing.T) {
	assert := assrt.NewAssert(t)

	ch := make(chan string)
	var output bytes.Buffer
	go func() {
		ch <- "asdf"
		ch <- ""
		ch <- "\nwakawaka"
		ch <- "\tz"
		close(ch)
	}()
	io.Copy(&output, ReaderFromChanString(ch))

	assert.Equal(
		"asdf\nwakawaka\tz",
		output.String(),
	)
}

func TestReaderFromChanByteSlice(t *testing.T) {
	assert := assrt.NewAssert(t)

	ch := make(chan []byte)
	var output bytes.Buffer
	go func() {
		ch <- []byte("asdf")
		ch <- []byte("")
		ch <- []byte("\nwakawaka")
		ch <- []byte("\tz")
		close(ch)
	}()
	io.Copy(&output, ReaderFromChanByteSlice(ch))

	assert.Equal(
		"asdf\nwakawaka\tz",
		output.String(),
	)
}

type clearlyNotAReader struct{}

func TestReaderUnrefinable(t *testing.T) {
	assert := assrt.NewAssert(t)

	var x clearlyNotAReader

	defer func() {
		err := recover()
		switch y := err.(type) {
		case error:
			assert.Equal(
				"ReaderFromInterface cannot refine type \"iox.clearlyNotAReader\" to a Reader",
				y.Error(),
			)
		default:
			t.Fatal("recover returned a non-error type!")
		}
	}()
	ReaderFromInterface(x)
}

func TestReaderFromChanByteSliceIsClosable(t *testing.T) {
	ch := make(chan []byte)
	reader := ReaderFromInterface(ch)
	if _, ok := reader.(io.ReadCloser); !ok {
		t.Fatalf("did not get a reader that supported close; did want")
	}
}

func TestReaderFromChanReadonlyByteSliceIsNotClosable(t *testing.T) {
	ch := make(<-chan []byte)
	reader := ReaderFromInterface(ch)
	if _, ok := reader.(io.ReadCloser); ok {
		t.Fatalf("got a reader that supported close; did not want")
	}
}
