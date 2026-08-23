-- Migracja 075 — silnik mowa→tekst: wpisy katalogu ustawień
-- sterujące silnikiem oraz dziennik wykonanych transkrypcji.
--
-- Schemat daje silnikowi mowy dwie rzeczy, których w bazie nie ma: sterowanie
-- i ślad.
--
-- Silnik jest lokalny i CPU-only. Dźwięk nie wychodzi z maszyny Operatora:
-- rozpoznanie robi proces Pythona uruchamiany obok rdzenia, na tej samej
-- maszynie, bez wywołania sieciowego. Dlatego w tej migracji nie ma ani kolumny
-- na klucz API, ani na adres usługi, ani na konto dostawcy — nie ma dostawcy.
-- Nie ma też kolumny na wybór urządzenia liczącego: silnik chodzi na procesorze
-- i wiersz „karta graficzna tak/nie" byłby ustawieniem, którego nikt nie czyta.
--
-- Nagranie jest ścieżką, nie bajtami. `AudioRef` niesie ścieżkę pliku
-- na dysku Operatora; rdzeń bajtów nie kopiuje — otwiera plik w miejscu, tak
-- samo jak biblioteka trzyma treść pliku na dysku i w kolumnie `tresc_odwolanie`
-- wyłącznie odwołanie. Kolumny BLOB w tej migracji nie ma i mieć
-- nie będzie: nagranie głosowe wniesione do bazy urosłoby o rząd wielkości
-- ponad wszystko inne w tym pliku, a kopia dźwięku obok oryginału byłaby drugą
-- prawdą o tym samym nagraniu.
--
-- Czas jest liczbą, a zegar jest jeden. Kolumna `utworzono` niesie
-- milisekundy epoki i wartość podaje wołający, a nie baza wyrażeniem
-- `strftime`. Tak samo robi `rozszerzenie.zaktualizowano`
-- i `metoda_uwierzytelnienia.utworzono`. Baza z własnym „teraz"
-- byłaby drugim zegarem obok zegara rdzenia i dwa ślady tego samego zlecenia
-- rozjeżdżałyby się o czas zapisu.
--
-- (a) Sterowanie silnikiem: wpisy katalogu, a nie stałe kodu
--
-- Cztery wpisy katalogu zamiast czterech stałych w Go. Gdyby ścieżka
-- interpretera, rozmiar modelu, język i katalog pobrania siedziały w kodzie,
-- zmiana każdego z nich byłaby wydaniem binarium. Ważniejsze jest jednak coś
-- innego: wpis katalogu jest sterowalny od zaraz. Komendy `config.set`,
-- `config.get` i `settings.definition.list` są już wpięte w rdzeń i czytają
-- tabele katalogu ustawień, więc Operator ustawia silnik, zanim powstanie
-- jakakolwiek nowa rodzina komend mowy. Nowa pozycja okna konfiguracji to nowy
-- wiersz, nie nowy kod.
--
-- Kategoria jest nowa, nie dopisana do `modele`. Kategoria `modele` opisuje
-- kanał modelu, czyli to, co odpowiada Operatorowi (kanał główny, kanał
-- zapasowy, nakład rozumowania). Silnik mowy niczego nie odpowiada; jest
-- przetwornikiem wejścia stojącym przed rozmową
-- i nie ma z kanałem modelu wspólnego ani wyboru dostawcy, ani nakładu, ani
-- konta. Wsunięcie go do `modele` zmieszałoby dwa różne byty w jednym oknie
-- konfiguracji.
--
-- Oś jest jedna — `platform`. Osie `model` i `account` przysługują pozycjom
-- opisującym sposób uruchomienia modelu. Rozmiar modelu mowy ani
-- ścieżka Pythona nie zależą od tego, który model odpowiada ani na czyim koncie:
-- to własność maszyny Operatora. Oś modelu byłaby tu pytaniem bez treści.
--
-- Zasięgi są dwojakie: dwa wpisy opisują instalację, a dwa — zlecenie.
-- `mowa_program` i `mowa_katalog_modeli` opisują jedną instalację na maszynie:
-- interpreter jest jeden i pobrane modele leżą w jednym miejscu, więc te dwa
-- wpisy stoją wyłącznie na poziomie `globalny` (tak samo jak
-- `katalog.roboczy.wzorzec_sesji`). `mowa_model` i `mowa_jezyk`
-- opisują pojedyncze zlecenie — jedno okno dyktuje po polsku modelem `small`,
-- drugie przepisuje nagranie obcojęzyczne modelem `large-v3` — więc są
-- ustawialne na wszystkich ośmiu poziomach.
--
-- Wartość pusta coś znaczy. Pusty `mowa_program` znaczy „szukaj
-- `python3` na ścieżce wyszukiwania systemu", a pusty `mowa_katalog_modeli` —
-- „katalog domyślny biblioteki". Żadna z tych pustek nie jest brakiem danych ani
-- powodem odmowy uruchomienia; obie są rozstrzygnięciem „decyduje warstwa
-- niżej". Dlatego nie ma tu wpisów `wymagane`.
--
-- Idempotencja. Wstawienia kończą się klauzulą ON CONFLICT DO NOTHING, więc
-- krok przechodzi także na bazie, w której część
-- wierszy już stoi. Klauzula `WHERE true` przed ON CONFLICT jest wymogiem
-- składni SQLite dla INSERT ... SELECT z upsertem.
--
-- (b) Dziennik transkrypcji
--
-- Odmowy zapisujemy razem z powodzeniami. Kolumna `stan` niesie dwie wartości,
-- a nie jedną: transkrypcja, która się nie udała, jest zdarzeniem, o które
-- Operator zapyta jako pierwsze („dlaczego nic nie przepisało"). Odmowa bez
-- śladu jest nie do zdiagnozowania — zostaje pusty ekran i żadnej odpowiedzi na
-- pytanie, czy zabrakło Pythona, silnika, czy nagrania. Kolumna `powod` niesie
-- nazwany powód i ma sens wyłącznie przy odmowie; wiąże je CHECK tabelowy,
-- żeby „gotowa z powodem" i „odmowa bez powodu" były niemożliwe w bazie,
-- a nie tylko niezalecane w kodzie.
--
-- `trwanie_ms` to długość nagrania, nie czas przetwarzania. Nazwa myli się
-- łatwo i dlatego stoi tu wprost: kolumna mówi, ile trwa dźwięk, a nie ile
-- sekund liczył procesor. Czasu przetwarzania ta migracja nie zapisuje — byłby
-- miarą maszyny i chwili, a nie faktem o nagraniu, i nikt w rdzeniu go nie
-- czyta.
--
-- `okno_id` bywa puste i jest napisem, nie kluczem obcym. Transkrypcję wołają
-- także ścieżki spoza okna komunikacji (kolejka wykonawcy asystenta), więc
-- kolumna dopuszcza NULL. Klucza obcego nie ma z tego samego powodu, dla
-- którego nie ma go `sesja_bramki.urzadzenie_kod`: odmowa zapisu
-- śladu z powodu nieznanego okna kasowałaby dowód zdarzenia, które się wydarzyło.

-- ── Kategoria okna konfiguracji ──────────────────────────────────────────────
INSERT INTO kategoria_ustawien (kod, nazwa, opis, ikona, kolejnosc) VALUES
    ('mowa', 'Mowa na tekst', 'Lokalny silnik rozpoznawania mowy: interpreter, rozmiar modelu, język i katalog modeli', 'mikrofon', 9)
ON CONFLICT(kod) DO NOTHING;

-- ── Cztery pozycje katalogu ──────────────────────────────────────────────────
-- Rodzaje wartości dobrane z wyliczenia SettingValueType kontraktu (kolumna
-- `definicja_ustawienia.rodzaj_wartosci`): dwie ścieżki są rodzaju `path` — tak
-- samo jak `harness.program_claude` — kod języka jest zwykłym
-- napisem `string`, a rozmiar modelu wyliczeniem `enum` z opcjami niżej.
WITH katalog(klucz, kategoria, nazwa, opis, rodzaj, domyslna, podpowiedz,
             wymaga_restartu, kolejnosc) AS (
    VALUES
        ('mowa_program', 'mowa', 'Interpreter Pythona',
         'Ścieżka interpretera Pythona uruchamiającego pomocnika transkrypcji. Pusta znaczy: szukaj `python3` na ścieżce wyszukiwania systemu, a nie odmawiaj uruchomienia.',
         'path', '', 'python3', 1, 1),
        ('mowa_model', 'mowa', 'Rozmiar modelu',
         'Rozmiar modelu rozpoznawania — od najszybszego do najdokładniejszego. Model większy przepisuje wierniej i wolniej; silnik liczy na procesorze, więc wybór jest wprost wymianą czasu na dokładność.',
         'enum', 'small', '', 0, 2),
        ('mowa_jezyk', 'mowa', 'Język rozpoznawania',
         'Kod języka, w którym rozpoznawana jest mowa. Wskazanie języka z góry jest szybsze i wierniejsze niż zdanie się na rozpoznanie automatyczne.',
         'string', 'pl', 'pl', 0, 3),
        ('mowa_katalog_modeli', 'mowa', 'Katalog modeli',
         'Katalog, do którego pobierane są modele rozpoznawania. Pusty znaczy: katalog domyślny biblioteki silnika. Modele bywają wielkie, więc Operator może je przenieść na dysk pojemniejszy.',
         'path', '', '', 1, 4)
)
INSERT INTO definicja_ustawienia (klucz, kategoria_id, nazwa, opis, rodzaj_wartosci,
                                  wartosc_domyslna, podpowiedz, wymaga_restartu, kolejnosc)
SELECT k.klucz, kat.id, k.nazwa, k.opis, k.rodzaj, k.domyslna, k.podpowiedz,
       k.wymaga_restartu, k.kolejnosc
  FROM katalog k
  JOIN kategoria_ustawien kat ON kat.kod = k.kategoria
 WHERE true
ON CONFLICT(klucz) DO NOTHING;

-- ── Dopuszczalne rozmiary modelu ─────────────────────────────────────────────
-- Pięć rozmiarów wydawanych przez bibliotekę silnika, w kolejności rosnącej
-- wierności. Wartości pustej nie ma — inaczej niż przy `naklad_rozumowania`,
-- bo pod spodem nie ma warstwy, która by rozmiar rozstrzygnęła:
-- silnik musi dostać nazwę modelu, więc „bez wskazania" znaczyłoby tylko
-- „odmowa uruchomienia", a to nie jest opcja katalogu.
WITH opcje(wartosc, etykieta, opis, kolejnosc) AS (
    VALUES
        ('tiny',     'Najmniejszy', 'Najszybciej, najmniej dokładnie',            1),
        ('base',     'Mały',        'Szybciej, kosztem dokładności',              2),
        ('small',    'Średni',      'Równowaga szybkości i dokładności',          3),
        ('medium',   'Duży',        'Dokładniej, wolniej',                        4),
        ('large-v3', 'Największy',  'Najdokładniej, najwolniej i najwięcej pamięci', 5)
)
INSERT INTO opcja_ustawienia (definicja_id, wartosc, etykieta, opis, kolejnosc)
SELECT d.id, o.wartosc, o.etykieta, o.opis, o.kolejnosc
  FROM opcje o
  JOIN definicja_ustawienia d ON d.klucz = 'mowa_model'
 WHERE true
ON CONFLICT(definicja_id, wartosc) DO NOTHING;

-- ── Dopuszczalne poziomy zasięgu ─────────────────────────────────────────────
-- Instalacja: jeden interpreter i jeden katalog modeli na maszynę.
INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
 CROSS JOIN poziom_zasiegu p
 WHERE d.klucz IN ('mowa_program', 'mowa_katalog_modeli')
   AND p.kod = 'globalny'
ON CONFLICT DO NOTHING;

-- Zlecenie: rozmiar modelu i język bywają inne w każdym oknie.
INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
 CROSS JOIN poziom_zasiegu p
 WHERE d.klucz IN ('mowa_model', 'mowa_jezyk')
ON CONFLICT DO NOTHING;

-- ── Dopuszczalne osie ────────────────────────────────────────────────────────
-- Wyłącznie `platform` — patrz nagłówek.
INSERT INTO definicja_ustawienia_os (definicja_id, os)
SELECT d.id, 'platform'
  FROM definicja_ustawienia d
 WHERE d.klucz IN ('mowa_program', 'mowa_model', 'mowa_jezyk', 'mowa_katalog_modeli')
ON CONFLICT DO NOTHING;

-- ── Dziennik wykonanych transkrypcji ─────────────────────────────────────────
CREATE TABLE transkrypcja (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Tożsamość wpisu widoczna na zewnątrz, tak samo jak
    -- `rozszerzenie.identyfikator_zewnetrzny`. Klucz sztuczny
    -- zostaje sprawą bazy i nigdy nie wychodzi w odpowiedzi.
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- Okno, w którym padło zlecenie. NULL znaczy „spoza okna" — patrz nagłówek.
    okno_id                  TEXT,
    -- Ścieżka pliku na dysku Operatora, nie bajty dźwięku.
    nagranie_odnosnik        TEXT    NOT NULL,
    -- Rozmiar modelu i język zapisane tak, jak obowiązywały w chwili zlecenia.
    -- Nie są odczytem bieżącego ustawienia: Operator zmieni je jutro, a ślad ma
    -- mówić, czym przepisano to nagranie, a nie czym przepisano by je dziś.
    model                    TEXT    NOT NULL,
    jezyk                    TEXT    NOT NULL,
    znakow                   INTEGER NOT NULL DEFAULT 0,
    -- Długość nagrania w milisekundach, nie czas przetwarzania — patrz nagłówek.
    trwanie_ms               INTEGER NOT NULL DEFAULT 0,
    -- Trzy stany, nie dwa. Poza 'gotowa' i 'odmowa' jest 'bez_mowy': nagranie
    -- ciszy albo szumu, gdzie przetworzenie się odbyło i zmierzyło, że mowy
    -- nie ma. Zapisanie go jako 'gotowa' ze znakow=0 kazałoby czytającemu
    -- zgadywać z liczby; zapisanie jako 'odmowa' byłoby nieprawdą, bo niczego
    -- nie odmówiono. Stan nazywa więc rzecz wprost.
    stan                     TEXT    NOT NULL CHECK(stan IN ('gotowa', 'bez_mowy', 'odmowa')),
    -- Nazwany powód odmowy; przy stanach udanych nie ma czego opisywać.
    powod                    TEXT,
    -- Milisekundy epoki podane przez wołającego — patrz nagłówek.
    utworzono                INTEGER NOT NULL,
    -- Wiąże stan z powodem w obie strony: odmowa MUSI nieść powód, a wpis
    -- gotowy nie ma prawa go nieść. Bez tego więzu w tabeli stanęłaby odmowa
    -- bez wyjaśnienia — czyli dokładnie ten wiersz, dla którego dziennik
    -- powstaje.
    CHECK((stan = 'odmowa' AND powod IS NOT NULL AND powod <> '')
          OR (stan IN ('gotowa', 'bez_mowy') AND powod IS NULL))
);

-- Odpowiedź na pytanie „co przepisano w tym oknie, najnowsze najpierw".
-- Kolumny biegną dokładnie porządkiem zapytania wykazu z `mowa/dziennik.go`
-- (`okno_id` zawęża, `utworzono, id` porządkują), więc przy wskazanym oknie
-- SQLite czyta indeks wstecz zamiast sortować wynik. Kolumny `stan` w indeksie
-- nie ma celowo: stanęłaby między zawężeniem a porządkiem i zepsułaby ten drugi,
-- a wykaz pokazuje odmowy razem z powodzeniami — po to je zapisujemy.
CREATE INDEX idx_transkrypcja_wykaz ON transkrypcja(okno_id, utworzono, id);
