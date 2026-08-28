-- Migracja 372 wprowadza do katalogu ustawień konto nadawcze platformy
-- wraz z nastawami serwera poczty wychodzącej.

INSERT INTO kategoria_ustawien (kod, nazwa, opis, ikona, kolejnosc) VALUES
    ('nadawca', 'Konto nadawcze platformy',
     'Serwer poczty wychodzącej, którym platforma wysyła potwierdzenie adresu i drogę odzyskania konta',
     'koperta', 9)
ON CONFLICT(kod) DO NOTHING;

INSERT INTO definicja_ustawienia (klucz, kategoria_id, nazwa, opis, rodzaj_wartosci,
                                  wartosc_domyslna, podpowiedz, wymaga_restartu, kolejnosc)
SELECT wykaz.klucz, kat.id, wykaz.nazwa, wykaz.opis, wykaz.rodzaj,
       wykaz.domyslna, wykaz.podpowiedz, 0, wykaz.kolejnosc
  FROM (
    SELECT 'mailer.host' AS klucz, 'Serwer poczty wychodzącej' AS nazwa,
           'Adres serwera SMTP, którym platforma nadaje listy systemowe. '
           || 'Bez niego rejestracja i odzyskanie konta odmawiają, bo droga potwierdzenia '
           || 'nie ma czym dojść do Operatora.' AS opis,
           'string' AS rodzaj, '' AS domyslna,
           'np. smtp.twojadomena.pl' AS podpowiedz, 10 AS kolejnosc
    UNION ALL SELECT 'mailer.port', 'Port serwera poczty',
           'Port SMTP. Pusto znaczy 587 — port zgłoszenia z szyfrowaniem STARTTLS.',
           'int', '', '587', 20
    UNION ALL SELECT 'mailer.address', 'Adres nadawcy',
           'Adres w kopercie zwrotnej listu. Serwer poczty odrzuca list bez niego.',
           'string', '', 'np. platforma@twojadomena.pl', 30
    UNION ALL SELECT 'mailer.displayName', 'Nazwa wyświetlana nadawcy',
           'Nazwa widoczna przy adresie nadawcy w skrzynce odbiorcy.',
           'string', '', 'np. Danaco Console', 40
    UNION ALL SELECT 'mailer.username', 'Użytkownik serwera poczty',
           'Nazwa logowania do serwera poczty wychodzącej. Pusto znaczy serwer bez logowania '
           || '— tak stoi przekaźnik na tej samej maszynie.',
           'string', '', '', 50
    UNION ALL SELECT 'mailer.secret', 'Hasło do serwera poczty',
           'Hasło konta nadawczego. Wartość jest tajemnicą: okno Konfiguracji pokazuje ją zakrytą.',
           'secret', '', '', 60
    UNION ALL SELECT 'mailer.startTLS', 'Szyfruj połączenie (STARTTLS)',
           'Czy rozmowa z serwerem poczty ma zostać zaszyfrowana. Zdjęcie tego ma sens wyłącznie '
           || 'dla przekaźnika na tej samej maszynie i jest zapisem jawnym, nie przeoczeniem.',
           'bool', 'true', '', 70
  ) AS wykaz
  JOIN kategoria_ustawien kat ON kat.kod = 'nadawca'
ON CONFLICT(klucz) DO NOTHING;

INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
  JOIN poziom_zasiegu p ON p.kod = 'aplikacja'
 WHERE d.klucz IN ('mailer.host', 'mailer.port', 'mailer.address', 'mailer.displayName',
                   'mailer.username', 'mailer.secret', 'mailer.startTLS')
ON CONFLICT DO NOTHING;
