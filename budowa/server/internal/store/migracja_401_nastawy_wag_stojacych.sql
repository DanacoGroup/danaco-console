-- Migracja 401 zmienia wartości domyślne wskaźnika znaczenia tak, aby
-- wskazywały wagi już rozłożone na maszynie.

UPDATE definicja_ustawienia
   SET wartosc_domyslna = 'BAAI/bge-m3',
       opis = 'Nazwa modelu zamieniającego tekst na wektor znaczenia. Domyślny jest wielojęzyczny i radzi sobie z polszczyzną — to warunek, nie życzenie: wiedza Operatora jest po polsku. Jego wagi leżą już na maszynie, w katalogu wskazanym nastawą `wiedza_katalog_modeli`, więc pierwsze indeksowanie niczego nie pobiera. Zmiana modelu unieważnia dotychczasowe wektory, bo dwa modele opisują znaczenie w dwóch różnych przestrzeniach; po zmianie trzeba przebudować wskaźnik.'
 WHERE klucz = 'wiedza_model';

UPDATE definicja_ustawienia
   SET wartosc_domyslna = '/opt/danaco-modele/embedder',
       podpowiedz = '/opt/danaco-modele/embedder',
       opis = 'Katalog wag modelu osadzeń. Domyślnie wskazuje wagi rozłożone na maszynie obok rdzenia — leżą tam gotowe, więc nic się nie pobiera. Katalog, w którym wag nie ma, służy za miejsce ich pobrania; pusty znaczy podkatalog `wiedza/modele` katalogu danych rdzenia. Wagi ważą rząd gigabajtów, więc Operator może je przenieść na dysk pojemniejszy.'
 WHERE klucz = 'wiedza_katalog_modeli';
