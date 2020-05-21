CREATE TABLE namespace
(
  name character varying(63) NOT NULL,
  PRIMARY KEY (name)
);

CREATE TABLE environment
(
  name character varying(63) NOT NULL,
  doc jsonb NOT NULL,
  PRIMARY KEY(name)
);

CREATE TABLE app
(
  id uuid NOT NULL,
  name character varying(63) NOT NULL,
  namespace character varying(63) NOT NULL REFERENCES namespace(name),
  PRIMARY KEY(id)
);

CREATE UNIQUE INDEX ix_app_name ON app(name, namespace);

CREATE TABLE deployment_reservation
(
  id uuid NOT NULL,
  app_id uuid NOT NULL REFERENCES app(id),
  name character varying(63) NOT NULL,
  namespace character varying(63) NOT NULL REFERENCES namespace(name),
  PRIMARY KEY(id)
);

CREATE UNIQUE INDEX ix_deployment_reservation_namespace ON deployment_reservation(name, namespace);

CREATE TABLE deployment
(
  id uuid NOT NULL,
  deployment_reservation_id uuid NOT NULL REFERENCES deployment_reservation(id),
  environment_name character varying(63) NOT NULL REFERENCES environment(name),
  riser_revision integer NOT NULL DEFAULT(0),
  deleted_at TIMESTAMP WITH TIME ZONE,
  doc jsonb NOT NULL,
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX ix_deployment_environment_name ON deployment(deployment_reservation_id, environment_name);
CREATE INDEX ix_deployment_riser_revision ON deployment(riser_revision);
CREATE INDEX ix_deployment_deleted_at ON deployment(deleted_at);

CREATE TABLE secret_meta (
  name character varying(63) NOT NULL,
  app_id uuid NOT NULL REFERENCES app(id),
  environment_name character varying(63) NOT NULL REFERENCES environment(name),
  revision integer NOT NULL DEFAULT(0),
  committed_revision integer NOT NULL DEFAULT(0),
  PRIMARY KEY (name, app_id, environment_name)
);

CREATE INDEX ix_secretmeta_revision ON secret_meta(revision);
CREATE INDEX ix_secretmeta_committed_revision ON secret_meta(committed_revision);

/* "user" is a reserved word in Postgress. Easier to just use riser_user. The domain will still call this resource a "user" */
CREATE TABLE riser_user
(
  id uuid NOT NULL,
  username character varying(32) NOT NULL,
  doc jsonb NOT NULL,
  PRIMARY KEY(id)
);

CREATE UNIQUE INDEX ix_riser_user_username ON riser_user(username);

CREATE TABLE apikey
(
  riser_user_id uuid NOT NULL REFERENCES riser_user(id),
  key_hash bytea NOT NULL
);

CREATE INDEX ix_apikey_riser_user_id ON apikey(riser_user_id);
-- The hash must be unique across users since we find the user based on the hash alone
CREATE UNIQUE INDEX ix_apikey_key_hash ON apikey(key_hash);




