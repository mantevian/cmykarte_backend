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

Unity error, assertion, and exception logs are queued locally and submitted to
`POST /crash-reports`; failed uploads remain queued for the next launch. A
marker also reports an unclean previous shutdown, which can indicate a native
crash or forced termination but cannot provide a native crash dump.
Crash reports include the message, stack trace, app/Unity versions, platform,
and active scene, but do not include player IDs or credentials. The provisioner
creates a Crash Reports dashboard with report trends, build/platform breakdowns,
frequent messages, and recent stack traces.
Exception messages and stack traces can still contain runtime data; do not put
secrets or sensitive user input in application error messages.

For an existing PostgreSQL volume, apply
`postgres/migrations/001_crash_reports.sql` before deploying the updated game
server, then rerun the `metabase-provisioner` service to create the dashboard.
The Unity server defaults to `http://localhost:8080` for local development.
Pass `-server-url=https://your-game-server.example.com` in the player launch
arguments to use a deployed backend; the WebSocket connection automatically
uses `wss` for HTTPS. Put a TLS-terminating reverse proxy in front of the Go
server before accepting public Internet traffic.
