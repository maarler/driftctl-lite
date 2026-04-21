// Package paginate provides cursor-free, offset-based pagination for slices
// of drift.Result values.
//
// Usage:
//
//	page, err := paginate.Apply(results, paginate.Options{Page: 2, PageSize: 10})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("page %d/%d\n", page.Page, page.TotalPages)
package paginate
