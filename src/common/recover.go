package common

// HandlePanics function is a shorthand for recovering from panics and putting handled
// error values into provided error pointer
// Example usage:
//
func HandlePanics(err *error, handle func(recovered interface{}) error) {
	if rec := recover(); rec != nil {
		*err = handle(rec)
	}
}
