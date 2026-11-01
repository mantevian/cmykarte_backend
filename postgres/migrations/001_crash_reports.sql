create table if not exists crash_reports (
	report_id uuid primary key,
	received_at timestamptz not null default now(),
	report_type text not null check (report_type in ('Error', 'Assert', 'Exception', 'Crash')),
	message text not null,
	stack_trace text not null,
	app_version text not null,
	unity_version text not null,
	platform text not null,
	scene text not null
);

create index if not exists crash_reports_received_at_idx on crash_reports (received_at);
grant select on crash_reports to metabase_readonly;
