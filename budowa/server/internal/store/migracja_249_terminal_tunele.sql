-- Migracja 249 — przekierowania portów modułu Terminal (okno Session Manager).
--
-- Do tej pory tunel dało się założyć wyłącznie poleceniem wydanym w karcie:
-- `ssh -L …` biegł wtedy jako zwykły proces, a jego stan i to, czy port po
-- stronie rdzenia w ogóle nasłuchuje, były niewidoczne. Tunel dostaje własny byt,
-- bo jest bytem długożyjącym o własnym stanie — inaczej niż polecenie, które ma
-- początek i koniec.
--
-- Dlaczego tunel ma wiersz, skoro po restarcie rdzenia nie biegnie. Bo Operator
-- ma zobaczyć, że tunel BYŁ i dlaczego go nie ma. Przy montażu rdzeń przestawia
-- tunele zostawione w stanie `active` na `inactive` (tak samo jak osierocone
-- procesy w migracji 041) — wykazywanie ich jako czynnych byłoby nieprawdą.
--
-- Kolumn `bajty_we` i `bajty_wy` nie ma. Kontrakt ma na nie pola opcjonalne,
-- a licznik przepustowości jest wielkością chwili, nie dziennika: zapisany
-- w bazie starzeje się między dwoma odczytami. Rdzeń liczy je przy żywym tunelu
-- i tam, gdzie ich nie zna, nie wpisuje zera — zero znaczyłoby „nic nie
-- przeszło”, a nie „nie wiadomo”.

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
    -- Powód niepowodzenia; pusty poza stanem `failed`.
    powod         TEXT    NOT NULL DEFAULT '',
    zalozono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zamknieto     TEXT
);
CREATE INDEX idx_terminal_tunel_okno ON terminal_tunel(okno_kod, zalozono DESC);
CREATE INDEX idx_terminal_tunel_stan ON terminal_tunel(stan, zalozono DESC);
