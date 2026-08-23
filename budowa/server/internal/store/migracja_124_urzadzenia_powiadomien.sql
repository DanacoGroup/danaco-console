-- Migracja 124 — powiadomienia: rejestracja urządzenia, kolejka i doręczenie.
--
-- Nośnik doręczenia już istnieje: `transport.Serwer.Rozglos` (rozgloszenie.go)
-- wysyła kopertę do otwartych gniazd i zwraca liczbę urządzeń, które ją przyjęły.
-- Rejestr połączeń zna tożsamość każdego gniazda wraz z `clientId`
-- (transport/tozsamosc.go, pole `IdKlienta`). Ta migracja nie zakłada własnego
-- kanału — opisuje wyłącznie to, czego nośnikowi brakuje: kogo wołać, czym
-- i czy doszło.
--
-- Nie ma tu drugiej tabeli urządzeń. Urządzenie opisuje tabela `urzadzenie`
-- (`nazwa_hosta`, `biezace`); druga tabela maszyn byłaby drugą prawdą o tym, czym
-- Operator dysponuje. `urzadzenie_powiadomien` nie opisuje maszyny — opisuje zgodę
-- tej maszyny na wołanie i drogę, którą wołanie idzie. Jedna maszyna może mieć
-- kilka takich dróg (pulpit i przeglądarka to dwa osobne gniazda o dwóch
-- `clientId`), więc relacja jest jeden-do-wielu, a nie kolumną doklejoną do
-- `urzadzenie`.
--
-- Słownik kanałów jest zamknięty na to, co rdzeń dziś potrafi doręczyć. `kanal`
-- dopuszcza jedną wartość: 'polaczenie'. CHECK dopuszczający kanał usługi
-- zewnętrznej (APNs, FCM) bez kodu, który go obsłuży, przepuściłby rejestrację,
-- której żaden takt nie doręczy — wiersz stanąłby w kolejce na zawsze. Założenie
-- kont i kluczy usługi zewnętrznej jest rozstrzygnięciem Właściciela; kanały
-- wejdą osobną migracją, razem z kodem.
--
-- `klucz_kanalu` dla kanału 'polaczenie' jest `clientId` — tym samym napisem,
-- który transport już dziś niesie jako `Tozsamosc.IdKlienta`. Nie zakładamy
-- nowego identyfikatora urządzenia: byłby drugą tożsamością tego samego gniazda.
--
-- Kolejka zamiast wysyłki wprost. Bez tabeli powiadomienie zgłoszone przy
-- zamkniętej aplikacji znika. Kolejka sprawia, że brak odbiorcy jest stanem,
-- a nie ciszą: wiersz `oczekuje` z licznikiem prób mówi wprost „mieliśmy zawołać
-- i nie było komu".
--
-- `wygasa` jest NOT NULL z rozmysłu. Powiadomienie bez terminu ważności wisi
-- wiecznie i po tygodniu postoju rdzenia Operator dostaje lawinę budzików
-- o sprawach dawno nieaktualnych. Termin wymuszony schematem znaczy, że każdy
-- wołacz musi odpowiedzieć na pytanie „do kiedy to ma sens".
--
-- `byt_rodzaj` + `byt_id` to kotwica miękka, bez klucza obcego. Powiadomienie
-- dotyczy czegoś — dziś przede wszystkim kroku wstrzymanego, czekającego na słowo
-- Operatora (`wstrzymanie_kroku`). Klucza obcego do jednej tabeli tu nie ma, bo
-- powiadomienie ma z założenia dotyczyć różnych bytów — kroku, zlecenia, biegu
-- automatyki — a klucz obcy zamknąłby je na jeden byt i wymusił kolumnę na każdy
-- następny. Ceną jest brak kaskady: powiadomienie o bycie usuniętym zostaje
-- w kolejce i wygasa własnym terminem.
--
-- Doręczenie jest osobną tabelą. „Doszło" bez wskazania, do którego urządzenia,
-- jest odpowiedzią nie do sprawdzenia: Operator ma pulpit i telefon naraz.
-- Kolumna `dostarczono` w `powiadomienie` mówi „dotarło gdziekolwiek" i to
-- wystarcza kolejce do zamknięcia sprawy; `powiadomienie_dostarczenie` mówi gdzie
-- i kiedy, i to jest odpowiedź, którą można pokazać Operatorowi bez zmyślania.

