package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ETNA10-7/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]*types.Product, 0)

	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (s *Store) GetProductByID(productID int) (*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE id = $1 ", productID)
	if err != nil {
		return nil, err
	}

	p := new(types.Product)

	for rows.Next() {
		p, err = scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

// func (s *Store) GetProductsByID(productIds []int) ([]types.ProductStock, error) {
// 	// placeholders := strings.Repeat(",?", len(productIDs)-1)
// 	// query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)
// 	placeholders := make([]string, len(productIds))
// 	for i := range productIds {
// 		placeholders[i] = fmt.Sprintf("$%d", i+1)
// 	}

// 	query := fmt.Sprintf("SELECT * FROM product_stock WHERE id IN (%s)", strings.Join(placeholders, ","))

// 	//Convert productIDs to []interface{}
// 	args := make([]interface{}, len(productIds))
// 	for i, v := range productIds {
// 		args[i] = v
// 	}

// 	rows, err := s.db.Query(query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	productsq := []types.ProductStock{}

// 	for rows.Next() {
// 		p, err := scanRowsIntoProductStock(rows)

// 		if err != nil {
// 			return nil, err
// 		}
// 		productsq = append(productsq, *p)
// 	}
// 	return productsq, nil

// }

func (s *Store) GetProductsByID(productIds []int) ([]types.Product, []types.ProductStock, error) {
	// placeholders := strings.Repeat(",?", len(productIDs)-1)
	// query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)
	placeholders := make([]string, len(productIds))
	for i := range productIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (%s)", strings.Join(placeholders, ","))
	queryq := fmt.Sprintf("SELECT * FROM product_stock WHERE product_id IN (%s)", strings.Join(placeholders, ","))

	//Convert productIDs to []interface{}
	args := make([]interface{}, len(productIds))
	for i, v := range productIds {
		args[i] = v
	}

	// products types.Product
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}

	// productq types.ProductStock
	rowsq, err := s.db.Query(queryq, args...)
	if err != nil {
		return nil, nil, err
	}

	products := []types.Product{}
	productsq := []types.ProductStock{}

	//To get Products by ID types.Products
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)

		if err != nil {
			return nil, nil, err
		}
		products = append(products, *p)
	}

	//To get ProductStock by ID types.ProductStock
	for rowsq.Next() {
		p, err := scanRowsIntoProductStock(rowsq)

		if err != nil {
			return nil, nil, err
		}
		productsq = append(productsq, *p)
	}
	//log.Printf()
	return products, productsq, nil

}

// func (s *Store) CreateProduct(product types.CreateProductPayload) error {
// 	_, err := s.db.Exec("INSERT INTO products (name, price, image, description, quantity) VALUES ($1, $2, $3, $4, $5)", product.Name, product.Price, product.Image, product.Description, product.Quantity)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (s *Store) CreateProduct(product types.CreateProductPayload) error {

	stock := product.Quantity

	tx, err := s.db.Begin()
	if err != nil {
		//http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Check if the product already exists
	var existingProductID int
	err = tx.QueryRow("SELECT id FROM products WHERE name = $1", product.Name).Scan(&existingProductID)
	if err != nil && err != sql.ErrNoRows {
		//http.Error(w, "Database error", http.StatusInternalServerError)
		return err
	}

	if existingProductID == 0 {
		// Product doesn't exist, create it
		_, err = tx.Exec("INSERT INTO products (name, description, image, price) VALUES ($1, $2, $3, $4)",
			product.Name, product.Description, product.Image, product.Price)
		if err != nil {
			//http.Error(w, "Failed to create product", http.StatusInternalServerError)
			return err
		}

		// Get the new product ID
		err = tx.QueryRow("SELECT LASTVAL()").Scan(&existingProductID)
		if err != nil {
			//http.Error(w, "Failed to retrieve product ID", http.StatusInternalServerError)
			return err
		}

		// Insert initial stock
		_, err = tx.Exec("INSERT INTO product_stock (product_id, stock) VALUES ($1, $2)", existingProductID, stock)
		if err != nil {
			//http.Error(w, "Failed to initialize stock", http.StatusInternalServerError)
			return err
		}
	} else {
		// Product exists, update stock (restocking)
		_, err = tx.Exec("UPDATE product_stock SET stock = stock + $1 WHERE product_id = $2", stock, existingProductID)
		if err != nil {
			//http.Error(w, "Failed to update stock", http.StatusInternalServerError)
			return err
		}
	}

	//w.WriteHeader(http.StatusOK)
	return nil
}

// func (s *Store) UpdateProductStock(product types.ProductStock) error {
// 	_, err := s.db.Exec("UPDATE product_stock SET stock = $1 WHERE product_id = $2", product.Stock, product.ProductID)
// 	if err != nil {
// 		return err
// 	}

