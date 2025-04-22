// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !windows
// +build !windows

package renameio

import (
	"io"
	"os"
)

// WriteFile mirrors ioutil.WriteFile, replacing an existing file with the same
// name atomically.
func WriteFile(filename string, data []byte, perm os.FileMode, opts ...Option) error {
	opts = append([]Option{
		WithPermissions(perm),
		WithExistingPermissions(),
	}, opts...)

	t, err := NewPendingFile(filename, opts...)
	if err != nil {
		return err
	}
	defer t.Cleanup()

	if _, err := t.Write(data); err != nil {
		return err
	}

	return t.CloseAtomicallyReplace()
}

// WriteFileReader write to a tmp file replacing an existing file with the same name atomically.
func WriteFileReader(filename string, reader io.Reader, perm os.FileMode, opts ...Option) (int64, error) {
	opts = append([]Option{
		WithPermissions(perm),
		WithExistingPermissions(),
	}, opts...)

	t, err := NewPendingFile(filename, opts...)
	if err != nil {
		return 0, err
	}
	defer t.Cleanup()

	written, err := io.Copy(t, reader)
	if err != nil {
		return 0, err
	}

	return written, t.CloseAtomicallyReplace()
}
