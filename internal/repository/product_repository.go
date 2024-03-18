package repository

import (
	"errors"
	"fmt"
	"golang_socmed/internal/entity"
	"golang_socmed/internal/model"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ProductRepository struct {
	DB *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) List(filter *model.ProductFilter, userId string) ([]entity.Product, error) {
	tx, _ := r.DB.Beginx()
	defer tx.Rollback()
	// product := []entity.Product{}

	query := `SELECT * FROM products`
	var filterValues []interface{}

	// Conditionally append filters
	if *filter.Keyword != "" {
		query += ` WHERE name LIKE CONCAT('%', $1::TEXT, '%')`
		filterValues = append(filterValues, *filter.Keyword)
	}

	if *filter.Condition != "" {
		if len(filterValues) > 0 {
			query += ` AND `
		} else {
			query += ` WHERE `
		}
		query += ` condition = $` + strconv.Itoa(len(filterValues)+1)
		filterValues = append(filterValues, *filter.Condition)
	}

	// Add sorting if SortField and SortOrder are provided
	if *filter.SortBy != "" && *filter.OrderBy != "" {
		query += ` ORDER BY ` + *filter.SortBy + ` ` + *filter.OrderBy
	}

	// Add price range if MinPrice and MaxPrice are provided
	if *filter.MinPrice != 0 && *filter.MaxPrice != 0 {
		if len(filterValues) > 0 {
			query += ` AND `
		} else {
			query += ` WHERE `
		}
		query += `price >= $` + strconv.Itoa(len(filterValues)+1) + ` AND price <= $` + strconv.Itoa(len(filterValues)+2)
		filterValues = append(filterValues, filter.MinPrice, filter.MaxPrice)
	}

	if filter.Tags != nil && len(*filter.Tags) > 0 {
		if len(filterValues) > 0 {
			query += ` AND `
		} else {
			query += ` WHERE `
		}

		tags := strings.Join(*filter.Tags, "','")
		query += "ARRAY['" + tags + "']::text[] <@ tags::text[]"
	}

	if filter.UserOnly != nil && *filter.UserOnly {
		if len(filterValues) > 0 {
			query += ` AND `
		} else {
			query += ` WHERE `
		}
		query += ` userId = $` + strconv.Itoa(len(filterValues)+1)
		filterValues = append(filterValues, userId)
	}

	if filter.ShowEmptyStock != nil && *filter.ShowEmptyStock {
		if len(filterValues) > 0 {
			query += ` AND `
		} else {
			query += ` WHERE `
		}
		query += ` stock = $` + strconv.Itoa(len(filterValues)+1)
		filterValues = append(filterValues, 0)
	}

	if *filter.Limit != 0 {
		query += fmt.Sprintf(" LIMIT $%s", strconv.Itoa(len(filterValues)+1))
		filterValues = append(filterValues, filter.Limit)
	}

	if *filter.Offset != 0 {
		query += fmt.Sprintf(" OFFSET $%s", strconv.Itoa(len(filterValues)+1))
		filterValues = append(filterValues, filter.Offset)
	}

	// Execute the query
	rows, err := r.DB.Query(query, filterValues...)
	if err != nil {
		// Handle error
		return nil, err
	}

	// Slice to hold the fetched products
	var products []entity.Product

	// Loop through the rows and scan each product into the slice
	for rows.Next() {
		var product entity.Product
		// Use pq.Array to scan the Tags column into the product.Tags slice
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.ImageUrl, &product.Stock, &product.Condition, &product.IsPurchasable, pq.Array(&product.Tags), &product.UserId)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// if err := query.Find(&product).Error; err != nil {
	// 	return nil, err
	// }

	return products, nil
}

func (r *ProductRepository) Get(id string, product *entity.Product) (entity.Product, error) {
	tx, _ := r.DB.Beginx()
	defer tx.Rollback()
	query := `SELECT * FROM products WHERE id = $1`
	productData := entity.Product{}

	// Execute the query
	rows, err := r.DB.Query(query, id)
	if err != nil {
		// Handle error
		return productData, err
	}

	// Loop through the rows and scan each product into the slice
	for rows.Next() {
		var product entity.Product
		// Use pq.Array to scan the Tags column into the product.Tags slice
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.ImageUrl, &product.Stock, &product.Condition, &product.IsPurchasable, pq.Array(&product.Tags), &product.UserId)
		if err != nil {
			return productData, err
		}
		productData = product
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return productData, err
	}
	// Check for no results
	if productData.Name == "" {
		return productData, errors.New("No Data Found")
	}

	return productData, nil
}

func (r *ProductRepository) Create(request *entity.Product) error {
	query := `INSERT INTO products VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.DB.Exec(query, request.ID, request.Name, request.Price, request.ImageUrl, request.Stock, request.Condition, request.IsPurchasable, pq.Array(request.Tags), request.UserId)
	return err
}

func (r *ProductRepository) Update(id string, request *entity.Product) error {
	query := `UPDATE products SET name = $2, price = $3, imageUrl = $4, condition = $5, isPurchasable = $6,tags = $7 WHERE id = $1`

	_, err := r.DB.Exec(query, id, request.Name, request.Price, request.ImageUrl, request.Condition, request.IsPurchasable, pq.Array(request.Tags))
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) UpdateStock(id string, request *entity.Product) error {
	query := `UPDATE products SET stock = $2 WHERE id = $1`

	_, err := r.DB.Exec(query, id, request.Stock)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) Delete(id string) error {
	query := `DELETE FROM products WHERE id = $1`
	// Send query to database.
	_, err := r.DB.Exec(query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

func (r *ProductRepository) Buy(request *entity.Payment) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Lock the table
	// if _, err := tx.Exec("LOCK TABLE products IN SHARE MODE"); err != nil {
	// 	return err
	// }

	product := new(entity.Product)
	query := `SELECT stock FROM products WHERE id = $1 FOR NO KEY UPDATE`
	err = tx.Get(product, query, request.ProductId)
	if err != nil {
		return err
	}

	if request.Quantity > int(product.Stock) {
		return &fiber.Error{
			Code:    400,
			Message: "insufficient quantity",
		}
	}

	insertQuery := `INSERT INTO payments VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(insertQuery, request.ID, request.BankAccountId, request.PaymentProofImageUrl, request.ProductId, request.Quantity)
	if err != nil {
		return err
	}

	// newStock := int(product.Stock) - request.Quantity
	updateQuery := `UPDATE products SET stock = stock - $1 WHERE id = $2`
	// _, err = tx.Exec(updateQuery, newStock, request.ProductId)
	_, err = tx.Exec(updateQuery, request.Quantity, request.ProductId)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
