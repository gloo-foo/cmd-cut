// Package alias provides short names for cut command flags.
package alias

import command "github.com/gloo-foo/cmd-cut"

// Cut is the cut command constructor.
var Cut = command.Cut

// Delimiter sets the field delimiter (-d flag).
var Delimiter = command.CutDelimiter

// Fields selects fields by 1-based position (-f flag).
var Fields = command.CutFields

// Bytes selects bytes by position (-b flag).
var Bytes = command.CutBytes

// Chars selects characters (runes) by position (-c flag).
var Chars = command.CutChars

// Complement inverts the selection (--complement flag).
const Complement = command.CutComplement

// NoComplement disables complement (default).
const NoComplement = command.CutNoComplement
