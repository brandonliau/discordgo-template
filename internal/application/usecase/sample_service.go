package usecase

import (
	"fmt"
)

type SampleResult struct {
	Count   int
	Title   string
	Message string
}

type SampleService struct{}

func NewSampleService() *SampleService {
	return &SampleService{}
}

func (s *SampleService) Get(count int) (*SampleResult, error) {
	if count < 0 {
		return nil, fmt.Errorf("count cannot be negative")
	}
	return &SampleResult{
		Count:   count,
		Title:   "Sample interaction",
		Message: fmt.Sprintf("The button has been clicked %d times.", count),
	}, nil
}
