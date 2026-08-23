-- Migracja 209 — rodzina `extension.*`, warstwa integracji zewnętrznych:
-- transport i poświadczenie integracji, webhooki oraz odwzorowania danych.
--
-- TRANSPORT JEST WŁASNOŚCIĄ INTEGRACJI, NIE PUNKTU DOSTĘPU.
-- `extension.transport.set` przyjmuje transport wraz z adresem albo poleceniem
-- procesu, a `punkt_dostepu` (migracja 013) opisuje most maszyny, nie sposób
-- rozmowy z serwerem MCP tej pozycji. Wiersz niżej trzyma to, czym rdzeń woła
-- serwer: wybrany transport, adres nasłuchu i polecenie procesu lokalnego.
--
-- POŚWIADCZENIE TO ODWOŁANIE, NIGDY TREŚĆ. `extension.credential.bind` przyjmuje
-- `credentialRef` — klucz jawny warstwy sekretów — i zakres uprawnień OAuth2.
-- Hasła, tokenu ani klucza API w tej tabeli nie ma i nie będzie; wprowadzenie
-- ich zostaje po stronie Operatora, a rdzeń trzyma wyłącznie nazwę, po której
-- sejf je wydaje.
--
-- WEBHOOK MA JEDNĄ TABELĘ NA OBA KIERUNKI. Przychodzący niesie adres nasłuchu
-- i odwołanie do sekretu podpisu HMAC, wychodzący — adres docelowy i wykaz
-- zdarzeń platformy. Kształt jest wspólny, więc dwie tabele byłyby dwiema
-- prawdami; kolumna `kierunek` rozstrzyga, które pola mają znaczenie.
--
-- ODWZOROWANIE TRZYMA REGUŁY JAKO SUROWY JSON. Kontrakt niesie je polem `rules`
-- typu `json`, a kształt reguły należy do systemu zewnętrznego. Rozłożenie ich
-- na kolumny byłoby wyborem cudzego kształtu zrobionym przez platformę.

-- ── Transport i poświadczenie integracji ─────────────────────────────────────
-- Jeden wiersz na pozycję: `extension.transport.set` i `extension.credential.bind`
-- zmieniają dwie strony tej samej rozmowy z jednym serwerem, więc zapis jest
-- UPSERT-em po kodzie pozycji, nie dziennikiem kolejnych nastaw.
CREATE TABLE integracja_rozszerzenia (
    rozszerzenie_kod  TEXT PRIMARY KEY,
    -- Wartości kontraktu (McpTransport).
    transport         TEXT CHECK(transport IN ('stdio','sse','http','streamableHttp')),
    adres             TEXT,
    polecenie         TEXT,
    -- Wartości kontraktu (ExtensionAuthKind).
    sposob_logowania  TEXT CHECK(sposob_logowania IN ('oauth2','apiKey','token','basic','none')),
    -- Klucz jawny warstwy sekretów — patrz czoło pliku.
    odwolanie_sekretu TEXT,
    -- Zakresy OAuth2 rozdzielone znakiem nowej linii; nikt nie filtruje po
    -- pojedynczym zakresie, a wychodzą zawsze w komplecie.
    zakresy           TEXT,
    zaktualizowano    INTEGER NOT NULL DEFAULT 0
);

-- ── Webhook przychodzący albo wychodzący ─────────────────────────────────────
CREATE TABLE webhook_rozszerzenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rozszerzenie_kod         TEXT    NOT NULL,
    -- Wartości kontraktu (ExtensionWebhookDirection).
    kierunek                 TEXT    NOT NULL CHECK(kierunek IN ('inbound','outbound')),
    adres                    TEXT,
    adres_nasluchu           TEXT,
    -- Zdarzenia platformy rozdzielone znakiem nowej linii.
    zdarzenia                TEXT,
    odwolanie_sekretu        TEXT,
    czynny                   INTEGER NOT NULL DEFAULT 0 CHECK(czynny IN (0,1)),
    zaktualizowano           INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_webhook_rozszerzenia ON webhook_rozszerzenia(rozszerzenie_kod, kierunek, id);

-- ── Odwzorowanie i transformacja danych ──────────────────────────────────────
CREATE TABLE mapowanie_rozszerzenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rozszerzenie_kod         TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    -- Surowy JSON reguł — patrz czoło pliku.
    reguly                   TEXT    NOT NULL DEFAULT '{}',
    zaktualizowano           INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_mapowanie_rozszerzenia ON mapowanie_rozszerzenia(rozszerzenie_kod, nazwa);
