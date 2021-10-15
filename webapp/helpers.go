package webapp

import (
	"net/http"
	"strconv"
)

func orderByFromRequest(r *http.Request, orderMapping map[string]string, defaultOrder string) (string, error) {
	var orderBy *string
	qOrderBy, hasOrderBy := r.URL.Query()["order_by"]
	if hasOrderBy {
		for apiOrder, dbOrder := range orderMapping {
			if qOrderBy[0] == apiOrder {
				orderBy = &dbOrder
				break
			}
		}

		// didn't find the attribute in the mapping
		if orderBy == nil {
			return "", ErrUnknownAttribute
		}
	} else {
		orderBy = &defaultOrder
	}
	return *orderBy, nil
}

func orderAscFromRequest(r *http.Request) (bool, error) {
	orderAsc := true
	qOrderAsc, hasOrderAsc := r.URL.Query()["order_asc"]
	if hasOrderAsc {
		var err error
		orderAsc, err = strconv.ParseBool(qOrderAsc[0])
		if err != nil {
			return false, err
		}
	}
	return orderAsc, nil
}

func paginationFromRequest(r *http.Request, defaultCount int) (int, int, bool, error) {
	// get pages
	qPage, hasPage := r.URL.Query()["page"]
	page := int64(1)
	if hasPage {
		var err error
		page, err = strconv.ParseInt(qPage[0], 10, 64)
		if err != nil {
			return 0, 0, false, err
		}
		if page <= 0 {
			return 0, 0, false, ErrInvalidPage
		}
	}

	// get count
	qCount, hasCount := r.URL.Query()["count"]
	count := int64(defaultCount)
	if hasCount {
		var err error
		count, err = strconv.ParseInt(qCount[0], 10, 64)
		if err != nil {
			return 0, 0, false, err
		}
		if count <= 0 {
			return 0, 0, false, ErrInvalidCount
		}
	}

	return int(page), int(count), hasPage || hasCount, nil
}
