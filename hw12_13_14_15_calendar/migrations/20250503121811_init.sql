-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS event (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    delay_notification INT DEFAULT NULL,
    delay_notification_type TEXT DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_event_user_uuid ON event(user_uuid);
CREATE INDEX IF NOT EXISTS idx_event_start_date ON event(start_date);
CREATE INDEX IF NOT EXISTS idx_event_end_date ON event(end_date);

ALTER TABLE event
    ADD CONSTRAINT fk_event_user_uuid__user FOREIGN KEY (user_uuid) REFERENCES "user"(uuid) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS notification (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_uuid UUID NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_notification_event_uuid ON notification(event_uuid);
ALTER TABLE notification
    ADD CONSTRAINT fk_notification_event_uuid__event FOREIGN KEY (event_uuid) REFERENCES event(uuid) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification;
DROP TABLE IF EXISTS event;
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd
