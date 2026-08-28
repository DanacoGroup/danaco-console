-- Migracja 075 daje silnikowi mowy sterowanie wpisami katalogu ustawień oraz
-- ślad w dzienniku transkrypcji. Silnik jest lokalny i liczy na procesorze;
-- pełne uzasadnienie stoi w dokumentacji architektury zaplecza.

INSERT INTO kategoria_ustawien (kod, nazwa, opis, ikona, kolejnosc) VALUES
    ('mowa', 'Mowa na tekst', 'Lokalny silnik rozpoznawania mowy: interpreter, rozmiar modelu, język i katalog modeli', 'mikrofon', 9)
ON CONFLICT(kod) DO NOTHING;

-- Rodzaj wartości każdej pozycji pochodzi z wyliczenia SettingValueType kontraktu:
-- ścieżki są rodzaju path, kod języka string, a rozmiar modelu enum z opcjami niżej.
WITH katalog(klucz, kategoria, nazwa, opis, rodzaj, domyslna, podpowiedz,
             wymaga_restartu, kolejnosc) AS (
    VALUES
        ('mowa_program', 'mowa', 'Interpreter Pythona',
         'Ścieżka interpretera Pythona uruchamiającego pomocnika transkrypcji. Pusta znaczy: szukaj `python3` na ścieżce wyszukiwania systemu, a nie odmawiaj uruchomienia.',
         'path', '', 'python3', 1, 1),
        ('mowa_model', 'mowa', 'Rozmiar modelu',
         'Rozmiar modelu rozpoznawania — od najszybszego do najdokładniejszego. Model większy przepisuje wierniej i wolniej; silnik liczy na procesorze, więc wybór jest wprost wymianą czasu na dokładność.',
         'enum', 'small', '', 0, 2),
        ('mowa_jezyk', 'mowa', 'Język rozpoznawania',
         'Kod języka, w którym rozpoznawana jest mowa. Wskazanie języka z góry jest szybsze i wierniejsze niż zdanie się na rozpoznanie automatyczne.',
         'string', 'pl', 'pl', 0, 3),
        ('mowa_katalog_modeli', 'mowa', 'Katalog modeli',
         'Katalog, do którego pobierane są modele rozpoznawania. Pusty znaczy: katalog domyślny biblioteki silnika. Modele bywają wielkie, więc Operator może je przenieść na dysk pojemniejszy.',
         'path', '', '', 1, 4)
)
INSERT INTO definicja_ustawienia (klucz, kategoria_id, nazwa, opis, rodzaj_wartosci,
                                  wartosc_domyslna, podpowiedz, wymaga_restartu, kolejnosc)
SELECT k.klucz, kat.id, k.nazwa, k.opis, k.rodzaj, k.domyslna, k.podpowiedz,
       k.wymaga_restartu, k.kolejnosc
  FROM katalog k
  JOIN kategoria_ustawien kat ON kat.kod = k.kategoria
 WHERE true
ON CONFLICT(klucz) DO NOTHING;

-- Pięć rozmiarów wydawanych przez bibliotekę silnika. Wartości pustej nie ma:
-- silnik musi dostać nazwę modelu, więc brak wskazania byłby odmową uruchomienia.
WITH opcje(wartosc, etykieta, opis, kolejnosc) AS (
    VALUES
        ('tiny',     'Najmniejszy', 'Najszybciej, najmniej dokładnie',            1),
        ('base',     'Mały',        'Szybciej, kosztem dokładności',              2),
        ('small',    'Średni',      'Równowaga szybkości i dokładności',          3),
        ('medium',   'Duży',        'Dokładniej, wolniej',                        4),
        ('large-v3', 'Największy',  'Najdokładniej, najwolniej i najwięcej pamięci', 5)
)
INSERT INTO opcja_ustawienia (definicja_id, wartosc, etykieta, opis, kolejnosc)
SELECT d.id, o.wartosc, o.etykieta, o.opis, o.kolejnosc
  FROM opcje o
  JOIN definicja_ustawienia d ON d.klucz = 'mowa_model'
 WHERE true
ON CONFLICT(definicja_id, wartosc) DO NOTHING;

INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
 CROSS JOIN poziom_zasiegu p
 WHERE d.klucz IN ('mowa_program', 'mowa_katalog_modeli')
   AND p.kod = 'globalny'
ON CONFLICT DO NOTHING;

INSERT INTO definicja_ustawienia_zasieg (definicja_id, poziom_zasiegu_id)
SELECT d.id, p.id
  FROM definicja_ustawienia d
 CROSS JOIN poziom_zasiegu p
 WHERE d.klucz IN ('mowa_model', 'mowa_jezyk')
ON CONFLICT DO NOTHING;

-- Oś jest wyłącznie platform: rozmiar modelu i ścieżka Pythona są własnością
-- maszyny operatora, niezależną od modelu odpowiadającego ani od konta.
INSERT INTO definicja_ustawienia_os (definicja_id, os)
SELECT d.id, 'platform'
  FROM definicja_ustawienia d
 WHERE d.klucz IN ('mowa_program', 'mowa_model', 'mowa_jezyk', 'mowa_katalog_modeli')
ON CONFLICT DO NOTHING;

CREATE TABLE transkrypcja (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Tożsamość wpisu widoczna na zewnątrz; klucz sztuczny nigdy nie wychodzi w odpowiedzi.
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- Okno, w którym padło zlecenie; wartość pusta znaczy zlecenie spoza okna.
    okno_id                  TEXT,
    -- Ścieżka pliku na dysku operatora, nie bajty dźwięku.
    nagranie_odnosnik        TEXT    NOT NULL,
    -- Rozmiar modelu i język zapisane tak, jak obowiązywały w chwili zlecenia.
    model                    TEXT    NOT NULL,
    jezyk                    TEXT    NOT NULL,
    znakow                   INTEGER NOT NULL DEFAULT 0,
    -- Długość nagrania w milisekundach, nie czas przetwarzania.
    trwanie_ms               INTEGER NOT NULL DEFAULT 0,
    -- Trzeci stan bez_mowy nazywa nagranie ciszy albo szumu, odrębnie od gotowa i odmowa.
    stan                     TEXT    NOT NULL CHECK(stan IN ('gotowa', 'bez_mowy', 'odmowa')),
    -- Nazwany powód odmowy; przy stanach udanych nie ma czego opisywać.
    powod                    TEXT,
    -- Milisekundy epoki podane przez wołającego, nie wyrażeniem czasu bazy.
    utworzono                INTEGER NOT NULL,
    -- Wiąże stan z powodem w obie strony: odmowa musi nieść powód, wpis gotowy nie ma prawa.
    CHECK((stan = 'odmowa' AND powod IS NOT NULL AND powod <> '')
          OR (stan IN ('gotowa', 'bez_mowy') AND powod IS NULL))
);

-- Kolumny biegną porządkiem zapytania wykazu z mowa/dziennik.go, więc SQLite
-- czyta indeks wstecz zamiast sortować; kolumny stan w indeksie nie ma celowo.
CREATE INDEX idx_transkrypcja_wykaz ON transkrypcja(okno_id, utworzono, id);
