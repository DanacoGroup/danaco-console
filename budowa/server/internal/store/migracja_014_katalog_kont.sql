-- Migracja 014 — katalog kont modeli i kont programów code CLI.
--
-- Migracja rozszerza tabelę `konto` o pola wymagane przez strukturę Account
-- kontraktu: dostawcę, model domyślny, adres punktu końcowego, oznaczenie konta
-- domyślnego oraz stan konta w rotacji. Konta prowadzi jedna tabela — drugiego
-- bytu na to samo nie ma.
--
-- Rodzaj konta. CHECK(rodzaj IN ('cli','api','sdk')) pokrywa się ze słownikami
-- kontraktu shared.WartosciBazyAccountKind / WartosciKontraktuAccountKind, które
-- wprost wskazują kolumnę `konto.rodzaj`. Konto programu code CLI to wiersz
-- rodzaju 'cli', konto modelu — 'api' albo 'sdk'.
--
-- Rotacja. Algorytm rotacji mieszka w `injection.PulaKont`. Kolumny `stan`
-- i `wyczerpane_do` są zapisem tego, co pula już rozpoznała, oraz przełącznikiem
-- Operatora (stan 'zawieszone'), a nie drugą implementacją rotacji. Katalog
-- rozstrzyga, które konta wolno brać i w jakiej kolejności; kiedy przejść na
-- następne, rozstrzyga pula. `aktywne` mówi, czy konto w ogóle bierze udział
-- w pracy, `stan` — czy jest w tej chwili zdatne; to dwa różne pytania, stąd dwie
-- kolumny.
--
-- Poświadczenia. `poswiadczenie_odwolanie` jest jedynym miejscem danych
-- dostępowych i niesie wyłącznie nazwę wpisu w magazynie sekretów albo ścieżkę
-- profilu. Migracja nie dokłada kolumny na sekret, a odczyt katalogu w warstwie
-- `dane` zwraca wyłącznie znacznik „poświadczenie jest zapisane".
--
-- Relacja konto ↔ kanał modelu. `kanal_modelu.konto_id` (ON DELETE SET NULL) jest
-- jedynym wiązaniem tych dwóch bytów. Jeden kanał wskazuje konto preferowane,
-- jedno konto może obsługiwać wiele kanałów. Kanał jest definicją rozmowy
-- z modelem (jak wołać, jakim modelem), konto — profilem uwierzytelnienia.
-- Usunięcie konta odłącza kanały, ale ich nie kasuje; kontrakt oddaje to polem
-- detachedChannelIds odpowiedzi account.remove. Rotacja bierze konta tego samego
-- rodzaju, więc kanał pracuje dalej także wtedy, gdy jego konto preferowane
-- wyczerpało limit.

ALTER TABLE konto ADD COLUMN dostawca       TEXT    NOT NULL DEFAULT '';
ALTER TABLE konto ADD COLUMN model_domyslny TEXT;
ALTER TABLE konto ADD COLUMN adres_bazowy   TEXT;
ALTER TABLE konto ADD COLUMN domyslne       INTEGER NOT NULL DEFAULT 0
                                            CHECK(domyslne IN (0,1));
ALTER TABLE konto ADD COLUMN stan           TEXT    NOT NULL DEFAULT 'aktywne'
                                            CHECK(stan IN ('aktywne','wyczerpane','zawieszone'));
ALTER TABLE konto ADD COLUMN wyczerpane_do  TEXT;

-- Konto domyślne jest dokładnie jedno na rodzaj — pilnuje tego indeks częściowy
-- w bazie, nie warunek w kodzie aplikacji.
CREATE UNIQUE INDEX idx_konto_domyslne_rodzaj ON konto(rodzaj) WHERE domyslne = 1;

-- Odczyt puli rotacji: rodzaj → konta czynne w kolejności rotacji.
CREATE INDEX idx_konto_rotacja ON konto(rodzaj, aktywne, kolejnosc);
