-- Migracja 106 — pułap kosztu okna: nastawa zapobiegawcza zamiast wykrywania
-- po fakcie.
--
-- Rdzeń umie o koszcie dwie rzeczy i obie są wsteczne: czyta pole `Koszt`
-- z `costUsd` zdarzenia `result` (`injection/fragment.go`) — koszt tury, która
-- już się odbyła — i wykrywa granicę, o którą tura już uderzyła
-- (`injection/wyczerpanie.go`: „usage limit reached", „429", `rate_limit_event`).
-- Czego nie umie: powiedzieć programowi, ile wolno wydać. Wykaz przełączników
-- składanych w `injection/argumenty.go` nie niesie żadnego ograniczenia, choć
-- program `claude` przełącznik taki ma (`--max-budget-usd <amount>`).
--
-- Wykrycie po fakcie mówi Operatorowi, że pieniądze wydano; pułap sprawia, że nie
-- zostaną wydane. Pierwsze jest meldunkiem, drugie — nastawą.
--
-- Pułap jest nastawą okna, nie stałą produktu. Wchodzi katalogiem ustawień
-- (migracja 012), więc rozstrzyga się tymi samymi poziomami zasięgu, co nakład
-- rozumowania i kanał zapasowy: Operator ustawia go globalnie, na projekcie, na
-- karcie sesji albo na jednym oknie — i okno wygrywa. Drugiej drogi do tej
-- wartości nie ma.
--
-- Wartość domyślna 0 znaczy brak pułapu i tak ma znaczyć: bazy zastane po tej
-- migracji zachowują się jak przed nią — `injection/argumenty.go` nie dopisuje
-- wtedy przełącznika, więc wiersz argv nie zmienia się o bajt. Pułap zaczyna
-- działać w chwili, w której Operator wpisze liczbę dodatnią, i wyłącznie na tym
-- poziomie zasięgu, na którym ją wpisał.
--
-- Nastawa jest typu 'float', nie 'int' w groszach: przełącznik programu przyjmuje
-- kwotę w dolarach, a pole `costUsd` zdarzenia `result` jest liczbą
-- zmiennoprzecinkową w tej samej jednostce. Trzymanie nastawy w groszach
-- wymagałoby przeliczania w dwie strony na granicy z programem — a każde
-- przeliczenie jest miejscem, w którym jednostka może się zgubić. Jednostka jest
-- zapisana wprost w kolumnie `jednostka`, żeby panel jej nie zgadywał.
--
-- Odwracalność: krok jest dokładający — jeden wiersz katalogu i jego przypięcia
-- do poziomów i osi zasięgu. Cofnięcie to `DELETE FROM definicja_ustawienia
-- WHERE klucz='pulap_kosztu_usd'` — przypięcia znikają kaskadą, a wartości
-- zapisane przez Operatora w tabeli `ustawienie` przestają być czytane, bo
-- rozstrzygacz chodzi po definicjach. Każde wstawienie kończy się klauzulą
-- ON CONFLICT DO NOTHING, więc krok przechodzi także na bazie, w której wiersz
-- już jest.

-- ── 1. Pozycja katalogu ─────────────────────────────────────────────────────
--
-- Kategoria `modele`, bo pułap dotyczy wywołania modelu i stoi obok nakładu
-- rozumowania oraz kanału zapasowego — Operator szuka go tam, gdzie ustawia
-- resztę parametrów wywołania. Kolejność 4, zaraz za nakładem (3).
--
-- Minimum 0 i skok 0.5 opisują suwak panelu; maksimum zostaje puste, bo górnej
-- granicy wydatku nie wyznacza produkt — wyznacza ją Operator.
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

-- ── 2. Poziomy zasięgu, na których wolno pułap ustawić ──────────────────────
--
-- Wszystkie, którym wolno ustawić nakład rozumowania. Wybór nie jest tu
-- niezależnym rozstrzygnięciem: pułap jest parametrem tego samego wywołania,
-- więc rozdzielenie zasięgów zrobiłoby dwie różne geometrie dla dwóch pól
-- jednej komendy. Przepisanie z istniejącego wiersza jest zarazem odporne na
-- przyszłość — poziomy dopisane później nadążą same.
INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT nowa.id, p.poziom_zasiegu_id
  FROM definicja_ustawienia nowa
  JOIN definicja_ustawienia wzor ON wzor.klucz = 'naklad_rozumowania'
  JOIN definicja_ustawienia_zasieg p ON p.definicja_id = wzor.id
 WHERE nowa.klucz = 'pulap_kosztu_usd'
ON CONFLICT(definicja_id, poziom_zasiegu_id) DO NOTHING;

-- ── 3. Osie zasięgu — tym samym rozumowaniem, co poziomy ────────────────────
INSERT INTO definicja_ustawienia_os (definicja_id, os)
SELECT nowa.id, o.os
  FROM definicja_ustawienia nowa
  JOIN definicja_ustawienia wzor ON wzor.klucz = 'naklad_rozumowania'
  JOIN definicja_ustawienia_os o ON o.definicja_id = wzor.id
 WHERE nowa.klucz = 'pulap_kosztu_usd'
ON CONFLICT(definicja_id, os) DO NOTHING;
