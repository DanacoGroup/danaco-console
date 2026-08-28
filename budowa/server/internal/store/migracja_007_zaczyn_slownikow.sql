-- Migracja 007 — zaczyn słowników platformy.
--
-- Łańcuch `srodowisko → modul → karta_sesji → sesja → okno_komunikacji →
-- wiadomosc` stoi na więzach klucza obcego. Bez wierszy w `srodowisko` i `modul`
-- nie da się zapisać żadnego okna komunikacji, a więc i żadnej wiadomości.
--
-- Każde wstawienie kończy się ON CONFLICT DO NOTHING, więc migracja przechodzi
-- także na bazie, w której część wierszy już jest. Klauzula `WHERE true` przed
-- ON CONFLICT jest wymogiem składni SQLite dla INSERT ... SELECT z upsertem.

-- ── Cztery środowiska ──────────────────────────────────────────────────────────
INSERT INTO srodowisko (kod, nazwa, opis, kolejnosc, aktywne)
VALUES
    ('talkin',         'TalkIn',         'Wiedza, komunikacja i praca z treścią',        1, 1),
    ('workspace',      'WorkSpace',      'Produktywność, organizacja i realizacja projektów', 2, 1),
    ('codestudio',     'CodeStudio',     'Programowanie',                                3, 1),
    ('multitaskingai', 'MultitaskingAI', 'Orkiestracja autonomicznej pracy ciągłej',     4, 1)
ON CONFLICT(kod) DO NOTHING;

-- ── Piętnaście modułów ─────────────────────────────────────────────────────────
INSERT INTO modul (kod, nazwa, opis, aktywny)
VALUES
    ('studio',      'Studio',      'Praca z dokumentem: edycja, operacje kontekstowe, wersjonowanie sesji', 1),
    ('workspace',   'Workspace',   'Projekt: zadania, zasoby, instrukcje systemowe, pamięć kontekstu',      1),
    ('automations', 'Automations', 'Workflow, harmonogram, kolejki i orkiestracja',                         1),
    ('browser',     'Browser',     'Przeglądanie stron z gromadzeniem źródeł i notatek',                    1),
    ('research',    'Research',    'Badanie: źródła, ustalenia, raport',                                    1),
    ('library',     'Library',     'Repozytorium plików, wersji, etykiet i kolekcji',                       1),
    ('translate',   'Translate',   'Tłumaczenie równoległe z glosariuszem',                                 1),
    ('roundtable',  'Roundtable',  'Debata wielu modeli z moderacją',                                       1),
    ('design',      'Design',      'Praca wizualna: kompozycja, zasoby, generowanie',                       1),
    ('assistant',   'Assistant',   'Asystent głosowy z historią działań',                                   1),
    ('terminal',    'Terminal',    'Powłoki, polecenia i procesy urządzenia',                               1),
    ('developer',   'Developer',   'Kod, budowanie, repozytorium Git',                                      1),
    ('diagnostics', 'Diagnostics', 'Stan systemu, logi i analiza błędów',                                   1),
    ('apps',        'Apps',        'Budowa produktu od architektury po wdrożenie',                          1),
    ('agents',      'Agents',      'Agenci: tożsamość, model, uprawnienia, konektory',                      1)
ON CONFLICT(kod) DO NOTHING;

-- ── Macierz widoczności ────────────────────────────────────────────────────────
-- Kolejność odpowiada kolejności pozycji w nawigacji bocznej. Automations nie ma
-- okna modułowego w żadnym środowisku, a MultitaskingAI nie udostępnia modułów —
-- oba są nieobecne w macierzy.
WITH macierz(srodowisko_kod, modul_kod, kolejnosc) AS (
    VALUES
        ('talkin',     'studio',      0),
        ('talkin',     'workspace',   1),
        ('talkin',     'browser',     2),
        ('talkin',     'research',    3),
        ('talkin',     'library',     4),
        ('talkin',     'translate',   5),
        ('talkin',     'roundtable',  6),
        ('talkin',     'assistant',   7),
        ('talkin',     'agents',      8),
        ('workspace',  'studio',      0),
        ('workspace',  'workspace',   1),
        ('workspace',  'browser',     2),
        ('workspace',  'research',    3),
        ('workspace',  'library',     4),
        ('workspace',  'roundtable',  5),
        ('workspace',  'design',      6),
        ('workspace',  'apps',        7),
        ('workspace',  'agents',      8),
        ('codestudio', 'workspace',   0),
        ('codestudio', 'roundtable',  1),
        ('codestudio', 'design',      2),
        ('codestudio', 'terminal',    3),
        ('codestudio', 'developer',   4),
        ('codestudio', 'diagnostics', 5),
        ('codestudio', 'apps',        6),
        ('codestudio', 'agents',      7)
)
INSERT INTO srodowisko_modul (srodowisko_id, modul_id, kolejnosc, widoczny)
SELECT s.id, m.id, macierz.kolejnosc, 1
  FROM macierz
  JOIN srodowisko s ON s.kod = macierz.srodowisko_kod
  JOIN modul m      ON m.kod = macierz.modul_kod
 WHERE true
