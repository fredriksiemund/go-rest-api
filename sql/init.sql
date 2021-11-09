CREATE TABLE albums (
    id SERIAL,
    title varchar(50) NOT NULL,
    artist varchar(50) NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO albums (
    title,
    artist
)
VALUES
    ('American Idiot', 'Green Day'),
    ('Lemon Tree ', 'Fools Garden'),
    ('Come Away With Me', 'Norah Jones');
