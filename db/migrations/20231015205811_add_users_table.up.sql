create table users(
    Id bigserial primary key,
    Username varchar not null,
    Email varchar not null ,
    Password varchar not null
)