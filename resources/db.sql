create table if not exists user
(
    user_id            bigint      not null
    primary key,
    area_code          int         null,
    current_login_time bigint      null,
    game_mode          varchar(30) not null
    );

create index nc_area_code
    on user (area_code);

create index nc_game_mode
    on user (game_mode);

create index nc_user_id
    on user (user_id);

