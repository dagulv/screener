package domain

type IsZeroer interface {
	IsZero() bool
}

func HasValid(zs ...IsZeroer) bool {
	for _, z := range zs {
		if !z.IsZero() {
			return true
		}
	}

	return false
}
