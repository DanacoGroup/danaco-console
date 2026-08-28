-- Migracja 230 — kolekcje zasobów modułu Design (Assets Panel).
--
-- Kolekcja jest bytem osobnym od etykiety, choć obie grupują zasoby. Etykieta
-- (`etykieta_zasobu_design`) jest wolnym słowem bez tożsamości: nie ma nazwy
-- własnej, opisu ani porządku, a dwa zasoby z tym samym słowem należą do siebie
-- wyłącznie przez to słowo. Kolekcja ma nazwę, opis i kolejność, więc ma własny
-- wiersz i własny identyfikator zewnętrzny — inaczej zmiana nazwy kolekcji
-- oznaczałaby przepisanie etykiety w każdym zasobie z osobna.
--
-- Przypisanie jest dokładką albo odjęciem, nigdy zastąpieniem
-- (`design.collection.assign`), więc para (kolekcja, zasób) jest kluczem bez
-- surogatu: powtórzone dołożenie tego samego zasobu nie zakłada drugiego
-- wiersza i nie zmienia liczby zasobów kolekcji.
--
-- Zasób wskazywany jest identyfikatorem zewnętrznym (TEXT), nie więzem obcym —
-- ta sama decyzja co przy `warstwa_kompozycji_design.zasob_id` i z tego samego
-- powodu: kolekcja przeżywa usunięcie zasobu, który się w niej znalazł, a
-- `design.collection.list` liczy pozycje, nie zasoby żywe.

CREATE TABLE kolekcja_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- design.collection.list czyta kolekcje okna, od ostatnio zmienianej.
CREATE INDEX idx_kolekcja_design_okno ON kolekcja_design(okno, zaktualizowano DESC, id DESC);

CREATE TABLE pozycja_kolekcji_design (
    kolekcja_id INTEGER NOT NULL REFERENCES kolekcja_design(id) ON DELETE CASCADE,
    zasob_id    TEXT    NOT NULL,
    kolejnosc   INTEGER NOT NULL DEFAULT 0,
    utworzono   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (kolekcja_id, zasob_id)
);
-- design.collection.list zawęża wykaz do kolekcji zawierających wskazany zasób.
CREATE INDEX idx_pozycja_kolekcji_design_zasob ON pozycja_kolekcji_design(zasob_id);
