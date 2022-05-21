package models

import (
	"belajar_golang/middlewares"
	"math"
	"net/http"
)

type Pagination struct {
	Page       int64 `json:"page"`       // current page
	Size       int64 `json:"size"`       // limit
	Offset     int64 `json:"offset"`     // limit
	TotalPages int64 `json:"totalPages"` // total page
	Total      int64 `json:"total"`      // total row
	Visible    int64 `json:"visible"`    // total row in current page
	Last       bool  `json:"last"`       // is last page
	First      bool  `json:"first"`      // is first page
}

func (pagination *Pagination) CreatePagination(r *http.Request) Pagination {
	meta_param := r.Context().Value("metaParam").(middlewares.MetaParam)

	int_limit := meta_param.Limit
	int_page := meta_param.Page
	int_offset := meta_param.Offset

	pagination.Page = int64(int_page)
	pagination.Size = int64(int_limit)
	pagination.Offset = int64(int_offset)
	if pagination.Total <= pagination.Size {
		pagination.Visible = pagination.Total
	}else if pagination.Total > pagination.Size{
		current_total := pagination.Page * pagination.Size
		if pagination.Total > current_total {
			pagination.Visible = pagination.Size
		}else{
			mod_total := pagination.Total % pagination.Size
			pagination.Visible = mod_total
		}
	}
	total_pages := math.Ceil(float64(pagination.Total / pagination.Size))
	pagination.TotalPages = int64(total_pages)
	if pagination.Page == 1 {
		pagination.First = true
		pagination.Last = false
	} else if pagination.Page == pagination.TotalPages {
		pagination.First = false
		pagination.Last = true
	} else {
		pagination.First = false
		pagination.Last = false
	}
	return *pagination
}
