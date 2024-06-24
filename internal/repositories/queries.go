package repositories

const (
	GetCategoryByID     = "SELECT name, product_id, created_at, updated_at FROM category WHERE id = $1"
	UpdateCategory      = "UPDATE category SET name = $2, product_id = $3, updated_at = $4 WHERE id = $1"
	CreateCategory      = "INSERT INTO category (name, product_id) VALUES ($1, $2)"
	CheckCategoryExists = "SELECT EXISTS (SELECT 1 FROM category WHERE name = $1)"
	DeleteCategory      = "DELETE FROM category WHERE id = $1"
	GetProduct          = "SELECT name, category_id, created_at, updated_at FROM products WHERE id = $1"
	UpdateProduct       = "UPDATE products SET name = $2, category_id = $3, updated_at = $4 WHERE id = $1"
	CreateProduct       = "INSERT INTO products (name, category_id) VALUES ($1, $2)"
	CheckProductExists  = "SELECT EXISTS (SELECT 1 FROM products WHERE name = $1)"
	DeleteProduct       = "DELETE FROM products WHERE id = $1"
	AddCustomer         = "INSERT INTO customer (username, email, password, role) VALUES ($1, $2, $3, $4)"
	GetUserByEmail      = "SELECT id, username, email, password, role FROM customer WHERE email = $1"
	GetUserByUsername   = "SELECT username, email, role FROM customer WHERE username = $1"
	GetAllUsers         = "SELECT email, role FROM customer"
	UpdateUser          = "UPDATE customer SET username = $2, email = $3, role = $4 WHERE id = $1"
	DeleteUser          = "DELETE FROM customer WHERE id = $1"
	CheckUserExists     = "SELECT EXISTS (SELECT 1 FROM customer WHERE email = $1)"
)
