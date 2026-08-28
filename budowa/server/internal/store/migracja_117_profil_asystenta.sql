-- Migracja 117 — profil asystenta jako byt trwały.
--
-- `assistant.voice.command` niesie `profileId` od pierwszej wersji kontraktu, ale
-- migracja 050 (zlecenie + dziennik) nie zna ani profilu, ani kolumny, w której
-- dałoby się go zapamiętać. Adapter przyjmował więc wskazanie profilu i porzucał
-- je bez śladu — pole żądania było etykietą bez bytu po drugiej stronie.
--
-- Profil nie jest ozdobą wykazu okien. Niesie warstwę promptu, która mówi
-- modelowi „jesteś klawiaturą Operatora, nie autorem" — a to ona rozstrzyga, czy
-- model sięgnie po narzędzia platformy, czy odpowie tekstem we własnym oknie. Bez
-- miejsca na jej zapis asystent wraca do bycia czatem przy pierwszym uruchomieniu,
-- bo warstwy nie ma gdzie postawić. Stąd byt, nie etykieta.
--
-- Profil nie jest zestawem uprawnień blokujących. Żadna kolumna tej tabeli nie
-- dopuszcza ani nie odmawia czynności; wszystkie są nastawą zasięgu pracy,
-- a zasięg domyślny jest pełny. Dlatego `tryb_uprawnien` startuje z wartości
-- najszerszej, a nie najwęższej — profil świeżo założony ma pracować, nie prosić
-- o zgodę.
--
-- Skąd pola. Wyłącznie z tego, co niesie kontrakt i widok okna asystenta — pola
-- bez odbiorcy nie ma tu ani jednego:
--
--   `Profil: Asystent biurowy`                      → nazwa
--   `Konfiguruj profil` (formularz warstwy promptu)  → warstwa_promptu
--   kanał modelu okna asystenta                      → kanal_modelu
--   `Synteza mowy — … · profil formalny`            → glos_syntezy
--   chip `Ten komputer`                              → srodowisko_wykonania
--   chip `uprawnienia: praca`                        → tryb_uprawnien
--   profil wskazywany oknu bez wyboru za każdym razem → domyslny
--
-- Osobne rozmowy przypisane do profilu wymagałyby kolumny w `okno_komunikacji`
-- albo w wykazie rozmów — to cudza tabela, więc ten krok jej nie dotyka. Nie ma
-- też przełącznika „synteza wł./wył.": pusty `glos_syntezy` znaczy brak syntezy,
-- a osobna flaga byłaby drugą prawdą o tym samym.
--
-- Czas jest liczbą — milisekundy epoki, jak w migracji 050 i 043.

-- ── Profil asystenta ─────────────────────────────────────────────────────────
-- `srodowisko_wykonania` i `tryb_uprawnien` trzymają wartości kontraktu wprost
-- (shared.ExecutionEnv, shared.PermissionMode), bez tłumaczenia — tak jak `stan`
-- i `droga` w migracji 050. Dzięki temu nastawa profilu wchodzi do
-- `models.Zapytanie` bez ani jednego przekładu po drodze; przekład byłby
-- miejscem, w którym „Ten komputer" mogłoby po cichu stać się czymś innym.
--
-- Uwaga na `tryb_uprawnien`: kolumna `okno_komunikacji.tryb_uprawnien` trzyma
-- wartości polskie (shared.WartosciBazyPermissionMode). Tutaj świadomie nie —
-- profil nie jest oknem i nie jedzie tamtym przekładem; wartość kontraktu jest
-- tym, co czyta rdzeń składający turę.
CREATE TABLE profil_asystenta (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,

    -- warstwa_promptu jest powodem istnienia tej tabeli: instrukcja systemowa,
    -- którą profil dokłada do wywołania kanału. Pusta znaczy profil bez zdania
    -- do powiedzenia — zlecenie idzie wtedy tak, jak szło przed tą migracją
    --, a nie jest odrzucane.
    warstwa_promptu          TEXT,

    -- kanal_modelu wskazuje kanał, którym profil pracuje. NULL znaczy „kanał
    -- okna" — profil nie ma obowiązku przykrywać wskazania okna, a przykrycie
    -- pustym napisem zabrałoby turze jedyny kanał, jaki miała.
    kanal_modelu             TEXT,

    -- glos_syntezy niesie nastawę głosu odczytu odpowiedzi (`profil formalny`
    -- z kompilacji). Rdzeń syntezy dziś nie wykonuje — kolumna jest miejscem
    -- zapisu nastawy, nie obietnicą, że coś przemówi.
    glos_syntezy             TEXT,

    -- srodowisko_wykonania to chip `Ten komputer`: zasięg maszyny, na której
    -- profil pracuje. `local` = urządzenie Operatora, i to jest wartość
    -- domyślna, bo dokładnie tak brzmi chip na kompilacji.
    srodowisko_wykonania     TEXT    NOT NULL DEFAULT 'local'
                                     CHECK(srodowisko_wykonania IN ('local','core','remote')),

    -- tryb_uprawnien to chip `uprawnienia: praca` — nastawa zasięgu pracy, nie
    -- bramka. Domyślnie pełny (`bypassPermissions`), bo zero bramek znaczy, że
    -- profil świeżo założony nie zatrzymuje się na pytaniu o zgodę.
    tryb_uprawnien           TEXT    NOT NULL DEFAULT 'bypassPermissions'
                                     CHECK(tryb_uprawnien IN ('manual','acceptEdits','plan',
                                                              'auto','dontAsk','bypassPermissions')),

    -- domyslny wskazuje profil brany, gdy żądanie nie poda `profileId`.
    domyslny                 INTEGER NOT NULL DEFAULT 0 CHECK(domyslny IN (0,1)),

    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL
);

-- Jeden domyślny albo żaden. Indeks częściowy pilnuje tego w schemacie, a nie
-- w kodzie: dwa profile domyślne znaczyłyby, że o warstwie promptu asystenta
-- rozstrzyga kolejność wierszy, czyli przypadek.
CREATE UNIQUE INDEX idx_profil_asystenta_domyslny ON profil_asystenta(domyslny) WHERE domyslny = 1;

-- ── Zlecenie pamięta, którym profilem zostało wydane ─────────────────────────
-- Bez tej kolumny profil nie przeżyłby zlecenia, i to jest jej całe uzasadnienie.
-- `assistant.voice.command` potwierdza przyjęcie od razu, a tura biegnie
-- gorutyną dłużej niż komenda: wykonawca czyta zlecenie z bazy, nie
-- z żądania, którego już nie ma. Profil trzymany wyłącznie w pamięci procesu
-- zniknąłby też przy restarcie rdzenia w trakcie zlecenia.
--
-- Klucz obcy jest miękki (ON DELETE SET NULL): skasowanie profilu nie ma prawa
-- wymazać zlecenia, które naprawdę się odbyło — traci ono wtedy wskazanie
-- warstwy, nie swój ślad.
ALTER TABLE zlecenie_asystenta ADD COLUMN profil_kod TEXT
    REFERENCES profil_asystenta(identyfikator_zewnetrzny) ON DELETE SET NULL;
CREATE INDEX idx_zlecenie_asystenta_profil ON zlecenie_asystenta(profil_kod, zaktualizowano DESC);
