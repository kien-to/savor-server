-- Fix any remaining NULL store_type values
-- This ensures all stores have a default store_type value

UPDATE stores 
SET store_type = 'restaurant' 
WHERE store_type IS NULL;

-- Optional: Add a NOT NULL constraint if desired
-- ALTER TABLE stores 
-- ALTER COLUMN store_type SET NOT NULL;

-- Verify the update
SELECT COUNT(*) as null_store_type_count 
FROM stores 
WHERE store_type IS NULL;
