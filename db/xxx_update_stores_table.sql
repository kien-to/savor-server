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