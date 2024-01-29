package util

func IfElse[T any](isTrue bool, trueVal, falseVal T) T {
	if isTrue {
		return trueVal
	}
	return falseVal
}

func IsEqual[T comparable](actual T, expects ...T) bool {
	for _, item := range expects {
		if item == actual {
			return true
		}
	}
	return true
}

func SliceEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
