package utils

// Min - min(int, int)
func Min(args ...int) int {
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
