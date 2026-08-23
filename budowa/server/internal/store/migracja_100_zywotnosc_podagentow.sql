-- Migracja 100 — żywotność podagentów: co się z nimi dzieje po awarii rdzenia.
--
-- Migracja 077 postawiła `podagent` z pełnym cyklem stanów
-- ('pending','running','done','failed','stopped'), ale wiersz nie wie, które
-- uruchomienie rdzenia prowadzi jego pracę. Rdzeń ubity w trakcie pracy
-- podagenta zostawia wiersz w stanie 'running' na zawsze — proces modelu zginął
-- razem z rdzeniem, a panel zadań w tle po restarcie pokazuje pracę, której nikt
-- nie wykonuje. Stan 'running' kłamie, a kłamiący stan jest gorszy niż brak stanu.
--
-- Podagent należy do uruchomienia rdzenia, nie do pliku bazy. Trwała jest jego
-- tożsamość i jego wynik — wykaz zakończonych ma przeżywać restart i daje się
-- przeglądać (migracja 077). Nietrwała jest jego praca: żyje w procesie modelu,
-- więc ginie z rdzeniem. Kolumna `uruchomienie_rdzenia` zapisuje tę granicę
-- wprost: wiersz niezakończony ze znacznikiem innym niż znacznik rdzenia, który
-- właśnie wstał, jest podagentem osieroconym i rdzeń zamyka go przy starcie.
--
-- Znacznik, a nie sam czas: oznaka życia starsza niż próg mówi „może padł,
-- a może po prostu długo myśli" — to domysł. Znacznik uruchomienia odpowiada bez
-- domysłu: jeśli prowadzi go rdzeń, którego już nie ma, praca nie jest wykonywana
-- przez nikogo. Oznaka życia zostaje obok, ale jako zapis ostatniego dotknięcia
-- wiersza, nie jako kryterium osierocenia.
--
-- Żadna z tych kolumn niczego nie zabrania: powołanie bez znacznika przechodzi
-- bez zmian, a wiersze zastane dostają znacznik pusty i przy pierwszym starcie po
-- tej migracji zostaną zamknięte jako osierocone — bo dokładnie tym są.
--
-- ALTER TABLE, nie przebudowa (wzorem migracji 014 i 017): `podagent` jest celem
-- indeksu częściowego `idx_podagent_pozycja` i nośnikiem wyniku, a przebudowa
-- musiałaby te indeksy odtworzyć bez potrzeby. Żaden klucz obcy nie wskazuje na
-- `podagent`, a nic tu nie ubywa, więc przebudowa nie miałaby czego naprawić.
--
-- Odwracalność: krok jest dokładający — trzy kolumny i jeden indeks częściowy.
-- Cofnięcie to `DROP INDEX idx_podagent_osieroceni` i przebudowa `podagent`
-- bez trzech kolumn (SQLite zdejmuje kolumnę wyłącznie przez przepisanie
-- tabeli). Dane sprzed migracji nie tracą znaczenia — kolumny mają wartości
-- domyślne, których odczyt bez tej migracji nikomu nie jest potrzebny.

-- ── 1. Znacznik uruchomienia rdzenia prowadzącego pracę ─────────────────────
--
-- Puste znaczy „prowadzenia nie odnotowano". Dla wiersza niezakończonego jest
-- to stan przejściowy między zapisem powołania a zasileniem kolejki — i tak
-- samo jak znacznik obcy oznacza sierotę, bo pracy nie prowadzi wtedy nikt.
ALTER TABLE podagent ADD COLUMN uruchomienie_rdzenia TEXT NOT NULL DEFAULT '';

-- ── 2. Oznaka życia — ostatnie dotknięcie wiersza przez rdzeń prowadzący ────
--
-- Nie jest kryterium osierocenia (patrz nagłówek), jest dowodem w meldunku:
-- pozwala powiedzieć, kiedy urwała się praca zamkniętego sieroty.
ALTER TABLE podagent ADD COLUMN oznaka_zycia TEXT;

-- ── 3. Powód zakończenia — dlaczego wiersz stanął w stanie końcowym ─────────
--
-- Stan 'failed' odpowiada „nie udało się", ale nie odróżnia błędu modelu od
-- pracy przerwanej awarią rdzenia. Panel zadań w tle i Operator potrzebują tej
-- różnicy: pierwsze jest wynikiem pracy, drugie — jej brakiem.
--
-- Słownik jest polski, bo to nazwa nadana tutaj, a nie wzięta z kontraktu:
-- kontrakt zna wyłącznie `SubagentStatus`, powodu nie niesie.
ALTER TABLE podagent ADD COLUMN powod_zakonczenia TEXT
    CHECK(powod_zakonczenia IS NULL OR powod_zakonczenia IN
          ('ukonczony','blad','zatrzymany','osierocony'));

-- ── 4. Odnalezienie sierot przy starcie ─────────────────────────────────────
--
-- Zapytanie startowe pyta o wiersze niezakończone. Indeks częściowy obejmuje
-- wyłącznie je, więc nie rośnie z wykazem zakończonych — a ten rośnie zawsze
-- („ZAKOŃCZONE 48" i dalej). Znacznik idzie drugą kolumną, bo warunek na stan
-- jest w indeksie stały, a po znaczniku odsiewa się rdzeń bieżący.
CREATE INDEX idx_podagent_osieroceni
    ON podagent(stan, uruchomienie_rdzenia)
    WHERE stan IN ('pending','running');
