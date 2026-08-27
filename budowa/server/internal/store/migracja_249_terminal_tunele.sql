-- Migracja 249 zakłada tabelę tuneli SSH modułu Terminal jako byt długożyjący
-- o własnym stanie, niezależny od procesu polecenia, które go zakłada.

CREATE TABLE terminal_tunel (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    kod           TEXT    NOT NULL UNIQUE,
    -- Okno terminala, do którego tunel należy; z niego bierze się zasięg izolacji.
    okno_kod      TEXT    NOT NULL,
    rodzaj        TEXT    NOT NULL CHECK(rodzaj IN ('local','remote','dynamic')),
    -- Wpis książki hostów, przez który idzie tunel; pusty znaczy cel podany wprost.
    host_kod      TEXT REFERENCES terminal_host(kod) ON DELETE SET NULL,
    -- Adres celu, gdy tunel nie idzie przez wpis książki.
    cel           TEXT    NOT NULL DEFAULT '',
    port_lokalny  INTEGER CHECK(port_lokalny IS NULL OR port_lokalny BETWEEN 1 AND 65535),
    host_docelowy TEXT    NOT NULL DEFAULT '',
    port_docelowy INTEGER CHECK(port_docelowy IS NULL OR port_docelowy BETWEEN 1 AND 65535),
    stan          TEXT    NOT NULL DEFAULT 'inactive'
                          CHECK(stan IN ('inactive','active','failed')),
    -- Powód niepowodzenia; pusty poza stanem failed.
    powod         TEXT    NOT NULL DEFAULT '',
    zalozono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zamknieto     TEXT
);
CREATE INDEX idx_terminal_tunel_okno ON terminal_tunel(okno_kod, zalozono DESC);
CREATE INDEX idx_terminal_tunel_stan ON terminal_tunel(stan, zalozono DESC);
