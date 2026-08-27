-- Tworzy tabelę wstrzymanie_kroku, zapisującą wstrzymanie pojedynczego kroku kolejki wraz z decyzją Operatora, aby wznowienie wracało do stanu sprzed wstrzymania.

CREATE TABLE wstrzymanie_kroku (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    pozycja_kolejki_id     INTEGER NOT NULL
                                   REFERENCES pozycja_kolejki(id) ON DELETE CASCADE,

    -- Stan sterowania krokiem: wstrzymany oznacza, że krok czeka na decyzję Operatora.
    stan                   TEXT    NOT NULL DEFAULT 'wstrzymany'
                                   CHECK(stan IN ('wstrzymany','zatwierdzony','odrzucony')),

    -- Stan pracy kroku z chwili wstrzymania; po decyzji krok wraca do niego, nie zaczyna od nowa.
    stan_pozycji_przed     TEXT    NOT NULL
                                   CHECK(stan_pozycji_przed IN ('oczekuje','przydzielona',
                                                                'wykonywana','do_weryfikacji')),

    powod                  TEXT,   -- dlaczego krok wstrzymano (wolny opis)
    uzasadnienie           TEXT,   -- treść decyzji Operatora — jedzie do wykonawcy

    wstrzymano_o           TEXT    NOT NULL
                                   DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zdecydowano_o          TEXT,   -- Operator zdecydował
    zastosowano_o          TEXT,   -- decyzja zmieniła los kroku
    doreczono_o            TEXT,   -- decyzja dojechała do wykonawcy

    -- Decyzja i jej znacznik czasu chodzą parą, aby nie dało się zapisać decyzji bez daty.
    CHECK((stan = 'wstrzymany') = (zdecydowano_o IS NULL)),
    -- Porządek faktów: zastosować można tylko zapadłą decyzję, doręczyć tylko zastosowaną.
    CHECK(zastosowano_o IS NULL OR zdecydowano_o IS NOT NULL),
    CHECK(doreczono_o   IS NULL OR zastosowano_o IS NOT NULL)
);

CREATE INDEX idx_wstrzymanie_kroku_pozycja
    ON wstrzymanie_kroku(pozycja_kolejki_id, id);

-- Tworzy unikalny indeks czynnego wstrzymania kroku, obejmujący krok czekający na decyzję oraz krok z decyzją jeszcze niezastosowaną.
CREATE UNIQUE INDEX idx_wstrzymanie_kroku_czynne
    ON wstrzymanie_kroku(pozycja_kolejki_id) WHERE zastosowano_o IS NULL;
