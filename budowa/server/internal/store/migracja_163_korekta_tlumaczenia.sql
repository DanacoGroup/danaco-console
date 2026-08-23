-- Migracja 163 — ustalenia korekty językowej modułu Translate
-- (`translate.proofread.run`, `translate.proofread.apply`).
--
-- Ustalenie korekty musi być trwałe, inaczej `proofread.apply` nie ma czego
-- zastosować: kontrakt każe wskazać ustalenia identyfikatorami
-- (`findingIds`) w osobnym wywołaniu, a identyfikator wydany w odpowiedzi
-- i zapomniany po niej byłby identyfikatorem donikąd. To odróżnia korektę od
-- kontroli jakości (`quality.check`), której zastrzeżenia są migawką
-- wymienianą w całości i nikt ich nie adresuje pojedynczo.
--
-- Ustalenie ma trzy stany życia i wszystkie trzy są tu widoczne: nowe
-- (`zastosowano` i `odrzucono` puste), zastosowane (Operator przyjął
-- poprawkę) i odrzucone (Operator ustalenie oddalił). Kasowanie odrzuconych
-- byłoby zgubieniem odpowiedzi „nie, tak ma być" — kolejny przebieg korekty
-- zgłosiłby to samo po raz drugi.
--
-- Wynik czytelności (`ReadabilityScore`) nie ma tu tabeli. Jest funkcją treści
-- panelu w chwili pomiaru: liczba zdań, długość słowa, długość zdania. Wiersz
-- trzeba by unieważniać przy każdej korekcie panelu, a `proofread.run` i tak
-- liczy go od nowa; nikt też nie adresuje wyniku czytelności identyfikatorem.

CREATE TABLE ustalenie_korekty (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    panel_id                 INTEGER NOT NULL REFERENCES panel_tlumaczenia(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('grammar','spelling','punctuation','style',
                                                      'register','readability','falseFriends','typography')),
    waga                     TEXT    NOT NULL CHECK(waga IN ('hint','warning','error')),
    segment                  TEXT,
    szczegol                 TEXT    NOT NULL,
    -- Propozycja poprawki. Pusta znaczy „widzę usterkę, ale nie mam czym jej
    -- zastąpić" — `proofread.apply` odmawia wtedy zamiast wstawić pustkę.
    propozycja               TEXT,
    zastosowano              INTEGER,
    odrzucono                INTEGER,
    utworzono                INTEGER NOT NULL
);
CREATE INDEX idx_ustalenie_korekty_panel ON ustalenie_korekty(panel_id, utworzono DESC);
