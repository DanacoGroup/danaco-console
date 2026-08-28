-- Migracja 362 — obiekty osadzone w dokumencie, aparat dokumentu i pola.
--
-- ── Dlaczego obiekt ma wiersz, a nie miejsce w drzewie postaci ───────────────
-- Obraz, kształt, ikona i pole tekstowe są bytami, po których się PYTA
-- niezależnie od tego, gdzie w treści wiszą: „które obiekty w tym dokumencie
-- przyszły z modułu Design", „co jest osadzone z bazy zdjęciowej", „pokaż
-- obiekty bez tekstu zastępczego". Drzewo postaci trzyma tylko zakotwiczenie
-- obiektu, a jego opis leży tutaj — dlatego wstawienie obiektu nie przepisuje
-- całego drzewa, a wykaz obiektów nie wymaga jego rozbierania.
--
-- Bajty obiektu leżą w magazynie zasobów pod sumą kontrolną (`zasob_kod`) —
-- tym samym, którym jedzie warsztat PDF i moduł Design. Drugiego magazynu
-- Studio nie zakłada; kolumna trzyma odwołanie, nie zawartość.
--
-- ── Dlaczego pochodzenie obiektu stoi przy obiekcie ─────────────────────────
-- Fragment i obraz wciągnięty ze strony albo z Biblioteki ma nieść zapis, skąd
-- jest — inaczej za tydzień nikt nie odtworzy, na czym pismo się opiera. Dla
-- treści służy temu osobna tabela (migracja 368); dla obiektu wystarczą trzy
-- kolumny tutaj, bo obiekt ma jedno źródło i nie dzieli się na fragmenty.
--
-- ── Dlaczego aparat dokumentu jest jedną tabelą, a nie trzynastoma ──────────
-- Spis treści, spis ilustracji, przypis dolny i końcowy, podpis, zakładka,
-- odwołanie wzajemne, odsyłacz, powołanie, bibliografia, hasło i indeks różnią
-- się tym, CO NIOSĄ, ale nie tym, JAK się nimi pracuje: każdy jest przypięty do
-- miejsca w treści, każdy ma numer nadawany przy odświeżeniu i każdy może być
-- nieświeży. Trzynaście tabel o tym samym kształcie znaczyłoby trzynaście
-- zapytań przy każdym odświeżeniu aparatu. Rodzaj rozstrzyga kolumna `rodzaj`,
-- a to, co swoiste, leży w `dane_json`.
--
-- ── Dlaczego pole ma tabelę osobną od aparatu ──────────────────────────────
-- Pole jest przeliczane, nie zbierane: numer strony, liczba stron, data
-- i pole obliczane liczą się z układu i z właściwości dokumentu, a nie
-- z nagłówków. Odświeżenie pól i odświeżenie aparatu są dwiema czynnościami
-- o różnym koszcie i różnej porze — pola odświeżają się przy każdym podglądzie,
-- spis treści wtedy, gdy Operator o to poprosi.

CREATE TABLE obiekt_dokumentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('image','shape','icon','textbox','logo','chart')),
    zrodlo                   TEXT
                                     CHECK(zrodlo IS NULL OR zrodlo IN ('file','coreAsset','designModule',
                                           'photoBank','libraryFile','web','drawn')),
    zasob_kod                TEXT,
    design_wezel_kod         TEXT,
    biblioteka_plik_kod      TEXT,
    adres_zrodla             TEXT,
    zakotwiczenie            TEXT    NOT NULL DEFAULT 'paragraph'
                                     CHECK(zakotwiczenie IN ('character','paragraph','page')),
    zakotwiczenie_pozycja    INTEGER NOT NULL DEFAULT 0,
    warstwa                  INTEGER NOT NULL DEFAULT 0,
    tekst_zastepczy          TEXT,
    tekst_wewnetrzny         TEXT,
    podpis                   TEXT,
    postac_json              TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_obiekt_dokumentu_studio_dokument
    ON obiekt_dokumentu_studio(dokument_id, zakotwiczenie_pozycja, id);
CREATE INDEX idx_obiekt_dokumentu_studio_rodzaj
    ON obiekt_dokumentu_studio(dokument_id, rodzaj);

CREATE TABLE element_aparatu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('toc','figureIndex','tableIndex','footnote','endnote',
                                           'caption','bookmark','crossReference','hyperlink','citation',
                                           'bibliography','indexEntry','index')),
    kotwica_od               INTEGER NOT NULL DEFAULT 0,
    kotwica_do               INTEGER NOT NULL DEFAULT 0,
    numer                    TEXT,
    etykieta                 TEXT,
    tresc                    TEXT,
    cel_kod                  TEXT,
    cel_adres                TEXT,
    -- Kolejność w obrębie rodzaju. Przypis wstawiony PRZED innym przypisem ma
    -- przenumerować oba — a numeracja liczy się z kolejności, nie z klucza
    -- wiersza, bo klucz rośnie w porządku zapisu, nie w porządku czytania.
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    nieswiezy                INTEGER NOT NULL DEFAULT 0,
    dane_json                TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_element_aparatu_studio_dokument
    ON element_aparatu_studio(dokument_id, rodzaj, kolejnosc, id);
CREATE INDEX idx_element_aparatu_studio_kotwica
    ON element_aparatu_studio(dokument_id, kotwica_od);
CREATE INDEX idx_element_aparatu_studio_nieswieze
    ON element_aparatu_studio(dokument_id, nieswiezy);

CREATE TABLE pole_dokumentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('pageNumber','pageCount','date','time','documentTitle',
                                           'documentAuthor','documentProperty','calculated','templateField')),
    kotwica                  INTEGER NOT NULL DEFAULT 0,
    format                   TEXT,
    wyrazenie                TEXT,
    nazwa_wlasciwosci        TEXT,
    wartosc                  TEXT,
    nieswieze                INTEGER NOT NULL DEFAULT 1,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_pole_dokumentu_studio_dokument
    ON pole_dokumentu_studio(dokument_id, kotwica, id);
