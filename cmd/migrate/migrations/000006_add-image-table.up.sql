CREATE TABLE IF NOT EXISTS Image (
    id SERIAL PRIMARY KEY,
    scan_id INT REFERENCES Scanned_Tire(id),
    data_url VARCHAR(255),
    type VARCHAR(255),
    size INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
