-- Dodaje do katalogu ustawień pozycję pułapu kosztu okna, ograniczającą z góry koszt jednego wywołania modelu, wraz z poziomami i osiami zasięgu tego wywołania.

-- Wstawia do katalogu ustawień pozycję pułapu kosztu w kategorii modeli, z domyślną wartością zero oznaczającą brak ograniczenia wydatku.
INSERT INTO definicja_ustawienia (klucz, kategoria_id, nazwa, opis, rodzaj_wartosci,
                                  wartosc_domyslna, minimum, skok, jednostka,
                                  podpowiedz, wymaga_restartu, kolejnosc)
SELECT 'pulap_kosztu_usd', kat.id, 'Pułap kosztu okna',
       'Górna granica kosztu jednego wywołania modelu, przekazywana programowi '
       || 'przełącznikiem --max-budget-usd. Zero znaczy BRAK PUŁAPU — wywołanie '
       || 'idzie bez ograniczenia, dokładnie tak jak dotąd. Wartość dodatnia '
       || 'sprawia, że program przerywa turę PRZED przekroczeniem kwoty; rdzeń '
       || 'odróżnia takie przerwanie od awarii i zgłasza je jako wstrzymanie '
       || 'na pułapie, nie jako błąd.',
       'float', '0', 0, 0.5, 'USD', 'na przykład 2.5', 0, 4
  FROM kategoria_ustawien kat
 WHERE kat.kod = 'modele'
ON CONFLICT(klucz) DO NOTHING;

-- Wstawia dla pozycji pułapu kosztu te same poziomy zasięgu, na których wolno ustawić nakład rozumowania, przez przepisanie z istniejącego wiersza.
INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT nowa.id, p.poziom_zasiegu_id
  FROM definicja_ustawienia nowa
  JOIN definicja_ustawienia wzor ON wzor.klucz = 'naklad_rozumowania'
  JOIN definicja_ustawienia_zasieg p ON p.definicja_id = wzor.id
 WHERE nowa.klucz = 'pulap_kosztu_usd'
ON CONFLICT(definicja_id, poziom_zasiegu_id) DO NOTHING;

-- Wstawia dla pozycji pułapu kosztu te same osie zasięgu, na których wolno ustawić nakład rozumowania, tym samym rozumowaniem co poziomy zasięgu.
INSERT INTO definicja_ustawienia_os (definicja_id, os)
SELECT nowa.id, o.os
  FROM definicja_ustawienia nowa
  JOIN definicja_ustawienia wzor ON wzor.klucz = 'naklad_rozumowania'
  JOIN definicja_ustawienia_os o ON o.definicja_id = wzor.id
 WHERE nowa.klucz = 'pulap_kosztu_usd'
ON CONFLICT(definicja_id, os) DO NOTHING;
