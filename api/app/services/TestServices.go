package services

import (
	"combined-crawler/api/app/models"
	"combined-crawler/api/app/repositories"
)

type TestService struct {
	Repository *repositories.Repository
}

func NewTestService(repo *repositories.Repository) *TestService {
	return &TestService{Repository: repo}
}

// GetAll returns all records from the model using the repository
func (s *TestService) GetAllSiteCollections() ([]models.SiteCollection, error) {
	return s.Repository.GetAllSiteCollections()
}
