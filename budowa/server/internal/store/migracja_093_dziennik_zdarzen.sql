-- Migracja zakłada dziennik zdarzeń zaczepów, zapisujący surową kopertę ze strumienia jako
-- dowód pierwotny wykonania zaczepu.

CREATE TABLE dziennik_zdarzen (
    id            INTEGER PRIMARY KEY,
    chwila        INTEGER NOT NULL,                -- epoka w milisekundach
    okno_kod      TEXT    NOT NULL DEFAULT '',
    wiadomosc_kod TEXT    NOT NULL DEFAULT '',
    -- Rodzaj zdarzenia zamknięty na to, co rdzeń dziś zapisuje; poszerzenie jest osobną migracją.
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
