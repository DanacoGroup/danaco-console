-- Migracja 081 — tryb nałożenia instrukcji eksperta jako pole eksperta.
--
-- Kontrakt niesie `Agent.mode` (`IdentityMode`), a `agent.create`
-- i `agent.update` je przyjmują, ale tabela `agent` nie miała pola trybu — pole
-- było przyjmowane i milcząco wyrzucane. Ta kolumna daje mu nośnik.
--
-- Dwie wartości i asymetria między nimi:
--   'DOLACZ'  — opcja domyślna. Prompt globalny obowiązuje pierwszy, instrukcja
--               eksperta dopisuje się do niego jako zakres użytkownika.
--   'ZASTAP'  — odstępstwo jawne, oznaczane przez Operatora świadomie.
--               Instrukcja eksperta staje się promptem systemowym, a globalny
--               przestaje obowiązywać dla okien tego eksperta.
-- Wartość domyślna kolumny jest odwrotna niż w tożsamości osi, gdzie kontrakt
-- mówi „pominięte znaczy ZASTAP" (`tozsamosc.tryb_domyslny`). To nie jest
-- niekonsekwencja: tożsamość osi jest konfiguracją platformy, modelu albo konta,
-- a ekspert jest warstwą nałożoną na wywołanie. Byt nakładany nie ma prawa
-- zastępować podkładu domyślnie.
--
-- Pole należy do eksperta, nie do warstwy: tryb dotyczy instrukcji modelu jako
-- całości — albo prompt globalny obowiązuje, albo nie. Kolumna `agent_warstwa.tryb`
-- była per warstwa, więc ekspert bez warstw nie miał w niej wiersza wcale,
-- a ekspert z trzema miał trzy kopie jednej prawdy.
--
-- `UPDATE` niżej nadaje odstępstwo tym ekspertom, których choć jedna czynna
-- warstwa z treścią żąda zastąpienia — tą samą regułą, którą liczy rdzeń
-- (`core/tozsamosc_agenta_zrodlo.go`). Warstwa nieczynna albo pusta nie wnosi
-- treści, więc nie wnosi też odstępstwa.
--
-- Kolumna `agent_warstwa.tryb` zostaje i nie jest kasowana. Nie ma już
-- właściciela: kontrakt nie pozwala jej ustawić (`AgentLayer.mode` zdjęte), rdzeń
-- wpisuje w nią stałą, a tryb mieszka przy ekspercie. Jej zdjęcie jest osobnym
-- rozstrzygnięciem i osobną migracją.
--
-- Każdy ekspert założony wcześniej dostaje 'DOLACZ': `ADD COLUMN` z `DEFAULT`
-- nie dotyka ani jednego wiersza danych. Ekspert bez ani jednej warstwy z treścią
-- nie zastępuje niczego niezależnie od tego pola — rdzeń sprawdza to przed
-- rozgałęzieniem (`nakladkaZAgentem`), bo zastąpienie promptu globalnego pustką
-- zostawiłoby model bez konstytucji platformy.

-- ── Tryb nałożenia ───────────────────────────────────────────────────────────
ALTER TABLE agent ADD COLUMN tryb_nakladki TEXT NOT NULL DEFAULT 'DOLACZ'
    CHECK (tryb_nakladki IN ('ZASTAP', 'DOLACZ'));

-- ── Przeniesienie nośnika przejściowego ──────────────────────────────────────
UPDATE agent
   SET tryb_nakladki = 'ZASTAP'
 WHERE id IN (
       SELECT agent_id
         FROM agent_warstwa
        WHERE tryb = 'ZASTAP'
          AND aktywna = 1
          AND TRIM(tresc) <> ''
 );
