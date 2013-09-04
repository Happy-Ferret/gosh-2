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

package psh

import (
	"fmt"
	"reflect"
)

/**
 * Error when Sh() or its family of functions is called with arguments of an unexpected
 * type.  Sh() functions only expect arguments of the public types declared in the
 * sh_modifiers.go file when setting up a command.
 *
 * This should mostly be a compile-time problem as long as you write your
 * script to not actually pass unchecked types of interface{} into Sh() commands.
 */
type IncomprehensibleCommandModifier struct {
	wat *interface{}
}

func (err IncomprehensibleCommandModifier) Error() string {
	return fmt.Sprintf("sh: incomprehensible command modifier: do not want type \"%v\"", whoru(reflect.ValueOf(*err.wat)))
}

func whoru(val reflect.Value) string {
	kind := val.Kind()
	typ := val.Type()

	if kind == reflect.Ptr {
		return fmt.Sprintf("*%s", whoru(val.Elem()))
	} else if kind == reflect.Interface {
		return whoru(val.Elem())
	} else {
		return typ.Name()
	}
}

/**
 * Error for commands run by Sh that exited with a non-successful status.
 *
 * What exactly qualifies as an unsuccessful status can be defined per command,
 * but by default is any exit code other than zero.
 */
type FailureExitCode struct {
	cmdname string
	code    int
}

func (err FailureExitCode) Error() string {
	return fmt.Sprintf("sh: command \"%s\" exited with unexpected status %d", err.cmdname, err.code)
}
