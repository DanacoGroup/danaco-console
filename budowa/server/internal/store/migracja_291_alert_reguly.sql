-- Migracja 291 — reguły wyzwalania i rejestr wyzwoleń (rodzina `alert.*`).
--
-- Reguła mówi, KIEDY produkt ma zawołać; wyzwolenie jest zapisem tego, że
-- zawołał — z wartością zmierzoną w chwili wyzwolenia. Wartość obserwowana
-- (`wartosc_obserwowana`) leży w wierszu wyzwolenia, a nie w regule, bo reguła
-- trwa, a pomiar dotyczy jednej chwili.
--
-- Drogi dostarczenia (AlertChannel) idą zapisem strukturalnym w jednej kolumnie
-- (`kanaly_json`), bo są wykazem wartości bez własnych atrybutów; ta sama forma
-- co przy `deliveredChannels` wyzwolenia. Osobna tabela wiązania dawałaby
-- złączenie dla listy trzech napisów.
--
-- Wyzwolenie zna sondę i wywołanie modelu wyłącznie kodem, bez klucza obcego:
-- usunięcie sondy nie ma prawa wymazać śladu alertu, który realnie się odbył.

CREATE TABLE regula_alertu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    rodzaj                   TEXT    NOT NULL,
    miara                    TEXT    NOT NULL,
    porownanie               TEXT,
    prog                     REAL,
    okno_ms                  INTEGER NOT NULL,
    waga                     TEXT    NOT NULL,
    kanaly_json              TEXT    NOT NULL DEFAULT '[]',
    zasieg                   TEXT,
    zasieg_kod               TEXT,
    czynna                   INTEGER NOT NULL DEFAULT 1,
    wyciszona_do             INTEGER,
    eskalacja_po_ms          INTEGER,
    adres_zwrotny            TEXT,
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER,
    ostatnie_wyzwolenie      INTEGER
);
CREATE INDEX idx_regula_alertu_miara ON regula_alertu(miara, id);
CREATE INDEX idx_regula_alertu_rodzaj ON regula_alertu(rodzaj, id);

CREATE TABLE wyzwolenie_alertu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    regula_kod               TEXT    NOT NULL
                                     REFERENCES regula_alertu(identyfikator_zewnetrzny) ON DELETE CASCADE,
    stan                     TEXT    NOT NULL,
    waga                     TEXT    NOT NULL,
    miara                    TEXT    NOT NULL,
    wartosc_obserwowana      REAL    NOT NULL,
    prog                     REAL,
    komunikat                TEXT    NOT NULL,
    wyzwolono                INTEGER NOT NULL,
    potwierdzono             INTEGER,
    rozwiazano               INTEGER,
    notatka                  TEXT,
    blad_kod                 TEXT,
    sonda_kod                TEXT,
    wywolanie_kod            TEXT,
    kanaly_dostarczone_json  TEXT    NOT NULL DEFAULT '[]'
);
CREATE INDEX idx_wyzwolenie_alertu_regula ON wyzwolenie_alertu(regula_kod, wyzwolono DESC, id DESC);
CREATE INDEX idx_wyzwolenie_alertu_stan ON wyzwolenie_alertu(stan, wyzwolono DESC, id DESC);
CREATE INDEX idx_wyzwolenie_alertu_czas ON wyzwolenie_alertu(wyzwolono DESC, id DESC);
