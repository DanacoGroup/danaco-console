-- Dodaje poziom zasięgu aplikacja jako najszerszy, poprawką warunku CHECK schematu tabeli poziom_zasiegu, dla nastaw dotyczących samego programu, niezależnie od treści platformy.

PRAGMA writable_schema = ON;

UPDATE sqlite_master
   SET sql = replace(sql,
                     '''globalny'',''srodowisko''',
                     '''aplikacja'',''globalny'',''srodowisko''')
 WHERE type = 'table' AND name = 'poziom_zasiegu';

PRAGMA schema_version = 1000119;

PRAGMA writable_schema = OFF;

INSERT INTO poziom_zasiegu (kod, nazwa, pierwszenstwo) VALUES
    ('aplikacja', 'Aplikacja', 0);
