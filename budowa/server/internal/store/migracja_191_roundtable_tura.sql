-- Migracja 191 — tura i wypowiedź doprowadzone do kształtu kontraktu.
--
-- Trzy braki naraz, wszystkie w warunku CHECK albo w brakującej kolumnie:
--
--  1. Format. Migracja 044 dopuszczała cztery formaty, a kontrakt ma sześć:
--     doszły `delphi` (rundy anonimowe) i `expertPanel` (panel ekspercki).
--     Warunku CHECK nie da się poszerzyć poleceniem ALTER — stąd przebudowa.
--  2. Tura nadrzędna. `roundtable.debate.branch` zakłada wariant tury,
--     a `roundtable.debate.followup` wątek boczny. Obie potrzebują wskazania
--     tury, przy której stoją; bez niego wariant byłby zwykłą kolejną turą.
--  3. Granice tury. Opracowanie modułu ma zegar tury i granicę długości
--     wypowiedzi (2.8.2), a kontrakt pola `timeLimitMs`, `maxStatementChars`
--     i `anonymous`.
--
-- Wypowiedź dostaje redakcję (`roundtable.statement.regenerate` zastępuje
-- treść, a numer redakcji odróżnia zastąpienie od pierwszego głosu), akt mowy
-- (klasyfikacja z `roundtable.analysis.run`), pewność deklarowaną (wejście do
-- kalibracji) i odwołanie do wypowiedzi, na którą odpowiada.
--
-- Przebudowa idzie parami, bo `debata_wypowiedz` wiąże się kluczem obcym
-- z `debata_tura`: przemianowanie tabeli wskazywanej przeciąga za sobą
-- deklarację klucza w tabeli wskazującej, więc obie muszą powstać na nowo
-- w jednym kroku.

-- ── 1. Nowe tabele ───────────────────────────────────────────────────────────

CREATE TABLE debata_tura_nowa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    numer                    INTEGER NOT NULL,
    zagadnienie              TEXT,
    pytanie                  TEXT    NOT NULL DEFAULT '',
    format                   TEXT    NOT NULL DEFAULT 'free'
                                     CHECK(format IN ('free','structured','oxford',
                                                      'roundRobin','delphi','expertPanel')),
    stan                     TEXT    NOT NULL DEFAULT 'open'
                                     CHECK(stan IN ('open','closed')),
    granica_tur              INTEGER NOT NULL DEFAULT 0,
    -- Tura nadrzędna wskazywana kodem, nie kluczem: wariant wskazuje turę
    -- w tym samym oknie, a kod jest tym, czym posługuje się kontrakt.
    tura_nadrzedna           TEXT    NOT NULL DEFAULT '',
    granica_czasu_ms         INTEGER NOT NULL DEFAULT 0,
    granica_znakow           INTEGER NOT NULL DEFAULT 0,
    anonimowa                INTEGER NOT NULL DEFAULT 0 CHECK(anonimowa IN (0,1)),
    rozpoczeto               TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zamknieto                TEXT,
    UNIQUE (okno, numer)
);

INSERT INTO debata_tura_nowa
    (id, identyfikator_zewnetrzny, okno, numer, zagadnienie, pytanie, format, stan,
     granica_tur, rozpoczeto, zamknieto)
SELECT id, identyfikator_zewnetrzny, okno, numer, zagadnienie, pytanie, format, stan,
       granica_tur, rozpoczeto, zamknieto
FROM debata_tura;

CREATE TABLE debata_wypowiedz_nowa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    tura_id                  INTEGER NOT NULL REFERENCES debata_tura_nowa(id) ON DELETE CASCADE,
    uczestnik                TEXT    NOT NULL,
    tresc                    TEXT    NOT NULL DEFAULT '',
    -- Wypowiedź, na którą ta odpowiada (wątek boczny, przesłuchanie krzyżowe).
    odpowiedz_na             TEXT    NOT NULL DEFAULT '',
    -- Akt mowy nadaje analiza; pusty znaczy „jeszcze nieklasyfikowany".
    akt_mowy                 TEXT    NOT NULL DEFAULT ''
                                     CHECK(akt_mowy IN ('','claim','argument','counterArgument',
                                                        'question','concession','rebuttal')),
    -- Pewność ujemna znaczy „nie deklarowano"; zakres deklarowany to 0..1.
    pewnosc                  REAL    NOT NULL DEFAULT -1.0,
    redakcja                 INTEGER NOT NULL DEFAULT 1,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO debata_wypowiedz_nowa
    (id, identyfikator_zewnetrzny, tura_id, uczestnik, tresc, utworzono)
SELECT id, identyfikator_zewnetrzny, tura_id, uczestnik, tresc, utworzono
FROM debata_wypowiedz;

-- ── 2. Usunięcie tabel starych ───────────────────────────────────────────────

DROP TABLE debata_wypowiedz;
DROP TABLE debata_tura;

-- ── 3. Podmiana nazw ─────────────────────────────────────────────────────────

ALTER TABLE debata_tura_nowa      RENAME TO debata_tura;
ALTER TABLE debata_wypowiedz_nowa RENAME TO debata_wypowiedz;

-- ── 4. Odtworzenie indeksów ──────────────────────────────────────────────────

CREATE INDEX idx_debata_tura_okno ON debata_tura(okno, numer DESC);
CREATE INDEX idx_debata_tura_nadrzedna ON debata_tura(tura_nadrzedna);
CREATE INDEX idx_debata_wypowiedz_tura ON debata_wypowiedz(tura_id, id);
CREATE INDEX idx_debata_wypowiedz_uczestnik ON debata_wypowiedz(uczestnik, id);
