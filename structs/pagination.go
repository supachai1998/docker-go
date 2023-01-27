package structs

import (
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PanigationBase struct {
	Value     *string `json:"keyword,omitempty" query:"keyword" swagger:"desc(keyword for search)" `
	DateStart *string `json:"date_start,omitempty" query:"date_start" swagger:"desc(format : dd-mm-yyyy hh:mm:ss)"`
	DateEnd   *string `json:"date_end,omitempty" query:"date_end" swagger:"desc(format : dd-mm-yyyy hh:mm:ss)"`
	Page      int64   `json:"page,omitempty" query:"page"`
	Limit     int64   `json:"limit,omitempty" query:"limit"`
}

type PanigationSwagger struct {
	PanigationBase
	OrderBy   *string `json:"order_by,omitempty" query:"order_by" swagger:"enum(id|date|rand)"`
	OrderType *string `json:"order_type,omitempty" query:"order_type" swagger:"enum(asc|desc)"`
}

/*
Example:

	input :
		api?total=0&
		page=1&
		limit=10&
		order_by=id&
		order_type=asc&
		keyword=keyword&
		date_start=01-01-2021 00:00:00& or date_start=01-01-2021&
		date_end=01-01-2021 00:00:00 or date_end=01-01-2021
	output :
		{
			data : ...
			pagination : {
				total : 0,
				page : 1,
				limit : 10,
				order_by : id,
		}
*/
type Pagination struct {
	PanigationBase
	Total     int64   `json:"total,omitempty" query:"total"`
	OrderBy   *string `json:"order_by,omitempty" query:"order_by" swagger:"enum(id|date)"`
	OrderType *string `json:"order_type,omitempty" query:"order_type" swagger:"enum(asc|desc)"`
}

func (p *Pagination) GetOffset() int64 {
	if p.Limit == 0 {
		p.Limit = 50
	}
	if p.Page == 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.Limit
}

func (p *Pagination) InitPanigation(c echo.Context, queryDB *gorm.DB, columnNames []string) (tx *gorm.DB, err error) {
	if p.Limit == 0 {
		p.Limit = 50
	}
	if p.Page == 0 {
		p.Page = 1
	}
	c.Bind(p)
	tx = queryDB.Limit(int(p.Limit)).Offset(int(p.GetOffset()))

	if p.OrderBy != nil {
		if p.OrderType == nil {
			var asc string = "asc"
			p.OrderType = &asc
		}
		switch *p.OrderBy {
		case "date":
			*p.OrderBy = "created_at"
		}
		if *p.OrderBy == "rand" {
			tx = tx.Order("random()")
		} else {
			tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: *p.OrderBy}, Desc: *p.OrderType == "desc"})
		}

	}
	if p.Value != nil {
		if len(columnNames) > 0 {
			var whereClause strings.Builder
			for i, v := range columnNames {
				if i > 0 {
					whereClause.WriteString(" OR ")
				}
				whereClause.WriteString(v + " LIKE '%" + *p.Value + "%'")
			}
			tx = tx.Where(whereClause.String())
		} else {
			tx = tx.Where("id LIKE ?", fmt.Sprintf("%%%s%%", *p.Value))
		}
	}
	if p.DateStart != nil && p.DateEnd != nil {
		var dateStart, dateEnd string
		// if dateStart or dateEnd not time hh:mm:ss in query set default time
		// layoutDDMMYYYY := "02-01-2006"
		layoutDDMMYYYYHHMMSS := "02-01-2006 15:04:05"

		if len(*p.DateStart) <= 10 {
			dateStart = fmt.Sprintf("%s 00:00:00", *p.DateStart)
		} else {
			dateStart = *p.DateStart
		}
		if len(*p.DateEnd) <= 10 {
			dateEnd = fmt.Sprintf("%s 23:59:59", *p.DateEnd)
		} else {
			dateEnd = *p.DateEnd
		}
		// ddmmyyyyhhmmss to yyyy-mm-dd hh:mm:ss
		dateStartSQL, err := time.Parse(layoutDDMMYYYYHHMMSS, dateStart)
		if err != nil {
			return nil, err
		}
		dateEndSQL, err := time.Parse(layoutDDMMYYYYHHMMSS, dateEnd)
		if err != nil {
			return nil, err
		}

		tx = tx.Where("created_at BETWEEN ? AND ?", dateStartSQL, dateEndSQL)
	}
	return tx, nil
}

// GetTotalCount get total count of query
func (p *Pagination) GetTotalCount(queryDB *gorm.DB) (int64, error) {
	var count int64
	err := queryDB.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
