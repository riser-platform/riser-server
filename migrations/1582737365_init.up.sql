CREATE TABLE namespace (
  name character varying(63) NOT NULL,
  PRIMARY KEY(name)
);

CREATE TABLE stage (
  name character varying(63) NOT NULL,
  doc jsonb NOT NULL,
  PRIMARY KEY(name)
);

CREATE TABLE app
(
  id uuid NOT NULL,
  name character varying(63) NOT NULL,
-- TODO: namespace
  --namespace character varying(63) NOT NULL REFERENCES namespace(name),
  PRIMARY KEY(id)
);

-- TODO: namespace
-- CREATE UNIQUE INDEX ix_app_name ON app(name, namespace);
CREATE UNIQUE INDEX ix_app_name ON app(name);

CREATE TABLE deployment
(
  name character varying(63) NOT NULL,
  stage_name character varying(63) NOT NULL REFERENCES stage(name),
  app_id uuid NOT NULL REFERENCES app(id),
  riser_revision integer NOT NULL DEFAULT(0),
  deleted_at TIMESTAMP WITH TIME ZONE,
  doc jsonb NOT NULL,
  PRIMARY KEY (name,stage_name)
);

CREATE UNIQUE INDEX ix_deployment ON deployment(name,stage_name);
CREATE INDEX ix_deployment_riser_revision ON deployment(riser_revision);
CREATE INDEX ix_deployment_deleted_at ON deployment(deleted_at);

CREATE TABLE secret_meta (
  name character varying(63) NOT NULL,
  app_id uuid NOT NULL REFERENCES app(id),
  stage_name character varying(63) NOT NULL REFERENCES stage(name),
  revision integer NOT NULL DEFAULT(0),
  committed_revision integer NOT NULL DEFAULT(0),
  PRIMARY KEY (name, app_id, stage_name)
);

CREATE INDEX ix_secretmeta_revision ON secret_meta(revision);
CREATE INDEX ix_secretmeta_committed_revision ON secret_meta(committed_revision);

/* "user" is a reserved word in Postgress. Easier to just use riser_user. The domain will still call this resource a "user" */
CREATE TABLE riser_user
(
-- TODO: uuid
  id serial,
  username character varying(32) NOT NULL,
  doc jsonb NOT NULL,
  PRIMARY KEY(id)
);

CREATE UNIQUE INDEX ix_riser_user_username ON riser_user(username);

CREATE TABLE apikey
(
-- TODO: uuid
  id serial,
  riser_user_id integer NOT NULL REFERENCES riser_user(id),
  key_hash bytea NOT NULL,
  PRIMARY KEY(id)
);

CREATE INDEX ix_userlogin_user_id ON apikey(riser_user_id);
CREATE UNIQUE INDEX ix_userlogin_key_hash ON apikey(key_hash);




