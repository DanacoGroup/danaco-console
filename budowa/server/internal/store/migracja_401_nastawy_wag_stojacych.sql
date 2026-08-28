-- Migracja 401 — wskaźnik znaczenia startuje na wagach, które już leżą
-- na maszynie, zamiast pobierać drugą kopię tego samego.
--
-- Migracja 115 założyła cztery nastawy wskaźnika znaczenia z wartościami
-- domyślnymi opisującymi stan, w którym maszyna nie ma jeszcze żadnych wag:
-- model `sentence-transformers/paraphrase-multilingual-mpnet-base-v2` i katalog
-- pusty, czyli „pobierz do podkatalogu katalogu danych rdzenia". Od tamtej pory
-- wagi na maszynie stoją: 4,3 GB w `/opt/danaco-modele/embedder`, z czego
-- 2,2 GB to wydanie ONNX, którym liczy pomocnik osadzeń. Przy wartościach
-- z migracji 115 świeże wdrożenie ich nie widziało i przy pierwszym
-- `knowledge.index` sięgało po drugi model do sieci — pobranie liczone
-- w gigabajtach zamiast wystartowania z tego, co jest na dysku.
--
-- Migracji 115 się nie zmienia. Bazy założone wcześniej mają ją wykonaną
-- i odnotowaną wraz z sumą kontrolną, a suma kroku już zastosowanego musi się
-- zgadzać (`store/migracje.go`); zmiana treści tamtego pliku wywróciłaby start
-- rdzenia u każdego, kto go dziś ma. Zmiana wartości domyślnej jest więc
-- osobnym krokiem, tak samo jak każda inna poprawka danych katalogu.
--
-- UPDATE, a nie INSERT ... ON CONFLICT: wiersze istnieją, bo założyła je
-- migracja 115, i zmienia się w nich wyłącznie wartość domyślna wraz z opisem,
-- który tę wartość tłumaczy Operatorowi. Warunek na klucz czyni krok
-- idempotentnym i nieszkodliwym na bazie, gdzie tych wierszy z jakiegoś powodu
-- nie ma.
--
-- Zapis Operatora zostaje nietknięty. Wartość domyślna jest ostatnim ogniwem
-- rozstrzygania (`konfig/rozstrzyganie.go`): każdy zapis na dowolnym poziomie
-- zasięgu ma przed nią pierwszeństwo, więc maszyna, na której wdrożeniowiec
-- ustawił katalog wag ręcznie, dalej pracuje na jego wskazaniu. Ten krok zmienia
-- wyłącznie to, co widzi Operator, który niczego nie ustawił.
--
-- Nazwa modelu wchodzi do każdej pozycji wskaźnika i zawęża odczyt przy szukaniu
-- (`fragment_wiedzy.model`, migracja 115), więc na maszynie, która zdążyła
-- zbudować wskaźnik modelem poprzednim, wiersze tamtego modelu zostają i nie
-- mieszają się z nowymi — sprząta je `knowledge.index` z `rebuild`. Kasowanie
-- ich tutaj byłoby usunięciem pracy Operatora krokiem schematu.
--
-- Katalog wag jest ścieżką bezwzględną i to jest świadome: nastawa opisuje
-- rozmieszczenie wag na maszynie, a nie coś, co da się wyliczyć z katalogu
-- danych rdzenia. Katalog, w którym pliku ONNX nie ma, pomocnik traktuje jak
-- pamięć podręczną, czyli zachowuje się dokładnie tak jak przy wartości pustej
-- (`wiedza/pomocnik_osadzen.py`); wskazanie nieistniejącego katalogu nie jest
-- więc odmową, tylko powrotem do zachowania sprzed tego kroku.

UPDATE definicja_ustawienia
   SET wartosc_domyslna = 'BAAI/bge-m3',
       opis = 'Nazwa modelu zamieniającego tekst na wektor znaczenia. Domyślny jest wielojęzyczny i radzi sobie z polszczyzną — to warunek, nie życzenie: wiedza Operatora jest po polsku. Jego wagi leżą już na maszynie, w katalogu wskazanym nastawą `wiedza_katalog_modeli`, więc pierwsze indeksowanie niczego nie pobiera. Zmiana modelu unieważnia dotychczasowe wektory, bo dwa modele opisują znaczenie w dwóch różnych przestrzeniach; po zmianie trzeba przebudować wskaźnik.'
 WHERE klucz = 'wiedza_model';

UPDATE definicja_ustawienia
   SET wartosc_domyslna = '/opt/danaco-modele/embedder',
       podpowiedz = '/opt/danaco-modele/embedder',
       opis = 'Katalog wag modelu osadzeń. Domyślnie wskazuje wagi rozłożone na maszynie obok rdzenia — leżą tam gotowe, więc nic się nie pobiera. Katalog, w którym wag nie ma, służy za miejsce ich pobrania; pusty znaczy podkatalog `wiedza/modele` katalogu danych rdzenia. Wagi ważą rząd gigabajtów, więc Operator może je przenieść na dysk pojemniejszy.'
 WHERE klucz = 'wiedza_katalog_modeli';
