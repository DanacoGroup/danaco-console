-- Ciąg dalszy kroków 490 i 491. Nazwa hosta zdalnego, ścieżka klucza terminala,
-- odwołanie sekretu rozszerzenia oraz para rozszerzenie-wersja i nadanie
-- uprawnienia to klucze wpisywane albo wskazywane przez Operatora. Więz
-- obejmujący całą tabelę odbierał je kontu młodszemu i pozwalał upsertowi
-- sięgnąć wiersza konta cudzego.
--
-- Postępowanie i uzasadnienie wskaźnika po wyrażeniu — jak w kroku 490.

CREATE TABLE IF NOT EXISTS kontrola_przejazdu (
    wersja  INTEGER NOT NULL,
    tabela  TEXT    NOT NULL,
    roznica INTEGER NOT NULL CHECK (roznica = 0)
);

PRAGMA foreign_keys = off;

-- ── host zdalny ─────────────────────────────────────────────────────────────

CREATE TABLE host_zdalny_nowa (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    nazwa          TEXT    NOT NULL,
    adres          TEXT    NOT NULL DEFAULT '',
    uzytkownik     TEXT    NOT NULL DEFAULT '',
    port           INTEGER NOT NULL DEFAULT 22 CHECK(port BETWEEN 1 AND 65535),
    zgoda          INTEGER NOT NULL DEFAULT 0 CHECK(zgoda IN (0,1)),
    zgode_wydano   TEXT    NOT NULL DEFAULT '',
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    klucz_hosta    TEXT    NOT NULL DEFAULT '',
    konto_id       INTEGER
);

INSERT INTO host_zdalny_nowa
    (id, nazwa, adres, uzytkownik, port, zgoda, zgode_wydano, utworzono,
     zaktualizowano, klucz_hosta, konto_id)
SELECT id, nazwa, adres, uzytkownik, port, zgoda, zgode_wydano, utworzono,
       zaktualizowano, klucz_hosta, konto_id
  FROM host_zdalny;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 492, 'host_zdalny',
       (SELECT COUNT(*) FROM host_zdalny) - (SELECT COUNT(*) FROM host_zdalny_nowa);

DROP TABLE host_zdalny;
ALTER TABLE host_zdalny_nowa RENAME TO host_zdalny;

CREATE UNIQUE INDEX idx_host_zdalny_nazwa
    ON host_zdalny (nazwa, COALESCE(konto_id, 0));

-- ── klucz SSH terminala ─────────────────────────────────────────────────────

CREATE TABLE terminal_klucz_nowa (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    kod         TEXT    NOT NULL UNIQUE,
    nazwa       TEXT    NOT NULL,
    rodzaj      TEXT    NOT NULL CHECK(rodzaj IN ('ed25519','rsa','ecdsa')),
    odcisk      TEXT    NOT NULL DEFAULT '',
    klucz_jawny TEXT    NOT NULL DEFAULT '',
    sciezka     TEXT    NOT NULL,
    haslo       INTEGER NOT NULL DEFAULT 0 CHECK(haslo IN (0,1)),
    utworzono   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    konto_id    INTEGER
);

INSERT INTO terminal_klucz_nowa
    (id, kod, nazwa, rodzaj, odcisk, klucz_jawny, sciezka, haslo, utworzono, konto_id)
SELECT id, kod, nazwa, rodzaj, odcisk, klucz_jawny, sciezka, haslo, utworzono, konto_id
  FROM terminal_klucz;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 492, 'terminal_klucz',
       (SELECT COUNT(*) FROM terminal_klucz) - (SELECT COUNT(*) FROM terminal_klucz_nowa);

DROP TABLE terminal_klucz;
ALTER TABLE terminal_klucz_nowa RENAME TO terminal_klucz;

CREATE UNIQUE INDEX idx_terminal_klucz_sciezka
    ON terminal_klucz (sciezka, COALESCE(konto_id, 0));

-- ── sekret rozszerzenia ─────────────────────────────────────────────────────

CREATE TABLE sekret_rozszerzenia_nowa (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    odwolanie        TEXT    NOT NULL,
    etykieta         TEXT,
    sposob_logowania TEXT CHECK(sposob_logowania IN ('oauth2','apiKey','token','basic','none')),
    wygasa           INTEGER,
    zaktualizowano   INTEGER NOT NULL DEFAULT 0,
    konto_id         INTEGER
);

