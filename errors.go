package warcraft

// Error represents an error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
