-- Zamknięcie tury jako zdarzenie wykonawcze.
--
-- Stan odpowiedzi (`wiadomosc.stan`) rozstrzyga zdarzenie `result` programu CLI,
-- nie przebieg strumienia: zdarzenie potrafi przyjść z `subtype:"success"`
-- i jednocześnie `is_error:true`, więc brak błędu kanału nie znaczy „tura
-- zakończona poprawnie".
--
-- Wiersz powstaje po turze i niczym nie steruje. Jest zapisem pochodzenia
-- rozstrzygnięcia: jakie zdarzenie zamknęło turę, co niosło i jaki stan
-- wiadomości z niego wyprowadzono.
--
-- `okno_kod` i `wiadomosc_kod` są identyfikatorami kontraktowymi (napisy),
-- nie kluczami obcymi: zapis zamknięcia ma przeżyć okno, którego dotyczył,
-- a usunięcie sesji nie wymazuje śladu rozstrzygnięć. Tekstu odpowiedzi tu nie
-- ma — mieszka w `wiadomosc`, tu mieszka wyłącznie zdarzenie.

CREATE TABLE zamkniecie_tury (
    id                   INTEGER PRIMARY KEY,
    chwila               INTEGER NOT NULL,              -- epoka w milisekundach
    okno_kod             TEXT    NOT NULL,
    wiadomosc_kod        TEXT    NOT NULL UNIQUE,       -- jedna tura = jedno zamknięcie
    podtyp               TEXT    NOT NULL DEFAULT '',   -- dosłowne `subtype` zdarzenia result
    blad                 INTEGER NOT NULL DEFAULT 0
                         CHECK (blad IN (0, 1)),        -- dosłowne `is_error`
    -- stan wiadomości wyprowadzony ze zdarzenia — słownik ten sam co
    -- `wiadomosc.stan`; zapisany, żeby rozstrzygnięcie było odtwarzalne bez
    -- odgadywania reguły z kodu.
    stan_nadany          TEXT    NOT NULL
                         CHECK (stan_nadany IN ('zakonczona','zatrzymana','bledna')),
    koszt_usd            REAL    NOT NULL DEFAULT 0,
    tury                 INTEGER NOT NULL DEFAULT 0,
    czas_ms              INTEGER NOT NULL DEFAULT 0,
    konto                TEXT    NOT NULL DEFAULT '',   -- kod konta, nigdy sekret
    sesja_cli            TEXT    NOT NULL DEFAULT '',   -- identyfikator rozmowy CLI
    -- wykaz typów linii przechwyconych w turze (CSV w kolejności pierwszego
    -- wystąpienia) i liczba linii nieczytelnych: widać, co kanał zobaczył,
    -- także to, czego nie umiał przełożyć.
    typy_zdarzen         TEXT    NOT NULL DEFAULT '',
    linie_nierozpoznane  INTEGER NOT NULL DEFAULT 0
);

-- Odczyt idzie po oknie (ślad rozmowy) i po czasie (ślad dnia pracy).
CREATE INDEX idx_zamkniecie_okno   ON zamkniecie_tury (okno_kod, chwila);
CREATE INDEX idx_zamkniecie_chwila ON zamkniecie_tury (chwila);
