-- Migracja 208 — rodzina `extension.*`, warstwa protokołu: narzędzia odkryte
-- u integracji, dziennik ramek JSON-RPC, wywołania wraz z ich czasem oraz
-- wyniki sprawdzeń kondycji.
--
-- NARZĘDZIE JEST WIERSZEM, BO ODKRYCIE MA PRZEŻYĆ ODPOWIEDŹ. `extension.tool.list`
-- ma pole `refresh`: bez niego oddaje to, co odkryto wcześniej, z nim — pyta
-- serwer na nowo. Wykaz trzymany wyłącznie w pamięci znikałby przy restarcie
-- rdzenia i pierwsze otwarcie konsoli po starcie zawsze musiałoby czekać na
-- handshake, także wtedy, gdy serwer akurat nie odpowiada.
--
-- TRZY WYKAZY PROTOKOŁU MCP W JEDNEJ TABELI. `ExtensionToolKind` niesie `tool`,
-- `resource` i `prompt` — trzy wykazy `tools/list`, `resources/list`,
-- `prompts/list` — a kształt wpisu jest dla nich wspólny (nazwa, opis, schemat
-- wejścia, adres). Trzy tabele o tych samych kolumnach byłyby trzema prawdami
-- o jednym bycie.
--
-- RAMKA PROTOKOŁU JEST DZIENNIKIEM DIAGNOSTYCZNYM, NIE STANEM POŁĄCZENIA.
-- `extension.protocol.log.list` czyta ją przy diagnozie błędu integracji; wiersz
-- powstaje przy każdym wywołaniu narzędzia i przy każdym powitaniu, po jednym
-- na kierunek, z korelacją żądania z odpowiedzią.
--
-- WYWOŁANIE MA WIERSZ, BO METRYKA UŻYCIA MUSI SIĘ Z CZEGOŚ LICZYĆ.
-- `extension.usage.get` oddaje liczbę wywołań, liczbę niepowodzeń i średni czas
-- odpowiedzi w oknie czasu — żadnej z tych trzech wartości nie da się policzyć
-- z licznika nadpisywanego w miejscu. `extension.audit.list` czyta te same
-- wiersze od strony eksperta, który wywołania dokonał.
--
-- KONDYCJA JEST DZIENNIKIEM SPRAWDZEŃ. `extension.health.check` dopisuje wynik
-- każdego sprawdzenia; panel zdrowia Integrations Hub pokazuje najnowszy, a
-- historia zostaje, bo z niej widać, kiedy integracja zaczęła się psuć.

-- ── Narzędzie, zasób albo prompt odkryty u integracji ────────────────────────
CREATE TABLE narzedzie_rozszerzenia (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    rozszerzenie_kod TEXT    NOT NULL,
    -- Wartości kontraktu (ExtensionToolKind) wprost, bez tłumaczenia.
    rodzaj           TEXT    NOT NULL CHECK(rodzaj IN ('tool','resource','prompt')),
    nazwa            TEXT    NOT NULL,
    opis             TEXT,
    -- Surowy JSON Schema wejścia narzędzia; rdzeń go nie rozkłada, bo kształt
    -- należy do serwera MCP, nie do platformy.
    schemat_wejscia  TEXT,
    adres            TEXT,
    odkryto          INTEGER NOT NULL DEFAULT 0,
    -- Wersja protokołu podana w powitaniu; ta sama dla całego wykazu jednego
    -- odkrycia, ale trzymana przy wpisie, bo wykaz oddaje się w częściach.
    wersja_protokolu TEXT,
    UNIQUE (rozszerzenie_kod, rodzaj, nazwa)
);
CREATE INDEX idx_narzedzie_rozszerzenia ON narzedzie_rozszerzenia(rozszerzenie_kod, rodzaj, nazwa);

-- ── Ramka protokołu JSON-RPC ─────────────────────────────────────────────────
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

-- ── Wywołanie narzędzia integracji ───────────────────────────────────────────
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

-- ── Wynik sprawdzenia kondycji integracji ────────────────────────────────────
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
