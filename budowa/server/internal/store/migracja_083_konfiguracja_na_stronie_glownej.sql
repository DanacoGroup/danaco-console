-- `modul.konfigurowany_na_stronie_glownej` — które moduły nastawia się
-- w Strefie 2 Strony głównej.
--
-- Fakt jest niezależny od `modul.rodzaj`: rodzaj orzeka, CZYM moduł jest
-- (`srodowisko_robocze`, `kompozytor`, `repozytorium_plikow`,
-- `sekcja_konfiguracyjna`), ta kolumna — GDZIE moduł się nastawia. Moduły
-- niosące ją prawdziwie nie mają wspólnego rodzaju, więc kolumna nie jest
-- widokiem tamtej.
--
-- SQLite nie ma typu logicznego, więc niesie go INTEGER, a zbioru wartości
-- pilnuje CHECK. `ALTER TABLE ADD COLUMN` zabrania w SQLite wyłącznie PRIMARY
-- KEY i UNIQUE — CHECK przyjmuje i egzekwuje. NOT NULL wymaga wartości
-- domyślnej, bo kolumna dokładana do tabeli z wierszami musi powiedzieć, co
-- znaczy dla wierszy już stojących: moduł bez rozstrzygnięcia na Stronie
-- głównej się nie nastawia.

-- ── Kolumna ─────────────────────────────────────────────────────────────────
ALTER TABLE modul
    ADD COLUMN konfigurowany_na_stronie_glownej INTEGER NOT NULL DEFAULT 0
        CHECK(konfigurowany_na_stronie_glownej IN (0, 1));

-- ── Moduły nastawiane na Stronie głównej ────────────────────────────────────
--
-- Zapis idzie po kodzie, nie po kluczu sztucznym: `modul.kod` jest tożsamością
-- modułu stałą między bazami, a `id` nadaje AUTOINCREMENT i nie musi być takie
-- samo nigdzie indziej. Moduł nieobecny w słowniku po prostu nie zostanie
-- trafiony — krok nie zakłada wiersza, którego zaczyn nie wniósł.
UPDATE modul
   SET konfigurowany_na_stronie_glownej = 1
 WHERE kod IN ('automations', 'agents', 'workspace', 'assistant');
