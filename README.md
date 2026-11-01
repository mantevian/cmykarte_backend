# CMYKarte backend

## Run the backend and analytics

Set the existing database and pgAdmin values, plus `METABASE_DB_PASSWORD`,
`METABASE_ADMIN_EMAIL`, and `METABASE_ADMIN_PASSWORD` in `.env`. The
`.env.example` file lists the required variables; replace every example
password before starting the services. Use a valid email address and a unique,
uncommon Metabase admin password of at least 12 characters.

Start the services with `docker compose -f compose.yaml -f compose.dev.yaml up --build`
for development with Go hot reload. To run the production-style container
locally, use `docker compose up --build`. Metabase is available at
http://localhost:3000 and the game server at http://localhost:8080. Metabase
binds to localhost and requires the configured admin login. On first startup,
the provisioner creates a read-only PostgreSQL connection and the Match Health,
Card Balance, and Gameplay Progression dashboards.

## Deploying to an Alpine server

The GitHub Actions workflow runs the Go tests on pushes to `main`, then deploys
over SSH only when those tests pass. On Alpine, install the required packages
and enable Docker and SSH at boot:

```sh
apk add docker docker-cli-compose git openssh
rc-update add docker default
rc-update add sshd default
service docker start
service sshd start
```

Clone this repository into a deployment directory, check out `main`, and create
that directory's `.env` from `.env.example` with unique, strong passwords. Keep
`.env` on the server; it is ignored by Git. The SSH deployment user needs the
matching public key in `~/.ssh/authorized_keys`, access to the repository, and
membership in the Docker group.

Add these repository Actions secrets under **Settings → Secrets and variables
→ Actions**:

- `DEPLOY_HOST`: the server's reachable hostname or IP.
- `DEPLOY_USER`: the SSH account used for deployment.
- `DEPLOY_SSH_KEY`: that account's private SSH key.
- `DEPLOY_KNOWN_HOSTS`: the server's SSH host-key entry, obtained and verified
  through a trusted channel.
- `DEPLOY_PATH`: the absolute path to the server's repository checkout (use a
  path without spaces or single quotes).

The server must be reachable over SSH from GitHub-hosted runners, and its
checkout must have an `origin` remote with a `main` branch. Each deployment
fast-forwards that checkout, builds the production Go image, and starts the
Compose services. Database and Metabase volumes are retained. The deployed
game's `BUILD_VERSION` is set to the short Git commit hash.

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
