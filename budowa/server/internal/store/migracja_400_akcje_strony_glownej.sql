-- Migracja 400 dodaje do katalogu akcji pozycje kafli komponentów własnych
-- oraz funkcji listwy ustawień Strony głównej.

WITH katalog(kod, nazwa, opis, ikona, komenda, kolejnosc) AS (
    VALUES
        -- Strefa 2 — kafle komponentów własnych, w kolejności z opracowania.
        ('component.create.automations', 'Zbuduj automatykę',
         'Zakłada komponent własny rodzaju automations — automatykę modułu Automations',
         'automatyzacja',  'component.create', 14),
        ('component.create.agents',      'Skonfiguruj agenta',
         'Zakłada komponent własny rodzaju agents — eksperta modułu Agents',
         'agent',          'component.create', 15),
        ('component.create.workspace',   'Załóż projekt',
         'Zakłada komponent własny rodzaju workspace — projekt przestrzeni roboczej',
         'folder',         'component.create', 16),
        ('component.create.assistant',   'Ustaw profil asystenta',
         'Zakłada komponent własny rodzaju assistant — profil asystenta',
         'uzytkownik',     'component.create', 17),
        ('component.list',               'Komponenty własne',
         'Zwraca komponenty własne zapisane w Strefie 2 Strony głównej',
         'tabela-danych',  'component.list',   18),
        ('component.update',             'Zmień komponent własny',
         'Zmienia komponent własny; pola pominięte zostają bez zmian',
         'olowek',         'component.update', 19),
        ('component.delete',             'Usuń komponent własny',
         'Usuwa komponent własny',
         'kosz',           'component.delete', 20),

        -- Strefa 3 — dwie funkcje globalne listwy ustawień.
        ('mobile.process.list',          'Mobile',
         'Procesy platformy widoczne z urządzenia mobilnego',
         'telefon',        'mobile.process.list', 21),
        ('aod.status.get',               'Always On Display',
         'Stan nakładki globalnego agenta towarzyszącego',
         'aktywnosc',      'aod.status.get',      22)
)
INSERT INTO akcja (kod, nazwa, opis, ikona, poziom_zasiegu_id, klucz_zasiegu,
                   komenda, warunek_dostepnosci, kolejnosc, aktywna)
SELECT katalog.kod, katalog.nazwa, katalog.opis, katalog.ikona, p.id, '',
       katalog.komenda, '', katalog.kolejnosc, 1
  FROM katalog
  JOIN poziom_zasiegu p ON p.kod = 'globalny'
 WHERE true
ON CONFLICT(kod) DO NOTHING;
