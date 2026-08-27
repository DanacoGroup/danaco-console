-- Migracja 015 zakłada katalog kategorii tożsamości modelu oraz treść tych
-- kategorii osobno dla platformy, modelu i konta, sterując dwuwarstwową
-- nakładką promptu systemowego.

-- Katalog kategorii zasad i tożsamości modelu naśladuje wzorzec innych
-- słowników platformy: kod jest trwałym identyfikatorem, kolejność ustala
-- porządek składania.
CREATE TABLE kategoria_tozsamosci (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                    TEXT    NOT NULL UNIQUE,
    nazwa                  TEXT    NOT NULL,
    opis                   TEXT    NOT NULL DEFAULT '',
    warstwa                TEXT    NOT NULL
                                   CHECK(warstwa IN ('constitution','profile','expertise')),
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    obowiazkowa            INTEGER NOT NULL DEFAULT 0 CHECK(obowiazkowa IN (0,1)),
    tryb_domyslny          TEXT    NOT NULL DEFAULT 'ZASTAP'
                                   CHECK(tryb_domyslny IN ('ZASTAP','DOLACZ')),
    aktywna                INTEGER NOT NULL DEFAULT 1 CHECK(aktywna IN (0,1)),
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Kolejność wewnątrz warstwy jest jednoznaczna, nie ozdobna: składacz ma dawać
-- ten sam prompt przy tej samej konfiguracji, a dwie kategorie o tej samej
-- kolejności czyniłyby porządek składania zależnym od przypadku.
CREATE UNIQUE INDEX idx_kategoria_tozsamosci_porzadek
    ON kategoria_tozsamosci(warstwa, kolejnosc);

-- Treść kategorii per oś rozdziela wartość dla platformy, dla konkretnego
-- modelu i dla konkretnego konta, bo są to trzy niezależne poziomy
-- obowiązywania tej samej kategorii.
CREATE TABLE dokument_tozsamosci (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kategoria_id           INTEGER NOT NULL REFERENCES kategoria_tozsamosci(id) ON DELETE CASCADE,
    os                     TEXT    NOT NULL DEFAULT 'platform'
                                   CHECK(os IN ('platform','model','account')),
    os_byt                 TEXT    NOT NULL DEFAULT '',
    tryb                   TEXT    NOT NULL DEFAULT 'ZASTAP'
                                   CHECK(tryb IN ('ZASTAP','DOLACZ')),
    tresc                  TEXT    NOT NULL DEFAULT '',
    odcisk_tresci          TEXT    NOT NULL DEFAULT '',
    aktywny                INTEGER NOT NULL DEFAULT 1 CHECK(aktywny IN (0,1)),
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    CHECK((os = 'platform') = (os_byt = '')),
    UNIQUE(kategoria_id, os, os_byt)
);
CREATE INDEX idx_dokument_tozsamosci_os ON dokument_tozsamosci(os, os_byt, aktywny);
CREATE INDEX idx_dokument_tozsamosci_kategoria ON dokument_tozsamosci(kategoria_id, os);

-- Katalog niesie wyłącznie metadane kategorii — nazwę, opis przeznaczenia,
-- warstwę, porządek i obowiązkowość; treść promptu pochodzi wyłącznie
-- z okna konfiguracji, nigdy z kodu ani ze schematu.
WITH katalog(kod, nazwa, opis, warstwa, kolejnosc, obowiazkowa) AS (
    VALUES
        -- Warstwa konstytucji — zasady nadrzędne; stoją najwyżej.
        ('konstytucja',           'Konstytucja',
         'Zasady nadrzędne tożsamości modelu; warstwa najwyższa nakładki. Jej odcisk niesie panel prowenancji',
         'constitution', 1, 1),
        ('bezpieczenstwo',        'Bezpieczeństwo',
         'Zasady bezpieczeństwa pracy modelu; kategoria wskazana wprost w katalogu kontraktu (security)',
         'constitution', 2, 0),
        ('poufnosc_i_sekrety',    'Poufność i sekrety',
         'Postępowanie z danymi poufnymi i danymi dostępowymi; sekret nigdy nie wchodzi do treści ani do bazy — wyłącznie odwołanie',
         'constitution', 3, 0),
        ('granica_uprawnien',     'Granica uprawnień',
         'Zakres czynności dozwolonych modelowi. Uprawnienia są konfiguracją możliwości, nie kontrolą dostępu; kontrolą dostępu jest uwierzytelnienie',
         'constitution', 4, 0),
        ('kompletnosc_rezultatu', 'Kompletność rezultatu',
         'Zakaz atrap i rezultatu pozornego: kod niepodłączony do punktu wejścia nie liczy się jako wykonany',
         'constitution', 5, 0),

        -- Warstwa profilu — kim model jest przy tej pracy i jak prowadzi narzędzia.
        ('profil_roli',           'Profil roli',
         'Rola, w jakiej model występuje w tym oknie; warstwa środkowa nakładki',
         'profile', 1, 0),
        ('harness',               'Harness',
         'Zasady pracy powłoki wykonawczej modelu wraz z hookami; wyłącznie z konfiguracji, nic zaszyte',
         'profile', 2, 0),
        ('narzedzia_modelu',      'Narzędzia modelu',
         'Zasady sterowania platformą przez narzędzia generowane z kontraktu; brak osobnego parsera intencji',
         'profile', 3, 0),
        ('rola_w_petli',          'Rola w pętli',
         'Zasady pracy koordynatora i wykonawcy w pętli na silniku kolejek: przekazanie zlecenia, wgląd w strumień, zatrzymanie biegu',
         'profile', 4, 0),
        ('katalog_roboczy',       'Katalog roboczy i dostępy',
         'Zasady korzystania z nadanych punktów dostępu oraz z katalogu roboczego sesji; dostęp i katalog roboczy są dwoma odrębnymi ustawieniami',
         'profile', 5, 0),

        -- Warstwa ekspertyzy — najniższa; zadanie, dziedzina, sposób odpowiadania.
        ('ekspertyza',            'Ekspertyza zadaniowa',
         'Wiedza i wytyczne dla bieżącego zadania; warstwa najniższa nakładki',
         'expertise', 1, 0),
        ('dziedzina_modulu',      'Dziedzina modułu',
         'Ekspertyza właściwa modułowi, w którym pracuje okno; kontekst poleceń paska promptu (docs/inwentarz/moduly.md)',
         'expertise', 2, 0),
        ('zachowanie_modelu',     'Zachowanie modelu',
         'Sposób prowadzenia rozmowy i redagowania odpowiedzi; zakres „Zachowanie modeli" okna konfiguracji (docs/inwentarz/okna-operacyjne.md, stany-i-menu.md)',
         'expertise', 3, 0)
)
INSERT INTO kategoria_tozsamosci (kod, nazwa, opis, warstwa, kolejnosc, obowiazkowa)
SELECT kod, nazwa, opis, warstwa, kolejnosc, obowiazkowa FROM katalog
WHERE true
ON CONFLICT(kod) DO NOTHING;
