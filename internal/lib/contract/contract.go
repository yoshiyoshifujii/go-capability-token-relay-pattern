package contract

type (
	Validator interface {
		Validate() error
	}
)

func AssertValidatable(v Validator) {
	if v == nil {
		panic("Validator is nil")
	}
	if err := v.Validate(); err != nil {
		panic(err)
	}
}
