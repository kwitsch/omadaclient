package model

type Paged struct {
	PageInterface
	TotalRows   uint `json:"totalRows"`
	CurrentPage uint `json:"currentPage"`
	CurrentSize uint `json:"currentSize"`
}

func (p *Paged) GetTotalRows() uint {
	return p.TotalRows
}

func (p *Paged) GetCurrentPage() uint {
	return p.CurrentPage
}

type PageInterface interface {
	GetTotalRows() uint
	GetCurrentPage() uint
}
