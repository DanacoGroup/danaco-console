-- Migracja tworzy trwałość modułu Design: prompt strukturalny ma własną tabelę,
-- bo wiele zasobów i wariantów powstaje z tego samego promptu, a warstwy kompozycji
-- mają własną tabelę zamiast zapisu JSON.

-- Tabela prompt_design trzyma prompt strukturalny modułu Prompt Builder jako byt
-- trwały, niezależny od pojedynczego wywołania generowania zasobu.
CREATE TABLE prompt_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    temat                    TEXT    NOT NULL,
    styl                     TEXT,
    kompozycja               TEXT,
    oswietlenie              TEXT,
    paleta                   TEXT,
    proporcje_kadru          TEXT,
    wykluczenia              TEXT,
    ziarno                   INTEGER,
    warianty                 INTEGER,
    silnik                   TEXT,
    kreatywnosc              REAL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- ── Zasób wizualny — Assets Panel ──────────────────────────────────────────────
-- Wartości kolumny `rodzaj` są wartościami kontraktu (DesignAssetKind), nie
-- ich tłumaczeniem.
CREATE TABLE zasob_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('image','vector','composition')),
    format                   TEXT,
    uri                      TEXT,
    prompt_id                INTEGER REFERENCES prompt_design(id) ON DELETE SET NULL,
    -- Wariant zasobu jest polem danych, nie więzem obcym: nie wymusza kolejności wstawiania wierszy.
    wariant_zasobu_id        TEXT,
    ulubiony                 INTEGER NOT NULL DEFAULT 0 CHECK(ulubiony IN (0,1)),
    szerokosc                INTEGER,
    wysokosc                 INTEGER,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- design.asset.list filtruje po oknie, rodzaju i ulubionych, sortując od
-- najnowszych (kontrakt nie daje jawnego pola sortowania — CreatedAt jest
-- jedynym polem czasu na DesignAsset).
CREATE INDEX idx_zasob_design_okno ON zasob_design(okno, utworzono DESC, id DESC);
CREATE INDEX idx_zasob_design_prompt ON zasob_design(prompt_id);

-- Tabela etykieta_zasobu_design trzyma etykiety zasobu panelu Assets Panel jako
-- parę zasób-etykieta bez własnego identyfikatora.
CREATE TABLE etykieta_zasobu_design (
    zasob_id   INTEGER NOT NULL REFERENCES zasob_design(id) ON DELETE CASCADE,
    etykieta   TEXT    NOT NULL,
    utworzono  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (zasob_id, etykieta)
);
-- Indeks etykiety wspiera zapytanie design.asset.list, które filtruje wykaz zasobów po etykiecie żądania.
CREATE INDEX idx_etykieta_zasobu_design_etykieta ON etykieta_zasobu_design(etykieta, zasob_id);

-- ── Kompozycja — Design Board ───────────────────────────────────────────────────
-- Kompozycja jest bytem okna: `design.board.update` kieruje przez `windowId`
-- i `boardId` opcjonalny (brak zakłada nową).
CREATE TABLE kompozycja_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_kompozycja_design_okno ON kompozycja_design(okno, zaktualizowano DESC);

-- Tabela warstwa_kompozycji_design trzyma warstwy kompozycji Design Board, każda
-- z własną pozycją, rozmiarem i kolejnością renderowania.
CREATE TABLE warstwa_kompozycji_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    -- Zasób warstwy jest polem danych, nie więzem obcym; może wskazywać zasób już usunięty z panelu.
    zasob_id                 TEXT,
    x                        REAL,
    y                        REAL,
    szerokosc                REAL,
    wysokosc                 REAL,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    zablokowana              INTEGER NOT NULL DEFAULT 0 CHECK(zablokowana IN (0,1)),
    adnotacja                TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_warstwa_kompozycji_design_kompozycja ON warstwa_kompozycji_design(kompozycja_id, kolejnosc, id);
