-- Migracja 209 zakłada tabele warstwy integracji zewnętrznych rodziny extension: transport i poświadczenie, webhooki oraz odwzorowania danych.

-- Zakłada tabelę integracja_rozszerzenia niosącą transport, poświadczenie i zakresy OAuth2 jednym wierszem nadpisywanym przy każdej zmianie pozycji.
CREATE TABLE integracja_rozszerzenia (
    rozszerzenie_kod  TEXT PRIMARY KEY,
    -- Wartości kontraktu (McpTransport).
    transport         TEXT CHECK(transport IN ('stdio','sse','http','streamableHttp')),
    adres             TEXT,
    polecenie         TEXT,
    -- Wartości kontraktu (ExtensionAuthKind).
    sposob_logowania  TEXT CHECK(sposob_logowania IN ('oauth2','apiKey','token','basic','none')),
    -- Klucz jawny warstwy sekretów, nigdy treść samego sekretu.
    odwolanie_sekretu TEXT,
    -- Zakresy OAuth2 rozdzielone znakiem nowej linii, wychodzące zawsze w komplecie.
    zakresy           TEXT,
    zaktualizowano    INTEGER NOT NULL DEFAULT 0
);

-- Zakłada tabelę webhook_rozszerzenia niosącą webhook przychodzący albo wychodzący w jednym kształcie rozstrzyganym kolumną kierunku.
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

-- Zakłada tabelę mapowanie_rozszerzenia niosącą regułę odwzorowania danych jako surowy JSON systemu zewnętrznego.
CREATE TABLE mapowanie_rozszerzenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rozszerzenie_kod         TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    -- Surowy JSON reguł odwzorowania, w kształcie należącym do systemu zewnętrznego.
    reguly                   TEXT    NOT NULL DEFAULT '{}',
    zaktualizowano           INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_mapowanie_rozszerzenia ON mapowanie_rozszerzenia(rozszerzenie_kod, nazwa);
