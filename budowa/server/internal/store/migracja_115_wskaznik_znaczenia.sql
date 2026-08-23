-- Wskaźnik znaczenia: fragmenty wiedzy wraz z ich osadzeniami oraz wpisy
-- katalogu ustawień sterujące silnikiem osadzeń.
--
-- Pełnotekstowy indeks treści biblioteki dopasowuje słowa, ten wskaźnik —
-- znaczenia: trzyma wektor, czyli ciąg liczb, w którym bliskość odpowiada
-- bliskości sensu, więc łączy zdania niemające wspólnego wyrazu. Zakresy obu
-- struktur są różne: indeks pełnotekstowy widzi wyłącznie pliki biblioteki,
-- wskaźnik obejmuje ponadto historię rozmów i pliki przestrzeni roboczej okna.
--
-- Fragment i wektor są bytem wtórnym, odtwarzalnym przebiegiem
-- `knowledge.index`; bajty treści leżą w magazynie biblioteki pod sumą sha256
-- (`core/adapter_modul_library_magazyn.go`), a `tresc_odwolanie` pozostaje
-- jedyną drogą do nich.
--
-- Wektor jest BLOB-em, nie tabelą współrzędnych: czyta się go zawsze w całości
-- i zawsze po to, żeby policzyć jeden iloczyn skalarny, a wiersz na współrzędną
-- dałby przy dziesięciu tysiącach fragmentów blisko osiem milionów wierszy,
-- o które nikt nie pyta pojedynczo. Zapis opisuje `wiedza/podobienstwo.go`:
-- liczby pojedynczej precyzji, porządek bajtów mało-końcowy, wektor podzielony
-- przez własną długość. Kolumna `wymiar` pozwala rozpoznać wiersz uszkodzony,
-- zanim trafi do porównania.
--
-- `model` stoi w wierszu i w warunku jednoznaczności, bo wektory dwóch modeli
-- leżą w różnych przestrzeniach i ich iloczyn skalarny nie jest trafnością.
-- Po zmianie ustawienia `wiedza_model` w tabeli stoją dwa komplety, a odczyt
-- zawęża się po modelu (`Pozycje`, `wiedza/skladnica.go`); wiersze poprzedniego
-- modelu zostają i czyści je `knowledge.index` z `rebuild`.
--
-- Kluczy obcych do źródła nie ma, bo jeden z trzech zakresów — plik przestrzeni
-- roboczej — nie ma wiersza w żadnej tabeli, a trzy wzajemnie wykluczające się
-- kolumny byłyby kształtem gorszym od jednego napisu. Fragmenty po skasowanym
-- źródle sprząta przebieg `knowledge.index`, który zna dziś istniejące źródła,
-- a nie kaskada bazy, która ich nie zna.
--
-- `utworzono` niesie milisekundy epoki podane przez wołającego, a nie wyrażenie
-- `strftime`.

