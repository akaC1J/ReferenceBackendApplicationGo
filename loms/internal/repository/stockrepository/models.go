// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package stockrepository

type Stock struct {
	Sku           int64
	TotalCount    int64
	ReservedCount int64
}