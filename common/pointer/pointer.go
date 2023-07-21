package pointer

func Pointer[Type any](value Type) *Type {
	return &value
}
