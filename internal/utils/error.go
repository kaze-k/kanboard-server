package utils

func HandleError[T any](some *T, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	return some, nil
}
