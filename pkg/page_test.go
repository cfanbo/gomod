package pkg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func Test_Page(t *testing.T) {
	debug := func(data []interface{}) {
		for k, v := range data {
			fmt.Printf("%d\t %#v\n", k, v)
		}
	}

	// {0,1,2} {3,4,5} {5,6,7}
	testData := []int{0, 1, 2, 3, 4, 5, 6, 7}

	var data = make([]interface{}, len(testData))
	for k, v := range testData {
		data[k] = v
	}

	page := NewPager(data)
	page.PageSize(3)

	got := func(data []interface{}) string {
		s := make([]string, 0, len(data))
		for _, v := range data {
			s = append(s, strconv.Itoa(v.(int)))
		}
		return strings.Join(s, ",")
	}

	// page 1
	// init
	assert.Equal(t, "0,1,2", got(page.Result()), "they should be equal[init]")

	//Next 1
	page.Next()
	assert.Equal(t, "0,1,2", got(page.Result()), "they should be equal[Next page1]")

	page.Next()
	assert.Equal(t, "0,1,2", got(page.Result()), "they should be equal[Next page1]")

	// page 2
	page.Next()
	assert.Equal(t, "3,4,5", got(page.Result()), "they should be equal[Next page2-1]")
	page.Next()
	page.Next()
	assert.Equal(t, "3,4,5", got(page.Result()), "they should be equal[Next page2-2]")

	// page 3
	page.Next()
	assert.Equal(t, "5,6,7", got(page.Result()), "they should be equal[Next page3-1]")

	page.Next()
	page.Next()
	t.Log(page)
	debug(page.Result())

	page.Next()
	t.Log(page)
	debug(page.Result())

	page.Next()
	assert.Equal(t, "5,6,7", got(page.Result()), "they should be equal[Next page3-4]")
	t.Log(page)
	debug(page.Result())

	page.Prev()
	t.Log(page)
	debug(page.Result())
	assert.Equal(t, "5,6,7", got(page.Result()), "they should be equal[Prev page3-1]")
	//page.Prev()
	page.Prev()
	t.Log(page)
	debug(page.Result())
	assert.Equal(t, "2,3,4", got(page.Result()), "they should be equal[Prev page3-2]")

}
