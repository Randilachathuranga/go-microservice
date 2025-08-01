package dto

type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryId  uint   `json:"category_id"`
	ImageUrl    string `json:"image_url"`
	Price       uint   `json:"price"`
	Stock       uint   `json:"stock"`
}

type UpdateStockRequest struct {
	Stock uint `json:"stock"`
}
