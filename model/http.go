package model

type Paged struct {
	PageInterface
	TotalRows   int `json:"totalRows"`
	CurrentPage int `json:"currentPage"`
	CurrentSize int `json:"currentSize"`
}

func (p *Paged) GetTotalRows() int {
	return p.TotalRows
}

func (p *Paged) GetCurrentPage() int {
	return p.CurrentPage
}

type PageInterface interface {
	GetTotalRows() int
	GetCurrentPage() int
}
