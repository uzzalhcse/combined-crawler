package services

import (
	"combined-crawler/api/app/models"
	"combined-crawler/api/app/repositories"
)

type SiteCollectionService struct {
	Repository *repositories.Repository
}

func NewSiteCollectionService(repo *repositories.Repository) *SiteCollectionService {
	return &SiteCollectionService{Repository: repo}
}
func (s *SiteCollectionService) GetAllSiteCollections() ([]models.SiteCollection, error) {
	return s.Repository.GetAllSiteCollections()
}
func (s *SiteCollectionService) Create(siteCollection *models.SiteCollection) error {
	return s.Repository.CreateSiteCollection(siteCollection)
}

func (s *SiteCollectionService) GetByID(siteID string) (*models.SiteCollection, error) {
	return s.Repository.GetSiteCollectionByID(siteID)
}

func (s *SiteCollectionService) Update(siteID string, update map[string]interface{}) error {
	return s.Repository.UpdateSiteCollection(siteID, update)
}

func (s *SiteCollectionService) Delete(siteID string) error {
	return s.Repository.DeleteSiteCollection(siteID)
}
