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

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"antrea.io/website/scripts/pkg"
)

var (
	version string
)

type Version struct {
	major  int
	minor  int
	patch  int
	suffix string
}

func (v Version) IsNewMinorVersion() bool {
	return v.patch == 0
}

func (v Version) String() string {
	if v.suffix != "" {
		return fmt.Sprintf("v%v.%v.%v-%s", v.major, v.minor, v.patch, v.suffix)
	}
	return fmt.Sprintf("v%v.%v.%v", v.major, v.minor, v.patch)
}

func (v Version) MinorVersion() Version {
	return Version{
		major:  v.major,
		minor:  v.minor,
		patch:  0,
		suffix: "",
	}
}

func (v Version) PrevPatchVersion() Version {
	if v.patch == 0 {
		return Version{
			major:  v.major,
			minor:  v.minor,
			patch:  0,
			suffix: "",
		}
	}
	return Version{
		major:  v.major,
		minor:  v.minor,
		patch:  v.patch - 1,
		suffix: "",
	}
}

func (v Version) HasSuffix() bool {
	return v.suffix != ""
}

func (v Version) LessThan(other Version) bool {
	if v.major < other.major {
		return true
	}
	if v.major > other.major {
		return false
	}
	// v.major == other.major
	if v.minor < other.minor {
		return true
	}
	if v.minor > other.minor {
		return false
	}
	// v.minor == other.minor
	if v.patch < other.patch {
		return true
	}
	if v.patch > other.patch {
		return false
	}
	// v.patch == other.patch
	if v.suffix == "" && other.suffix != "" {
		return false
	}
	if v.suffix != "" && other.suffix == "" {
		return true
	}
	// v.suffix != "" && other.suffix != ""
	// good enough for now...
	return (strings.Compare(v.suffix, other.suffix) < 0)
}

func parseVersion(versionString string) (Version, error) {
	version := Version{}
	re := regexp.MustCompile(`^v([\d]+)\.([\d]+)\.([\d]+)(-(.*))?$`)
	match := re.FindStringSubmatch(versionString)
	if match == nil {
		return version, fmt.Errorf("not a valid version string: %s", versionString)
	}
	var err error
	if version.major, err = strconv.Atoi(match[1]); err != nil {
		return version, fmt.Errorf("cannot convert major version number to int")
	}
	if version.minor, err = strconv.Atoi(match[2]); err != nil {
		return version, fmt.Errorf("cannot convert minor version number to int")
	}
	if version.patch, err = strconv.Atoi(match[3]); err != nil {
		return version, fmt.Errorf("cannot convert patch version number to int")
	}
	version.suffix = match[5]
	return version, nil
}

func mustParseVersion(versionString string) Version {
	version, err := parseVersion(versionString)
	if err != nil {
		panic(fmt.Sprintf("Not a semver: %s", versionString))
	}
	return version
}

func updateVersionInFrontMatter(destDocsPath string, version string) error {
	versionVariable := regexp.MustCompile(`version:\s*([\w\.\-]+)`)
	return filepath.WalkDir(destDocsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Base(path) != "_index.md" {
			return nil
		}
		md, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		match := versionVariable.FindSubmatch(md)
		if match == nil {
			log.Printf("No version variable in front matter for %s\n", path)
			return nil
		}
		if string(match[1]) == version {
			log.Printf("Version front matter variable in %s matches expected version\n", path)
			return nil
		}
		log.Printf("Version front matter variable in %s doesn't match expected version\n", path)
		if pkg.DryRun {
			return nil
		}
		log.Printf("Updating version front matter variable in %s\n", path)
		md = versionVariable.ReplaceAll(md, []byte(fmt.Sprintf("version: %s", version)))
		return os.WriteFile(path, md, 0644)
	})
}

