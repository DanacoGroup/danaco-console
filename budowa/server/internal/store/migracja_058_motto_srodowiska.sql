-- Migracja dokłada kolumnę motta w tabeli środowiska wraz z treścią czterech mott,
-- osobną od opisu środowiska.

ALTER TABLE srodowisko ADD COLUMN motto TEXT;

UPDATE srodowisko SET motto = 'Pracuj z wiedzą, dokumentami i modelami AI w jednym miejscu.'
 WHERE kod = 'talkin';
UPDATE srodowisko SET motto = 'Organizuj projekty, procesy i zadania wspierane przez sztuczną inteligencję.'
 WHERE kod = 'workspace';
UPDATE srodowisko SET motto = 'Projektuj, twórz i rozwijaj oprogramowanie wspólnie z AI.'
 WHERE kod = 'codestudio';
UPDATE srodowisko SET motto = 'Zarządzaj zespołem modeli AI realizujących złożone procesy.'
 WHERE kod = 'multitaskingai';
