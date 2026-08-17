package osutil

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestWalkFilesSecure(t *testing.T) {
	dir := t.TempDir()

	mustWrite := func(rel, content string) {
		full := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("MkdirAll(%q): %v", filepath.Dir(full), err)
		}
		if err := os.WriteFile(full, []byte(content), 0o600); err != nil {
			t.Fatalf("WriteFile(%q): %v", full, err)
		}
	}

	mustWrite("a.go", "package a")
	mustWrite("sub/b.go", "package sub")
	mustWrite("vendor/c.go", "package vendor")
	mustWrite("testdata/d.go", "package testdata")

	var gotPaths []string
	err := WalkFilesSecure(dir, func(name string) bool {
		return name == "vendor" || name == "testdata"
	}, func(path string, content []byte) error {
		gotPaths = append(gotPaths, path)
		return nil
	})
	if err != nil {
		t.Fatalf("WalkFilesSecure() error = %v", err)
	}

	sort.Strings(gotPaths)
	want := []string{"a.go", "sub/b.go"}
	if len(gotPaths) != len(want) {
		t.Fatalf("WalkFilesSecure() paths = %v, want %v", gotPaths, want)
	}
	for i, p := range want {
		if gotPaths[i] != p {
			t.Errorf("WalkFilesSecure() paths[%d] = %q, want %q", i, gotPaths[i], p)
		}
	}
}

func TestWalkFilesSecure_NilSkipDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "x.txt"), []byte("hi"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var got []byte
	err := WalkFilesSecure(dir, nil, func(path string, content []byte) error {
		if path == "x.txt" {
			got = content
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkFilesSecure() error = %v", err)
	}
	if string(got) != "hi" {
		t.Errorf("WalkFilesSecure() content = %q, want %q", got, "hi")
	}
}

func TestReadDirFilesSecure(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "x.txt"), []byte("hi"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sub", "y.txt"), []byte("bye"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	files, err := ReadDirFilesSecure(dir)
	if err != nil {
		t.Fatalf("ReadDirFilesSecure() error = %v", err)
	}
	if string(files["x.txt"]) != "hi" {
		t.Errorf("files[x.txt] = %q, want %q", files["x.txt"], "hi")
	}
	if string(files["sub/y.txt"]) != "bye" {
		t.Errorf("files[sub/y.txt] = %q, want %q", files["sub/y.txt"], "bye")
	}
}
