-- Trzy korzenie pracy Operatora zostały poza migracją 484 i nie mają czym
-- wskazać właściciela wiersza: kolejka zadań, raport badania i przebieg wsadu
-- Studia. Każdy z nich powstaje na żądanie Operatora i żaden nie wisi kluczem
-- obcym na korzeniu już zawężonym — kolejka globalna nie ma sesji, raport
-- badania stoi samodzielnie, przebieg wsadu wiąże się ze stroną, nie z kontem.
--
-- Kolumna dopuszcza NULL i nie ma klucza obcego, tak samo jak w migracji 484:
-- wiersze zastane powstały przed rozdzieleniem kont i czyta się je jako konto
-- najstarsze. Wskaźnika krok nie zakłada — wejdzie tam, gdzie wskaże go pomiar.

ALTER TABLE kolejka ADD COLUMN konto_id INTEGER;
ALTER TABLE raport_badania ADD COLUMN konto_id INTEGER;
ALTER TABLE przebieg_wsadu_studio ADD COLUMN konto_id INTEGER;
