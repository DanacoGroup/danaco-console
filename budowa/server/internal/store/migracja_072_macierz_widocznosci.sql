-- Migracja dopełnia macierz widoczności modułów do pełnych czterdziestu pięciu wierszy,
-- dopisując brakujące pary jako niewidoczne.

-- Siedemnaście par niewidocznych dla trzech środowisk, wypisanych jawnie parami
-- środowiska i modułu wykazu.
WITH macierz(srodowisko_kod, modul_kod) AS (
    VALUES
        -- Środowisko pracy z treścią; poza wykazem zostają narzędzia wytwórcze i programistyczne.
        ('talkin',     'automations'),
        ('talkin',     'design'),
        ('talkin',     'terminal'),
        ('talkin',     'developer'),
        ('talkin',     'diagnostics'),
        ('talkin',     'apps'),
        -- Środowisko prowadzenia projektu; poza wykazem zostaje warsztat kodu i dwa moduły TalkIn.
        ('workspace',  'translate'),
        ('workspace',  'assistant'),
        ('workspace',  'terminal'),
        ('workspace',  'developer'),
        ('workspace',  'diagnostics'),
        -- Środowisko programowania; poza wykazem zostaje praca z dokumentem i źródłami.
        ('codestudio', 'studio'),
        ('codestudio', 'browser'),
        ('codestudio', 'research'),
        ('codestudio', 'library'),
        ('codestudio', 'translate'),
        ('codestudio', 'assistant')
)
INSERT INTO srodowisko_modul (srodowisko_id, modul_id, kolejnosc, widoczny)
SELECT s.id, m.id, 0, 0
  FROM macierz
  JOIN srodowisko s ON s.kod = macierz.srodowisko_kod
  JOIN modul m      ON m.kod = macierz.modul_kod
 WHERE true
ON CONFLICT(srodowisko_id, modul_id) DO NOTHING;
