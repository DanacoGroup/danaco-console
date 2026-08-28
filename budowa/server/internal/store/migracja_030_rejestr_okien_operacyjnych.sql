-- Migracja 030 — rejestr okien operacyjnych: definicja okna osobno, przypisanie
-- okna do modułu osobno.
--
-- Okno jest bytem katalogu, a obecność okna w module jest relacją między dwoma
-- bytami. Rozdzielamy jedno od drugiego: `okno_operacyjne` niesie definicję —
-- jeden wiersz na okno; macierz `okno_operacyjne_modul` (migracja 031) niesie
-- wystąpienia wraz z kolejnością właściwą danemu modułowi. Dzięki temu okno
-- wspólne dla wielu modułów ma jeden wiersz definicji i wiele przypięć, a jego
-- pozycja może być inna w każdym module — czego jedna kolumna `kolejnosc` na
-- definicji nie potrafi wyrazić.
--
-- Kolumna `modul_id` odchodzi, więc tabela powstaje na nowo i przejmuje nazwę
-- starej. Katalog jest słownikiem wnoszonym migracją — nie ma w nim danych
-- operatora, dlatego wiersze zastępujemy kompletem inwentarza. Kolumna
-- `kolejnosc` definicji porządkuje okno w obrębie jego kategorii; kolejność
-- w module należy do macierzy.

CREATE TABLE okno_operacyjne_nowe (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    kod        TEXT    NOT NULL UNIQUE,
    nazwa      TEXT    NOT NULL,
    rola       TEXT    NOT NULL
                       CHECK(rola IN ('wiodace','pomocnicze','monitor','kreator','zarzadca')),
    kategoria  TEXT    NOT NULL
                       CHECK(kategoria IN ('komunikacja','edycja','podglad','repozytorium',
                                           'konstruktor','kolejka','monitor','narzedzia',
                                           'zrodla','konfiguracja')),
    kolejnosc  INTEGER NOT NULL DEFAULT 0,
    aktywne    INTEGER NOT NULL DEFAULT 1 CHECK(aktywne IN (0,1))
);

DROP TABLE okno_operacyjne;
ALTER TABLE okno_operacyjne_nowe RENAME TO okno_operacyjne;
CREATE INDEX idx_okno_operacyjne_kategoria ON okno_operacyjne(kategoria, kolejnosc);

