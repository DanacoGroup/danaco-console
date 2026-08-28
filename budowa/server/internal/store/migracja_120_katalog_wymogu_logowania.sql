-- Migracja 120 — przełącznik „Wymóg logowania" wchodzi do katalogu ustawień.
--
-- Migracja 119 dołożyła poziom zasięgu `aplikacja` — miejsce, w którym nastawa
-- może zamieszkać. Sam poziom nie wystarcza: `config.set` sprawdza klucz wobec
-- katalogu (`core/walidacja_klucza_ustawienia.go`), a katalog żyje w tabeli
-- `definicja_ustawienia` z migracji 012 i przesłania rejestr wbudowany rdzenia.
-- Bez wiersza katalogu `config.set` z kluczem `gateway.requireLogin` odmawia
-- („klucz nie należy do katalogu ustawień"); odmowa jest prawdziwa, bo cichego
-- zapisu bez skutku nie ma. Ten krok dokłada brakujący wiersz.
--
-- Kategoria `bezpieczenstwo`, bo tam już mieszka uwierzytelnianie: migracja 012
-- opisała ją wprost („Zakres zgody wydanej modelowi oraz egzekwowanie
-- uwierzytelniania"). Nowa kategoria byłaby drugim miejscem na to samo pojęcie
-- w oknie Ustawień.
--
-- Wartość domyślna pusta, a nie „false": wymóg logowania ma trzy stany, nie dwa
-- (`transport/ustawienia.go`): „nie wskazałem — rozstrzyga adres nasłuchu",
-- „wskazałem: tak", „wskazałem: nie". Wartość domyślna `false` skasowałaby stan
-- pierwszy i po cichu zniosłaby wymóg na nasłuchu wystawionym poza pętlę zwrotną
-- — dźwignia zamieniłaby się w otwarcie drzwi przez zapomnienie. Pusto znaczy
-- „bez wskazania".
--
-- Wymaga restartu: nie. Warstwa nasłuchu składa straż bramki raz na połączenie
-- (`transport/nawiazanie.go`), a rdzeń czyta nastawę przy powitaniu
-- `connection.hello` (`core/nastawy_aplikacji.go`). Nowa wartość obowiązuje więc
-- od następnego połączenia, bez zatrzymywania rdzenia — po to ten poziom powstał.
--
-- Jeden dozwolony zasięg: wyłącznie `aplikacja`. Zapis tego klucza na poziomie
-- okna albo sesji nie znaczyłby nic: nasłuch jest jeden dla całego programu
-- i nie ma czym go zawęzić.

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
