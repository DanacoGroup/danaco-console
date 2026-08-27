-- Migracja 402 — cztery nastawy przesiewu i osi obrazu dostają wiersz
-- w katalogu ustawień, więc Operator ustawi je z okna konfiguracji.
--
-- Klucze `wiedza_model_przesiewu`, `wiedza_katalog_przesiewu`,
-- `wiedza_model_obrazu` i `wiedza_katalog_obrazu` stoją w kodzie od czasu
-- dołożenia przesiewu i osi obrazu (`wiedza/ustawienia.go`), ale wiersza
-- `definicja_ustawienia` nie miały. Zdolność mimo to działa: rozstrzyganie
-- nastawy czyta zapis niezależnie od katalogu definicji
-- (`konfig/rozstrzyganie.go`), więc wartość zapisana wprost w tabeli
-- `ustawienie` dochodzi do silnika. Czego bez wiersza katalogu nie ma, to drogi
-- Operatora: `config.set` odmawia klucza spoza katalogu, a okno konfiguracji
-- wystawia wyłącznie pozycje katalogu. Cztery nastawy były więc ustawialne
-- ręcznym zapisem do bazy i tylko nim.
--
-- Migracji 115 się nie zmienia — jej suma kontrolna stoi w rejestrze `migracja`
-- u każdego, kto rdzeń postawił, a niezgodność sumy wywraca start rdzenia
-- (`store/migracje.go`). Cztery wiersze dokłada więc osobny krok, wzorem tego,
-- jak migracja 401 zmieniła wartości domyślne dwóch nastaw z migracji 115.
--
-- Kategoria, zasięg i oś są te same, co u czterech nastaw migracji 115, i to
-- z tych samych powodów. Kategoria `wiedza`, bo to ten sam silnik. Zasięg
-- wyłącznie globalny, bo wskaźnik znaczenia jest jeden na maszynę: przesiew
-- układający kolejność dwoma różnymi koderami w dwóch oknach dawałby dwie
-- nieporównywalne kolejności tego samego wyniku. Oś wyłącznie `platform`, bo
-- katalog wag i nazwa modelu liczącego lokalnie są własnością maszyny, a nie
-- konta ani kanału modelu.
--
-- Wartości domyślne są kopią stałych `wiedza/ustawienia.go` co do znaku —
-- `ModelPrzesiewuDomyslny`, `ModelObrazuDomyslny` i `katalogNiewskazany`.
-- Rozjazd znaczyłby dwie prawdy o tym, czym rdzeń liczy, zależne od drogi
-- wywołania: rozstrzygacz zasięgu oddaje wartość z tej kolumny, a stała pakietu
-- wchodzi tam, gdzie rozstrzygacza nie ma (`core/adapter_modul_wiedza.go`).
--
-- Oba katalogi wag zostają PUSTE, choć wagi obu modeli leżą na maszynie
-- (`/opt/danaco-modele/reranker`, `/opt/danaco-modele/clip`). Wartość niepusta
-- byłaby tutaj drugą prawdą wobec stałej `katalogNiewskazany`, a
-- `internal/wiedza/` leży poza terenem tej zmiany. Wskazanie wag stojących jest
-- osobnym krokiem, obejmującym zarazem stałą i ten wiersz — dokładnie tak, jak
-- migracja 401 zrobiła to dla osadzarki.
--
-- `ON CONFLICT DO NOTHING` czyni krok idempotentnym i nieszkodliwym na bazie,
-- gdzie te wiersze z jakiegoś powodu już stoją.

WITH katalog(klucz, kategoria, nazwa, opis, rodzaj, domyslna, podpowiedz,
             wymaga_restartu, kolejnosc) AS (
    VALUES
        ('wiedza_model_przesiewu', 'wiedza', 'Model przesiewu',
         'Nazwa krzyżowego kodera układającego na nowo kandydatów pierwszego przebiegu. Koder czyta pytanie i fragment razem, więc ocenia trafność wierniej niż porównanie dwóch wektorów — ale liczy się dla każdego kandydata z osobna, więc wchodzi dopiero po zawężeniu wyniku. Domyślny jest wielojęzyczny z tego samego powodu, dla którego wielojęzyczna jest osadzarka: wiedza Operatora jest po polsku, a koder jednojęzyczny oceniałby polskie fragmenty przez podobieństwo do angielskiego pytania.',
         'string', 'BAAI/bge-reranker-v2-m3', '', 0, 5),
        ('wiedza_katalog_przesiewu', 'wiedza', 'Katalog wag przesiewu',
         'Katalog wag krzyżowego kodera. Pusty znaczy: podkatalog `wiedza/modele/przesiew` katalogu danych rdzenia. Katalog, w którym wag nie ma, służy za miejsce ich pobrania — wskazanie ścieżki nieistniejącej nie jest odmową. Wagi ważą rząd gigabajtów, więc Operator może je przenieść na dysk pojemniejszy.',
         'path', '', '', 1, 6),
        ('wiedza_model_obrazu', 'wiedza', 'Model osi obrazu',
         'Nazwa modelu dwuwieżowego wiążącego obraz ze zdaniem w jednej przestrzeni — to nim wyszukiwanie odpowiada na pytanie zadane słowami o rzecz widoczną na obrazie. Oś obrazu wchodzi na żądanie i liczy się raz na zapytanie, więc rozstrzyga o trafności, a nie o czasie odpowiedzi; dlatego domyślne jest wydanie duże, nie podstawowe.',
         'string', 'openai/clip-vit-large-patch14', '', 0, 7),
        ('wiedza_katalog_obrazu', 'wiedza', 'Katalog wag osi obrazu',
         'Katalog wag modelu osi obrazu. Pusty znaczy: podkatalog `wiedza/modele/obraz` katalogu danych rdzenia. Katalog osobny od katalogu wag osadzarki i przesiewu, bo są to trzy różne modele — nazwy podkatalogów odpowiadają zdolnościom, nie wydawcom, więc zmiana modelu jest wartością ustawienia, a nie przeprowadzką katalogu.',
         'path', '', '', 1, 8)
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
INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
 CROSS JOIN poziom_zasiegu p
 WHERE d.klucz IN ('wiedza_model_przesiewu', 'wiedza_katalog_przesiewu',
                   'wiedza_model_obrazu', 'wiedza_katalog_obrazu')
   AND p.kod = 'globalny'
ON CONFLICT DO NOTHING;

-- ── Dopuszczalne osie ────────────────────────────────────────────────────────
INSERT INTO definicja_ustawienia_os (definicja_id, os)
SELECT d.id, 'platform'
  FROM definicja_ustawienia d
 WHERE d.klucz IN ('wiedza_model_przesiewu', 'wiedza_katalog_przesiewu',
                   'wiedza_model_obrazu', 'wiedza_katalog_obrazu')
ON CONFLICT DO NOTHING;