-- ── Definicje okien. Kod jest bezmodułowy, bo definicja nie należy do modułu. ─
INSERT INTO okno_operacyjne (kod, nazwa, rola, kategoria, kolejnosc)
VALUES
    ('chat-window',            'Chat Window',                        'wiodace',    'komunikacja',  1),
    ('voice-console',          'Voice Console',                      'wiodace',    'komunikacja',  2),
    ('executor-chat-1',        'Executor Chat (Executor 1)',         'wiodace',    'komunikacja',  3),
    ('executor-chat-2',        'Executor Chat (Executor 2)',         'wiodace',    'komunikacja',  4),
    ('coordinator-chat',       'Coordinator Chat',                   'zarzadca',   'komunikacja',  5),

    ('studio-editor',          'Studio Editor',                      'wiodace',    'edycja',       1),
    ('code-editor',            'Code Editor',                        'wiodace',    'edycja',       2),
    ('design-board',           'Design Board',                       'wiodace',    'edycja',       3),
    ('source-panel',           'Source Panel',                       'wiodace',    'edycja',       4),
    ('translation-panels',     'Translation Panels',                 'wiodace',    'edycja',       5),
    ('frontend-workspace',     'Frontend Workspace',                 'wiodace',    'edycja',       6),
    ('backend-workspace',      'Backend Workspace',                  'wiodace',    'edycja',       7),

    ('diff-grep-panel',        'Diff/Grep Panel',                    'pomocnicze', 'podglad',      1),
    ('preview-window',         'Preview Window',                     'pomocnicze', 'podglad',      2),
    ('file-preview',           'File Preview',                       'pomocnicze', 'podglad',      3),
    ('build-output',           'Build Output',                       'monitor',    'podglad',      4),
    ('output-console',         'Output Console',                     'monitor',    'podglad',      5),
    ('export-panel',           'Export Panel',                       'pomocnicze', 'podglad',      6),

    ('library-explorer',       'Library Explorer',                   'wiodace',    'repozytorium', 1),
    ('project-library',        'Project Library',                    'zarzadca',   'repozytorium', 2),
    ('session-repository',     'Session Repository',                 'zarzadca',   'repozytorium', 3),
    ('versioning-panel',       'Versioning Panel',                   'zarzadca',   'repozytorium', 4),
    ('assets-panel',           'Assets Panel',                       'zarzadca',   'repozytorium', 5),
    ('context-memory',         'Context Memory',                     'zarzadca',   'repozytorium', 6),

    ('workflow-builder',       'Workflow Builder',                   'kreator',    'konstruktor',  1),
    ('agent-builder',          'Agent Builder',                      'kreator',    'konstruktor',  2),
    ('product-builder',        'Product Builder',                    'wiodace',    'konstruktor',  3),
    ('prompt-builder',         'Prompt Builder',                     'kreator',    'konstruktor',  4),
    ('architecture-designer',  'Architecture Designer',              'kreator',    'konstruktor',  5),
    ('report-builder',         'Report Builder',                     'kreator',    'konstruktor',  6),

    ('scheduler',              'Scheduler',                          'zarzadca',   'kolejka',      1),
    ('queue-manager',          'Queue Manager',                      'zarzadca',   'kolejka',      2),
    ('orchestrator',           'Orchestrator',                       'kreator',    'kolejka',      3),

    ('execution-monitor',      'Execution Monitor',                  'monitor',    'monitor',      1),
    ('process-monitor',        'Process Monitor',                    'monitor',    'monitor',      2),
    ('actions-monitor',        'Actions Monitor',                    'monitor',    'monitor',      3),
    ('activity-feed',          'Activity Feed',                      'monitor',    'monitor',      4),
    ('diagnostics-center',     'Diagnostics Center',                 'wiodace',    'monitor',      5),
    ('results-analyzer',       'Results Analyzer',                   'monitor',    'monitor',      6),
    ('logs-viewer',            'Logs Viewer',                        'monitor',    'monitor',      7),
    ('debate-panel',           'Debate Panel',                       'monitor',    'monitor',      8),

    ('tools-panel',            'Tools Panel',                        'pomocnicze', 'narzedzia',    1),
    ('glossary-manager',       'Glossary Manager',                   'zarzadca',   'narzedzia',    2),
    ('tags-collections',       'Tags & Collections',                 'zarzadca',   'narzedzia',    3),
    ('permissions-center',     'Permissions Center',                 'zarzadca',   'narzedzia',    4),
    ('skills-manager',         'Skills Manager',                     'zarzadca',   'narzedzia',    5),
    ('moderator-panel',        'Moderator Panel',                    'zarzadca',   'narzedzia',    6),
    ('git-panel',              'Git Panel',                          'zarzadca',   'narzedzia',    7),
    ('agent-manager',          'Agent Manager',                      'zarzadca',   'narzedzia',    8),
    ('deployment-panel',       'Deployment Panel',                   'zarzadca',   'narzedzia',    9),
    ('terminal-tabs',          'Terminal Tabs',                      'wiodace',    'narzedzia',   10),
    ('project-tree',           'Project Tree',                       'pomocnicze', 'narzedzia',   11),

    ('sources-panel',          'Sources Panel',                      'pomocnicze', 'zrodla',       1),
    ('sources-manager',        'Sources Manager',                    'zarzadca',   'zrodla',       2),
    ('findings-panel',         'Findings Panel',                     'pomocnicze', 'zrodla',       3),
    ('notes-panel',            'Notes Panel',                        'pomocnicze', 'zrodla',       4),
    ('browser-window',         'Browser Window',                     'wiodace',    'zrodla',       5),
    ('research-workspace',     'Research Workspace',                 'wiodace',    'zrodla',       6),
    ('errors-panel',           'Errors Panel',                       'pomocnicze', 'zrodla',       7),
    ('recommendations-panel',  'Recommendations Panel',              'pomocnicze', 'zrodla',       8),
    ('consensus-panel',        'Consensus Panel',                    'pomocnicze', 'zrodla',       9),
    ('model-panels',           'Model Panels',                       'wiodace',    'zrodla',      10),

    ('instructions-panel',     'Instructions Panel',                 'pomocnicze', 'konfiguracja', 1),
    ('model-configuration',    'Model Configuration',                'pomocnicze', 'konfiguracja', 2),
    ('connectors-manager',     'Connectors Manager',                 'zarzadca',   'konfiguracja', 3),
    ('punkty-izolacji',        'Okno konfiguracji punktów izolacji', 'zarzadca',   'konfiguracja', 4),
    ('project-dashboard',      'Project Dashboard',                  'wiodace',    'konfiguracja', 5),
    ('okno-konfiguracji',      'Okno Konfiguracji',                  'zarzadca',   'konfiguracja', 6),
    ('okno-ustawien',          'Okno Ustawień',                      'zarzadca',   'konfiguracja', 7),
    ('strona-glowna',          'Strona główna',                      'wiodace',    'konfiguracja', 8)
ON CONFLICT(kod) DO NOTHING;
