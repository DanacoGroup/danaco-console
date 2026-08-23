-- Migracja 231 — szablony promptu strukturalnego (Prompt Builder).
--
-- Szablon ma własną tabelę, a nie flagę w `prompt_design`. Prompt wydany jest
-- zdarzeniem historii: powstał raz, poszedł kanałem, urodził zasoby i jest
-- zapisem tego, co się stało. Szablon jest bytem żywym: Operator go nadpisuje,
-- przemianowuje i używa wielokrotnie. Wspólna tabela kazałaby odróżniać jedno
-- od drugiego kolumną, po której trzeba by filtrować każdy odczyt historii,
-- a nadpisanie szablonu przepisywałoby wiersz, na który wskazują zasoby
-- (`zasob_design.prompt_id`) — czyli zmieniałoby prompt, z którego one powstały.
--
-- Pola promptu są przepisane z `prompt_design`, bo opisują ten sam kształt
-- kontraktu (`DesignPrompt`). Odwołanie do wiersza promptu zamiast kopii pól
-- związałoby szablon z jednym wydaniem promptu i usunięcie tamtego wydania
-- zabrałoby Operatorowi szablon.

CREATE TABLE szablon_promptu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    temat                    TEXT    NOT NULL,
    styl                     TEXT,
    kompozycja               TEXT,
    oswietlenie              TEXT,
    paleta                   TEXT,
    proporcje_kadru          TEXT,
    wykluczenia              TEXT,
    ziarno                   INTEGER,
    warianty                 INTEGER,
    silnik                   TEXT,
    kreatywnosc              REAL,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- design.prompt.template.list czyta szablony okna, od ostatnio zmienianego.
CREATE INDEX idx_szablon_promptu_design_okno
    ON szablon_promptu_design(okno, zaktualizowano DESC, id DESC);
