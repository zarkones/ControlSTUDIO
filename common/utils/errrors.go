package utils

// MaybeFatal handles errors of the main function only.
// It handles unrecoverable errors and for most cases you don't
// need to use this function.
func MaybeFatal(err error) {
	if err == nil {
		return
	}
	panic(err)
}
