package main

type engineIDError struct {
	s string
}

func (e *engineIDError) Error() string {
	return "Error parsing engine ID: " + e.s
}
func New(text string) error {
	return &engineIDError{text}
}
