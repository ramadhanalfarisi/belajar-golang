package models

import (

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ProductId    uuid.UUID `json:"ID,omitempty" validate:"required"`
	UserId       uuid.UUID `json:"userId,omitempty" validate:"required"`
	ProductName  string    `json:"productName,omitempty" validate:"required"`
	ProductDesc  string    `json:"productDesc,omitempty"`
	ProductPrice float32    `json:"productPrice,omitempty" gorm:"type:decimal(15,2)" validate:"required,min=0,number"`
	ProductImage string    `json:"productImage,omitempty"`
	CreatedAt     string    `json:"createdAt,omitempty"`
	UpdatedAt     *string    `json:"updatedAt,omitempty"`
}

type GetProduct struct {
	ProductId    uuid.UUID `json:"ID,omitempty"`
	ProductName  string    `json:"productName.omitempty"`
	ProductDesc  string    `json:"productDesc,omitempty"`
	ProductPrice float32    `json:"productPrice,omitempty" gorm:"type:decimal(15,2)"`
	ProductImage string    `json:"productImage,omitempty"`
	CreatedAt    string    `json:"createdAt,omitempty"`
	UpdateAt     string    `json:"updateAt,omitempty"`
}

func (product *Product) SelectAllProducts(db *gorm.DB, fields []string, limit int64, offset int64) ([]GetProduct, error) {
	var result []GetProduct
	get_product := db.Table("products").Select(fields).Where("user_id = ?", product.UserId).Order("product_name asc").Limit(int(limit)).Offset(int(offset)).Find(&result)
	if get_product.Error != nil {
		return nil, get_product.Error
	}
	return result, nil
}

func (product *Product) SelectRowProducts(db *gorm.DB, fields []string) int64 {
	var result []GetProduct
	get_product := db.Table("products").Select(fields).Where("user_id = ?", product.UserId).Find(&result)
	if get_product.Error != nil {
		return 0
	}
	num_rows := get_product.RowsAffected
	return num_rows
}

func (product *Product) SelectOneProduct(db *gorm.DB, fields []string) (GetProduct, error) {
	var result GetProduct
	get_product := db.Table("products").Select(fields).Where("user_id = ?", product.UserId).Find(&result)
	if get_product.Error != nil {
		return GetProduct{}, get_product.Error
	}
	return result, nil
}

func (product *Product) InsertProduct(db *gorm.DB) error {
	create := db.Create(&product)
	if create.Error != nil {
		return create.Error
	}
	return nil
}
func (product *Product) UpdateProduct(data Product, db *gorm.DB) (GetProduct, error) {
	updates := db.Model(&Product{}).Where("product_id = ? AND user_id = ?", product.ProductId, product.UserId).Updates(data)
	if updates.Error != nil {
		return GetProduct{}, updates.Error
	}
	return product.SelectOneProduct(db, []string{"product_id", "product_name", "product_desc", "product_price", "product_image", "created_at", "updated_at"})
}

func (product *Product) DeleteProduct(db *gorm.DB) error {
	delete := db.Delete(&Product{}, "product_id = ? AND user_id = ?", product.ProductId, product.UserId)
	if delete.Error != nil {
		return delete.Error
	}
	return nil
}
