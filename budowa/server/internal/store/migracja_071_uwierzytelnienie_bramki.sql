-- Migracja 071 — bramka uwierzytelnienia Operatora: metody wejścia i sesje bramki.
--
-- Bramka jest bytem osobnym od zastanych tabel: `konto` niesie poświadczenie do
-- kanału modelu, a `sesja` — kartę pracy, która istnieje niezależnie od tego,
-- czy ktokolwiek się zalogował. Operator jest jeden i bezimienny, konta
-- użytkownika ani adresu e-mail w zakresie nie ma.
--
-- Sekretu w bazie nie ma. Kolumna `sekret_odwolanie` niesie odwołanie do wpisu
-- w sejfie poświadczeń — tym samym, którym jadą konta i punkty dostępu
-- (`dane.SejfPlikowy`, plik `poswiadczenia.sejf` katalogu danych). W sejfie leży
-- skrót hasła, nie hasło; postać zapisu ustala `core/adapter_modul_auth_sekret.go`,
-- a schemat traktuje odwołanie jak nieprzezroczysty napis. Kolumny na hasło, PIN
-- ani materiał klucza tu nie ma, więc odczyt katalogu metod nie ma czego wynieść.
--
-- Urządzenie jest napisem, nie kluczem obcym. Kontrakt niesie `deviceId` jako
-- identyfikator nadawany przez klienta i nie wymaga, żeby urządzenie stało
-- wcześniej w katalogu `urzadzenie`; katalog ten wypełnia rozpoznanie maszyny
-- bieżącej i punkty dostępu rodzaju `localDirectory`, czyli inny zbiór maszyn
-- niż ten, z którego Operator wchodzi. Klucz obcy odmawiałby założenia PIN-u
-- z powodu, którego kontrakt nie zna.
--
-- Kotwica jest najwyżej jedna i pilnuje tego baza, nie warstwa wyżej. Kotwicą
-- jest hasło bramki ustawiane raz przez `auth.register`; ono rozstrzyga
-- o trwałej odmowie powtórzenia tej komendy i jest jedyną metodą, której
-- `auth.method.remove` zdjąć nie może.
--
-- PIN jest właściwy urządzeniu, więc para (urządzenie, rodzaj) jest
-- jednoznaczna: drugi PIN na tej samej maszynie byłby drugą prawdą o tym samym
-- wejściu. Ponowne założenie PIN-u na urządzeniu ma być zmianą PIN-u, a nie
-- cichym dołożeniem drugiego — indeks poniżej wymusza, żeby rdzeń rozstrzygnął
-- to jawnie.
--
-- Sesja bramki trzyma skrót tokenu, nie token: token jest poświadczeniem na
-- okaziciela i zapisany wprost byłby sekretem w bazie. Powstaje ze źródła
-- losowości kryptograficznej, więc do rozpoznania wystarczy skrót SHA-256 bez
-- soli i bez rozciągania — materiału do zgadywania nie ma. Kolumna jest UNIQUE,
-- bo skrót jest tożsamością sesji.
--
-- Unieważnienie jest datą, nie skasowaniem wiersza: `auth.password.reset` oddaje
-- liczbę sesji unieważnionych zmianą hasła i musi ją mieć skąd wziąć.

CREATE TABLE metoda_uwierzytelnienia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('password','pin','hello')),
    etykieta                 TEXT,
    urzadzenie_kod           TEXT,
    nazwa_urzadzenia         TEXT,
    kotwica                  INTEGER NOT NULL DEFAULT 0 CHECK(kotwica IN (0,1)),
    sekret_odwolanie         TEXT    NOT NULL DEFAULT '',
    utworzono                INTEGER NOT NULL,
    ostatnio_uzyto           INTEGER,
    -- Kotwicą jest wyłącznie hasło i nie należy ono do żadnego urządzenia:
    -- hasło otwiera bramkę z każdej maszyny, PIN i klucz Hello — tylko ze swojej.
    CHECK(kotwica = 0 OR rodzaj = 'password'),
    CHECK(rodzaj <> 'password' OR urzadzenie_kod IS NULL)
);

CREATE UNIQUE INDEX idx_metoda_uwierzytelnienia_kotwica
    ON metoda_uwierzytelnienia(kotwica) WHERE kotwica = 1;

CREATE UNIQUE INDEX idx_metoda_uwierzytelnienia_urzadzenie
    ON metoda_uwierzytelnienia(urzadzenie_kod, rodzaj) WHERE urzadzenie_kod IS NOT NULL;

CREATE INDEX idx_metoda_uwierzytelnienia_wykaz
    ON metoda_uwierzytelnienia(rodzaj, utworzono, id);

CREATE TABLE sesja_bramki (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    token_skrot    TEXT    NOT NULL UNIQUE,
    metoda_rodzaj  TEXT    CHECK(metoda_rodzaj IS NULL
                                 OR metoda_rodzaj IN ('password','pin','hello')),
    urzadzenie_kod TEXT,
    wygasa         INTEGER NOT NULL,
    utworzono      INTEGER NOT NULL,
    uniewazniono   INTEGER
);

CREATE INDEX idx_sesja_bramki_czynne ON sesja_bramki(uniewazniono, wygasa);