CREATE TABLE fragment_wiedzy (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Zakres kontraktowy: 'library', 'history' albo 'workspace'. Wartość
    -- 'all' nie wchodzi do wiersza — jest pytaniem o sumę trzech zakresów,
    -- a nie czwartym miejscem, z którego coś pochodzi.
    zakres     TEXT    NOT NULL CHECK (zakres IN ('library', 'history', 'workspace')),
    -- Nazwa źródła dla człowieka: nazwa pliku biblioteki, tytuł okna rozmowy,
    -- ścieżka względna pliku przestrzeni. Wchodzi wprost do pola `source`
    -- trafienia — bez wskazania źródła model cytowałby bez możliwości
    -- sprawdzenia, a to jest dokładnie ten rodzaj odpowiedzi, którego kontrakt
    -- tej rodziny komend zakazuje.
    zrodlo     TEXT    NOT NULL,
    -- Identyfikator, którym da się po źródło sięgnąć (`sourceId`). Napis pusty
    -- jest stanem poprawnym dla źródeł bez identyfikatora — dlatego
    -- kolumna jest NOT NULL z wartością domyślną pustą, a nie NULL-owalna:
    -- wchodzi do warunku jednoznaczności, a NULL w SQLite nie równa się NULL-owi
    -- i dwa wgrania tego samego bezimiennego źródła założyłyby dwa komplety.
    zrodlo_kod TEXT    NOT NULL DEFAULT '',
    -- Numer fragmentu w dokumencie źródłowym, liczony od zera.
    kolejnosc  INTEGER NOT NULL,
    -- Sam fragment — ten, który wraca Operatorowi jako cytat (`text`).
    tresc      TEXT    NOT NULL,
    -- Nazwa modelu osadzeń; patrz nagłówek.
    model      TEXT    NOT NULL,
    -- Liczba współrzędnych wektora.
    wymiar     INTEGER NOT NULL CHECK (wymiar > 0),
    -- Współrzędne w zapisie opisanym w nagłówku.
    wektor     BLOB    NOT NULL,
    -- Milisekundy epoki podane przez wołającego.
    utworzono  INTEGER NOT NULL,

    -- Powtórne indeksowanie tego samego dokumentu nadpisuje fragmenty zamiast
    -- dokładać drugi komplet. Bez tego warunku dokument zaindeksowany dwa razy
    -- miałby w wskaźniku dwie kopie każdego zdania i obie wracałyby w odpowiedzi
    -- jako dwa różne trafienia o identycznej treści.
    UNIQUE (zakres, zrodlo_kod, kolejnosc, model)
);

-- Odczyt idzie zawsze po zakresie i modelu, a potem przegląda wszystko, co
-- z tego wyszło (przegląd zupełny — uzasadnienie liczbami w
-- `wiedza/podobienstwo.go`). Ten indeks skraca więc dokładnie ten jeden krok,
-- który da się skrócić: wybór wierszy, których w ogóle wolno dotknąć.
CREATE INDEX idx_fragment_wiedzy_zakres_model ON fragment_wiedzy (model, zakres);

-- Sprzątanie po źródle (`UsunZrodlo`) pyta o parę zakres+kod. Bez tego indeksu
-- każde ponowne zaindeksowanie jednego pliku przeglądałoby całą tabelę.
CREATE INDEX idx_fragment_wiedzy_zrodlo ON fragment_wiedzy (zakres, zrodlo_kod);

-- ── Kategoria okna konfiguracji ──────────────────────────────────────────────
-- Kategoria jest nowa. Kategoria `mowa` nie pasuje, bo opisuje przetwornik
-- dźwięku, z którym wskaźnik znaczenia nie dzieli ani modelu, ani biblioteki,
-- ani powodu istnienia. Kategoria `modele` też nie: opisuje kanał modelu, czyli
-- to, co odpowiada Operatorowi. Silnik osadzeń niczego nie odpowiada; zamienia
-- tekst na liczby i stoi przed rozmową.
INSERT INTO kategoria_ustawien (kod, nazwa, opis, ikona, kolejnosc) VALUES
    ('wiedza', 'Wyszukiwanie po znaczeniu',
     'Lokalny silnik osadzeń budujący wskaźnik znaczenia wiedzy Operatora: interpreter, model, katalog wag i długość fragmentu',
     'kompas', 10)
ON CONFLICT(kod) DO NOTHING;

