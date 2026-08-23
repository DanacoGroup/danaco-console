-- Migracja 248 — wykaz kluczy SSH znanych rdzeniowi (okno Session Manager).
--
-- Wykaz, nie sejf. W wierszu stoi ŚCIEŻKA klucza prywatnego na maszynie rdzenia,
-- treść klucza PUBLICZNEGO i odcisk — nic z materiału tajnego. Klucz prywatny
-- nie opuszcza dysku maszyny rdzenia ani przy wytworzeniu (`terminal.key.generate`
-- pisze plik i oddaje sam odcisk), ani przy wciągnięciu do wykazu
-- (`terminal.key.import` wskazuje się ścieżką, nie treścią).
--
-- `haslo` jest jedną wartością logiczną „klucz jest chroniony hasłem”, a nie
-- hasłem. Odwołanie do hasła w sejfie idzie osobno, poza tą tabelą; wpisanie go
-- tutaj zamieniłoby wykaz w skrytkę, którą baza nie jest.
--
-- `sciezka` ma warunek UNIQUE: dwa wpisy wskazujące ten sam plik byłyby nie do
-- rozróżnienia w chwili, gdy `terminal.key.remove` ma usunąć pliki z dysku.

CREATE TABLE terminal_klucz (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    kod         TEXT    NOT NULL UNIQUE,
    nazwa       TEXT    NOT NULL,
    rodzaj      TEXT    NOT NULL CHECK(rodzaj IN ('ed25519','rsa','ecdsa')),
    -- Odcisk klucza publicznego w postaci SHA256:… — ta sama, którą pokazuje
    -- OpenSSH, żeby Operator mógł porównać wprost.
    odcisk      TEXT    NOT NULL DEFAULT '',
    -- Treść klucza publicznego do przeniesienia na host docelowy.
    klucz_jawny TEXT    NOT NULL DEFAULT '',
    -- Ścieżka klucza prywatnego na maszynie rdzenia.
    sciezka     TEXT    NOT NULL UNIQUE,
    -- Czy klucz jest chroniony hasłem; samo hasło nie ma tu kolumny.
    haslo       INTEGER NOT NULL DEFAULT 0 CHECK(haslo IN (0,1)),
    utworzono   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
