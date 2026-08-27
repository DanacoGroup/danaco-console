-- Migracja 403 — silnik mowy startuje na wagach, które leżą obok rdzenia,
-- zamiast szukać ich w pamięci podręcznej konta, na którym akurat pracuje.
--
-- Migracja 075 założyła `mowa_katalog_modeli` z wartością pustą, czyli
-- „katalog domyślny biblioteki". Dla faster-whisper znaczy to pamięć podręczna
-- Huba konta uruchamiającego rdzeń — dziś `~/.cache/huggingface/hub`
-- (`pomocniki/transkrypcja/silnik.py`, `katalogi_cache`). Wagi modelu
-- domyślnego tam właśnie stały, więc transkrypcja działała, ale zależnie od
-- tego, na czyim koncie stoi proces: usługa systemowa albo konto serwisowe ma
-- inny katalog domowy i tych samych wag nie widzi — pobierze drugą kopię.
--
-- Wagi zostały przeniesione do `/opt/danaco-modele/mowa`, tam gdzie stoją wagi
-- pozostałych zdolności liczących lokalnie (`embedder`, `reranker`, `clip`),
-- i ta ścieżka jest odtąd wartością domyślną. Układ katalogu jest układem
-- pamięci podręcznej Huba (`models--Systran--faster-whisper-small/`), bo
-- pomocnik przekazuje tę ścieżkę biblioteki jako `download_root` i tego układu
-- w niej szuka.
--
-- Migracji 075 się nie zmienia — jej suma kontrolna stoi w rejestrze `migracja`
-- u każdego, kto rdzeń postawił, a niezgodność sumy wywraca start rdzenia
-- (`store/migracje.go`). Zmiana wartości domyślnej jest osobnym krokiem,
-- wzorem migracji 401 dla osadzarki.
--
-- UPDATE, a nie INSERT ... ON CONFLICT: wiersz istnieje od migracji 075
-- i zmienia się w nim wyłącznie wartość domyślna wraz z podpowiedzią i opisem,
-- który tę wartość tłumaczy Operatorowi. Warunek na klucz czyni krok
-- idempotentnym i nieszkodliwym na bazie, gdzie tego wiersza nie ma.
--
-- Zapis Operatora zostaje nietknięty. Wartość domyślna jest ostatnim ogniwem
-- rozstrzygania (`konfig/rozstrzyganie.go`), więc maszyna, na której katalog
-- wag ustawiono ręcznie, dalej pracuje na tym wskazaniu.
--
-- Wskazanie katalogu, w którym wag nie ma, nie jest odmową: pomocnik traktuje
-- go wtedy jak miejsce pobrania, czyli zachowuje się tak, jak przy wartości
-- pustej. Ten krok nie może więc zepsuć maszyny, na której wag pod tą ścieżką
-- nie rozłożono — najgorsze, co na niej zrobi, to przywrócenie stanu sprzed
-- siebie.

UPDATE definicja_ustawienia
   SET wartosc_domyslna = '/opt/danaco-modele/mowa',
       podpowiedz = '/opt/danaco-modele/mowa',
       opis = 'Katalog wag modelu rozpoznawania mowy. Domyślnie wskazuje wagi rozłożone na maszynie obok rdzenia — leżą tam gotowe, więc nic się nie pobiera, niezależnie od tego, na czyim koncie stoi proces rdzenia. Katalog, w którym wag nie ma, służy za miejsce ich pobrania; pusty znaczy pamięć podręczną biblioteki w katalogu domowym konta. Wagi ważą setki megabajtów, więc Operator może je przenieść na dysk pojemniejszy.'
 WHERE klucz = 'mowa_katalog_modeli';
