package service

import (
	"context"
	"mysql/constant/apperror"
	"mysql/model"
	"mysql/repository"
	"mysql/utils"
)

type LocationService interface {
	GetProvince(ctx context.Context) ([]model.Province, error)
	GetDistrict(ctx context.Context, id int) ([]model.District, error)
	GetCommune(ctx context.Context, id int) ([]model.Communce, error)
	GetVillage(ctx context.Context, id int) ([]model.Village, error)
}

type locationService struct {
	repo repository.LocationRepository
}

func NewLocationService(repo repository.LocationRepository) LocationService {
	return &locationService{
		repo: repo,
	}
}

func (s *locationService) GetProvince(ctx context.Context) ([]model.Province, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	data, err := s.repo.GetProvince(ctx)
	if err != nil {
		return nil, apperror.Internal("failed to fetch provinces", err)
	}

	return data, nil
}

func (s *locationService) GetDistrict(ctx context.Context, id int) ([]model.District, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	data, err := s.repo.GetDistrict(ctx, id)
	if err != nil {
		return nil, apperror.Internal("failed to fetch districts", err)
	}

	return data, nil
}

func (s *locationService) GetCommune(ctx context.Context, id int) ([]model.Communce, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	data, err := s.repo.GetCommune(ctx, id)
	if err != nil {
		return nil, apperror.Internal("failed to fetch communes", err)
	}

	return data, nil
}

func (s *locationService) GetVillage(ctx context.Context, id int) ([]model.Village, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultQueryTimeout)
	defer cancel()
	data, err := s.repo.GetVillage(ctx, id)
	if err != nil {
		return nil, apperror.Internal("failed to fetch villages", err)
	}

	return data, nil
}
