-- First, add the owner_id column
ALTER TABLE stores 
ADD COLUMN owner_id TEXT;

-- Set a default owner for existing records (replace 'default_owner_id' with an actual user ID)
UPDATE stores 
SET owner_id = 'default_owner_id' 
WHERE owner_id IS NULL;

-- Make owner_id NOT NULL
ALTER TABLE stores 
ALTER COLUMN owner_id SET NOT NULL;

-- Add foreign key constraint to users table
ALTER TABLE stores 
ADD CONSTRAINT fk_owner 
FOREIGN KEY (owner_id) 
REFERENCES users(id);

-- Add store_type column to stores table
ALTER TABLE stores 
ADD COLUMN store_type VARCHAR(255);

-- Update existing records with a default value
UPDATE stores 
SET store_type = 'restaurant' 
WHERE store_type IS NULL;

-- Add business_hours column to stores table
ALTER TABLE stores 
ADD COLUMN business_hours JSONB DEFAULT '[]'::jsonb;


-- add background_url and image_url values for null columns
UPDATE stores 
SET background_url = 'https://vietnamnomad.com/wp-content/uploads/2023/05/What-is-bun-dau-mam-tom.jpg'
WHERE background_url IS NULL;

UPDATE stores 
SET image_url = 'https://vietnamnomad.com/wp-content/uploads/2023/05/What-is-bun-dau-mam-tom.jpg'
WHERE image_url IS NULL;

-- add price value for null columns
UPDATE stores 
SET price = 0
WHERE price IS NULL;

-- add original_price value for null columns
UPDATE stores 
SET original_price = 5
WHERE original_price IS NULL;

-- add rating value for null columns
UPDATE stores 
SET rating = 4.5
WHERE rating IS NULL;

-- add rating_count value for null columns
UPDATE stores 
SET rating_count = 100
WHERE rating_count IS NULL;

-- add is_favorite value for null columns
UPDATE stores 
SET is_favorite = false
WHERE is_favorite IS NULL;

-- add is_open value for null columns
UPDATE stores 
SET is_open = true
WHERE is_open IS NULL;

-- add is_closed value for null columns
UPDATE stores 
SET is_closed = false
WHERE is_closed IS NULL;    

-- add reviews value for null columns
UPDATE stores 
SET reviews = 5
WHERE reviews IS NULL;

-- add reviews_count value for null columns
UPDATE stores 
SET reviews_count = 100
WHERE reviews_count IS NULL;
