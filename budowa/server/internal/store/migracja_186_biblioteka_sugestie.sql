-- Migracja 186 — moduł Library: sugestie porządkujące.
--
-- Domyślnym zachowaniem modułu jest sugestia z akceptacją Operatora, nie zapis
-- bez pytania: klasyfikacja wsadowa (`library.classify.run`) wytwarza wiersze
-- tej tabeli, a dopiero `library.suggestion.apply` zamienia je w zmianę zasobu.
-- Sugestia musi więc przeżyć między jednym żądaniem a drugim — stąd tabela,
-- a nie wynik oddany i zapomniany.
--
-- Uzasadnienie jest kolumną obowiązkową. Sugestia bez powodu jest poleceniem
-- podanym bez podstawy, a Operator ma decydować, nie zgadywać, skąd wzięła się
-- propozycja.
--
-- Decyzja zostaje przy wierszu (`stan`), zamiast kasować go przy odrzuceniu:
-- sugestia odrzucona ma nie wracać przy kolejnym przebiegu klasyfikacji, więc
-- rdzeń musi wiedzieć, że raz już padła i została odsunięta.

CREATE TABLE sugestia_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    plik_kod                 TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL
                             CHECK(rodzaj IN ('etykieta','kolekcja','duplikat',
                                              'osierocony','wrazliwy')),
    wartosc                  TEXT,
    uzasadnienie             TEXT    NOT NULL,
    -- Pewność w setnych, jeśli model ją podał; brak znaczy „model nie mówi".
    pewnosc                  INTEGER,
    stan                     TEXT    NOT NULL DEFAULT 'oczekujaca'
                             CHECK(stan IN ('oczekujaca','przyjeta','odrzucona')),
    rozstrzygnieto           TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_sugestia_biblioteki_stan ON sugestia_biblioteki(stan, utworzono DESC);
CREATE INDEX idx_sugestia_biblioteki_plik ON sugestia_biblioteki(plik_kod, stan);
