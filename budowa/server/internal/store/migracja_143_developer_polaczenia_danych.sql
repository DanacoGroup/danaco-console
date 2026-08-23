-- Migracja 143 — połączenia bazodanowe okna Data Console modułu Developer.
--
-- Wiersz opisuje, DO CZEGO się łączyć, i nie niesie hasła. Hasło stoi w sejfie
-- pod odwołaniem zapisanym w `poswiadczenie` — ta sama zasada, którą kanały
-- modeli stosują do kluczy dostawców. Trzymanie hasła w tej tabeli oznaczałoby
-- kopię sekretu w drugim miejscu, którego nikt nie rotuje.
--
-- `tylko_odczyt` jest nastawą połączenia, nie podpowiedzią interfejsu.
-- Połączenie oznaczone jako tylko do odczytu odrzuca zapytanie zmieniające dane
-- w rdzeniu, zanim cokolwiek wyjdzie do silnika — konsola SQL nad produkcją bez
-- takiej nastawy jest jednym nieuważnym poleceniem od szkody.
CREATE TABLE developer_polaczenie_danych (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    kod           TEXT    NOT NULL UNIQUE,
    okno_kod      TEXT    NOT NULL,
    nazwa         TEXT    NOT NULL,
    silnik        TEXT    NOT NULL CHECK(silnik IN ('postgres','mysql','sqlite')),
    host          TEXT,
    port          INTEGER,
    baza          TEXT    NOT NULL,
    uzytkownik    TEXT,
    poswiadczenie TEXT,
    tylko_odczyt  INTEGER NOT NULL DEFAULT 0,
    zmieniono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_developer_polaczenie_okno ON developer_polaczenie_danych(okno_kod, nazwa);
