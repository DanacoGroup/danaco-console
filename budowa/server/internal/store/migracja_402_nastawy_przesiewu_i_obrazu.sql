-- Migracja 402 dokłada do katalogu ustawień cztery nastawy przesiewu i osi
-- obrazu, które kod już obsługiwał, lecz były ustawialne wyłącznie ręcznym
-- zapisem do bazy.

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

-- Cztery nastawy przesiewu i osi obrazu mają zasięg wyłącznie globalny, ponieważ
-- wskaźnik znaczenia jest jeden na maszynę, a nie na konto czy kanał modelu.
INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
 CROSS JOIN poziom_zasiegu p
 WHERE d.klucz IN ('wiedza_model_przesiewu', 'wiedza_katalog_przesiewu',
                   'wiedza_model_obrazu', 'wiedza_katalog_obrazu')
   AND p.kod = 'globalny'
ON CONFLICT DO NOTHING;

-- Cztery nastawy przesiewu i osi obrazu obowiązują wyłącznie oś platform,
-- ponieważ katalog wag i nazwa modelu są własnością maszyny, a nie konta.
INSERT INTO definicja_ustawienia_os (definicja_id, os)
SELECT d.id, 'platform'
  FROM definicja_ustawienia d
 WHERE d.klucz IN ('wiedza_model_przesiewu', 'wiedza_katalog_przesiewu',
                   'wiedza_model_obrazu', 'wiedza_katalog_obrazu')
ON CONFLICT DO NOTHING;