-- ── Cztery pozycje katalogu ──────────────────────────────────────────────────
-- Wpisy katalogu zamiast stałych w Go — zmiana modelu ma być wierszem, nie
-- wydaniem binarium.
WITH katalog(klucz, kategoria, nazwa, opis, rodzaj, domyslna, podpowiedz,
             wymaga_restartu, kolejnosc) AS (
    VALUES
        ('wiedza_program', 'wiedza', 'Interpreter Pythona',
         'Ścieżka interpretera Pythona liczącego osadzenia. Pusta znaczy: szukaj `python3` na ścieżce wyszukiwania systemu, a nie odmawiaj. Interpreter musi mieć zainstalowaną bibliotekę `fastembed`; jej brak jest odmową NAZYWAJĄCĄ brak, nigdy cichym zejściem na wyszukiwanie po słowach.',
         'path', '', 'python3', 1, 1),
        ('wiedza_model', 'wiedza', 'Model osadzeń',
         'Nazwa modelu zamieniającego tekst na wektor znaczenia. Domyślny jest wielojęzyczny i radzi sobie z polszczyzną — to warunek, nie życzenie: wiedza Operatora jest po polsku. Zmiana modelu unieważnia dotychczasowe wektory, bo dwa modele opisują znaczenie w dwóch różnych przestrzeniach; po zmianie trzeba przebudować wskaźnik.',
         'string', 'sentence-transformers/paraphrase-multilingual-mpnet-base-v2', '', 0, 2),
        ('wiedza_katalog_modeli', 'wiedza', 'Katalog wag modelu',
         'Katalog, do którego pobierane są wagi modelu osadzeń. Pusty znaczy: podkatalog `wiedza/modele` katalogu danych rdzenia — tam, gdzie leżą baza i magazyn biblioteki. Wagi ważą rząd gigabajta, więc Operator może je przenieść na dysk pojemniejszy.',
         'path', '', '', 1, 3),
        ('wiedza_dlugosc_fragmentu', 'wiedza', 'Długość fragmentu w znakach',
         'Docelowa długość fragmentu, na jakie dzielony jest dokument przed osadzeniem. Fragment za długi przekracza okno modelu i zostaje po cichu obcięty; za krótki przestaje nieść kontekst i trafienia zlewają się ze sobą. Wartość spoza przedziału 200–1100 wraca na domyślną.',
         'int', '700', '700', 0, 4)
)
INSERT INTO definicja_ustawienia (klucz, kategoria_id, nazwa, opis, rodzaj_wartosci,
                                  wartosc_domyslna, podpowiedz, wymaga_restartu, kolejnosc)
SELECT k.klucz, kat.id, k.nazwa, k.opis, k.rodzaj, k.domyslna, k.podpowiedz,
       k.wymaga_restartu, k.kolejnosc
  FROM katalog k
  JOIN kategoria_ustawien kat ON kat.kod = k.kategoria
 WHERE true
ON CONFLICT(klucz) DO NOTHING;

-- ── Dopuszczalne poziomy zasięgu ─────────────────────────────────────────────
-- Wszystkie cztery stoją wyłącznie na poziomie globalnym, bo wskaźnik jest jeden
-- na maszynę i wspólny dla wszystkich okien. Gdyby okno A indeksowało modelem X,
-- a okno B pytało modelem Y, oba pracowałyby na tej samej tabeli w dwóch
-- nieporównywalnych przestrzeniach — a Operator zobaczyłby po prostu brak
-- trafień, bez śladu mówiącego dlaczego. Długość fragmentu ma ten sam problem:
-- fragmenty dwóch długości w jednym wskaźniku dają ranking, w którym dłuższe
-- wygrywają z powodu długości, a nie treści.
INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
 CROSS JOIN poziom_zasiegu p
 WHERE d.klucz IN ('wiedza_program', 'wiedza_model', 'wiedza_katalog_modeli',
                   'wiedza_dlugosc_fragmentu')
   AND p.kod = 'globalny'
ON CONFLICT DO NOTHING;

-- ── Dopuszczalne osie ────────────────────────────────────────────────────────
-- Wyłącznie `platform`: ścieżka interpretera i katalog wag są własnością
-- maszyny Operatora, a nie tego, który model odpowiada ani na czyim koncie.
INSERT INTO definicja_ustawienia_os (definicja_id, os)
SELECT d.id, 'platform'
  FROM definicja_ustawienia d
 WHERE d.klucz IN ('wiedza_program', 'wiedza_model', 'wiedza_katalog_modeli',
                   'wiedza_dlugosc_fragmentu')
ON CONFLICT DO NOTHING;
