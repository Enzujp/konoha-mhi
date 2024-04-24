-- create user table
CREATE TABLE IF NOT EXISTS users
(
    id             uuid    NOT NULL PRIMARY KEY,
    first_name     varchar NOT NULL DEFAULT '',
    last_name      varchar NOT NULL DEFAULT '',
    email          varchar NOT NULL DEFAULT '',
    wallet_balance numeric NOT NULL DEFAULT 0,
    CONSTRAINT uq_email UNIQUE (email)
);

-- seed user table
INSERT INTO users (id, first_name, last_name, email, wallet_balance)
VALUES (gen_random_uuid(), 'Yuji', 'Itadori', 'yuji.itadori@example.com', FLOOR(RANDOM() * 15001 + 5000)),
       (gen_random_uuid(), 'Sukuna', 'Shinazugawa', 'sukuna.shinazugawa@example.com', FLOOR(RANDOM() * 15001 + 5000)),
       (gen_random_uuid(), 'Satoru', 'Gojo', 'satoru.gojo@example.com', FLOOR(RANDOM() * 15001 + 5000)),
       (gen_random_uuid(), 'Kenji', 'Kamado', 'kenji.kamado@example.com', FLOOR(RANDOM() * 15001 + 5000)),
       (gen_random_uuid(), 'Tanjiro', 'Kamado', 'tanjiro.kamado@example.com', FLOOR(RANDOM() * 15001 + 5000)),
       (gen_random_uuid(), 'Zenitsu', 'Agatsuma', 'zenitsu.agatsuma@example.com', FLOOR(RANDOM() * 15001 + 5000)),
       (gen_random_uuid(), 'Inosuke', 'Hashibira', 'inosuke.hashibira@example.com', FLOOR(RANDOM() * 15001 + 5000)),
       (gen_random_uuid(), 'Mikasa', 'Ackerman', 'mikasa.ackerman@example.com', FLOOR(RANDOM() * 15001 + 5000)),
       (gen_random_uuid(), 'Levi', 'Ackerman', 'levi.ackerman@example.com', FLOOR(RANDOM() * 15001 + 5000)),
       (gen_random_uuid(), 'Senjougahara', 'Kyoko', 'senjougahara.kyoko@example.com', FLOOR(RANDOM() * 15001 + 5000)),
       (gen_random_uuid(), 'Kurapika', 'Zoldyck', 'kurapika.zoldyck@example.com', FLOOR(RANDOM() * 15001 + 5000)),
       (gen_random_uuid(), 'Kurapika', 'Moriya', 'kurapika.moriya@example.com', FLOOR(RANDOM() * 15001 + 5000)),
       (gen_random_uuid(), 'Yuno', 'Gonzalez', 'yuno.gonzalez@example.com', FLOOR(RANDOM() * 15001 + 5000)),
       (gen_random_uuid(), 'Yami', 'Kurapika', 'yami.kurapika@example.com', FLOOR(RANDOM() * 15001 + 5000));
