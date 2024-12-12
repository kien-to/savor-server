-- First, delete existing highlights to avoid duplicates
DELETE FROM store_highlights;

-- Insert exactly 3 highlights for each store
INSERT INTO store_highlights (store_id, highlight)
SELECT 
    id,
    unnest(ARRAY[
        CASE 
            WHEN random() < 0.33 THEN 'Fresh baked goods'
            WHEN random() < 0.66 THEN 'Local favorite'
            ELSE 'Best seller'
        END,
        CASE 
            WHEN random() < 0.33 THEN 'Great value'
            WHEN random() < 0.66 THEN 'Popular item'
            ELSE 'Staff pick'
        END,
        CASE
            WHEN random() < 0.33 THEN 'Trending now'
            WHEN random() < 0.66 THEN 'Limited time'
            ELSE 'Must try'
        END
    ])
FROM stores
WHERE id IN ('1', '2', '3', '4', '5', '6', 
            'sf001', 'sf002', 'sf003', 'sf004', 'sf005');
