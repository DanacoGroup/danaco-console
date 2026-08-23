-- Migracja 015 — tożsamość modelu: katalog kategorii i treść per oś.
--
-- Tożsamość modelu jest zamieniana, nie dołączana do ustawień fabrycznych.
-- Silnik nakładki (`server/internal/injection/nakladka.go`) niesie trzy warstwy
-- (konstytucja → profil → ekspertyza) oraz dwa tryby podania: TrybZastap
-- (`--system-prompt`) i TrybDopisz (`--append-system-prompt`). Te dwie tabele
-- niosą sterowanie silnikiem — skąd bierze się treść warstw i który tryb
-- obowiązuje — i ani jednego zdania promptu: treść wchodzi wyłącznie
-- z konfiguracji.
--
-- Wzorzec. Katalog naśladuje `akcja` i `kanal_modelu`: wiersz opisuje byt
-- w całości, `kod` jest jego trwałym identyfikatorem, `aktywna` rozstrzyga
-- widoczność, `kolejnosc` porządek składania. Nowa kategoria to nowy wiersz,
-- nie zmiana kodu.
--
-- Warstwa. Kolumna `warstwa` przyjmuje dosłownie wartości kontraktu
-- shared.IdentityLayer ('constitution','profile','expertise'), które kontrakt
-- sam wiąże ze stałymi silnika (injection.WarstwaKonstytucja / WarstwaProfil /
-- WarstwaEkspertyza). Kontrakt nie daje dla tego wyliczenia słownika przekładu
-- bazy — jedynym takim słownikiem nowych wyliczeń jest WartosciBazyAccountKind
-- dla kolumny `konto.rodzaj` — więc kolumna trzyma wartość kontraktu wprost,
-- bez drugiego, równoległego nazewnictwa.
--
-- Tryb. `tryb_domyslny` kategorii i `tryb` zapisu przyjmują wartości
-- shared.IdentityMode ('ZASTAP','DOLACZ'). Domyślną jest 'ZASTAP' — wartość
-- kolumny odpowiada wartości domyślnej klucza `tozsamosc.tryb_domyslny`. Baza
-- nie powtarza reguły składania trybu nakładki: rozstrzyga ją składacz
-- w `core/tozsamosc_prompt.go`, tak samo jak kolejność poziomów zasięgu
-- rozstrzyga `internal/konfig`, a nie schemat.
--
-- Oś. `os` przyjmuje wartości shared.ConfigAxis ('platform','model','account')
-- i mówi, dla czego treść obowiązuje; jest prostopadła do poziomu zasięgu.
-- `os_byt` wskazuje konkretny byt osi: identyfikator modelu dla osi 'model',
-- identyfikator konta (`konto.id`) dla osi 'account'. Dla osi 'platform'
-- `os_byt` jest pusty — pilnuje tego CHECK, bo pusty i niepusty byt to dwa
-- różne klucze rozstrzygania i pomyłka dawałaby dwa wiersze o tym samym
-- znaczeniu.
--
-- `os_byt` nie jest kluczem obcym: jedna kolumna niosłaby odwołania do dwóch
-- różnych rodziców — identyfikatora modelu, który w ogóle nie ma tabeli (model
-- jest napisem w `kanal_modelu.identyfikator_modelu` oraz w polu okna),
-- i identyfikatora konta z tabeli `konto`. Klucz obcy do dwóch tabel naraz nie
-- istnieje, a wybór jednej z nich zamykałby drogę drugiej osi. Spójność osi
-- 'account' sprawdza warstwa `dane` przy zapisie; brak bytu osi nie wywraca
-- odczytu, tylko znaczy „ten wiersz nie obowiązuje".
--
-- `odcisk_tresci` jest skrótem SHA-256 treści liczonym przez warstwę wyższą.
-- Służy diagnostyce prowenancji i rozpoznaniu, czy nakładka się zmieniła —
-- niczego nie dopuszcza i niczego nie blokuje. Kolumna pusta jest dopuszczalna.
--
-- Idempotencja. Zaczyn kończy się ON CONFLICT(kod) DO NOTHING, więc migracja
-- przechodzi także na bazie, w której część wierszy katalogu już jest.

-- ── Katalog kategorii zasad i tożsamości modelu ─────────────────────
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

-- Kolejność wewnątrz warstwy jest jednoznaczna. To nie jest ozdoba: składacz ma
-- dawać bajtowo ten sam prompt przy tej samej konfiguracji (pamięć podręczna
-- promptu), a dwie kategorie o tej samej kolejności w jednej warstwie czyniłyby
-- porządek składania zależnym od przypadku.
CREATE UNIQUE INDEX idx_kategoria_tozsamosci_porzadek
    ON kategoria_tozsamosci(warstwa, kolejnosc);

-- ── Treść kategorii per oś: platforma · model · konto ─────────────────────────
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

-- ── Zaczyn katalogu kategorii ─────────────────────────────────────────────────
--
-- Katalog niesie wyłącznie metadane kategorii — nazwę, opis przeznaczenia,
-- warstwę, porządek i obowiązkowość. Treść promptu nie pochodzi ani z kodu, ani
-- ze schematu: wnosi ją okno konfiguracji wierszem tabeli `dokument_tozsamosci`.
--
-- Obowiązkowa jest wyłącznie konstytucja — to jej odcisk niesie panel
-- prowenancji jako osobny blok „hash konstytucji". Obowiązkowość jest
-- oznaczeniem dla okna konfiguracji, nie bramą: brak treści nie wstrzymuje
-- uruchomienia, wraca wykazem `missingRequiredCategoryIds`.

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
