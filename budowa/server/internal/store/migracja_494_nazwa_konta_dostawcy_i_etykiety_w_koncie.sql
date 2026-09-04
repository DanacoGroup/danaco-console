-- Domknięcie serii 490–493. Nazwa konta dostawcy modeli i nazwa etykiety
-- słownika biblioteki to napisy wpisywane przez Operatora. Jednoznaczność
-- obejmująca całą instalację odbierała drugiemu Operatorowi nazwę „robocze"
-- czy „pilne", jeśli pierwszy jej użył.
--
-- Etykieta miała nazwę kluczem głównym, więc wskazanie konta wchodzi razem
-- z kluczem sztucznym; wiersze zastane zachowują nazwy, a jednoznaczność
-- wraca wskaźnikiem po wyrażeniu, wzorem kroku 490.

CREATE TABLE IF NOT EXISTS kontrola_przejazdu (
    wersja  INTEGER NOT NULL,
    tabela  TEXT    NOT NULL,
    roznica INTEGER NOT NULL CHECK (roznica = 0)
);

PRAGMA foreign_keys = off;

-- ── konto dostawcy modeli ───────────────────────────────────────────────────

CREATE TABLE konto_nowa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    nazwa                    TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL CHECK(rodzaj IN ('cli','api','sdk')),
    identyfikator_zewnetrzny TEXT,
    poswiadczenie_odwolanie  TEXT,
    katalog_konfiguracji     TEXT,
    aktywne                  INTEGER NOT NULL DEFAULT 1 CHECK(aktywne IN (0,1)),
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    dostawca                 TEXT    NOT NULL DEFAULT '',
    model_domyslny           TEXT,
    adres_bazowy             TEXT,
    domyslne                 INTEGER NOT NULL DEFAULT 0 CHECK(domyslne IN (0,1)),
    stan                     TEXT    NOT NULL DEFAULT 'aktywne'
                                     CHECK(stan IN ('aktywne','wyczerpane','zawieszone')),
    wyczerpane_do            TEXT,
    konto_id                 INTEGER
);

INSERT INTO konto_nowa
    (id, nazwa, rodzaj, identyfikator_zewnetrzny, poswiadczenie_odwolanie,
     katalog_konfiguracji, aktywne, kolejnosc, utworzono, zaktualizowano, dostawca,
     model_domyslny, adres_bazowy, domyslne, stan, wyczerpane_do, konto_id)
SELECT id, nazwa, rodzaj, identyfikator_zewnetrzny, poswiadczenie_odwolanie,
       katalog_konfiguracji, aktywne, kolejnosc, utworzono, zaktualizowano, dostawca,
       model_domyslny, adres_bazowy, domyslne, stan, wyczerpane_do, konto_id
  FROM konto;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 494, 'konto', (SELECT COUNT(*) FROM konto) - (SELECT COUNT(*) FROM konto_nowa);

DROP TABLE konto;
ALTER TABLE konto_nowa RENAME TO konto;

CREATE INDEX idx_konto_rotacja ON konto(rodzaj, aktywne, kolejnosc);

CREATE UNIQUE INDEX idx_konto_domyslne_rodzaj
    ON konto (rodzaj, COALESCE(konto_id, 0)) WHERE domyslne = 1;

CREATE UNIQUE INDEX idx_konto_nazwa
    ON konto (nazwa, COALESCE(konto_id, 0));

-- ── etykieta słownika biblioteki ────────────────────────────────────────────

CREATE TABLE etykieta_slownika_biblioteki_nowa (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    nazwa     TEXT NOT NULL,
    -- Barwa jest kodem żetonu interfejsu, nie wartością szesnastkową.
    barwa     TEXT,
    utworzono TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    konto_id  INTEGER
);

INSERT INTO etykieta_slownika_biblioteki_nowa (nazwa, barwa, utworzono, konto_id)
SELECT nazwa, barwa, utworzono, konto_id FROM etykieta_slownika_biblioteki;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 494, 'etykieta_slownika_biblioteki',
       (SELECT COUNT(*) FROM etykieta_slownika_biblioteki)
     - (SELECT COUNT(*) FROM etykieta_slownika_biblioteki_nowa);

DROP TABLE etykieta_slownika_biblioteki;
ALTER TABLE etykieta_slownika_biblioteki_nowa RENAME TO etykieta_slownika_biblioteki;

CREATE UNIQUE INDEX idx_etykieta_slownika_biblioteki_nazwa
    ON etykieta_slownika_biblioteki (nazwa, COALESCE(konto_id, 0));

PRAGMA foreign_keys = on;
