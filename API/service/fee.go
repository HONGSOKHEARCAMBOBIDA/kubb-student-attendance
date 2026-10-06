package service

import (
	"context"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/repository"
	"mysql/request"
	"mysql/utils"
)

type FeeService interface {
	GetFeeSchedule(ctx context.Context) ([]model.FeeSchedule, error)
	AddFee(ctx context.Context, input request.FeeRequestCreate) error
}

type feeservice struct {
	repo repository.FeeRepository
}

func NewFeeService(repo repository.FeeRepository) FeeService {
	return &feeservice{
		repo: repo,
	}
}

func (s *feeservice) GetFeeSchedule(ctx context.Context) ([]model.FeeSchedule, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	data, err := s.repo.GetFeeSchedule(ctx)
	if err != nil {
		return nil, apperror.Internal("failed to fetch provinces", err)
	}

	return data, nil
}

func (s *feeservice) AddFee(ctx context.Context, input request.FeeRequestCreate) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	if err := s.repo.AddFee(ctx, input); err != nil {
		return apperror.Internal("failed to create fee", err)
	}
	return nil
}
