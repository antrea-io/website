// Copyright 2021 Antrea Authors
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
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type DocDir struct {
	path      string
	filter    string
	recursive bool
}

var (
	AntreaRepo  string
	WebsiteRepo string
	DryRun      bool
)

func CopyFile(source, dest string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func CopyDir(source, dest string) error {
	return filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.Mkdir(filepath.Join(dest, relPath), os.ModePerm)
		}
		return CopyFile(path, filepath.Join(dest, relPath))
	})
}

func ignoreDestFile(path string) bool {
	if filepath.Base(path) == "_index.md" {
		return true
	}
	return false
}

func syncDirs(sourceDocsPath string, destDocsPath string, docDir *DocDir) error {
	re := regexp.MustCompile(docDir.filter)
	sourceFiles := make(map[string]string)
	if err := filepath.WalkDir(filepath.Join(sourceDocsPath, docDir.path), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if !docDir.recursive && path != sourceDocsPath {
				return fs.SkipDir
			}
			return nil
		}
		if docDir.filter == "" || re.MatchString(path) {
			relPath, err := filepath.Rel(sourceDocsPath, path)
			if err != nil {
				return err
			}
			sourceFiles[relPath] = path
		}
		return nil
	}); err != nil {
		return fmt.Errorf("error when walking directory '%s': %w", sourceDocsPath, err)
	}

	destFiles := make(map[string]string)
	if err := filepath.WalkDir(filepath.Join(destDocsPath, docDir.path), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if !docDir.recursive && path != destDocsPath {
				return fs.SkipDir
			}
			return nil
		}
		if ignoreDestFile(path) {
			return nil
		}
		relPath, err := filepath.Rel(destDocsPath, path)
		if err != nil {
			return err
		}
		destFiles[relPath] = path
		return nil
	}); err != nil {
		return fmt.Errorf("error when walking directory '%s': %w", destDocsPath, err)
	}

	for id, sourcePath := range sourceFiles {
		destPath, exists := destFiles[id]
		if exists {
			log.Printf("Syncing %s -> %s\n", sourcePath, destPath)
		} else {
			destPath = filepath.Join(destDocsPath, id)
			log.Printf("Syncing %s -> %s [NEW]\n", sourcePath, destPath)
		}
		if !DryRun {
			if err := CopyFile(sourcePath, destPath); err != nil {
				return fmt.Errorf("error when copying file: %w", err)
			}
		}
	}

	for id, destPath := range destFiles {
		_, exists := sourceFiles[id]
		if exists {
			continue
		}
		log.Printf("Deleting %s\n", destPath)
		if !DryRun {
			if err := os.Remove(destPath); err != nil {
				return fmt.Errorf("error when deleting file: %w", err)
			}
		}
	}

	return nil
}

func fixupMarkdownFile(path string) error {
	// Handle HTML <img> tags in Markdown
	imgTag := regexp.MustCompile(`<(img (?s).*?)>`)
	md, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	md = imgTag.ReplaceAllFunc(md, func(m []byte) []byte {
		return []byte(strings.ReplaceAll(string(m), "\n", ""))
	})
	md = imgTag.ReplaceAll(md, []byte("{{< $1 >}}"))
	if DryRun {
		return nil
	}
	if err := ioutil.WriteFile(path, md, 0644); err != nil {
		return err
	}
	return nil
}

func fixupMarkdown(destDocsPath string) error {
	re := regexp.MustCompile("^.*md$")
	destFiles := make([]string, 0)
	if err := filepath.WalkDir(destDocsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if ignoreDestFile(path) {
			return nil
		}
		if !re.MatchString(path) {
			return nil
		}
		destFiles = append(destFiles, path)
		return nil
	}); err != nil {
		return fmt.Errorf("error when walking directory '%s': %w", destDocsPath, err)
	}

	for _, destFile := range destFiles {
		log.Printf("Fixing up markdown file %s\n", destFile)
		if err := fixupMarkdownFile(destFile); err != nil {
			return fmt.Errorf("error when fixing up markdown file: %w", err)
		}
	}

	return nil
}

func init() {
	flag.StringVar(&AntreaRepo, "antrea-repo", "", "Path to the Antrea repo")
	flag.StringVar(&WebsiteRepo, "website-repo", "", "Path to the Antrea website")
	flag.BoolVar(&DryRun, "dry-run", false, "Do a dry-run (do not modify any website source files")
}

func UpdateDocs(destDocsPath string) error {
	if AntreaRepo == "" || WebsiteRepo == "" {
		return fmt.Errorf("flags -antrea-repo and -website-repo are required")
	}

	sourceDocsPath := filepath.Join(AntreaRepo)

	if stat, err := os.Stat(sourceDocsPath); err != nil || !stat.IsDir() {
		return fmt.Errorf("'%s' is not a valid directory", sourceDocsPath)
	}

	if stat, err := os.Stat(destDocsPath); err != nil || !stat.IsDir() {
		return fmt.Errorf("'%s' is not a valid directory", destDocsPath)
	}

	docDirs := []DocDir{
		DocDir{path: "", filter: "^.*md$", recursive: false},
		DocDir{path: "docs", filter: "", recursive: true},
	}
	for _, docDir := range docDirs {
		if err := syncDirs(sourceDocsPath, destDocsPath, &docDir); err != nil {
			return fmt.Errorf("error when syncing doc dir %v: %w", docDir, err)
		}
	}

	if err := fixupMarkdown(destDocsPath); err != nil {
		return fmt.Errorf("error when fixing-up markdown files: %w", err)
	}

	return nil
}
