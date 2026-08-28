-- Migracja 058 — kolumna `motto` w tabeli `srodowisko` wraz z treścią czterech mott.
--
-- `motto` stoi osobno od `opis`, bo to dwa różne zdania o środowisku, drukowane
-- w dwóch różnych miejscach: `opis` jest jednozdaniowym opisem trybu pracy na
-- karcie wejścia strony głównej, a `motto` — podtytułem w nagłówku kolumny
-- matrycy Mission Control (klasa `mc-kolumna__motto`), łamanym do jednego
-- wiersza obok nazwy. Jedna kolumna na oba fakty zmusiłaby jedną ze stron do
-- skracania cudzego zdania w locie.
--
-- Kolumna dopuszcza NULL i nie ma DEFAULT: środowisko, dla którego hasła nie
-- ustalono, stoi z NULL-em — odróżnialnym od pustego napisu, który znaczy
-- „motta nie ma". Warunek NOT NULL wymuszałby hasło wpisane byle jak przy
-- zakładaniu wiersza, a warstwa widoku pokazuje brak źródła jako wartość pustą
-- (`MOTTO_BEZ_ZRODLA = ''`), nie jako tekst ułożony przez siebie.
--
-- Dopasowanie idzie po `kod`, nie po `nazwa`: `kod` jest jedyną kolumną
-- z warunkiem UNIQUE i CHECK na zamkniętym zbiorze czterech wartości, więc
-- UPDATE trafia w dokładnie jeden wiersz albo w żaden. Kolumna `nazwa` nie ma
-- warunku na wartość i bywa zmieniana.
--
-- `ALTER TABLE ... ADD COLUMN` w SQLite przepisuje sam nagłówek schematu, nie
-- tabelę; kolumna dopuszczająca NULL i bez DEFAULT nie dotyka ani jednego
-- wiersza danych. ALTER i UPDATE-y idą w jednej transakcji (migracje.go stosuje
-- krok jednym `Exec`), więc schemat nie zostanie zastosowany bez treści ani
-- treść bez schematu. Na bazie bez wiersza o danym kodzie UPDATE zmienia zero
-- wierszy i nie jest błędem — wierszy ten krok nie zakłada.

ALTER TABLE srodowisko ADD COLUMN motto TEXT;

UPDATE srodowisko SET motto = 'Pracuj z wiedzą, dokumentami i modelami AI w jednym miejscu.'
 WHERE kod = 'talkin';
UPDATE srodowisko SET motto = 'Organizuj projekty, procesy i zadania wspierane przez sztuczną inteligencję.'
 WHERE kod = 'workspace';
UPDATE srodowisko SET motto = 'Projektuj, twórz i rozwijaj oprogramowanie wspólnie z AI.'
 WHERE kod = 'codestudio';
UPDATE srodowisko SET motto = 'Zarządzaj zespołem modeli AI realizujących złożone procesy.'
 WHERE kod = 'multitaskingai';
