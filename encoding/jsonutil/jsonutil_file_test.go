package jsonutil

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type fileTestDoc struct {
	Foo []string `json:"foo"`
}

func TestMarshalFileRoundTrip(t *testing.T) {
	name := filepath.Join(t.TempDir(), "doc.json")
	in := fileTestDoc{Foo: []string{"bar", "baz"}}
	if err := MarshalFile(name, in, "", "  ", 0o600); err != nil {
		t.Fatal(err)
	}

	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	want := "{\n  \"foo\": [\n    \"bar\",\n    \"baz\"\n  ]\n}\n"
	if string(b) != want {
		t.Errorf("MarshalFile() wrote %q, want %q", b, want)
	}

	var got fileTestDoc
	if err := UnmarshalFile(name, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, in) {
		t.Errorf("UnmarshalFile() = %+v, want %+v", got, in)
	}

	var got2 fileTestDoc
	raw, err := UnmarshalFileWithBytes(name, &got2)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != want || !reflect.DeepEqual(got2, in) {
		t.Errorf("UnmarshalFileWithBytes() = (%q, %+v), want (%q, %+v)", raw, got2, want, in)
	}
}

func TestMarshalFilePrefix(t *testing.T) {
	name := filepath.Join(t.TempDir(), "doc.json")
	if err := MarshalFile(name, fileTestDoc{Foo: []string{"a"}}, "//", "\t", 0o600); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	if want := "{\n//\t\"foo\": [\n//\t\t\"a\"\n//\t]\n//}\n"; string(b) != want {
		t.Errorf("MarshalFile() wrote %q, want %q", b, want)
	}
}

func TestMarshalFileTruncates(t *testing.T) {
	name := filepath.Join(t.TempDir(), "doc.json")
	if err := MarshalFile(name, fileTestDoc{Foo: []string{"a long value that leaves trailing bytes"}}, "", "", 0o600); err != nil {
		t.Fatal(err)
	}
	if err := MarshalFile(name, fileTestDoc{}, "", "", 0o600); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	if want := "{\"foo\":null}\n"; string(b) != want {
		t.Errorf("after rewrite file = %q, want %q", b, want)
	}
}

func TestMarshalFileMissingDir(t *testing.T) {
	name := filepath.Join(t.TempDir(), "missing", "doc.json")
	if err := MarshalFile(name, fileTestDoc{}, "", "", 0o600); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("MarshalFile() into missing dir error = %v, want fs.ErrNotExist", err)
	}
}

func TestUnmarshalFileErrors(t *testing.T) {
	dir := t.TempDir()
	var doc fileTestDoc

	missing := filepath.Join(dir, "missing.json")
	if err := UnmarshalFile(missing, &doc); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("UnmarshalFile(missing) error = %v, want fs.ErrNotExist", err)
	}
	if _, err := UnmarshalFileWithBytes(missing, &doc); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("UnmarshalFileWithBytes(missing) error = %v, want fs.ErrNotExist", err)
	}

	invalid := filepath.Join(dir, "invalid.json")
	if err := os.WriteFile(invalid, []byte(`{"foo": [`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := UnmarshalFile(invalid, &doc); err == nil {
		t.Error("UnmarshalFile(invalid) error = nil, want error")
	}
	raw, err := UnmarshalFileWithBytes(invalid, &doc)
	if err == nil {
		t.Error("UnmarshalFileWithBytes(invalid) error = nil, want error")
	}
	if string(raw) != `{"foo": [` {
		t.Errorf("UnmarshalFileWithBytes(invalid) bytes = %q, want the file contents", raw)
	}
}
