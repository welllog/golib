package strz

import (
	"bytes"
	"io"
	"testing"
)

func TestReader_Read(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		size     int
		expected string
	}{
		{
			name:     "read full string",
			input:    "hello world",
			size:     11,
			expected: "hello world",
		},
		{
			name:     "read partial string",
			input:    "hello world",
			size:     5,
			expected: "hello",
		},
		{
			name:     "read empty string",
			input:    "",
			size:     1,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.input)
			buf := make([]byte, tt.size)
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := string(buf[:n]); got != tt.expected {
				t.Errorf("unexpected result: got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestReader_ReadAt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		offset   int64
		size     int
		expected string
	}{
		{
			name:     "read full string",
			input:    "hello world",
			offset:   0,
			size:     11,
			expected: "hello world",
		},
		{
			name:     "read partial string",
			input:    "hello world",
			offset:   6,
			size:     5,
			expected: "world",
		},
		{
			name:     "read empty string",
			input:    "",
			offset:   0,
			size:     1,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.input)
			buf := make([]byte, tt.size)
			n, err := r.ReadAt(buf, tt.offset)
			if err != nil && err != io.EOF {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := string(buf[:n]); got != tt.expected {
				t.Errorf("unexpected result: got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestReader_ReadByte(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected byte
	}{
		{
			name:     "read first byte",
			input:    "hello world",
			expected: 'h',
		},
		{
			name:     "read empty string",
			input:    "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.input)
			b, err := r.ReadByte()
			if err != nil && err != io.EOF {
				t.Fatalf("unexpected error: %v", err)
			}
			if b != tt.expected {
				t.Errorf("unexpected result: got %q, want %q", b, tt.expected)
			}
		})
	}
}

func TestReader_UnreadByte(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected byte
	}{
		{
			name:     "unread first byte",
			input:    "hello world",
			expected: 'h',
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.input)
			b, err := r.ReadByte()
			if err != nil && err != io.EOF {
				t.Fatalf("unexpected error: %v", err)
			}
			if b != tt.expected {
				t.Errorf("unexpected result: got %q, want %q", b, tt.expected)
			}
			err = r.UnreadByte()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			b, err = r.ReadByte()
			if err != nil && err != io.EOF {
				t.Fatalf("unexpected error: %v", err)
			}
			if b != tt.expected {
				t.Errorf("unexpected result: got %q, want %q", b, tt.expected)
			}
		})
	}
}

func TestReader_WriteTo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "write full string",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "write empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.input)
			var buf bytes.Buffer
			_, err := r.WriteTo(&buf)
			if err != nil && err != io.EOF {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := buf.String(); got != tt.expected {
				t.Errorf("unexpected result: got %q, want %q", got, tt.expected)
			}
		})
	}
}
