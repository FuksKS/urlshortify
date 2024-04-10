package pg

const (
	existDBQuery = `
SELECT EXISTS (SELECT FROM pg_database WHERE datname = 'shortener');
`

	createDBQuery = `
CREATE DATABASE shortener
`

	createTableQuery = `
create table if not exists shortener
(
    id           text not null,
    short_url    text not null,
    original_url text not null primary key,
	user_id 	 text not null
);
`

	getAllURLsQuery = `
select id, short_url, original_url, user_id from shortener;
`

	saveOneURLQuery = `
insert into shortener (id, short_url, original_url, user_id)
values ($1, $2, $3, $4)
on conflict do nothing;
`

	getUsersURLsQuery = `
select id, short_url, original_url, user_id from shortener where user_id = $1;
`
)
