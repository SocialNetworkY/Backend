package service

type Service struct {
	Tag *TagService
}

func New(tagStorage TagStorage) *Service {
	return &Service{
		Tag: NewTagService(tagStorage),
	}
}
