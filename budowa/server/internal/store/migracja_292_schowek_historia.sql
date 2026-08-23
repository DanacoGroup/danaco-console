-- Migracja 292 — historia schowka Operatora (rodzina `clipboard.*`).
--
-- Schowek należy do maszyny Operatora i rdzeń go NIE czyta. Ta tabela nie
-- udaje dostępu do cudzego schowka: okno, w którym Operator skopiował treść,
-- oddaje ją rdzeniowi wprost (`clipboard.push`), a rdzeń daje jej trwałość —
-- historia przestaje ginąć razem z kartą. Wklejenie jest ruchem powrotnym:
-- `clipboard.list` oddaje treść oknu, a wstawia ją okno, u siebie.
--
-- Odcisk treści (`odcisk`) jest kluczem powtórzenia. Kontrakt mówi wprost:
-- powtórzenie tej samej treści nie mnoży wpisów, tylko podnosi wpis zastany na
-- czoło wykazu. Bez UNIQUE na odcisku rozstrzygałby to odczyt-i-zapis, czyli
-- wyścig dwóch okien kopiujących naraz.
--
-- Wpis wrażliwy jest wierszem jak każdy inny — różni się polityką, nie
-- schematem: nie wchodzi do eksportu i wygasa wedle zasady retencji. Rdzeń nie
-- szyfruje go tutaj, bo baza rdzenia jest tym samym magazynem co reszta stanu;
-- osobne szyfrowanie jednej kolumny dawałoby poczucie ochrony bez ochrony.

CREATE TABLE wpis_schowka (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rodzaj                   TEXT    NOT NULL DEFAULT 'text',
    tresc                    TEXT    NOT NULL,
    odcisk                   TEXT    NOT NULL UNIQUE,
    zajawka                  TEXT,
    rozmiar_bajtow           INTEGER NOT NULL DEFAULT 0,
    przypiety                INTEGER NOT NULL DEFAULT 0,
    wrazliwy                 INTEGER NOT NULL DEFAULT 0,
    okno_zrodlowe            TEXT,
    utworzono                INTEGER NOT NULL,
    uzyto                    INTEGER
);
-- Wykaz idzie przypiętymi na czele, potem od najnowszego — dokładnie tak, jak
-- opisuje `clipboard.list`.
CREATE INDEX idx_wpis_schowka_wykaz ON wpis_schowka(przypiety DESC, utworzono DESC, id DESC);
CREATE INDEX idx_wpis_schowka_rodzaj ON wpis_schowka(rodzaj, utworzono DESC);
