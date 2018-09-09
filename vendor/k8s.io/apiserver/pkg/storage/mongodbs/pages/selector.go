/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package pages

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

type Selector interface {
	//present page.default first page must be 1
	PresentPage() (has bool, page, perPage uint64)
	//previous page
	PreviousPage() (has bool, page, perPage uint64)
	//previous page
	NextPage() (has bool, page, perPage uint64)
	//previous page
	LastPage() (has bool, page, perPage uint64)
	//set count with all of itmes
	SetItemTotal(itemsSum uint64)
	// Empty returns true if this pagination does not restrict the pagination space.
	Empty() bool
	// String returns a human readable string that represents this pager.
	String() string
	// RequirePage returns RequirePage, RequirePerPage
	RequirePage() (uint64, uint64)
	//determine condition for this request,return:hasPage, perpageCount, skipcount
	Condition(itemsSum uint64) (bool, uint64, uint64)
}

const (
	pageKeyName    string = "page"
	perPageKeyName string = "perPage"
	perPageDefault uint64 = 30
)

// Everything will be return
func Everything() Selector {
	return &hasPage{}
}

type hasPage struct {
	itemTotal                                      uint64
	requirePagination, presentPagination           [2]uint64
	prevPagination, nextPagination, lastPagination [2]uint64
}

func (h *hasPage) PresentPage() (has bool, page, perPage uint64) {
	if h.presentPagination[0] == 0 {
		return false, 0, 0
	}

	return true, h.presentPagination[0], h.presentPagination[1]
}

func (h *hasPage) PreviousPage() (has bool, page, perPage uint64) {
	if h.prevPagination[0] == 0 {
		return false, 0, 0
	}

	return true, h.prevPagination[0], h.prevPagination[1]
}

func (h *hasPage) NextPage() (has bool, page, perPage uint64) {
	if h.nextPagination[0] == 0 {
		return false, 0, 0
	}

	return true, h.nextPagination[0], h.nextPagination[1]
}

func (h *hasPage) LastPage() (has bool, page, perPage uint64) {
	if h.lastPagination[0] == 0 {
		return false, 0, 0
	}
	return true, h.lastPagination[0], h.lastPagination[1]
}

func (h *hasPage) SetItemTotal(itemsSum uint64) {
	h.itemTotal = itemsSum
	remainder := h.itemTotal % h.requirePagination[1]
	page := h.itemTotal / h.requirePagination[1]

	glog.Infof("pagination:remainder:%v,page:%v", remainder, page)

	if page == 0 {
		if page+remainder < h.requirePagination[0] {
			h.presentPagination = [...]uint64{0, 0}
			h.prevPagination = [...]uint64{0, 0}
			h.nextPagination = [...]uint64{0, 0}
			h.lastPagination = [...]uint64{0, 0}
		} else if remainder > 0 {
			h.presentPagination[0] = page + 1
			h.presentPagination[1] = remainder

			h.lastPagination[0] = page + 1
			h.lastPagination[1] = remainder
		} else {
			h.presentPagination = [...]uint64{0, 0}
			h.prevPagination = [...]uint64{0, 0}
			h.nextPagination = [...]uint64{0, 0}
			h.lastPagination = [...]uint64{0, 0}
		}
	} else {

		emptyPagination := [...]uint64{0, 0}
		requirePage := h.requirePagination[0]
		requirePerPage := h.requirePagination[1]
		nextPage := requirePerPage + 1
		prevPage := requirePage - 1

		if remainder == 0 {
			h.lastPagination[0] = page
			h.lastPagination[1] = requirePerPage
		} else {
			h.lastPagination[0] = page + 1
			h.lastPagination[1] = remainder
		}

		//is last page
		if requirePage == h.lastPagination[0] {
			h.presentPagination = h.lastPagination
			h.nextPagination = emptyPagination
		} else {
			h.presentPagination = [...]uint64{requirePage, requirePerPage}
			nextPage = h.presentPagination[0] + 1
			if nextPage == h.lastPagination[0] {
				h.nextPagination = h.lastPagination
			} else {
				h.nextPagination = [...]uint64{nextPage, requirePerPage}
			}
		}

		if h.presentPagination[0] == 1 {
			h.prevPagination = [...]uint64{0, 0}
		} else {
			prevPage = h.presentPagination[0] - 1
			h.prevPagination = [...]uint64{prevPage, requirePerPage}
		}
	}

	glog.Infof("pagination:%+v", *h)
}

