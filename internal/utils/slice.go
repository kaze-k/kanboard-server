package utils

func MergeByKey[T any, K comparable](a, b []T, keyFunc func(T) K, keepFirst bool) []T {
	resultMap := make(map[K]T)

	for _, item := range a {
		k := keyFunc(item)
		resultMap[k] = item
	}

	for _, item := range b {
		k := keyFunc(item)
		if _, exists := resultMap[k]; !exists || !keepFirst {
			resultMap[k] = item
		}
	}

	result := make([]T, 0, len(resultMap))
	for _, v := range resultMap {
		result = append(result, v)
	}

	return result
}

func UniqueUintSlice(input []uint) []uint {
	seen := make(map[uint]struct{})
	result := make([]uint, 0, len(input))

	for _, v := range input {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}
