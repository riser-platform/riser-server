# Riser Server

<p align="left">
  <a href="https://github.com/riser-platform/riser-server"><img alt="GitHub Actions status" src="https://github.com/riser-platform/riser-server/workflows/Build/badge.svg"></a>
</p>

This is the server for the [Riser Platform](https://github.com/riser-platform/riser).

## Development

This section will be improved as we're closer to accepting PRs from the community.

### DB Migrations

Migration Tool: `go get -u -d github.com/golang-migrate/migrate/cmd/migrate`

#### Create a new migration

`migrate create -dir migrations -format unix -ext sql addstuff`
> Note that this creates both an `up` and a `down` script. You should delete the `down` script. The philosophy is that migrations should be should always "roll forward".

#### Applying migrations
By default the riser server will apply database migrations if needed during startup. You may disable this behavior by setting the environment
variable `RISER_POSTGRES_MIGRATE_ON_STARTUP=false`. To manually migrate, modify the connection string below and run:
```
migrate -source=file://migrations -database="postgres://user:password@postgreshost/riserdb" up
```

See the [migrate](github.com/golang-migrate/migrate) documentation for more details.
