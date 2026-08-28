-- Migracja 277 — osiem zakresów izolacji technicznej zapisanych PRZY EKSPERCIE
-- (Permissions Center, grupa „Izolacja techniczna").
--
-- Tabela nie powiela okna konfiguracji punktów izolacji. Tamto okno zapisuje
-- regułę na jednym z siedmiu poziomów zasięgu (tabele `konfiguracja_osi`
-- i `profil_izolacji`); tutaj zapisuje się WARTOŚĆ WYJŚCIOWA WŁAŚCIWA
-- EKSPERTOWI — ta, która obowiązuje wszędzie, gdzie ekspert działa, dopóki nie
-- nadpisze jej reguła z poziomu bardziej szczegółowego. Bez własnego nośnika
-- ustawienie eksperta musiałoby udawać regułę zasięgu, a wtedy nie dałoby się
-- odróżnić „ekspert ma tak ustawione" od „platforma ma tak ustawione".
--
-- Brak wiersza znaczy zakres nieodcięty — stan wyjściowy z rozdziału 6.4
-- Specyfikacji agentów. Dlatego tabela nie zakłada ośmiu wierszy przy
-- założeniu eksperta: komplet ośmiu składa odczyt, a baza trzyma wyłącznie to,
-- co Operator naprawdę przestawił.

CREATE TABLE agent_izolacja_techniczna (
    agent_id INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    -- Wartość wyliczenia `IsolationTechnicalScope` kontraktu. Katalog ośmiu
    -- zakresów stoi w kontrakcie i rdzeń sprawdza wskazanie przy zapisie —
    -- powtarzanie go warunkiem CHECK dałoby drugie źródło tej samej prawdy,
    -- rozjeżdżające się przy każdym rozszerzeniu kontraktu.
    zakres   TEXT    NOT NULL,
    odciety  INTEGER NOT NULL DEFAULT 0 CHECK (odciety IN (0, 1)),
    zapisano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (agent_id, zakres)
);
