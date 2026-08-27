-- Migracja przenosi do danych fakt, że terminal jest oknem pomocniczym trzech modułów,
-- znikając z bocznej nawigacji jako osobny moduł.

-- Terminal traci widoczność w środowisku programowania, zostając w macierzy z
-- widocznością wyłączoną zamiast usunięcia wiersza.
UPDATE srodowisko_modul
   SET widoczny  = 0,
       kolejnosc = 0
 WHERE srodowisko_id = (SELECT id FROM srodowisko WHERE kod = 'codestudio')
   AND modul_id      = (SELECT id FROM modul      WHERE kod = 'terminal');

-- Siedem pozostałych pozycji środowiska programowania dostaje kolejność bez przerwy
-- w numeracji, wartościami bezwzględnymi.
WITH nowa_kolejnosc(modul_kod, kolejnosc) AS (
    VALUES
        ('workspace',   0),
        ('roundtable',  1),
        ('design',      2),
        ('developer',   3),
        ('diagnostics', 4),
        ('apps',        5),
        ('agents',      6)
)
UPDATE srodowisko_modul
   SET kolejnosc = (SELECT n.kolejnosc
                      FROM nowa_kolejnosc n
                      JOIN modul m ON m.kod = n.modul_kod
                     WHERE m.id = srodowisko_modul.modul_id)
 WHERE srodowisko_id = (SELECT id FROM srodowisko WHERE kod = 'codestudio')
   AND modul_id IN (SELECT m.id
                      FROM modul m
                      JOIN nowa_kolejnosc n ON n.modul_kod = m.kod);
