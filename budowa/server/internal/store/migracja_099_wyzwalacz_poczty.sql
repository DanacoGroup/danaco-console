-- Przebudowuje tabelę wyzwalacz_automatyki, aby ograniczenie CHECK dopuściło dodatkowy rodzaj wyzwalacza mail, ponieważ SQLite nie zmienia ograniczeń CHECK w istniejącej tabeli.

CREATE TABLE wyzwalacz_automatyki_nowa (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    harmonogram_id           INTEGER NOT NULL
                                     REFERENCES harmonogram_automatyki(id) ON DELETE CASCADE,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('cron','webhook','file','condition','mail')),
    wyrazenie                TEXT    NOT NULL,
    czynny                   INTEGER NOT NULL DEFAULT 1 CHECK(czynny IN (0,1)),
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO wyzwalacz_automatyki_nowa
    (id, harmonogram_id, identyfikator_zewnetrzny, rodzaj, wyrazenie,
     czynny, kolejnosc, utworzono)
SELECT
     id, harmonogram_id, identyfikator_zewnetrzny, rodzaj, wyrazenie,
     czynny, kolejnosc, utworzono
  FROM wyzwalacz_automatyki;

DROP TABLE wyzwalacz_automatyki;

ALTER TABLE wyzwalacz_automatyki_nowa RENAME TO wyzwalacz_automatyki;

-- Odtwarza indeks tabeli wyzwalacz_automatyki po jej przebudowie, zachowując porządek zapytań według harmonogramu i kolejności.
CREATE INDEX idx_wyzwalacz_automatyki_harmonogram
    ON wyzwalacz_automatyki(harmonogram_id, kolejnosc, id);

-- Doręczenie listu pyta: „które czynne wyzwalacze rodzaju mail pasują do tej
-- skrzynki". Indeks częściowy, bo pozostałe rodzaje obsługują inne odbiorniki.
CREATE INDEX idx_wyzwalacz_automatyki_mail
    ON wyzwalacz_automatyki(wyrazenie)
    WHERE rodzaj = 'mail' AND czynny = 1;
