-- Migration: Allow NULL user_id for guest reservations
-- This migration modifies the reservations table to allow guest users (NULL user_id)

-- Drop the NOT NULL constraint on user_id to allow guest reservations
ALTER TABLE reservations 
ALTER COLUMN user_id DROP NOT NULL;

-- Update the comment to reflect that user_id can be NULL for guests
COMMENT ON COLUMN reservations.user_id IS 'Firebase UID of the user who made the reservation (NULL for guest users)';

-- Create a partial index for guest reservations
CREATE INDEX idx_reservations_guest 
ON reservations (created_at DESC) 
WHERE user_id IS NULL;

-- Update the active reservations index to handle NULL user_id
DROP INDEX IF EXISTS idx_reservations_active;
CREATE INDEX idx_reservations_active 
ON reservations (COALESCE(user_id, ''), created_at DESC) 
WHERE status IN ('pending', 'confirmed');

