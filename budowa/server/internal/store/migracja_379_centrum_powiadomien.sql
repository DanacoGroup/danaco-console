-- Migracja 379 dodaje rejestr centrum powiadomień jako odrębny byt od
-- kolejki doręczeń funkcji Mobile.

CREATE TABLE powiadomienie_centrum (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,

    klasa                  TEXT    NOT NULL
                                   CHECK(klasa IN ('zakonczenie','decyzja','blad','wzmianka',
                                                   'termin','automatyka','system')),
    waga                   TEXT    NOT NULL
                                   CHECK(waga IN ('informacyjna','normalna','wymagajaca_decyzji')),

    tresc                  TEXT    NOT NULL CHECK(TRIM(tresc) <> ''),

    zrodlo_typ             TEXT    CHECK(zrodlo_typ IN ('srodowisko','modul','projekt','sesja',
                                                        'zadanie','przebieg_petli','automatyka')),
    zrodlo_id              TEXT,

    -- Środowisko i sesja idą kodem, nie kluczem obcym, aby zdarzenie
    -- przeżyło skasowanie karty sesji.
    srodowisko_kod         TEXT    NOT NULL DEFAULT '',
    sesja_id               TEXT    NOT NULL DEFAULT '',

    -- Działania dostępne z poziomu pozycji; kształt NotificationAction[].
    akcje                  TEXT    NOT NULL DEFAULT '[]',

    stan                   TEXT    NOT NULL DEFAULT 'nowe'
                                   CHECK(stan IN ('nowe','odczytane','obsluzone','odlozone')),

    kanal_dostarczenia     TEXT    NOT NULL DEFAULT 'centrum'
                                   CHECK(kanal_dostarczenia IN ('centrum','centrum_i_push')),

    -- Urządzenie, na które zdarzenie zostało wypchnięte; puste znaczy „nigdzie".
    urzadzenie_id          INTEGER REFERENCES urzadzenie_powiadomien(id) ON DELETE SET NULL,

    odlozone_do            TEXT,

    znacznik_czasu         TEXT    NOT NULL
                                   DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),

    -- Pół wskazania nie wskazuje niczego.
    CHECK((zrodlo_typ IS NULL) = (zrodlo_id IS NULL)),
    -- Odłożenie bez chwili powrotu jest cichym skasowaniem.
    CHECK((stan = 'odlozone') = (odlozone_do IS NOT NULL))
);

-- Odczyt gorący centrum pobiera to, co należy pokazać w kolumnie
-- powiadomień, od najnowszego zdarzenia.
CREATE INDEX idx_powiadomienie_centrum_widok
    ON powiadomienie_centrum(stan, znacznik_czasu DESC, id DESC);

-- Licznik plakietki liczy wyłącznie zdarzenia w stanie nowe, stąd indeks
-- częściowy pomijający pozostałe wiersze.
CREATE INDEX idx_powiadomienie_centrum_nowe
    ON powiadomienie_centrum(id) WHERE stan = 'nowe';

-- Indeks zawęża odczyt do danego środowiska albo karty sesji, wspierając
-- filtrowanie źródła zdarzeń centrum.
CREATE INDEX idx_powiadomienie_centrum_zrodlo
    ON powiadomienie_centrum(srodowisko_kod, sesja_id, znacznik_czasu DESC);

-- Ten indeks częściowy wspiera powrót zdarzeń odłożonych, wskazując, co
-- wraca do stanu nowe w tej chwili.
CREATE INDEX idx_powiadomienie_centrum_odlozone
    ON powiadomienie_centrum(odlozone_do) WHERE stan = 'odlozone';
