package domain

type Warehouse struct {
	Name         string `json:"name"`
	Availability bool   `json:"availability"`
}

type GetFromWarehouse struct {
	WarehouseID int64 `json:"warehouse_id"`
}
