-- Migracja zakłada trwałość modułu tłumaczenia: okno tłumaczenia z tekstem źródłowym
-- oraz panel tłumaczenia jako byt centralny modułu.

-- Treść źródłowa bywa obszerna, więc trafia do pliku na dysku, a baza trzyma tylko
-- odwołanie do niego.
CREATE TABLE okno_tlumaczenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    tekst_zrodlowy           TEXT,
    tekst_zrodlowy_odwolanie TEXT,
    jezyk_zrodlowy           TEXT,
    liczba_segmentow         INTEGER,
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL
);

-- Wartości kolumny stanu są wartościami kontraktu wprost; ton i treść zwrotna są polami
-- tego samego panelu, nie osobnym bytem.
CREATE TABLE panel_tlumaczenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno_id                  INTEGER NOT NULL REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    jezyk                    TEXT    NOT NULL,
    tresc                    TEXT,
    tresc_odwolanie          TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'pending'
                                     CHECK(stan IN ('pending','translating','ready','error')),
    ton                      TEXT,
    tresc_zwrotna            TEXT,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_panel_tlumaczenia_okno ON panel_tlumaczenia(okno_id, jezyk);
CREATE INDEX idx_panel_tlumaczenia_stan ON panel_tlumaczenia(stan, zaktualizowano DESC);
