-- Migracja 376 wprowadza do katalogu ustawień okna Konfiguracji zakres Pamięć
-- oraz reguły wyciszania nakładki Always On Display.

INSERT INTO kategoria_ustawien (kod, nazwa, opis, ikona, kolejnosc) VALUES
    ('pamiec', 'Pamięć',
     'Poziom pamięci kontekstowej, stan jej włączenia na poziomie i odłączenie pamięci w sesji',
     'warstwy', 10),
    ('aod', 'Always On Display',
     'Nakładka stałej obecności asystenta: reguły wyciszania sugestii',
     'karta-okna', 11)
ON CONFLICT(kod) DO NOTHING;

INSERT INTO definicja_ustawienia (klucz, kategoria_id, nazwa, opis, rodzaj_wartosci,
                                  wartosc_domyslna, podpowiedz, wymaga_restartu, kolejnosc)
SELECT wykaz.klucz, kat.id, wykaz.nazwa, wykaz.opis, wykaz.rodzaj,
       wykaz.domyslna, wykaz.podpowiedz, 0, wykaz.kolejnosc
  FROM (
    SELECT 'memory.level' AS klucz, 'pamiec' AS kategoria, 'Poziom pamięci' AS nazwa,
           'Określa poziom, na którym zasób pamięci obowiązuje: globalny, środowisko, projekt '
           || 'albo karta sesji. Poziom węższy wygrywa z szerszym, a brak ustawienia znaczy '
           || 'dziedziczenie z poziomu bezpośrednio szerszego. Zmiana poziomu nie rusza treści '
           || 'żadnego wpisu — mówi wyłącznie, gdzie wpis obowiązuje.' AS opis,
           'enum' AS rodzaj, 'session' AS domyslna, '' AS podpowiedz, 10 AS kolejnosc
    UNION ALL SELECT 'memory.enabled', 'pamiec', 'Pamięć włączona na tym poziomie',
           'Określa, czy pamięć na tym poziomie wchodzi do kontekstu modelu. Wyłączenie NIE USUWA '
           || 'ani jednego wpisu: treść zostaje nietknięta, wpisy przestają wchodzić do kontekstu '
           || 'i wracają w całości po ponownym włączeniu. Wyłączenie zapisane na węższym zasięgu '
           || '— w środowisku, projekcie, module, parze modułów albo karcie sesji — dotyczy '
           || 'wyłącznie tego zasięgu; wszędzie indziej pamięć działa dalej.',
           'bool', 'true', '', 20
    UNION ALL SELECT 'memory.session.detached', 'pamiec', 'Pamięć odłączona w tej sesji',
           'Określa, czy ta karta sesji ma dostęp do pamięci. Odłączenie jest możliwe przed '
           || 'pierwszym promptem karty; po nim karta pracuje z kontekstem, który już zebrała. '
           || 'Odłączenie nie usuwa wpisów i nie zmienia niczego w pozostałych kartach.',
           'bool', 'false', '', 30
    UNION ALL SELECT 'aod.mute.rules', 'aod', 'Reguły wyciszania',
           'Określa zakresy i czasy wyciszenia dostępne w menu nakładki: wyciszenie czasowe '
           || '(kwadrans, godzina, do końca dnia), wyciszenie bieżącego modułu, wyciszenie '
           || 'bieżącej karty sesji i wyciszenie klasy zdarzeń. Zdjęcie pozycji z tego wykazu '
           || 'zabiera ją z menu i nie znosi wyciszeń już czynnych — te znosi się w nakładce '
           || 'jednym kliknięciem. Punkt decyzyjny wstrzymujący proces ujawnia się mimo każdego '
           || 'wyciszenia i tego nie zmienia żadna reguła.',
           'enumList', 'timed,module,session,eventClass', '', 10
  ) AS wykaz
  JOIN kategoria_ustawien kat ON kat.kod = wykaz.kategoria
ON CONFLICT(klucz) DO NOTHING;

-- Dopuszczalne wartości pozycji wyliczeniowych. Wartości poziomu pamięci są
-- wartościami kontraktu (`MemoryLevel`), nie polskimi kodami: przekład nazw
-- mieszka wyłącznie w kontrakcie.
INSERT INTO opcja_ustawienia (definicja_id, wartosc, etykieta, opis, kolejnosc)
SELECT d.id, wykaz.wartosc, wykaz.etykieta, wykaz.opis, wykaz.kolejnosc
  FROM (
    SELECT 'memory.level' AS klucz, 'global' AS wartosc, 'Globalna' AS etykieta,
           'Pamięć wspólna wszystkim projektom i oknom' AS opis, 10 AS kolejnosc
    UNION ALL SELECT 'memory.level', 'environment', 'Środowisko',
           'Pamięć środowiska: TalkIn, WorkSpace, CodeStudio, MultitaskingAI', 20
    UNION ALL SELECT 'memory.level', 'project', 'Projekt',
           'Pamięć projektu, do którego należy okno', 30
    UNION ALL SELECT 'memory.level', 'session', 'Karta sesji',
           'Pamięć jednej rozmowy', 40
    UNION ALL SELECT 'aod.mute.rules', 'timed', 'Wyciszenie czasowe',
           'Kwadrans, godzina, do końca dnia', 10
    UNION ALL SELECT 'aod.mute.rules', 'module', 'Wyciszenie bieżącego modułu',
           'Sugestie wskazanego modułu nie ujawniają się', 20
    UNION ALL SELECT 'aod.mute.rules', 'session', 'Wyciszenie bieżącej karty sesji',
           'Sugestie wskazanej karty sesji nie ujawniają się', 30
    UNION ALL SELECT 'aod.mute.rules', 'eventClass', 'Wyciszenie klasy zdarzeń',
           'Wskazana klasa nie tworzy sugestii do chwili zniesienia', 40
  ) AS wykaz
  JOIN definicja_ustawienia d ON d.klucz = wykaz.klucz
ON CONFLICT DO NOTHING;

-- Pamięć obejmuje cztery warstwy ogólne modelu konfiguracji, poszerzone
-- o moduł i parę modułów przy stanie włączenia, aby dało się wyłączyć pamięć
-- tylko dla wybranych modułów; odłączenie w sesji obejmuje wyłącznie kartę
-- sesji.
INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
  JOIN poziom_zasiegu p ON p.kod IN ('globalny', 'srodowisko', 'projekt', 'karta_sesji')
 WHERE d.klucz = 'memory.level'
ON CONFLICT DO NOTHING;

INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
  JOIN poziom_zasiegu p ON p.kod IN ('globalny', 'srodowisko', 'modul', 'para_modulow',
                                     'projekt', 'karta_sesji')
 WHERE d.klucz = 'memory.enabled'
ON CONFLICT DO NOTHING;

INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
  JOIN poziom_zasiegu p ON p.kod = 'karta_sesji'
 WHERE d.klucz = 'memory.session.detached'
ON CONFLICT DO NOTHING;

INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
  JOIN poziom_zasiegu p ON p.kod = 'globalny'
 WHERE d.klucz = 'aod.mute.rules'
ON CONFLICT DO NOTHING;
