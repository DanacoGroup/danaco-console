-- Migracja 200 tworzy stronę budowy produktu modułu apps: produkt okna, etapy
-- budowy, kamienie milowe i ich związek z etapami.

-- Tabela produkt_apps przechowuje jeden produkt na okno; kolumna okno jest
-- UNIQUE, a zapis jest UPSERT-em po niej.
CREATE TABLE produkt_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- Jedyny produkt okna; drugi zapis dla tego samego okna nadpisuje ten sam wiersz.
    okno                     TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    -- Wartości kontraktu (AppProductPlatform) rozdzielone znakiem nowej linii.
    platformy                TEXT,
    repozytorium             TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- ── Etap budowy produktu — tracker Product Buildera ──────────────────────────
-- Wartości kolumny `stan` są wartościami kontraktu (AppStageStatus).
CREATE TABLE etap_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    stan                     TEXT    NOT NULL DEFAULT 'pending'
                                     CHECK(stan IN ('pending','active','done','blocked')),
    -- Wykonawca etapu; pusty łańcuch żądania zdejmuje przypisanie, więc kolumna
    -- dopuszcza brak wartości.
    wykonawca                TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_etap_apps_okno ON etap_apps(okno, kolejnosc, id);

-- Tabela kamien_milowy_apps przechowuje kamień milowy produktu; termin niesie
-- milisekundy epoki wprost z kontraktu, bez przekładu.
CREATE TABLE kamien_milowy_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    termin                   INTEGER,
    stan                     TEXT    NOT NULL DEFAULT 'planned'
                                     CHECK(stan IN ('planned','active','reached','missed')),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_kamien_milowy_apps_okno ON kamien_milowy_apps(okno, termin, id);

-- Tabela kamien_milowy_etap_apps wiąże kamień milowy z etapami tabelą
-- złącznikową, wymienianą w całości przy każdym zapisie.
CREATE TABLE kamien_milowy_etap_apps (
    kamien_id INTEGER NOT NULL REFERENCES kamien_milowy_apps(id) ON DELETE CASCADE,
    -- Kod zewnętrzny etapu, nie więz obcy: kamień może wskazywać etap usunięty.
    etap_kod  TEXT    NOT NULL,
    PRIMARY KEY (kamien_id, etap_kod)
);
