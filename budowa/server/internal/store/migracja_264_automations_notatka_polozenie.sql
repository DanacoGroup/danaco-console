-- Migracja 264 — notatka przy kroku i położenie węzła na kanwie
-- (`automation.step.note.set`, `automation.step.layout.set`).
--
-- Jedna tabela na oba, kluczowana KODEM kroku, a nie kolumny w tabeli
-- `krok_automatyki`. Powód jest w zapisie definicji: `ZapiszKroki` podmienia
-- komplet wierszy kroków przy każdym zapisie Workflow Buildera. Kolumna
-- kasowałaby się więc przy każdej zmianie nazwy kroku — Operator zastawałby
-- kanwę ułożoną od nowa i notatki zdjęte, choć niczego takiego nie żądał.
--
-- Adnotacja kroku, którego już nie ma, zostaje w tabeli i nikogo nie boli:
-- odczyt idzie od kroków do adnotacji, a nie odwrotnie. Krok przywrócony
-- z wersji wcześniejszej zastaje wtedy swoje położenie na miejscu.
CREATE TABLE adnotacja_kroku_automatyki (
    automatyka_id INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    krok_kod      TEXT    NOT NULL,
    notatka       TEXT,
    wspolrzedna_x INTEGER NOT NULL DEFAULT 0,
    wspolrzedna_y INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (automatyka_id, krok_kod)
);

-- Odwołania do skarbca idą kolumną KROKU, a nie adnotacją.
--
-- Rozdział jest rozmyślny i przeciwny temu wyżej. Notatka i położenie są
-- własnością Operatora nad węzłem — mają przeżyć podmianę wierszy kroków.
-- Odwołanie do poświadczenia jest częścią SAMEJ definicji kroku: krok, który
-- przestał wołać interfejs zewnętrzny, przestaje potrzebować klucza. Kolumna
-- kroku znika więc razem z krokiem i to jest zachowanie właściwe.
--
-- Bez tej kolumny `automation.secret.remove` nie miał jak nazwać kroków, które
-- straciły pokrycie (pole `referencingStepIds` kontraktu): odwołania przychodziły
-- żądaniem zapisu definicji i ginęły w locie. Skarbiec zdejmowałby wtedy klucz
-- w ciszy, a Operator dowiadywałby się o skutku dopiero z nieudanego przebiegu.
--
-- Wykaz leży zapisem strukturalnym, bo kontrakt niesie go wykazem tekstów
-- ustalanym w całości wraz z krokiem.
ALTER TABLE krok_automatyki ADD COLUMN odwolania_sekretow TEXT;
