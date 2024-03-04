CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(127) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO projects (name) VALUES ('test-service');

CREATE TABLE goods (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    priority SERIAL NOT NULL,
    removed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE INDEX idx_project_id_hash ON goods USING HASH (project_id);
CREATE INDEX idx_goods_name ON goods (name COLLATE pg_catalog."default" ASC);