package domain

type Product struct {
	Name     string `json:"name"`
	Size     string `json:"size"`
	Code     string `json:"code"`
	Quantity uint64 `json:"quantity"`
}

type WarehouseProduct struct {
	WarehouseID int64  `json:"warehouse_id"`
	Code        string `json:"code"`
	Quantity    uint64 `json:"quantity"`
	Status      string `json:"status"`
}

type TransferProduct struct {
	WarehouseFromID int64  `json:"warehouse_from_id"`
	WarehouseToID   int64  `json:"warehouse_to_id"`
	Code            string `json:"code"`
	Quantity        uint64 `json:"quantity"`
}

type AddProduct struct {
	Code        string `json:"code"`
	Quantity    uint64 `json:"quantity"`
	WarehouseID int64  `json:"warehouse_id"`
}

type DeleteProduct struct {
	Code string `json:"code"`
}
