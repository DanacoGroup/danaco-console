-- Migracja 201 — moduł Apps, dobudowa Architecture Designera: historia wersji
-- układu i adnotacje projektowe.
--
-- Migracja 051 dała architekturze kolumnę `wersja`, ale nie dała historii:
-- numer rósł przy każdym zapisie, a poprzedni układ przestawał istnieć.
-- `apps.architecture.version.list` oddaje `AppArchitectureVersion[]` z liczbą
-- komponentów i zdaniem o różnicy wobec wersji poprzedniej — żadnej z tych
-- dwóch wartości nie da się odtworzyć z wiersza bieżącego, bo poprzedni układ
-- został nadpisany. Wiersz historii powstaje więc przy zapisie definicji, w tej
-- samej transakcji co wymiana komponentów.
--
-- Historia trzyma liczbę i zdanie, nie kopię układu. Kontrakt
-- `AppArchitectureVersion` niesie `componentCount` i `diffSummary`, i nic
-- ponadto — kopia kompletu komponentów przy każdym zapisie rosłaby bez granicy,
-- a kontrakt nie ma czym jej oddać. Odtworzenia wersji wcześniejszej ta tabela
-- nie obiecuje i obiecywać nie może.
--
-- Adnotacja przypina się do komponentu ALBO do zależności, nigdy do obu naraz.
-- Kontrakt daje trzy pola opcjonalne (`componentId`, `dependencyFrom`,
-- `dependencyTo`) i jedno wymagane (`text`); warunek CHECK niżej pilnuje, żeby
-- adnotacja niosła jedno przypięcie albo żadne (notatka wolna kanwy), a nie dwa
-- sprzeczne. Wskazania są kodami zewnętrznymi komponentów, nie więzami obcymi:
-- notatka ma prawo przeżyć komponent zdjęty z układu — to samo rozstrzygnięcie,
-- co przy `plik_warsztatu_apps.komponent_id` (migracja 051).

-- ── Wersja układu komponentów ────────────────────────────────────────────────
CREATE TABLE wersja_architektury_apps (
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    architektura_id    INTEGER NOT NULL REFERENCES architektura_apps(id) ON DELETE CASCADE,
    wersja             INTEGER NOT NULL,
    liczba_komponentow INTEGER NOT NULL DEFAULT 0,
    -- Zdanie o różnicy wobec wersji poprzedniej; pierwsza wersja go nie ma.
    roznica            TEXT,
    utworzono          TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (architektura_id, wersja)
);
CREATE INDEX idx_wersja_architektury_apps ON wersja_architektury_apps(architektura_id, wersja DESC);

-- ── Adnotacja projektowa kanwy ───────────────────────────────────────────────
CREATE TABLE adnotacja_architektury_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    architektura_id          INTEGER NOT NULL REFERENCES architektura_apps(id) ON DELETE CASCADE,
    komponent_kod            TEXT,
    zaleznosc_z              TEXT,
    zaleznosc_do             TEXT,
    tresc                    TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Jedno przypięcie albo żadne — patrz czoło pliku. Zależność wymaga obu
    -- końców: pół krawędzi nie wskazuje niczego na kanwie.
    CHECK (
        (komponent_kod IS NOT NULL AND zaleznosc_z IS NULL AND zaleznosc_do IS NULL)
        OR (komponent_kod IS NULL AND zaleznosc_z IS NOT NULL AND zaleznosc_do IS NOT NULL)
        OR (komponent_kod IS NULL AND zaleznosc_z IS NULL AND zaleznosc_do IS NULL)
    )
);
CREATE INDEX idx_adnotacja_architektury_apps ON adnotacja_architektury_apps(architektura_id, id);
