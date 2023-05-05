CREATE TABLE IF NOT EXISTS categories (
	category_id INT PRIMARY KEY,
	category_name VARCHAR (255) NOT NULL
);
CREATE TABLE IF NOT EXISTS brands (
	brand_id INT PRIMARY KEY,
	brand_name VARCHAR (255) NOT NULL
);
CREATE TABLE IF NOT EXISTS customers (
	customer_id SERIAL PRIMARY KEY,
	first_name VARCHAR (255) NOT NULL,
	last_name VARCHAR (255) NOT NULL,
	phone VARCHAR (25),
	email VARCHAR (255) NOT NULL,
	street VARCHAR (255),
	city VARCHAR (50),
	state VARCHAR (25),
	zip_code NUMERIC
);

-- select products.list_price from orders left join order_items ON orders.order_id=order_items.order_id left JOIN products on order_items.product_id=products.product_id where orders.order_id=1;
CREATE TABLE IF NOT EXISTS products (
	product_id INT PRIMARY KEY,
	product_name VARCHAR (255) NOT NULL,
	brand_id INT NOT NULL REFERENCES brands (brand_id) ON DELETE CASCADE ON UPDATE NO ACTION,
	model_year SMALLINT NOT NULL,
	list_price DECIMAL (10, 2) NOT NULL,
	category_id int REFERENCES categories (category_id) ON DELETE CASCADE ON UPDATE NO ACTION
);
CREATE TABLE IF NOT EXISTS stores (
	store_id SERIAL PRIMARY KEY,
	store_name VARCHAR (255) NOT NULL,
	phone VARCHAR (25),
	email VARCHAR (255),
	street VARCHAR (255),
	city VARCHAR (255),
	state VARCHAR (10),
	zip_code VARCHAR (5)
);
CREATE TABLE IF NOT EXISTS stocks (
	store_id INT REFERENCES stores (store_id) ON DELETE CASCADE ON UPDATE NO ACTION,
	product_id INT REFERENCES products (product_id) ON DELETE CASCADE ON UPDATE NO ACTION,
	quantity INT,
	PRIMARY KEY (store_id, product_id)
);
CREATE TABLE IF NOT EXISTS staffs (
	staff_id INT PRIMARY KEY,
	first_name VARCHAR (50) NOT NULL,
	last_name VARCHAR (50) NOT NULL,
	email VARCHAR (255) NOT NULL UNIQUE,
	phone VARCHAR (25),
	active SMALLINT NOT NULL,
	store_id INT NOT NULL REFERENCES stores (store_id) ON DELETE CASCADE ON UPDATE NO ACTION,
	manager_id INT REFERENCES staffs (staff_id) ON DELETE CASCADE ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS promocode(
	id INT PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	discount INT NOT NULL,
	discount_type VARCHAR(100) NOT NULL,
	order_limit_price FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
	order_id INT PRIMARY KEY,
	customer_id INT REFERENCES customers (customer_id) ON DELETE CASCADE ON UPDATE NO ACTION,
	order_status SMALLINT NOT NULL,
	-- Order status: 1 = Pending; 2 = Processing; 3 = Rejected; 4 = Completed
	order_date DATE NOT NULL,
	required_date DATE NOT NULL,
	shipped_date DATE,
	promocode_id INTEGER REFERENCES promocode(id) ON DELETE CASCADE ON UPDATE NO ACTION,
	store_id INT NOT NULL REFERENCES stores (store_id) ON DELETE CASCADE ON UPDATE NO ACTION,
	staff_id INT NOT NULL REFERENCES staffs (staff_id) ON DELETE CASCADE ON UPDATE NO ACTION
);
CREATE TABLE IF NOT EXISTS order_items (
	order_id INT REFERENCES orders (order_id) ON DELETE CASCADE ON UPDATE NO ACTION,
	item_id INT,
	product_id INT NOT NULL REFERENCES products(product_id) ON DELETE CASCADE ON UPDATE NO ACTION,
	quantity INT NOT NULL,
	list_price DECIMAL (10, 2) NOT NULL,
	discount DECIMAL (4, 2) NOT NULL DEFAULT 0
);


