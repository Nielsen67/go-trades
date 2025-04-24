package errormessages

var (
	ErrAlreadyPaid                = "order has been paid"
	ErrPaymentAmountInsufficient  = "payment amount insufficient"
	ErrPaymentNotFound            = "payment not found"
	ErrOrderNotFound              = "order not found"
	ErrOrderDuplicateProduct      = "order contains duplicate product"
	ErrOrderUncancelable          = "order unable to be canceled"
	ErrInvalidOrderStatus         = "invalid order status"
	ErrInvalidInventoryId         = "invalid inventory id"
	ErrInventoryNotFound          = "inventory not found"
	ErrInventoryInvalidStock      = "invalid inventory stock"
	ErrInventoryInsufficientStock = "insufficient stock"
	ErrInventoryStockUpdate       = "failed to update inventory stock"
	ErrInvalidProductId           = "invalid product id"
	ErrProductNotFound            = "product not found"
	ErrProductNameExists          = "product name exists"
	ErrInvalidCategoryId          = "invalid category id"
	ErrCategoryNotFound           = "category not found"
	ErrCategoryNameExists         = "category name exists"
	ErrCategoryCodeExists         = "category code exists"
	ErrUserNotExists              = "user not exists"
)
