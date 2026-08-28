-- Migracja 293 dodaje słownik skrótów tekstowych rozwijanych we wszystkich
-- polach tekstowych platformy.

CREATE TABLE skrot_tekstowy (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    profil_kod               TEXT    NOT NULL DEFAULT '',
    skrot                    TEXT    NOT NULL,
    tresc                    TEXT    NOT NULL,
    opis                     TEXT,
    pola_json                TEXT    NOT NULL DEFAULT '[]',
    czynny                   INTEGER NOT NULL DEFAULT 1,
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL,
    UNIQUE (profil_kod, skrot)
);
CREATE INDEX idx_skrot_tekstowy_wykaz ON skrot_tekstowy(profil_kod, skrot);
