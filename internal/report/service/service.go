package service

type Service struct {
	Report *Report
}

func New(reportRepo ReportRepo) *Service {
	r := NewReport(reportRepo)

	return &Service{
		Report: r,
	}
}
