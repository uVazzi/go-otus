-- +goose Up
-- +goose StatementBegin
INSERT INTO "user" (uuid, username) VALUES ('8f3d77cc-1234-4567-890a-abcdefabcdef', 'otus');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "user" WHERE uuid = '8f3d77cc-1234-4567-890a-abcdefabcdef';
-- +goose StatementEnd
