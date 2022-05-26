package model

type Paged struct {
	PageInterface
	TotalRows   int `json:"totalRows"`
	CurrentPage int `json:"currentPage"`
	CurrentSize int `json:"currentSize"`
}

func (p *Paged) GetPageData() *Paged {
	return p
}

type PageInterface interface {
	GetPageData() *Paged
}
