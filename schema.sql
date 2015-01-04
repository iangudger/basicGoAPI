CREATE TABLE users (
   id               SERIAL                   PRIMARY KEY,
   email            TEXT           NOT NULL  UNIQUE,
   password         TEXT           NOT NULL,
   date             TIMESTAMP      NOT NULL  DEFAULT now()
);

CREATE TABLE sessions (
   id               TEXT                     PRIMARY KEY,
   userid           INT            NOT NULL  REFERENCES users(id),
   start            TIMESTAMP      NOT NULL  DEFAULT now(),
   last             TIMESTAMP      NOT NULL  DEFAULT now()
);
