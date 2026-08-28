-- Migracja 192 — graf argumentów, katalog błędów logicznych i ich oznaczenia.
--
-- Graf jest trwały, nie liczony przy każdym odczycie. Powód jest jeden i
-- rozstrzygający: węzeł da się oznaczyć jako kluczowy (`roundtable.argument.pin`),
-- a oznaczenie postawione przez Operatora na węźle wyliczanym w locie znikałoby
-- przy następnym odczycie, bo węzeł dostawałby nowy identyfikator.
--
-- Katalog błędów jest dwuwarstwowy. Definicje wnosi migracja — są wspólne dla
-- całej platformy i nie należą do okna. Zakres wykrywania należy do okna
-- (`roundtable.fallacy.catalog.set`), więc leży w osobnej tabeli wiążącej kod
-- błędu z oknem. Bez tego rozdziału wyłączenie błędu w jednym oknie wyłączałoby
-- go wszystkim.

-- ── Węzeł grafu argumentów ───────────────────────────────────────────────────
CREATE TABLE debata_wezel (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    -- Wypowiedź i uczestnik wskazywani kodem: węzeł scalony (2.6.4) pochodzi
    -- z wielu wypowiedzi naraz i wtedy nie wskazuje żadnej.
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

-- ── Krawędź grafu ────────────────────────────────────────────────────────────
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

-- ── Katalog błędów logicznych — definicje wspólne ────────────────────────────
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

-- ── Zakres wykrywania w oknie ────────────────────────────────────────────────
-- Brak wiersza dla okna znaczy „katalog w całości włączony" — okno, w którym
-- nikt zakresu nie zawężał, wykrywa wszystko.
CREATE TABLE debata_katalog_okna (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    okno     TEXT    NOT NULL,
    kod      TEXT    NOT NULL,
    UNIQUE (okno, kod)
);

-- ── Oznaczenie błędu na węźle ────────────────────────────────────────────────
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
