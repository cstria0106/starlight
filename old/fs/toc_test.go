/*
   Copyright The starlight Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   file created by maverick in 2021
*/

package fs

import (
	"fmt"
	"testing"
)

func TestListxattr(t *testing.T) {
	m := map[string]string{
		"asdfasdf": "asdfasdf",
		"123":      "asdfasdf",
		"asdfe":    "asdfasdf",
		"as23dfe":  "asdfasdf",
	}

	for i := range m {
		fmt.Println(i)
	}
}
