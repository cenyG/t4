package usecase

import "T4_test_case/internal/restserver/services"

type Provider interface {
	GetDownloadFileUseCase() DownloadFileUseCase
	GetUploadFileUseCase() UploadFileUseCase
}

func NewUseCaseProvider(storageServersProvider services.StorageServersProvider) Provider {
	return &useCaseProvider{
		downloadUseCase:   NewDownloadFileUseCase(),
		uploadFileUseCase: NewUploadFileUseCase(storageServersProvider),
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
