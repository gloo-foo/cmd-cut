// Package alias provides short names for cut command flags.
package alias

import (
	gloo "github.com/gloo-foo/framework"

	command "github.com/gloo-foo/cmd-cut"
)

// Cut selects fields, bytes, or characters from each input line; see the
// command package for the flag set.
func Cut(opts ...any) gloo.Command[[]byte, []byte] { return command.Cut(opts...) }

// Delimiter sets the field delimiter (-d flag).
type Delimiter = command.CutDelimiter

// Fields selects fields by 1-based position (-f flag).
func Fields(f ...command.CutField) command.CutFieldsOpt { return command.CutFields(f...) }

// Bytes selects bytes by position (-b flag).
type Bytes = command.CutBytes

// Chars selects characters (runes) by position (-c flag).
type Chars = command.CutChars

// Complement inverts the selection (--complement flag).
const Complement = command.CutComplement

// NoComplement disables complement (default).
const NoComplement = command.CutNoComplement
