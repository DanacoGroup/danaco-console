-- Przestrzeń badania stała jednym wierszem na instalację: klucz główny niósł
-- CHECK(id = 1) (migracja 049). Konto młodsze nie dostawało własnej przestrzeni,
-- tylko odmowę — wiersz o kluczu 1 należał już do konta starszego.
--
-- Pamięć tłumaczeń nie miała czym wskazać właściciela: kolumny `konto_id` nie
-- dostała krokiem 484, a jej `panel_id` bywa pusty, więc drogi przez panel nie
-- ma. Segmenty dwóch Operatorów leżały w jednym zbiorze i podpowiadały się
-- nawzajem.
--
-- Przestrzeń dostaje klucz sztuczny i jeden wiersz na konto; pamięć tłumaczeń
-- dostaje wskazanie konta wzorem kroku 484 — kolumna dopuszcza NULL, bo wiersze
-- zastane powstały przed rozdzieleniem kont.

CREATE TABLE IF NOT EXISTS kontrola_przejazdu (
    wersja  INTEGER NOT NULL,
    tabela  TEXT    NOT NULL,
    roznica INTEGER NOT NULL CHECK (roznica = 0)
);

ALTER TABLE pamiec_tlumaczen ADD COLUMN konto_id INTEGER;

PRAGMA foreign_keys = off;

CREATE TABLE przestrzen_badania_nowa (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    zakres         TEXT    NOT NULL DEFAULT '',
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    odbiorca       TEXT,
    protokol       TEXT,
    granice        TEXT,
    notatka        TEXT,
    konto_id       INTEGER
);

INSERT INTO przestrzen_badania_nowa
    (id, zakres, zaktualizowano, odbiorca, protokol, granice, notatka, konto_id)
SELECT id, zakres, zaktualizowano, odbiorca, protokol, granice, notatka, konto_id
  FROM przestrzen_badania;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 495, 'przestrzen_badania',
       (SELECT COUNT(*) FROM przestrzen_badania)
     - (SELECT COUNT(*) FROM przestrzen_badania_nowa);

DROP TABLE przestrzen_badania;
ALTER TABLE przestrzen_badania_nowa RENAME TO przestrzen_badania;

CREATE UNIQUE INDEX idx_przestrzen_badania_konto
    ON przestrzen_badania (COALESCE(konto_id, 0));

PRAGMA foreign_keys = on;
