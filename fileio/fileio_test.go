package fileio

import (
	"testing"
)

// MockReadWriteCloser implements io.ReadWriteCloser and allows us to control the Close() behavior for testing
type MockReadWriteCloser struct {
	closed     bool
	closeError error
}

func (m *MockReadWriteCloser) Read([]byte) (n int, err error) {
	return 0, nil
}

func (m *MockReadWriteCloser) Write([]byte) (n int, err error) {
	return 0, nil
}

func (m *MockReadWriteCloser) Close() error {
	m.closed = true
	return m.closeError
}

func TestClose(t *testing.T) {
	m := &MockReadWriteCloser{}
	Close(m)
	if !m.closed {
		t.Error("Expected Close() to be called on the ReadWriteCloser")
	}
}