ON CONFLICT(srodowisko_id, modul_id) DO NOTHING;

-- ── Chat Window — okno wspólne wszystkim piętnastu modułom ─────────────────────
INSERT INTO okno_operacyjne (kod, modul_id, nazwa, rola, kategoria, kolejnosc)
SELECT m.kod || '.chat-window', m.id, 'Chat Window', 'wiodace', 'komunikacja', 0
  FROM modul m
 WHERE true
ON CONFLICT(kod) DO NOTHING;

-- ── Pozostałe okna operacyjne, modułowe i globalne.
--    Pusty kod modułu oznacza okno globalne. ────────────────────────────────────
WITH katalog(kod, modul_kod, nazwa, rola, kategoria, kolejnosc) AS (
    VALUES
        ('assistant.voice-console',          'assistant',   'Voice Console',                       'wiodace',    'komunikacja',  1),
        ('multitaskingai.executor-chat-1',   '',            'Executor Chat (Executor 1)',          'wiodace',    'komunikacja',  2),
        ('multitaskingai.executor-chat-2',   '',            'Executor Chat (Executor 2)',          'wiodace',    'komunikacja',  3),
        ('multitaskingai.coordinator-chat',  '',            'Coordinator Chat',                    'zarzadca',   'komunikacja',  4),

        ('studio.studio-editor',             'studio',      'Studio Editor',                       'wiodace',    'edycja',       1),
        ('developer.code-editor',            'developer',   'Code Editor',                         'wiodace',    'edycja',       1),
        ('design.design-board',              'design',      'Design Board',                        'wiodace',    'edycja',       1),
        ('translate.source-panel',           'translate',   'Source Panel',                        'wiodace',    'edycja',       1),
        ('translate.translation-panels',     'translate',   'Translation Panels',                  'wiodace',    'edycja',       2),
        ('apps.frontend-workspace',          'apps',        'Frontend Workspace',                  'wiodace',    'edycja',       1),
        ('apps.backend-workspace',           'apps',        'Backend Workspace',                   'wiodace',    'edycja',       2),

        ('studio.diff-grep-panel',           'studio',      'Diff/Grep Panel',                     'pomocnicze', 'podglad',      2),
        ('studio.preview-window',            'studio',      'Preview Window',                      'pomocnicze', 'podglad',      3),
        ('design.preview-window',            'design',      'Preview Window',                      'pomocnicze', 'podglad',      2),
        ('library.file-preview',             'library',     'File Preview',                        'pomocnicze', 'podglad',      1),
        ('developer.build-output',           'developer',   'Build Output',                        'monitor',    'podglad',      2),
        ('terminal.output-console',          'terminal',    'Output Console',                      'monitor',    'podglad',      1),
        ('research.export-panel',            'research',    'Export Panel',                        'pomocnicze', 'podglad',      1),

        ('library.library-explorer',         'library',     'Library Explorer',                    'wiodace',    'repozytorium', 2),
        ('workspace.project-library',        'workspace',   'Project Library',                     'zarzadca',   'repozytorium', 1),
        ('studio.session-repository',        'studio',      'Session Repository',                  'zarzadca',   'repozytorium', 4),
        ('library.versioning-panel',         'library',     'Versioning Panel',                    'zarzadca',   'repozytorium', 3),
        ('design.assets-panel',              'design',      'Assets Panel',                        'zarzadca',   'repozytorium', 3),
        ('workspace.context-memory',         'workspace',   'Context Memory',                      'zarzadca',   'repozytorium', 2),

        ('automations.workflow-builder',     'automations', 'Workflow Builder',                    'kreator',    'konstruktor',  1),
        ('agents.agent-builder',             'agents',      'Agent Builder',                       'kreator',    'konstruktor',  1),
        ('apps.product-builder',             'apps',        'Product Builder',                     'wiodace',    'konstruktor',  3),
        ('design.prompt-builder',            'design',      'Prompt Builder',                      'kreator',    'konstruktor',  4),
        ('apps.architecture-designer',       'apps',        'Architecture Designer',               'kreator',    'konstruktor',  4),
        ('research.report-builder',          'research',    'Report Builder',                      'kreator',    'konstruktor',  2),

        ('automations.scheduler',            'automations', 'Scheduler',                           'zarzadca',   'kolejka',      2),
        ('automations.queue-manager',        'automations', 'Queue Manager',                       'zarzadca',   'kolejka',      3),
        ('automations.orchestrator',         'automations', 'Orchestrator',                        'kreator',    'kolejka',      4),

        ('automations.execution-monitor',    'automations', 'Execution Monitor',                   'monitor',    'monitor',      5),
        ('terminal.process-monitor',         'terminal',    'Process Monitor',                     'monitor',    'monitor',      2),
        ('assistant.actions-monitor',        'assistant',   'Actions Monitor',                     'monitor',    'monitor',      2),
        ('assistant.activity-feed',          'assistant',   'Activity Feed',                       'monitor',    'monitor',      3),
        ('diagnostics.diagnostics-center',   'diagnostics', 'Diagnostics Center',                  'wiodace',    'monitor',      1),
        ('multitaskingai.results-analyzer',  '',            'Results Analyzer',                    'monitor',    'monitor',      5),
        ('diagnostics.logs-viewer',          'diagnostics', 'Logs Viewer',                         'monitor',    'monitor',      2),
        ('roundtable.debate-panel',          'roundtable',  'Debate Panel',                        'monitor',    'monitor',      1),

        ('studio.tools-panel',               'studio',      'Tools Panel',                         'pomocnicze', 'narzedzia',    5),
        ('translate.glossary-manager',       'translate',   'Glossary Manager',                    'zarzadca',   'narzedzia',    3),
        ('library.tags-collections',         'library',     'Tags & Collections',                  'zarzadca',   'narzedzia',    4),
        ('agents.permissions-center',        'agents',      'Permissions Center',                  'zarzadca',   'narzedzia',    2),
        ('agents.skills-manager',            'agents',      'Skills Manager',                      'zarzadca',   'narzedzia',    3),
        ('roundtable.moderator-panel',       'roundtable',  'Moderator Panel',                     'zarzadca',   'narzedzia',    2),
        ('developer.git-panel',              'developer',   'Git Panel',                           'zarzadca',   'narzedzia',    3),
        ('workspace.agent-manager',          'workspace',   'Agent Manager',                       'zarzadca',   'narzedzia',    3),
        ('apps.deployment-panel',            'apps',        'Deployment Panel',                    'zarzadca',   'narzedzia',    5),
        ('terminal.terminal-tabs',           'terminal',    'Terminal Tabs',                       'wiodace',    'narzedzia',    3),
        ('developer.project-tree',           'developer',   'Project Tree',                        'pomocnicze', 'narzedzia',    4),

        ('browser.sources-panel',            'browser',     'Sources Panel',                       'pomocnicze', 'zrodla',       1),
        ('research.sources-manager',         'research',    'Sources Manager',                     'zarzadca',   'zrodla',       3),
        ('research.findings-panel',          'research',    'Findings Panel',                      'pomocnicze', 'zrodla',       4),
        ('browser.notes-panel',              'browser',     'Notes Panel',                         'pomocnicze', 'zrodla',       2),

        ('workspace.instructions-panel',     'workspace',   'Instructions Panel',                  'pomocnicze', 'konfiguracja', 4),
        ('agents.model-configuration',       'agents',      'Model Configuration',                 'pomocnicze', 'konfiguracja', 4),
        ('agents.connectors-manager',        'agents',      'Connectors Manager',                  'zarzadca',   'konfiguracja', 5),
        ('workspace.project-dashboard',      'workspace',   'Project Dashboard',                   'wiodace',    'konfiguracja', 5),
        ('globalne.punkty-izolacji',         '',            'Okno konfiguracji punktów izolacji',  'zarzadca',   'konfiguracja', 6),
        ('globalne.konfiguracja',            '',            'Okno Konfiguracji',                   'zarzadca',   'konfiguracja', 7),
        ('globalne.ustawienia',              '',            'Okno Ustawień',                       'zarzadca',   'konfiguracja', 8),
        ('globalne.strona-glowna',           '',            'Strona główna',                       'wiodace',    'konfiguracja', 9)
)
INSERT INTO okno_operacyjne (kod, modul_id, nazwa, rola, kategoria, kolejnosc)
SELECT katalog.kod, m.id, katalog.nazwa, katalog.rola, katalog.kategoria, katalog.kolejnosc
  FROM katalog
  LEFT JOIN modul m ON m.kod = katalog.modul_kod
 WHERE katalog.modul_kod = '' OR m.id IS NOT NULL
ON CONFLICT(kod) DO NOTHING;
