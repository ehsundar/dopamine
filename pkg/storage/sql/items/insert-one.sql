insert into ?(contents, created_at)
values (?, ?)
returning id, contents, created_at;
