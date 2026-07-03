package alias_test

import (
	"slices"
	"testing"

	"github.com/gloo-foo/testable"

	cut "github.com/gloo-foo/cmd-cut/alias"
)

// The alias package re-exports the constructor and flag helpers under short
// names. A mis-wired re-export (say, Bytes bound to Chars, or Complement bound
// to the disabled constant) compiles cleanly, so only behaviour can prove the
// wiring. Each test exercises one re-export and asserts the cut output it must
// produce.

func assertLines(t *testing.T, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAlias_FieldsWithDelimiter(t *testing.T) {
	// Delimiter (-d) plus Fields (-f): select fields 1 and 3.
	got, err := testable.TestLines(cut.Cut(cut.Delimiter(","), cut.Fields(1, 3)), "a,b,c\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, got, []string{"a,c"})
}

func TestAlias_Bytes(t *testing.T) {
	// Bytes (-b) selects byte positions, not runes.
	got, err := testable.TestLines(cut.Cut(cut.Bytes("1,3,5")), "abcdef\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, got, []string{"ace"})
}

func TestAlias_Chars(t *testing.T) {
	// Chars (-c) selects runes: two multi-byte characters, by position.
	got, err := testable.TestLines(cut.Cut(cut.Chars("1-2")), "日本語\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, got, []string{"日本"})
}

func TestAlias_Complement(t *testing.T) {
	// Complement (--complement) inverts the byte selection.
	got, err := testable.TestLines(cut.Cut(cut.Bytes("2,4"), cut.Complement), "abcdef\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, got, []string{"acef"})
}

func TestAlias_NoComplement(t *testing.T) {
	// NoComplement is the disabled form: it must behave like the default.
	got, err := testable.TestLines(cut.Cut(cut.Bytes("2,4"), cut.NoComplement), "abcdef\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, got, []string{"bd"})
}
