-- Migracja dokłada odwracalne usunięcie sesji jako znacznik czasu w wierszu, zamiast
-- natychmiastowego usunięcia całego wiersza.

ALTER TABLE sesja ADD COLUMN usunieto_o TEXT;

-- Czyszczenie trwałe przy starcie rdzenia pyta wyłącznie o wiersze sesji niosące ten
-- znacznik usunięcia w koszu.
CREATE INDEX idx_sesja_kosz ON sesja (usunieto_o) WHERE usunieto_o IS NOT NULL;
