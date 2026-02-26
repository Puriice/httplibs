package iterable

func Map[E, R any](slice []E, mapFn func(E) R) []R {
	result := make([]R, len(slice), cap(slice))

	for i, item := range slice {
		result[i] = mapFn(item)
	}

	return result
}
