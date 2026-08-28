-- Migracja 013 — punkty dostępu i nadania dostępu.
--
-- Punkt dostępu nie jest środowiskiem. `srodowisko` pozostaje profilem
-- widoczności modułów w bocznej nawigacji, więc dopisanie maszyny VPS do tamtej
-- tabeli rozbiłoby nawigację platformy; ta migracja tamtej tabeli nie dotyka.
--
-- Nie jest też katalogiem roboczym modelu. Katalog roboczy to ustawienie
-- (klucze `katalog.roboczy.podstawa`, `katalog.roboczy.wzorzec_sesji`) mówiące,
-- gdzie model zostawia własne pliki. Punkt dostępu mówi, do czego model ma
-- wgląd. Model może czytać `C:\Projekty\Lex`, a pliki zostawiać w
-- `<instalacja>/sesje/<id>/` — to dwa niezależne ustawienia.
--
-- Odwzorowany most MCP to `mcp-danaco-pulpit-console`:
--   /opt/danaco/most-konsoli/uruchom-stdio.sh — tryb bierze z argumentu albo
--     ze zmiennej DANACO_MOST_TRYB; brak obu znaczy `odczyt`,
--   /opt/danaco/most-konsoli/most_pulpit.py — korzenie z DANACO_MOST_KORZENIE
--     rozdzielonych `:`; poza korzenie most nie wychodzi, a w trybie odczytu
--     narzędzia zapisu nie są ogłaszane w tools/list.
--
-- Wartości słownikowe. Kolumny `rodzaj`, `tryb_domyslny`, `tryb` i `stan` niosą
-- wartości kontraktu wprost (`mcpBridge`, `localDirectory`, `read`, `write`,
-- `unknown`, `reachable`, `unreachable`). Kontrakt nie deklaruje przy tych
-- wyliczeniach pola `baza` — w odróżnieniu od AccountKind, które je ma. Własny
-- przekład polsko-angielski w warstwie trwałości byłby drugim źródłem przekładu
-- obok `shared/contract.json`.
--
-- Słownictwo mostu. `odczyt` i `zapis` to słowa argumentu uruchamiającego most,
-- nie wartości wyliczenia platformy. Są własnością konkretnego mostu, więc leżą
-- w danych (tabela `argument_trybu_mostu`), nie w kodzie: inny most może żądać
-- innych słów, a to ma być nowy wiersz, nie nowa gałąź w kodzie.

-- ── Punkt dostępu: maszyna przez most MCP albo katalog lokalny urządzenia ─────
-- `kod` jest trwałym identyfikatorem punktu i odpowiada polu `id` struktury
-- AccessPoint kontraktu. Adres mostu (`endpoint`) nie ma własnej kolumny —
-- składa się go z `uzytkownik`, `host` i `port`, żeby wiersz nie mógł mieć
-- adresu sprzecznego z własnymi polami.
CREATE TABLE punkt_dostepu (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                    TEXT    NOT NULL UNIQUE,
    nazwa                  TEXT    NOT NULL,
    opis                   TEXT    NOT NULL DEFAULT '',
    rodzaj                 TEXT    NOT NULL
                                   CHECK(rodzaj IN ('mcpBridge','localDirectory')),
    -- Katalog lokalny istnieje tylko na jednej maszynie, więc bez wskazania
    -- urządzenia nie byłby wpisem sprawdzalnym. Most MCP jest wpisem
    -- platformy — urządzenia nie wymaga.
    urzadzenie_id          INTEGER REFERENCES urzadzenie(id) ON DELETE CASCADE,
    host                   TEXT    NOT NULL DEFAULT '',
    port                   INTEGER NOT NULL DEFAULT 22,
    uzytkownik             TEXT    NOT NULL DEFAULT '',
    sciezka_klucza         TEXT    NOT NULL DEFAULT '',
    polecenie_startu       TEXT    NOT NULL DEFAULT '',
    nazwa_mostu            TEXT    NOT NULL DEFAULT '',
    -- wyłącznie odwołanie do wpisu w magazynie sekretów, nigdy treść.
    poswiadczenie_odwolanie TEXT,
    tryb_domyslny          TEXT    NOT NULL DEFAULT 'read'
                                   CHECK(tryb_domyslny IN ('read','write')),
    stan                   TEXT    NOT NULL DEFAULT 'unknown'
                                   CHECK(stan IN ('unknown','reachable','unreachable')),
    sprawdzono             TEXT,
    aktywny                INTEGER NOT NULL DEFAULT 1 CHECK(aktywny IN (0,1)),
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    CHECK(rodzaj <> 'localDirectory' OR urzadzenie_id IS NOT NULL),
    CHECK(rodzaj <> 'mcpBridge'      OR host <> '')
);
CREATE INDEX idx_punkt_dostepu_rodzaj ON punkt_dostepu(rodzaj, kolejnosc, kod);
CREATE INDEX idx_punkt_dostepu_urzadzenie ON punkt_dostepu(urzadzenie_id);

