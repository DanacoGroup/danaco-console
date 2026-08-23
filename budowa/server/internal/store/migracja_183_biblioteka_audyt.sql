-- Migracja 183 — moduł Library: dziennik audytu repozytorium.
--
-- Dziennik jest przyrostowy: wiersz raz dopisany nie jest zmieniany ani
-- kasowany żadną komendą rodziny `library.*`. Stąd brak kolumny zmiany i brak
-- kolumny stanu — wpis opisuje zdarzenie, które już zaszło, a zdarzenie zajść
-- nie przestaje.
--
-- Wskazanie zasobu jest luźne z zamysłu: `plik_kod` niesie identyfikator
-- zewnętrzny zasobu, a nie klucz obcy do `plik_biblioteki(id)`. Klucz obcy
-- z kaskadą zabrałby wpisy razem z zasobem usuniętym trwale — czyli zdjąłby
-- ślad dokładnie tej czynności, dla której dziennik istnieje. Wpis o usunięciu
-- ma przeżyć usunięcie.
--
-- Czynności obejmujące całość repozytorium (wywóz paczki, skanowanie
-- duplikatów) wchodzą bez wskazania zasobu — kolumna jest pusta i to jest
-- stan opisany kontraktem (`LibraryAuditEntry.fileId` opcjonalne).

CREATE TABLE wpis_audytu_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT NOT NULL UNIQUE,
    plik_kod                 TEXT,
    czynnosc                 TEXT NOT NULL
                             CHECK(czynnosc IN ('dostep','zmiana','archiwizacja',
                                                'przywrocenie','eksport','usuniecie',
                                                'utrwalenie')),
    -- Sprawca: Operator, moduł albo model. Wolny tekst, bo sprawcą bywa byt
    -- spoza katalogu kont — moduł źródłowy nie ma wiersza w `konto`.
    sprawca                  TEXT NOT NULL,
    opis                     TEXT,
    chwila                   TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- Dziennik czyta się od najnowszego, zwykle zawężony do jednego zasobu albo do
-- jednej czynności.
CREATE INDEX idx_wpis_audytu_biblioteki_chwila ON wpis_audytu_biblioteki(chwila DESC, id DESC);
CREATE INDEX idx_wpis_audytu_biblioteki_plik ON wpis_audytu_biblioteki(plik_kod, chwila DESC);
CREATE INDEX idx_wpis_audytu_biblioteki_czynnosc ON wpis_audytu_biblioteki(czynnosc, chwila DESC);
