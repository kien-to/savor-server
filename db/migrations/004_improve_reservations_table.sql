-- Migration: Improve reservations table structure
-- This migration enhances the existing reservations table with proper constraints and indexes

-- Drop the existing reservations table if it exists (since it's incomplete)
DROP TABLE IF EXISTS reservations;

-- Create the improved reservations table
CREATE TABLE reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    store_id VARCHAR(36) NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    total_amount DECIMAL(10,2) NOT NULL CHECK (total_amount >= 0),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    payment_id VARCHAR(255),
    pickup_time TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint to stores table
    CONSTRAINT fk_reservations_store 
        FOREIGN KEY (store_id) REFERENCES stores(id) 
        ON DELETE CASCADE
);

-- Create indexes for efficient queries
CREATE INDEX idx_reservations_user_id ON reservations(user_id);
CREATE INDEX idx_reservations_store_id ON reservations(store_id);
CREATE INDEX idx_reservations_status ON reservations(status);
CREATE INDEX idx_reservations_created_at ON reservations(created_at DESC);

-- Create a partial index for active reservations (for better performance)
CREATE INDEX idx_reservations_active 
ON reservations (user_id, created_at DESC) 
WHERE status IN ('pending', 'confirmed');

-- Add constraint for valid status values
ALTER TABLE reservations 
ADD CONSTRAINT check_status 
CHECK (status IN ('pending', 'confirmed', 'completed', 'cancelled', 'expired'));

-- Create trigger function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_reservations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER update_reservations_updated_at
    BEFORE UPDATE ON reservations
    FOR EACH ROW
    EXECUTE FUNCTION update_reservations_updated_at();

-- Add comments for documentation
COMMENT ON TABLE reservations IS 'Stores user reservations/orders for surprise bags';
COMMENT ON COLUMN reservations.id IS 'Unique identifier for the reservation';
COMMENT ON COLUMN reservations.user_id IS 'Firebase UID of the user who made the reservation';
COMMENT ON COLUMN reservations.store_id IS 'ID of the store where the reservation is made';
COMMENT ON COLUMN reservations.quantity IS 'Number of surprise bags reserved';
COMMENT ON COLUMN reservations.total_amount IS 'Total amount paid for the reservation';
COMMENT ON COLUMN reservations.status IS 'Current status of the reservation';
COMMENT ON COLUMN reservations.payment_id IS 'Stripe payment intent ID or custom payment identifier';
COMMENT ON COLUMN reservations.pickup_time IS 'Scheduled pickup time for the reservation';
