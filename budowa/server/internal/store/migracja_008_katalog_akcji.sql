-- Migracja 008 — katalog akcji sterowany danymi.
--
-- Panel akcji powstaje z rejestru: nowa akcja to nowy wiersz, nie zmiana w kodzie.
-- Ten sam rejestr zasila narzędzia modelu, więc kod nie zna ani jednej akcji.
-- Tabela naśladuje `kanal_modelu`: wiersz opisuje byt w całości, kod jest jego
-- trwałym identyfikatorem, `aktywna` rozstrzyga widoczność, `kolejnosc` porządek
-- prezentacji.
--
-- Zasięg. Akcja należy do jednego z ośmiu poziomów zasięgu. `poziom_zasiegu_id`
-- wskazuje poziom, `klucz_zasiegu` konkretny byt tego poziomu — kod modułu dla
-- poziomu `modul`, kod środowiska dla poziomu `srodowisko`, identyfikator okna
-- dla poziomu `okno`. Pusty `klucz_zasiegu` znaczy „każdy byt tego poziomu",
-- tak samo jak w tabeli `ustawienie`.
--
-- Komenda. `komenda` niesie nazwę komendy kontraktu wywoływanej przez akcję.
-- Kolumna nie jest kluczem obcym — kontrakt mieszka w `shared/contract.json`,
-- nie w bazie. Akcja wskazująca komendę bez obsługiwacza dostaje odpowiedź
-- `*.unknown`, więc rozjazd katalogu z rdzeniem jest widoczny, a nie wywracający.
--
-- Warunek dostępności. Zamknięty zbiór bytów, których akcja wymaga, żeby dało
-- się ją zaoferować: pusty (nic nie wymaga), `srodowisko`, `sesja`, `okno`,
-- `kolejka`, `kanal`. Wartość wyprowadzona z pól obowiązkowych żądania komendy
-- w `shared/contract.json`. Warunek nie wygasza kontrolki — rozstrzyga, czy
-- akcja trafia do panelu danego kontekstu.

CREATE TABLE akcja (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                    TEXT    NOT NULL UNIQUE,
    nazwa                  TEXT    NOT NULL,
    opis                   TEXT    NOT NULL DEFAULT '',
    ikona                  TEXT    NOT NULL DEFAULT '',
    poziom_zasiegu_id      INTEGER NOT NULL REFERENCES poziom_zasiegu(id) ON DELETE CASCADE,
    klucz_zasiegu          TEXT    NOT NULL DEFAULT '',
    komenda                TEXT    NOT NULL,
    warunek_dostepnosci    TEXT    NOT NULL DEFAULT ''
                                   CHECK(warunek_dostepnosci IN ('','srodowisko','sesja',
                                                                 'okno','kolejka','kanal')),
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    aktywna                INTEGER NOT NULL DEFAULT 1 CHECK(aktywna IN (0,1)),
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_akcja_zasieg ON akcja(poziom_zasiegu_id, klucz_zasiegu, kolejnosc);
CREATE INDEX idx_akcja_komenda ON akcja(komenda);
