package validator

type Option func(*Validator)

func WithDefaultValidators() Option {
	return func(v *Validator) {
		v.funcs = builtInFunctions()
	}
}

func WithStructTag(tag string) Option {
	return func(v *Validator) {
		v.structTag = tag
	}
}
