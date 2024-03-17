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
select short_url, original_url from shortener
`
	saveCashQuery = `
insert into shortener (short_url, original_url)
select
    url_info.short_url,
    url_info.original_url
    from jsonb_to_recordset($1::jsonb) as url_info (
                                                  "short_url" text,
                                                  "original_url" text
        ) on conflict do nothing;
`
)
