CREATE TABLE IF NOT EXISTS users
(
    user_id    serial primary key,
    full_name  varchar            default '',
    phone      varchar            default '',
    email      varchar            default '' UNIQUE,
    password   varchar            default '',
    verified   bool               default false,
    hash       varchar            default '',
    type       varchar            default 'user',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TABLE IF NOT EXISTS film_company (
--     company_id serial primary key ,
--     name varchar default '',
--     city varchar default '',
--     contacts varchar default '',
--     created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
-- );

CREATE TABLE IF NOT EXISTS films
(
    film_id    serial primary key,
    name       varchar            default '',
    price      float              default 0.0,
    rating     int                default 0,
    user_id    integer REFERENCES users (user_id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS subscriptions
(
    subscription_id serial primary key,
    user_id         int REFERENCES users (user_id),
    film_id         int REFERENCES films (film_id),
    expires         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);