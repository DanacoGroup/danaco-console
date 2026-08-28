-- Migracja 252 rozszerza warunek kolumny rodzaju powłoki karty terminala
-- o cztery kolejne rodzaje powłoki, odtwarzając tabelę karty od nowa.

CREATE TABLE terminal_karta_nowa (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    kod             TEXT    NOT NULL UNIQUE,
    okno_kod        TEXT    NOT NULL,
    powloka         TEXT    NOT NULL
                            CHECK(powloka IN ('powershell','cmd','bash','node','python','ssh',
                                              'container','pod','serial','telnet')),
    tytul           TEXT,
    katalog_roboczy TEXT,
    stan            TEXT    NOT NULL DEFAULT 'running' CHECK(stan IN ('running','exited')),
    utworzono       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    cel_zdalny      TEXT    NOT NULL DEFAULT '',
    port_zdalny     INTEGER,
    host_kod        TEXT
);

INSERT INTO terminal_karta_nowa
    (id, kod, okno_kod, powloka, tytul, katalog_roboczy, stan, utworzono,
     cel_zdalny, port_zdalny, host_kod)
SELECT id, kod, okno_kod, powloka, tytul, katalog_roboczy, stan, utworzono,
       cel_zdalny, port_zdalny, host_kod
FROM terminal_karta;

DROP TABLE terminal_karta;
ALTER TABLE terminal_karta_nowa RENAME TO terminal_karta;

CREATE INDEX idx_terminal_karta_okno ON terminal_karta(okno_kod, utworzono DESC);
