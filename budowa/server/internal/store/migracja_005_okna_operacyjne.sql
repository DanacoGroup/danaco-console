-- Migracja 005 — katalog okien operacyjnych.
--
-- Okno operacyjne to pozycja katalogu funkcji modułu, a nie okno komunikacji.
-- Okno komunikacji jest bytem wykonania jednej sesji; okno operacyjne jest
-- wpisem słownikowym mówiącym, jakie okna robocze niesie moduł.
--
-- `modul_id` dopuszcza NULL: trzy okna aplikacji (Konfiguracja, Ustawienia,
-- Strona główna) oraz okno konfiguracji punktów izolacji i okna ról
-- MultitaskingAI nie należą do żadnego modułu.
--
-- Role okna: wiodące · pomocnicze · monitor · kreator · zarządca.

CREATE TABLE okno_operacyjne (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                    TEXT    NOT NULL UNIQUE,
    modul_id               INTEGER REFERENCES modul(id) ON DELETE CASCADE,
    nazwa                  TEXT    NOT NULL,
    rola                   TEXT    NOT NULL
                                   CHECK(rola IN ('wiodace','pomocnicze','monitor','kreator','zarzadca')),
    kategoria              TEXT    NOT NULL
                                   CHECK(kategoria IN ('komunikacja','edycja','podglad','repozytorium',
                                                       'konstruktor','kolejka','monitor','narzedzia',
                                                       'zrodla','konfiguracja')),
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    aktywne                INTEGER NOT NULL DEFAULT 1 CHECK(aktywne IN (0,1))
);
CREATE INDEX idx_okno_operacyjne_modul ON okno_operacyjne(modul_id, kolejnosc);
CREATE INDEX idx_okno_operacyjne_kategoria ON okno_operacyjne(kategoria, kolejnosc);
