-- Migracja 210 zakłada tabele warstwy zaufania rodziny extension: uprawnienia deklarowane i nadane, podpis pozycji oraz referencje sekretów z zakresem współdzielenia.

-- Zakłada tabelę uprawnienie_rozszerzenia niosącą uprawnienie deklarowane w manifeście albo nadane przez Operatora jako dwa fakty jednej kolumny.
CREATE TABLE uprawnienie_rozszerzenia (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    rozszerzenie_kod TEXT    NOT NULL,
    -- Wartości kontraktu (ExtensionPermissionScope).
    zakres           TEXT    NOT NULL
                             CHECK(zakres IN ('network','fileRead','fileWrite',
                                              'processSpawn','secretRead','modelCall')),
    -- Byt objęty zakresem: domena, korzeń katalogu, nazwa zasobu; pustka znaczy brak zawężenia.
    byt              TEXT,
    -- Wartości kontraktu (AccessMode).
    tryb             TEXT CHECK(tryb IN ('read','write')),
    objasnienie      TEXT,
    -- 0 — uprawnienie deklarowane w manifescie, 1 — nadane przez Operatora.
    nadane           INTEGER NOT NULL DEFAULT 0 CHECK(nadane IN (0,1)),
    -- Ekspert, któremu uprawnienie nadano; puste znaczy „całej platformie".
    agent_kod        TEXT,
    nadano           INTEGER,
    UNIQUE (rozszerzenie_kod, zakres, byt, nadane, agent_kod)
);
CREATE INDEX idx_uprawnienie_rozszerzenia ON uprawnienie_rozszerzenia(rozszerzenie_kod, nadane);

-- Zakłada tabelę podpis_rozszerzenia niosącą wynik weryfikacji podpisu pozycji wraz z bajtami podpisu i kluczem publicznym.
CREATE TABLE podpis_rozszerzenia (
    rozszerzenie_kod TEXT PRIMARY KEY,
    algorytm         TEXT,
    suma_kontrolna   TEXT,
    wydawca          TEXT,
    -- Bajty podpisu i klucz publiczny w postaci base64.
    podpis_base64    TEXT,
    klucz_base64     TEXT,
    zaktualizowano   INTEGER NOT NULL DEFAULT 0
);

-- Zakłada tabele referencji sekretu i zakresu jej współdzielenia z pozycją katalogu albo rolą, bez treści poświadczenia.
CREATE TABLE sekret_rozszerzenia (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Klucz jawny odwołania — `SecretRef.ref`.
    odwolanie        TEXT    NOT NULL UNIQUE,
    etykieta         TEXT,
    -- Wartości kontraktu (ExtensionAuthKind).
    sposob_logowania TEXT CHECK(sposob_logowania IN ('oauth2','apiKey','token','basic','none')),
    wygasa           INTEGER,
    zaktualizowano   INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE udostepnienie_sekretu_rozszerzenia (
    sekret_id INTEGER NOT NULL REFERENCES sekret_rozszerzenia(id) ON DELETE CASCADE,
    -- Rodzaj bytu dopuszczonego do odwołania: pozycja katalogu albo rola.
    rodzaj    TEXT    NOT NULL CHECK(rodzaj IN ('rozszerzenie','rola')),
    byt_kod   TEXT    NOT NULL,
    PRIMARY KEY (sekret_id, rodzaj, byt_kod)
);
