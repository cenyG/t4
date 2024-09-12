package usecase

type DownloadFileUseCase interface {
	Download(filename string) error
}

type downloadFileUseCase struct {
}

func (d downloadFileUseCase) Download(filename string) error {
	//TODO implement me
	panic("implement me")
}

func NewDownloadFileUseCase() DownloadFileUseCase {
	return &downloadFileUseCase{}
}
