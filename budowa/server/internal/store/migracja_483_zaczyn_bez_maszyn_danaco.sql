-- Punkty dostępu bez adresów trzech maszyn produkcyjnych Danaco.
--
-- Adres, użytkownik i ścieżka klucza maszyny to dane jednej sieci, nie słownik
-- platformy; zaczyn migracji 013 wnosi je do każdej instalacji. Krok zdejmuje
-- adres, nie pracę Operatora: punkt z wpisem spoza zaczynu zostaje i traci
-- pola maszyny. Punkt, któremu Operator zostawił adres zaczynu, a zmienił
-- użytkownika, ścieżkę klucza albo port, wpisu nie ma i schodzi w całości.

-- Pary kod-adres zaczynu stoją w kroku raz.
CREATE TEMP TABLE zaczyn_maszyny (kod TEXT NOT NULL, host TEXT NOT NULL);
INSERT INTO zaczyn_maszyny (kod, host) VALUES
    ('danaco-system', '57.128.253.74'),
    ('danaco-data',   '162.19.226.145'),
    ('danaco-web',    '137.74.41.149');

-- Ścieżki wpisane punktom migracją 013 (013:172-175). Ścieżka spoza tego wykazu
-- pochodzi od Operatora.
CREATE TEMP TABLE zaczyn_sciezki (kod TEXT NOT NULL, sciezka TEXT NOT NULL);
INSERT INTO zaczyn_sciezki (kod, sciezka) VALUES
    ('danaco-system', '/opt/danaco'),
    ('danaco-data',   '/opt/danaco'),
    ('danaco-data',   '/opt/danaco-backups'),
    ('danaco-web',    '/opt/danaco');

-- Punkt z adresem zmienionym przez Operatora do wykazu nie wchodzi. Kolumna
-- `obciazony` mówi, czy przy punkcie stoi wpis spoza zaczynu: odwołanie do
-- poświadczenia albo wskazanie urządzenia we własnym wierszu, nadanie dostępu,
-- konektor eksperta, pozycja katalogu rozszerzeń, ścieżka albo słowo trybu.
-- Zaczyn 013 obu kolumn punktu nie wstawia, więc wypełnione są wpisem Operatora.
CREATE TEMP TABLE punkt_zaczynu AS
SELECT p.id AS id,
       (p.poswiadczenie_odwolanie IS NOT NULL
     OR p.urzadzenie_id IS NOT NULL
     OR EXISTS (SELECT 1 FROM nadanie_dostepu n WHERE n.punkt_dostepu_id = p.id)
     OR EXISTS (SELECT 1 FROM agent_konektor a WHERE a.punkt_dostepu_id = p.id)
     OR EXISTS (SELECT 1 FROM rozszerzenie r WHERE r.punkt_dostepu_id = p.id)
     OR EXISTS (SELECT 1 FROM korzen_punktu_dostepu k
                 WHERE k.punkt_dostepu_id = p.id
                   AND k.sciezka NOT IN (SELECT w.sciezka FROM zaczyn_sciezki w
                                          WHERE w.kod = p.kod))
     OR EXISTS (SELECT 1 FROM argument_trybu_mostu t
                 WHERE t.punkt_dostepu_id = p.id
                   AND NOT ((t.tryb = 'read'  AND t.argument = 'odczyt')
                         OR (t.tryb = 'write' AND t.argument = 'zapis')))) AS obciazony
FROM punkt_dostepu p
JOIN zaczyn_maszyny z ON z.kod = p.kod AND z.host = p.host;

-- Punkt bez wpisu Operatora znika wraz z własnym zaczynem. Więzy kluczy obcych
-- są na czas kroku wygaszone (migracje.go), więc kaskada nie odpala sama.
DELETE FROM korzen_punktu_dostepu
 WHERE punkt_dostepu_id IN (SELECT id FROM punkt_zaczynu WHERE obciazony = 0);

DELETE FROM argument_trybu_mostu
 WHERE punkt_dostepu_id IN (SELECT id FROM punkt_zaczynu WHERE obciazony = 0);

DELETE FROM punkt_dostepu
 WHERE id IN (SELECT id FROM punkt_zaczynu WHERE obciazony = 0);

-- Punkt z wpisem Operatora zostaje wraz z nim; maszynę traci przez wyzerowanie
-- pól. Opis zaczynu powtarza nazwę mostu prozą (013:150, 013:160), więc idzie
-- do zera razem z kolumną `nazwa_mostu`. Za adres wchodzi `kod` punktu — więz
-- CHECK(rodzaj <> 'mcpBridge' OR host <> '') z 013:66 pustego nie dopuszcza.
UPDATE punkt_dostepu
   SET host             = kod,
       opis             = '',
       uzytkownik       = '',
       sciezka_klucza   = '',
       polecenie_startu = '',
       nazwa_mostu      = '',
       zaktualizowano   = strftime('%Y-%m-%dT%H:%M:%fZ','now')
 WHERE id IN (SELECT id FROM punkt_zaczynu WHERE obciazony = 1);

-- Tabele tymczasowe żyją na połączeniu, a to wraca do puli po kroku.
DROP TABLE punkt_zaczynu;
DROP TABLE zaczyn_sciezki;
DROP TABLE zaczyn_maszyny;
