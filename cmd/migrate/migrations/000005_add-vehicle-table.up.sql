CREATE TABLE IF NOT EXISTS Vehicle (
  id SERIAL PRIMARY KEY,
  scan_id INT REFERENCES Tire_Inventory(id) NOT NULL,
  license_plate VARCHAR(255) NOT NULL,
  color VARCHAR(255),
  brand VARCHAR(255),
  model VARCHAR(255),
  year VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);
