-- Migracja 270 — dziennik przebiegu automatyki (`automation.execution.log`).
--
-- Log przebiegu musi przeżyć restart rdzenia. Przeglądarka logów Execution
-- Monitora pokazuje „pełny zapis zdarzeń pojedynczego uruchomienia, także
-- sprzed otwarcia okna” — a strumień WebSocket niesie wyłącznie to, co padło
-- przy otwartym oknie. Log trzymany w pamięci procesu byłby logiem, którego po
-- restarcie nie ma, choć przebieg dalej widnieje w historii.
--
-- Poziom idzie słownikiem bazy (kolumna `baza` wyliczenia `AutomationLogLevel`).
CREATE TABLE log_przebiegu (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    przebieg_id INTEGER NOT NULL REFERENCES przebieg_automatyki(id) ON DELETE CASCADE,
    krok_kod    TEXT,
    poziom      TEXT    NOT NULL DEFAULT 'informacja'
                        CHECK(poziom IN ('informacja','ostrzezenie','blad')),
    tresc       TEXT    NOT NULL,
    chwila      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- automation.execution.log czyta wiersze przebiegu w kolejności czasu,
-- zawężając po poziomie i po kroku.
CREATE INDEX idx_log_przebiegu_czas
    ON log_przebiegu(przebieg_id, id);
