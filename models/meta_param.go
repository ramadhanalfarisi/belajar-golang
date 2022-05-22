package models

import (
	"fmt"
	"net/http"
	"strconv"
)

type MetaParam struct {
	Limit  int64
	Page   int64
	Offset int64
}

func GetMetaParam(r *http.Request) (MetaParam, error) {
	var int_limit int64
	var int_page int64
	var int_offset int64

	page, ok := r.URL.Query()["page"]
	limit, ok := r.URL.Query()["limit"]

	climit, err1 := strconv.Atoi(limit[0])
	cpage, err2 := strconv.Atoi(page[0])

	if err1 != nil {
		return MetaParam{}, fmt.Errorf("Limit have to numeric")
	}

	if err2 != nil {
		return MetaParam{}, fmt.Errorf("Page have to numeric")
	}

	if !ok || len(limit[0]) < 1 {
		int_limit = 10
	} else {
		int_limit = int64(climit)
	}

	if !ok || len(page[0]) < 1 {
		int_page = 1
	} else {
		if int64(climit) < 1 {
			return MetaParam{}, fmt.Errorf("Page have to more than 0")
		} else {
			int_page = int64(cpage)
		}
	}

	if int_page == 1 {
		int_offset = 0
	} else {
		int_offset = (int64(int_page) - 1) * int64(int_limit)
	}

	var meta_param MetaParam
	meta_param.Page = int_page
	meta_param.Limit = int_limit
	meta_param.Offset = int_offset
	return meta_param, nil
}
