-- Migracja 363 — blokady fragmentów dokumentu: fragment, którego model nie tknie.
--
-- ── Dlaczego blokada jest w bazie, a nie w oknie ─────────────────────────────
-- Wymaganie rozstrzygające Właściciela: blokada obowiązuje W RDZENIU, nie
-- w oknie. Model woła komendy rdzenia tak samo jak klient, więc blokada
-- pilnowana wyłącznie przez okno jest pozorna — ominąłby ją bez wysiłku. To nie
-- jest wygodniejsze miejsce na to sprawdzenie; to jedyne. Sprawdzenie stoi na
-- drodze każdej komendy zmieniającej dokument, PRZED dotknięciem treści, i musi
-- mieć skąd czytać zakresy — stąd ten wiersz.
--
-- ── Dlaczego zasięg jest kolumną, nie założeniem ────────────────────────────
-- Blokada jest skierowana przeciw modelowi, nie przeciw właścicielowi
-- dokumentu: Operator zmienia fragment zablokowany bez przeszkód. Gdy ma
-- działać także na Operatora, jest to OSOBNE, JAWNE ustawienie blokady, a nie
-- zachowanie domyślne — dlatego `zasieg` domyślnie stoi na `model`.
--
-- ── Dlaczego blokada przechodzi przez wersje i przez szablon ─────────────────
-- Przywrócenie wcześniejszej wersji nie gubi blokad: blokada wisi przy
-- dokumencie, nie przy wersji, więc przywrócenie treści jej nie rusza. Szablon
-- niesie swoje blokady do dokumentów z niego zakładanych — fragmenty wzorcowe
-- pisma mają zostać wzorcowe — i takie blokady znaczy `z_szablonu`, żeby
-- Operator widział, że nie założył ich sam.
--
-- ── Dlaczego przesunięcie zakresu jest tu odnotowywane ──────────────────────
-- Zakres blokady jest liczony w znakach. Wpis PRZED blokadą przesuwa ją w prawo;
-- gdyby zakres został nietknięty, blokada zaczęłaby po pierwszej edycji chronić
-- nie ten fragment, co miała — i to bez słowa. `przesuniecia` liczy, ile razy
-- rdzeń zakres przeliczył: to jest miara zaufania do zakresu, którą Operator
-- widzi w wykazie blokad.

CREATE TABLE blokada_fragmentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    nazwa                    TEXT    NOT NULL,
    powod                    TEXT,
    zakres_od                INTEGER NOT NULL,
    zakres_do                INTEGER NOT NULL,
    zasieg                   TEXT    NOT NULL DEFAULT 'model'
                                     CHECK(zasieg IN ('model','everyone')),
    -- Kto blokadę założył. Zdejmuje ją WYŁĄCZNIE Operator, więc kolumna służy
    -- wykazowi i sprawozdaniu, nie rozstrzyganiu prawa do zdjęcia.
    zalozyl_rodzaj           TEXT    NOT NULL DEFAULT 'uzytkownik'
                                     CHECK(zalozyl_rodzaj IN ('uzytkownik','model')),
    zalozyl_agent_kod        TEXT,
    zalozyl_agent_nazwa      TEXT,
    z_szablonu               INTEGER NOT NULL DEFAULT 0,
    szablon_kod              TEXT,
    przesuniecia             INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    CHECK(zakres_do >= zakres_od)
);
-- Indeks po zakresie, bo najczęstsze pytanie tej tabeli brzmi: „czy zakres
-- [od, do) tego żądania dotyka którejkolwiek blokady tego dokumentu" — i pada
-- ono PRZED każdą zmianą dokumentu, więc musi być tanie.
CREATE INDEX idx_blokada_fragmentu_studio_zakres
    ON blokada_fragmentu_studio(dokument_id, zakres_od, zakres_do);

-- Blokady wzorcowe szablonu pisma. Leżą osobno od blokad dokumentu, bo szablon
-- nie jest dokumentem: nie ma treści bieżącej, a jego blokady są WZORCEM
-- kopiowanym do każdego dokumentu z niego zakładanego. Trzymanie ich w tabeli
-- dokumentu wymagałoby dokumentu-widma dla każdego szablonu.
CREATE TABLE blokada_szablonu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    szablon_kod              TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    powod                    TEXT,
    zakres_od                INTEGER NOT NULL,
    zakres_do                INTEGER NOT NULL,
    zasieg                   TEXT    NOT NULL DEFAULT 'model'
                                     CHECK(zasieg IN ('model','everyone')),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    CHECK(zakres_do >= zakres_od)
);
CREATE INDEX idx_blokada_szablonu_studio_szablon
    ON blokada_szablonu_studio(szablon_kod, zakres_od);
