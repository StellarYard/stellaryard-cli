package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Format represents the output format mode.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

// Writer handles formatted output.
type Writer struct {
	format Format
	w      io.Writer
}

// New creates a new Writer with the given format.
func New(format string) *Writer {
	f := Format(format)
	if f != FormatJSON {
		f = FormatTable
	}
	return &Writer{format: f, w: os.Stdout}
}

// Write outputs data in the configured format.
// For JSON, marshals the data and writes it.
// For table, the caller should use WriteTable for tabular data.
func (w *Writer) Write(data interface{}) error {
	if w.format == FormatJSON {
		enc := json.NewEncoder(w.w)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	}
	// Table format — caller should use WriteTable instead
	fmt.Fprintf(w.w, "%v\n", data)
	return nil
}

// WriteTable outputs tabular data using tabwriter.
func (w *Writer) WriteTable(headers []string, rows [][]string) error {
	if w.format == FormatJSON {
		// Convert to map for JSON output
		var results []map[string]string
		for _, row := range rows {
			m := make(map[string]string)
			for i, h := range headers {
				if i < len(row) {
					m[h] = row[i]
				}
			}
			results = append(results, m)
		}
		return w.Write(results)
	}

	tw := tabwriter.NewWriter(w.w, 0, 0, 2, ' ', 0)
	for i, h := range headers {
		if i > 0 {
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprint(tw, h)
	}
	fmt.Fprintln(tw)

	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				fmt.Fprint(tw, "\t")
			}
			fmt.Fprint(tw, cell)
		}
		fmt.Fprintln(tw)
	}

	return tw.Flush()
}

// WriteError outputs an error in the configured format.
// This is critical for CI — stderr must respect --format json.
func (w *Writer) WriteError(err error) error {
	if w.format == FormatJSON {
		return json.NewEncoder(w.w).Encode(map[string]string{
			"error": err.Error(),
		})
	}
	fmt.Fprintf(w.w, "error: %v\n", err)
	return nil
}
