package services

import (
	"combined-crawler/api/app/models"
	"combined-crawler/api/app/repositories"
)

type UrlCollectionService struct {
	Repository *repositories.Repository
}

func NewUrlCollectionService(repo *repositories.Repository) *UrlCollectionService {
	return &UrlCollectionService{Repository: repo}
}

func (s *UrlCollectionService) Create(urlCollection *models.UrlCollection) error {
	return s.Repository.CreateUrlCollection(urlCollection)
}

func (s *UrlCollectionService) GetByID(collectionID string) (*models.UrlCollection, error) {
	return s.Repository.GetUrlCollectionByID(collectionID)
}

func (s *UrlCollectionService) Update(collectionID string, update map[string]interface{}) error {
	return s.Repository.UpdateUrlCollection(collectionID, update)
}

func (s *UrlCollectionService) Delete(collectionID string) error {
	return s.Repository.DeleteUrlCollection(collectionID)
}
