package models

import (
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
	var climit, cpage int

	page, ok1 := r.URL.Query()["page"]
	limit, ok2 := r.URL.Query()["limit"]
	if ok1 {
		cpage, _ = strconv.Atoi(page[0])
	}
	if ok2 {
		climit, _ = strconv.Atoi(limit[0])
	}

	if !ok1 || len(page[0]) < 1 {
		int_page = 1
	} else {
		if int64(climit) < 1 {
			int_page = 1
		} else {
			int_page = int64(cpage)
		}
	}

	if !ok2 || len(limit[0]) < 1 {
		int_limit = 10
	} else {
		int_limit = int64(climit)
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