INSERT INTO sekret_rozszerzenia_nowa
    (id, odwolanie, etykieta, sposob_logowania, wygasa, zaktualizowano, konto_id)
SELECT id, odwolanie, etykieta, sposob_logowania, wygasa, zaktualizowano, konto_id
  FROM sekret_rozszerzenia;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 492, 'sekret_rozszerzenia',
       (SELECT COUNT(*) FROM sekret_rozszerzenia)
     - (SELECT COUNT(*) FROM sekret_rozszerzenia_nowa);

DROP TABLE sekret_rozszerzenia;
ALTER TABLE sekret_rozszerzenia_nowa RENAME TO sekret_rozszerzenia;

CREATE UNIQUE INDEX idx_sekret_rozszerzenia_odwolanie
    ON sekret_rozszerzenia (odwolanie, COALESCE(konto_id, 0));

-- ── wersja rozszerzenia ─────────────────────────────────────────────────────

CREATE TABLE wersja_rozszerzenia_nowa (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    rozszerzenie_kod TEXT    NOT NULL,
    wersja           TEXT    NOT NULL,
    dziennik_zmian   TEXT,
    paczka_odwolanie TEXT,
    utworzono        INTEGER NOT NULL DEFAULT 0,
    konto_id         INTEGER
);

INSERT INTO wersja_rozszerzenia_nowa
    (id, rozszerzenie_kod, wersja, dziennik_zmian, paczka_odwolanie, utworzono, konto_id)
SELECT id, rozszerzenie_kod, wersja, dziennik_zmian, paczka_odwolanie, utworzono, konto_id
  FROM wersja_rozszerzenia;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 492, 'wersja_rozszerzenia',
       (SELECT COUNT(*) FROM wersja_rozszerzenia)
     - (SELECT COUNT(*) FROM wersja_rozszerzenia_nowa);

DROP TABLE wersja_rozszerzenia;
ALTER TABLE wersja_rozszerzenia_nowa RENAME TO wersja_rozszerzenia;

CREATE UNIQUE INDEX idx_wersja_rozszerzenia_wydanie
    ON wersja_rozszerzenia (rozszerzenie_kod, wersja, COALESCE(konto_id, 0));

-- ── uprawnienie rozszerzenia ────────────────────────────────────────────────

CREATE TABLE uprawnienie_rozszerzenia_nowa (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    rozszerzenie_kod TEXT    NOT NULL,
    zakres           TEXT    NOT NULL
                             CHECK(zakres IN ('network','fileRead','fileWrite',
                                              'processSpawn','secretRead','modelCall')),
    byt              TEXT,
    tryb             TEXT CHECK(tryb IN ('read','write')),
    objasnienie      TEXT,
    nadane           INTEGER NOT NULL DEFAULT 0 CHECK(nadane IN (0,1)),
    agent_kod        TEXT,
    nadano           INTEGER,
    konto_id         INTEGER
);

INSERT INTO uprawnienie_rozszerzenia_nowa
    (id, rozszerzenie_kod, zakres, byt, tryb, objasnienie, nadane, agent_kod,
     nadano, konto_id)
SELECT id, rozszerzenie_kod, zakres, byt, tryb, objasnienie, nadane, agent_kod,
       nadano, konto_id
  FROM uprawnienie_rozszerzenia;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica)
SELECT 492, 'uprawnienie_rozszerzenia',
       (SELECT COUNT(*) FROM uprawnienie_rozszerzenia)
     - (SELECT COUNT(*) FROM uprawnienie_rozszerzenia_nowa);

DROP TABLE uprawnienie_rozszerzenia;
ALTER TABLE uprawnienie_rozszerzenia_nowa RENAME TO uprawnienie_rozszerzenia;

CREATE INDEX idx_uprawnienie_rozszerzenia
    ON uprawnienie_rozszerzenia(rozszerzenie_kod, nadane);

CREATE UNIQUE INDEX idx_uprawnienie_rozszerzenia_nadanie
    ON uprawnienie_rozszerzenia (rozszerzenie_kod, zakres, byt, nadane, agent_kod,
                                 COALESCE(konto_id, 0));

PRAGMA foreign_keys = on;
