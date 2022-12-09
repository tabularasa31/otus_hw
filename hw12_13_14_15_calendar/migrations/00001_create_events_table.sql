-- +goose Up
CREATE TABLE events
(
    id           serial primary key,
    title        text not null ,
    descr        text,
    user_id      bigint not null ,
    event_time   timestamp(0) with time zone not null,
    duration     time,
    notification timestamp(0) with time zone
);

CREATE INDEX user_id_idx
    ON events (user_id);

CREATE INDEX event_time_idx
    ON events (event_time);

INSERT INTO events VALUES
    (42, 'new year', 'Happy New Year!', '2019-12-31 23:59:59', '01:00:00', '2019-12-31 23:59:50'),
    (43, 'new year 2', 'Happy New Year 22!', '2019-12-31 23:59:59', '01:00:00', '2019-12-31 23:59:50'),
    (47, 'happy birthday', 'Its a new date', '2020-05-31 09:00:00', '01:30:00', '2020-05-31 08:00:00'),
    (48, 'happy birthday 2', 'Its a new date 22', '2020-06-05 09:00:00', '01:30:00', '2020-06-05 08:00:00');


-- +goose Down
DROP TABLE events;