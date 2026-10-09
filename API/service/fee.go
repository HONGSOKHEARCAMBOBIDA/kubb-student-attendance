package service

import (
	"context"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/repository"
	"mysql/request"
	"mysql/response"
	"mysql/utils"
)

type FeeService interface {
	GetFeeSchedule(ctx context.Context) ([]model.FeeSchedule, error)
	AddFee(ctx context.Context, input request.FeeRequestCreate) error
	GetSchoolarship(ctx context.Context) ([]model.Schoolarship, error)
	GetUserClass(ctx context.Context, userID int) ([]response.UserClass, error)
	AddFeeTransaction(ctx context.Context, input request.FeeTransaction) error
	PrintInvoice(ctx context.Context, id int) (response.PrintInvoiceResponse, error)
	DeleteFeeTransaction(ctx context.Context, id int) error
	DeleteFee(ctx context.Context, id int) error
}

type feeservice struct {
	repo repository.FeeRepository
}

func NewFeeService(repo repository.FeeRepository) FeeService {
	return &feeservice{
		repo: repo,
	}
}

func (s *feeservice) GetUserClass(ctx context.Context, userID int) ([]response.UserClass, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	data, err := s.repo.GetUserClass(ctx, userID)
	if err != nil {
		return nil, apperror.Internal("failed to fetch provinces", err)
	}

	return data, nil
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

func (s *feeservice) GetSchoolarship(ctx context.Context) ([]model.Schoolarship, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	data, err := s.repo.GetSchoolarship(ctx)
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

func (s *feeservice) AddFeeTransaction(ctx context.Context, input request.FeeTransaction) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	if err := s.repo.AddFeeTransaction(ctx, input); err != nil {
		return apperror.Internal("failed to create fee", err)
	}
	return nil
}

func (s *feeservice) PrintInvoice(ctx context.Context, id int) (response.PrintInvoiceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()

	data, err := s.repo.PrintInvoice(ctx, id)
	if err != nil {
		return response.PrintInvoiceResponse{}, err
	}

	return data, nil
}

func (s *feeservice) DeleteFeeTransaction(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	if err := s.repo.DeleteFeeTransaction(ctx, id); err != nil {
		return apperror.Internal("failed to delete fee", err)
	}
	return nil
}

func (s *feeservice) DeleteFee(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	if err := s.repo.DeleteFee(ctx, id); err != nil {
		return apperror.Internal("failed to delete fee", err)
	}
	return nil
}
