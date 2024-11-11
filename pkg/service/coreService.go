package service

type coreService struct{}

func NewServiceData(s int64) *coreService {
	return &coreService{}
}
