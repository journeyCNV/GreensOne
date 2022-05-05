package test

import (
	"fmt"
	"github.com/journeycnv/greensone/gsweb/container"
)

/**
Demo
*/

const key = "greens:hhh"

type DService interface {
	MustSmile() Cheer
}

type RDService struct {
	DService
	c container.GContainer
}

func (s *RDService) MustSmile() Cheer {
	return Cheer{
		happy: "smile",
	}
}

type DServiceProvider struct {
	container.ServiceProvider
}

func (d *DServiceProvider) Name() string {
	return key
}

func (d *DServiceProvider) Register(c container.GContainer) container.NewInstance {
	return NewService
}

func (d *DServiceProvider) IsDefer() bool {
	return true
}

func (d *DServiceProvider) Params(c container.GContainer) []interface{} {
	return []interface{}{c}
}

func (d *DServiceProvider) Boot(c container.GContainer) error {
	fmt.Println("DserviceProvider call Boot() : service boot")
	return nil
}

func NewService(params ...interface{}) (interface{}, error) {
	c := params[0].(container.GContainer)
	return &RDService{c: c}, nil
}
