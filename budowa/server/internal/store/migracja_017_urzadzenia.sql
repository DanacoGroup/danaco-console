-- Migracja 017 — dopełnienie tabeli `urzadzenie` na potrzeby rozpoznania maszyny.
--
-- Czemu ta migracja istnieje. Tabela `urzadzenie` stoi w schemacie od migracji
-- 001, lecz nie miała ścieżki zapisu: żadne repozytorium jej nie obsługiwało,
-- więc wiersz dało się wnieść wyłącznie poleceniem wprost w bazie (tak robią
-- dzisiaj rusztowania testów). Blokowało to punkty dostępu rodzaju
-- `localDirectory`, bo więz migracji 013
-- (`CHECK(rodzaj <> 'localDirectory' OR urzadzenie_id IS NOT NULL)`) wymaga
-- wskazania urządzenia, a żadnego urządzenia nie dało się założyć.
--
-- Co dokłada. Dwie kolumny, obie wynikające z rozpoznania maszyny, na której
-- działa rdzeń:
--   `nazwa_hosta` — nazwa hosta maszyny; fakt maszyny, nie napis Operatora.
--       Kolumna `nazwa` pozostaje nazwą do pokazania i wolno ją zmienić —
--       rozpoznanie startowe jej nie nadpisuje, żeby nie kasować pracy
--       Operatora przy każdym uruchomieniu rdzenia.
--   `biezace`    — oznaczenie „to jest maszyna, na której działa ten rdzeń".
--       Bez niego pytanie „na czym stoi rdzeń" musiałoby być rozstrzygane
--       ponownym rozpoznaniem w każdej warstwie, która chce zaproponować
--       urządzenie dla katalogu lokalnego; byłoby to drugie źródło tej samej
--       odpowiedzi.
--
-- ALTER TABLE, nie przebudowa. Do `urzadzenie` odwołują się `polaczenie`
-- (migracja 001) oraz `punkt_dostepu` (migracja 013). Odtworzenie tabeli
-- wymagałoby wyłączenia więzów kluczy obcych, a `PRAGMA foreign_keys` nie
-- działa wewnątrz transakcji — każdy krok migracji idzie tu w jednej
-- transakcji. Zostaje dokładanie kolumn ze stałą wartością domyślną, na co
-- SQLite pozwala.

ALTER TABLE urzadzenie ADD COLUMN nazwa_hosta TEXT NOT NULL DEFAULT '';
ALTER TABLE urzadzenie ADD COLUMN biezace INTEGER NOT NULL DEFAULT 0
    CHECK(biezace IN (0,1));

-- Maszyna bieżąca jest najwyżej jedna i pilnuje tego baza, nie warstwa wyżej —
-- tak samo jak głównego nadania okna w migracji 013.
CREATE UNIQUE INDEX idx_urzadzenie_biezace ON urzadzenie(biezace) WHERE biezace = 1;

-- Odczyt po identyfikatorze sprzętowym idzie już po więzie UNIQUE z migracji 001,
-- więc osobnego indeksu nie zakładamy.
--
-- Wierszy zastanych migracja nie zgaduje: pustej `nazwa_hosta` nie wolno wypełnić
-- wartością kolumny `nazwa`, bo nazwa do pokazania nie jest dowodem nazwy hosta.
-- Puste znaczy „nierozpoznane" i wypełni je pierwsze rozpoznanie maszyny.
