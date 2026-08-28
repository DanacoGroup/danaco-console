-- Migracja zakłada trwałość poczty rdzenia: skrzynkę pocztową oraz dziennik listów
-- odebranych i wysłanych przez tę skrzynkę.

-- Host i port są konfiguracją skrzynki, nie stałą kodu, bo schemat nie przesądza, gdzie
-- skrzynka faktycznie stoi.
CREATE TABLE skrzynka_pocztowa (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    kod             TEXT    NOT NULL UNIQUE,
    -- Adres skrzynki: nadawca listów wychodzących i adresat przychodzących.
    adres           TEXT    NOT NULL,
    host_odbioru    TEXT    NOT NULL,
    port_odbioru    INTEGER NOT NULL DEFAULT 995,
    host_wysylki    TEXT    NOT NULL,
    port_wysylki    INTEGER NOT NULL DEFAULT 587,
    uzytkownik      TEXT    NOT NULL,
    -- Odwołanie sejfu poświadczeń; brak wartości znaczy skrzynkę bez hasła.
    haslo_odwolanie TEXT,
    tls_weryfikacja INTEGER NOT NULL DEFAULT 1 CHECK(tls_weryfikacja IN (0,1)),
    -- Takt obserwatora poczty w sekundach jest konfiguracją skrzynki, nie stałą wspólną.
    takt_sekundy    INTEGER NOT NULL DEFAULT 60 CHECK(takt_sekundy > 0),
    aktywna         INTEGER NOT NULL DEFAULT 1 CHECK(aktywna IN (0,1)),
    utworzono       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Znacznik przetworzenia stawia krok odbioru listu, nie sam obserwator poczty, żeby było
-- jedno miejsce prawdy o tym, co automatyka już widziała.
CREATE TABLE list_odebrany (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    skrzynka_id         INTEGER NOT NULL
                                REFERENCES skrzynka_pocztowa(id) ON DELETE CASCADE,
    identyfikator_listu TEXT    NOT NULL,
    nadawca             TEXT    NOT NULL,
    temat               TEXT    NOT NULL DEFAULT '',
    -- Chwila z nagłówka listu; brak wartości, gdy list daty nie niesie.
    chwila              TEXT,
    tresc               TEXT    NOT NULL DEFAULT '',
    odebrano            TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    przetworzony        INTEGER NOT NULL DEFAULT 0 CHECK(przetworzony IN (0,1)),
    UNIQUE(skrzynka_id, identyfikator_listu)
);

-- Krok odbioru listu pyta wyłącznie o listy nieprzetworzone własnej skrzynki pocztowej,
-- bez przeglądu pozostałych.
CREATE INDEX idx_list_odebrany_nieprzetworzone
    ON list_odebrany(skrzynka_id, id)
    WHERE przetworzony = 0;

-- Ślad listów wysłanych niesie temat, treść i wynik wysyłki, udanej albo nieudanej,
-- wraz z treścią błędu wysyłki.
CREATE TABLE list_wyslany (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    skrzynka_id     INTEGER NOT NULL
                            REFERENCES skrzynka_pocztowa(id) ON DELETE CASCADE,
    adresat         TEXT    NOT NULL,
    temat           TEXT    NOT NULL DEFAULT '',
    tresc           TEXT    NOT NULL DEFAULT '',
    -- Identyfikator listu odebranego, na który ten odpowiada; NULL dla listu
    -- nadanego bez odpowiedzi.
    w_odpowiedzi_na TEXT,
    chwila          TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    powodzenie      INTEGER NOT NULL CHECK(powodzenie IN (0,1)),
    blad            TEXT,
    -- Wysyłka udana nie ma błędu; nieudana musi go mieć.
    CHECK ((powodzenie = 1 AND blad IS NULL) OR (powodzenie = 0 AND blad IS NOT NULL))
);

CREATE INDEX idx_list_wyslany_skrzynka ON list_wyslany(skrzynka_id, id DESC);
