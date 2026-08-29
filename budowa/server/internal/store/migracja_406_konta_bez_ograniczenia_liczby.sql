-- Konta bez ograniczenia liczby: jedno konto na użytkownika, nie jedno konto
-- na instalację.
--
-- Migracja 125 postawiła tabelę z warunkiem `CHECK (id = 1)`, czytając zamiar
-- „jedno konto na jeden adres" jako „jeden wiersz w całej bazie". Skutkiem
-- rejestracja odmawiała każdemu kolejnemu użytkownikowi, niezależnie od tego,
-- jaki adres podał. Rozstrzygnięcie Właściciela z 29.08.2026 (`decyzje.md`,
-- poz. 22) uchyla tamten model: platforma przyjmuje dowolną liczbę kont,
-- a jednoznaczność ma pilnować adres i login, nie identyfikator wiersza.
--
-- Warunku ze starej tabeli nie da się zdjąć poleceniem ALTER — SQLite nie zna
-- usuwania CHECK. Tabela powstaje więc od nowa, a wiersz istniejący przechodzi
-- do niej bez zmian: instalacja z jednym kontem zostaje z tym samym kontem.

CREATE TABLE konto_nowe (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	login        TEXT    NOT NULL,
	email        TEXT    NOT NULL,
	potwierdzone INTEGER NOT NULL DEFAULT 0,
	utworzono    INTEGER NOT NULL
);

INSERT INTO konto_nowe (id, login, email, potwierdzone, utworzono)
SELECT id, login, email, potwierdzone, utworzono FROM konto_wlasciciela;

DROP TABLE konto_wlasciciela;
ALTER TABLE konto_nowe RENAME TO konto_wlasciciela;

-- Jednoznaczność tożsamości. Porównanie bez względu na wielkość liter, bo
-- „Operator" i „operator" to ten sam login, a adres e-mail w części domenowej
-- wielkości liter nie rozróżnia. Bez tych wskaźników dwa konta o tym samym
-- adresie byłyby stanem osiągalnym, a odzyskanie konta listem nie wiedziałoby,
-- które z nich otworzyć.
CREATE UNIQUE INDEX idx_konto_login  ON konto_wlasciciela (login COLLATE NOCASE);
CREATE UNIQUE INDEX idx_konto_email  ON konto_wlasciciela (email COLLATE NOCASE);

-- Metoda wejścia, sesja bramki i droga potwierdzenia powstały, gdy konto było
-- jedno, więc żadna z nich nie niesie wskazania konta. Przy wielu kontach
-- hasło bez wskazania właściciela otwierałoby cudzą bramkę, a droga
-- potwierdzenia nie wiedziałaby, które konto potwierdza.
--
-- Kolumna dopuszcza NULL i nie ma klucza obcego: wiersze zastane powstały przed
-- tą migracją i nie mają jak wskazać konta wstecz. Wiązanie ich z kontem
-- jedynym byłoby zgadywaniem — instalacja z jednym kontem i tak rozstrzyga
-- jednoznacznie, a rdzeń uzupełnia wskazanie przy pierwszym użyciu.
ALTER TABLE metoda_uwierzytelnienia  ADD COLUMN konto_id INTEGER;
ALTER TABLE sesja_bramki             ADD COLUMN konto_id INTEGER;
ALTER TABLE potwierdzenie_tozsamosci ADD COLUMN konto_id INTEGER;

CREATE INDEX idx_metoda_konto        ON metoda_uwierzytelnienia  (konto_id);
CREATE INDEX idx_sesja_konto         ON sesja_bramki             (konto_id);
CREATE INDEX idx_potwierdzenie_konto ON potwierdzenie_tozsamosci (konto_id);

-- Hasło jest jedno NA KONTO, nie jedno w tabeli. Wskaźnik z migracji 071
-- dopuszczał wyłącznie jeden wiersz z `kotwica = 1` w całej bazie, więc drugie
-- konto nie miało gdzie zapisać swojego hasła i rejestracja odbijała się
-- o kolizję, choć login i adres były wolne.
DROP INDEX idx_metoda_uwierzytelnienia_kotwica;

-- Wiersz zastany nie ma wskazania konta; `COALESCE` sprowadza go do zera, więc
-- jedna instalacja sprzed migracji ma nadal dokładnie jedno hasło bez konta.
CREATE UNIQUE INDEX idx_metoda_kotwica_konta
    ON metoda_uwierzytelnienia (COALESCE(konto_id, 0)) WHERE kotwica = 1;
