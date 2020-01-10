package utils

import "io"

// M is a type representing a map
type M map[string]interface{}

// CloseTheCloser closes the closer :P
func CloseTheCloser(c io.Closer) {
	_ = c.Close()
}
