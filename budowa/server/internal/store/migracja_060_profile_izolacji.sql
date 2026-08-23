-- Migracja 060 — profile izolacji, ich przypisanie do poziomu zasięgu oraz
-- warstwa izolacji wybrana dla karty sesji albo okna.
--
-- Wartości jedenastu punktów izolacji (trzy wymiary kontekstu i osiem zakresów
-- technicznych) leżą w tabeli `ustawienie` pod adresem złożonym z poziomu
-- zasięgu i osi; rozstrzyga je pakiet `internal/konfig`, a egzekwuje
-- `internal/session`. Ta migracja dokłada trzy rzeczy, których w tym układzie
-- nie ma:
--   1. profil — nazwany zestaw przełączników, zapisywany raz i przypisywany
--      wielokrotnie (`isolation.profile.save/list/load/assign/delete`);
--   2. przypisanie profilu do poziomu — odpowiedź na pytanie, z którego profilu
--      pochodzi polityka obowiązująca na tym poziomie (`IsolationPolicy.profileId`);
--   3. wybór warstwy izolacji (`isolation.layer.set`).
--
-- Profil nie jest drugim miejscem wartości izolacji. Wartość obowiązująca leży
-- wyłącznie w tabeli `ustawienie` i tylko tam czyta ją rozstrzygacz. Profil jest
-- szablonem: `isolation.profile.assign` przepisuje jego przełączniki do
-- `ustawienie` pod wskazany adres i dopiero to daje skutek wykonawczy. Profil
-- jako drugie źródło prawdy zmuszałby egzekutora do czytania dwóch miejsc.
--
-- Przełącznik jest wierszem, nie kolumną: tabela niesie parę `klucz` + `wartosc`
-- dokładnie w postaci, w jakiej czyta je rozstrzygacz (`izolacja_historia` =
-- `odrebna`, `izolacja_dostep_sieciowy` = `wylaczony` — patrz
-- `server/internal/konfig/definicje_izolacji.go`). Przypisanie profilu jest więc
-- przepisaniem wiersza jeden do jednego, bez przekładu nazw, a nowy punkt
-- izolacji wymaga wiersza, nie zmiany schematu.
--
-- Profil bez kompletu punktów jest poprawny: kontrakt oznacza `contextSwitches`
-- i `technicalSwitches` jako pola nieobowiązkowe. Punkt nieujęty w profilu nie
-- jest przy przypisaniu ruszany — obowiązuje to, co na poziomie stoi, a w
-- ostateczności wartość domyślna.
--
-- Przypisanie jest jedno na adres: warunek UNIQUE stoi na parze
-- (poziom_zasiegu_id, klucz_zasiegu), bo pytanie o źródło polityki danego bytu
-- ma dokładnie jedną odpowiedź. Kolejne przypisanie pod ten sam adres zastępuje
-- poprzednie, tak jak zastępują się wartości w `ustawienie`.
--
-- Warstwy nie ma w adresie zapisu i ta migracja jej tam nie wnosi. Kontrakt zna
-- dwie warstwy: `default` (domyślna platformy) i `session` (warstwa karty
-- sesji). Adresem zapisu ustawienia jest czwórka (poziom zasięgu, byt poziomu,
-- oś, byt osi), a karta sesji jest jednym z poziomów zasięgu (`karta_sesji`).
-- Kolumna `warstwa` w tabeli `ustawienie` zbudowałaby drugi wymiar adresu,
-- którego rozstrzygacz nie czyta: wartość zapisana na warstwie sesyjnej nie
-- doszłaby do egzekutora, a zapis zostałby potwierdzony bez skutku. Tabela
-- `warstwa_izolacji` trzyma więc wyłącznie wybór — którą warstwę pokazuje
-- i zapisuje panel izolacji dla danej karty sesji albo okna.
--
-- Usunięcie profilu zabiera jego przełączniki i przypisania (ON DELETE CASCADE),
-- ale nie cofa wartości już przepisanych do `ustawienie` — te są odtąd
-- wartościami poziomu, nie własnością profilu. Cofa je `config.reset` albo
-- ponowny zapis przełączników.

CREATE TABLE profil_izolacji (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    kod            TEXT    NOT NULL UNIQUE,
    nazwa          TEXT    NOT NULL,
    opis           TEXT    NOT NULL DEFAULT '',
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE przelacznik_profilu_izolacji (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    profil_id INTEGER NOT NULL REFERENCES profil_izolacji(id) ON DELETE CASCADE,
    klucz     TEXT    NOT NULL,
    wartosc   TEXT    NOT NULL,
    UNIQUE(profil_id, klucz)
);

CREATE TABLE przypisanie_profilu_izolacji (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    profil_id         INTEGER NOT NULL REFERENCES profil_izolacji(id) ON DELETE CASCADE,
    poziom_zasiegu_id INTEGER NOT NULL REFERENCES poziom_zasiegu(id) ON DELETE CASCADE,
    klucz_zasiegu     TEXT    NOT NULL DEFAULT '',
    przypisano        TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(poziom_zasiegu_id, klucz_zasiegu)
);

CREATE TABLE warstwa_izolacji (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    poziom_zasiegu_id INTEGER NOT NULL REFERENCES poziom_zasiegu(id) ON DELETE CASCADE,
    klucz_zasiegu     TEXT    NOT NULL DEFAULT '',
    warstwa           TEXT    NOT NULL CHECK(warstwa IN ('default','session')),
    zaktualizowano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(poziom_zasiegu_id, klucz_zasiegu)
);