func (h *hasPage) Empty() bool {
	if h.requirePagination[0] == 0 {
		return true
	}
	return false
}

func (h *hasPage) String() string {
	return fmt.Sprintf("page=%v,perPage=%v", h.requirePagination[0], h.requirePagination[1])
}

func (h *hasPage) RequirePage() (uint64, uint64) {
	return h.requirePagination[0], h.requirePagination[1]
}

//PagerToCondition build pager condition by total count
//return value:have pager, perPagecount,skipitem
func (h *hasPage) Condition(total uint64) (bool, uint64, uint64) {

	//update current item sum
	h.SetItemTotal(total)

	//if there have not present page do nothing
	has, _, perPage := h.PresentPage()
	if !has {
		return false, 0, 0
	}

	var skip uint64
	hasPrev, prevPage, prevPerPage := h.PreviousPage()
	if hasPrev {
		skip = prevPage * prevPerPage
	} else {
		skip = 0
	}

	return true, perPage, skip
}

func try(pagination, op string) (lhs, rhs string, ok bool) {
	pieces := strings.Split(pagination, op)
	if len(pieces) == 2 {
		return pieces[0], pieces[1], true
	}
	return "", "", false
}

func SelectorFromSet(page, perPage uint64) Selector {
	return &hasPage{
		itemTotal:         0,
		requirePagination: [...]uint64{page, perPage},
	}
}

//parsePagination accept format like : "page=1,perPage=10"
func parsePagination(pagination string) (Selector, error) {
	parts := strings.Split(pagination, ",")
	sort.StringSlice(parts).Sort()
	var page, perPage uint64
	page = 0
	perPage = 0
	var err error

	glog.Infof("parse pagination:%s", pagination)

	if len(pagination) == 0 {
		return nil, nil
	}

	for _, part := range parts {
		if part == "" {
			continue
		}
		if lhs, rhs, ok := try(part, "=="); ok {
			if lhs == pageKeyName {
				page, err = strconv.ParseUint(rhs, 10, 64)
			} else if lhs == perPageKeyName {
				perPage, err = strconv.ParseUint(rhs, 10, 64)
			} else {
				return nil, fmt.Errorf("invalid pagination: '%s'; can't understand '%s'", pagination, part)
			}
		} else if lhs, rhs, ok := try(part, "="); ok {
			if lhs == pageKeyName {
				page, err = strconv.ParseUint(rhs, 10, 64)
			} else if lhs == perPageKeyName {
				perPage, err = strconv.ParseUint(rhs, 10, 64)
			} else {
				return nil, fmt.Errorf("invalid pagination: '%s'; can't understand '%s'", pagination, part)
			}
		} else {
			return nil, fmt.Errorf("invalid pagination: '%s'; can't understand '%s'", pagination, part)
		}

		if err != nil {
			return nil, fmt.Errorf("invalid pagination: '%s'; can't understand '%s'", pagination, part)
		}
	}

	if page <= 0 || perPage <= 0 {
		return nil, fmt.Errorf("invalid pagination: '%s'; can't understand the value of page or perPage", pagination)
	}

	// if perPage <= 0 {
	// 	perPage = perPageDefault
	// }

	return &hasPage{
		itemTotal:         0,
		requirePagination: [...]uint64{page, perPage},
	}, nil
}

func ParsePaginaton(pagination string) (Selector, error) {
	return parsePagination(pagination)
}
