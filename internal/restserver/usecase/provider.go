package usecase

import (
	"T4_test_case/internal/restserver/repo"
	"T4_test_case/internal/restserver/services"
)

type Provider interface {
	GetDownloadFileUseCase() DownloadFileUseCase
	GetUploadFileUseCase() UploadFileUseCase
}

func NewUseCaseProvider(storageServersProvider services.StorageServersProvider, fileRepo repo.FileRepository) Provider {
	return &useCaseProvider{
		downloadUseCase:   NewDownloadFileUseCase(storageServersProvider, fileRepo),
		uploadFileUseCase: NewUploadFileUseCase(storageServersProvider, fileRepo),
	}
}

type useCaseProvider struct {
	downloadUseCase   DownloadFileUseCase
	uploadFileUseCase UploadFileUseCase
}

func (p *useCaseProvider) GetDownloadFileUseCase() DownloadFileUseCase {
	return p.downloadUseCase
}

func (p *useCaseProvider) GetUploadFileUseCase() UploadFileUseCase {
	return p.uploadFileUseCase
}
