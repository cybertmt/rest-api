DROP TABLE IF EXISTS locations;

CREATE TABLE locations
(
    id      SERIAL PRIMARY KEY,
    title   TEXT   NOT NULL,
    content TEXT   NOT NULL,
    link    TEXT   NOT NULL,
    latitude FLOAT NOT NULL,
    longitude  FLOAT NOT NULL

);

-- INSERT INTO posts (title, content, pubtime, link) VALUES ('Статья 1', 'Содержание статьи 1', 1,'http//http1');
-- INSERT INTO posts (title, content, pubtime, link) VALUES ('Статья 2', 'Содержание статьи 2', 0,'http//http2');
-- INSERT INTO posts (title, content, pubtime, link) VALUES ('Статья 3', 'Содержание статьи 3', 2,'http//http3');
-- INSERT INTO posts (author_id, title, content, created_at) VALUES (1, 'Статья 2', 'Содержание статьи 2', 0);
-- INSERT INTO posts (author_id, title, content, created_at) VALUES (1, 'Статья 3', 'Содержание статьи 3', 0);
-- INSERT INTO posts (author_id, title, content, created_at) VALUES (0, 'Статья 4', 'Содержание статьи 4', 11);

INSERT INTO locations (title, content, link, latitude, longitude) VALUES ('Msc Apt', 'Moscow', 'https://ya.ru', 55.751244,37.618423);
INSERT INTO locations (title, content, link, latitude, longitude) VALUES ('NY Apt', 'New York', 'https://google.com', 40.650002,-73.949997);
INSERT INTO locations (title, content, link, latitude, longitude) VALUES ('Syd Apt', 'Sydney', 'https://yahoo.com', 151.268865,37.618423);
