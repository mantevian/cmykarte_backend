create table users (
	id uuid,
	name varchar(64),
	password varchar(256),
	created_at timestamptz,

	primary key (id),
	unique (name)
);

create table matches (
	id text primary key,
	player_rgb_id uuid not null references users(id),
	player_cmy_id uuid not null references users(id),
	first_player_id uuid not null references users(id),
	status text not null check (status in ('in_progress', 'completed', 'abandoned')),
	winner_id uuid references users(id),
	winner_side text,
	started_at timestamptz not null,
	ended_at timestamptz,
	duration_seconds double precision,
	rounds_played integer not null default 1,
	final_score_rgb integer not null default 0,
	final_score_cmy integer not null default 0,
	end_reason text,
	build_version text not null,
	ruleset_version text not null,
	rules jsonb not null
);

create table match_events (
	id bigint generated always as identity primary key,
	match_id text not null references matches(id),
	event_type text not null check (event_type in ('card_played', 'turn_ended', 'match_abandoned')),
	round integer not null,
	turn_number integer,
	player_id uuid not null references users(id),
	player_side text not null,
	card_cost integer,
	card_type text,
	card_color jsonb,
	card_pattern jsonb,
	target_x integer,
	target_y integer,
	score_rgb integer not null,
	score_cmy integer not null,
	outcome jsonb not null default '{}'::jsonb,
	created_at timestamptz not null default now()
);

create index match_events_match_order_idx on match_events (match_id, id);
create index match_events_card_idx on match_events (event_type, card_type, card_cost);
create index matches_started_at_idx on matches (started_at);

create view analytics_match_health as
select
	status,
	end_reason,
	started_at,
	ended_at,
	duration_seconds,
	rounds_played,
	final_score_rgb,
	final_score_cmy,
	winner_side,
	first_player_id,
	player_rgb_id,
	player_cmy_id,
	build_version,
	ruleset_version
from matches;

create view analytics_card_usage as
select
	e.id as event_id,
	e.match_id,
	e.player_id,
	e.player_side,
	e.round,
	e.turn_number,
	e.card_cost,
	e.card_type,
	e.card_color,
	e.card_pattern,
	e.target_x,
	e.target_y,
	e.score_rgb,
	e.score_cmy,
	e.outcome,
	m.status as match_status,
	m.winner_side,
	m.first_player_id,
	m.build_version,
	m.ruleset_version,
	e.created_at
from match_events e
join matches m on m.id = e.match_id
where e.event_type = 'card_played';

create view analytics_turns as
select
	e.match_id,
	e.player_id,
	e.player_side,
	e.round,
	e.turn_number,
	e.score_rgb,
	e.score_cmy,
	e.created_at,
	m.status as match_status,
	m.winner_side,
	m.first_player_id
from match_events e
join matches m on m.id = e.match_id
where e.event_type = 'turn_ended';

\getenv metabase_db_password METABASE_DB_PASSWORD
\getenv postgres_db POSTGRES_DB
create role metabase_readonly login password :'metabase_db_password';
grant connect on database :"postgres_db" to metabase_readonly;
grant usage on schema public to metabase_readonly;
grant select on matches, match_events, analytics_match_health, analytics_card_usage, analytics_turns
	to metabase_readonly;
