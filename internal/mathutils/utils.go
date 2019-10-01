package mathutils

const maxInt32 = int32(^uint32(0) >> 1)
const minInt32 = -maxInt32 - 1

// Min computes minimum value from arbitrary number of arguments provided
func Min(values ...int32) int32 {
	minValue := maxInt32

	for _, val := range values {
		if val < minValue {
			minValue = val
		}
	}

	return minValue
}

// Max computes maximum value from arbitrary number of arguments provided
func Max(values ...int32) int32 {
	maxValue := minInt32

	for _, val := range values {
		if val > maxValue {
			maxValue = val
		}
	}

	return maxValue
}
