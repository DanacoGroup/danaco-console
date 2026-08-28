-- Migracja 176 — menedżer pobrań modułu Browser (`browser.download.*`).
--
-- Pobranie jest bytem o własnym cyklu życia: czeka w kolejce, biegnie, bywa
-- wstrzymane, kończy się powodzeniem, błędem albo przerwaniem. Postęp
-- (`odebrano_bajtow` wobec `razem_bajtow`) jest liczbą mierzoną w trakcie, nie
-- opisem — wykaz pobrań ma pokazywać, ile naprawdę leży na dysku.
--
-- Komunikat błędu stoi w kolumnie obok stanu, bo „nie udało się" bez powodu
-- każe Operatorowi zgadywać, czy ponowienie ma sens. Ta sama zasada rządzi
-- odmowami rdzenia.
CREATE TABLE pobranie_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    nazwa_pliku              TEXT,
    sciezka_docelowa         TEXT,
    typ_mime                 TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'queued',
    odebrano_bajtow          INTEGER,
    razem_bajtow             INTEGER,
    komunikat_bledu          TEXT,
    rozpoczeto               TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono               TEXT
);
CREATE INDEX idx_pobranie_przegladania_okno ON pobranie_przegladania(okno, rozpoczeto DESC, id);
