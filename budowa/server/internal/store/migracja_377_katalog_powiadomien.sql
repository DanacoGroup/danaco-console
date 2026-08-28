-- Migracja 377 — sekcja „Powiadomienia" okna Ustawień wchodzi do katalogu.
--
-- Podstawa: `interfejs-uzytkownika/ustawienia.md` rozdz. 7 (przełącznik główny,
-- siedem klas zdarzeń, kanały dostarczenia) oraz `architektura/model-danych.md`
-- rozdz. 18.4, gdzie uwaga projektowa mówi wprost: „Zakres klas zdarzeń
-- zgłaszanych przez centrum powiadomień oraz kanał dostarczenia są USTAWIENIAMI
-- KONFIGURACYJNYMI warstwy globalnej i warstwy środowiska".
--
-- To zdanie rozstrzyga o kształcie tej roboty. Sekcja Powiadomień nie potrzebuje
-- nowej rodziny kontraktu: jest macierzą nastaw i idzie tą samą drogą, co każde
-- inne ustawienie platformy — `config.get`, `config.set`, `settings.definition.list`,
-- zdarzenie `config.changed`. Nowej komendy nie wnosimy, bo nie ma czego wnosić.
--
-- Kody klas są kodami z modelu danych (`powiadomienie.klasa`: zakonczenie,
-- decyzja, blad, wzmianka, termin, automatyka, system) — jeden zapis na całą
-- platformę. Drugi zestaw nazw, choćby ładniejszy, rozjechałby ustawienie
-- z wierszem, którego dotyczy.
--
-- Kanał „centrum" nie jest tu wyborem i nie ma dla niego kolumny. Rozdz. 7.3
-- mówi, że centrum powiadomień jest kanałem PODSTAWOWYM każdej klasy: „każde
-- zdarzenie objęte ustawieniem trafia do rejestru centrum niezależnie od
-- pozostałych kanałów". Nastawa wybiera więc kanały DODATKOWE; centrum stoi
-- zawsze, dopóki klasa jest czynna.
--
-- Zasięgi: `globalny`, `srodowisko` i `karta_sesji`. Rozdz. 7.1 wymienia warstwę
-- globalną i sesji, rozdz. 18.4 — globalną i środowiska. Dopuszczenie trzech
-- poziomów spełnia oba zapisy; zawężenie do dwóch wybierałoby, który z dwóch
-- dokumentów dostawy jest ważniejszy, a to nie jest rozstrzygnięcie wykonawcy.
--
-- Wymaga restartu: nie. Nastawy czyta się drogą `config.get` przy otwarciu
-- sekcji, a zmiany dolatują zdarzeniem `config.changed`.

INSERT INTO kategoria_ustawien (kod, nazwa, opis, ikona, kolejnosc) VALUES
    ('powiadomienia', 'Powiadomienia',
     'Zakres klas zdarzeń zgłaszanych przez centrum powiadomień oraz kanały dostarczenia',
     'dzwonek', 10)
ON CONFLICT(kod) DO NOTHING;

-- ── przełącznik główny ──────────────────────────────────────────────────────
-- Wyłączenie wygasza wszystkie klasy naraz, ZACHOWUJĄC ich ustawienia (rozdz.
-- 7.5). Dlatego jest osobnym kluczem, a nie zapisem „false" do siedmiu kluczy
-- klas: tamto skasowałoby wybór Operatora i ponowne włączenie przywróciłoby
-- stan domyślny zamiast poprzedniego.
INSERT INTO definicja_ustawienia (klucz, kategoria_id, nazwa, opis, rodzaj_wartosci,
                                  wartosc_domyslna, podpowiedz, wymaga_restartu, kolejnosc)
SELECT 'powiadomienia.wlaczone', kat.id, 'Powiadomienia aktywne',
       'Przełącznik główny centrum powiadomień. Wyłączenie wygasza wszystkie klasy zdarzeń '
       || 'naraz i zachowuje ich ustawienia — ponowne włączenie przywraca stan poprzedni, '
       || 'nie domyślny. Wyłączenie nie ogranicza dostępu do żadnej funkcji platformy.',
       'bool', 'true', '', 0, 10
  FROM kategoria_ustawien kat
 WHERE kat.kod = 'powiadomienia'
ON CONFLICT(klucz) DO NOTHING;

-- ── siedem klas zdarzeń ─────────────────────────────────────────────────────
-- Kolejność i opisy wprost z tabeli rozdz. 7.2. Wszystkie czynne domyślnie.
INSERT INTO definicja_ustawienia (klucz, kategoria_id, nazwa, opis, rodzaj_wartosci,
                                  wartosc_domyslna, podpowiedz, wymaga_restartu, kolejnosc)
