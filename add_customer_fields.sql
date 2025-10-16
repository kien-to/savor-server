-- Add customer fields to reservations table if they don't exist
DO $$ 
BEGIN
    -- Add customer_name column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'reservations' AND column_name = 'customer_name') THEN
        ALTER TABLE reservations ADD COLUMN customer_name VARCHAR(255);
    END IF;
    
    -- Add customer_email column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'reservations' AND column_name = 'customer_email') THEN
        ALTER TABLE reservations ADD COLUMN customer_email VARCHAR(255);
    END IF;
    
    -- Add phone_number column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'reservations' AND column_name = 'phone_number') THEN
        ALTER TABLE reservations ADD COLUMN phone_number VARCHAR(50);
    END IF;
    
    -- Add pickup_timestamp column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'reservations' AND column_name = 'pickup_timestamp') THEN
        ALTER TABLE reservations ADD COLUMN pickup_timestamp TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Update existing reservations to use email as customer_name if customer_name is null
UPDATE reservations 
SET customer_name = COALESCE(customer_name, (SELECT email FROM users WHERE users.id = reservations.user_id))
WHERE customer_name IS NULL OR customer_name = '';

-- Update existing reservations to use email as customer_email if customer_email is null
UPDATE reservations 
SET customer_email = COALESCE(customer_email, (SELECT email FROM users WHERE users.id = reservations.user_id))
WHERE customer_email IS NULL OR customer_email = '';

-- Set pickup_timestamp to created_at + 2 hours for existing reservations if null
UPDATE reservations 
SET pickup_timestamp = COALESCE(pickup_timestamp, created_at + INTERVAL '2 hours')
WHERE pickup_timestamp IS NULL;
