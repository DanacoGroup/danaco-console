-- Migracja 096 — pełna treść rozmowy: bloki wiadomości i szukanie w treści.
--
-- Bez tych bloków tok rozumowania, wywołania narzędzi z wynikami, prowenancja
-- wywołania i metadane konta żyłyby tylko w strumieniu i ginęłyby z restartem
-- rdzenia — okno po odświeżeniu pokazywałoby sam tekst, choć widoki transkryptu
-- klienta ('rozumowanie', 'pelny') istnieją (widok-zapisu.ts).
--
-- ── blok_wiadomosci ─────────────────────────────────────────────────────────
-- Jeden wiersz = jeden fragment strumienia rodzaju nietekstowego, zapisany
-- w trakcie tury przez rejestrator bloków (core/rejestrator_blokow.go).
--
-- Tekstu tu nie ma i nie wolno go tu pisać. Treść tekstowa odpowiedzi mieszka
-- w `wiadomosc.tresc` (domyka ją dziennik po turze) — powtórzenie jej w blokach
-- byłoby drugą prawdą o tej samej wypowiedzi. Słownik `rodzaj` jest więc
-- słownikiem ChunkKind pomniejszonym o 'text'.
--
-- Rodzaj niesie wartość kontraktu (angielską), nie przekład. Kolumna
-- `wiadomosc.rodzaj_tresci` ma polski słownik, bo kontrakt odwzorowuje jej siedem
-- wartości polem `baza`; dla 'provenance' i 'account' kontrakt żadnego przekładu
-- nie zapisuje, a przekład wolno trzymać wyłącznie w kontrakcie. Dopisanie
-- własnego słownika tutaj byłoby drugim źródłem odwzorowania — zamiast tego
-- kolumna trzyma wartość kontraktu dosłownie, a round-trip zapis → odczyt →
-- klient obywa się bez tłumaczenia.
--
-- `wiadomosc_kod` i `okno_kod` są identyfikatorami kontraktowymi (napisy),
-- nie kluczami obcymi: rejestrator strumienia zna wyłącznie identyfikatory
-- rdzenia, a sięganie po klucz wiersza w środku tury dokładałoby odczyt bazy do
-- każdego fragmentu. Skutkiem braku klucza obcego kaskada usunięcia sesji bloków
-- nie zabiera — sprząta je czyszczenie kosza (dane/sesje_kosz.go), ta sama droga,
-- którą znika sama sesja.
--
-- `ladunek` niesie surowe pole `data` fragmentu (dowód pierwotny, jak
-- `dziennik_zdarzen.ladunek`); NULL znaczy fragment bez ładunku. `tresc` niesie
-- pole `text` fragmentu — tak nadchodzi tok rozumowania.

CREATE TABLE blok_wiadomosci (
    id            INTEGER PRIMARY KEY,
    chwila        INTEGER NOT NULL,                -- epoka w milisekundach
    okno_kod      TEXT    NOT NULL,
    wiadomosc_kod TEXT    NOT NULL,
    kolejnosc     INTEGER NOT NULL,                -- porządek w obrębie wiadomości
    rodzaj        TEXT    NOT NULL
                  CHECK (rodzaj IN ('thinking','tool_use','tool_result',
                                    'image','audio','error',
                                    'provenance','account')),
    tresc         TEXT    NOT NULL DEFAULT '',
    ladunek       TEXT
);

-- Odczyt idzie całym oknem (odtworzenie historii) i pojedynczą wiadomością.
CREATE INDEX idx_blok_wiadomosci_okno      ON blok_wiadomosci (okno_kod, wiadomosc_kod, kolejnosc);
CREATE INDEX idx_blok_wiadomosci_wiadomosc ON blok_wiadomosci (wiadomosc_kod, kolejnosc);

-- ── wiadomosc_szukanie ──────────────────────────────────────────────────────
-- Indeks pełnotekstowy FTS5 nad `wiadomosc.tresc` na sterowniku
-- modernc.org/sqlite: CREATE VIRTUAL TABLE USING fts5, MATCH, snippet() i rank
-- działają; tabela zewnętrzna (content='wiadomosc') z triggerami spójności
-- przechodzi wstawienie, aktualizację treści (droga ZapiszWynik) i usunięcie
-- kaskadą.
--
-- Tabela jest zewnętrzna (content=), więc treści nie powiela — jedyną prawdą
-- o słowie pozostaje `wiadomosc.tresc`, a indeks trzyma wyłącznie słownik
-- trafień. Spójność utrzymują triggery poniżej; są częścią schematu, nie warstwy
-- dane, bo indeks ma nadążać także za zapisem, który przyjdzie inną drogą niż
-- repozytorium wiadomości.
--
-- Ograniczenie zapisane jawnie: tokenizator unicode61 z remove_diacritics 2
-- sprowadza ż/ź/ó/ą/ę/ć/ń/ś do liter podstawowych, ale „ł" nie jest znakiem
-- składanym i pozostaje osobną literą — zapytanie „lodz" nie trafi w „łódź",
-- trafi „łodz" i „łódź". Warstwa dane powtarza to ograniczenie przy repozytorium
-- szukania (dane/szukanie_rozmow.go).

CREATE VIRTUAL TABLE wiadomosc_szukanie USING fts5(
    tresc,
    content='wiadomosc',
    content_rowid='id',
    tokenize="unicode61 remove_diacritics 2"
);

-- Zasiew: wiadomości zapisane przed tą migracją też mają być odnajdywalne.
INSERT INTO wiadomosc_szukanie(rowid, tresc)
    SELECT id, COALESCE(tresc, '') FROM wiadomosc;

CREATE TRIGGER wiadomosc_szukanie_po_wstawieniu AFTER INSERT ON wiadomosc BEGIN
    INSERT INTO wiadomosc_szukanie(rowid, tresc)
        VALUES (new.id, COALESCE(new.tresc, ''));
END;

CREATE TRIGGER wiadomosc_szukanie_po_usunieciu AFTER DELETE ON wiadomosc BEGIN
    INSERT INTO wiadomosc_szukanie(wiadomosc_szukanie, rowid, tresc)
        VALUES ('delete', old.id, COALESCE(old.tresc, ''));
END;

CREATE TRIGGER wiadomosc_szukanie_po_zmianie AFTER UPDATE OF tresc ON wiadomosc BEGIN
    INSERT INTO wiadomosc_szukanie(wiadomosc_szukanie, rowid, tresc)
        VALUES ('delete', old.id, COALESCE(old.tresc, ''));
    INSERT INTO wiadomosc_szukanie(rowid, tresc)
        VALUES (new.id, COALESCE(new.tresc, ''));
END;
