CREATE TABLE IF NOT EXISTS Image (
    id SERIAL PRIMARY KEY,
    inventory_id INT REFERENCES Tire_Inventory(id),
    data_url VARCHAR(255),
    type VARCHAR(255),
    size INT,
    filename VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
