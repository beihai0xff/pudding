package entity

// PageQuery is a helper struct to page query
type PageQuery struct {
	// Offset is the offset of the query
	// Offset must be greater than or equal to 0
	// Limit is the limit of the query
	// Limit must be greater than 0
	Offset, Limit int
}
