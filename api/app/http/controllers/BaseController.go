// controllers/BaseController.go

package controllers

import "combined-crawler/api/bootstrap"

// BaseController contains the application instance
type BaseController struct {
	*bootstrap.Application
}

// NewBaseController initializes a new BaseController
func NewBaseController() *BaseController {
	return &BaseController{bootstrap.App()}
}
