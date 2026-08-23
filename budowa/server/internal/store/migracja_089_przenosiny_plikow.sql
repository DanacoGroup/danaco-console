-- Migracja 089 — rejestr przenosin plików torem zdalnym.
--
-- Przy zasięgu `remote` proces modelu
-- pracuje na hoście zdalnym, a pliki bywają gdzie indziej: nagranie dyktowania
-- leży na maszynie Operatora, katalog roboczy okna — na maszynie rdzenia.
-- Rozstrzygnięcie jest dwudzielne:
--
-- 1. Pliki zwykłe — przenosiny są, pod tą samą zgodą per host co tor procesu.
--    Nośnikiem jest `scp` po tym samym SSH: te same klucze, ta
--    sama zgoda, żaden nowy protokół. Przeniesienie wykonuje
--    `zdalne.PrzeniesPlik`, a każdy wykonany ruch bajtów zostawia wiersz w tej
--    tabeli — prowenancję przenosin: skąd, dokąd, ile i kiedy.
--    Wiersz nie jest kolejką ani zleceniem — zapis następuje po przeniesieniu
--    i niczego nie steruje (zero bramek).
--
-- 2. Nagrania dźwięku — odmowa, z dwóch niezależnych powodów:
--    · reguła zapisana w kliencie wprost: dźwięk nie opuszcza maszyny Operatora
--      (client/src/okno-komunikacji/dyktowanie/dostarczenie-nagrania.ts) —
--      wysłanie nagrania na host zdalny byłoby dokładnie tym, czego ta reguła
--      zakazuje;
--    · potrzeby nie ma: silnik mowy (`internal/mowa`) jest usługą rdzenia
--      i bierze ścieżkę na maszynie silnika (`AudioRef` w mowa/nagranie.go);
--      transkrypcja domyka się przed torem zdalnym, a do procesu na hoście
--      zdalnym jedzie wyłącznie tekst.
--    Odmowę egzekwuje `zdalne.PrzeniesPlik` po wykazie rozszerzeń
--    `mowa.FormatyNagran` (jeden wykaz) — zdaniem trójczęściowym.
--
-- `okno_id` jest identyfikatorem kontraktowym okna (napis `okn_…`), nie kluczem
-- obcym: przenosiny bywają wykonywane poza cyklem życia okna (przygotowanie
-- katalogu przed startem, sprzątanie po zamknięciu), a prowenancja ma przeżyć
-- okno, którego dotyczyła. Pusty napis znaczy: ruch poza kontekstem okna.

CREATE TABLE zdalne_przeniesienie (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    host_id          INTEGER NOT NULL REFERENCES host_zdalny(id) ON DELETE CASCADE,
    okno_id          TEXT    NOT NULL DEFAULT '',
    -- Kierunek ruchu bajtów względem maszyny rdzenia.
    kierunek         TEXT    NOT NULL CHECK(kierunek IN ('wyslanie','pobranie')),
    sciezka_zrodlowa TEXT    NOT NULL,
    sciezka_docelowa TEXT    NOT NULL,
    -- Rozmiar w bajtach zmierzony po stronie lokalnej; 0, gdy pomiar się nie
    -- powiódł — brak pomiaru nie unieważnia zapisu ruchu.
    rozmiar          INTEGER NOT NULL DEFAULT 0,
    przeniesiono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zdalne_przeniesienie_host ON zdalne_przeniesienie(host_id, przeniesiono);
