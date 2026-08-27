-- Migracja zakłada trwałość niezgodności kontroli jakości panelu, śladu syntezy mowy
-- i śladu eksportu panelu tłumaczenia.

-- Klucz obcy twardy: niezgodność bez panelu, którego dotyczy, nie ma sensu, więc
-- usunięcie panelu zabiera ze sobą jego zastrzeżenia.
CREATE TABLE panel_tlumaczenia_niezgodnosc (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    panel_id  INTEGER NOT NULL REFERENCES panel_tlumaczenia(id) ON DELETE CASCADE,
    rodzaj    TEXT    NOT NULL
                       CHECK(rodzaj IN ('number','date','currency','placeholder','length','omission')),
    segment   TEXT,
    szczegol  TEXT,
    utworzono INTEGER NOT NULL
);
CREATE INDEX idx_panel_tlumaczenia_niezgodnosc_panel ON panel_tlumaczenia_niezgodnosc(panel_id, utworzono DESC);

-- Kolumna odwołania do nagrania jest odwołaniem do pliku dostarczonego z zewnątrz
-- i zostaje pusta, dopóki nagranie realnie nie istnieje.
CREATE TABLE panel_tlumaczenia_synteza_mowy (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    panel_id          INTEGER NOT NULL REFERENCES panel_tlumaczenia(id) ON DELETE CASCADE,
    nagranie_odnosnik TEXT,
    utworzono         INTEGER NOT NULL
);
CREATE INDEX idx_panel_tlumaczenia_synteza_mowy_panel ON panel_tlumaczenia_synteza_mowy(panel_id, utworzono DESC);

-- Odwołanie do pliku wyniku jest wypełnione tylko wtedy, gdy plik realnie powstał,
-- nigdy treścią zmyśloną.
CREATE TABLE panel_tlumaczenia_eksport (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    panel_id       INTEGER NOT NULL REFERENCES panel_tlumaczenia(id) ON DELETE CASCADE,
    format         TEXT    NOT NULL CHECK(format IN ('pdf','docx','markdown','html','txt')),
    plik_odnosnik  TEXT,
    utworzono      INTEGER NOT NULL
);
CREATE INDEX idx_panel_tlumaczenia_eksport_panel ON panel_tlumaczenia_eksport(panel_id, utworzono DESC);
