-- Migracja 173 zakłada tabele kanałów RSS, Atom i JSON Feed oraz ich wpisów, z kaskadowym usuwaniem wpisów i warunkiem unikalności adresu w oknie.

CREATE TABLE kanal_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    tytul                    TEXT,
    postac                   TEXT,
    interwal_sekund          INTEGER,
    pobrano                  TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(okno, url)
);
CREATE INDEX idx_kanal_przegladania_okno ON kanal_przegladania(okno, utworzono DESC, id);

CREATE TABLE wpis_kanalu_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kanal_zewnetrzny_id      TEXT    NOT NULL REFERENCES kanal_przegladania(identyfikator_zewnetrzny) ON DELETE CASCADE,
    url                      TEXT    NOT NULL,
    tytul                    TEXT,
    streszczenie             TEXT,
    przeczytany              INTEGER NOT NULL DEFAULT 0 CHECK(przeczytany IN (0,1)),
    opublikowano             TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(kanal_zewnetrzny_id, url)
);
CREATE INDEX idx_wpis_kanalu_przegladania ON wpis_kanalu_przegladania(kanal_zewnetrzny_id, opublikowano DESC, id);
