-- Migracja 009 — zaczyn katalogu akcji.
--
-- Źródło wierszy. `shared/contract.json` — sekcja `komendy` (nazwa komendy,
-- opis, pola obowiązkowe żądania) oraz sekcja `narzedzia`. Zaczyn obejmuje te
-- komendy, które kontrakt wskazuje jako sterowanie platformą; poza wykazem
-- zostają `connection.hello` i `session.bind` — czynności warstwy połączenia
-- klienta, nie akcje panelu. Osobno wchodzi pasek narzędzi promptu każdego
-- modułu: każdy moduł niesie Chat Window, a jego pasek promptu niesie wysłanie
-- polecenia, zatrzymanie odpowiedzi i historię poleceń.
--
-- Akcja wskazująca komendę spoza `shared/contract.json` byłaby pozycją, której
-- nie da się wywołać, więc do zaczynu nie wchodzi. Takie akcje dochodzą
-- wierszami, bez zmiany kodu, gdy ich komendy wejdą do kontraktu.
--
-- Zasięg. Poziom wyprowadzony z bytu, na którym komenda działa:
-- `environment.enter` działa na środowisku, komendy sesji i kolejek na karcie
-- sesji, komendy okna na oknie komunikacji, komendy wiadomości na oknie
-- czatu modułu, reszta na całej platformie.
--
-- Każde wstawienie kończy się ON CONFLICT(kod) DO NOTHING — migracja przechodzi
-- także na bazie, w której część wierszy już jest. Klauzula `WHERE true` przed
-- ON CONFLICT jest wymogiem składni SQLite dla INSERT ... SELECT z upsertem.

-- ── Akcje o zasięgu poziomu (globalny · środowisko · karta sesji · okno) ──────
WITH katalog(kod, nazwa, opis, ikona, poziom, komenda, warunek, kolejnosc) AS (
    VALUES
        -- Poziom globalny — nawigacja platformy, sesje, konfiguracja, kanały.
        ('home.enter',        'Strona główna',          'Wejście na stronę główną; zwraca karty środowisk i sesje czynne konta', 'dom',           'globalny',   'home.enter',        '',           1),
        ('environment.list',  'Środowiska',             'Zwraca środowiska platformy',                                           'warstwy',       'globalny',   'environment.list',  '',           2),
        ('module.list',       'Moduły',                 'Zwraca moduły wraz z katalogiem ich okien operacyjnych',                'menu',          'globalny',   'module.list',       '',           3),
        ('session.create',    'Nowa sesja',             'Zakłada sesję',                                                         'plus',          'globalny',   'session.create',    '',           4),
        ('session.list',      'Sesje',                  'Zwraca listę sesji',                                                    'tabela-danych', 'globalny',   'session.list',      '',           5),
        ('window.list',       'Okna komunikacji',       'Zwraca okna komunikacji',                                               'karta-okna',    'globalny',   'window.list',       '',           6),
        ('config.get',        'Konfiguracja',           'Odczytuje konfigurację',                                                'ustawienia',    'globalny',   'config.get',        '',           7),
        ('config.set',        'Zapisz ustawienie',      'Zapisuje ustawienie na wskazanym poziomie zasięgu',                     'olowek',        'globalny',   'config.set',        '',           8),
        ('config.reset',      'Przywróć domyślne',      'Przywraca wartość domyślną; brak ustawienia znaczy wartość domyślną',   'odswiez',       'globalny',   'config.reset',      '',           9),
        ('channel.add',       'Nowy kanał modelu',      'Dodaje wiersz do rejestru kanałów modelu',                              'plus',          'globalny',   'channel.add',       '',          10),
        ('channel.update',    'Zmień kanał modelu',     'Zmienia wiersz rejestru kanałów',                                       'olowek',        'globalny',   'channel.update',    'kanal',     11),
        ('channel.remove',    'Usuń kanał modelu',      'Usuwa wiersz rejestru kanałów',                                         'kosz',          'globalny',   'channel.remove',    'kanal',     12),
        ('channel.list',      'Kanały modelu',          'Zwraca rejestr kanałów modelu',                                         'siec',          'globalny',   'channel.list',      '',          13),

        -- Poziom środowiska.
        ('environment.enter', 'Wejdź do środowiska',    'Wchodzi do środowiska; zwraca jego nawigację i karty sesji',            'strzalka-prawo', 'srodowisko', 'environment.enter', 'srodowisko', 1),

        -- Poziom karty sesji.
        ('session.focus',     'Przenieś ognisko',       'Przenosi ognisko na wskazaną kartę sesji i opcjonalnie na okno w jej wnętrzu', 'oko',      'karta_sesji', 'session.focus',   'sesja',      1),
        ('session.open',      'Otwórz sesję',           'Otwiera sesję wraz z jej oknami',                                       'folder',        'karta_sesji', 'session.open',    'sesja',      2),
        ('session.close',     'Zamknij sesję',          'Zamyka sesję',                                                          'zamknij',       'karta_sesji', 'session.close',   'sesja',      3),
        ('session.delete',    'Usuń sesję',             'Usuwa sesję',                                                           'kosz',          'karta_sesji', 'session.delete',  'sesja',      4),
        ('workspace.enter',   'Przełącz moduł',         'Przeładowuje przestrzeń roboczą karty sesji na wskazany moduł',         'karta-okna',    'karta_sesji', 'workspace.enter', 'sesja',      5),
        ('window.create',     'Nowe okno komunikacji',  'Zakłada okno komunikacji w sesji',                                      'plus',          'karta_sesji', 'window.create',   'sesja',      6),
        ('queue.create',      'Nowa kolejka',           'Zakłada kolejkę na jednym silniku pętli',                               'automatyzacja', 'karta_sesji', 'queue.create',    'sesja',      7),
        ('queue.action',      'Działanie na kolejce',   'Wykonuje działanie na kolejce',                                         'uruchom',       'karta_sesji', 'queue.action',    'kolejka',    8),

        -- Poziom okna komunikacji.
        ('window.state.get',  'Stan okna',              'Zwraca stan okna komunikacji: parametry wykonania, stan procesu i konfigurację efektywną', 'monitor', 'okno', 'window.state.get', 'okno', 1),
        ('window.update',     'Zmień parametry okna',   'Zmienia parametry okna komunikacji',                                    'olowek',        'okno',       'window.update',     'okno',       2),
        ('window.close',      'Zamknij okno',           'Zamyka okno komunikacji',                                               'zamknij',       'okno',       'window.close',      'okno',       3),
        ('context.transfer',  'Przenieś kontekst',      'Przenosi komplet kontekstu między modułami jedną komendą',              'wezly',         'okno',       'context.transfer',  'okno',       4)
)
INSERT INTO akcja (kod, nazwa, opis, ikona, poziom_zasiegu_id, klucz_zasiegu,
                   komenda, warunek_dostepnosci, kolejnosc, aktywna)
