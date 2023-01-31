package utils

func Reduce[Tin, Tout any](s []Tin, f func(Tout, Tin) Tout, initValue Tout) Tout {
	acc := initValue
	for _, v := range s {
		acc = f(acc, v)
	}

	return acc
}