-- ── Korzenie punktu — odpowiednik DANACO_MOST_KORZENIE ────────────────────────
-- Lista, nie pojedyncze pole, tak samo jak katalog roboczy okna.
-- Poza te ścieżki punkt nie sięga; nadanie może je jedynie zawęzić.
CREATE TABLE korzen_punktu_dostepu (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    punkt_dostepu_id       INTEGER NOT NULL REFERENCES punkt_dostepu(id) ON DELETE CASCADE,
    sciezka                TEXT    NOT NULL,
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    UNIQUE(punkt_dostepu_id, sciezka)
);
CREATE INDEX idx_korzen_punktu ON korzen_punktu_dostepu(punkt_dostepu_id, kolejnosc);

-- ── Słownictwo trybu konkretnego mostu ────────────────────────────────────────
-- Wiersz mówi, jakim słowem dany most nazywa tryb kontraktu przy uruchomieniu.
-- Brak wiersza nie jest awarią: uruchamiający poda tryb domyślny mostu, a most
-- `mcp-danaco-pulpit-console` bez argumentu i tak wchodzi w odczyt.
CREATE TABLE argument_trybu_mostu (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    punkt_dostepu_id       INTEGER NOT NULL REFERENCES punkt_dostepu(id) ON DELETE CASCADE,
    tryb                   TEXT    NOT NULL CHECK(tryb IN ('read','write')),
    argument               TEXT    NOT NULL,
    UNIQUE(punkt_dostepu_id, tryb)
);

-- ── Nadanie dostępu: okno komunikacji ↔ punkt dostępu ─────────────────────────
-- Nadanie żyje per okno komunikacji, nie per sesja i nie per platforma. Okno ma
-- zbiór nadań — stąd brak unikalności na samym oknie, a `kolejnosc` i `glowne`
-- niosą znaczenie: kolejność ustala pierwszeństwo prezentacji i przekazania
-- mostom, `glowne` wskazuje punkt domyślny okna.
CREATE TABLE nadanie_dostepu (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_komunikacji_id    INTEGER NOT NULL REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    punkt_dostepu_id       INTEGER NOT NULL REFERENCES punkt_dostepu(id) ON DELETE CASCADE,
    tryb                   TEXT    NOT NULL DEFAULT 'read' CHECK(tryb IN ('read','write')),
    kolejnosc              INTEGER NOT NULL DEFAULT 1,
    glowne                 INTEGER NOT NULL DEFAULT 0 CHECK(glowne IN (0,1)),
    aktywne                INTEGER NOT NULL DEFAULT 1 CHECK(aktywne IN (0,1)),
    -- Rdzeń posługuje się identyfikatorem tekstowym; NULL oznacza nadanie
    -- założone wprost w bazie, bez odpowiednika w pamięci rdzenia.
    identyfikator_zewnetrzny TEXT,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(okno_komunikacji_id, punkt_dostepu_id)
);
CREATE INDEX idx_nadanie_okno ON nadanie_dostepu(okno_komunikacji_id, kolejnosc, id);
CREATE INDEX idx_nadanie_punkt ON nadanie_dostepu(punkt_dostepu_id);
-- Głównych nadań okna jest najwyżej jedno — pilnuje tego baza, nie warstwa wyżej.
CREATE UNIQUE INDEX idx_nadanie_glowne ON nadanie_dostepu(okno_komunikacji_id)
    WHERE glowne = 1;
CREATE UNIQUE INDEX idx_nadanie_identyfikator_zewnetrzny
    ON nadanie_dostepu(identyfikator_zewnetrzny)
    WHERE identyfikator_zewnetrzny IS NOT NULL;

