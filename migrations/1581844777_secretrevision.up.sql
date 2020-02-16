ALTER TABLE secret_meta ADD CONSTRAINT stagefk FOREIGN KEY (stage_name) REFERENCES stage(name);
ALTER TABLE secret_meta ADD COLUMN revision integer NOT NULL DEFAULT(0);
ALTER TABLE secret_meta ADD COLUMN committed_revision integer NOT NULL DEFAULT(0);
ALTER TABLE secret_meta DROP COLUMN doc;
CREATE INDEX ix_secretmeta_revision ON secret_meta(revision);
CREATE INDEX ix_secretmeta_committed_revision ON secret_meta(committed_revision);