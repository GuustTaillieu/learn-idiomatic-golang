-- insert a couple of fruits into the items table
INSERT OR IGNORE INTO items (id, name)
VALUES ('00000000-0000-0000-0000-000000000001', 'Gopher Plushie'),
       ('00000000-0000-0000-0000-000000000002', 'Gopher Sticker'),
       ('00000000-0000-0000-0000-000000000003', 'Gopher Mug');

INSERT OR IGNORE INTO inventory (id, item_id, quantity, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    1000,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);
