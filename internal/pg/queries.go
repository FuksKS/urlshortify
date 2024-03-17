package pg

const (
	existDBQuery = `
SELECT EXISTS (SELECT FROM pg_database WHERE datname = $1);
`

	createDBQuery = `
CREATE DATABASE $1
`
	createTableQuery = `
create table if not exists $1
(
    short_url    text not null primary key,
    original_url text not null
);
`
	getAllURLsQuery = `
select short_url, original_url from shortener.shortener
`
)
