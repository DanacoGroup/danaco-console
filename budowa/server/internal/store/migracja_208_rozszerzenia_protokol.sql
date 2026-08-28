-- Migracja 208 zakłada tabele warstwy protokołu rodziny extension: narzędzia odkryte u integracji, dziennik ramek JSON-RPC, wywołania i wyniki sprawdzeń kondycji.

-- Zakłada tabelę narzedzie_rozszerzenia niosącą wykaz narzędzi, zasobów i promptów odkrytych u integracji protokołem MCP.
CREATE TABLE narzedzie_rozszerzenia (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    rozszerzenie_kod TEXT    NOT NULL,
    -- Wartości kontraktu (ExtensionToolKind) wprost, bez tłumaczenia.
    rodzaj           TEXT    NOT NULL CHECK(rodzaj IN ('tool','resource','prompt')),
    nazwa            TEXT    NOT NULL,
    opis             TEXT,
    -- Surowy JSON Schema wejścia narzędzia, nierozkładany, bo kształt należy do serwera MCP.
    schemat_wejscia  TEXT,
    adres            TEXT,
    odkryto          INTEGER NOT NULL DEFAULT 0,
    -- Wersja protokołu z powitania, wspólna wykazowi jednego odkrycia, trzymana przy każdym wpisie.
    wersja_protokolu TEXT,
    UNIQUE (rozszerzenie_kod, rodzaj, nazwa)
);
CREATE INDEX idx_narzedzie_rozszerzenia ON narzedzie_rozszerzenia(rozszerzenie_kod, rodzaj, nazwa);

-- Zakłada tabelę ramka_protokolu_rozszerzenia niosącą dziennik diagnostyczny ramek JSON-RPC wymienianych z integracją.
CREATE TABLE ramka_protokolu_rozszerzenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rozszerzenie_kod         TEXT    NOT NULL,
    -- Wartości kontraktu (ProtocolFrameDirection).
    kierunek                 TEXT    NOT NULL CHECK(kierunek IN ('outgoing','incoming')),
    metoda                   TEXT,
    korelacja                TEXT,
    tresc                    TEXT,
    kod_bledu                TEXT,
    zaszlo                   INTEGER NOT NULL
);
CREATE INDEX idx_ramka_protokolu_rozszerzenia
    ON ramka_protokolu_rozszerzenia(rozszerzenie_kod, zaszlo DESC, id DESC);

-- Zakłada tabelę wywolanie_rozszerzenia niosącą wiersz każdego wywołania narzędzia integracji jako podstawę metryki użycia.
CREATE TABLE wywolanie_rozszerzenia (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    rozszerzenie_kod TEXT    NOT NULL,
    -- Ekspert, który wywołania dokonał; puste dla wywołania Operatora
    -- z konsoli inspektora.
    agent_kod        TEXT,
    narzedzie        TEXT,
    udane            INTEGER NOT NULL DEFAULT 1 CHECK(udane IN (0,1)),
    czas_ms          INTEGER NOT NULL DEFAULT 0,
    szczegol         TEXT,
    zaszlo           INTEGER NOT NULL
);
CREATE INDEX idx_wywolanie_rozszerzenia ON wywolanie_rozszerzenia(rozszerzenie_kod, zaszlo DESC);
CREATE INDEX idx_wywolanie_rozszerzenia_agent ON wywolanie_rozszerzenia(agent_kod, zaszlo DESC);

-- Zakłada tabelę kondycja_rozszerzenia niosącą dziennik wyników sprawdzeń kondycji integracji w czasie.
CREATE TABLE kondycja_rozszerzenia (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    rozszerzenie_kod TEXT    NOT NULL,
    -- Wartości kontraktu (ExtensionHealthStatus).
    stan             TEXT    NOT NULL
                             CHECK(stan IN ('unknown','healthy','degraded','unavailable')),
    czas_ms          INTEGER,
    liczba_narzedzi  INTEGER,
    blad_powitania   TEXT,
    sprawdzono       INTEGER NOT NULL
);
CREATE INDEX idx_kondycja_rozszerzenia ON kondycja_rozszerzenia(rozszerzenie_kod, sprawdzono DESC);
