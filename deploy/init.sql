CREATE TABLE IF NOT EXISTS warehouses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    availability BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    size TEXT NOT NULL,
    code UUID UNIQUE NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0)
);

CREATE TABLE IF NOT EXISTS warehouse_products(
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER REFERENCES warehouses(id) ON DELETE CASCADE,
    product_code UUID REFERENCES products(code) ON DELETE CASCADE,
    available_quantity INTEGER NOT NULL DEFAULT 0 CHECK(available_quantity >= 0),
    reserved_quantity INTEGER NOT NULL DEFAULT 0 CHECK(reserved_quantity >= 0),
    CONSTRAINT unique_warehouse_product UNIQUE (warehouse_id, product_code)
);

CREATE FUNCTION wareproducts_availability()
RETURNS TRIGGER AS $$
DECLARE
    available BOOLEAN;
    warehouse_id INTEGER;
BEGIN
    CASE TG_OP
        WHEN 'INSERT', 'UPDATE' THEN
            warehouse_id := NEW.warehouse_id;
        WHEN 'DELETE' THEN
            warehouse_id := OLD.warehouse_id;
    END CASE;

    SELECT availability INTO available
    FROM warehouses 
    WHERE id = warehouse_id;

    IF available = false then
        RAISE EXCEPTION USING ERRCODE = 70001, MESSAGE = 'insert/update/delete in no available warehouse';
    END IF;

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;

    RETURN NEW;
END $$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER tr_wareproducts_availability
BEFORE INSERT OR UPDATE OR DELETE ON warehouse_products
FOR EACH ROW
EXECUTE FUNCTION wareproducts_availability();

CREATE OR REPLACE PROCEDURE insertWarehouseProducts(wi INTEGER, pc UUID)
AS $$
DECLARE
    available_quantity INTEGER;
BEGIN
    SELECT quantity INTO available_quantity
    FROM products
    WHERE code = pc;

    IF available_quantity IS NULL THEN
        RAISE EXCEPTION USING ERRCODE = 70002, MESSAGE = 'Dont find proudct with that id'; 
    END IF;

    INSERT INTO warehouse_products (warehouse_id, product_code, available_quantity, reserved_quantity)
    VALUES (wi, pc, available_quantity, 0);
END $$ LANGUAGE plpgsql;
