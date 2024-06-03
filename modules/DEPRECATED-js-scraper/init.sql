CREATE TABLE trending_github (
  name        text,
  href        text,
  description text,
  stars       integer,
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE trending_movies (
  title       text,
  link        text,
  img         text,
  year        integer,
  rank        integer,
  velocity    integer,
  rating      varchar(3),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE trending_npm (
  subject     text,
  page        integer,
  rank        integer,
  title       text,
  description text,
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE trending_pypi (
  downloads   integer,
  name        text,
  description text,
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE trending_tv (
  link        text,
  img         text,
  title       text,
  rank        integer,
  velocity    integer,
  rating      varchar(3),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE trending_games (
  link             text,
  img              text,
  title            text,
  regular_price    DECIMAL(19,4),
  discounted_price DECIMAL(19,4),
  is_free          BOOLEAN,
  rating           text,
  release          text,
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE upcoming_movies (
  title             text,
  release           date,
  genre             text,
  actors            text,
  img               text,
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE trending_go (
  name        text,
  href        text,
  description text,
  stars       integer,
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- for mongo like object data I can use raw_json JSONB,

-- INSERT INTO trending_github(name, href, description, stars)
-- VALUES ('freeCodeCamp', 'https://github.com/freeCodeCamp/freeCodeCamp', 'freeCodeCamp.org''s open-source codebase and curriculum. Learn to code ...', 334801);

-- INSERT INTO trending_github(name, href, description, stars) VALUES($1, $2, $3, $4);

-- INSERT INTO upcoming_movies(raw_json) VALUES('{"thing": "wow"}');

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp_trending_github
BEFORE UPDATE ON trending_github
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_trending_go
BEFORE UPDATE ON trending_go
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_trending_movies
BEFORE UPDATE ON trending_movies
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_trending_npm
BEFORE UPDATE ON trending_npm
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_trending_pypi
BEFORE UPDATE ON trending_pypi
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_trending_tv
BEFORE UPDATE ON trending_tv
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_trending_games
BEFORE UPDATE ON trending_games
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_upcoming_movies
BEFORE UPDATE ON upcoming_movies
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();