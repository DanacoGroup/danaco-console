-- Karty sesji należą do konta.
--
-- Migracja 406 dopuściła wiele kont, ale praca pozostała wspólna: karta sesji
-- wisiała wyłącznie na środowisku, więc każde konto widziało karty wszystkich
-- pozostałych. Rozstrzygnięcie Właściciela z 29.08.2026 (`decyzje.md`, poz. 22)
-- stanowi, że konta nie mieszają między sobą sesji.
--
-- Wskazanie stoi na karcie sesji, a nie na każdej tabeli pracy z osobna, bo
-- karta jest korzeniem: sesja wisi na karcie, a dokumenty, okna, pamięć
-- i pozostała praca wiszą na sesji przez klucze obce z kasowaniem kaskadowym.
-- Odcięcie przy korzeniu odcina całe poddrzewo i nie zostawia miejsca, w którym
-- rozdzielenie dałoby się ominąć.
--
-- Środowisko zostaje wspólne: cztery środowiska pracy są własnością produktu,
-- nie konta — konto zakłada w nich własne karty.
--
-- Kolumna dopuszcza NULL i nie ma klucza obcego: karty zastane powstały przed
-- rozdzieleniem i nie mają jak wskazać konta wstecz. Instalacja z jednym kontem
-- rozstrzyga jednoznacznie, więc rdzeń czyta karty bez wskazania jako należące
-- do konta najstarszego.
ALTER TABLE karta_sesji ADD COLUMN konto_id INTEGER;

CREATE INDEX idx_karta_sesji_konto ON karta_sesji (konto_id);

-- Wyszukiwanie kart idzie parą (środowisko, konto) — obie wartości naraz, bo
-- okno pyta zawsze o karty jednego konta w jednym środowisku.
CREATE INDEX idx_karta_sesji_srodowisko_konto ON karta_sesji (srodowisko_id, konto_id);
