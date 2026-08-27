-- Migracja tworzy profil asystenta jako byt trwały: kontrakt niesie profileId od
-- pierwszej wersji, lecz do tej pory nie było kolumny, w której dałoby się profil
-- zapamiętać.

-- Kolumny srodowisko_wykonania i tryb_uprawnien trzymają wartości kontraktu wprost,
-- bez tłumaczenia na polskie nazwy stosowane w tabeli okno_komunikacji.
CREATE TABLE profil_asystenta (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,

    -- warstwa_promptu jest powodem istnienia tabeli: instrukcja systemowa dokładana do wywołania kanału.
    warstwa_promptu          TEXT,

    -- kanal_modelu wskazuje kanał profilu; pusty znaczy kanał okna, nie brak kanału.
    kanal_modelu             TEXT,

    -- glos_syntezy niesie nastawę głosu odczytu; rdzeń syntezy dziś jej nie wykonuje.
    glos_syntezy             TEXT,

    -- srodowisko_wykonania to zasięg maszyny profilu; local domyślnie znaczy urządzenie Operatora.
    srodowisko_wykonania     TEXT    NOT NULL DEFAULT 'local'
                                     CHECK(srodowisko_wykonania IN ('local','core','remote')),

    -- tryb_uprawnien to zasięg pracy, nie bramka; domyślnie pełny, żeby świeży profil nie pytał o zgodę.
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

-- Zlecenie pamięta profil, którym zostało wydane, bo wykonawca czyta je z bazy
-- dłużej niż trwa żądanie; klucz obcy jest miękki, żeby skasowanie profilu nie
-- kasowało zlecenia.
ALTER TABLE zlecenie_asystenta ADD COLUMN profil_kod TEXT
    REFERENCES profil_asystenta(identyfikator_zewnetrzny) ON DELETE SET NULL;
CREATE INDEX idx_zlecenie_asystenta_profil ON zlecenie_asystenta(profil_kod, zaktualizowano DESC);
