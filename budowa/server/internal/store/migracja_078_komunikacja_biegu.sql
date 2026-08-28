-- Migracja 078 — komunikacja między modelami: co następny widzi z pracy poprzednika.
--
-- Model kończy turę, następny dostaje jego wynik jako wejście. Zakres tego, co
-- widzi następny, wybiera się osobno dla każdego wiązania; domyślnie `artefakt`:
--
--   artefakt      wyłącznie wytwór poprzednika — plik, poprawka, odpowiedź.
--                 Najmniejszy i najtańszy w żetonach, a walidator ocenia wynik,
--                 nie tok myślenia.
--   streszczenie  skrót wypowiedzi. Wymaga osobnej tury modelu, więc ma własne
--                 miejsce, w którym może zawieść.
--   wypowiedz     cała wypowiedź poprzednika. Najbogatszy i najdroższy; przy
--                 czterech stanowiskach i długiej pętli zalewa okno kontekstu.
--
-- Jeden sztywny zakres byłby albo zbyt ubogi dla koordynatora przekazującego
-- zlecenie, albo zbyt drogi dla walidatora oceniającego wynik.
--
-- Ślad narzędzi poprzednika jest polem osobnym (`ze_sladem_narzedzi`),
-- domyślnie wyłączonym: bywa większy od samej wypowiedzi i najczęściej jest
-- szumem, ale walidator sprawdzający, czy wykonawca uruchomił testy, bez niego
-- nie ma czego sprawdzić. Pole jest prostopadłe do `zakres` — wolno chcieć
-- samego artefaktu ze śladem i całej wypowiedzi bez śladu.
--
-- Pola wyłączającego podgląd tej rozmowy nie ma: ruch idzie tą samą drogą co
-- każda tura — fragmentami strumienia do okien — a wiersz w `wiadomosc_biegu`
-- zostaje jako ślad.
--
-- Ta migracja nie buduje drugiego mechanizmu widoczności kontekstu. Rodzina
-- `isolation.*` rozstrzyga warstwowo, co okno widzi z historii, pamięci
-- i kontekstu sąsiada — jedenaście punktów przez osiem poziomów zasięgu.
-- Podział jest ostry:
--
--   ta migracja        co poprzednik jawnie przekazuje następnemu
--   profil izolacji    co następny może zobaczyć z kontekstu poprzednika
--                      poza tym, co mu przekazano
--
-- Stanowisko obsady wskazuje swój profil izolacji kolumną `profil_izolacji_id`.

-- ── 1. Przekazanie — wiązanie, czyli plan rozmowy ───────────────────────────
--
-- To jest projekt, nie ruch: „gdy skończy stanowisko 2, jego artefakt ze śladem
-- narzędzi idzie do stanowiska 3". Operator układa te wiązania, projektując bieg
-- a silnik je czyta, gdy stanowisko kończy pracę.
--
-- Miejsce, nie identyfikator stanowiska. Wiązanie wskazuje `od_miejsca`
-- i `do_miejsca`, czyli numery 1-4 z obsady — tak samo jak `zaleznosc_kroku_automatyki`
-- wiąże kroki po identyfikatorze zewnętrznym, a nie po kluczu wiersza. Powód
-- ten sam: obsadę wolno podmienić w całości, a plan rozmowy ma to przeżyć.

CREATE TABLE przekazanie_biegu (
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    bieg_id            INTEGER NOT NULL REFERENCES bieg_orkiestracji(id) ON DELETE CASCADE,
    od_miejsca         INTEGER NOT NULL CHECK(od_miejsca BETWEEN 1 AND 4),
    do_miejsca         INTEGER NOT NULL CHECK(do_miejsca BETWEEN 1 AND 4),
    zakres             TEXT    NOT NULL DEFAULT 'artefakt'
                               CHECK(zakres IN ('artefakt','streszczenie','wypowiedz')),
    ze_sladem_narzedzi INTEGER NOT NULL DEFAULT 0 CHECK(ze_sladem_narzedzi IN (0,1)),
    -- Warunek przekazania; puste znaczy „zawsze". Ten sam kształt co `warunek`
    -- kroku automatyki — nie wprowadzam drugiego języka warunków.
    warunek            TEXT,
    kolejnosc          INTEGER NOT NULL DEFAULT 0,
    utworzono          TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Stanowisko nie przekazuje samo sobie — to byłaby pętla bez wyjścia.
    CHECK (od_miejsca <> do_miejsca),
    UNIQUE (bieg_id, od_miejsca, do_miejsca)
);

CREATE INDEX idx_przekazanie_biegu_od
    ON przekazanie_biegu(bieg_id, od_miejsca, kolejnosc);

-- ── 2. Wiadomość biegu — ruch, czyli co naprawdę poszło ─────────────────────
--
-- Ślad tego, co jedno stanowisko przekazało drugiemu. Wiersz powstaje w chwili
-- przekazania i zostaje — to on odpowiada Operatorowi na pytanie „dlaczego
-- walidator tak orzekł", gdy bieg dawno się skończył.
--
-- Dwa pola treści, nie jedno. `tresc` niesie to, co objął `zakres`; `slad_narzedzi`
-- niesie ślad, jeśli wiązanie go żądało. Sklejenie ich w jedno uniemożliwiłoby
-- późniejsze pokazanie samej treści bez szumu.
--
-- Pozycja kolejki jako źródło. Wiadomość wie, która pozycja ją wytworzyła —
-- dzięki temu panel zadań w tle umie pokazać, że wynik podagenta wszedł do
-- rozmowy, a nie zginął.

CREATE TABLE wiadomosc_biegu (
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    bieg_id            INTEGER NOT NULL REFERENCES bieg_orkiestracji(id) ON DELETE CASCADE,
    od_miejsca         INTEGER NOT NULL CHECK(od_miejsca BETWEEN 1 AND 4),
    do_miejsca         INTEGER NOT NULL CHECK(do_miejsca BETWEEN 1 AND 4),
    zakres             TEXT    NOT NULL
                               CHECK(zakres IN ('artefakt','streszczenie','wypowiedz')),
    tresc              TEXT    NOT NULL,
    slad_narzedzi      TEXT,
    -- Pozycja kolejki, której wykonanie wytworzyło tę treść. Puste dla
    -- przekazania spoza kolejki (np. zlecenie początkowe od koordynatora).
    pozycja_kolejki_id INTEGER REFERENCES pozycja_kolejki(id) ON DELETE SET NULL,
    -- Podagent, który treść wytworzył; puste, gdy wytworzyło ją samo stanowisko.
    podagent_id        INTEGER REFERENCES podagent(id) ON DELETE SET NULL,
    utworzono          TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    CHECK (od_miejsca <> do_miejsca)
);

-- Odczyt rozmowy biegu w kolejności — to jest widok, który Operator ogląda.
CREATE INDEX idx_wiadomosc_biegu_bieg
    ON wiadomosc_biegu(bieg_id, id);

-- „Co dostało stanowisko 3" — wejście stanowiska przy rozpoczęciu tury.
CREATE INDEX idx_wiadomosc_biegu_do
    ON wiadomosc_biegu(bieg_id, do_miejsca, id);
