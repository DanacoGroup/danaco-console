-- Migracja zakłada trwałość modułu asystenta: zlecenie wieloetapowe wydane poleceniem
-- oraz chronologiczny dziennik działań asystenta.

-- ── Zlecenie wieloetapowe asystenta (Actions Monitor) ────────────────────────
-- Wartości kolumn `stan` i `droga` są wartościami kontraktu
-- (AssistantActionStatus, AssistantOrigin) wprost, bez tłumaczenia.
CREATE TABLE zlecenie_asystenta (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno_kod                 TEXT    NOT NULL,
    tytul                    TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'queued'
                                     CHECK(stan IN ('queued','running','paused','done','failed','cancelled')),
    droga                    TEXT    NOT NULL DEFAULT 'text'
                                     CHECK(droga IN ('voice','text')),
    etap_biezacy             INTEGER,
    liczba_etapow            INTEGER,
    priorytet                INTEGER,
    wynik                    TEXT,
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_zlecenie_asystenta_okno ON zlecenie_asystenta(okno_kod, zaktualizowano DESC);
CREATE INDEX idx_zlecenie_asystenta_stan ON zlecenie_asystenta(stan, zaktualizowano DESC);

-- ── Dziennik działań asystenta (Actions Monitor / rozmowa) ───────────────────
-- Klucz obcy jest miękki (`ON DELETE SET NULL`): usunięcie zlecenia nie ma
-- prawa wymazać śladu rozmowy, która realnie się odbyła.
CREATE TABLE wpis_dziennika_asystenta (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    kod               TEXT    NOT NULL UNIQUE,
    okno_kod          TEXT    NOT NULL,
    zlecenie_kod      TEXT    REFERENCES zlecenie_asystenta(identyfikator_zewnetrzny) ON DELETE SET NULL,
    rodzaj            TEXT    NOT NULL DEFAULT 'note'
                               CHECK(rodzaj IN ('command','result','note')),
    tresc             TEXT    NOT NULL,
    nagranie_odnosnik TEXT,
    utworzono         INTEGER NOT NULL
);
CREATE INDEX idx_wpis_dziennika_asystenta_okno ON wpis_dziennika_asystenta(okno_kod, utworzono DESC);
CREATE INDEX idx_wpis_dziennika_asystenta_zlecenie ON wpis_dziennika_asystenta(zlecenie_kod, utworzono DESC);
