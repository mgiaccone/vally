package validator

type Option func(*Validator)

func WithDefaultValidators() Option {
	return func(v *Validator) {
		v.funcs = defaultFuncs()
	}
}
