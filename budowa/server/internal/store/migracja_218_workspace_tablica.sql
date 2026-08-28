-- Migracja 218 — kolumny tablicy kanban projektu.
--
-- Kolumna jest nastawą Operatora i odwzorowuje się na stan zadania z zestawu
-- wyjściowego platformy. Kolumna własna nie zakłada stanu nowego — dlatego
-- kolumna `stan` ma ten sam warunek co zadanie, a nie własne słownictwo.
--
-- Granica prac w toku jest sygnalizowana, nie egzekwowana: zero znaczy brak
-- granicy, a przekroczenie wraca ostrzeżeniem odpowiedzi, nie odmową.

CREATE TABLE kolumna_tablicy_projektu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    projekt_id               INTEGER NOT NULL REFERENCES projekt(id) ON DELETE CASCADE,
    nazwa                    TEXT    NOT NULL,
    stan                     TEXT    NOT NULL
                                     CHECK(stan IN ('todo','inProgress','inReview',
                                                    'done','blocked','cancelled')),
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    granica_wip              INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_kolumna_tablicy_projekt ON kolumna_tablicy_projektu(projekt_id, kolejnosc);
