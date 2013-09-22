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
	. "fmt"
)

/*
	Error raised by ReaderFromInterface() when is called with an argument of an unexpected type.
 */
type ReaderUnrefinableFromInterface struct {
	wat interface{}
}

func (err ReaderUnrefinableFromInterface) Error() string {
	return Sprintf("ReaderFromInterface cannot refine type \"%T\" to a Reader", err.wat)
}
