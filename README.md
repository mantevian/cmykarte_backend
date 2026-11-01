# CMYKarte backend

## Run the backend and analytics

Set the existing database and pgAdmin values, plus `METABASE_DB_PASSWORD`,
`METABASE_ADMIN_EMAIL`, and `METABASE_ADMIN_PASSWORD` in `.env`. The
`.env.example` file lists the required variables; replace every example
password before starting the services. Use a valid email address and a unique,
uncommon Metabase admin password of at least 12 characters.

Start the services with `docker compose up --build`. Metabase is available at
http://localhost:3000 and the game server at http://localhost:8080. Metabase
binds to localhost and requires the configured admin login. On first startup,
the provisioner creates a read-only PostgreSQL connection and the Match Health,
Card Balance, and Gameplay Progression dashboards.

If Metabase rejects the admin credentials during initial setup, update
`METABASE_ADMIN_EMAIL` and `METABASE_ADMIN_PASSWORD` in `.env`, then rerun the
provisioner with
`docker compose --file compose.yaml up --force-recreate metabase-provisioner`.
If setup has already completed, those variables must match the existing
Metabase admin account; changing them does not reset Metabase's stored login.

The PostgreSQL init script creates the analytics tables, reporting views, and
read-only Metabase role only when initializing an empty PostgreSQL data volume.
For an existing volume, apply the new schema and grants using your normal
database migration process before starting the backend and Metabase. Existing
matches are not backfilled.

Completed matches and disconnect-abandoned matches are recorded in `matches`.
Accepted card plays and end-turn actions are recorded in `match_events`.
`BUILD_VERSION` can be set on the game service to tag analytics with a deployed
build; the ruleset is tagged with the current version in the game code.
