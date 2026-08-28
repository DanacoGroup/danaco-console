-- Migracja 003 — jeden silnik kolejek dla pętli sesyjnej i MultitaskingAI.
-- Bieg naprawczy bez limitu: licznik obiegów rośnie, przerwanie należy do użytkownika.

CREATE TABLE kolejka (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    nazwa                  TEXT    NOT NULL,
    rodzaj                 TEXT    NOT NULL DEFAULT 'sesyjna'
                                   CHECK(rodzaj IN ('sesyjna','multitasking')),
    sesja_id               INTEGER REFERENCES sesja(id) ON DELETE CASCADE,
    okno_koordynatora_id   INTEGER REFERENCES okno_komunikacji(id) ON DELETE SET NULL,
    stan                   TEXT    NOT NULL DEFAULT 'bezczynna'
                                   CHECK(stan IN ('bezczynna','pracuje','wstrzymana','zatrzymana')),
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_kolejka_sesja ON kolejka(sesja_id);
CREATE INDEX idx_kolejka_koordynator ON kolejka(okno_koordynatora_id);

CREATE TABLE pozycja_kolejki (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kolejka_id             INTEGER NOT NULL REFERENCES kolejka(id) ON DELETE CASCADE,
    okno_wykonawcy_id      INTEGER REFERENCES okno_komunikacji(id) ON DELETE SET NULL,
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    tytul                  TEXT    NOT NULL,
    tresc_zlecenia         TEXT,
    tresc_odwolanie        TEXT,
    stan                   TEXT    NOT NULL DEFAULT 'oczekuje'
                                   CHECK(stan IN ('oczekuje','przydzielona','wykonywana',
                                                  'do_weryfikacji','ukonczona','bledna','anulowana')),
    werdykt_weryfikacji    TEXT    CHECK(werdykt_weryfikacji IS NULL OR
                                         werdykt_weryfikacji IN ('przyjete','do_poprawy','odrzucone')),
    licznik_obiegow        INTEGER NOT NULL DEFAULT 0,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_pozycja_kolejki_kolejka ON pozycja_kolejki(kolejka_id, kolejnosc, id);
CREATE INDEX idx_pozycja_kolejki_wykonawca ON pozycja_kolejki(okno_wykonawcy_id);

-- Dziennik akcji kolejki: zapis przejść stanu pozycji wraz z numerem obiegu.
CREATE TABLE log_akcji_kolejki (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kolejka_id             INTEGER NOT NULL REFERENCES kolejka(id) ON DELETE CASCADE,
    pozycja_kolejki_id     INTEGER REFERENCES pozycja_kolejki(id) ON DELETE CASCADE,
    akcja                  TEXT    NOT NULL,
    stan_przed             TEXT,
    stan_po                TEXT,
    numer_obiegu           INTEGER NOT NULL DEFAULT 0,
    szczegoly              TEXT,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_log_akcji_kolejki_kolejka ON log_akcji_kolejki(kolejka_id, id);
CREATE INDEX idx_log_akcji_kolejki_pozycja ON log_akcji_kolejki(pozycja_kolejki_id);
