-- Migracja 210 — rodzina `extension.*`, warstwa zaufania: uprawnienia
-- deklarowane i nadane, podpis pozycji oraz rejestr referencji sekretów wraz
-- z zakresem współdzielenia.
--
-- UPRAWNIENIE DEKLAROWANE I NADANE TO DWA RÓŻNE FAKTY W JEDNEJ TABELI.
-- `extension.permission.list` oddaje osobno `declared` (co manifest deklaruje)
-- i `granted` (co Operator naprawdę nadał), a różnicę między nimi nazywa polem
-- `excessive`. Kolumna `nadane` rozstrzyga, którym z dwóch faktów wiersz jest;
-- dwie tabele byłyby dwiema prawdami o jednym uprawnieniu i rozjechałyby się
-- przy pierwszej zmianie manifestu.
--
-- ŻADNE UPRAWNIENIE NICZEGO NIE BLOKUJE. Warstwa rozszerzeń mówi wprost, że
-- kontrola idzie przez stan wyjściowy i zakres, nie przez bramę; ta tabela jest
-- zapisem tego, co Operator wie i co nadał, a nie strażnikiem wywołania.
--
-- PODPIS MA WIERSZ, BO WERYFIKACJA MA WYNIK. `extension.signature.verify` oddaje
-- `ExtensionSignature` — rozstrzygnięcia, nie bajty. Bajty podpisu i klucz
-- publiczny leżą w kolumnach obok, bo bez nich nie da się zweryfikować niczego
-- po raz drugi, a weryfikacja przepisująca zapamiętane „tak" nie byłaby
-- weryfikacją.
--
-- REFERENCJA SEKRETU TO KLUCZ JAWNY. Wiersz trzyma nazwę odwołania, jego
-- etykietę, sposób uwierzytelnienia i termin ważności — a nie treść
-- poświadczenia, która żyje w sejfie rdzenia. Zakres współdzielenia ma tabelę
-- złącznikową, bo `extension.secret.share` nadsyła oba wykazy (rozszerzenia
-- i role) w komplecie.

-- ── Uprawnienie pozycji katalogu ─────────────────────────────────────────────
CREATE TABLE uprawnienie_rozszerzenia (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    rozszerzenie_kod TEXT    NOT NULL,
    -- Wartości kontraktu (ExtensionPermissionScope).
    zakres           TEXT    NOT NULL
                             CHECK(zakres IN ('network','fileRead','fileWrite',
                                              'processSpawn','secretRead','modelCall')),
    -- Byt objęty zakresem: domena, korzeń katalogu, nazwa zasobu. Pusty znaczy
    -- „bez zawężenia" i to właśnie wychwytuje skaner manifestu.
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

-- ── Podpis i pochodzenie pozycji ─────────────────────────────────────────────
CREATE TABLE podpis_rozszerzenia (
    rozszerzenie_kod TEXT PRIMARY KEY,
    algorytm         TEXT,
    suma_kontrolna   TEXT,
    wydawca          TEXT,
    -- Bajty podpisu i klucz publiczny w postaci base64 — patrz czoło pliku.
    podpis_base64    TEXT,
    klucz_base64     TEXT,
    zaktualizowano   INTEGER NOT NULL DEFAULT 0
);

-- ── Referencja sekretu i zakres jej współdzielenia ───────────────────────────
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
