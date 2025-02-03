package fileio

import "io"

// Close the given closable reader or writer, ignoring any close errors that may arise.
// This wrapper is essential in defer statements,
// when we actually don't need to care whether the file was closed properly
func Close(closable io.Closer) {
	_ = closable.Close()
}
