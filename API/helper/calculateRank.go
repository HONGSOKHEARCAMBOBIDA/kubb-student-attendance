package helper

import (
	"mysql/response"
	"sort"
)

func CalculateRank(data []response.ScoreResponse) {
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Total > data[j].Total
	})

	rank := 0
	lastTotal := -1.0

	for i := range data {
		if data[i].Total != lastTotal {
			rank = i + 1
			lastTotal = data[i].Total
		}

		data[i].Rank = rank
	}
}
