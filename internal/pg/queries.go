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
    short_url    text not null primary key,
    original_url text not null
);
`
	getAllURLsQuery = `
select short_url, original_url from shortener.shortener
`
)
