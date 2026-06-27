package command

// cutDelimiterOpt sets the field delimiter (-d flag). Default is tab.
type cutDelimiterOpt string

// CutDelimiter returns a delimiter option for Cut.
func CutDelimiter(d string) cutDelimiterOpt { return cutDelimiterOpt(d) }

// cutFieldsOpt selects fields by 1-based position (-f flag).
type cutFieldsOpt []int

// CutFields returns a fields option for Cut.
func CutFields(f ...int) cutFieldsOpt { return cutFieldsOpt(f) }

// cutBytesOpt selects bytes by position (-b flag).
// The spec string uses 1-based positions and ranges: "1-3,5,7-9".
type cutBytesOpt string

// CutBytes returns a byte-position option for Cut.
func CutBytes(spec string) cutBytesOpt { return cutBytesOpt(spec) }

// cutCharsOpt selects characters (runes) by position (-c flag).
// The spec string uses 1-based positions and ranges: "1-3,5,7-9".
type cutCharsOpt string

// CutChars returns a character-position option for Cut.
func CutChars(spec string) cutCharsOpt { return cutCharsOpt(spec) }

// cutComplementFlag inverts the selection (--complement flag).
type cutComplementFlag bool

const (
	// CutComplement inverts the selection: emit positions NOT in the set.
	CutComplement cutComplementFlag = true
	// CutNoComplement is the default: emit positions in the set.
	CutNoComplement cutComplementFlag = false
)
