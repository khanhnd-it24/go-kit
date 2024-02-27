package utility

func CvtToListInterface[T any](p []T) []interface{} {
	out := make([]interface{}, 0)
	for _, t := range p {
		out = append(out, t)
	}
	return out
}
