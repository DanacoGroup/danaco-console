-- Migracja 194 zakłada tabele wariantów, głosów i wyniku agregacji głosowania nad stanowiskami, z jednym głosem na wyborcę zastępowanym przy zmianie zdania.

CREATE TABLE debata_glosowanie (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tura                     TEXT    NOT NULL DEFAULT '',
    metoda                   TEXT    NOT NULL
                                     CHECK(metoda IN ('approval','irv','schulze','score','quadratic')),
    stan                     TEXT    NOT NULL DEFAULT 'open'
                                     CHECK(stan IN ('open','closed','tied')),
    -- Próg zgody ujemny znaczy „bez progu".
    prog                     REAL    NOT NULL DEFAULT -1.0,
    -- Uprawnieni do głosu, rozdzieleni znakiem nowego wiersza; pusto znaczy
    -- „cały skład".
    uprawnieni               TEXT    NOT NULL DEFAULT '',
    rozpoczeto               TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zamknieto                TEXT
);
CREATE INDEX idx_debata_glosowanie_okno ON debata_glosowanie(okno, id DESC);

CREATE TABLE debata_wariant (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    glosowanie               TEXT    NOT NULL,
    etykieta                 TEXT    NOT NULL,
    wypowiedz                TEXT    NOT NULL DEFAULT '',
    kolejnosc                INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_debata_wariant_glosowanie ON debata_wariant(glosowanie, kolejnosc, id);

CREATE TABLE debata_glos (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    glosowanie               TEXT    NOT NULL,
    wyborca                  TEXT    NOT NULL,
    aprobaty                 TEXT    NOT NULL DEFAULT '',
    ranking                  TEXT    NOT NULL DEFAULT '',
    punkty_json              TEXT    NOT NULL DEFAULT '',
    oddano                   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (glosowanie, wyborca)
);
CREATE INDEX idx_debata_glos_glosowanie ON debata_glos(glosowanie, id);
