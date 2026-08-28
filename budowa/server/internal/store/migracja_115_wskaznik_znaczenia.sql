-- Wskaźnik znaczenia trzyma fragmenty wiedzy wraz z ich osadzeniami oraz
-- wpisy katalogu ustawień sterujące silnikiem osadzeń; wektor liczb wyraża
-- bliskość sensu, nie dopasowanie słów jak indeks pełnotekstowy.

CREATE TABLE fragment_wiedzy (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Zakres kontraktowy źródła fragmentu; wartość oznaczająca wszystkie
    -- zakresy nie wchodzi do wiersza.
    zakres     TEXT    NOT NULL CHECK (zakres IN ('library', 'history', 'workspace')),
    -- Nazwa źródła dla człowieka — plik, tytuł rozmowy, ścieżka pliku —
    -- wchodzi wprost do cytatu.
    zrodlo     TEXT    NOT NULL,
    -- Identyfikator źródła; napis pusty jest stanem poprawnym dla źródeł bez
    -- własnego identyfikatora.
    zrodlo_kod TEXT    NOT NULL DEFAULT '',
    -- Numer fragmentu w dokumencie źródłowym, liczony od zera.
    kolejnosc  INTEGER NOT NULL,
    -- Sam fragment — ten, który wraca Operatorowi jako cytat (`text`).
    tresc      TEXT    NOT NULL,
    -- Nazwa modelu osadzeń; wektory dwóch modeli leżą w różnych
    -- przestrzeniach, nieporównywalnych wprost.
    model      TEXT    NOT NULL,
    -- Liczba współrzędnych wektora.
    wymiar     INTEGER NOT NULL CHECK (wymiar > 0),
    -- Współrzędne pojedynczej precyzji, bajty mało-końcowe, wektor podzielony
    -- przez własną długość.
    wektor     BLOB    NOT NULL,
    -- Milisekundy epoki podane przez wołającego.
    utworzono  INTEGER NOT NULL,

    -- Powtórne indeksowanie dokumentu nadpisuje fragmenty zamiast dokładać
    -- drugi komplet.
    UNIQUE (zakres, zrodlo_kod, kolejnosc, model)
);

-- Odczyt idzie zawsze po zakresie i modelu, a potem przegląda przeglądem
-- zupełnym wszystko, co z tego wyszło; ten indeks skraca właśnie wybór
-- wierszy, których wolno dotknąć.
CREATE INDEX idx_fragment_wiedzy_zakres_model ON fragment_wiedzy (model, zakres);

-- Sprzątanie po źródle (`UsunZrodlo`) pyta o parę zakres+kod. Bez tego indeksu
-- każde ponowne zaindeksowanie jednego pliku przeglądałoby całą tabelę.
CREATE INDEX idx_fragment_wiedzy_zrodlo ON fragment_wiedzy (zakres, zrodlo_kod);

-- Kategoria okna konfiguracji jest nowa: opisuje silnik osadzeń, który nie
-- dzieli modelu ani biblioteki z kategorią mowy ani z kategorią modeli
-- odpowiadających operatorowi.
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

-- Wszystkie cztery ustawienia stoją wyłącznie na poziomie globalnym, bo
-- wskaźnik jest jeden na maszynę i wspólny dla wszystkich okien; różne
-- modele albo długości fragmentu w jednym wskaźniku popsułyby trafienia.
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
