ALTER TABLE deployment ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
CREATE INDEX ix_deployment_deleted_at ON deployment(deleted_at);