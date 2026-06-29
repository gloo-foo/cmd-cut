package command

// CutDelimiterOpt sets the field delimiter (-d flag). Default is tab.
type CutDelimiterOpt string

// CutDelimiter returns a delimiter option for Cut.
func CutDelimiter(d string) CutDelimiterOpt { return CutDelimiterOpt(d) }

// CutFieldsOpt selects fields by 1-based position (-f flag).
type CutFieldsOpt []int

// CutFields returns a fields option for Cut.
func CutFields(f ...int) CutFieldsOpt { return CutFieldsOpt(f) }

// CutBytesOpt selects bytes by position (-b flag).
// The spec string uses 1-based positions and ranges: "1-3,5,7-9".
type CutBytesOpt string

// CutBytes returns a byte-position option for Cut.
func CutBytes(spec string) CutBytesOpt { return CutBytesOpt(spec) }

// CutCharsOpt selects characters (runes) by position (-c flag).
// The spec string uses 1-based positions and ranges: "1-3,5,7-9".
type CutCharsOpt string

// CutChars returns a character-position option for Cut.
func CutChars(spec string) CutCharsOpt { return CutCharsOpt(spec) }

// cutComplementFlag inverts the selection (--complement flag).
type cutComplementFlag bool

const (
	// CutComplement inverts the selection: emit positions NOT in the set.
	CutComplement cutComplementFlag = true
	// CutNoComplement is the default: emit positions in the set.
	CutNoComplement cutComplementFlag = false
)
