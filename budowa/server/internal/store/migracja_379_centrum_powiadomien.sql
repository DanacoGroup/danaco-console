-- Migracja 379 — rejestr centrum powiadomień.
--
-- Encja wprost z `architektura/model-danych.md` rozdz. 18.4: klasa, waga,
-- źródło polimorficzne, środowisko prezentacji, sesja pochodzenia, treść,
-- działania, stan, kanał dostarczenia i znacznik czasu.
--
-- ── czemu nie tabela `powiadomienie` ────────────────────────────────────────
-- Nazwa `powiadomienie` jest w tej bazie ZAJĘTA od migracji 124 przez kolejkę
-- doręczeń funkcji Mobile: wiersz tamtej tabeli to jedna próba wypchnięcia na
-- urządzenie, z ponowieniami, terminem ważności i dziennikiem doręczeń. To jest
-- inny byt niż rejestr centrum — jeden opisuje ZDARZENIE, drugie DORĘCZENIE tego
-- zdarzenia. Dołożenie kolumn 18.4 do tamtej tabeli zlepiłoby dwa cykle życia
-- w jeden wiersz; zmiana jej nazwy ruszyłaby działający silnik kolejki. Rejestr
-- dostaje więc własną tabelę, a kolumna `kanal_dostarczenia` mówi, czy zdarzenie
-- poszło również tamtą drogą.
--
-- ── stan `odlozone` ma chwilę powrotu ──────────────────────────────────────
-- Odłożenie bez terminu jest cichym skasowaniem: zdarzenie znika z widoku i nie
-- wraca. Więz niżej nie pozwala zapisać jednego bez drugiego.
--
-- ── źródło polimorficzne bez klucza obcego ─────────────────────────────────
-- Para `zrodlo_typ` / `zrodlo_id` wskazuje siedem różnych tabel, więc klucza
-- obcego mieć nie może — tak samo jak kotwica w tabeli `powiadomienie`
-- z migracji 124 i w `sugestia_aod`. Więz pilnuje przynajmniej tego, żeby nie
-- zapisać połowy wskazania.

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

    -- Środowisko prezentacji i sesja pochodzenia idą kodem, nie kluczem obcym:
    -- zdarzenie ma przeżyć skasowanie karty sesji, o której mówi. Rejestr, który
    -- gubi zdarzenie razem z jego źródłem, gubi je zawsze wtedy, gdy jest
    -- najbardziej potrzebne — po awarii.
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

-- Odczyt gorący centrum: „co pokazać w kolumnie", od najnowszego.
CREATE INDEX idx_powiadomienie_centrum_widok
    ON powiadomienie_centrum(stan, znacznik_czasu DESC, id DESC);

-- Licznik plakietki liczy wyłącznie stan `nowe` — indeks częściowy, bo pozostałe
-- wiersze nigdy do niego nie wchodzą.
CREATE INDEX idx_powiadomienie_centrum_nowe
    ON powiadomienie_centrum(id) WHERE stan = 'nowe';

-- Zawężenie kolumny do środowiska albo karty sesji (filtr źródła, rozdz. 11.6).
CREATE INDEX idx_powiadomienie_centrum_zrodlo
    ON powiadomienie_centrum(srodowisko_kod, sesja_id, znacznik_czasu DESC);

-- Powrót zdarzeń odłożonych: „co wraca do stanu nowe w tej chwili".
CREATE INDEX idx_powiadomienie_centrum_odlozone
    ON powiadomienie_centrum(odlozone_do) WHERE stan = 'odlozone';
