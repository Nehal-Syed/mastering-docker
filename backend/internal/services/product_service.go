package services

import (
    "database/sql"
    "fmt"

    "mastering-docker/internal/models"
)

type ProductService struct {
    db *sql.DB
}

func NewProductService(db *sql.DB) *ProductService {
    return &ProductService{db: db}
}

func (s *ProductService) CreateProduct(req *models.CreateProductRequest) (*models.Product, error) {
    query := `INSERT INTO products (name, description, price, quantity) VALUES (?, ?, ?, ?)`
    result, err := s.db.Exec(query, req.Name, req.Description, req.Price, req.Quantity)
    if err != nil {
        return nil, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }

    return s.GetProductByID(int(id))
}

func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
    query := `SELECT id, name, description, price, quantity, created_at, updated_at FROM products WHERE id = ?`
    var product models.Product
    err := s.db.QueryRow(query, id).Scan(
        &product.ID, &product.Name, &product.Description,
        &product.Price, &product.Quantity, &product.CreatedAt, &product.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("product not found")
    }
    if err != nil {
        return nil, err
    }
    return &product, nil
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
    query := `SELECT id, name, description, price, quantity, created_at, updated_at FROM products ORDER BY created_at DESC`
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var product models.Product
        err := rows.Scan(
            &product.ID, &product.Name, &product.Description,
            &product.Price, &product.Quantity, &product.CreatedAt, &product.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        products = append(products, product)
    }
    return products, nil
}

func (s *ProductService) UpdateProduct(id int, req *models.UpdateProductRequest) (*models.Product, error) {
    query := `UPDATE products SET name = COALESCE(?, name), description = COALESCE(?, description), price = COALESCE(?, price), quantity = COALESCE(?, quantity) WHERE id = ?`
    _, err := s.db.Exec(query, req.Name, req.Description, req.Price, req.Quantity, id)
    if err != nil {
        return nil, err
    }
    return s.GetProductByID(id)
}

func (s *ProductService) DeleteProduct(id int) error {
    query := `DELETE FROM products WHERE id = ?`
    result, err := s.db.Exec(query, id)
    if err != nil {
        return err
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return fmt.Errorf("product not found")
    }
    return nil
}