SELECT katalog.kod, katalog.nazwa, katalog.opis, katalog.ikona, p.id, '',
       katalog.komenda, katalog.warunek, katalog.kolejnosc, 1
  FROM katalog
  JOIN poziom_zasiegu p ON p.kod = katalog.poziom
 WHERE true
ON CONFLICT(kod) DO NOTHING;

-- ── Akcje paska narzędzi promptu — po jednym komplecie na każdy moduł ────────
-- Wiersze powstają złączeniem z tabelą `modul`: dopisanie modułu wierszem daje
-- jego akcje bez zmiany tego pliku.
WITH pasek(przyrostek, nazwa, opis, ikona, komenda, warunek, kolejnosc) AS (
    VALUES
        ('message.send', 'Wyślij',              'Wysyła wiadomość do okna komunikacji',                            'wyslij',    'message.send', 'okno', 1),
        ('message.stop', 'Zatrzymaj odpowiedź', 'Zatrzymuje bieżącą odpowiedź; przycisk zatrzymania zawsze czynny', 'zatrzymaj', 'message.stop', 'okno', 2),
        ('message.list', 'Historia okna',       'Zwraca wiadomości okna komunikacji',                              'historia',  'message.list', 'okno', 3)
)
INSERT INTO akcja (kod, nazwa, opis, ikona, poziom_zasiegu_id, klucz_zasiegu,
                   komenda, warunek_dostepnosci, kolejnosc, aktywna)
SELECT m.kod || '.' || pasek.przyrostek, pasek.nazwa, pasek.opis, pasek.ikona,
       p.id, m.kod, pasek.komenda, pasek.warunek, pasek.kolejnosc, 1
  FROM pasek
  CROSS JOIN modul m
  JOIN poziom_zasiegu p ON p.kod = 'modul'
 WHERE true
ON CONFLICT(kod) DO NOTHING;
