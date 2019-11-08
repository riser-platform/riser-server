CREATE TABLE app
(
  name character varying(63) NOT NULL,
  hashid bytea NOT NULL,
	PRIMARY KEY(name)
);

CREATE UNIQUE INDEX ix_app_hashid ON app(hashid);

CREATE TABLE deployment
(
  name character varying(63) NOT NULL,
  stage_name character varying(63) NOT NULL,
  app_name character varying(63) NOT NULL REFERENCES app(name),
  riser_generation integer NOT NULL DEFAULT(0),
  doc jsonb NOT NULL,
  PRIMARY KEY (name,stage_name)
);

CREATE UNIQUE INDEX ix_deployment ON deployment(name,stage_name);
CREATE INDEX ix_deployment_riser_generation ON deployment(riser_generation);

CREATE TABLE secret_meta (
  app_name character varying(63) NOT NULL REFERENCES app(name),
  stage_name character varying(63) NOT NULL,
  secret_name character varying(63) NOT NULL,
  doc jsonb NOT NULL,
  PRIMARY KEY (app_name, stage_name, secret_name)
);

CREATE UNIQUE INDEX ix_secret_meta ON secret_meta(app_name, stage_name, secret_name);

/* "user" is a reserved word in Postgress. Easier to just use riser_user. The domain will still call this resource a "user" */
CREATE TABLE riser_user
(
  id serial,
  username character varying(32) NOT NULL,
  doc jsonb NOT NULL,
  PRIMARY KEY(id)
);

CREATE UNIQUE INDEX ix_riser_user_username ON riser_user(username);

CREATE TABLE apikey
(
  id serial,
  riser_user_id integer NOT NULL REFERENCES riser_user(id),
  key_hash bytea NOT NULL,
  PRIMARY KEY(id)
);

CREATE INDEX ix_userlogin_user_id ON apikey(riser_user_id);
CREATE UNIQUE INDEX ix_userlogin_key_hash ON apikey(key_hash);

CREATE TABLE stage (
  name character varying(63) NOT NULL,
  doc jsonb NOT NULL,
  PRIMARY KEY(name)
)


