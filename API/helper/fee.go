package helper

import "mysql/model"

func CalculateDiscountBySchoolarship(amount float64, group *model.Schoolarship) float64 {
	if group == nil || !group.Active {
		return 0
	}
	switch group.DiscountType {
	case model.DiscountPercentage:
		return amount * (group.DiscountPercentage / 100)
	case model.DiscountAmount:
		return group.DiscountAmount
	default:
		return 0
	}
}

func GetFeeAmountPerYear(
	majorPrice model.MajorPrice,
	feeSchedule model.FeeSchedule,
) float64 {

	switch feeSchedule.FeeInterval {

	case model.FeeScheduleIntervalMonthlyFee:
		return majorPrice.MonthlyFee * 12

	case model.FeeScheduleIntervalQuarterFee:
		return majorPrice.QuarterFee * 4

	case model.FeeScheduleIntervalSemesterFee:
		return majorPrice.SemesterFee * 2

	case model.FeeScheduleIntervalYearlyFee:
		return majorPrice.YearFee

	default:
		return 0
	}
}

func GetNextDueDate(feeschedule model.FeeSchedule) (count int, months int) {
	switch feeschedule.FeeInterval {
	case model.FeeScheduleIntervalMonthlyFee:
		return 12, 1

	case model.FeeScheduleIntervalQuarterFee:
		return 4, 3

	case model.FeeScheduleIntervalSemesterFee:
		return 2, 6

	case model.FeeScheduleIntervalYearlyFee:
		return 1, 12

	default:
		return 0, 0
	}
}
