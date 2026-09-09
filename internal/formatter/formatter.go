package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// Format outputs data in the specified format.
func Format(w io.Writer, format string, headers []string, rows [][]string) error {
	switch format {
	case "json":
		return formatJSON(w, headers, rows)
	case "table":
		return formatTable(w, headers, rows)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func formatTable(w io.Writer, headers []string, rows [][]string) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	// Header
	for i, h := range headers {
		if i > 0 {
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprint(tw, h)
	}
	fmt.Fprintln(tw)

	// Separator
	for i := range headers {
		if i > 0 {
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprint(tw, "----")
	}
	fmt.Fprintln(tw)

	// Rows
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

func formatJSON(w io.Writer, headers []string, rows [][]string) error {
	// Convert to map format
	var result []map[string]string
	for _, row := range rows {
		m := make(map[string]string)
		for i, h := range headers {
			if i < len(row) {
				m[h] = row[i]
			}
		}
		result = append(result, m)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
