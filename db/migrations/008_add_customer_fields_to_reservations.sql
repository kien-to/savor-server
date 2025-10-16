-- Add customer fields to reservations table
ALTER TABLE reservations 
ADD COLUMN customer_name VARCHAR(255),
ADD COLUMN customer_email VARCHAR(255),
ADD COLUMN phone_number VARCHAR(50);

-- Add pickup_timestamp if it doesn't exist (from previous migration)
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'reservations' AND column_name = 'pickup_timestamp') THEN
        ALTER TABLE reservations ADD COLUMN pickup_timestamp TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;
