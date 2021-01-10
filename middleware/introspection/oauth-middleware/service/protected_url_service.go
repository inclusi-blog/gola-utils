package service

//go:generate mockgen -source=protected_url_service.go -destination=./../mocks/mock_protected_url_service.go -package=mocks

type ProtectedUrlService interface {
	IsProtected(url string) bool
}
