// Copyright 2023 Antrea Authors
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

package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const version = "v1.11.0"

func Test_translateRelativeLinks(t *testing.T) {
	// "synced" word in test name means that the file is present in this website repository
	// as it has been copied(synced) from Antrea repository while updating docs.
	tests := []struct {
		name     string
		fullPath string
		content  []byte
		expected []byte
	}{
		{
			name:     "no translation-absolute external link",
			fullPath: "website/content/docs/v1.11.0/somefile.md",
			content:  []byte("texts[filename](https://github.com/antrea-io/antrea/blob/main/test/test.yaml)texts"),
			expected: []byte("texts[filename](https://github.com/antrea-io/antrea/blob/main/test/test.yaml)texts"),
		},
		{
			name:     "no translation-relative link to a file which is synced",
			fullPath: "website/content/docs/v1.11.0/somefile.md",
			content:  []byte("texts[filename](./docs/ci/filename.md)texts"),
			expected: []byte("texts[filename](./docs/ci/filename.md)texts"),
		},
		{
			name:     "no translation-relative link to a file present inside a sub-dir named same as one of the prefixes",
			fullPath: "website/content/docs/v1.11.0/docs/somefile.md",
			content:  []byte("texts[filename](build/filename.md)texts"),
			expected: []byte("texts[filename](build/filename.md)texts"),
		},
		{
			name:     "no translation-relative link to a file at parent dir and synced",
			fullPath: "website/content/docs/v1.11.0/docs/somefile.md",
			content:  []byte("texts[filename](../filename.md)texts"),
			expected: []byte("texts[filename](../filename.md)texts"),
		},
		{
			name:     "translation-relative link from root to a file inside build dir at root of repository and dir not synced",
			fullPath: "website/content/docs/v1.11.0/somefile.md",
			content:  []byte("texts[filename](build/yamls/filename.yml)texts"),
			expected: []byte("texts[filename](https://github.com/antrea-io/antrea/blob/v1.11.0/build/yamls/filename.yml)texts"),
		},
		{
			name:     "translation-absolute link to a file inside build dir at root of repository and dir not synced",
			fullPath: "website/content/docs/v1.11.0/somefile.md",
			content:  []byte("texts[filename](/build/yamls/filename.yaml)texts"),
			expected: []byte("texts[filename](https://github.com/antrea-io/antrea/blob/v1.11.0/build/yamls/filename.yaml)texts"),
		},
		{
			name:     "translation-relative link from child dir to a file inside build dir at root of repository and dir not synced",
			fullPath: "website/content/docs/v1.11.0/docs/somefile.md",
			content:  []byte("texts[filename](../build/yamls/filename.yml)texts"),
			expected: []byte("texts[filename](https://github.com/antrea-io/antrea/blob/v1.11.0/build/yamls/filename.yml)texts"),
		},
		{
			name:     "translation-relative link from child dir to a file at root of repository and not synced",
			fullPath: "website/content/docs/v1.11.0/docs/maintainers/release.md",
			content:  []byte("texts[VERSION](../../VERSION)texts"),
			expected: []byte("texts[VERSION](https://github.com/antrea-io/antrea/blob/v1.11.0/VERSION)texts"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedContent, err := translateRelativeLinks(tt.fullPath, version, tt.content)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, updatedContent)
		})
	}
}
