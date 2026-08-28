-- Dziennik zdarzeń zaczepów.
--
-- Konfiguracja zaczepów idzie do CLI sekcją `hooks` pliku ustawień, a przy
-- przełączniku `--include-hook-events` strumień niesie na tym samym torze
-- koperty:
--   {"type":"system","subtype":"hook_started","hook_id","hook_name",
--    "hook_event","uuid","session_id"}
--   {"type":"system","subtype":"hook_response", … ,"output","stdout",
--    "stderr","exit_code","outcome"}
--
-- Kolumna `ladunek` niesie surową kopertę ze strumienia — zapis źródłowy,
-- nie przekład; bez niego nie da się odtworzyć, na jakiej podstawie rdzeń
-- uznał zaczep za wykonany.
--
-- Wiersz niczym nie steruje: zaczep jest narzędziem, dziennik śladem po nim.
-- Odmowy zaczepów idą osobno do diagnostyki — tam jako fakty, tu jako zapis
-- źródłowy tego samego zdarzenia.
--
-- `okno_kod` i `wiadomosc_kod` są identyfikatorami kontraktowymi, nie kluczami
-- obcymi: ślad ma przeżyć byt, którego dotyczy.

CREATE TABLE dziennik_zdarzen (
    id            INTEGER PRIMARY KEY,
    chwila        INTEGER NOT NULL,                -- epoka w milisekundach
    okno_kod      TEXT    NOT NULL DEFAULT '',
    wiadomosc_kod TEXT    NOT NULL DEFAULT '',
    -- rodzaj po polsku, jak wartości `wiadomosc.rodzaj_tresci`; słownik
    -- zamknięty na to, co rdzeń dziś naprawdę zapisuje — poszerzenie słownika
    -- jest osobną migracją, nie cichym INSERT-em.
    rodzaj        TEXT    NOT NULL
                  CHECK (rodzaj IN ('zaczep_start','zaczep_odpowiedz')),
    zaczep_id     TEXT    NOT NULL DEFAULT '',      -- `hook_id`: wiąże start z odpowiedzią
    zaczep        TEXT    NOT NULL DEFAULT '',      -- `hook_name`
    zdarzenie     TEXT    NOT NULL DEFAULT '',      -- `hook_event`, np. UserPromptSubmit
    wynik         TEXT    NOT NULL DEFAULT '',      -- `outcome`; puste dla startu
    kod_wyjscia   INTEGER,                          -- `exit_code`; NULL dla startu
    tresc         TEXT    NOT NULL DEFAULT '',      -- `output` + `stderr` po przycięciu
    ladunek       TEXT    NOT NULL,                 -- surowa koperta JSON — dowód pierwotny
    sesja_cli     TEXT    NOT NULL DEFAULT ''
);

CREATE INDEX idx_dziennik_zdarzen_okno   ON dziennik_zdarzen (okno_kod, chwila);
CREATE INDEX idx_dziennik_zdarzen_zaczep ON dziennik_zdarzen (zaczep_id);
