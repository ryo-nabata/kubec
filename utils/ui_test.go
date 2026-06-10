package utils

import (
	"bytes"
	"testing"

	"github.com/chzyer/readline"
)

func TestNoBellStdoutSuppressesSingleBell(t *testing.T) {
	buf := &bytes.Buffer{}
	w := &noBellStdout{w: buf}

	n, err := w.Write([]byte{readline.CharBell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("n = %d, want 1", n)
	}
	if got := buf.String(); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestNoBellStdoutPassesThroughNormalOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	w := &noBellStdout{w: buf}

	input := "Select a context"
	n, err := w.Write([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(input) {
		t.Errorf("n = %d, want %d", n, len(input))
	}
	if got := buf.String(); got != input {
		t.Errorf("got %q, want %q", got, input)
	}
}

func TestNoBellStdoutPassesThroughMultiByteWithBell(t *testing.T) {
	buf := &bytes.Buffer{}
	w := &noBellStdout{w: buf}

	input := []byte{'a', readline.CharBell, 'b'}
	n, err := w.Write(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(input) {
		t.Errorf("n = %d, want %d", n, len(input))
	}
	if got := buf.Bytes(); !bytes.Equal(got, input) {
		t.Errorf("got %q, want %q", got, input)
	}
}

func TestNoBellStdoutCloseIsNoop(t *testing.T) {
	w := &noBellStdout{w: &bytes.Buffer{}}
	if err := w.Close(); err != nil {
		t.Errorf("Close() = %v, want nil", err)
	}
}
