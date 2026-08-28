-- Migracja dokłada do tabeli urządzenie dwie kolumny wynikające z rozpoznania maszyny: nazwę
-- hosta i oznaczenie maszyny bieżącej.

ALTER TABLE urzadzenie ADD COLUMN nazwa_hosta TEXT NOT NULL DEFAULT '';
ALTER TABLE urzadzenie ADD COLUMN biezace INTEGER NOT NULL DEFAULT 0
    CHECK(biezace IN (0,1));

-- Maszyna bieżąca jest najwyżej jedna i pilnuje tego baza, nie warstwa wyżej —
-- tak samo jak głównego nadania okna w migracji 013.
CREATE UNIQUE INDEX idx_urzadzenie_biezace ON urzadzenie(biezace) WHERE biezace = 1;

-- Osobnego indeksu nie zakładamy: odczyt idzie już po więzie unikalności sprzętu.
