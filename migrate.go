package paystore

var createTableBalance = `CREATE TABLE balance (
		uuid VARCHAR(255) PRIMARY KEY,
		randid VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		balance BIGINT NOT NULL DEFAULT 0,
		last_receive TIMESTAMP WITH TIME ZONE,
		last_withdraw TIMESTAMP WITH TIME ZONE,
		income_accumulation BIGINT NOT NULL DEFAULT 0,
		withdraw_accumulation BIGINT NOT NULL DEFAULT 0,
		currency VARCHAR(3) NOT NULL,
		active BOOL NOT NULL DEFAULT true,
		external_id VARCHAR(255),
		organization_uuid UUID NOT NULL
    );

    -- Indexes for better query performance
    CREATE INDEX idx_accounts_organization_uuid (organization_uuid);
    CREATE INDEX idx_accounts_external_id (external_id);
`

var createTableQuery = `
		CREATE TABLE payment (
			-- Fields from Record
			uuid VARCHAR(255) PRIMARY KEY,
			randid VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		
			-- Fields from Payment
			amount BIGINT NOT NULL,
			balance_before_payment BIGINT NOT NULL,
			balance_after_payment BIGINT NOT NULL,
			balance_uuid VARCHAR(255) NOT NULL,
			organization_uuid VARCHAR(255) NOT NULL,
			vendor_record_id VARCHAR(255) NOT NULL,
			status VARCHAR(20) NOT NULL,
			hash VARCHAR(255) NOT NULL
		);
		
		-- Indexes for common queries
		CREATE INDEX idx_payments_balance_uuid ON payment(balance_uuid);
		CREATE INDEX idx_payments_created_at ON payment(created_at);
		CREATE INDEX idx_payments_hash ON payment(hash);
`

var createTableOrganization = `CREATE TABLE organization (
		uuid VARCHAR(255) PRIMARY KEY,
		randid VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		name VARCHAR(255) NOT NULL,
		slug VARCHAR(255) NOT NULL
	);
`

var createTableTransaction = `CREATE TABLE transaction (
		uuid VARCHAR(255) PRIMARY KEY, 
		randid VARCHAR(255) NOT NULL, 
		created_at TIMESTAMP NOT NULL DEFAULT NOW(), 
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(), 
		transaction_type VARCHAR(255) NOT NULL, 
		record_uuid VARCHAR(255) NOT NULL, 
		balance_uuid VARCHAR(255) NOT NULL
 	);`
