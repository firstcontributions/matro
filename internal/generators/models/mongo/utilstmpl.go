package mongo

const utilsTmpl = `
package utils

const (
	defaultLimit int64 = 10
	defaultOrder int   = -1
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

func GetSortOrder(order int) bson.D {
	order = order * defaultOrder
	return bson.D{
		{"time_created", order},
		{"_id", order},
	}
}
`
