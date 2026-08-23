-- Migracja 172 — monitory zmian strony (`browser.monitor.*`).
--
-- Monitor pilnuje strony między sprawdzeniami, więc musi trzymać odniesienie:
-- treść z chwili założenia albo z chwili ostatniego przyjętego sprawdzenia.
-- Bez odniesienia „zmieniło się" nie ma względem czego być prawdą, a komenda
-- `browser.monitor.check` oddawałaby zawsze `changed:false` albo zawsze
-- `changed:true` — i jedno, i drugie jest meldunkiem bez pomiaru.
--
-- Odniesienie idzie odwołaniem do pliku, nie treścią w kolumnie: strona bywa
-- setkami kilobajtów tekstu, a wzorem jest `migawka_strony.tekst_odwolanie`
-- z migracji 047. Obok odwołania stoi długość odniesienia — dzięki niej wykaz
-- monitorów mówi o rozmiarze pilnowanej treści bez sięgania po plik.
--
-- Stan sprawdzenia (`pending`, `unchanged`, `changed`, `failed`) jest kolumną,
-- a nie wyliczeniem z dat: sprawdzenie nieudane (`failed`) też jest
-- sprawdzeniem, ma swój czas i nie wolno go pomylić z „bez zmian".
CREATE TABLE monitor_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    selektor                 TEXT,
    interwal_sekund          INTEGER NOT NULL DEFAULT 3600,
    prog_zmiany              INTEGER,
    kanal_powiadomienia      TEXT,
    wlaczony                 INTEGER NOT NULL DEFAULT 1 CHECK(wlaczony IN (0,1)),
    stan                     TEXT    NOT NULL DEFAULT 'pending',
    odniesienie_odwolanie    TEXT,
    odniesienie_dlugosc      INTEGER,
    sprawdzono               TEXT,
    zmieniono                TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_monitor_przegladania_okno ON monitor_przegladania(okno, utworzono DESC, id);
