-- Migracja 001 zakłada fundament modelu danych Danaco Console; sekrety nigdy
-- nie trafiają do bazy, a kolumny odwołania niosą wyłącznie wskazanie na nie.

-- ── Konta kanałów modelu (rotacja kont Code CLI oraz API) ──────────────
-- Rotacja pozwala rozkładać wywołania modelu między kontami bez zmiany
-- konfiguracji klienta.
CREATE TABLE konto (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    nazwa                  TEXT    NOT NULL UNIQUE,
    rodzaj                 TEXT    NOT NULL CHECK(rodzaj IN ('cli','api','sdk')),
    identyfikator_zewnetrzny TEXT,
    poswiadczenie_odwolanie  TEXT,
    katalog_konfiguracji   TEXT,
    aktywne                INTEGER NOT NULL DEFAULT 1 CHECK(aktywne IN (0,1)),
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- ── Urządzenia z zainstalowanym klientem i agentem lokalnym ──────────
-- Wiersz łączy sprzęt z historią połączeń i zaufaniem nadanym przez
-- Operatora.
CREATE TABLE urzadzenie (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    nazwa                  TEXT    NOT NULL,
    identyfikator_sprzetowy TEXT   NOT NULL UNIQUE,
    system_operacyjny      TEXT,
    wersja_klienta         TEXT,
    zaufane                INTEGER NOT NULL DEFAULT 1 CHECK(zaufane IN (0,1)),
    ostatnio_widziane      TEXT,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- ── Połączenia WebSocket klienta z rdzeniem ──────────────────────────
-- Wiersz trzyma jedno połączenie od nawiązania do zamknięcia lub zerwania
-- łącza.
CREATE TABLE polaczenie (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    urzadzenie_id          INTEGER NOT NULL REFERENCES urzadzenie(id) ON DELETE CASCADE,
    konto_id               INTEGER REFERENCES konto(id) ON DELETE SET NULL,
    identyfikator_polaczenia TEXT  NOT NULL UNIQUE,
    adres_zdalny           TEXT,
    stan                   TEXT    NOT NULL DEFAULT 'otwarte'
                                   CHECK(stan IN ('otwarte','zamkniete','zerwane')),
    nawiazano              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono             TEXT
);
CREATE INDEX idx_polaczenie_urzadzenie ON polaczenie(urzadzenie_id);
CREATE INDEX idx_polaczenie_konto ON polaczenie(konto_id);

-- ── Środowiska produktu: TalkIn, WorkSpace, CodeStudio, MultitaskingAI ─────────
-- Każdy kod środowiska jest zamkniętą wartością wyliczeniową rdzenia.
CREATE TABLE srodowisko (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                    TEXT    NOT NULL UNIQUE
                                   CHECK(kod IN ('talkin','workspace','codestudio','multitaskingai')),
    nazwa                  TEXT    NOT NULL,
    opis                   TEXT,
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    aktywne                INTEGER NOT NULL DEFAULT 1 CHECK(aktywne IN (0,1))
);

-- ── Moduły (jednostki funkcjonalne osadzane w środowiskach) ────────────────────
-- Moduł istnieje niezależnie od środowiska; osadzenie opisuje osobna tabela
-- złącznikowa.
CREATE TABLE modul (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                    TEXT    NOT NULL UNIQUE,
    nazwa                  TEXT    NOT NULL,
    opis                   TEXT,
    ikona                  TEXT,
    aktywny                INTEGER NOT NULL DEFAULT 1 CHECK(aktywny IN (0,1))
);

-- ── Osadzenie modułu w środowisku (macierz widoczności) ────────────────────────
-- Wiersz określa, czy dany moduł jest widoczny w danym środowisku i w jakiej
-- kolejności.
CREATE TABLE srodowisko_modul (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    srodowisko_id          INTEGER NOT NULL REFERENCES srodowisko(id) ON DELETE CASCADE,
    modul_id               INTEGER NOT NULL REFERENCES modul(id) ON DELETE CASCADE,
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    widoczny               INTEGER NOT NULL DEFAULT 1 CHECK(widoczny IN (0,1)),
    UNIQUE(srodowisko_id, modul_id)
);
CREATE INDEX idx_srodowisko_modul_modul ON srodowisko_modul(modul_id);

-- ── Otwarty rejestr kanałów modelu sterowany danymi ─────────
-- Nowy dostawca lub model wchodzi wierszem danych, bez zmiany kodu rdzenia.
CREATE TABLE kanal_modelu (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                    TEXT    NOT NULL UNIQUE,
    nazwa                  TEXT    NOT NULL,
    dostawca               TEXT    NOT NULL,
    identyfikator_modelu   TEXT    NOT NULL,
    rodzaj_kanalu          TEXT    NOT NULL
                                   CHECK(rodzaj_kanalu IN ('cli','api','sdk','lokalny')),
    konto_id               INTEGER REFERENCES konto(id) ON DELETE SET NULL,
    poswiadczenie_odwolanie TEXT,
    parametry_json         TEXT    NOT NULL DEFAULT '{}',
    multimodalny           INTEGER NOT NULL DEFAULT 0 CHECK(multimodalny IN (0,1)),
    aktywny                INTEGER NOT NULL DEFAULT 1 CHECK(aktywny IN (0,1)),
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_kanal_modelu_konto ON kanal_modelu(konto_id);

-- ── Osiem poziomów zasięgu konfiguracji ──────────────────────────────
-- Wyższe pierwszeństwo wygrywa; okno komunikacji jest poziomem najwęższym.
CREATE TABLE poziom_zasiegu (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                    TEXT    NOT NULL UNIQUE
                                   CHECK(kod IN ('globalny','srodowisko','modul','para_modulow',
                                                 'projekt','karta_sesji','rola','okno')),
    nazwa                  TEXT    NOT NULL,
    pierwszenstwo          INTEGER NOT NULL UNIQUE
);

INSERT INTO poziom_zasiegu (kod, nazwa, pierwszenstwo) VALUES
    ('globalny',     'Globalny',            1),
    ('srodowisko',   'Środowisko',          2),
    ('modul',        'Moduł',               3),
    ('para_modulow', 'Para modułów',        4),
    ('projekt',      'Projekt',             5),
    ('karta_sesji',  'Karta sesji',         6),
    ('rola',         'Rola',                7),
    ('okno',         'Okno komunikacji',    8);

-- ── Ustawienia rozstrzygane wg poziomu zasięgu ───────────────────────
-- klucz_zasiegu wskazuje konkretny byt poziomu ('' dla poziomu globalnego).
CREATE TABLE ustawienie (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    poziom_zasiegu_id      INTEGER NOT NULL REFERENCES poziom_zasiegu(id) ON DELETE CASCADE,
    klucz_zasiegu          TEXT    NOT NULL DEFAULT '',
    klucz                  TEXT    NOT NULL,
    wartosc                TEXT,
    rodzaj_wartosci        TEXT    NOT NULL DEFAULT 'tekst'
                                   CHECK(rodzaj_wartosci IN ('tekst','liczba','logiczna','json')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(poziom_zasiegu_id, klucz_zasiegu, klucz)
);
CREATE INDEX idx_ustawienie_klucz ON ustawienie(klucz);

-- ── Karta sesji — kontener sesji w obrębie środowiska ──────────────────────────
-- Karta grupuje sesje należące do tego samego wątku pracy w jednym
-- środowisku.
CREATE TABLE karta_sesji (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    srodowisko_id          INTEGER NOT NULL REFERENCES srodowisko(id) ON DELETE CASCADE,
    nazwa                  TEXT    NOT NULL,
    opis                   TEXT,
    przypieta              INTEGER NOT NULL DEFAULT 0 CHECK(przypieta IN (0,1)),
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_karta_sesji_srodowisko ON karta_sesji(srodowisko_id);

-- ── Sesja — wspólna dla plików, pamięci, projektu i agentów.
-- Sesja nie ma przypisania do modułu; moduł jest atrybutem okna komunikacji.
CREATE TABLE sesja (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    karta_sesji_id         INTEGER NOT NULL REFERENCES karta_sesji(id) ON DELETE CASCADE,
    tytul                  TEXT    NOT NULL,
    projekt                TEXT,
    stan                   TEXT    NOT NULL DEFAULT 'aktywna'
                                   CHECK(stan IN ('aktywna','wstrzymana','zakonczona','archiwalna')),
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono             TEXT
);
CREATE INDEX idx_sesja_karta ON sesja(karta_sesji_id);
