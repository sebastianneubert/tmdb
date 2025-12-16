package display_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/sebastianneubert/tmdb/internal/display"
)

// captureOutput captures stdout during a function call
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintSearchStartMessage(t *testing.T) {
	output := captureOutput(func() {
		display.PrintSearchStartMessage("Test Movies", 7.5, 500, "Netflix,Amazon", "US")
	})

	if !strings.Contains(output, "Test Movies") {
		t.Error("Output should contain search type")
	}
	if !strings.Contains(output, "7.5") {
		t.Error("Output should contain min rating")
	}
	if !strings.Contains(output, "500") {
		t.Error("Output should contain min votes")
	}
	if !strings.Contains(output, "Netflix,Amazon") {
		t.Error("Output should contain providers")
	}
	if !strings.Contains(output, "US") {
		t.Error("Output should contain region")
	}
}

func TestPrintSearchResultsSummaryWithResults(t *testing.T) {
	output := captureOutput(func() {
		display.PrintSearchResultsSummary("top-rated movies", 5)
	})

	if !strings.Contains(output, "5") {
		t.Error("Output should contain result count")
	}
	if !strings.Contains(output, "top-rated movies") {
		t.Error("Output should contain search type")
	}
	if strings.Contains(output, "No") && strings.Contains(output, "found") {
		t.Error("Output should not say 'No results' when results > 0")
	}
}

func TestPrintSearchResultsSummaryNoResults(t *testing.T) {
	output := captureOutput(func() {
		display.PrintSearchResultsSummary("popular movies", 0)
	})

	if !strings.Contains(output, "No") {
		t.Error("Output should contain 'No'")
	}
	if !strings.Contains(output, "popular movies") {
		t.Error("Output should contain search type")
	}
}

func TestPrintSearchNoResults(t *testing.T) {
	output := captureOutput(func() {
		display.PrintSearchNoResults("Matrix", 15, 7.0, 500)
	})

	if !strings.Contains(output, "Matrix") {
		t.Error("Output should contain query")
	}
	if !strings.Contains(output, "15") {
		t.Error("Output should contain movies checked count")
	}
	if !strings.Contains(output, "7.0") {
		t.Error("Output should contain min rating")
	}
	if !strings.Contains(output, "500") {
		t.Error("Output should contain min votes")
	}
	if !strings.Contains(output, "--min-rating") {
		t.Error("Output should contain suggestion about min-rating flag")
	}
	if !strings.Contains(output, "--min-votes") {
		t.Error("Output should contain suggestion about min-votes flag")
	}
}

func TestPrintSearchCompleteMessage(t *testing.T) {
	output := captureOutput(func() {
		display.PrintSearchCompleteMessage(3, 20)
	})

	if !strings.Contains(output, "3") {
		t.Error("Output should contain results found count")
	}
	if !strings.Contains(output, "20") {
		t.Error("Output should contain total movies checked count")
	}
	if !strings.Contains(output, "complete") {
		t.Error("Output should indicate search is complete")
	}
}

func TestPrintSearchCompleteMessageSingleResult(t *testing.T) {
	output := captureOutput(func() {
		display.PrintSearchCompleteMessage(1, 1)
	})

	if !strings.Contains(output, "1") {
		t.Error("Output should contain count 1")
	}
}

func TestPrintSearchCompleteMessageManyResults(t *testing.T) {
	output := captureOutput(func() {
		display.PrintSearchCompleteMessage(100, 500)
	})

	if !strings.Contains(output, "100") {
		t.Error("Output should contain results count")
	}
	if !strings.Contains(output, "500") {
		t.Error("Output should contain checked count")
	}
}
