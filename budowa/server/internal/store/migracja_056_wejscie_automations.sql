-- Migracja dokłada wejście do modułu automatyzacji w macierzy środowisko-moduł dla
-- dwóch środowisk, na końcu wykazu.

WITH macierz(srodowisko_kod, modul_kod, kolejnosc) AS (
    VALUES
        ('workspace',  'automations', 9),
        ('codestudio', 'automations', 8)
)
INSERT INTO srodowisko_modul (srodowisko_id, modul_id, kolejnosc, widoczny)
SELECT s.id, m.id, macierz.kolejnosc, 1
  FROM macierz
  JOIN srodowisko s ON s.kod = macierz.srodowisko_kod
  JOIN modul m      ON m.kod = macierz.modul_kod
 WHERE true
ON CONFLICT(srodowisko_id, modul_id) DO NOTHING;
