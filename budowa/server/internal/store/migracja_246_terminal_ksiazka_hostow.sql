-- Migracja 246 — książka hostów modułu Terminal (okno Session Manager).
--
-- Do tej pory książka żyła wyłącznie w widoku klienta: wpisy ginęły przy
-- odświeżeniu strony, a jedyną drogą ich wyniesienia poza jedno posiedzenie był
-- wywóz do pliku. Wpis hosta jest jednak nastawą Operatora, nie stanem widoku —
-- ma przeżyć zamknięcie okna, zamknięcie przeglądarki i restart rdzenia.
--
-- Czego w wpisie NIE MA. Nie ma hasła, frazy klucza ani żadnego materiału
-- tajnego. Kolumna `klucz_kod` wskazuje wpis wykazu kluczy (migracja 248), a ten
-- niesie wyłącznie ŚCIEŻKĘ klucza prywatnego na maszynie rdzenia. Baza nie jest
-- sejfem i nie stanie się nim przez dołożenie kolumny.
--
-- `host_posredni_kod` wskazuje inny wiersz tej samej tabeli — host, przez który
-- idzie połączenie (ProxyJump OpenSSH). Klucz obcy jest na siebie samą i kasuje
-- się na NULL: usunięcie hosta pośredniego nie ma prawa unieważnić wpisów, które
-- przez niego szły; one wracają do połączenia bezpośredniego, a Operator to
-- widzi.

CREATE TABLE terminal_host (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    kod               TEXT    NOT NULL UNIQUE,
    -- Nazwa widoczna w wykazie; klucz rozpoznania dla Operatora.
    nazwa             TEXT    NOT NULL,
    -- Adres celu w postaci `użytkownik@host` albo alias konfiguracji OpenSSH
    -- maszyny RDZENIA — to na niej uruchamia się program `ssh`.
    cel               TEXT    NOT NULL,
    -- Port połączenia; pusty bierze port domyślny protokołu.
    port              INTEGER CHECK(port IS NULL OR port BETWEEN 1 AND 65535),
    -- Folder porządkujący wykaz; pusty znaczy „bez folderu”.
    grupa             TEXT    NOT NULL DEFAULT '',
    -- Katalog roboczy karty zakładanej z tego wpisu.
    katalog_roboczy   TEXT    NOT NULL DEFAULT '',
    -- Wpis wykazu kluczy SSH; pusty bierze klucz domyślny konfiguracji maszyny.
    klucz_kod         TEXT,
    -- Wpis hosta pośredniego, przez który idzie połączenie.
    host_posredni_kod TEXT REFERENCES terminal_host(kod) ON DELETE SET NULL,
    -- Notatka Operatora albo ślad pochodzenia wpisu (np. wczytanie z pliku).
    notatka           TEXT    NOT NULL DEFAULT '',
    utworzono         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_terminal_host_grupa ON terminal_host(grupa, nazwa);
