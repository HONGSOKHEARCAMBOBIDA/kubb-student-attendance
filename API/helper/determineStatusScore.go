package helper

import "time"

func DetermineStatusScore(checkTime, scheduledTime string, isCheckIn bool) (string, float64) {
	ct, err1 := time.Parse("15:04:05", checkTime)
	st, err2 := time.Parse("15:04:05", scheduledTime)
	if err1 != nil || err2 != nil {
		return "A", 0
	}

	diff := ct.Sub(st).Minutes()
	if !isCheckIn {
		diff = -diff // for check-out, being early is "late" in reverse
	}

	switch {
	case diff <= 0:
		return "PR", 1.0 // on time or early
	case diff <= 15:
		return "P", 0.5 // late but within grace period
	default:
		return "A", 0.0 // too late / absent
	}
}
