-- Migracja 403 zmienia wartość domyślną katalogu wag silnika mowy tak, aby
-- wskazywała wagi rozłożone obok rdzenia.

UPDATE definicja_ustawienia
   SET wartosc_domyslna = '/opt/danaco-modele/mowa',
       podpowiedz = '/opt/danaco-modele/mowa',
       opis = 'Katalog wag modelu rozpoznawania mowy. Domyślnie wskazuje wagi rozłożone na maszynie obok rdzenia — leżą tam gotowe, więc nic się nie pobiera, niezależnie od tego, na czyim koncie stoi proces rdzenia. Katalog, w którym wag nie ma, służy za miejsce ich pobrania; pusty znaczy pamięć podręczną biblioteki w katalogu domowym konta. Wagi ważą setki megabajtów, więc Operator może je przenieść na dysk pojemniejszy.'
 WHERE klucz = 'mowa_katalog_modeli';
