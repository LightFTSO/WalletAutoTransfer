package utils

import (
	"regexp"
)

func IsValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}

func GetAverageBlockTime(blockTimes []uint64) float64 {
	if len(blockTimes) < 2 {
		return 1
	}
	var total float64 = 0
	for i, time := range blockTimes[1:] {
		diff := time - blockTimes[i]
		total += float64(diff)
	}
	//fmt.Println(blockTimes)
	avg := total / float64(len(blockTimes)-1)
	return avg
}

func AppendToFIFOSlice(slice []uint64, newItem uint64, capacity int) []uint64 {
	if len(slice) == capacity {
		slice = slice[1:]
	}
	slice = append(slice, newItem)

	return slice
}
