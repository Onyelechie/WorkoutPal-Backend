-- Create follow_requests table for managing pending follow requests to private profiles
CREATE TABLE IF NOT EXISTS follow_requests (
    id SERIAL PRIMARY KEY,
    requester_id INTEGER NOT NULL,
    requested_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, accepted, rejected
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(requester_id, requested_id),
    FOREIGN KEY (requester_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (requested_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create index for faster queries
CREATE INDEX idx_follow_requests_requested_id ON follow_requests(requested_id);
CREATE INDEX idx_follow_requests_requester_id ON follow_requests(requester_id);
CREATE INDEX idx_follow_requests_status ON follow_requests(status);
