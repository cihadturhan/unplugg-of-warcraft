package warcraft

// Error represents an error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }

// General errors.
const (
	ErrGetDump    Error = "failed to retrieve api dump"
	ErrDumpExists Error = "API dump already exists"
)