-- ── Korzenie nadania — podzbiór korzeni punktu ────────────────────────────────
-- Brak wierszy znaczy „komplet korzeni punktu", zgodnie z opisem pola `roots`
-- struktury AccessGrant kontraktu.
CREATE TABLE korzen_nadania (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    nadanie_dostepu_id     INTEGER NOT NULL REFERENCES nadanie_dostepu(id) ON DELETE CASCADE,
    sciezka                TEXT    NOT NULL,
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    UNIQUE(nadanie_dostepu_id, sciezka)
);
CREATE INDEX idx_korzen_nadania ON korzen_nadania(nadanie_dostepu_id, kolejnosc);

-- ── Zaczyn: trzy maszyny mostu `mcp-danaco-pulpit-console` ────────────────────
-- Tryb domyślny. `danaco-data` dostaje `read`, bo działa tam produkcyjna
-- platforma LEX i tryb zapisu wymaga tam wyraźnej potrzeby. Tryb domyślny jest
-- propozycją przy zakładaniu nadania i podnosi się go świadomie w oknie
-- konfiguracji.
--
-- Idempotencja. Wstawienia kończą się ON CONFLICT DO NOTHING, więc migracja
-- przechodzi także na bazie, w której część wierszy już jest. Klauzula
-- `WHERE true` przed ON CONFLICT jest wymogiem składni SQLite dla
-- INSERT ... SELECT z upsertem.
INSERT INTO punkt_dostepu
    (kod, nazwa, opis, rodzaj, host, port, uzytkownik, sciezka_klucza,
     polecenie_startu, nazwa_mostu, tryb_domyslny, kolejnosc, aktywny)
VALUES
    ('danaco-system', 'danaco-system',
     'Maszyna systemowa Danaco; most MCP mcp-danaco-pulpit-console',
     'mcpBridge', '57.128.253.74', 22, 'ubuntu', '~/.ssh/id_ed25519',
     '/opt/danaco/most-konsoli/uruchom-stdio.sh', 'mcp-danaco-pulpit-console-danaco-system',
     'write', 1, 1),
    ('danaco-data', 'danaco-data',
     'Maszyna danych; produkcyjna platforma LEX — tryb zapisu wyłącznie świadomie',
     'mcpBridge', '162.19.226.145', 22, 'ubuntu', '~/.ssh/id_ed25519',
     '/opt/danaco/most-konsoli/uruchom-stdio.sh', 'mcp-danaco-pulpit-console-danaco-data',
     'read', 2, 1),
    ('danaco-web', 'danaco-web',
     'Maszyna warstwy webowej Danaco; most MCP mcp-danaco-pulpit-console',
     'mcpBridge', '137.74.41.149', 22, 'ubuntu', '~/.ssh/id_ed25519',
     '/opt/danaco/most-konsoli/uruchom-stdio.sh', 'mcp-danaco-pulpit-console-danaco-web',
     'read', 3, 1)
ON CONFLICT(kod) DO NOTHING;

-- Korzenie: katalog instalacyjny Danaco na każdej z maszyn oraz kopie zapasowe
-- na maszynie danych.
INSERT INTO korzen_punktu_dostepu (punkt_dostepu_id, sciezka, kolejnosc)
SELECT p.id, k.sciezka, k.kolejnosc
FROM punkt_dostepu p
JOIN (
    SELECT 'danaco-system' AS kod, '/opt/danaco'          AS sciezka, 1 AS kolejnosc
    UNION ALL SELECT 'danaco-data', '/opt/danaco',          1
    UNION ALL SELECT 'danaco-data', '/opt/danaco-backups',  2
    UNION ALL SELECT 'danaco-web',  '/opt/danaco',          1
) k ON k.kod = p.kod
WHERE true
ON CONFLICT(punkt_dostepu_id, sciezka) DO NOTHING;

-- Słownictwo trybu mostu `mcp-danaco-pulpit-console` — te same dwa słowa dla
-- wszystkich trzech maszyn, bo wszystkie uruchamiają ten sam skrypt.
INSERT INTO argument_trybu_mostu (punkt_dostepu_id, tryb, argument)
SELECT p.id, s.tryb, s.argument
FROM punkt_dostepu p
JOIN (
    SELECT 'read' AS tryb, 'odczyt' AS argument
    UNION ALL SELECT 'write', 'zapis'
) s
WHERE p.rodzaj = 'mcpBridge'
ON CONFLICT(punkt_dostepu_id, tryb) DO NOTHING;
