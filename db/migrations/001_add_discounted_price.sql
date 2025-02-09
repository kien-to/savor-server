-- Add discounted_price column with default value
ALTER TABLE stores 
ADD discounted_price DECIMAL(10,2) DEFAULT 0.00;

-- First, ensure all stores have a price value
UPDATE stores 
SET price = COALESCE(price, original_price, 5.99)
WHERE price IS NULL;

-- Then, set discounted_price based on price for existing records
UPDATE stores 
SET discounted_price = price
WHERE discounted_price = 0.00;

-- For stores with original_price, calculate a reasonable discount
UPDATE stores 
SET discounted_price = ROUND(original_price * 0.65, 2)
WHERE original_price IS NOT NULL 
AND id NOT IN ('demo-1', 'demo-2');

-- Update the demo data with specific discounted prices
UPDATE stores
SET discounted_price = 19.99,
    original_price = 29.99
WHERE id = 'demo-1';

UPDATE stores
SET discounted_price = 9.99,
    original_price = 15.99
WHERE id = 'demo-2';

-- Finally, make discounted_price NOT NULL
ALTER TABLE stores 
ADD CONSTRAINT discounted_price_not_null 
CHECK (discounted_price IS NOT NULL); 