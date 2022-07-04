package mongo

const utilsTmpl = `
package utils

const (
	defaultLimit int64 = 10
	defaultOrder int   = 1
)

func GetLimitAndSortOrderAndCursor(first, last *int64, after, before *string) (int64, int, *string) {
	if first != nil || after != nil {
		if first == nil {
			return defaultLimit, 1, after
		}
		return *first, 1, after
	}
	if last == nil {
		return defaultLimit, -1, before
	}
	return *last, -1, before
}

func GetSortOrder(sortBy, requestedSortOrder  *string, order int) bson.D {
	order = order * defaultOrder * getSortOrderFromString(requestedSortOrder)
	
	field := "time_created"
	if sortBy != nil {
		field = *sortBy
	}
	return bson.D{
		{field, order},
		{"_id", order},
	}
}

func getSortOrderFromString(order *string) int {
	if order != nil && *order == "desc" {
		return -1
	}
	return 1
}

func ReverseList[T interface{}](list []*T) []*T {
	ln := len(list)
	for i := 0; i< ln/2; i++ {
		list[i], list[ln-1-i] =  list[ln-1-i], list[i]
	}
	return list
} 
`
