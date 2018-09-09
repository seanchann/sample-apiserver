/*

Copyright 2018 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package pages

import (
	"fmt"

	"github.com/golang/glog"
)

//BuildPageLink build a string 'Link' that like as :
// "Link: /api/v1beta1/namespace/default/users?pageSelector=page=1,perPage=1; rel=prev,
// "/api/v1beta1/users?pageSelector=page=3,perPage=1; rel= next,"
// "/api/v1beta1/users?pageSelector=page=5,perPage=1; rel=last"
func BuildDefPageLink(pager Selector, baseLink string) (string, error) {

	glog.V(5).Infof("Got base link %v\r\n", baseLink)
	var link string

	if pager == nil {
		return "", nil
	}

	if !pager.Empty() {
		var prevPageLink, nextPageLink, lastPageLink string

		has, page, perPage := pager.PreviousPage()
		if has {
			prevPageLink = fmt.Sprintf("%v?pageSelector=page=%v,perPage=%v; rel=prev", baseLink, page, perPage)
		}

		has, page, perPage = pager.NextPage()
		if has {
			nextPageLink = fmt.Sprintf("%v?pageSelector=page=%v,perPage=%v; rel=next", baseLink, page, perPage)
		}

		has, page, perPage = pager.LastPage()
		if has {
			lastPageLink = fmt.Sprintf("%v?pageSelector=page=%v,perPage=%v; rel=last", baseLink, page, perPage)
		}

		separator := string("")
		if len(prevPageLink) > 0 {
			link += prevPageLink
			separator = string(",")
		}

		if len(nextPageLink) > 0 {
			link += separator + nextPageLink
			separator = string(",")
		}

		if len(lastPageLink) > 0 {
			link += separator + lastPageLink
		}
	}

	return link, nil
}
