package helper

import (
	"mysql/model"

	"gorm.io/gorm"
)

func ManageClassFilter(query *gorm.DB, db *gorm.DB, user model.User) *gorm.DB {
	if user.Role.Level >= 7 {
		return query
	}

	var classIDs []int64

	if err := db.
		Table("user_class uc").
		Joins("INNER JOIN class c ON c.id = uc.class_id").
		Where("uc.user_id = ? AND uc.is_active = 1 AND c.is_active = 1", user.ID).
		Pluck("uc.class_id", &classIDs).Error; err != nil {
		return query.Where("1 = 0")
	}

	if len(classIDs) == 0 {
		return query.Where("1 = 0")
	}

	return query.Where("c.id IN ?", classIDs)
}
