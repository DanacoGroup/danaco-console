-- Wyzwalacz „przyszedł list": rodzaj `mail` w słowniku wyzwalaczy automatyki.
--
-- Odbiór listu robi dwie rzeczy jednym mechanizmem:
--
--   1. uruchamia automatyki z czynnym wyzwalaczem rodzaju `mail` — wyzwalacz
--      startuje bieg od początku;
--   2. doręcza sygnał `mail:<kod skrzynki>` silnikowi wybudzeń — sygnał wznawia
--      konkretny bieg zawieszony krokiem oczekiwania.
--
-- Wyrażenie wyzwalacza `mail` wskazuje skrzynkę: kod skrzynki pocztowej albo
-- `*` (puste znaczy to samo). `automation.schedule.set` pomija wyzwalacze
-- z pustym wyrażeniem, więc drogą kontraktu „każda skrzynka" zapisuje się jako
-- `*`.
--
-- SQLite nie zmienia CHECK w miejscu, więc tabela idzie przez przebudowę.
-- Na `wyzwalacz_automatyki` nie wskazuje żaden klucz obcy, więc przebudowa
-- obejmuje jedną tabelę.

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

-- Indeks odtworzony po przebudowie tabeli.
CREATE INDEX idx_wyzwalacz_automatyki_harmonogram
    ON wyzwalacz_automatyki(harmonogram_id, kolejnosc, id);

-- Doręczenie listu pyta: „które czynne wyzwalacze rodzaju mail pasują do tej
-- skrzynki". Indeks częściowy, bo pozostałe rodzaje obsługują inne odbiorniki.
CREATE INDEX idx_wyzwalacz_automatyki_mail
    ON wyzwalacz_automatyki(wyrazenie)
    WHERE rodzaj = 'mail' AND czynny = 1;
