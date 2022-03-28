package pkg

import "fmt"

type Pager struct {
	currentIndex int
	data         []interface{}
	pageSize     int

	total int

	pageStartIndex, pageEndIndex int
}

func NewPager(data []interface{}) *Pager {
	page := &Pager{
		data:         data,
		pageSize:     10,
		currentIndex: 0,
		total:        len(data),
	}

	if page.pageEndIndex > page.total {
		page.pageEndIndex = page.total
	}

	return page
}

func (p *Pager) SetData(data []interface{}) *Pager {
	p.data = data
	return p
}

func (p Pager) Data() []interface{} {
	return p.data
}

func (p *Pager) PageSize(size int) *Pager {
	p.pageSize = size

	p.pageEndIndex = p.pageStartIndex + p.pageSize
	if p.pageEndIndex > p.total {
		p.pageEndIndex = p.total
	}
	return p
}

func (p *Pager) SelectedIndex() int {
	return p.currentIndex
}

func (p *Pager) SelectedRecord() interface{} {
	return p.data[p.currentIndex]
}

func (p *Pager) StartIndex() int {
	return p.pageStartIndex
}

func (p *Pager) EndIndex() int {
	return p.pageEndIndex
}

func (p *Pager) TotalRecord() int {
	return p.total
}

func (p *Pager) resetEndIndex() {
	p.pageEndIndex = p.pageStartIndex + p.pageSize
	if p.pageEndIndex > p.total {
		p.pageEndIndex = p.total
	}
}

func (p *Pager) resetStartIndex() {
	p.pageStartIndex = p.pageEndIndex - p.pageSize
	if p.pageStartIndex < 0 {
		p.pageStartIndex = 0
	}
}

func (p *Pager) Next() {
	if p.currentIndex >= p.total-1 {
		return
	}

	p.currentIndex++

	// endIndex/endIndex
	if p.currentIndex >= p.pageEndIndex {
		p.pageStartIndex = p.currentIndex

		p.pageEndIndex = p.pageStartIndex + p.pageSize
		if p.pageEndIndex > p.total {
			p.pageEndIndex = p.total

			// 向前倒近推计算
			p.resetStartIndex()
		}
	}
}

func (p *Pager) Prev() {
	if p.currentIndex <= 0 {
		return
	}
	p.currentIndex--

	// startIndex/endIndex
	if p.currentIndex <= p.pageStartIndex {
		p.pageEndIndex = p.currentIndex

		// startIndex
		//p.resetStartIndex()
		p.pageStartIndex = p.pageEndIndex - p.pageSize
		if p.pageStartIndex < 0 {
			p.pageStartIndex = 0
			p.resetEndIndex()
		}
	}
}

func (p *Pager) PgUp() {
	// current is FirstPage
	if p.pageStartIndex == 0 {
		return
	}

	if p.currentIndex-p.pageSize < 0 {
		p.currentIndex = 0
	} else {
		p.currentIndex -= p.pageSize
	}

	// startIndex
	p.pageStartIndex = p.pageStartIndex - p.pageSize
	if p.pageStartIndex < 0 {
		p.pageStartIndex = 0
	}

	// endIndex
	p.pageEndIndex = p.pageStartIndex + p.pageSize
	if p.pageEndIndex > p.total {
		p.pageEndIndex = p.total
	}
}

func (p *Pager) PgDown() {
	// current is lastPage
	if p.pageEndIndex == p.total {
		return
	}

	// 超出当前页的显示范围，自动从数据最后向前计算出 pageSize
	if p.currentIndex+p.pageSize >= p.total {
		startIndex := p.total - p.pageSize
		if startIndex < 0 {
			p.currentIndex = 0
		} else {
			p.currentIndex = startIndex
		}
	} else {
		p.currentIndex += p.pageSize
	}

	// startIndex
	p.pageStartIndex = p.pageStartIndex + p.pageSize
	if p.pageStartIndex > p.total {
		p.pageEndIndex = p.total
		p.resetStartIndex()
	} else {
		p.resetEndIndex()
		p.resetStartIndex()
	}
}

func (p *Pager) Result() []interface{} {
	return p.data[p.pageStartIndex:p.pageEndIndex]
}

func (p *Pager) String() string {
	return fmt.Sprintf("total=%d pagesize=%d startIndex=%d currentIndex=%d [value=%#v] endIndex=%d\n",
		p.total, p.pageSize, p.pageStartIndex, p.currentIndex, p.data[p.currentIndex], p.pageEndIndex)
}
