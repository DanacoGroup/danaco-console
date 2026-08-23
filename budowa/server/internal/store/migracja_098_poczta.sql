-- Poczta rdzenia: skrzynka pocztowa, listy odebrane i wysłane.
--
-- Host i port są konfiguracją skrzynki, nie stałą kodu — schemat nie przesądza,
-- czy skrzynka stoi u dostawcy zewnętrznego, na własnym serwerze, czy na
-- stanowisku probierczym.
--
-- Protokoły: odbiór POP3 over TLS, wysyłka SMTP z STARTTLS — oba z biblioteki
-- standardowej Go (`crypto/tls` z `net/textproto` oraz `net/smtp`). Biblioteka
-- standardowa nie ma IMAP, a jego użycie wymagałoby nowej zależności w go.mod.
--
-- Hasło skrzynki nie leży w bazie: kolumna `haslo_odwolanie` niesie wyłącznie
-- odwołanie do sejfu poświadczeń (`sejf:poczta:<kod>`, wzorem kont
-- w `core/handlers_konta_poswiadczenia.go`). Kolumny na treść sekretu w schemacie
-- nie ma.
--
-- `tls_weryfikacja` jest polem konfiguracji: stanowisko probiercze mówi
-- prawdziwym protokołem przez prawdziwe gniazdo, lecz jego certyfikat jest
-- samopodpisany. Wyłączenie weryfikacji łańcucha jest jawnym zapisem w wierszu
-- skrzynki, domyślnie weryfikacja jest włączona.
--
-- `identyfikator_listu` to UIDL serwera POP3, unikalny w obrębie skrzynki —
-- stąd UNIQUE(skrzynka_id, identyfikator_listu): ten sam list odczytany
-- w dwóch taktach obserwatora daje jeden wiersz. Rdzeń nie kasuje listów
-- z serwera (DELE zostaje niewysłane), bo odbiór nie może być jedynym
-- istnieniem listu.
--
-- `list_wyslany` zapisuje także wysyłkę nieudaną — `powodzenie = 0` z treścią
-- błędu w kolumnie `blad`.

-- ── Skrzynka pocztowa — konfiguracja odbioru i wysyłki ───────────────────────
CREATE TABLE skrzynka_pocztowa (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    kod             TEXT    NOT NULL UNIQUE,
    -- Adres skrzynki: nadawca listów wychodzących i adresat przychodzących.
    adres           TEXT    NOT NULL,
    host_odbioru    TEXT    NOT NULL,
    port_odbioru    INTEGER NOT NULL DEFAULT 995,
    host_wysylki    TEXT    NOT NULL,
    port_wysylki    INTEGER NOT NULL DEFAULT 587,
    uzytkownik      TEXT    NOT NULL,
    -- Odwołanie sejfu poświadczeń; NULL znaczy skrzynkę bez hasła
    -- (stanowisko probiercze) — brak nie jest błędem.
    haslo_odwolanie TEXT,
    tls_weryfikacja INTEGER NOT NULL DEFAULT 1 CHECK(tls_weryfikacja IN (0,1)),
    -- Takt obserwatora poczty w sekundach — konfiguracja skrzynki, bo skrzynka
    -- odpytywana często i skrzynka rzadko czynna nie mają wspólnego taktu.
    takt_sekundy    INTEGER NOT NULL DEFAULT 60 CHECK(takt_sekundy > 0),
    aktywna         INTEGER NOT NULL DEFAULT 1 CHECK(aktywna IN (0,1)),
    utworzono       TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- ── Dziennik listów odebranych ───────────────────────────────────────────────
-- `przetworzony` znaczy: krok `poczta.odbierz` zabrał list do streszczenia.
-- Obserwator poczty zapisuje list i wyzwala automatyki, ale znacznika nie
-- stawia — postawienie go jest pracą kroku, nie odbiornika (jedno
-- miejsce prawdy o tym, co automatyka już widziała).
CREATE TABLE list_odebrany (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    skrzynka_id         INTEGER NOT NULL
                                REFERENCES skrzynka_pocztowa(id) ON DELETE CASCADE,
    identyfikator_listu TEXT    NOT NULL,
    nadawca             TEXT    NOT NULL,
    temat               TEXT    NOT NULL DEFAULT '',
    -- Chwila z nagłówka Date listu (ISO 8601 UTC); NULL, gdy list daty nie
    -- niesie — rdzeń nie zmyśla chwili nadania, chwilę odbioru
    -- trzyma `odebrano`.
    chwila              TEXT,
    tresc               TEXT    NOT NULL DEFAULT '',
    odebrano            TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    przetworzony        INTEGER NOT NULL DEFAULT 0 CHECK(przetworzony IN (0,1)),
    UNIQUE(skrzynka_id, identyfikator_listu)
);

-- Krok `poczta.odbierz` pyta wyłącznie o listy nieprzetworzone swojej skrzynki.
CREATE INDEX idx_list_odebrany_nieprzetworzone
    ON list_odebrany(skrzynka_id, id)
    WHERE przetworzony = 0;

-- ── Ślad listów wysłanych ────────────────────────────────────────────────────
CREATE TABLE list_wyslany (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    skrzynka_id     INTEGER NOT NULL
                            REFERENCES skrzynka_pocztowa(id) ON DELETE CASCADE,
    adresat         TEXT    NOT NULL,
    temat           TEXT    NOT NULL DEFAULT '',
    tresc           TEXT    NOT NULL DEFAULT '',
    -- Identyfikator listu odebranego, na który ten odpowiada; NULL dla listu
    -- nadanego bez odpowiedzi.
    w_odpowiedzi_na TEXT,
    chwila          TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    powodzenie      INTEGER NOT NULL CHECK(powodzenie IN (0,1)),
    blad            TEXT,
    -- Wysyłka udana nie ma błędu; nieudana musi go mieć.
    CHECK ((powodzenie = 1 AND blad IS NULL) OR (powodzenie = 0 AND blad IS NOT NULL))
);

CREATE INDEX idx_list_wyslany_skrzynka ON list_wyslany(skrzynka_id, id DESC);