//		return nil
//	}
func (s *Store) UpdateAndRestock(cartItems []types.CartCheckoutItem) error {
	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Check stock for each cart item
	for _, item := range cartItems {
		var stock int
		err := tx.QueryRow("SELECT stock FROM product_stock WHERE product_id = $1", item.ProductID).Scan(&stock)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("product ID %d not found", item.ProductID)
			}
			return fmt.Errorf("error checking stock for product ID %d: %v", item.ProductID, err)
		}
		if stock < item.Quantity {
			return fmt.Errorf("insufficient stock for product ID %d", item.ProductID)
		}

		// Update the stock
		_, err = tx.Exec("UPDATE product_stock SET stock = stock - $1 WHERE product_id = $2 AND stock >= $1", item.Quantity, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to update stock for product ID %d: %v", item.ProductID, err)
		}
	}

	return nil
}

// func (s *Store) UpdateProductStockTx(tx *sql.Tx, productID, reduceBy int) error {
// 	_, err := tx.Exec("UPDATE product_stock SET stock = stock - $1 WHERE product_id = $2 AND stock >= $1", reduceBy, productID)
// 	if err != nil {
// 		return fmt.Errorf("failed to update stock for product ID %d: %v", productID, err)
// 	}
// 	return nil
// }

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		//&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func scanRowsIntoProductStock(rows *sql.Rows) (*types.ProductStock, error) {
	product := new(types.ProductStock)

	err := rows.Scan(
		&product.ProductID,
		&product.Stock,
	)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// func (s *Store) CreateOrRestockProduct(w http.ResponseWriter, r *http.Request) {
// 	var payload ProductWithStock
// 	err := json.NewDecoder(r.Body).Decode(&payload)
// 	if err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	product := payload.Product
// 	stock := payload.Stock

// 	tx, err := s.db.Begin()
// 	if err != nil {
// 		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
// 		return
// 	}
// 	defer func() {
// 		if p := recover(); p != nil {
// 			tx.Rollback()
// 			panic(p)
// 		} else if err != nil {
// 			tx.Rollback()
// 		} else {
// 			err = tx.Commit()
// 		}
// 	}()

// 	// Check if the product already exists
// 	var existingProductID int
// 	err = tx.QueryRow("SELECT id FROM products WHERE id = $1", product.ID).Scan(&existingProductID)
// 	if err != nil && err != sql.ErrNoRows {
// 		http.Error(w, "Database error", http.StatusInternalServerError)
// 		return
// 	}

// 	if existingProductID == 0 {
// 		// Product doesn't exist, create it
// 		_, err = tx.Exec("INSERT INTO products (name, description, image, price, created_at) VALUES ($1, $2, $3, $4, $5)",
// 			product.Name, product.Description, product.Image, product.Price, product.CreatedAt)
// 		if err != nil {
// 			http.Error(w, "Failed to create product", http.StatusInternalServerError)
// 			return
// 		}

// 		// Get the new product ID
// 		err = tx.QueryRow("SELECT LASTVAL()").Scan(&existingProductID)
// 		if err != nil {
// 			http.Error(w, "Failed to retrieve product ID", http.StatusInternalServerError)
// 			return
// 		}

// 		// Insert initial stock
// 		_, err = tx.Exec("INSERT INTO product_stock (product_id, stock) VALUES ($1, $2)", existingProductID, stock)
// 		if err != nil {
// 			http.Error(w, "Failed to initialize stock", http.StatusInternalServerError)
// 			return
// 		}
// 	} else {
// 		// Product exists, update stock (restocking)
// 		_, err = tx.Exec("UPDATE product_stock SET stock = stock + $1 WHERE product_id = $2", stock, existingProductID)
// 		if err != nil {
// 			http.Error(w, "Failed to update stock", http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

// func (s *Store) CheckoutCart(w http.ResponseWriter, r *http.Request) {
// 	var items []CartCheckoutItem
// 	err := json.NewDecoder(r.Body).Decode(&items)
// 	if err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	tx, err := s.db.Begin()
// 	if err != nil {
// 		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
// 		return
// 	}
// 	defer func() {
// 		if p := recover(); p != nil {
// 			tx.Rollback()
// 			panic(p)
// 		} else if err != nil {
// 			tx.Rollback()
// 		} else {
// 			err = tx.Commit()
// 		}
// 	}()

// 	for _, item := range items {
// 		var availableStock int

// 		// Check current stock
// 		err = tx.QueryRow("SELECT stock FROM product_stock WHERE product_id = $1", item.ProductID).Scan(&availableStock)
// 		if err != nil {
// 			http.Error(w, "Product not found or database error", http.StatusBadRequest)
// 			return
// 		}

// 		if availableStock < item.Quantity {
// 			http.Error(w, "Insufficient stock for product ID "+strconv.Itoa(item.ProductID), http.StatusBadRequest)
// 			return
// 		}

// 		// Deduct stock
// 		_, err = tx.Exec("UPDATE product_stock SET stock = stock - $1 WHERE product_id = $2", item.Quantity, item.ProductID)
// 		if err != nil {
// 			http.Error(w, "Failed to update stock", http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	w.WriteHeader(http.StatusOK)
// }
