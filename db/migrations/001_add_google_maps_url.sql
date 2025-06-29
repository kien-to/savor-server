-- Migration: Add google_maps_url column to stores table
-- Date: 2024-01-01

ALTER TABLE stores 
ADD COLUMN google_maps_url TEXT;

-- Add index for better performance when querying by google_maps_url
CREATE INDEX idx_stores_google_maps_url ON stores(google_maps_url) WHERE google_maps_url IS NOT NULL; 