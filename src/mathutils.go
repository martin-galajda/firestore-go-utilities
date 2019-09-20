package main

const maxInt32 = int32(^uint32(0) >> 1)
const minInt32 = -maxInt32 - 1

func min(values ...int32) int32 {
	minValue := maxInt32

	for _, val := range values {
		if val < minValue {
			minValue = val
		}
	}

	return minValue
}

func max(values ...int32) int32 {
	maxValue := minInt32

	for _, val := range values {
		if val > maxValue {
			maxValue = val
		}
	}

	return maxValue
}
