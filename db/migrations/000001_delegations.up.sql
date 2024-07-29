CREATE TABLE delegations (
    id BIGINT NOT NULL,
    delegator VARCHAR(64) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    amount BIGINT NOT NULL,
    level BIGINT NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (id, delegator, timestamp, amount, level)
);

CREATE INDEX idx_delegations_timestamp ON delegations USING btree (timestamp);

