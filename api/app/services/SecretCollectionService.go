package services

import (
	"combined-crawler/api/app/models"
	"combined-crawler/api/app/repositories"
)

type SecretCollectionService struct {
	Repository *repositories.Repository
}

func NewSecretCollectionService(repo *repositories.Repository) *SecretCollectionService {
	return &SecretCollectionService{Repository: repo}
}
func (s *SecretCollectionService) GetAllSiteSecret() ([]models.SiteSecret, error) {
	return s.Repository.GetAllSiteSecretCollections()
}
func (s *SecretCollectionService) GetAllGlobalSecret() ([]models.GlobalSecret, error) {
	return s.Repository.GetAllGlobalSecretCollections()
}
func (s *SecretCollectionService) Create(siteSecret *models.SiteSecret) error {
	return s.Repository.CreateSecretCollection(siteSecret)
}

func (s *SecretCollectionService) GetByID(siteID string) (*models.SiteSecret, error) {
	return s.Repository.GetSiteSecretCollectionByID(siteID)
}

func (s *SecretCollectionService) Update(siteID string, update map[string]interface{}) error {
	return s.Repository.UpdateSiteCollection(siteID, update)
}

func (s *SecretCollectionService) Delete(siteID string) error {
	return s.Repository.DeleteSiteCollection(siteID)
}
