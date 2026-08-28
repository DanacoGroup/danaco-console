-- Punkt dostępu opisuje maszynę osiąganą przez most MCP albo katalog lokalny
-- urządzenia, odrębnie od środowiska nawigacji i katalogu roboczego modelu.
-- Kolumna kod odpowiada polu id struktury AccessPoint kontraktu.
CREATE TABLE punkt_dostepu (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                    TEXT    NOT NULL UNIQUE,
    nazwa                  TEXT    NOT NULL,
    opis                   TEXT    NOT NULL DEFAULT '',
    rodzaj                 TEXT    NOT NULL
                                   CHECK(rodzaj IN ('mcpBridge','localDirectory')),
    -- Katalog lokalny wymaga urządzenia — istnieje na jednej maszynie; most
    -- MCP go nie wymaga.
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

-- Korzenie punktu odpowiadają zmiennej DANACO_MOST_KORZENIE ograniczającej
-- ścieżki mostu. Lista, nie pojedyncze pole, tak samo jak katalog roboczy
-- okna; nadanie może te ścieżki jedynie zawęzić.
CREATE TABLE korzen_punktu_dostepu (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    punkt_dostepu_id       INTEGER NOT NULL REFERENCES punkt_dostepu(id) ON DELETE CASCADE,
    sciezka                TEXT    NOT NULL,
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    UNIQUE(punkt_dostepu_id, sciezka)
);
CREATE INDEX idx_korzen_punktu ON korzen_punktu_dostepu(punkt_dostepu_id, kolejnosc);

-- Wiersz podaje słowo, którym dany most nazywa tryb kontraktu przy
-- uruchomieniu. Brak wiersza nie jest usterką: uruchamiający poda tryb
-- domyślny mostu, a most bez argumentu wchodzi w tryb odczytu.
CREATE TABLE argument_trybu_mostu (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    punkt_dostepu_id       INTEGER NOT NULL REFERENCES punkt_dostepu(id) ON DELETE CASCADE,
    tryb                   TEXT    NOT NULL CHECK(tryb IN ('read','write')),
    argument               TEXT    NOT NULL,
    UNIQUE(punkt_dostepu_id, tryb)
);

-- Nadanie żyje przy oknie komunikacji, nie przy sesji ani przy platformie.
-- Kolejność ustala pierwszeństwo prezentacji i przekazania mostom, a pole
-- główne wskazuje punkt domyślny okna.
CREATE TABLE nadanie_dostepu (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_komunikacji_id    INTEGER NOT NULL REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    punkt_dostepu_id       INTEGER NOT NULL REFERENCES punkt_dostepu(id) ON DELETE CASCADE,
    tryb                   TEXT    NOT NULL DEFAULT 'read' CHECK(tryb IN ('read','write')),
    kolejnosc              INTEGER NOT NULL DEFAULT 1,
    glowne                 INTEGER NOT NULL DEFAULT 0 CHECK(glowne IN (0,1)),
    aktywne                INTEGER NOT NULL DEFAULT 1 CHECK(aktywne IN (0,1)),
    -- Wartość pusta oznacza nadanie założone wprost w bazie danych, bez
    -- odpowiednika w pamięci rdzenia.
    identyfikator_zewnetrzny TEXT,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(okno_komunikacji_id, punkt_dostepu_id)
);
CREATE INDEX idx_nadanie_okno ON nadanie_dostepu(okno_komunikacji_id, kolejnosc, id);
CREATE INDEX idx_nadanie_punkt ON nadanie_dostepu(punkt_dostepu_id);
-- Głównych nadań okna jest najwyżej jedno; unikalność pilnowana jest przez
-- bazę danych, nie przez warstwę wyższą aplikacji.
CREATE UNIQUE INDEX idx_nadanie_glowne ON nadanie_dostepu(okno_komunikacji_id)
    WHERE glowne = 1;
CREATE UNIQUE INDEX idx_nadanie_identyfikator_zewnetrzny
    ON nadanie_dostepu(identyfikator_zewnetrzny)
    WHERE identyfikator_zewnetrzny IS NOT NULL;

-- Brak wierszy oznacza pełny zbiór korzeni punktu, zgodnie z opisem pola
-- korzeni struktury AccessGrant kontraktu.
CREATE TABLE korzen_nadania (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    nadanie_dostepu_id     INTEGER NOT NULL REFERENCES nadanie_dostepu(id) ON DELETE CASCADE,
    sciezka                TEXT    NOT NULL,
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    UNIQUE(nadanie_dostepu_id, sciezka)
);
CREATE INDEX idx_korzen_nadania ON korzen_nadania(nadanie_dostepu_id, kolejnosc);

-- Wstawienia kończą się klauzulą braku działania przy konflikcie, więc
-- migracja przechodzi także na bazie z częścią wierszy już obecną. Tryb
-- domyślny maszyny danych to odczyt, bo działa tam produkcyjna platforma.
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

-- Korzenie obejmują katalog instalacyjny Danaco na każdej z maszyn oraz
-- dodatkowo katalog kopii zapasowych na maszynie danych, gdzie działa
-- produkcyjna platforma.
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
