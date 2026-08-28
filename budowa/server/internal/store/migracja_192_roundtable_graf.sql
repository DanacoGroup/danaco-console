-- Migracja 192 wprowadza trwały graf argumentów oraz dwuwarstwowy katalog błędów
-- logicznych, którego definicje są wspólne dla platformy, a zakres wykrywania
-- jest przypisany do okna.

-- Tabela debata_wezel przechowuje węzeł grafu argumentów utworzony w oknie wraz
-- z aktem mowy, treścią wypowiedzi i licznikiem poparcia po scaleniu równoważnych
-- wypowiedzi.
CREATE TABLE debata_wezel (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    -- Węzeł scalony łączy wiele wypowiedzi, więc pola wypowiedź i uczestnik pozostają puste.
    wypowiedz                TEXT    NOT NULL DEFAULT '',
    uczestnik                TEXT    NOT NULL DEFAULT '',
    tura                     TEXT    NOT NULL DEFAULT '',
    akt_mowy                 TEXT    NOT NULL DEFAULT 'claim'
                                     CHECK(akt_mowy IN ('claim','argument','counterArgument',
                                                        'question','concession','rebuttal')),
    tresc                    TEXT    NOT NULL,
    -- Poparcie liczy uczestników, którzy wypowiedzieli ten sam węzeł po scaleniu.
    poparcie                 INTEGER NOT NULL DEFAULT 1,
    kluczowy                 INTEGER NOT NULL DEFAULT 0 CHECK(kluczowy IN (0,1)),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_wezel_okno ON debata_wezel(okno, id);
CREATE INDEX idx_debata_wezel_tura ON debata_wezel(tura, id);

-- Tabela debata_krawedz łączy węzły grafu relacją retoryczną — wsparciem, atakiem,
-- odpowiedzią, ustępstwem albo powtórzeniem — wraz z pewnością rozpoznania tej
-- relacji.
CREATE TABLE debata_krawedz (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    wezel_zrodlowy           TEXT    NOT NULL,
    wezel_wskazany           TEXT    NOT NULL,
    relacja                  TEXT    NOT NULL
                                     CHECK(relacja IN ('supports','attacks','answers',
                                                       'concedes','restates')),
    pewnosc                  REAL    NOT NULL DEFAULT -1.0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_krawedz_okno ON debata_krawedz(okno, id);

-- Tabela debata_katalog_bledu przechowuje definicje błędów logicznych wspólne dla
-- całej platformy, niezależnie od tego, które z nich wykrywa poszczególne okno.
CREATE TABLE debata_katalog_bledu (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    kod         TEXT    NOT NULL UNIQUE,
    nazwa       TEXT    NOT NULL,
    opis        TEXT    NOT NULL
);

INSERT INTO debata_katalog_bledu (kod, nazwa, opis) VALUES
    ('adHominem', 'Atak na osobę',
     'Zarzut kierowany wobec mówcy zamiast wobec jego twierdzenia.'),
    ('strawman', 'Kukła',
     'Podmiana stanowiska przeciwnika na wersję słabszą i zbicie tej wersji.'),
    ('falseDilemma', 'Fałszywa alternatywa',
     'Przedstawienie dwóch wariantów jako jedynych, choć istnieją inne.'),
    ('slipperySlope', 'Równia pochyła',
     'Wywodzenie skrajnego skutku z pierwszego kroku bez wykazania ogniw pośrednich.'),
    ('appealToAuthority', 'Odwołanie do autorytetu',
     'Uznanie twierdzenia za prawdziwe wyłącznie dlatego, że wypowiedział je autorytet.'),
    ('appealToEmotion', 'Odwołanie do emocji',
     'Zastąpienie przesłanki wzbudzeniem emocji u odbiorcy.'),
    ('circularReasoning', 'Błędne koło',
     'Wniosek użyty jako własna przesłanka.'),
    ('hastyGeneralization', 'Uogólnienie pochopne',
     'Wniosek ogólny wyprowadzony z liczby przypadków zbyt małej, by go poprzeć.'),
    ('falseCause', 'Fałszywa przyczyna',
     'Uznanie następstwa w czasie za związek przyczynowy.'),
    ('whataboutism', 'Odwrócenie zarzutu',
     'Odpowiedź na zarzut wskazaniem cudzego przewinienia zamiast odniesieniem się do zarzutu.'),
    ('cherryPicking', 'Dobór danych',
     'Przywołanie wyłącznie danych zgodnych z tezą przy pominięciu przeczących.'),
    ('equivocation', 'Nadużycie wieloznaczności',
     'Użycie tego samego słowa w dwóch znaczeniach w obrębie jednego rozumowania.');

-- Tabela debata_katalog_okna zawęża zakres wykrywania błędów w oknie; brak wiersza
-- dla okna oznacza, że katalog jest włączony w całości i okno wykrywa wszystkie
-- błędy.
CREATE TABLE debata_katalog_okna (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    okno     TEXT    NOT NULL,
    kod      TEXT    NOT NULL,
    UNIQUE (okno, kod)
);

-- Tabela debata_oznaczenie_bledu zapisuje oznaczenie błędu logicznego postawione
-- na węźle grafu wraz z uzasadnieniem i pewnością rozpoznania.
CREATE TABLE debata_oznaczenie_bledu (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    wezel                    TEXT    NOT NULL,
    kod                      TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    uzasadnienie             TEXT    NOT NULL DEFAULT '',
    pewnosc                  REAL    NOT NULL DEFAULT -1.0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_oznaczenie_wezel ON debata_oznaczenie_bledu(wezel, id);
