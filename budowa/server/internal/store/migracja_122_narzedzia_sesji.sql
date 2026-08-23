-- Migracja 122 — doraźne dołożenie narzędzia na czas sesji.
--
-- Raz załadowane narzędzie jest obecne agentowi do końca pracy w tej sesji:
-- nie przechodzi na stałe do zestawu agenta ani nie jest na jednorazowe użycie,
-- lecz trwa do zakończenia sesji albo usunięcia rozmowy.
--
-- Dlaczego własna tabela, a nie wiersz agenta. Dwa oczywiste miejsca zapisu
-- odpadają. Definicja eksperta
-- odpada wprost: „definicja pozostaje nietknięta" — dopisanie dołożenia do
-- `agent_umiejetnosc` albo `agent_konektor` (migracja 037/038) byłoby zmianą
-- definicji, a ta obowiązuje WSZYSTKIE sesje eksperta, nie tę jedną. Pamięć
-- tury odpada tak samo: dołożenie nie znika po jednym użyciu, więc nie może
-- żyć w bycie, który kończy się razem z turą.
--
-- Zostaje stan sesji. Ta tabela JEST stanem sesji: klucz obcy na `sesja(id)`
-- z kasowaniem kaskadowym daje dokładnie ten czas życia, który Właściciel
-- wskazał — dołożenia przeżywają rozłączenie klienta, bo wiersz leży
-- w bazie, a nie w pamięci gniazda, i giną razem z rozmową, bo `session.delete`
-- kasuje wiersz `sesja`, a kaskada zabiera resztę. Czasu życia nie pilnuje więc
-- ani zadanie sprzątające, ani warstwa wyższa: pilnuje go klucz obcy.
--
-- Zapisujemy pozycję, nie odwołanie do wykazu. Wiersz niesie nazwę pełną,
-- nazwę skróconą i opis — czyli te trzy pola nazewnicze — a nie sam klucz obcy do
-- pozycji wykazu. Powód: wykaz po ukośniku nie jest tabelą. Składa się na
-- bieżąco z komend kontraktu, katalogu akcji i katalogu rozszerzeń
-- (`dane/narzedzia_sesji.go`), więc nie ma czego wskazać kluczem obcym. Gdyby
-- trzymać tu samą nazwę i doczytywać resztę z wykazu przy każdym odczycie,
-- odinstalowanie rozszerzenia zamieniłoby dołożenie w wiersz bez opisu —
-- Operator zobaczyłby, że model ma narzędzie, ale nie zobaczyłby jakie.
-- Odpis pozycji jest tu więc zapisem faktu („to dołożono i tak to wtedy
-- wyglądało"), a nie zdublowanym katalogiem.
--
-- Rodzaj: dwa, bo po ukośniku idą dwa rodzaje wpisów. Obok powołania narzędzia
-- stoi komenda akcji, która
-- czynność wykonuje i zestawu narzędzi nie dotyka. Kolumna `rodzaj` niesie
-- wartość kontraktu wprost (shared.SlashEntryKind), tak samo jak
-- `agent_pamiec_poziom.poziom` z migracji 121 — kontrakt nie daje dla tego
-- wyliczenia słownika przekładu bazy, więc go tu nie ma.
--
-- CHECK wylicza oba rodzaje, ale dołożone bywają wyłącznie pozycje rodzaju
-- `tool` i tego pilnuje warstwa `dane` odmową, nie baza. Warunek bazy zawężony
-- do jednej wartości byłby drugą regułą tego samego wyboru, gotową rozejść się
-- z pierwszą, gdyby komenda akcji też miała kiedyś zostać przypięta do sesji.
--
-- Źródło: czyja ręka, a nie druga droga. Kolumna `zrodlo` niesie
-- shared.SessionToolSource. Nie jest to wyliczenie
-- dróg dołożenia — droga jest jedna. To zapis sprawcy: `slashCommand`, gdy Operator
-- wpisał komendę z klawiatury, `assistant`, gdy zrobił to za niego asystent.
-- Rozróżnienie ma to samo uzasadnienie, co pole `actor` przy `session.changed`:
-- bez niego ekran nie umiałby napisać „Asystent dołożył", a dołożenie zrobione
-- cudzą ręką wyglądałoby jak własne.
--
-- Czas jest liczbą. `SessionTool.attachedAt` niesie milisekundy epoki, więc
-- kolumna jest INTEGER
-- i przenosi wartość kontraktu w obie strony bez przekładu — tak samo jak
-- `rozszerzenie.zaktualizowano` z migracji 070. Baza nie wstawia własnego
-- „teraz": dwa zegary dla jednego pola byłyby dwiema prawdami.

CREATE TABLE narzedzie_sesji (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    sesja_id       INTEGER NOT NULL REFERENCES sesja(id) ON DELETE CASCADE,

    -- Nazwa pełna ze źródłem, np. `anthropic-skills:skill-creator`. To ona jest
    -- tożsamością dołożenia w obrębie sesji — stąd warunek UNIQUE niżej.
    nazwa_pelna    TEXT    NOT NULL,
    -- Nazwa skrócona — ta, którą Operator faktycznie wpisuje (rozdz. 7.2).
    nazwa_skrocona TEXT    NOT NULL,
    -- Opis pełnym zdaniem. Bez niego wykaz dołożeń cierpiałby na tę samą cichą
    -- degradację, co zestaw narzędzi bez opisów: `:design-audit` nie mówi nic.
    opis           TEXT    NOT NULL,

    rodzaj         TEXT    NOT NULL DEFAULT 'tool'
                           CHECK (rodzaj IN ('tool', 'action')),
    grupa          TEXT    NOT NULL DEFAULT '',
    -- Przedrostek źródła pozycji (`anthropic-skills`, `danaco`) — Operator widzi,
    -- skąd pozycja pochodzi, bez wchodzenia w szczegóły.
    zrodlo_pozycji TEXT    NOT NULL DEFAULT '',

    zrodlo         TEXT    NOT NULL DEFAULT 'slashCommand'
                           CHECK (zrodlo IN ('slashCommand', 'assistant')),
    dolozono       INTEGER NOT NULL DEFAULT 0
);

-- Jedno dołożenie na sesję i nazwę. Powtórzona komenda po ukośniku nie jest
-- błędem — Operator ma prawo nie pamiętać, że już to dołożył — ale nie może
-- rodzić drugiego wiersza. Bez tego warunku wykaz dołożeń pokazywałby to samo
-- narzędzie tyle razy, ile razy padła komenda, a zdjęcie jednego zostawiałoby
-- resztę: zestaw narzędzi rozjechałby się z tym, co Operator widzi.
CREATE UNIQUE INDEX idx_narzedzie_sesji_nazwa ON narzedzie_sesji(sesja_id, nazwa_pelna);

-- Odczyt jest zawsze „co dołożono do TEJ sesji, w kolejności dokładania".
CREATE INDEX idx_narzedzie_sesji_wykaz ON narzedzie_sesji(sesja_id, dolozono, id);
