CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);


CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NOT NULL REFERENCES users(id),
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    distance DECIMAL(10, 2) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    truck_size VARCHAR(10) NOT NULL,
    weight DECIMAL(10, 2) NOT NULL,
    volume DECIMAL(10, 2) NOT NULL,
    payout DECIMAL(10, 2) NOT NULL,
    is_new BOOLEAN DEFAULT TRUE,
    is_claimed BOOLEAN DEFAULT FALSE,
    is_verified BOOLEAN DEFAULT FALSE,
    is_protected BOOLEAN DEFAULT TRUE,
    is_escrow BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_jobs_user_id ON jobs(user_id);
CREATE INDEX idx_jobs_origin ON jobs(origin);
CREATE INDEX idx_jobs_destination ON jobs(destination);
CREATE INDEX idx_jobs_start_date ON jobs(start_date);
CREATE INDEX idx_jobs_truck_size ON jobs(truck_size);
CREATE INDEX idx_jobs_is_claimed ON jobs(is_claimed);


CREATE TABLE IF NOT EXISTS job_claims (
    id SERIAL PRIMARY KEY,
    job_id INTEGER NOT NULL REFERENCES jobs(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending' NOT NULL,
    UNIQUE(job_id, user_id)
);

CREATE INDEX idx_job_claims_job_id ON job_claims(job_id);
CREATE INDEX idx_job_claims_user_id ON job_claims(user_id);