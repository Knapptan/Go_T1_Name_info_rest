CREATE TABLE persons (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  surname TEXT NOT NULL,
  patronymic TEXT,
  age INT,
  gender TEXT,
  nationality TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_name ON persons (name);
CREATE INDEX idx_surname ON persons (surname);
CREATE INDEX idx_age ON persons (age);
CREATE INDEX idx_created_at ON persons (created_at);