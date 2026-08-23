-- Migracja 317 — ramki makiety, więzy responsywne i siatki układu modułu
-- Design (grupa D opracowania: `design.frame.*`, `design.constraint.set`,
-- `design.grid.set`).
--
-- ── Dlaczego ramka jest bytem obok kompozycji ───────────────────────────────
-- Kompozycja jest płótnem, ramka jest ekranem. To ramka, a nie płótno,
-- wyznacza obszar wydania i to jej rozmiar zmienia się przy sprawdzaniu układu
-- na innym urządzeniu. Bez własnego bytu makieta czterech ekranów obok siebie
-- byłaby jedną powierzchnią bez granic, a `design.frame.resize.apply` nie
-- miałoby czego przeliczyć.
--
-- ── Dlaczego przynależność warstwy do ramki leży osobno ─────────────────────
-- `design.board.update` przepisuje komplet warstw kompozycji od nowa (usuń
-- i wstaw). Kolumna `ramka_id` dołożona do `warstwa_kompozycji_design` ginęłaby
-- przy każdym zapisie planszy — Operator przesunąłby jedną warstwę i cała
-- makieta wypadłaby ze swoich ramek. Przynależność wiąże się więc
-- z identyfikatorem ZEWNĘTRZNYM warstwy, który klient nadsyła z powrotem
-- i który ten zapis przeżywa.
--
-- Warstwa należy najwyżej do jednej ramki — stąd UNIQUE na samym kodzie
-- warstwy, a nie na parze. Warstwa w dwóch ramkach naraz musiałaby przy
-- `design.frame.resize.apply` przyjąć dwa różne położenia.

CREATE TABLE ramka_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    nazwa                    TEXT    NOT NULL,
    szerokosc                REAL    NOT NULL,
    wysokosc                 REAL    NOT NULL,
    x                        REAL,
    y                        REAL,
    nastawa_urzadzenia       TEXT,
    siatka_json              TEXT,
    uklad_json               TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_ramka_design_kompozycja ON ramka_design(kompozycja_id, id);

CREATE TABLE warstwa_ramki_design (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    ramka_id    INTEGER NOT NULL REFERENCES ramka_design(id) ON DELETE CASCADE,
    warstwa_kod TEXT    NOT NULL UNIQUE,
    kolejnosc   INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_warstwa_ramki_design_ramka ON warstwa_ramki_design(ramka_id, kolejnosc, id);

-- Więz responsywny mówi, co warstwa robi przy zmianie rozmiaru ramki. Jeden
-- wiersz na warstwę w ramce: kotwica pozioma i pionowa razem, bo przeliczenie
-- bierze obie naraz i rozdzielenie ich na dwa wiersze dawałoby stan, w którym
-- warstwa ma pół więzu.
CREATE TABLE wiez_ramki_design (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    ramka_id    INTEGER NOT NULL REFERENCES ramka_design(id) ON DELETE CASCADE,
    warstwa_kod TEXT    NOT NULL,
    poziomo     TEXT    NOT NULL,
    pionowo     TEXT    NOT NULL,
    UNIQUE(ramka_id, warstwa_kod)
);
CREATE INDEX idx_wiez_ramki_design_ramka ON wiez_ramki_design(ramka_id, id);

-- Siatka obowiązująca całą kompozycję. Siatka ramki leży przy ramce
-- (`ramka_design.siatka_json`), bo należy do ekranu; siatka płótna nie ma do
-- czego przylgnąć i dostaje własny wiersz — jeden na kompozycję.
CREATE TABLE siatka_kompozycji_design (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    kompozycja_id  INTEGER NOT NULL UNIQUE REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    siatka_json    TEXT    NOT NULL,
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
