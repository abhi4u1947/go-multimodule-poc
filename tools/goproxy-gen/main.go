// Command goproxy-gen builds a static GOPROXY-compatible directory tree for
// every Go module found in this repository, using each module's real Git
// tags as its version history.
//
// The output layout follows the Go module proxy protocol
// (https://go.dev/ref/mod#goproxy-protocol) exactly:
//
//	<escaped-module-path>/@v/list             newline-separated known versions
//	<escaped-module-path>/@v/<version>.info   {"Version":"vX.Y.Z","Time":"..."}
//	<escaped-module-path>/@v/<version>.mod    that version's go.mod content
//	<escaped-module-path>/@v/<version>.zip    that version's module zip
//	<escaped-module-path>/@latest             .info content of the newest version
//
// Point GOPROXY at the result (e.g. `GOPROXY=file:///abs/path/to/dist/goproxy`)
// and `go get`/`go mod download` will resolve this repo's modules from it
// exactly as they would from a real proxy, with no network access and no
// dependency on GitHub's tag-push or module-proxy caching behavior.
//
// This is not a general-purpose proxy server: it's a generator that snapshots
// whatever is currently tagged in the repo into a directory a plain static
// file server (or GOPROXY's own file:// support) can serve as-is.
package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
	xzip "golang.org/x/mod/zip"
)

func main() {
	repoFlag := flag.String("repo", ".", "path to the Git repository root")
	outFlag := flag.String("out", "dist/goproxy", "output directory for the generated proxy tree")
	flag.Parse()

	if err := run(*repoFlag, *outFlag); err != nil {
		log.Fatalf("goproxy-gen: %v", err)
	}
}

func run(repoArg, outArg string) error {
	repoRoot, err := filepath.Abs(repoArg)
	if err != nil {
		return fmt.Errorf("resolving repo path: %w", err)
	}
	outRoot, err := filepath.Abs(outArg)
	if err != nil {
		return fmt.Errorf("resolving output path: %w", err)
	}

	mods, err := discoverModules(repoRoot)
	if err != nil {
		return fmt.Errorf("discovering modules: %w", err)
	}
	if len(mods) == 0 {
		return fmt.Errorf("no go.mod files found under %s", repoRoot)
	}

	for _, m := range mods {
		if err := generateModule(repoRoot, outRoot, m); err != nil {
			return fmt.Errorf("module %s: %w", m.path, err)
		}
	}
	return nil
}

// module describes one Go module found in the repository.
type modInfo struct {
	path string // module import path, from its go.mod
	dir  string // slash-separated path relative to repoRoot ("" for the repo root module)
}

