package service

import (
	"context"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/repository"
	"mysql/request"
	"mysql/utils"
)

type IncomeService interface {
	GetIncomeCategory(ctx context.Context) ([]model.IncomeCategory, error)
	AddIncome(ctx context.Context, input request.IncomeRequest) error
}

type incomeservice struct {
	repo repository.IncomeRepository
}

func NewIncomeService(repo repository.IncomeRepository) IncomeService {
	return &incomeservice{
		repo: repo,
	}
}

func (s *incomeservice) GetIncomeCategory(ctx context.Context) ([]model.IncomeCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	data, err := s.repo.GetIncomeCategory(ctx)
	if err != nil {
		return nil, apperror.Internal("failed to fetch provinces", err)
	}
	return data, nil
}

func (s *incomeservice) AddIncome(ctx context.Context, input request.IncomeRequest) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	if err := s.repo.AddIncome(ctx, input); err != nil {
		return apperror.Internal("failed to create fee", err)
	}
	return nil
}
