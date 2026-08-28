-- Migracja zakłada wykaz zdalnych hostów wykonania, na których rdzeń może uruchomić
-- proces przez SSH, z osobną zgodą dla każdego.

CREATE TABLE host_zdalny (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Nazwa hosta z ustawienia `host_wykonania`; klucz dopasowania wyboru.
    nazwa          TEXT    NOT NULL UNIQUE,
    -- Adres sieciowy dla SSH; pusty znaczy: łącz się z nazwą.
    adres          TEXT    NOT NULL DEFAULT '',
    -- Konto na hoście zdalnym; puste znaczy: konto procesu rdzenia.
    uzytkownik     TEXT    NOT NULL DEFAULT '',
    port           INTEGER NOT NULL DEFAULT 22 CHECK(port BETWEEN 1 AND 65535),
    -- Zgoda na inicjowanie połączeń z tym hostem; domyślnie brak.
    zgoda          INTEGER NOT NULL DEFAULT 0 CHECK(zgoda IN (0,1)),
    -- Znacznik chwili wydania zgody; pusty, póki zgody nie wydano.
    zgode_wydano   TEXT    NOT NULL DEFAULT '',
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
