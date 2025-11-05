create table if not exists test(id int);

create type status_type as enum ('active', 'inactive');

create table if not exists user_data (
    user_id uuid primary key,
    email text not null,
    password_hash text not null,
    crea timestamp not null,
    lupa timestamp not null,
    status status_type not null
);

create table if not exists user_analytics (
    analytic_id uuid primary key,
    keyname int not null,
    value text not null,
    user_id uuid not null,
    foreign key (user_id) references user_data (user_id)
);

create table if not exists generated_url (
    url_id uuid primary key,
    url text not null,
    short text not null,
    crea timestamp not null,
    lupa timestamp not null,
    status status_type not null,
    user_id uuid not null,
    foreign key (user_id) references user_data (user_id)
);