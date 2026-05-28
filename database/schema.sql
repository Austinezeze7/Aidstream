-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) CHECK (role IN ('donor', 'ngo', 'admin', 'verifier')) NOT NULL,
    wallet_address VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Cases Table
CREATE TABLE IF NOT EXISTS cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(50),
    target_amount NUMERIC(12,2) NOT NULL,
    amount_raised NUMERIC(12,2) DEFAULT 0,
    status VARCHAR(20) CHECK (status IN ('pending', 'verified', 'active', 'closed')) DEFAULT 'pending',
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Donations Table
CREATE TABLE IF NOT EXISTS donations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    donor_id UUID REFERENCES users(id),
    case_id UUID REFERENCES cases(id),
    amount NUMERIC(12,2) NOT NULL,
    transaction_hash VARCHAR(200),
    status VARCHAR(20) CHECK (status IN ('pending', 'confirmed', 'failed')) DEFAULT 'pending',
    donated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Verifications Table
CREATE TABLE IF NOT EXISTS verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id UUID REFERENCES cases(id),
    verifier_id UUID REFERENCES users(id),
    status VARCHAR(20) CHECK (status IN ('pending', 'approved', 'rejected')) DEFAULT 'pending',
    notes TEXT,
    verified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);