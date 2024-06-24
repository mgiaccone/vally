package errutil

// Must forces a panic if the given error is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
