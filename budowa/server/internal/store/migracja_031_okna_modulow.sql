-- Migracja 031 — macierz „okno operacyjne ↔ moduł”.
--
-- Uzupełnia `migracja_030_rejestr_okien_operacyjnych.sql`: definicja okna leży
-- w `okno_operacyjne`, a jego wystąpienie w module leży tutaj. Klucz główny na
-- parze (okno, moduł) nie pozwala przypiąć tego samego okna do modułu dwa razy,
-- więc powielanie wiersza definicji dla okien wspólnych nie ma jak wrócić.
--
-- KOLEJNOŚĆ NALEŻY DO PRZYPIĘCIA, NIE DO OKNA. Preview Window stoi w Studio na
-- pozycji trzeciej, a w Design na drugiej. Pozycja zero należy do okna rozmowy,
-- bo to ono jest punktem wejścia modułu; po nim idzie okno wiodące modułu,
-- a dalej okna w porządku pracy operatora.

CREATE TABLE okno_operacyjne_modul (
    okno_operacyjne_id INTEGER NOT NULL REFERENCES okno_operacyjne(id) ON DELETE CASCADE,
    modul_id           INTEGER NOT NULL REFERENCES modul(id) ON DELETE CASCADE,
    kolejnosc          INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (okno_operacyjne_id, modul_id)
);
CREATE INDEX idx_okno_operacyjne_modul_kolejnosc ON okno_operacyjne_modul(modul_id, kolejnosc);

-- ── Chat Window — jedna definicja w piętnastu modułach ─────────────────────────
-- Zapisane regułą, nie listą: „okno rozmowy jest w każdym module” to zdanie
-- inwentarza, a lista piętnastu kodów byłaby jego przepisaniem, które rozjedzie
-- się przy szesnastym module.
INSERT INTO okno_operacyjne_modul (okno_operacyjne_id, modul_id, kolejnosc)
SELECT o.id, m.id, 0
  FROM okno_operacyjne o
 CROSS JOIN modul m
 WHERE o.kod = 'chat-window'
ON CONFLICT(okno_operacyjne_id, modul_id) DO NOTHING;

-- ── Pozostałe przypięcia ─────────────────────────────────────────────────────
WITH macierz(okno_kod, modul_kod, kolejnosc) AS (
    VALUES
        ('studio-editor',         'studio',      1),
        ('diff-grep-panel',       'studio',      2),
        ('preview-window',        'studio',      3),
        ('session-repository',    'studio',      4),
        ('tools-panel',           'studio',      5),

        ('project-dashboard',     'workspace',   1),
        ('project-library',       'workspace',   2),
        ('context-memory',        'workspace',   3),
        ('agent-manager',         'workspace',   4),
        ('instructions-panel',    'workspace',   5),

        ('workflow-builder',      'automations', 1),
        ('scheduler',             'automations', 2),
        ('queue-manager',         'automations', 3),
        ('orchestrator',          'automations', 4),
        ('execution-monitor',     'automations', 5),

        ('browser-window',        'browser',     1),
        ('sources-panel',         'browser',     2),
        ('notes-panel',           'browser',     3),

        ('research-workspace',    'research',    1),
        ('sources-manager',       'research',    2),
        ('findings-panel',        'research',    3),
        ('report-builder',        'research',    4),
        ('export-panel',          'research',    5),

        ('library-explorer',      'library',     1),
        ('file-preview',          'library',     2),
        ('versioning-panel',      'library',     3),
        ('tags-collections',      'library',     4),

        ('source-panel',          'translate',   1),
        ('translation-panels',    'translate',   2),
        ('glossary-manager',      'translate',   3),

        ('model-panels',          'roundtable',  1),
        ('debate-panel',          'roundtable',  2),
        ('moderator-panel',       'roundtable',  3),
        ('consensus-panel',       'roundtable',  4),

        ('design-board',          'design',      1),
        ('preview-window',        'design',      2),
        ('assets-panel',          'design',      3),
        ('prompt-builder',        'design',      4),

        ('voice-console',         'assistant',   1),
        ('actions-monitor',       'assistant',   2),
        ('activity-feed',         'assistant',   3),

        ('terminal-tabs',         'terminal',    1),
        ('output-console',        'terminal',    2),
        ('process-monitor',       'terminal',    3),

        ('code-editor',           'developer',   1),
        ('project-tree',          'developer',   2),
        ('build-output',          'developer',   3),
        ('git-panel',             'developer',   4),

        ('diagnostics-center',    'diagnostics', 1),
        ('logs-viewer',           'diagnostics', 2),
        ('errors-panel',          'diagnostics', 3),
        ('recommendations-panel', 'diagnostics', 4),

        ('product-builder',       'apps',        1),
        ('architecture-designer', 'apps',        2),
        ('frontend-workspace',    'apps',        3),
        ('backend-workspace',     'apps',        4),
        ('deployment-panel',      'apps',        5),

        ('agent-builder',         'agents',      1),
        ('model-configuration',   'agents',      2),
        ('permissions-center',    'agents',      3),
        ('skills-manager',        'agents',      4),
        ('connectors-manager',    'agents',      5)
)
INSERT INTO okno_operacyjne_modul (okno_operacyjne_id, modul_id, kolejnosc)
SELECT o.id, m.id, macierz.kolejnosc
  FROM macierz
  JOIN okno_operacyjne o ON o.kod = macierz.okno_kod
  JOIN modul m           ON m.kod = macierz.modul_kod
 WHERE true
ON CONFLICT(okno_operacyjne_id, modul_id) DO NOTHING;
