-- Migracja dokłada kolumnę wskazującą, które moduły nastawia się w strefie drugiej
-- strony głównej, niezależną od rodzaju modułu.

-- Kolumna logiczna zapisana jako liczba, bo baza nie ma typu logicznego wprost, ze
-- zbiorem wartości pilnowanym warunkiem.
ALTER TABLE modul
    ADD COLUMN konfigurowany_na_stronie_glownej INTEGER NOT NULL DEFAULT 0
        CHECK(konfigurowany_na_stronie_glownej IN (0, 1));

-- Zapis idzie po kodzie modułu, tożsamości stałej między bazami, nie po kluczu sztucznym
-- nadawanym autonumeracją.
UPDATE modul
   SET konfigurowany_na_stronie_glownej = 1
 WHERE kod IN ('automations', 'agents', 'workspace', 'assistant');
