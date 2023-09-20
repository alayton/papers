create table users
(
    id             bigint auto_increment
        primary key,
    email          varchar(100)      not null,
    username       varchar(60)       null,
    password       varchar(80)       null,
    totp_secret    varchar(80)       null,
    confirmed      tinyint           not null,
    confirm_token  varchar(100)      null,
    recovery_token varchar(100)      null,
    locked_until   datetime          null,
    attempts       tinyint default 0 not null,
    last_attempt   datetime          null,
    created_at     datetime          not null,
    last_login     datetime          not null,
    constraint users_confirm_token_uindex
        unique (confirm_token),
    constraint users_email_uindex
        unique (email),
    constraint users_recovery_token_uindex
        unique (recovery_token)
);

create index users_created_at_index
    on users (created_at);

create index users_last_login_index
    on users (last_login);

create index users_username_index
    on users (username);

create table access_tokens
(
    user_id    bigint      not null,
    token      varchar(40) not null,
    valid      tinyint     not null,
    chain      varchar(40) not null,
    created_at datetime    not null,
    primary key (user_id, token),
    constraint access_tokens_users_id_fk
        foreign key (user_id) references users (id)
            on delete cascade
);

create index access_tokens_created_at_index
    on access_tokens (created_at);

create index access_tokens_user_id_chain_index
    on access_tokens (user_id, chain);

create table refresh_tokens
(
    user_id    bigint      not null,
    token      varchar(40) not null,
    valid      tinyint     not null,
    chain      varchar(40) not null,
    created_at datetime    not null,
    primary key (user_id, token),
    constraint refresh_tokens_users_id_fk
        foreign key (user_id) references users (id)
            on delete cascade
);

create index refresh_tokens_created_at_index
    on refresh_tokens (created_at);

create index refresh_tokens_user_id_chain_index
    on refresh_tokens (user_id, chain);

create table user_identities
(
    user_id  bigint       not null,
    provider varchar(40)  not null,
    identity varchar(100) not null,
    primary key (user_id, provider, identity),
    constraint user_identities_users_id_fk
        foreign key (user_id) references users (id)
            on delete cascade
);

create table user_permissions
(
    user_id    bigint      not null,
    permission varchar(40) not null,
    primary key (user_id, permission),
    constraint user_permissions_users_id_fk
        foreign key (user_id) references users (id)
            on delete cascade
);

create table user_roles
(
    user_id bigint      not null,
    role    varchar(40) not null,
    primary key (user_id, role),
    constraint user_roles_users_id_fk
        foreign key (user_id) references users (id)
            on delete cascade
);

