package entity

type Report struct {
	BestSellingProducts []BestSellingProduct `json:"bestSellingProducts"`
	LowestInventories   []LowInventoryItem   `json:"lowestInventories"`
	OrderSummary        *OrderSummary        `json:"orderSummary"`
}

type BestSellingProduct struct {
	ProductID   uint   `json:"productId"`
	ProductName string `json:"productName"`
	QtySold     uint   `json:"qtySold"`
}

type LowInventoryItem struct {
	ProductID   uint   `json:"productId"`
	ProductName string `json:"productName"`
	Stock       uint   `json:"stock"`
}

type OrderSummary struct {
	TotalOrders    uint `json:"totalOrders"`
	TotalAmount    uint `json:"totalQty"`
	TotalCustomers uint `json:"totalCustomers"`
}
