-- Migracja 365 — znakowanie fragmentów i rodzaje znaczników własnych Operatora.
--
-- ── Trzy rzeczy, których nie wolno pomieszać ────────────────────────────────
-- Komentarz nie niesie brzmienia. Zmiana śledzona jest JUŻ w treści. Propozycja
-- na marginesie niesie brzmienie i NIE JEST w treści — Operator ją przyjmuje,
-- odrzuca albo poprawia. To trzy różne czynności i okno ma je odróżniać
-- wyraźnie, żeby Operator wiedział, na co patrzy.
--
-- Komentarz ma już swoją tabelę (`komentarz_studio`) i zmiana śledzona swoją
-- (`zmiana_sledzona_studio`). Tutaj leży to, czego nie miały: wyróżnienie barwą,
-- znacznik własny Operatora i propozycja z brzmieniem. Trzy rodzaje w jednej
-- tabeli, bo pracuje się nimi tak samo — wszystkie są przypięte do fragmentu,
-- wszystkie mają autora i stan, i wszystkie wchodzą do jednego wykazu
-- znakowań, po którym Operator skacze i który odhacza.
--
-- ── Dlaczego wyróżnienie jest i tu, i w postaci dokumentu ───────────────────
-- Wyróżnienie tła jest cechą postaci znaku (`highlightColor` w drzewie postaci)
-- — i tam musi być, bo inaczej nie wyszłoby przy wydaniu do docx ani do PDF.
-- Wiersz tutaj jest czymś innym: jest ZNAKOWANIEM, czyli pozycją wykazu, po
-- której się przechodzi, którą się filtruje wedle autora i którą się zdejmuje
-- jednym poleceniem. Bez tego wiersza wyróżnienie modelu w długim dokumencie
-- ginie: postać wie, że tło jest żółte, ale nie wie, kto je nadał ani po co.
--
-- ── Dlaczego rodzaj znacznika jest osobną tabelą ────────────────────────────
-- Znacznik własny („do sprawdzenia", „wymaga źródła", „gotowe") ma nazwę, barwę
-- i wykaz. Trzymanie barwy przy każdym użyciu znaczyłoby, że zmiana barwy
-- rodzaju wymaga przejścia wszystkich znakowań — a Operator, który zmienia
-- barwę „do sprawdzenia", zmienia ją dla wszystkich, nie dla jednego miejsca.
--
-- Rodzaje są zasięgu Operatora, nie dokumentu: znacznik „wymaga źródła"
-- obowiązuje we wszystkich pismach, a nie zakłada się go od nowa w każdym.

CREATE TABLE rodzaj_znacznika_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    nazwa                    TEXT    NOT NULL UNIQUE,
    nazwa_widoczna           TEXT    NOT NULL,
    barwa                    TEXT,
    -- Rodzaju fabrycznego nie usuwa się — odpowiada odmowa nazywająca powód,
    -- wzorem `studio.operation.delete` dla operacji fabrycznych.
    fabryczny                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE znakowanie_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('highlight','mark','suggestion')),
    -- Tożsamością autora jest kod agenta z modułu Agents; rodzaj autora zostaje
    -- grubym rozróżnieniem. Znakowanie modelu jest podpisane jako model, a przy
    -- dwóch wykonawcach pracujących naraz kod agenta mówi, który z nich to
    -- postawił — bez tego obaj wyglądaliby jak jeden.
    autor_rodzaj             TEXT    NOT NULL DEFAULT 'uzytkownik'
                                     CHECK(autor_rodzaj IN ('uzytkownik','model')),
    autor_agent_kod          TEXT,
    autor_agent_nazwa        TEXT,
    autor_agent_wersja       TEXT,
    autor_podagent_kod       TEXT,
    zakres_od                INTEGER NOT NULL,
    zakres_do                INTEGER NOT NULL,
    barwa                    TEXT,
    znacznik_nazwa           TEXT REFERENCES rodzaj_znacznika_studio(nazwa) ON DELETE SET NULL,
    tresc                    TEXT,
    -- Brzmienie proponowane i brzmienie zastane stoją oba: bez zastanego nie da
    -- się pokazać, co propozycja zmienia, a po przyjęciu nie da się jej cofnąć.
    brzmienie_proponowane    TEXT,
    brzmienie_zastane        TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'open'
                                     CHECK(stan IN ('open','accepted','rejected','resolved')),
    czynnosc_kod             TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    CHECK(zakres_do >= zakres_od)
);
-- Wykaz znakowań idzie w kolejności WYSTĄPIENIA W TREŚCI, nie zapisu — Operator
-- przechodzi dokument od początku do końca, a nie od najnowszego znakowania.
CREATE INDEX idx_znakowanie_studio_dokument
    ON znakowanie_studio(dokument_id, zakres_od, id);
CREATE INDEX idx_znakowanie_studio_filtr
    ON znakowanie_studio(dokument_id, rodzaj, autor_rodzaj, stan);
CREATE INDEX idx_znakowanie_studio_agent
    ON znakowanie_studio(dokument_id, autor_agent_kod);

-- Rodzaje fabryczne, wymienione przez Właściciela wprost. Wchodzą tu, a nie
-- w kod rdzenia, żeby wykaz był jeden: gdyby rdzeń dokładał je w locie do
-- odczytu, Operator nie mógłby zmienić ich barwy, a zmiana barwy „do
-- sprawdzenia" jest normalną rzeczą.
INSERT INTO rodzaj_znacznika_studio (nazwa, nazwa_widoczna, barwa, fabryczny) VALUES
    ('doSprawdzenia', 'Do sprawdzenia', '#f5a623', 1),
    ('wymagaZrodla',  'Wymaga źródła',  '#d0021b', 1),
    ('gotowe',        'Gotowe',         '#2f9e44', 1);
