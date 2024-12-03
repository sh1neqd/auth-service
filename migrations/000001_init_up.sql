CREATE TABLE refresh_tokens (
                                id UUID PRIMARY KEY,
                                user_id UUID NOT NULL,
                                token_hash TEXT NOT NULL,
                                issued_at TIMESTAMPTZ NOT NULL,
                                expires_at TIMESTAMPTZ NOT NULL
);

