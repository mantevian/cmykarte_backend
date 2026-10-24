create table users (
	id uuid,
	name varchar(64),
	password varchar(256),
	created_at timestamptz,

	primary key (id),
	unique (name)
);
