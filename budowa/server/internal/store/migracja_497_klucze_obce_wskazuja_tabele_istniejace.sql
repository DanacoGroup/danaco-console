-- Cztery tabele wskazują kluczem obcym tabelę przejściową przebudowy, której już
-- nie ma: `pozycja_kolejki` i `log_akcji_kolejki` sięgają `kolejka_nowa`,
-- `debata_wypowiedz` — `debata_tura_nowa`, a `etykieta_zasobu_design` —
-- `zasob_design_nowy`. Nazwa przejściowa została w klauzuli REFERENCES po
-- przebudowie prowadzonej przy wyłączonych kluczach obcych.
--
-- Skutek nie jest teoretyczny: zapis wypowiedzi debaty kończy się odmową
-- „no such table: main.debata_tura_nowa" i tura nie zbiera ani jednej
-- wypowiedzi. Kontrola więzów tego nie wychwyciła, bo `foreign_key_check`
-- sprawdza dane wierszy, a nie istnienie tabeli wskazanej.
--
-- Naprawa przepisuje te tabele z klauzulą wskazującą tabelę istniejącą.
-- Kontrola przejazdu pilnuje liczebności każdej z nich.

CREATE TABLE IF NOT EXISTS kontrola_przejazdu (
    wersja  INTEGER NOT NULL,
    tabela  TEXT    NOT NULL,
    roznica INTEGER NOT NULL CHECK (roznica = 0)
);

PRAGMA foreign_keys = off;

-- ── pozycja kolejki ─────────────────────────────────────────────────────────

CREATE TABLE pozycja_kolejki_prosta (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    kolejka_id          INTEGER NOT NULL REFERENCES kolejka(id) ON DELETE CASCADE,
    okno_wykonawcy_id   INTEGER REFERENCES okno_komunikacji(id) ON DELETE SET NULL,
    kolejnosc           INTEGER NOT NULL DEFAULT 0,
    tytul               TEXT    NOT NULL,
    tresc_zlecenia      TEXT,
    tresc_odwolanie     TEXT,
    stan                TEXT    NOT NULL DEFAULT 'oczekuje'
                                CHECK(stan IN ('oczekuje','przydzielona','wykonywana',
                                               'do_weryfikacji','ukonczona','bledna','anulowana')),
    werdykt_weryfikacji TEXT    CHECK(werdykt_weryfikacji IS NULL OR
                                      werdykt_weryfikacji IN ('przyjete','do_poprawy','odrzucone')),
    licznik_obiegow     INTEGER NOT NULL DEFAULT 0,
    utworzono           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO pozycja_kolejki_prosta
    (id, kolejka_id, okno_wykonawcy_id, kolejnosc, tytul, tresc_zlecenia, tresc_odwolanie,
     stan, werdykt_weryfikacji, licznik_obiegow, utworzono, zaktualizowano)
SELECT id, kolejka_id, okno_wykonawcy_id, kolejnosc, tytul, tresc_zlecenia, tresc_odwolanie,
       stan, werdykt_weryfikacji, licznik_obiegow, utworzono, zaktualizowano
  FROM pozycja_kolejki;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 497, 'pozycja_kolejki',
       (SELECT COUNT(*) FROM pozycja_kolejki) - (SELECT COUNT(*) FROM pozycja_kolejki_prosta);

DROP TABLE pozycja_kolejki;
ALTER TABLE pozycja_kolejki_prosta RENAME TO pozycja_kolejki;

-- ── dziennik akcji kolejki ──────────────────────────────────────────────────

CREATE TABLE log_akcji_kolejki_prosta (
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    kolejka_id         INTEGER NOT NULL REFERENCES kolejka(id) ON DELETE CASCADE,
    pozycja_kolejki_id INTEGER REFERENCES pozycja_kolejki(id) ON DELETE CASCADE,
    akcja              TEXT    NOT NULL,
    stan_przed         TEXT,
    stan_po            TEXT,
    numer_obiegu       INTEGER NOT NULL DEFAULT 0,
    szczegoly          TEXT,
    utworzono          TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO log_akcji_kolejki_prosta
    (id, kolejka_id, pozycja_kolejki_id, akcja, stan_przed, stan_po, numer_obiegu,
     szczegoly, utworzono)
SELECT id, kolejka_id, pozycja_kolejki_id, akcja, stan_przed, stan_po, numer_obiegu,
       szczegoly, utworzono
  FROM log_akcji_kolejki;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 497, 'log_akcji_kolejki',
       (SELECT COUNT(*) FROM log_akcji_kolejki) - (SELECT COUNT(*) FROM log_akcji_kolejki_prosta);

DROP TABLE log_akcji_kolejki;
ALTER TABLE log_akcji_kolejki_prosta RENAME TO log_akcji_kolejki;

-- ── wypowiedź debaty ────────────────────────────────────────────────────────

CREATE TABLE debata_wypowiedz_prosta (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    tura_id                  INTEGER NOT NULL REFERENCES debata_tura(id) ON DELETE CASCADE,
    uczestnik                TEXT    NOT NULL,
    tresc                    TEXT    NOT NULL DEFAULT '',
    odpowiedz_na             TEXT    NOT NULL DEFAULT '',
    akt_mowy                 TEXT    NOT NULL DEFAULT ''
                                     CHECK(akt_mowy IN ('','claim','argument','counterArgument',
                                                        'question','concession','rebuttal')),
    pewnosc                  REAL    NOT NULL DEFAULT -1.0,
    redakcja                 INTEGER NOT NULL DEFAULT 1,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO debata_wypowiedz_prosta
    (id, identyfikator_zewnetrzny, tura_id, uczestnik, tresc, odpowiedz_na, akt_mowy,
     pewnosc, redakcja, utworzono)
SELECT id, identyfikator_zewnetrzny, tura_id, uczestnik, tresc, odpowiedz_na, akt_mowy,
       pewnosc, redakcja, utworzono
  FROM debata_wypowiedz;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 497, 'debata_wypowiedz',
       (SELECT COUNT(*) FROM debata_wypowiedz) - (SELECT COUNT(*) FROM debata_wypowiedz_prosta);

DROP TABLE debata_wypowiedz;
ALTER TABLE debata_wypowiedz_prosta RENAME TO debata_wypowiedz;

-- ── etykieta zasobu Designu ─────────────────────────────────────────────────

CREATE TABLE etykieta_zasobu_design_prosta (
    zasob_id  INTEGER NOT NULL REFERENCES zasob_design(id) ON DELETE CASCADE,
    etykieta  TEXT    NOT NULL,
    utworzono TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (zasob_id, etykieta)
);

INSERT INTO etykieta_zasobu_design_prosta (zasob_id, etykieta, utworzono)
SELECT zasob_id, etykieta, utworzono FROM etykieta_zasobu_design;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 497, 'etykieta_zasobu_design',
       (SELECT COUNT(*) FROM etykieta_zasobu_design)
     - (SELECT COUNT(*) FROM etykieta_zasobu_design_prosta);

DROP TABLE etykieta_zasobu_design;
ALTER TABLE etykieta_zasobu_design_prosta RENAME TO etykieta_zasobu_design;

CREATE INDEX idx_etykieta_zasobu_design_etykieta
    ON etykieta_zasobu_design(etykieta, zasob_id);

PRAGMA foreign_keys = on;
