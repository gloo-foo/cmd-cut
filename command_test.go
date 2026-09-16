package command_test

import (
	"slices"
	"testing"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/testable"
	"github.com/gloo-foo/testable/run"

	command "github.com/gloo-foo/cmd-cut"
)

// lines executes cmd over input and returns the split output lines.
func lines(t *testing.T, cmd gloo.Command[[]byte, []byte], input string) []string {
	t.Helper()
	got, err := testable.TestLines(cmd, run.Input(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return got
}

// assertLines fails unless got equals want.
func assertLines(t *testing.T, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

// ==============================================================================
// Field selection (-f / -d)
// ==============================================================================

func TestCut_Fields_Basic(t *testing.T) {
	got := lines(t, command.Cut(command.CutDelimiter(","), command.CutFields(1, 3)), "a,b,c\nd,e,f\n")
	assertLines(t, got, []string{"a,c", "d,f"})
}

func TestCut_Fields_DefaultTabDelimiter(t *testing.T) {
	// With no -d, fields split on tab.
	got := lines(t, command.Cut(command.CutFields(2)), "one\ttwo\tthree\n")
	assertLines(t, got, []string{"two"})
}

func TestCut_Fields_SelectedInInputOrderNotRequestOrder(t *testing.T) {
	// cut(1) emits fields in input order regardless of the order requested:
	// `-f3,1` yields field 1 then field 3, NOT field 3 then field 1.
	got := lines(t, command.Cut(command.CutDelimiter(","), command.CutFields(3, 1)), "a,b,c\n")
	assertLines(t, got, []string{"a,c"})
}

func TestCut_Fields_MissingFieldsOmitted(t *testing.T) {
	// Requesting only an out-of-range field yields an empty line.
	got := lines(t, command.Cut(command.CutDelimiter(","), command.CutFields(5)), "a,b,c\n")
	assertLines(t, got, []string{""})
}

func TestCut_Fields_NoDelimiterPassesLineThrough(t *testing.T) {
	got := lines(t, command.Cut(command.CutDelimiter(","), command.CutFields(1)), "no-comma-here\n")
	assertLines(t, got, []string{"no-comma-here"})
}

func TestCut_Fields_NoFieldsRequestedPassesThrough(t *testing.T) {
	got := lines(t, command.Cut(command.CutDelimiter(",")), "a,b,c\n")
	assertLines(t, got, []string{"a,b,c"})
}

func TestCut_Fields_EmptyInputYieldsNoLines(t *testing.T) {
	got := lines(t, command.Cut(command.CutDelimiter(","), command.CutFields(1)), "")
	assertLines(t, got, nil)
}

func TestCut_Fields_OpenEndedRangeViaMultiple(t *testing.T) {
	// Selecting fields 2 and 3 explicitly keeps the trailing fields in order.
	got := lines(t, command.Cut(command.CutDelimiter(","), command.CutFields(2, 3, 4)), "a,b,c,d\n")
	assertLines(t, got, []string{"b,c,d"})
}

// ==============================================================================
// Byte selection (-b)
// ==============================================================================

func TestCut_Bytes_Range(t *testing.T) {
	got := lines(t, command.Cut(command.CutBytes("1-3")), "abcdef\n")
	assertLines(t, got, []string{"abc"})
}

func TestCut_Bytes_DiscretePositions(t *testing.T) {
	got := lines(t, command.Cut(command.CutBytes("1,3,5")), "abcdef\n")
	assertLines(t, got, []string{"ace"})
}

func TestCut_Bytes_MixedRangeAndPosition(t *testing.T) {
	got := lines(t, command.Cut(command.CutBytes("1-3,5")), "abcdef\n")
	assertLines(t, got, []string{"abce"})
}

func TestCut_Bytes_OpenEndedToEnd(t *testing.T) {
	// "2-" selects byte 2 through the end of the line.
	got := lines(t, command.Cut(command.CutBytes("2-")), "abcdef\n")
	assertLines(t, got, []string{"bcdef"})
}

func TestCut_Bytes_OpenEndedFromStart(t *testing.T) {
	// "-3" selects bytes 1 through 3.
	got := lines(t, command.Cut(command.CutBytes("-3")), "abcdef\n")
	assertLines(t, got, []string{"abc"})
}

func TestCut_Bytes_InvalidPartsSkipped(t *testing.T) {
	// A non-numeric position part is ignored; "x" drops out, "2" remains.
	got := lines(t, command.Cut(command.CutBytes("x,2")), "abcdef\n")
	assertLines(t, got, []string{"b"})
}

func TestCut_Bytes_InvalidRangeBoundsSkipped(t *testing.T) {
	// Both an invalid low bound and an invalid high bound drop the whole part.
	got := lines(t, command.Cut(command.CutBytes("x-3,2-y,4")), "abcdef\n")
	assertLines(t, got, []string{"d"})
}

func TestCut_Bytes_AllInvalidSelectsNothing(t *testing.T) {
	// A non-empty spec whose every part is invalid selects no positions and so
	// yields an empty line (distinct from an unset spec, which passes through).
	got := lines(t, command.Cut(command.CutBytes("x")), "abcdef\n")
	assertLines(t, got, []string{""})
}

// ==============================================================================
// Character selection (-c)
// ==============================================================================

func TestCut_Chars_Range(t *testing.T) {
	got := lines(t, command.Cut(command.CutChars("1-3")), "abcdef\n")
	assertLines(t, got, []string{"abc"})
}

func TestCut_Chars_Unicode(t *testing.T) {
	// Each Japanese character is 3 bytes but 1 rune; -c counts runes.
	got := lines(t, command.Cut(command.CutChars("1-2")), "日本語\n")
	assertLines(t, got, []string{"日本"})
}

func TestCut_Chars_OpenEndedToEnd(t *testing.T) {
	got := lines(t, command.Cut(command.CutChars("2-")), "日本語\n")
	assertLines(t, got, []string{"本語"})
}

// ==============================================================================
// Complement (--complement)
// ==============================================================================

func TestCut_Complement_Fields(t *testing.T) {
	got := lines(t, command.Cut(command.CutDelimiter(","), command.CutFields(2), command.CutComplement), "a,b,c\n")
	assertLines(t, got, []string{"a,c"})
}

func TestCut_Complement_Bytes(t *testing.T) {
	got := lines(t, command.Cut(command.CutBytes("2,4"), command.CutComplement), "abcdef\n")
	assertLines(t, got, []string{"acef"})
}

func TestCut_Complement_Chars(t *testing.T) {
	got := lines(t, command.Cut(command.CutChars("1,3"), command.CutComplement), "abcde\n")
	assertLines(t, got, []string{"bde"})
}

func TestCut_Complement_NoComplementDefault(t *testing.T) {
	// The explicit NoComplement form behaves like the default.
	got := lines(t, command.Cut(command.CutBytes("1,2"), command.CutNoComplement), "abcdef\n")
	assertLines(t, got, []string{"ab"})
}
