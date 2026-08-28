-- Migracja 181 — moduł Library: słownik etykiet, tezaurus i hierarchia kolekcji.
--
-- Etykieta była dotąd wolnym tekstem w tabeli złącznikowej
-- (`etykieta_pliku_biblioteki`, migracja 045) i tym pozostaje przy zasobie.
-- Słownik jest bytem NAD tym tekstem: niesie barwę, czas założenia i sam fakt
-- istnienia etykiety, której dziś nie nosi żaden zasób. Bez słownika
-- `library.tag.list` mógłby pokazać wyłącznie etykiety użyte, więc etykieta
-- nieużywana — ta, którą komenda ma umieć wskazać i usunąć — nie istniałaby
-- dla rdzenia wcale.
--
-- Wpis słownika nie jest warunkiem noszenia etykiety. Zasób otagowany
-- `library.tag.set` etykietą nową dostaje ją natychmiast, a wpis słownika
-- powstaje przy okazji; klucz obcy w drugą stronę zamieniłby tagowanie
-- w dwuetapowy obrządek i wywrócił zapis przy wyścigu dwóch wgrań.
--
-- Tezaurus łączy dwie etykiety relacją modelu SKOS. Relacja jest bytem
-- symetrycznym w zapisie (para nazw + rodzaj), a odwrotność wyprowadza odczyt:
-- `nadrzedna` czytana od drugiej strony jest `podrzedna`, więc zapisywanie obu
-- kierunków dałoby dwa wiersze mówiące to samo i rozjazd, gdy zniknie jeden.
--
-- Kolekcje dostają rodzica i regułę. Rodzic daje hierarchię o dowolnej
-- głębokości (`LibraryCollection.parentId`), reguła — kolekcję inteligentną
-- (`LibraryCollection.ruleId`), której zawartość wynika z warunku, a nie
-- z ręcznego przypisania. Kolumna reguły nie ma klucza obcego, bo tabela reguł
-- powstaje krok dalej (182), a kolejność kroków jest jednokierunkowa.

-- ── Słownik etykiet ────────────────────────────────────────────────────────
CREATE TABLE etykieta_slownika_biblioteki (
    nazwa      TEXT PRIMARY KEY,
    -- Barwa jest kodem żetonu interfejsu, nie wartością szesnastkową: paleta
    -- należy do systemu wizualnego, a zapisany `#ff0000` przeżyłby zmianę
    -- palety jako plama obok niej.
    barwa      TEXT,
    utworzono  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- ── Relacja tezaurusa ──────────────────────────────────────────────────────
CREATE TABLE relacja_tezaurusa_biblioteki (
    etykieta_zrodlowa TEXT NOT NULL,
    etykieta_docelowa TEXT NOT NULL,
    rodzaj            TEXT NOT NULL
                      CHECK(rodzaj IN ('nadrzedna','podrzedna','pokrewna')),
    utworzono         TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (etykieta_zrodlowa, etykieta_docelowa, rodzaj)
);
CREATE INDEX idx_relacja_tezaurusa_biblioteki_cel
    ON relacja_tezaurusa_biblioteki(etykieta_docelowa, rodzaj);

-- ── Hierarchia kolekcji i kolekcja inteligentna ────────────────────────────
ALTER TABLE kolekcja_biblioteki
    ADD COLUMN rodzic_id INTEGER REFERENCES kolekcja_biblioteki(id) ON DELETE SET NULL;
ALTER TABLE kolekcja_biblioteki ADD COLUMN regula_kod TEXT;

CREATE INDEX idx_kolekcja_biblioteki_rodzic ON kolekcja_biblioteki(rodzic_id, nazwa);