-- ── REJESTRACJA URZĄDZENIA DO POWIADOMIEŃ ───────────────────────────────────
CREATE TABLE urzadzenie_powiadomien (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    urzadzenie_id          INTEGER NOT NULL
                                   REFERENCES urzadzenie(id) ON DELETE CASCADE,

    -- Droga, którą wołanie idzie do tego urządzenia. Słownik zamknięty na to,
    -- co rdzeń dziś naprawdę umie — uzasadnienie w nagłówku pliku.
    kanal                  TEXT    NOT NULL DEFAULT 'polaczenie'
                                   CHECK(kanal IN ('polaczenie')),

    -- Adres w obrębie kanału. Dla 'polaczenie' jest to `clientId` gniazda,
    -- czyli `transport.Tozsamosc.IdKlienta`. Pusty napis nie jest adresem.
    klucz_kanalu           TEXT    NOT NULL CHECK(TRIM(klucz_kanalu) <> ''),

    -- Napis dla Operatora („Pulpit w biurze"). NULL znaczy „nie nazwał" —
    -- warstwa wyżej pokaże wtedy nazwę urządzenia, a nie nazwę zmyśloną.
    etykieta               TEXT,

    -- Zgoda czynna. Wyrejestrowanie nie kasuje wiersza (patrz `wyrejestrowano`),
    -- bo ślad po tym, że urządzenie kiedyś było wołane, jest częścią
    -- przejrzystości kanału.
    aktywne                INTEGER NOT NULL DEFAULT 1 CHECK(aktywne IN (0,1)),

    zarejestrowano         TEXT    NOT NULL
                                   DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    wyrejestrowano         TEXT,
    ostatnio_dostarczono   TEXT,

    -- Zgoda i data jej cofnięcia chodzą parą w obie strony. Bez tego więzu dałoby
    -- się zapisać rejestrację czynną z datą wyrejestrowania — czyli wiersz, który
    -- twierdzi dwie rzeczy naraz.
    CHECK((aktywne = 1) = (wyrejestrowano IS NULL)),

    -- Jedna droga na urządzenie i adres. Rejestracja powtórzona ma odświeżyć
    -- wiersz, a nie założyć drugi — inaczej jedno urządzenie dostawałoby
    -- to samo powiadomienie tyle razy, ile razy się przedstawiło.
    UNIQUE(urzadzenie_id, kanal, klucz_kanalu)
);

-- Odczyt gorący: „komu mam to teraz wysłać". Indeks częściowy, bo wierszy
-- nieczynnych ta droga nie ogląda nigdy.
CREATE INDEX idx_urzadzenie_powiadomien_czynne
    ON urzadzenie_powiadomien(kanal, klucz_kanalu) WHERE aktywne = 1;

CREATE INDEX idx_urzadzenie_powiadomien_urzadzenie
    ON urzadzenie_powiadomien(urzadzenie_id);

-- ── KOLEJKA POWIADOMIEŃ ─────────────────────────────────────────────────────
CREATE TABLE powiadomienie (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,

    tytul                  TEXT    NOT NULL CHECK(TRIM(tytul) <> ''),
    tresc                  TEXT    NOT NULL DEFAULT '',

    -- 'pilny' znaczy „Operator ma to zobaczyć teraz", nie „wyślij dwa razy".
    -- Dwie wartości, bo trzeciej nie da się dziś odróżnić zachowaniem, a stopień
    -- bez skutku byłby ozdobą.
    priorytet              TEXT    NOT NULL DEFAULT 'zwykly'
                                   CHECK(priorytet IN ('zwykly','pilny')),

    -- Po co dzwonimy — zdanie dla Operatora, nie kod. Powiadomienie bez powodu
    -- jest budzikiem, którego nie da się ocenić.
    powod                  TEXT    NOT NULL DEFAULT '',

    -- CZEGO dotyczy. Kotwica miękka, bez klucza obcego — uzasadnienie w nagłówku.
    -- Rodzaj i klucz chodzą parą; więz stoi niżej, razem z pozostałymi, bo
    -- SQLite nie pozwala wrócić do definicji kolumn po pierwszym więzie tabeli.
    byt_rodzaj             TEXT,
    byt_id                 TEXT,

    -- 'porzucone' to stan po wyczerpaniu prób przed terminem ważności;
    -- 'wygasle' to stan po przekroczeniu `wygasa`. Rozróżnienie ma treść:
    -- pierwsze mówi „wołaliśmy i nie było komu", drugie „przestało mieć sens".
    stan                   TEXT    NOT NULL DEFAULT 'oczekuje'
                                   CHECK(stan IN ('oczekuje','dostarczone','odwolane',
                                                  'wygasle','porzucone')),

    prob_ile               INTEGER NOT NULL DEFAULT 0 CHECK(prob_ile >= 0),
    nastepna_proba         TEXT    NOT NULL
                                   DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),

    -- NOT NULL z rozmysłu — uzasadnienie w nagłówku pliku.
    wygasa                 TEXT    NOT NULL,

    utworzono              TEXT    NOT NULL
                                   DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    dostarczono            TEXT,
    odwolano               TEXT,
    powod_odwolania        TEXT,

    -- Pół wskazania nie wskazuje niczego: rodzaj bez klucza i klucz bez rodzaju
    -- są wierszami, po których nie da się nic odnaleźć ani odwołać.
    CHECK((byt_rodzaj IS NULL) = (byt_id IS NULL)),

    -- Stan końcowy i jego znacznik chodzą parą. Wiersz `dostarczone` bez chwili
    -- doręczenia twierdziłby, że coś doszło, nie umiejąc powiedzieć kiedy —
    -- czyli dokładnie to, czego wymaganie „czy doszło, kiedy doszło" zabrania.
    CHECK((stan = 'dostarczone') = (dostarczono IS NOT NULL)),
    CHECK((stan = 'odwolane')    = (odwolano    IS NOT NULL))
);

