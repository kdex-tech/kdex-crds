package predicate

func IfE[T any](predicate bool, trueVal T, elseVal T) T {
	if predicate {
		return trueVal
	}
	return elseVal
}

func IfEF[T any](predicate bool, trueFunc func() T, elseFunc func() T) T {
	if predicate {
		return trueFunc()
	}
	return elseFunc()
}
