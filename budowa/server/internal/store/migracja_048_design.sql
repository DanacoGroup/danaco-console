-- Migracja 048 — trwałość modułu Design: prompt strukturalny, zasób
-- wizualny Assets Panel oraz kompozycja Design Board wraz z jej warstwami.
--
-- Prompt dostaje własną tabelę. Kontrakt niesie `DesignPrompt.Id` jako pole
-- opcjonalne, a `DesignAsset.PromptId` odwołuje się do promptu osobnym polem —
-- dwa sygnały, że prompt bywa bytem trwałym, nie tylko parametrem jednego
-- wywołania `design.asset.generate`. Warianty (`DesignAsset.VariantOfAssetId`,
-- `DesignPrompt.Variants`) i generowanie obraz-do-obrazu (`ReferenceAssetId`)
-- zakładają wprost, że wiele zasobów powstaje z tego samego promptu — gdyby
-- prompt żył jako kolumny powielone w każdym wierszu zasobu, każdy wariant
-- niósłby własną kopię tych samych pól. Jedna tabela `prompt_design` jest jedną
-- prawdą o promptcie; `zasob_design.prompt_id` jest jedynym miejscem odwołania.
--
-- Warstwa kompozycji dostaje własną tabelę, nie zapis strukturalny w kolumnie.
-- `DesignBoardLayer` niesie własny `Id`, pozycję (X, Y, Width, Height), `Order`
-- i `Locked` — pola, po których trzeba by filtrować i sortować przy odczycie
-- („warstwa zablokowana”, kolejność renderowania), gdyby leżały w jednym polu
-- JSON. `design.board.update` nadsyła całą listę warstw na nowo (kontrakt:
-- `Layers []DesignBoardLayer` bez trybu częściowej zmiany), więc zapis jest
-- zawsze „usuń warstwy kompozycji, wstaw przysłane od nowa”.
--
-- Treść zasobu trzyma dysk lub usługa zewnętrzna, nie baza. `DesignAsset.Uri`
-- w kontrakcie już jest odnośnikiem, nie surową treścią — kolumna `uri`
-- przechowuje więc ten odnośnik wprost, bez pośredniej kolumny BLOB.
--
-- Etykiety zasobu mają własną tabelę złącznikową: etykieta jest wolnym tekstem
-- bez własnej tożsamości (kontrakt: `DesignAsset.Tags []string`), więc para
-- (zasób, etykieta) jest kluczem bez surogatu.

-- ── Prompt strukturalny — Prompt Builder ──────────────────────────────────────
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
    -- Odwołanie do wariantu jest wartością danych (identyfikator zewnętrzny),
    -- nie więzem obcym: wariant i zasób źródłowy współistnieją bez porządku
    -- wstawiania wymuszonego przez SQLite.
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

-- ── Etykieta zasobu — Assets Panel ─────────────────────────────────────────────
CREATE TABLE etykieta_zasobu_design (
    zasob_id   INTEGER NOT NULL REFERENCES zasob_design(id) ON DELETE CASCADE,
    etykieta   TEXT    NOT NULL,
    utworzono  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (zasob_id, etykieta)
);
-- design.asset.list filtruje po etykiecie (Tags []string w żądaniu).
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

-- ── Warstwa kompozycji — Design Board ────────────────────────────────────────
CREATE TABLE warstwa_kompozycji_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    -- Zasób warstwy jest wartością danych, nie więzem obcym: warstwa może
    -- wskazywać zasób usunięty z Assets Panel po zapisie kompozycji, a
    -- design.board.update nie ma trybu naprawy takiego wskazania.
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