func createDocsIfNeeded(destDocsPath string, version string, referenceVersion string) error {
	destDocsPathExists := true
	if _, err := os.Stat(destDocsPath); err != nil {
		destDocsPathExists = false
	}

	if !destDocsPathExists {
		referenceDocsPath := filepath.Join(pkg.WebsiteRepo, "content", "docs", referenceVersion)
		log.Printf("Copying %s -> %s\n", referenceDocsPath, destDocsPath)
		if !pkg.DryRun {
			if err := pkg.CopyDir(referenceDocsPath, destDocsPath); err != nil {
				return fmt.Errorf("error when copying reference docs from '%s': %w", referenceDocsPath, err)
			}

			if err := updateVersionInFrontMatter(destDocsPath, version); err != nil {
				return fmt.Errorf("error when updating version variable in front matter: %w", err)
			}
		}
	} else {
		log.Printf("Docs directory '%s' exists\n", destDocsPath)
	}

	return nil
}

func createTocFileIfNeeded(tocFile string, referenceVersion string) error {
	tocFileExists := true
	if _, err := os.Stat(tocFile); err != nil {
		tocFileExists = false
	}

	if !tocFileExists {
		referenceTocFile := filepath.Join(pkg.WebsiteRepo, "data", "docs", fmt.Sprintf("%s-toc.yml", referenceVersion))
		log.Printf("Copying %s -> %s\n", referenceTocFile, tocFile)
		if !pkg.DryRun {
			if err := pkg.CopyFile(referenceTocFile, tocFile); err != nil {
				return fmt.Errorf("error when copying reference TOC file from '%s': %w", referenceTocFile, err)
			}
		}
	} else {
		log.Printf("TOC file '%s' exists", tocFile)
	}

	return nil
}

func updateTocMapping(tocMappingPath string, version string) error {
	re := regexp.MustCompile(fmt.Sprintf(`%s:\s*%s-toc`, version, version))
	b, err := os.ReadFile(tocMappingPath)
	if err != nil {
		return err
	}
	if re.Find(b) != nil {
		log.Printf("TOC mapping up-to-date\n")
		return nil
	}
	log.Printf("Updating TOC mapping")
	if !pkg.DryRun {
		b = append(b, []byte(fmt.Sprintf("%s: %s-toc\n", version, version))...)
		return os.WriteFile(tocMappingPath, b, 0644)
	}
	return nil
}

func updateHugoConfig(hugoConfigPath string, version string) error {
	b, err := os.ReadFile(hugoConfigPath)
	if err != nil {
		return err
	}
	var n yaml.Node
	if err := yaml.Unmarshal(b, &n); err != nil {
		return err
	}
	if n.Kind != yaml.DocumentNode {
		return fmt.Errorf("wrong Node Kind")
	}
	if len(n.Content) != 1 {
		return fmt.Errorf("expected a single YAML document")
	}
	nMap := n.Content[0]
	if nMap.Kind != yaml.MappingNode {
		return fmt.Errorf("wrong Node Kind")
	}
	var nParams *yaml.Node
	for idx, n := range nMap.Content {
		if n.Value == "params" {
			nParams = nMap.Content[idx+1]
			break
		}
	}
	if nParams == nil {
		return fmt.Errorf("'params' not found")
	}
	if nParams.Kind != yaml.MappingNode {
		return fmt.Errorf("wrong Node Kind for 'params'")
	}
	var nDocsVersions *yaml.Node
	for idx, n := range nParams.Content {
		if n.Value == "docs_versions" {
			nDocsVersions = nParams.Content[idx+1]
			break
		}
	}
	if nDocsVersions == nil {
		return fmt.Errorf("'docs_versions' not found")
	}
	if nDocsVersions.Kind != yaml.SequenceNode {
		return fmt.Errorf("wrong Node Kind for 'docs_versions'")
	}
	var docsVersions []string
	if err := nDocsVersions.Decode(&docsVersions); err != nil {
		return fmt.Errorf("error when decoding 'docs_versions' list: %w", err)
	}
	hasVersion := false
	for _, v := range docsVersions {
		if v == version {
			hasVersion = true
		}
	}

	if !hasVersion {
		docsVersions = append(docsVersions, version)
	}
	if docsVersions[0] != "main" {
		return fmt.Errorf("first element of 'docs_versions' list should be 'main'")
	}
	semvers := docsVersions[1:]
	sort.Slice(semvers, func(i, j int) bool {
		v1 := mustParseVersion(semvers[i])
		v2 := mustParseVersion(semvers[j])
		return !v1.LessThan(v2)
	})

	if hasVersion {
		log.Printf("Hugo config already includes version %s\n", version)
	} else {
		log.Printf("Adding version %s to Hugo config\n", version)

		if err := nDocsVersions.Encode(&docsVersions); err != nil {
			return fmt.Errorf("error when encoding updated 'docs_versions' list: %w", err)
		}
	}

	latestVersion := docsVersions[1]

	var nDocsLatest *yaml.Node
	for idx, n := range nParams.Content {
		if n.Value == "docs_latest" {
			nDocsLatest = nParams.Content[idx+1]
			break
		}
	}
	if nDocsLatest == nil {
		return fmt.Errorf("'docs_latest' not found")
	}
	if nDocsLatest.Kind != yaml.ScalarNode {
		return fmt.Errorf("wrong Node Kind for 'docs_latest'")
	}
	var docsLatest string
	if err := nDocsLatest.Decode(&docsLatest); err != nil {
		return fmt.Errorf("error when decoding 'docs_latest' scalar: %w", err)
	}

	if docsLatest == latestVersion {
		log.Printf("Hugo config already has the correct latest version %s\n", latestVersion)
	} else {
		log.Printf("Setting latest version in Hugo config to %s\n", latestVersion)

		if err := nDocsLatest.Encode(&latestVersion); err != nil {
			return fmt.Errorf("error when encoding updated 'docs_latest' scalar: %w", err)
		}
	}

	if pkg.DryRun {
		return nil
	}

	f, err := os.Create(hugoConfigPath)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	encoder := yaml.NewEncoder(w)
	defer encoder.Close()
	encoder.SetIndent(2)
	return encoder.Encode(&n)
}

