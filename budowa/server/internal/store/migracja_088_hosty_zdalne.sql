-- Hosty zdalne wykonania: wykaz maszyn, na których rdzeń może uruchomić proces,
-- wraz z parametrami połączenia i zgodą wydaną osobno dla każdej z nich.
--
-- Torem jest SSH, nie własny agent sieciowy:
--   · SSH daje semantykę procesu w całości — strumienie wejścia, wyjścia
--     i diagnostyki, kod wyjścia, zakończenie po zerwaniu połączenia — czyli
--     kształt `session.UchwytProcesu`, bez ani jednej linii własnego protokołu;
--   · SSH daje uwierzytelnienie i szyfrowanie, których transport rdzenia nie ma,
--     więc tor nie otwiera nowej drogi wejścia, tylko jedzie usługą, którą host
--     zdalny już wystawia;
--   · proces transportu (`ssh`) startuje w tym samym jedynym spawnerze
--     (`injection.Wystartuj`), a po stronie zdalnej proces uruchamia `sshd` —
--     istniejąca usługa systemowa. Rola `agent` komponuje się z tym torem
--     wprost: `ssh host danaco-console --role agent` wykonuje żądania kontraktu
--     na hoście zdalnym tym samym produktem.
--
-- Kolumna `zgoda` ma DEFAULT 0: samo wpisanie hosta jeszcze niczego nie otwiera.
-- Rdzeń nie zainicjuje połączenia SSH z maszyną, której zgoda nie obejmuje,
-- i nie zapyta interaktywnie o hasło (tor jedzie z BatchMode=yes).
--
-- Wiersz hosta zakłada się poleceniem:
--
--   INSERT INTO host_zdalny (nazwa, adres, uzytkownik, port, zgoda, zgode_wydano)
--   VALUES ('danaco-system', 'danaco-system.example', 'operator', 22, 1,
--           strftime('%Y-%m-%dT%H:%M:%fZ','now'));
--
-- a zgodę cofa: UPDATE host_zdalny SET zgoda = 0 WHERE nazwa = '…';
-- Klucz publiczny rdzenia musi stać na hoście w authorized_keys — ta część
-- zgody zapada po stronie hosta, nie w tej tabeli; sekrety nie idą przez bazę.
--
-- `nazwa` odpowiada dosłownie wartości ustawienia `host_wykonania` — tak wybór
-- w oknie komunikacji dochodzi do tego wiersza. `adres` jest adresem sieciowym
-- dla SSH; pusty znaczy: adresem jest nazwa.

CREATE TABLE host_zdalny (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Nazwa hosta z ustawienia `host_wykonania`; klucz dopasowania wyboru.
    nazwa          TEXT    NOT NULL UNIQUE,
    -- Adres sieciowy dla SSH; pusty znaczy: łącz się z nazwą.
    adres          TEXT    NOT NULL DEFAULT '',
    -- Konto na hoście zdalnym; puste znaczy: konto procesu rdzenia.
    uzytkownik     TEXT    NOT NULL DEFAULT '',
    port           INTEGER NOT NULL DEFAULT 22 CHECK(port BETWEEN 1 AND 65535),
    -- Zgoda na inicjowanie połączeń z tym hostem; domyślnie brak.
    zgoda          INTEGER NOT NULL DEFAULT 0 CHECK(zgoda IN (0,1)),
    -- Znacznik chwili wydania zgody; pusty, póki zgody nie wydano.
    zgode_wydano   TEXT    NOT NULL DEFAULT '',
    utworzono      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
