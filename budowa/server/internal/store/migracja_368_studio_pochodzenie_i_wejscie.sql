-- Migracja 368 — pochodzenie fragmentów, autozamiana znaków i malarz formatów.
--
-- ── Dlaczego pochodzenie jest bytem trwałym ─────────────────────────────────
-- Fragment wciągnięty ze strony albo z Biblioteki niesie zapis, skąd jest —
-- adres, plik, wersja. Bez tego za tydzień nikt nie odtworzy, na czym pismo się
-- opiera. To jest też podstawa pod „podobieństwa" w panelu redaktora i pod
-- bibliografię: powołanie bibliograficzne dopisywane z ręki po tygodniu jest
-- zgadywaniem, a nie powołaniem.
--
-- Zakres jest liczony w znakach i przesuwa się razem z treścią, tak samo jak
-- zakres blokady. Fragment usunięty zostawia wiersz o zerowej długości —
-- świadomie: to, że Operator wniósł kiedyś fragment z danego źródła, jest
-- faktem, którego usunięcie tekstu nie unieważnia, a wiersz zdejmuje się
-- jawnie, nie w tle.
--
-- ── Dlaczego autozamiana jest zasięgu Operatora, nie dokumentu ──────────────
-- Skrót zamieniany na znak („--" na półpauzę, „(c)" na znak praw autorskich)
-- obowiązuje Operatora we wszystkich pismach. Zakładanie go od nowa w każdym
-- dokumencie byłoby pracą bez powodu. Nastawa jest JAWNA i ODWRACALNA — stąd
-- kolumna `czynna`, a nie usuwanie wiersza przy wyłączeniu: Operator, który
-- wyłączył zasadę fabryczną, ma ją zobaczyć wyłączoną, a nie stracić.
--
-- ── Dlaczego malarz formatów ma wiersz, a nie pamięć procesu ────────────────
-- Malarz kopiuje POSTAĆ, nie treść, i nanosi ją w innym miejscu — czyli między
-- pobraniem a naniesieniem stoją dwie osobne komendy. Postać trzymana w pamięci
-- procesu przepadłaby przy przeładowaniu rdzenia, a przy pracy modelu przepadła
-- by jeszcze łatwiej: model pobiera postać jednym narzędziem i nanosi ją drugim,
-- być może po kilku innych czynnościach. Wiersz z wygasaniem jest jedynym
-- miejscem, w którym oba wywołania widzą to samo.

CREATE TABLE pochodzenie_fragmentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('web','libraryFile','research','clipboard',
                                           'template','importedFile')),
    zakres_od                INTEGER NOT NULL,
    zakres_do                INTEGER NOT NULL,
    adres_zrodla             TEXT,
    biblioteka_plik_kod      TEXT,
    wersja_zrodla            TEXT,
    tytul_zrodla             TEXT,
    siegnieto                TEXT,
    autor_rodzaj             TEXT    NOT NULL DEFAULT 'uzytkownik'
                                     CHECK(autor_rodzaj IN ('uzytkownik','model')),
    autor_agent_kod          TEXT,
    autor_agent_nazwa        TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    CHECK(zakres_do >= zakres_od)
);
CREATE INDEX idx_pochodzenie_fragmentu_studio_dokument
    ON pochodzenie_fragmentu_studio(dokument_id, zakres_od, id);
CREATE INDEX idx_pochodzenie_fragmentu_studio_rodzaj
    ON pochodzenie_fragmentu_studio(dokument_id, rodzaj);

CREATE TABLE autozamiana_znaku_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    skrot                    TEXT    NOT NULL UNIQUE,
    zamiennik                TEXT    NOT NULL,
    czynna                   INTEGER NOT NULL DEFAULT 1,
    fabryczna                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Zasady fabryczne: znaki interpunkcyjne niedostępne z klawiatury, wymienione
-- przez Właściciela. Wchodzą do bazy, a nie do kodu, żeby Operator mógł je
-- wyłączyć — „--" zamieniane na półpauzę przeszkadza temu, kto pisze o wierszu
-- polecenia.
INSERT INTO autozamiana_znaku_studio (skrot, zamiennik, fabryczna) VALUES
    ('--',   '–', 1),
    ('---',  '—', 1),
    ('(c)',  '©', 1),
    ('(r)',  '®', 1),
    ('(tm)', '™', 1),
    ('...',  '…', 1),
    ('->',   '→', 1),
    ('<-',   '←', 1),
    ('+-',   '±', 1),
    ('!=',   '≠', 1),
    ('<=',   '≤', 1),
    ('>=',   '≥', 1);

CREATE TABLE postac_malarza_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    dokument_kod             TEXT,
    postac_znaku_json        TEXT,
    postac_akapitu_json      TEXT,
    styl_nazwany             TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Pobrana postać wygasa: malarz jest narzędziem jednej czynności, a postać
    -- pobrana wczoraj naniesiona dziś byłaby zaskoczeniem, nie pomocą.
    wygasa                   TEXT
);
CREATE INDEX idx_postac_malarza_studio_okno
    ON postac_malarza_studio(okno, utworzono DESC, id DESC);
