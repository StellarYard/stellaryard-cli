package formatter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestFormatTableRendersHeadersAndRows(t *testing.T) {
	var buf bytes.Buffer
	err := Format(&buf, "table",
		[]string{"NAME", "STATE"},
		[][]string{
			{"horizon", "running"},
			{"soroban-rpc", "exited"},
		},
	)
	if err != nil {
		t.Fatalf("Format(table) failed: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"NAME", "STATE", "horizon", "running", "soroban-rpc", "exited"} {
		if !strings.Contains(out, want) {
			t.Errorf("table output missing %q:\n%s", want, out)
		}
	}
}

func TestFormatTableEmptyRowsPrintsHeaderOnly(t *testing.T) {
	var buf bytes.Buffer
	if err := Format(&buf, "table", []string{"NAME", "STATE"}, nil); err != nil {
		t.Fatalf("Format(table) with no rows failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "NAME") {
		t.Errorf("header missing from empty table:\n%s", out)
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("expected header + separator (2 lines), got %d:\n%s", len(lines), out)
	}
}

func TestFormatTableAlignsColumns(t *testing.T) {
	var buf bytes.Buffer
	if err := Format(&buf, "table",
		[]string{"NAME", "STATE"},
		[][]string{
			{"horizon", "running"},
			{"soroban-rpc", "exited"},
		},
	); err != nil {
		t.Fatalf("Format(table) failed: %v", err)
	}

	// tabwriter pads to the widest cell, so the second column must start at
	// the same offset on every line.
	offsets := map[int]bool{}
	for _, line := range strings.Split(strings.TrimRight(buf.String(), "\n"), "\n") {
		if i := strings.Index(line, "STATE"); i >= 0 {
			offsets[i] = true
		}
		if i := strings.Index(line, "running"); i >= 0 {
			offsets[i] = true
		}
		if i := strings.Index(line, "exited"); i >= 0 {
			offsets[i] = true
		}
	}
	if len(offsets) != 1 {
		t.Errorf("second column is not aligned; offsets = %v\n%s", offsets, buf.String())
	}
}

func TestFormatJSONProducesArrayOfObjects(t *testing.T) {
	var buf bytes.Buffer
	err := Format(&buf, "json",
		[]string{"NAME", "STATE"},
		[][]string{
			{"horizon", "running"},
			{"soroban-rpc", "exited"},
		},
	)
	if err != nil {
		t.Fatalf("Format(json) failed: %v", err)
	}

	var got []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, buf.String())
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0]["NAME"] != "horizon" || got[0]["STATE"] != "running" {
		t.Errorf("first row = %v, want horizon/running", got[0])
	}
	if got[1]["NAME"] != "soroban-rpc" || got[1]["STATE"] != "exited" {
		t.Errorf("second row = %v, want soroban-rpc/exited", got[1])
	}
}

func TestFormatJSONRaggedRowDoesNotPanic(t *testing.T) {
	var buf bytes.Buffer
	err := Format(&buf, "json",
		[]string{"A", "B", "C"},
		[][]string{
			{"only-one", "only-two"}, // shorter than headers
			{"1", "2", "3", "4"},     // longer than headers
		},
	)
	if err != nil {
		t.Fatalf("Format(json) with ragged rows failed: %v", err)
	}

	var got []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, buf.String())
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if _, ok := got[0]["C"]; ok {
		t.Errorf("short row should omit the missing key, got %v", got[0])
	}
	if len(got[1]) != 3 {
		t.Errorf("long row should keep only the declared columns, got %v", got[1])
	}
}

func TestFormatJSONEmptyRowsProducesNullArray(t *testing.T) {
	var buf bytes.Buffer
	if err := Format(&buf, "json", []string{"A"}, nil); err != nil {
		t.Fatalf("Format(json) with no rows failed: %v", err)
	}
	if !json.Valid(buf.Bytes()) {
		t.Fatalf("output is not valid JSON: %s", buf.String())
	}
}

func TestFormatUnknownFormatReturnsError(t *testing.T) {
	var buf bytes.Buffer
	err := Format(&buf, "yaml", []string{"A"}, [][]string{{"1"}})
	if err == nil {
		t.Fatal("Format(yaml) succeeded, want an error")
	}
	if !strings.Contains(err.Error(), "unknown format") {
		t.Errorf("error = %q, want it to mention 'unknown format'", err.Error())
	}
	if !strings.Contains(err.Error(), "yaml") {
		t.Errorf("error = %q, want it to name the offending format", err.Error())
	}
}
