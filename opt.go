package command

// CutDelimiter sets the field delimiter (-d flag). Default is tab.
type CutDelimiter string

// CutField is a 1-based field position selected by the -f flag.
type CutField int

// CutFieldsOpt selects fields by 1-based position (-f flag).
type CutFieldsOpt []CutField

// CutFields returns a fields option for Cut.
func CutFields(f ...CutField) CutFieldsOpt { return CutFieldsOpt(f) }

// CutBytes selects bytes by position (-b flag).
// The spec string uses 1-based positions and ranges: "1-3,5,7-9".
type CutBytes string

// CutChars selects characters (runes) by position (-c flag).
// The spec string uses 1-based positions and ranges: "1-3,5,7-9".
type CutChars string

// cutComplementFlag inverts the selection (--complement flag).
type cutComplementFlag bool

const (
	// CutComplement inverts the selection: emit positions NOT in the set.
	CutComplement cutComplementFlag = true
	// CutNoComplement is the default: emit positions in the set.
	CutNoComplement cutComplementFlag = false
)
