package service

import (
	"context"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/utils"

	"gorm.io/gorm"
)

type SchoolarshipService interface {
	CreateSchoolarship(ctx context.Context, input request.SchoolarshipRequestCreate) error
	UpdateSchoolarship(ctx context.Context, id int, input request.SchoolarshipRequestUpdate) error
}

type schoolarshipservice struct {
	db *gorm.DB
}

func NewSchoolarshipService() SchoolarshipService {
	return &schoolarshipservice{
		db: config.DB,
	}
}

func (s *schoolarshipservice) CreateSchoolarship(ctx context.Context, input request.SchoolarshipRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newdata := model.Schoolarship{
			Code:               input.Code,
			Name:               input.Name,
			DiscountType:       input.DiscountType,
			DiscountPercentage: input.DiscountPercentage,
			DiscountAmount:     input.DiscountAmount,
			Description:        input.Description,
			Active:             true,
		}
		if err := tx.Create(&newdata).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *schoolarshipservice) UpdateSchoolarship(ctx context.Context, id int, input request.SchoolarshipRequestUpdate) error {
	updates := map[string]interface{}{}

	if input.Code != nil {
		updates["code"] = input.Code
	}
	if input.Name != nil {
		updates["name"] = input.Name
	}
	if input.DiscountType != nil {
		updates["discount_type"] = input.DiscountType
	}
	if input.DiscountPercentage != nil {
		updates["discount_percentage"] = input.DiscountPercentage
	}
	if input.DiscountAmount != nil {
		updates["discount_amount"] = input.DiscountAmount
	}
	if input.Description != nil {
		updates["description"] = input.Description
	}
	if input.Active != nil {
		updates["active"] = input.Active
	}
	if err := s.db.WithContext(ctx).Model(&model.Schoolarship{}).Where("id =?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}
