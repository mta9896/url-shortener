drop table if exists urls;
create table urls (
    id int auto_increment not null,
    short_url varchar(128) not null,
    main_url varchar(256) not null,
    primary key (`id`)
);

insert into urls (id, short_url, main_url)
values (1, 'short_url', 'main_url');