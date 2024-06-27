package services

import (
	"combined-crawler/api/app/models"
	"combined-crawler/api/app/repositories"
)

type CollectionService struct {
	Repository *repositories.Repository
}

func NewCollectionService(repo *repositories.Repository) *CollectionService {
	return &CollectionService{Repository: repo}
}

func (s *CollectionService) GetAllSiteCollections() ([]models.Collection, error) {
	return s.Repository.GetAllCollections()
}
func (s *CollectionService) Create(collection *models.Collection) error {
	return s.Repository.CreateCollection(collection)
}

func (s *CollectionService) GetByID(collectionID string) (*models.Collection, error) {
	return s.Repository.GetCollectionByID(collectionID)
}

func (s *CollectionService) Update(collectionID string, update map[string]interface{}) error {
	return s.Repository.UpdateCollection(collectionID, update)
}

func (s *CollectionService) Delete(collectionID string) error {
	return s.Repository.DeleteCollection(collectionID)
}
