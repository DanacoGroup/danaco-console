-- Migracja 123 — wstrzymanie kroku, nie całego zlecenia.
--
-- Krokiem w tym produkcie jest wiersz `pozycja_kolejki`. Silnik
-- `core/kolejka_silnik.go` nazywa go krokiem wprost: `krokNaprzod`,
-- `wykonawcaKroku`, metoda `krok`. Drugiego bytu „krok zlecenia" nie ma i ta
-- migracja go nie zakłada. `pozycja_kolejki.stan` zna siedem wartości (oczekuje,
-- przydzielona, wykonywana, do_weryfikacji, ukonczona, bledna, anulowana), ale
-- nie zna stanu „wstrzymany" i nie ma gdzie zapisać decyzji Operatora.
-- `queue.action pause` wstrzymuje kolejkę, a pozycji nie tyka. Ta tabela zapisuje
-- wstrzymanie pojedynczego kroku wraz z decyzją Operatora, żeby wznowienie
-- wracało do stanu sprzed wstrzymania, a nie traktowało pracy przerwanej jak
-- skończonej.
--
-- Sam schemat pilnuje trzech rzeczy:
--
--   1. Wstrzymać da się wyłącznie krok żywy. `stan_pozycji_przed` ma CHECK
--      dopuszczający cztery stany robocze i ani jednego końcowego. Wstrzymanie
--      kroku ukończonego, błędnego albo anulowanego odbija się o schemat,
--      a nie o pamiętliwość warstwy wyżej.
--   2. Jedno czynne wstrzymanie na krok (indeks częściowy niżej). Drugie
--      wstrzymanie tego samego kroku przed zastosowaniem pierwszej decyzji nie
--      powstanie — inaczej dwie decyzje Operatora ścigałyby się o ten sam krok.
--   3. Decyzja nie ginie. Kolumny `zdecydowano_o`, `zastosowano_o`,
--      `doreczono_o` rozdzielają trzy różne fakty: Operator zdecydował, decyzja
--      zmieniła los kroku, decyzja dojechała do wykonawcy. Wiersz z decyzją
--      i pustym `doreczono_o` mówi wprost „decyzja jest, wykonawcy nie było".
--
-- Stan sterowania a stan pracy. Tabela zapisuje trzy wartości sterowania:
-- wstrzymany, zatwierdzony, odrzucony. Nie zapisuje `czeka` ani `biegnie` — te
-- dwa czyta się z `pozycja_kolejki.stan` (oczekuje/przydzielona → czeka,
-- wykonywana → biegnie). Zapisanie ich tutaj byłoby drugą prawdą o tym samym
-- kroku.
--
-- Ślad zostaje. Wiersz wstrzymania nie kasuje się po zastosowaniu decyzji —
-- historia wstrzymań kroku jest częścią przejrzystości pętli, dokładnie jak
-- `log_akcji_kolejki`, do którego te zdarzenia też trafiają.

CREATE TABLE wstrzymanie_kroku (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    pozycja_kolejki_id     INTEGER NOT NULL
                                   REFERENCES pozycja_kolejki(id) ON DELETE CASCADE,

    -- Stan sterowania krokiem. 'wstrzymany' znaczy: krok stoi i czeka na decyzję
    -- Operatora — to jest zdarzenie, które AOD ma wypchnąć na telefon.
    stan                   TEXT    NOT NULL DEFAULT 'wstrzymany'
                                   CHECK(stan IN ('wstrzymany','zatwierdzony','odrzucony')),

    -- Stan pracy kroku z chwili wstrzymania. Po decyzji krok wraca do tego stanu,
    -- a nie na początek — proces wznawia się, nie zaczyna od nowa. Stany końcowe
    -- są poza CHECK (punkt 1 wyżej).
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

    -- Decyzja i znacznik decyzji chodzą parą w obie strony: krok wstrzymany nie
    -- ma daty decyzji, a krok zdecydowany ją ma. Bez tego więzu dałoby się
    -- zapisać „zatwierdzony", którego nikt nie zatwierdził.
    CHECK((stan = 'wstrzymany') = (zdecydowano_o IS NULL)),
    -- Porządek faktów: zastosować da się wyłącznie decyzję, która zapadła,
    -- a doręczyć wyłącznie decyzję, która została zastosowana.
    CHECK(zastosowano_o IS NULL OR zdecydowano_o IS NOT NULL),
    CHECK(doreczono_o   IS NULL OR zastosowano_o IS NOT NULL)
);

CREATE INDEX idx_wstrzymanie_kroku_pozycja
    ON wstrzymanie_kroku(pozycja_kolejki_id, id);

-- Jedno wstrzymanie czynne na krok. „Czynne" znaczy: decyzja jeszcze nie
-- została zastosowana — obejmuje więc i krok czekający na decyzję, i krok
-- z decyzją, której nie zdążono zastosować (np. rdzeń stanął w międzyczasie).
-- Ten drugi przypadek jest właśnie tym, czego nie wolno zgubić.
CREATE UNIQUE INDEX idx_wstrzymanie_kroku_czynne
    ON wstrzymanie_kroku(pozycja_kolejki_id) WHERE zastosowano_o IS NULL;
