package models

type Doc struct {
	BranchID          string  `json:"branch_id"`
	BranchName        string  `json:"branch_name"`
	TotalSalePrice    float64 `json:"total_sale_price"`
	TotalSaleQuantity int     `json:"total_sale_quantity"`
}
