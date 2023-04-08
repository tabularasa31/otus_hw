-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events
(
    id serial primary key,
    title text,
    descr text,
    user_id int,
    start_time timestamptz(0) not null,
    end_time timestamptz(0) not null,
    notification timestamptz(0)
);

CREATE INDEX user_id_idx
    ON events (user_id);

INSERT INTO events(user_id, title, descr, start_time, end_time, notification)
VALUES
    (42, 'new year', 'Happy New Year!', '2019-12-31 23:59:59', '2020-01-01 01:00:00', '2019-12-31 23:00:00'),
    (43, 'new year 2', 'Happy New Year 2020!', '2019-12-31 23:59:59',  '2020-01-01 01:30:00', '2019-12-31 21:50:50'),
    (45, 'new year 3', 'Happy New Year 33!', '2019-12-31 23:50:50',  '2020-01-01 02:00:00', '2019-12-31 13:00:00'),
    (47, 'new date', 'Its a new date', '2020-05-31 09:00:00', '2020-05-31 23:00:00', '2020-05-30 09:00:00'),
    (48, 'new date2', 'Its a new date2', '2020-06-05 09:00:00', '2020-06-05 10:00:00', '2020-06-05 07:00:00');
-- +goose StatementEnd


-- +goose Down
DROP TABLE IF EXISTS events;