// discoverModules walks the repository for go.mod files, skipping this
// tool's own module (tools/) and version control metadata.
func discoverModules(repoRoot string) ([]modInfo, error) {
	var mods []modInfo
	err := filepath.WalkDir(repoRoot, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "tools", "dist":
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() != "go.mod" {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		f, err := modfile.Parse(p, data, nil)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", p, err)
		}
		if f.Module == nil {
			return nil
		}
		rel, err := filepath.Rel(repoRoot, filepath.Dir(p))
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "." {
			rel = ""
		}
		mods = append(mods, modInfo{path: f.Module.Mod.Path, dir: rel})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(mods, func(i, j int) bool { return mods[i].path < mods[j].path })
	return mods, nil
}

// versionTagPattern matches only exact, unsuffixed semver tags for a given
// directory prefix (e.g. "entities/shared-lib/v1.1.0"), rejecting pseudo
// tags like the repo's own intentionally-broken "-mismatch" experiment tag.
func versionTagPattern(dir string) *regexp.Regexp {
	prefix := ""
	if dir != "" {
		prefix = regexp.QuoteMeta(dir) + "/"
	}
	return regexp.MustCompile(`^` + prefix + `v[0-9]+\.[0-9]+\.[0-9]+$`)
}

func generateModule(repoRoot, outRoot string, m modInfo) error {
	versions, err := taggedVersions(repoRoot, m.dir)
	if err != nil {
		return fmt.Errorf("listing tags: %w", err)
	}
	if len(versions) == 0 {
		log.Printf("goproxy-gen: %s: no semver tags found, skipping", m.path)
		return nil
	}

	escapedPath, err := module.EscapePath(m.path)
	if err != nil {
		return fmt.Errorf("escaping module path: %w", err)
	}
	vDir := filepath.Join(outRoot, filepath.FromSlash(escapedPath), "@v")
	if err := os.MkdirAll(vDir, 0o755); err != nil {
		return err
	}

	var listLines []string
	var latestInfo []byte

	for _, version := range versions {
		tag := version
		if m.dir != "" {
			tag = m.dir + "/" + version
		}

		escapedVersion, err := module.EscapeVersion(version)
		if err != nil {
			return fmt.Errorf("escaping version %s: %w", version, err)
		}

		zipBuf := &bytes.Buffer{}
		if err := xzip.CreateFromVCS(zipBuf, module.Version{Path: m.path, Version: version}, repoRoot, tag, m.dir); err != nil {
			return fmt.Errorf("zipping %s: %w", tag, err)
		}

		goModData, err := readZipFile(zipBuf.Bytes(), m.path+"@"+version+"/go.mod")
		if err != nil {
			return fmt.Errorf("reading go.mod from %s zip: %w", tag, err)
		}

		commitTime, err := tagCommitTime(repoRoot, tag)
		if err != nil {
			return fmt.Errorf("getting commit time for %s: %w", tag, err)
		}

		info, err := json.Marshal(struct {
			Version string
			Time    string
		}{Version: version, Time: commitTime.UTC().Format(time.RFC3339)})
		if err != nil {
			return err
		}

		if err := os.WriteFile(filepath.Join(vDir, escapedVersion+".info"), info, 0o644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(vDir, escapedVersion+".mod"), goModData, 0o644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(vDir, escapedVersion+".zip"), zipBuf.Bytes(), 0o644); err != nil {
			return err
		}

		listLines = append(listLines, version)
		latestInfo = info // versions are processed in ascending order; last write wins

		log.Printf("goproxy-gen: %s@%s -> %s", m.path, version, vDir)
	}

	listContent := strings.Join(listLines, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(vDir, "list"), []byte(listContent), 0o644); err != nil {
		return err
	}
	latestPath := filepath.Join(outRoot, filepath.FromSlash(escapedPath), "@latest")
	if err := os.WriteFile(latestPath, latestInfo, 0o644); err != nil {
		return err
	}
	return nil
}

// taggedVersions returns every exact semver version tagged for the module
// rooted at dir (relative to the repo root, "" for the repo-root module),
// sorted oldest to newest.
func taggedVersions(repoRoot, dir string) ([]string, error) {
	out, err := runGit(repoRoot, "tag", "--list")
	if err != nil {
		return nil, err
	}
	pattern := versionTagPattern(dir)
	prefixLen := 0
	if dir != "" {
		prefixLen = len(dir) + 1 // + "/"
	}

	var versions []string
	for _, tag := range strings.Split(strings.TrimSpace(out), "\n") {
		if tag == "" || !pattern.MatchString(tag) {
			continue
		}
		v := tag[prefixLen:]
		if !semver.IsValid(v) {
			continue
		}
		versions = append(versions, v)
	}
	sort.Slice(versions, func(i, j int) bool { return semver.Compare(versions[i], versions[j]) < 0 })
	return versions, nil
}

func tagCommitTime(repoRoot, tag string) (time.Time, error) {
	out, err := runGit(repoRoot, "log", "-1", "--format=%cI", tag)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, strings.TrimSpace(out))
}

func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, stderr.String())
	}
	return stdout.String(), nil
}

func readZipFile(zipData []byte, name string) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, err
	}
	for _, f := range zr.File {
		if path.Clean(f.Name) != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		return io.ReadAll(rc)
	}
	return nil, fmt.Errorf("%s not found in zip", name)
}
