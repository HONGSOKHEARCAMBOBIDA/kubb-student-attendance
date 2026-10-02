package helper

import (
	"fmt"
)

type Layout int

const (
	Bachelor Layout = iota
	Master
)

func (l Layout) String() string {
	if l == Master {
		return "Master"
	}
	return "Bachelor"
}

func LayoutForProgramme(programmID int64) (Layout, error) {
	switch programmID {
	case 1:
		return Bachelor, nil
	case 3:
		return Master, nil
	}
	return 0, fmt.Errorf("programme %d has no transcript layout (only Bachelor=1, Master=3)", programmID)
}

const minRowsPerSemester = 5

func ordinal(n int) string {
	switch {
	case n%100 >= 11 && n%100 <= 13:
		return fmt.Sprintf("%dth", n)
	case n%10 == 1:
		return fmt.Sprintf("%dst", n)
	case n%10 == 2:
		return fmt.Sprintf("%dnd", n)
	case n%10 == 3:
		return fmt.Sprintf("%drd", n)
	}
	return fmt.Sprintf("%dth", n)
}
