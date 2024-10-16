package service

type Service struct {
	Report *Report
}

func New(reportRepo ReportRepo, pg PostGateway) *Service {
	r := NewReport(reportRepo, pg)

	return &Service{
		Report: r,
	}
}
