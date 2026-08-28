-- Wstawia do katalogu ustawień pozycję wymogu logowania w kategorii bezpieczeństwa, dostępną wyłącznie na poziomie zasięgu aplikacja, z wartością domyślną pustą.

INSERT INTO definicja_ustawienia (klucz, kategoria_id, nazwa, opis, rodzaj_wartosci,
                                  wartosc_domyslna, podpowiedz, wymaga_restartu, kolejnosc)
SELECT 'gateway.requireLogin', kat.id, 'Wymóg logowania',
       'Czy gniazdo musi przedstawić token sesji bramki, zanim wykona cokolwiek poza '
       || 'connection.hello, auth.login i auth.register. Bez wskazania rozstrzyga adres '
       || 'nasłuchu: pętla zwrotna bez wymogu, nasłuch szerszy z wymogiem.',
       'bool', '',
       'Zmiana obowiązuje od następnego połączenia — rdzenia nie trzeba zatrzymywać.',
       0, 10
  FROM kategoria_ustawien kat
 WHERE kat.kod = 'bezpieczenstwo'
ON CONFLICT(klucz) DO NOTHING;

INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
  JOIN poziom_zasiegu p ON p.kod = 'aplikacja'
 WHERE d.klucz = 'gateway.requireLogin'
ON CONFLICT DO NOTHING;
