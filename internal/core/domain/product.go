package domain

type Product struct {
	BaseModel
	Name           string             `json:"name" gorm:"type:varchar(50);not null"`
	Description    string             `json:"description" gorm:"type:varchar(255);not null"`
	BasePrice      int                `json:"base_price" gorm:"type:int;not null"`
	PricingDetails *map[string]string `json:"pricing_details" gorm:"type:jsonb"`
}

type ProductMeta struct {
	BaseAutoIncModel
	ProductID string `json:"product_id" gorm:"type:varchar;not null"`
	ImageUrl  string `json:"image_url" gorm:"type:varchar;not null"`
}
