-- Migracja 182 — moduł Library: reguły repozytorium.
--
-- Jedna tabela na trzy rodzaje reguły (kolekcja inteligentna, reguła napływu,
-- folder obserwowany), bo wszystkie trzy mają ten sam kształt: warunek, cel
-- i przełącznik czynności. Rozdzielenie ich na trzy tabele powieliłoby warunek
-- i wykaz — a `library.rule.list` pyta o wszystkie naraz, zawężając rodzajem.
--
-- Warunek stoi w zapisie JSON, nie w kolumnach. Kontrakt niesie go jako `json`
-- (`LibraryRule.condition`: moduł źródłowy, etykiety, rodzaj treści, zakres
-- dat), a kolumna na każdy człon warunku zamieniłaby dodanie członu w migrację.
-- Rozbiór warunku należy do rdzenia, który go stosuje.
--
-- Przypisanie do kolekcji dostaje pochodzenie. Bez niego usunięcie reguły
-- kolekcji inteligentnej nie miałoby jak odróżnić zasobu wciągniętego regułą od
-- zasobu przypisanego ręką Operatora — a `library.rule.remove` z żądaniem
-- `detachFiles` ma zdjąć wyłącznie te pierwsze. Przypisanie ręczne jest
-- wartością domyślną, więc wiersze zastane opisują się same.

CREATE TABLE regula_biblioteki (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT NOT NULL UNIQUE,
    rodzaj               TEXT    NOT NULL
                         CHECK(rodzaj IN ('kolekcja','naplyw','obserwacja')),
    nazwa                TEXT    NOT NULL,
    -- Warunek w zapisie JSON; pusty obiekt znaczy regułę bez zawężenia.
    warunek              TEXT    NOT NULL DEFAULT '{}',
    kolekcja_docelowa_kod TEXT,
    -- Ścieżka folderu obserwowanego; wypełniona wyłącznie dla rodzaju
    -- `obserwacja`.
    sciezka_obserwowana  TEXT,
    czynna               INTEGER NOT NULL DEFAULT 1 CHECK(czynna IN (0,1)),
    ostatnie_przeliczenie TEXT,
    utworzono            TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_regula_biblioteki_rodzaj ON regula_biblioteki(rodzaj, nazwa);

ALTER TABLE przypisanie_kolekcji_biblioteki
    ADD COLUMN zrodlo TEXT NOT NULL DEFAULT 'reczne'
        CHECK(zrodlo IN ('reczne','regula'));
