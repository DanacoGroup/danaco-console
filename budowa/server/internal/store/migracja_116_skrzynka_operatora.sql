-- Dodaje do tabeli skrzynka_pocztowa kolumny protokołu, źródła nastaw, nazwy wyświetlanej, szyfrowania odbioru i wysyłki oraz domyślności skrzynki Operatora.

ALTER TABLE skrzynka_pocztowa ADD COLUMN protokol TEXT NOT NULL DEFAULT 'imap'
    CHECK(protokol IN ('imap','jmap','pop3'));

ALTER TABLE skrzynka_pocztowa ADD COLUMN zrodlo TEXT NOT NULL DEFAULT 'operator'
    CHECK(zrodlo IN ('operator','urzadzenie','chmura'));

-- Nazwa wyświetlana nadawcy. NULL znaczy „Operator jej nie podał" — list
-- wyjdzie wtedy z samym adresem w polu From, a nie z nazwą zmyśloną z adresu.
ALTER TABLE skrzynka_pocztowa ADD COLUMN nazwa_wyswietlana TEXT;

ALTER TABLE skrzynka_pocztowa ADD COLUMN szyfruj_odbior INTEGER NOT NULL DEFAULT 1
    CHECK(szyfruj_odbior IN (0,1));

ALTER TABLE skrzynka_pocztowa ADD COLUMN szyfruj_wysylke INTEGER NOT NULL DEFAULT 0
    CHECK(szyfruj_wysylke IN (0,1));

ALTER TABLE skrzynka_pocztowa ADD COLUMN domyslna INTEGER NOT NULL DEFAULT 0
    CHECK(domyslna IN (0,1));

CREATE UNIQUE INDEX idx_skrzynka_pocztowa_domyslna
    ON skrzynka_pocztowa(domyslna) WHERE domyslna = 1;
