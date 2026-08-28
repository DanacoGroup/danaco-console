-- Migracja 404 zmienia wartości domyślne katalogów wag przesiewu i osi
-- obrazu tak, aby wskazywały wagi już rozłożone na maszynie.

UPDATE definicja_ustawienia
   SET wartosc_domyslna = '/opt/danaco-modele/reranker',
       podpowiedz = '/opt/danaco-modele/reranker',
       opis = 'Katalog wag krzyżowego kodera. Domyślnie wskazuje wagi rozłożone na maszynie obok rdzenia — leżą tam gotowe, więc nic się nie pobiera. Katalog, w którym wag nie ma, służy za miejsce ich pobrania; pusty znaczy podkatalog `wiedza/modele/przesiew` katalogu danych rdzenia. Wagi ważą rząd gigabajtów, więc Operator może je przenieść na dysk pojemniejszy.'
 WHERE klucz = 'wiedza_katalog_przesiewu';

UPDATE definicja_ustawienia
   SET wartosc_domyslna = '/opt/danaco-modele/clip',
       podpowiedz = '/opt/danaco-modele/clip',
       opis = 'Katalog wag modelu osi obrazu. Domyślnie wskazuje wagi rozłożone na maszynie obok rdzenia — leżą tam gotowe, więc nic się nie pobiera. Katalog, w którym wag nie ma, służy za miejsce ich pobrania; pusty znaczy podkatalog `wiedza/modele/obraz` katalogu danych rdzenia. Katalog osobny od katalogu wag osadzarki i przesiewu, bo są to trzy różne modele — nazwy podkatalogów odpowiadają zdolnościom, nie wydawcom, więc zmiana modelu jest wartością ustawienia, a nie przeprowadzką katalogu.'
 WHERE klucz = 'wiedza_katalog_obrazu';
