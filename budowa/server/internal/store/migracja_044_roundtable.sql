-- Migracja 044 tworzy trwałość modułu roundtable: uczestnicy debaty, tury,
-- wypowiedzi i stanowisko końcowe, powiązane oknem, nie sesją.

-- Tabela debata_uczestnik przechowuje uczestnika debaty w oknie wraz z kanałem
-- modelu, tożsamością i kolejnością głosu w turze.
CREATE TABLE debata_uczestnik (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    kanal_modelu             TEXT    NOT NULL,
    nazwa_tozsamosci         TEXT,
    prompt_systemowy         TEXT,
    -- Wyciszony uczestnik zostaje w panelu, ale pytanie w turze do niego nie idzie.
    wyciszony                INTEGER NOT NULL DEFAULT 0 CHECK(wyciszony IN (0,1)),
    -- Zero w kolejności znaczy brak wskazania, wtedy porządkiem jest dołączenie; kluczowy tu nie ma.
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_uczestnik_okno ON debata_uczestnik(okno, kolejnosc, id);

-- Tabela debata_tura przechowuje turę debaty w oknie, jej format i stan jako
-- wartości kontraktu, bez słownika pośredniego.
CREATE TABLE debata_tura (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    numer                    INTEGER NOT NULL,
    zagadnienie              TEXT,
    -- Pytanie tury zostaje w zapisie, żeby transkrypt nie był zbiorem odpowiedzi na nieznane.
    pytanie                  TEXT    NOT NULL DEFAULT '',
    format                   TEXT    NOT NULL DEFAULT 'free'
                                     CHECK(format IN ('free','structured','oxford','roundRobin')),
    stan                     TEXT    NOT NULL DEFAULT 'open'
                                     CHECK(stan IN ('open','closed')),
    -- Granica liczby tur ustawiona przy uruchomieniu; zero znaczy „bez granicy”.
    granica_tur              INTEGER NOT NULL DEFAULT 0,
    rozpoczeto               TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zamknieto                TEXT,
    UNIQUE (okno, numer)
);
CREATE INDEX idx_debata_tura_okno ON debata_tura(okno, numer DESC);

-- Tabela debata_wypowiedz przechowuje wypowiedź uczestnika w turze debaty wraz
-- z treścią i chwilą jej powstania.
CREATE TABLE debata_wypowiedz (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    tura_id                  INTEGER NOT NULL REFERENCES debata_tura(id) ON DELETE CASCADE,
    uczestnik                TEXT    NOT NULL,
    tresc                    TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_wypowiedz_tura ON debata_wypowiedz(tura_id, id);

-- Tabela debata_stanowisko przechowuje stanowisko końcowe debaty w oknie,
-- z pustym kodem tury oznaczającym stanowisko całej debaty.
CREATE TABLE debata_stanowisko (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tura                     TEXT    NOT NULL DEFAULT '',
    tresc                    TEXT,
    -- Wersja rośnie przy każdym złożeniu stanowiska, żeby panel mógł pokazać jego historię.
    wersja                   INTEGER NOT NULL DEFAULT 1,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (okno, tura)
);
