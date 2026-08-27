-- Migracja dokłada tryb nałożenia instrukcji eksperta jako pole eksperta, zastępujące
-- dotąd milcząco odrzucane pole kontraktu.

-- Tryb nałożenia rozstrzyga, czy instrukcja eksperta dopisuje się do promptu globalnego,
-- czy go zastępuje w oknach tego eksperta.
ALTER TABLE agent ADD COLUMN tryb_nakladki TEXT NOT NULL DEFAULT 'DOLACZ'
    CHECK (tryb_nakladki IN ('ZASTAP', 'DOLACZ'));

-- Przeniesienie nośnika przejściowego: eksperci z choć jedną czynną warstwą treści
-- żądającą zastąpienia dostają odstępstwo.
UPDATE agent
   SET tryb_nakladki = 'ZASTAP'
 WHERE id IN (
       SELECT agent_id
         FROM agent_warstwa
        WHERE tryb = 'ZASTAP'
          AND aktywna = 1
          AND TRIM(tresc) <> ''
 );
