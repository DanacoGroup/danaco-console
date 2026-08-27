-- Migracja 012 przenosi katalog pozycji okna konfiguracji do danych i wprowadza
-- oś zasięgu ConfigAxis obok ośmiu istniejących poziomów zasięgu.

-- Domyślna oś 'platform' obowiązuje przy braku wskazania osi.

-- Wartości kolumny `kod` odpowiadają wyliczeniu ConfigAxis kontraktu.

CREATE TABLE os_zasiegu (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                    TEXT    NOT NULL UNIQUE
                                   CHECK(kod IN ('platform','model','account')),
    nazwa                  TEXT    NOT NULL,
    pierwszenstwo          INTEGER NOT NULL UNIQUE
);

INSERT INTO os_zasiegu (kod, nazwa, pierwszenstwo) VALUES
    ('platform', 'Platforma', 1),
    ('model',    'Model',     2),
    ('account',  'Konto',     3);

-- ── Przebudowa tabeli `ustawienie` o oś rozstrzygania ────────────────────────

CREATE TABLE ustawienie_z_osia (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    poziom_zasiegu_id      INTEGER NOT NULL REFERENCES poziom_zasiegu(id) ON DELETE CASCADE,
    klucz_zasiegu          TEXT    NOT NULL DEFAULT '',
    -- Oś wskazywana kodem `os_zasiegu.kod`; wartość domyślna 'platform' znaczy brak wskazania.
    os                     TEXT    NOT NULL DEFAULT 'platform'
                                   REFERENCES os_zasiegu(kod) ON UPDATE CASCADE,
    klucz_osi              TEXT    NOT NULL DEFAULT '',
    klucz                  TEXT    NOT NULL,
    wartosc                TEXT,
    rodzaj_wartosci        TEXT    NOT NULL DEFAULT 'tekst'
                                   CHECK(rodzaj_wartosci IN ('tekst','liczba','logiczna','json')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(poziom_zasiegu_id, klucz_zasiegu, os, klucz_osi, klucz)
);

INSERT INTO ustawienie_z_osia
    (id, poziom_zasiegu_id, klucz_zasiegu, os, klucz_osi, klucz, wartosc,
     rodzaj_wartosci, zaktualizowano)
SELECT u.id, u.poziom_zasiegu_id, u.klucz_zasiegu, 'platform', '',
       u.klucz, u.wartosc, u.rodzaj_wartosci, u.zaktualizowano
  FROM ustawienie u;

DROP TABLE ustawienie;
ALTER TABLE ustawienie_z_osia RENAME TO ustawienie;

CREATE INDEX idx_ustawienie_klucz ON ustawienie(klucz);
CREATE INDEX idx_ustawienie_os ON ustawienie(os, klucz_osi, klucz);

-- Kategoria może mieć rodzica: drzewo okna konfiguracji rośnie wierszami, nie kodem.

CREATE TABLE kategoria_ustawien (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kod                    TEXT    NOT NULL UNIQUE,
    nazwa                  TEXT    NOT NULL,
    opis                   TEXT    NOT NULL DEFAULT '',
    ikona                  TEXT    NOT NULL DEFAULT '',
    kategoria_nadrzedna_id INTEGER REFERENCES kategoria_ustawien(id) ON DELETE SET NULL,
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    aktywna                INTEGER NOT NULL DEFAULT 1 CHECK(aktywna IN (0,1)),
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Brak wiersza w katalogu nie jest awarią: rezolwer schodzi na rejestr wbudowany rdzenia.

CREATE TABLE definicja_ustawienia (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    klucz                  TEXT    NOT NULL UNIQUE,
    kategoria_id           INTEGER NOT NULL REFERENCES kategoria_ustawien(id) ON DELETE CASCADE,
    nazwa                  TEXT    NOT NULL,
    opis                   TEXT    NOT NULL DEFAULT '',
    rodzaj_wartosci        TEXT    NOT NULL
                                   CHECK(rodzaj_wartosci IN ('string','text','int','float','bool',
                                                             'enum','enumList','path','pathList',
                                                             'secret','json')),
    wartosc_domyslna       TEXT    NOT NULL DEFAULT '',
    minimum                REAL,
    maksimum               REAL,
    skok                   REAL,
    wzorzec                TEXT    NOT NULL DEFAULT '',
    jednostka              TEXT    NOT NULL DEFAULT '',
    podpowiedz             TEXT    NOT NULL DEFAULT '',
    wymagane               INTEGER NOT NULL DEFAULT 0 CHECK(wymagane IN (0,1)),
    wymaga_restartu        INTEGER NOT NULL DEFAULT 0 CHECK(wymaga_restartu IN (0,1)),
    widoczne_gdy_klucz     TEXT    NOT NULL DEFAULT '',
    widoczne_gdy_wartosc   TEXT    NOT NULL DEFAULT '',
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    aktywna                INTEGER NOT NULL DEFAULT 1 CHECK(aktywna IN (0,1)),
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_definicja_ustawienia_kategoria ON definicja_ustawienia(kategoria_id, kolejnosc);

-- Dopuszczalne wartości pozycji o rodzaju enum albo enumList (SettingOption).

CREATE TABLE opcja_ustawienia (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    definicja_id           INTEGER NOT NULL REFERENCES definicja_ustawienia(id) ON DELETE CASCADE,
    wartosc                TEXT    NOT NULL,
    etykieta               TEXT    NOT NULL,
    opis                   TEXT    NOT NULL DEFAULT '',
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    UNIQUE(definicja_id, wartosc)
);

-- Poziomy zasięgu, na których wolno zapisać pozycję (SettingDefinition.allowedScopes).

CREATE TABLE definicja_ustawienia_zasieg (
    definicja_id           INTEGER NOT NULL REFERENCES definicja_ustawienia(id) ON DELETE CASCADE,
    poziom_zasiegu_id      INTEGER NOT NULL REFERENCES poziom_zasiegu(id) ON DELETE CASCADE,
    PRIMARY KEY (definicja_id, poziom_zasiegu_id)
);

-- Osie, dla których wolno zapisać pozycję (SettingDefinition.allowedAxes).

CREATE TABLE definicja_ustawienia_os (
    definicja_id           INTEGER NOT NULL REFERENCES definicja_ustawienia(id) ON DELETE CASCADE,
    os                     TEXT    NOT NULL REFERENCES os_zasiegu(kod) ON UPDATE CASCADE,
    PRIMARY KEY (definicja_id, os)
);

-- Osiem kategorii katalogu ustawień odpowiada bytom czytanym przez rdzeń.

INSERT INTO kategoria_ustawien (kod, nazwa, opis, ikona, kolejnosc) VALUES
    ('modele',          'Ustawienia modeli', 'Kanał modelu, kanał zapasowy i nakład rozumowania',                                    'siec',       1),
    ('harness',         'Harness',           'Powłoka wykonawcza modelu: program, plik ustawień, mosty MCP, próg biegu naprawczego', 'terminal',   2),
    ('tozsamosc',       'Tożsamość modelu',  'Tryb podania tożsamości: zastąpienie promptu fabrycznego albo dopisanie do niego',     'tarcza',     3),
    ('bezpieczenstwo',  'Security',          'Zakres zgody wydanej modelowi oraz egzekwowanie uwierzytelniania',                     'klodka',     4),
    ('katalog_roboczy', 'Katalog roboczy',   'Miejsce, w którym powstają katalogi sesyjne i pliki robocze modelu',                   'folder',     5),
    ('wykonanie',       'Wykonanie okna',    'Rola okna w pętli koordynator–wykonawca oraz zasięg i host wykonania',                 'karta-okna', 6),
    ('izolacja',        'Izolacja zasięgu',  'Jedenaście punktów izolacji: trzy wymiary kontekstu i osiem zakresów technicznych',    'warstwy',    7),
    ('personalizacja',  'Personalizacja',    'Wygląd interfejsu Operatora',                                                          'paleta',     8)
ON CONFLICT(kod) DO NOTHING;

-- ── Pozycje katalogu ─────────────────────────────────────────────────────────

WITH katalog(klucz, kategoria, nazwa, opis, rodzaj, domyslna, minimum, skok,
             podpowiedz, wymaga_restartu, kolejnosc) AS (
    VALUES
        -- Osiem parametrów wykonania — definicje_wykonania.go.
        ('kanal_modelu', 'modele', 'Kanał modelu', 'Kanał modelu z rejestru kanałów. Brak wskazania znaczy: rozstrzyga rejestr kanałów, nie odmowa uruchomienia.', 'string', '', NULL, NULL, 'kod kanału z rejestru', 0, 1),
        ('kanal_modelu_zapasowy', 'modele', 'Kanał zapasowy', 'Kanał modelu używany po niepowodzeniu kanału głównego. Brak wskazania znaczy brak zapasu, nie blokadę wywołania.', 'string', '', NULL, NULL, 'kod kanału z rejestru', 0, 2),
        ('naklad_rozumowania', 'modele', 'Nakład rozumowania', 'Nakład rozumowania modelu — od szybciej do mądrzej. Brak wskazania zostawia rozstrzygnięcie kanałowi modelu.', 'enum', '', NULL, NULL, '', 0, 3),
        ('rola_okna', 'wykonanie', 'Rola okna', 'Rola okna w pętli koordynator–wykonawca. Okno samodzielne pozostaje poza pętlą.', 'enum', 'standalone', NULL, NULL, '', 0, 1),
        ('srodowisko_wykonania', 'wykonanie', 'Zasięg wykonania', 'Zasięg wykonania modelu: urządzenie użytkownika, host rdzenia albo host zdalny. Niezależny od umiejscowienia rdzenia.', 'enum', 'local', NULL, NULL, '', 0, 2),
        ('host_wykonania', 'wykonanie', 'Host wykonania', 'Nazwa hosta wykonania przy zasięgu zdalnym — na przykład danaco-system. Nazwa hosta jest ustawieniem poziomu zasięgu, nie wartością wyliczenia.', 'string', '', NULL, NULL, 'danaco-system', 0, 3),
        ('tryb_uprawnien', 'bezpieczenstwo', 'Tryb uprawnień', 'Tryb uprawnień okna komunikacji — zakres zgody wydanej modelowi przed zmianą w systemie. Wartości kontraktu PermissionMode.', 'enum', 'manual', NULL, NULL, '', 0, 1),
        ('katalogi_robocze', 'katalog_roboczy', 'Katalogi robocze okna', 'Lista katalogów roboczych okna komunikacji. Lista pusta znaczy: katalog wskaże Operator przy otwarciu okna.', 'pathList', '[]', NULL, NULL, '', 0, 1),

        -- 2. Jedenaście punktów izolacji — definicje_izolacji.go, co do znaku.
        ('izolacja_historia', 'izolacja', 'Historia wymiany', 'Zapis wymiany wiadomości. Odrębna: nowe okno zaczyna z pustą historią. Współdzielona: ten sam zapis widoczny w kilku oknach lub modułach.', 'enum', 'odrebna', NULL, NULL, '', 0, 1),
        ('izolacja_pamiec', 'izolacja', 'Pamięć długoterminowa', 'Pamięć długoterminowa zasięgu. Odrębna: pamięć jednego zasięgu niewidoczna w innym. Współdzielona: jeden zasób zasila kilka zasięgów.', 'enum', 'odrebna', NULL, NULL, '', 0, 2),
        ('izolacja_kontekst', 'izolacja', 'Kontekst roboczy', 'Bieżący stan roboczy: aktywne pliki, projekt, załączniki, zmienne. Współdzielony przenosi się między oknami bez przeładowania.', 'enum', 'odrebna', NULL, NULL, '', 0, 3),
        ('izolacja_katalog_roboczy_sesji', 'izolacja', 'Katalog roboczy sesji', 'Fizyczny katalog plików procesu. Włączony: własny katalog, niewidoczny dla innych sesji zasięgu. Stan wyjściowy platformy: wyłączony.', 'enum', 'wylaczony', NULL, NULL, '', 0, 4),
        ('izolacja_srodowisko_procesu', 'izolacja', 'Środowisko procesu', 'Zmienne środowiskowe i kontekst uruchomieniowy. Włączony: własny zestaw zmiennych zamiast wspólnego środowiska serwera. Stan wyjściowy platformy: wyłączony.', 'enum', 'wylaczony', NULL, NULL, '', 0, 5),
        ('izolacja_katalog_danych_modelu', 'izolacja', 'Katalog danych modelu', 'Dane pomocnicze kanału modelu: konfiguracja, dane tymczasowe, ustawienia dostawcy. Włączony: własny katalog danych modelu. Stan wyjściowy platformy: wyłączony.', 'enum', 'wylaczony', NULL, NULL, '', 0, 6),
        ('izolacja_dostep_sieciowy', 'izolacja', 'Dostęp sieciowy', 'Połączenia wychodzące: API, strony, serwery MCP, rozszerzenia. Włączony: odrębny, ograniczony dostęp sieciowy. Stan wyjściowy platformy: wyłączony.', 'enum', 'wylaczony', NULL, NULL, '', 0, 7),
        ('izolacja_odczyt_zapis_plikow', 'izolacja', 'Odczyt i zapis plików', 'Uprawnienia do plików poza własnym katalogiem roboczym. Włączony: dostęp wyłącznie do ścieżek jawnie dozwolonych. Stan wyjściowy platformy: wyłączony.', 'enum', 'wylaczony', NULL, NULL, '', 0, 8),
        ('izolacja_konto_i_token', 'izolacja', 'Konto i token', 'Token dostępu, klucz API, dane logowania kanału modelu. Włączony: własne dane dostępowe zamiast platformowych (odwołania do magazynu sekretów, nie same sekrety). Stan wyjściowy platformy: wyłączony.', 'enum', 'wylaczony', NULL, NULL, '', 0, 9),
        ('izolacja_model_procesu', 'izolacja', 'Instancja procesu modelu', 'Instancja procesu wykonawczego modelu. Włączony: własna, niezależna instancja zamiast wspólnej puli. Stan wyjściowy platformy: wyłączony.', 'enum', 'wylaczony', NULL, NULL, '', 0, 10),
        ('izolacja_serwer_wykonania', 'izolacja', 'Serwer wykonania', 'Serwer wykonania procesu — istotne przy kanale zdalnym. Włączony: serwer dedykowany zamiast współdzielonego. Stan wyjściowy platformy: wyłączony.', 'enum', 'wylaczony', NULL, NULL, '', 0, 11),

        -- 3. Klucze katalogu roboczego i tożsamości modelu.
        ('katalog.roboczy.podstawa', 'katalog_roboczy', 'Podstawa katalogu roboczego', 'Katalog, w którym powstają katalogi sesyjne i pliki robocze modelu. Pusty znaczy: miejsce instalacji aplikacji głównej. Dostęp modelu do maszyn i katalogów jest ustawieniem osobnym — nadaniem dostępu, nie tym kluczem.', 'path', '', NULL, NULL, 'miejsce instalacji aplikacji głównej', 0, 2),
        ('katalog.roboczy.wzorzec_sesji', 'katalog_roboczy', 'Wzorzec katalogu sesji', 'Wzorzec nazwy katalogu sesji powstającego wewnątrz podstawy katalogu roboczego. Znacznik identyfikatora zastępowany jest identyfikatorem sesji.', 'string', 'sesje/<identyfikator>', NULL, NULL, 'sesje/<identyfikator>', 0, 3),
        ('tozsamosc.tryb_domyslny', 'tozsamosc', 'Tryb tożsamości', 'Tryb podania tożsamości modelu. ZASTAP podmienia prompt fabryczny w całości, DOLACZ dokłada nakładkę do promptu programu.', 'enum', 'ZASTAP', NULL, NULL, '', 0, 1),

        -- 4. Pola powłoki wykonawczej — injection/ustawienia.go.
        ('harness.program_claude', 'harness', 'Program kanału głównego', 'Ścieżka pliku wykonywalnego kanału głównego. Wskazanie w wierszu rejestru kanałów jest węższe i wygrywa; ta wartość obowiązuje, gdy wiersz rejestru ścieżki nie niesie.', 'path', 'claude', NULL, NULL, 'claude', 1, 1),
        ('harness.plik_ustawien', 'harness', 'Plik ustawień powłoki', 'Ścieżka pliku albo treść JSON podawana powłoce przełącznikiem --settings. Pusta znaczy: przełącznika nie podajemy.', 'path', '', NULL, NULL, 'settings.json', 1, 2),
        ('harness.konfiguracja_mcp', 'harness', 'Konfiguracja mostów MCP', 'Wykaz plików albo treści JSON podawanych powłoce przełącznikiem --mcp-config. Lista pusta znaczy: przełącznika nie podajemy.', 'pathList', '[]', NULL, NULL, '', 1, 3),
        ('petla.prog_braku_postepu', 'harness', 'Próg braku postępu', 'Ile obiegów bez zmiany stanu z rzędu kończy bieg naprawczy. Bieg nie ma limitu obiegów — próg dotyczy wyłącznie powtarzalności wykonawcy, a zatrzymanie zawsze niesie nazwany powód.', 'int', '3', 1, 1, '', 0, 4),

        -- 5. Egzekwowanie uwierzytelniania oraz warstwa motywu klienta.
        ('bezpieczenstwo.egzekwowanie_uwierzytelniania', 'bezpieczenstwo', 'Egzekwowanie uwierzytelniania', 'Czy uwierzytelnianie jest egzekwowane. W fazie budowy wyłączone; egzekwowanie włącza Operator.', 'bool', 'false', NULL, NULL, '', 1, 2),
        ('personalizacja.motyw', 'personalizacja', 'Motyw interfejsu', 'Wybór motywu jasnego albo ciemnego. Brak wyboru zdaje się na preferencję systemu — oba motywy są równoprawne, żaden nie jest wartością domyślną produktu.', 'enum', '', NULL, NULL, '', 0, 1)
)
INSERT INTO definicja_ustawienia (klucz, kategoria_id, nazwa, opis, rodzaj_wartosci,
                                  wartosc_domyslna, minimum, skok, podpowiedz,
                                  wymaga_restartu, kolejnosc)
SELECT k.klucz, kat.id, k.nazwa, k.opis, k.rodzaj, k.domyslna, k.minimum, k.skok,
       k.podpowiedz, k.wymaga_restartu, k.kolejnosc
  FROM katalog k
  JOIN kategoria_ustawien kat ON kat.kod = k.kategoria
 WHERE true
ON CONFLICT(klucz) DO NOTHING;

-- Wartość pusta jest pełnoprawną opcją: znaczy, że rozstrzyga warstwa niżej.

WITH opcje(klucz, wartosc, etykieta, opis, kolejnosc) AS (
    VALUES
        ('tryb_uprawnien', 'manual',            'Ręczny',                'Pytanie o zgodę przed każdą zmianą',                  1),
        ('tryb_uprawnien', 'acceptEdits',       'Zgoda na zmiany plików','Automatyczna zgoda na zmiany plików',                 2),
        ('tryb_uprawnien', 'plan',              'Planowanie',            'Praca planistyczna bez zmian w systemie',             3),
        ('tryb_uprawnien', 'auto',              'Model rozstrzyga',      'Decyzje o uprawnieniach podejmuje model',             4),
        ('tryb_uprawnien', 'dontAsk',           'Bez zapytań',           'Bez zapytań, z zachowaniem ograniczeń',               5),
        ('tryb_uprawnien', 'bypassPermissions', 'Pominięcie kontroli',   'Pominięcie kontroli uprawnień',                       6),

        ('rola_okna', 'standalone',  'Samodzielne', 'Okno poza pętlą koordynator–wykonawca',                 1),
        ('rola_okna', 'coordinator', 'Koordynator', 'Okno koordynatora; przekazuje zlecenie przez narzędzie', 2),
        ('rola_okna', 'executor',    'Wykonawca',   'Okno wykonawcy; zakończenie tury wybudza koordynatora',  3),

        ('srodowisko_wykonania', 'local',  'Urządzenie użytkownika', 'Wykonanie przez agenta lokalnego urządzenia', 1),
        ('srodowisko_wykonania', 'core',   'Host rdzenia',           'Wykonanie na hoście rdzenia serwera',         2),
        ('srodowisko_wykonania', 'remote', 'Host zdalny',            'Wykonanie na hoście wskazanym ustawieniem host_wykonania', 3),

        ('naklad_rozumowania', '',       'Bez wskazania', 'Rozstrzyga kanał modelu', 1),
        ('naklad_rozumowania', 'low',    'Niski',         'Szybciej',                2),
        ('naklad_rozumowania', 'medium', 'Średni',        'Równowaga',               3),
        ('naklad_rozumowania', 'high',   'Wysoki',        'Mądrzej',                 4),
        ('naklad_rozumowania', 'xhigh',  'Bardzo wysoki', 'Mądrzej, dłużej',         5),
        ('naklad_rozumowania', 'max',    'Najwyższy',     'Najwyższy nakład rozumowania', 6),

        ('tozsamosc.tryb_domyslny', 'ZASTAP', 'Zastąpienie', 'Prompt fabryczny zostaje podmieniony w całości', 1),
        ('tozsamosc.tryb_domyslny', 'DOLACZ', 'Dołączenie',  'Nakładka dokłada się do promptu programu',       2),

        ('personalizacja.motyw', '',      'Preferencja systemu', 'Brak wyboru — rozstrzyga preferencja systemu operacyjnego', 1),
        ('personalizacja.motyw', 'light', 'Jasny',               'Motyw jasny',                                              2),
        ('personalizacja.motyw', 'dark',  'Ciemny',              'Motyw ciemny',                                             3)
)
INSERT INTO opcja_ustawienia (definicja_id, wartosc, etykieta, opis, kolejnosc)
SELECT d.id, o.wartosc, o.etykieta, o.opis, o.kolejnosc
  FROM opcje o
  JOIN definicja_ustawienia d ON d.klucz = o.klucz
 WHERE true
ON CONFLICT(definicja_id, wartosc) DO NOTHING;

-- Trzy wymiary izolacji kontekstu: odrębna albo współdzielona.

WITH wymiary(wartosc, etykieta, opis, kolejnosc) AS (
    VALUES
        ('odrebna',       'Odrębna',       'Zasięg ma własny zasób',            1),
        ('wspoldzielona', 'Współdzielona', 'Jeden zasób zasila kilka zasięgów', 2)
)
INSERT INTO opcja_ustawienia (definicja_id, wartosc, etykieta, opis, kolejnosc)
SELECT d.id, w.wartosc, w.etykieta, w.opis, w.kolejnosc
  FROM definicja_ustawienia d
 CROSS JOIN wymiary w
 WHERE d.klucz IN ('izolacja_historia','izolacja_pamiec','izolacja_kontekst')
ON CONFLICT(definicja_id, wartosc) DO NOTHING;

-- Osiem zakresów izolacji technicznej: włączony albo wyłączony.

WITH zakresy(wartosc, etykieta, opis, kolejnosc) AS (
    VALUES
        ('wylaczony', 'Wyłączony', 'Zasób wspólny zasięgu — stan wyjściowy platformy', 1),
        ('wlaczony',  'Włączony',  'Zasób odrębny, niewidoczny dla innych sesji',      2)
)
INSERT INTO opcja_ustawienia (definicja_id, wartosc, etykieta, opis, kolejnosc)
SELECT d.id, w.wartosc, w.etykieta, w.opis, w.kolejnosc
  FROM definicja_ustawienia d
 CROSS JOIN zakresy w
 WHERE d.klucz LIKE 'izolacja\_%' ESCAPE '\'
   AND d.klucz NOT IN ('izolacja_historia','izolacja_pamiec','izolacja_kontekst')
ON CONFLICT(definicja_id, wartosc) DO NOTHING;

-- Pozycja jest ustawialna na wszystkich ośmiu poziomach zasięgu, poza wyjątkami niżej.

-- Wyjątki globalne: motyw interfejsu, egzekwowanie uwierzytelniania, wzorzec katalogu sesji.

INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
 CROSS JOIN poziom_zasiegu p
 WHERE d.klucz NOT IN ('personalizacja.motyw',
                       'bezpieczenstwo.egzekwowanie_uwierzytelniania',
                       'katalog.roboczy.wzorzec_sesji')
ON CONFLICT DO NOTHING;

INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
 CROSS JOIN poziom_zasiegu p
 WHERE d.klucz IN ('personalizacja.motyw',
                   'bezpieczenstwo.egzekwowanie_uwierzytelniania',
                   'katalog.roboczy.wzorzec_sesji')
   AND p.kod = 'globalny'
ON CONFLICT DO NOTHING;

-- Oś `platform` przysługuje każdej pozycji katalogu ustawień.

-- Osie `model` i `account` przysługują pozycjom opisu uruchomienia modelu.

INSERT INTO definicja_ustawienia_os (definicja_id, os)
SELECT d.id, o.kod
  FROM definicja_ustawienia d
 CROSS JOIN os_zasiegu o
 WHERE o.kod = 'platform'
ON CONFLICT DO NOTHING;

INSERT INTO definicja_ustawienia_os (definicja_id, os)
SELECT d.id, o.kod
  FROM definicja_ustawienia d
 CROSS JOIN os_zasiegu o
 WHERE o.kod IN ('model','account')
   AND d.klucz IN ('kanal_modelu', 'kanal_modelu_zapasowy', 'naklad_rozumowania',
                   'tryb_uprawnien', 'tozsamosc.tryb_domyslny',
                   'harness.program_claude', 'harness.plik_ustawien',
                   'harness.konfiguracja_mcp', 'petla.prog_braku_postepu',
                   'katalogi_robocze', 'katalog.roboczy.podstawa')
ON CONFLICT DO NOTHING;
