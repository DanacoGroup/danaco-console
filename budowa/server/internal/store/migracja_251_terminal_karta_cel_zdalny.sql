-- Migracja 251 — adres celu karty powłoki zdalnej.
--
-- Kontrakt komendy `terminal.session.open` dostał pola `remoteTarget`,
-- `remotePort` i `hostId`, a byt `TerminalSession` — pole `remoteTarget`. Do tej
-- pory adres powłoki zdalnej wchodził zmienną środowiska `SSH_TARGET`, a zmienne
-- środowiska karty z zamysłu NIE MAJĄ kolumny (bywają nośnikiem poświadczeń,
-- a baza nie jest sejfem). Skutek był taki, że karta zdalna odtworzona po
-- restarcie rdzenia traciła adres i pierwsze polecenie kończyło się odmową.
--
-- Adres celu poświadczeniem nie jest — jest tym samym, co widnieje w wykazie
-- książki hostów i w nazwie karty — więc jego kolumna niczego z sejfu do bazy
-- nie przenosi. Hasło i klucz zostają poza tabelą tak samo jak dotąd.
--
-- Kolumny dokładają się osobnymi poleceniami ALTER, bez przebudowy tabeli:
-- warunek CHECK na kolumnie `powloka` zostaje nietknięty, a wiersze zastane
-- dostają wartości puste, czyli „karta lokalna”.

ALTER TABLE terminal_karta ADD COLUMN cel_zdalny  TEXT NOT NULL DEFAULT '';
ALTER TABLE terminal_karta ADD COLUMN port_zdalny INTEGER;
ALTER TABLE terminal_karta ADD COLUMN host_kod    TEXT;
