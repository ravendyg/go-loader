package utils

// Min - min(int64...)
func Min(args ...int64) int64 {
	if len(args) == 0 {
		return 0
	}
	min := args[0]
	for _, v := range args {
		if v < min {
			min = v
		}
	}

	return min
}
