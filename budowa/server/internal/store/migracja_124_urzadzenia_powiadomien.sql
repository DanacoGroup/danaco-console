-- Migracja 124 dodaje rejestrację urządzeń do powiadomień, kolejkę powiadomień
-- oraz doręczenia; nośnik połączenia już istnieje w warstwie transportu,
-- a migracja opisuje wyłącznie to, czego mu brakuje.

-- Rejestracja opisuje zgodę urządzenia na wołanie i drogę, którą wołanie idzie,
-- a nie samo urządzenie; jedna maszyna może mieć kilka takich dróg, więc
-- relacja jest jeden do wielu.
CREATE TABLE urzadzenie_powiadomien (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    urzadzenie_id          INTEGER NOT NULL
                                   REFERENCES urzadzenie(id) ON DELETE CASCADE,

    -- Droga, którą wołanie idzie do urządzenia; słownik dziś dopuszcza jedynie
    -- wartość połączenia.
    kanal                  TEXT    NOT NULL DEFAULT 'polaczenie'
                                   CHECK(kanal IN ('polaczenie')),

    -- Adres w obrębie kanału; dla połączenia jest to tożsamość gniazda
    -- transportu, nigdy pusty napis.
    klucz_kanalu           TEXT    NOT NULL CHECK(TRIM(klucz_kanalu) <> ''),

    -- Napis dla operatora; wartość pusta oznacza brak nazwania i pokazuje
    -- wtedy nazwę urządzenia.
    etykieta               TEXT,

    -- Zgoda czynna; wyrejestrowanie nie kasuje wiersza, bo ślad dawnego
    -- wołania zostaje częścią historii.
    aktywne                INTEGER NOT NULL DEFAULT 1 CHECK(aktywne IN (0,1)),

    zarejestrowano         TEXT    NOT NULL
                                   DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    wyrejestrowano         TEXT,
    ostatnio_dostarczono   TEXT,

    -- Zgoda i data jej cofnięcia chodzą parą, więc rejestracja czynna nie może
    -- nieść daty wyrejestrowania.
    CHECK((aktywne = 1) = (wyrejestrowano IS NULL)),

    -- Jedna droga na urządzenie i adres; powtórna rejestracja odświeża wiersz,
    -- nie zakłada drugiego.
    UNIQUE(urzadzenie_id, kanal, klucz_kanalu)
);

-- Odczyt gorący: „komu mam to teraz wysłać". Indeks częściowy, bo wierszy
-- nieczynnych ta droga nie ogląda nigdy.
CREATE INDEX idx_urzadzenie_powiadomien_czynne
    ON urzadzenie_powiadomien(kanal, klucz_kanalu) WHERE aktywne = 1;

CREATE INDEX idx_urzadzenie_powiadomien_urzadzenie
    ON urzadzenie_powiadomien(urzadzenie_id);

-- Kolejka zastępuje wysyłkę wprost: powiadomienie zgłoszone przy zamkniętej
-- aplikacji czeka w kolejce zamiast zniknąć, a licznik prób i termin
-- ważności rozstrzygają jego dalszy los.
CREATE TABLE powiadomienie (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,

    tytul                  TEXT    NOT NULL CHECK(TRIM(tytul) <> ''),
    tresc                  TEXT    NOT NULL DEFAULT '',

    -- Pilne znaczy, że operator ma to zobaczyć teraz, nie że wysyłka powtarza
    -- się dwa razy.
    priorytet              TEXT    NOT NULL DEFAULT 'zwykly'
                                   CHECK(priorytet IN ('zwykly','pilny')),

    -- Powód jest zdaniem dla operatora, nie kodem; powiadomienie bez powodu
    -- jest budzikiem bez treści.
    powod                  TEXT    NOT NULL DEFAULT '',

    -- Wskazuje, czego powiadomienie dotyczy; kotwica jest miękka, bez klucza
    -- obcego do jednej tabeli.
    byt_rodzaj             TEXT,
    byt_id                 TEXT,

    -- Porzucone znaczy wyczerpane próby przed terminem; wygasłe znaczy
    -- przekroczony termin ważności.
    stan                   TEXT    NOT NULL DEFAULT 'oczekuje'
                                   CHECK(stan IN ('oczekuje','dostarczone','odwolane',
                                                  'wygasle','porzucone')),

    prob_ile               INTEGER NOT NULL DEFAULT 0 CHECK(prob_ile >= 0),
    nastepna_proba         TEXT    NOT NULL
                                   DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),

    -- Termin ważności jest obowiązkowy, żeby powiadomienie nie wisiało bez
    -- końca po ciszy rdzenia.
    wygasa                 TEXT    NOT NULL,

    utworzono              TEXT    NOT NULL
                                   DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    dostarczono            TEXT,
    odwolano               TEXT,
    powod_odwolania        TEXT,

    -- Pół wskazania nie wskazuje niczego: rodzaj bez klucza i klucz bez rodzaju
    -- są wierszami bez adresu.
    CHECK((byt_rodzaj IS NULL) = (byt_id IS NULL)),

    -- Stan końcowy i jego znacznik chodzą parą; wiersz dostarczony musi nieść
    -- chwilę doręczenia.
    CHECK((stan = 'dostarczone') = (dostarczono IS NOT NULL)),
    CHECK((stan = 'odwolane')    = (odwolano    IS NOT NULL))
);

-- Odczyt gorący budzika: „co jest należne w tej chwili". Indeks częściowy —
-- wierszy zamkniętych pętla nie ogląda nigdy, a to one z czasem stanowią
-- większość tabeli.
CREATE INDEX idx_powiadomienie_nalezne
    ON powiadomienie(nastepna_proba, id) WHERE stan = 'oczekuje';

-- Odwołanie działa po bycie: gdy decyzja dotycząca bytu zapada, ten indeks
-- pozwala zgasić od razu wszystkie powiadomienia, które o niego pytają.
CREATE INDEX idx_powiadomienie_byt
    ON powiadomienie(byt_rodzaj, byt_id);

-- Doręczenie jest osobną tabelą, bo jedno powiadomienie może dotrzeć do
-- operatora kilkoma urządzeniami naraz, a każde dotarcie ma własną chwilę
-- i własne potwierdzenie.
CREATE TABLE powiadomienie_dostarczenie (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    powiadomienie_id       INTEGER NOT NULL
                                   REFERENCES powiadomienie(id) ON DELETE CASCADE,
    -- Wskazuje na rejestrację, nie na urządzenie: ślad zostaje przy drodze,
    -- którą doręczenie poszło.
    urzadzenie_powiadomien_id INTEGER NOT NULL
                                   REFERENCES urzadzenie_powiadomien(id) ON DELETE CASCADE,

    dostarczono            TEXT    NOT NULL
                                   DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),

    -- Kiedy operator potwierdził odbiór; wartość pusta znaczy, że doszło,
    -- ale nikt nie potwierdził.
    potwierdzono           TEXT,

    -- Jedno doręczenie na parę; ponowienie po powrocie urządzenia nie
    -- dopisuje drugiego wiersza.
    UNIQUE(powiadomienie_id, urzadzenie_powiadomien_id)
);

CREATE INDEX idx_powiadomienie_dostarczenie_urzadzenie
    ON powiadomienie_dostarczenie(urzadzenie_powiadomien_id, dostarczono);
