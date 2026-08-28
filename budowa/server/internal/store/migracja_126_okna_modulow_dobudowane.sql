-- Migracja 126 dopisuje do katalogu okno_operacyjne i przypięć okno_operacyjne_modul definicje siedmiu okien dobudowanych modułom, z rolą i kategorią dobraną wzorem pozycji już obecnych w katalogu.

INSERT INTO okno_operacyjne (kod, nazwa, rola, kategoria, kolejnosc)
VALUES
    ('metadata-archive-panel',   'Metadata & Archive Panel',   'zarzadca',   'repozytorium', 7),
    ('memory-context-manager',   'Memory & Context Manager',   'zarzadca',   'repozytorium', 8),
    ('script-library',           'Script Library',             'zarzadca',   'repozytorium', 9),

    ('automation-studio',        'Automation Studio',          'kreator',    'konstruktor',  7),

    ('task-schedule',            'Task & Schedule',            'zarzadca',   'kolejka',      4),

    ('observability-tools',      'Observability Tools',        'monitor',    'monitor',      9),
    ('capture-monitor-panel',    'Capture & Monitor Panel',    'monitor',    'monitor',     10),
    ('argument-map-analysis',    'Argument Map & Analysis',    'monitor',    'monitor',     11),
    ('voting-evaluation-center', 'Voting & Evaluation Center', 'zarzadca',   'monitor',     12),
    ('qa-review-center',         'QA & Review Center',         'zarzadca',   'monitor',     13),

    ('command-tools-hub',        'Command & Tools Hub',        'zarzadca',   'narzedzia',   12),
    ('translation-memory-panel', 'Translation Memory Panel',   'zarzadca',   'narzedzia',   13),
    ('format-studio',            'Format Studio',              'pomocnicze', 'narzedzia',   14),
    ('session-manager',          'Session Manager',            'zarzadca',   'narzedzia',   15)
ON CONFLICT(kod) DO NOTHING;

-- Wiąże w tabeli okno_operacyjne_modul nowo dodane okna z modułami, w których pracują, ustalając kolejność wyświetlenia w obrębie każdego modułu docelowego.
WITH macierz(okno_kod, modul_kod, kolejnosc) AS (
    VALUES
        ('automation-studio',        'browser',      4),
        ('capture-monitor-panel',    'browser',      5),

        ('argument-map-analysis',    'roundtable',   5),
        ('voting-evaluation-center', 'roundtable',   6),

        ('memory-context-manager',   'assistant',    4),
        ('command-tools-hub',        'assistant',    5),

        ('session-manager',          'terminal',     4),
        ('task-schedule',            'terminal',     5),
        ('script-library',           'terminal',     6),

        ('translation-memory-panel', 'translate',    4),
        ('format-studio',            'translate',    5),
        ('qa-review-center',         'translate',    6),

        ('metadata-archive-panel',   'library',      5),

        ('observability-tools',      'diagnostics',  5)
)
INSERT INTO okno_operacyjne_modul (okno_operacyjne_id, modul_id, kolejnosc)
SELECT o.id, m.id, macierz.kolejnosc
  FROM macierz
  JOIN okno_operacyjne o ON o.kod = macierz.okno_kod
  JOIN modul m           ON m.kod = macierz.modul_kod
 WHERE true
ON CONFLICT(okno_operacyjne_id, modul_id) DO NOTHING;
