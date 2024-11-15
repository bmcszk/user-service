CREATE TABLE users (
  id         BIGSERIAL PRIMARY KEY,
  name       text      NOT NULL,
  other      text,
  created_at timestamp NOT NULL,
  updated_at timestamp
);