SELECT klasy.klucz, kat.id, klasy.nazwa, klasy.opis, 'bool', 'true', '', 0, klasy.kolejnosc
  FROM (
    SELECT 'powiadomienia.klasa.zakonczenie' AS klucz, 'Zakończenie' AS nazwa,
           'Zakończenie przebiegu pętli wykonawczej, zakończenie procesu automatyki.' AS opis,
           20 AS kolejnosc
    UNION ALL SELECT 'powiadomienia.klasa.decyzja', 'Decyzja',
           'Krok procesu oczekujący na zatwierdzenie.', 30
    UNION ALL SELECT 'powiadomienia.klasa.blad', 'Błąd',
           'Niepowodzenie zadania, naruszenie zależności orkiestracji, powtarzające się '
           || 'niepowodzenie walidacji.', 40
    UNION ALL SELECT 'powiadomienia.klasa.wzmianka', 'Wzmianka',
           'Odwołanie do Operatora w treści, komentarz przypisany do artefaktu.', 50
    UNION ALL SELECT 'powiadomienia.klasa.termin', 'Termin',
           'Termin zadania projektu, zbliżający się cykl harmonogramu.', 60
    UNION ALL SELECT 'powiadomienia.klasa.automatyka', 'Automatyka',
           'Wpięcie automatyki, wynik cyklu harmonogramu.', 70
    UNION ALL SELECT 'powiadomienia.klasa.system', 'System',
           'Zmiany uwierzytelnienia i urządzeń, zmiany konfiguracji, pamięci i historii, '
           || 'wiadomości oraz zdolności telefonu.', 80
  ) AS klasy
  JOIN kategoria_ustawien kat ON kat.kod = 'powiadomienia'
ON CONFLICT(klucz) DO NOTHING;

-- ── kanały dodatkowe każdej klasy ───────────────────────────────────────────
-- Wartość domyślna wprost z makiety 5 rozdz. 7.4: Mobile przy trzech klasach
-- pilnych, e-mail przy klasie System, reszta samo centrum.
INSERT INTO definicja_ustawienia (klucz, kategoria_id, nazwa, opis, rodzaj_wartosci,
                                  wartosc_domyslna, podpowiedz, wymaga_restartu, kolejnosc)
SELECT kanaly.klucz, kat.id, kanaly.nazwa,
       'Kanały dodatkowe tej klasy zdarzeń. Centrum powiadomień jest kanałem podstawowym '
       || 'i stoi zawsze, dopóki klasa jest czynna — tu wybiera się drogi obok niego.',
       'enumList', kanaly.domyslna, '', 0, kanaly.kolejnosc
  FROM (
    SELECT 'powiadomienia.klasa.zakonczenie.kanaly' AS klucz,
           'Zakończenie — kanały' AS nazwa, '["mobile"]' AS domyslna, 21 AS kolejnosc
    UNION ALL SELECT 'powiadomienia.klasa.decyzja.kanaly', 'Decyzja — kanały', '["mobile"]', 31
    UNION ALL SELECT 'powiadomienia.klasa.blad.kanaly', 'Błąd — kanały', '["mobile"]', 41
    UNION ALL SELECT 'powiadomienia.klasa.wzmianka.kanaly', 'Wzmianka — kanały', '[]', 51
    UNION ALL SELECT 'powiadomienia.klasa.termin.kanaly', 'Termin — kanały', '[]', 61
    UNION ALL SELECT 'powiadomienia.klasa.automatyka.kanaly', 'Automatyka — kanały', '[]', 71
    UNION ALL SELECT 'powiadomienia.klasa.system.kanaly', 'System — kanały', '["email"]', 81
  ) AS kanaly
  JOIN kategoria_ustawien kat ON kat.kod = 'powiadomienia'
ON CONFLICT(klucz) DO NOTHING;

-- Dwa kanały dodatkowe, wspólne wszystkim klasom (rozdz. 7.3).
INSERT INTO opcja_ustawienia (definicja_id, wartosc, etykieta, opis, kolejnosc)
SELECT d.id, o.wartosc, o.etykieta, o.opis, o.kolejnosc
  FROM definicja_ustawienia d
  CROSS JOIN (
    SELECT 'mobile' AS wartosc, 'Mobile' AS etykieta,
           'Powiadomienie wypychane na telefon — kanałem własnym przy aplikacji na pierwszym '
           || 'planie, kanałem pośrednika poza nią.' AS opis, 10 AS kolejnosc
    UNION ALL SELECT 'email', 'E-mail',
           'List na adres e-mail uwierzytelniający konta.', 20
  ) AS o
 WHERE d.klucz LIKE 'powiadomienia.klasa.%.kanaly'
ON CONFLICT(definicja_id, wartosc) DO NOTHING;

INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
  JOIN poziom_zasiegu p ON p.kod IN ('globalny', 'srodowisko', 'karta_sesji')
 WHERE d.klucz = 'powiadomienia.wlaczone'
    OR d.klucz LIKE 'powiadomienia.klasa.%'
ON CONFLICT DO NOTHING;