func main() {
	flag.StringVar(&version, "version", "", "Version of the docs (must be a semver)")
	flag.Parse()

	if version == "" {
		log.Fatalf("flag -version is required")
	}

	semVer, err := parseVersion(version)
	if err != nil {
		log.Fatalf("Cannot parse semver: %v", err)
	}

	var referenceVersion string
	if semVer.IsNewMinorVersion() {
		referenceVersion = "main"
	} else {
		referenceVersion = ""
		v := semVer.PrevPatchVersion()
		for {
			docsPath := filepath.Join(pkg.WebsiteRepo, "content", "docs", v.String())
			if _, err := os.Stat(docsPath); err == nil {
				referenceVersion = v.String()
				break
			}
			if v.IsNewMinorVersion() {
				break
			}
			v = v.PrevPatchVersion()
		}
	}
	if referenceVersion == "" {
		log.Fatalf("Cannot determine reference version")
	}
	log.Printf("Reference version is %s\n", referenceVersion)

	destDocsPath := filepath.Join(pkg.WebsiteRepo, "content", "docs", version)
	if err := createDocsIfNeeded(destDocsPath, version, referenceVersion); err != nil {
		log.Fatalf("Failed to initialize docs: %v", err)
	}

	tocFile := filepath.Join(pkg.WebsiteRepo, "data", "docs", fmt.Sprintf("%s-toc.yml", version))
	if err := createTocFileIfNeeded(tocFile, referenceVersion); err != nil {
		log.Fatalf("Failed to initialize TOC file: %v", err)
	}

	tocMappingPath := filepath.Join(pkg.WebsiteRepo, "data", "docs", "toc-mapping.yml")
	if err := updateTocMapping(tocMappingPath, version); err != nil {
		log.Fatalf("Failed to update TOC mapping: %v", err)
	}

	if !semVer.HasSuffix() {
		// for pre-releases (Alpha / Beta / Release Candidates), we hide the generated docs,
		// by skipping the config.yaml update.
		hugoConfigPath := filepath.Join(pkg.WebsiteRepo, "config.yaml")
		if err := updateHugoConfig(hugoConfigPath, version); err != nil {
			log.Fatalf("Failed to update Hugo config: %v", err)
		}
	}

	if err := pkg.UpdateDocs(destDocsPath, version); err != nil {
		log.Fatalf("Failed to update docs: %v", err)
	}
}