-- Odczyt gorący budzika: „co jest należne w tej chwili". Indeks częściowy —
-- wierszy zamkniętych pętla nie ogląda nigdy, a to one z czasem stanowią
-- większość tabeli.
CREATE INDEX idx_powiadomienie_nalezne
    ON powiadomienie(nastepna_proba, id) WHERE stan = 'oczekuje';

-- Odwołanie idzie po bycie: „ta decyzja zapadła, zgaś wszystko, co o nią pyta".
CREATE INDEX idx_powiadomienie_byt
    ON powiadomienie(byt_rodzaj, byt_id);

-- ── DORĘCZENIE, PER URZĄDZENIE ──────────────────────────────────────────────
CREATE TABLE powiadomienie_dostarczenie (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    powiadomienie_id       INTEGER NOT NULL
                                   REFERENCES powiadomienie(id) ON DELETE CASCADE,
    -- Wskazanie na rejestrację, nie na urządzenie: doręczenie zaszło konkretną
    -- drogą i przy zmianie drogi ślad ma zostać przy tej, którą naprawdę poszło.
    urzadzenie_powiadomien_id INTEGER NOT NULL
                                   REFERENCES urzadzenie_powiadomien(id) ON DELETE CASCADE,

    dostarczono            TEXT    NOT NULL
                                   DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),

    -- Kiedy Operator potwierdził, że widział. NULL znaczy „doszło, ale nikt nie
    -- potwierdził" — i to jest odpowiedź uczciwa, nie brak danych.
    potwierdzono           TEXT,

    -- Jedno doręczenie na parę. Ponowienie po powrocie urządzenia nie ma
    -- dopisywać drugiego wiersza o tym samym fakcie.
    UNIQUE(powiadomienie_id, urzadzenie_powiadomien_id)
);

CREATE INDEX idx_powiadomienie_dostarczenie_urzadzenie
    ON powiadomienie_dostarczenie(urzadzenie_powiadomien_id, dostarczono);
