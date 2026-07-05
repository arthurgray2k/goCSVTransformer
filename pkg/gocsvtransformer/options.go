package gocsvtransformer

// Options contains configuration for transformations.
type Options struct {
	// Delimiter specifies the character used to separate fields in CSV.
	// Defaults to comma (',').
	Delimiter rune

	// HeaderRow specifies whether the first row of the CSV is a header.
	// When true, CSVToJSON uses headers as keys. When false, generates Column1, Column2, etc.
	HeaderRow bool

	// Pretty specifies whether JSON output should be formatted with indentation.
	Pretty bool

	// Indent specifies the number of spaces for indentation if Pretty is true.
	Indent int
}

// DefaultOptions returns the standard, sensible defaults for transformations.
func DefaultOptions() Options {
	return Options{
		Delimiter: ',',
		HeaderRow: true,
		Pretty:    false,
		Indent:    2,
	}
}
