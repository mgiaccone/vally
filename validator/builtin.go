package validator

func builtInFunctions() map[string]Function {
	return map[string]Function{
		"false":    False,
		"true":     True,
		"eq":       Equal,
		"required": Required,
		"email":    Email,
		"isbn10":   ISBN10,
		"isbn13":   ISBN13,
	}
}
