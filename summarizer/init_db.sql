CREATE TABLE embeddings (
    id SERIAL PRIMARY KEY,
    filename TEXT NOT NULL,
    chunk_id INT NOT NULL,
    text_chunk TEXT NOT NULL,
    embedding FLOAT8[] NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_filename_chunk ON embeddings(filename, chunk_id);



ALTER TABLE embeddings
ADD CONSTRAINT unique_file_chunk UNIQUE (filename, chunk_id);
