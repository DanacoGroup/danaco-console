import { GRANICA_NIEPODANA, narzedzie, profil, type ProfilModulu } from './profil-modulu';

/**
 * Profile pięciu modułów pracy z treścią: Studio, Workspace, Automations,
 * Browser, Research.
 *
 * Podział na trzy pliki idzie wzdłuż rodzajów pracy, a nie wzdłuż długości
 * pliku.
 *
 * Granica okien jest ustalona tylko dla Automations (dwa okna: wykonawca
 * i koordynator). Studio, Workspace, Browser i Research stoją na
 * `GRANICA_NIEPODANA`.
 */
export const PROFILE_TRESCI: readonly ProfilModulu[] = [
  profil(
    'studio',
    'Studio',
    'Praca z dokumentem: edycja, operacje kontekstowe, wersjonowanie sesji',
    ['Praca z dokumentem', 'Tools Panel', 'Session Repository', 'Ingest/OCR Panel'],
    [
      narzedzie(
        'studio.operacja',
        'Operacja kontekstowa',
        'Operacja Tools Panel na zaznaczeniu albo całym dokumencie (studio.contextual.op)',
        'Wykonaj operację kontekstową na dokumencie otwartym w oknie pracy z dokumentem — ' +
          'zakres i rodzaj operacji: ',
      ),
      narzedzie(
        'studio.wstaw',
        'Wstaw do dokumentu',
        'Wstawienie ostatniej odpowiedzi do dokumentu otwartego w oknie pracy z dokumentem',
        'Wstaw ostatnią odpowiedź do dokumentu otwartego w oknie pracy z dokumentem i zapisz wersję.',
      ),
      narzedzie(
        'studio.roznica',
        'Porównaj wersje',
        'Porównanie wersji przed i po zmianie, nakładane na treść (studio.diff.compare)',
        'Porównaj bieżącą wersję dokumentu z poprzednią i pokaż różnicę na treści dokumentu.',
      ),
      narzedzie(
        'studio.repozytorium',
        'Repozytorium sesji',
        'Historia wersji dokumentu narastająca w toku sesji (studio.repository.list)',
        'Pokaż historię wersji dokumentu z Session Repository i wskaż wersję do przywrócenia.',
      ),
    ],
    {
      // Moduł składa się z okna komunikacji i czterech okien operacyjnych.
      // Powierzchnia tekstowa jest JEDNA: okno pracy z dokumentem niesie treść,
      // podgląd wydania i różnicę jako tryby jednego widoku — Studio Editor,
      // kanwa tekstowa, Preview Window i Diff/Grep Panel zeszły się w nie.
      // Pozostałe trzy okna powierzchni tekstowej nie mają.
      postacRozmowy: 'okno',
      granicaOkien: GRANICA_NIEPODANA,
      pamiecSesyjna: true,
      oknaObowiazkowe: [
        'Praca z dokumentem',
        'Tools Panel',
        'Session Repository',
        'Ingest/OCR Panel',
      ],
    },
  ),
  profil(
    'workspace',
    'Workspace',
    'Projekt: zadania, zasoby, instrukcje systemowe, pamięć kontekstu',
    ['Project Dashboard', 'Instructions Panel', 'Context Memory', 'Project Library', 'Agent Manager'],
    [
      narzedzie(
        'workspace.wzmianka',
        'Wzmianka @',
        'Odwołanie do pliku, wpisu pamięci albo agenta projektu',
        'Weź do kontekstu zasób projektu (plik, wpis pamięci albo agenta): @',
      ),
      narzedzie(
        'workspace.wykonawca',
        'Wybór wykonawcy',
        'Przekazanie zadania agentowi przypisanemu w Agent Manager (workspace.agent.assign)',
        'Przekaż to zadanie agentowi przypisanemu do projektu i podaj, którego wybierasz.',
      ),
      narzedzie(
        'workspace.pamiec',
        'Przypnij do pamięci',
        'Zapis ustalenia w Context Memory projektu (workspace.context.set)',
        'Przypnij ostatnie ustalenie do Context Memory projektu i podaj jego zasięg.',
      ),
      narzedzie(
        'workspace.instrukcje',
        'Instrukcje efektywne',
        'Podgląd instrukcji warstwowych projektu i sesji (workspace.instructions.set)',
        'Pokaż instrukcje efektywne projektu i wskaż warstwę, z której pochodzi każda z nich.',
      ),
    ],
    {
      // Rozmowa jest centralnym elementem projektu; granicy okien nie ustalono.
      postacRozmowy: 'okno',
      granicaOkien: GRANICA_NIEPODANA,
      pamiecSesyjna: true,
    },
  ),
  profil(
    'automations',
    'Automations',
    'Workflow, harmonogram, kolejki i orkiestracja',
    ['Workflow Builder', 'Scheduler', 'Queue Manager', 'Orchestrator', 'Execution Monitor'],
    [
      narzedzie(
        'automations.krok',
        'Generuj krok',
        'Krok workflow zbudowany z opisu (automation.workflow.save)',
        'Zbuduj krok workflow z opisu i dołóż go do definicji automatyki: ',
      ),
      narzedzie(
        'automations.log',
        'Debugowanie logu',
        'Rozbiór logu z Execution Monitor i propozycja poprawki',
        'Rozbierz log przebiegu z Execution Monitor i wskaż poprawkę kroku, który zawiódł.',
      ),
      narzedzie(
        'automations.kolejka',
        'Działanie na kolejce',
        'Jedenaście działań silnika kolejek (automation.queue.action)',
        'Wykonaj działanie na kolejce (enqueue, dequeue, delay, retry, pause, resume, split, ' +
          'merge, route, branch, condition): ',
      ),
      narzedzie(
        'automations.harmonogram',
        'Harmonogram',
        'Cykliczność i wyzwalacze przebiegu (automation.schedule.set)',
        'Ustal harmonogram tej automatyki: częstotliwość, godziny i wyzwalacze.',
      ),
    ],
    {
      // Dwa okna, nie cztery: koordynator zarządza procesem, wykonawca
      // realizuje zadania. To uproszczona odmiana środowiska multitaskingu.
      postacRozmowy: 'okno',
      granicaOkien: 2,
      pamiecSesyjna: true,
    },
  ),
  profil(
    'browser',
    'Browser',
    'Przeglądanie stron z gromadzeniem źródeł i notatek',
    ['Browser Window', 'Sources Panel', 'Notes Panel'],
    [
      narzedzie(
        'browser.nawiguj',
        'Nawiguj',
        'Polecenie nawigacyjne wydane modelowi (browser.navigate)',
        'Otwórz stronę we wspólnym podglądzie Browser Window i opisz, co na niej widzisz: ',
      ),
      narzedzie(
        'browser.wyjasnij',
        'Wyjaśnij stronę',
        'Wyjaśnienie treści widocznej w oknie przeglądarki (browser.snapshot.get)',
        'Wyjaśnij treść widoczną teraz w Browser Window i wskaż fragmenty wymagające źródła.',
      ),
      narzedzie(
        'browser.zrodlo',
        'Dodaj źródło',
        'Dopisanie strony do Sources Panel (browser.source.add)',
        'Dopisz bieżącą stronę do Sources Panel wraz z metadanymi pochodzenia.',
      ),
      narzedzie(
        'browser.notatka',
        'Dodaj jako notatkę',
        'Zapis fragmentu odpowiedzi w Notes Panel (browser.note.add)',
        'Zapisz ostatnią odpowiedź jako notatkę w Notes Panel i powiąż ją ze źródłem.',
      ),
    ],
    {
      // Rozmowa pełna, powiązana z przeglądaną treścią; granicy okien nie
      // ustalono.
      postacRozmowy: 'okno',
      granicaOkien: GRANICA_NIEPODANA,
      pamiecSesyjna: true,
    },
  ),
  profil(
    'research',
    'Research',
    'Badanie: źródła, ustalenia, raport',
    ['Research Workspace', 'Sources Manager', 'Findings Panel', 'Report Builder', 'Export Panel'],
    [
      narzedzie(
        'research.ustalenie',
        'Dodaj jako ustalenie',
        'Zapis wniosku w Findings Panel wraz z powiązaniem źródła (research.finding.add)',
        'Zapisz ostatni wniosek jako ustalenie w Findings Panel i powiąż je ze źródłem.',
      ),
      narzedzie(
        'research.zrodlo',
        'Dodaj źródło',
        'Katalogowanie źródła w Sources Manager (research.source.add)',
        'Dopisz źródło do Sources Manager wraz z typem, pochodzeniem i oceną wiarygodności: ',
      ),
      narzedzie(
        'research.stan',
        'Podsumuj stan badania',
        'Zestawienie źródeł, ustaleń i luk badania',
        'Podsumuj stan badania: źródła, ustalenia, sprzeczności i luki wymagające domknięcia.',
      ),
      narzedzie(
        'research.raport',
        'Zbuduj raport',
        'Kompozycja raportu z ustaleń (research.report.build)',
        'Zbuduj raport z ustaleń Findings Panel i przygotuj streszczenie zarządcze.',
      ),
    ],
    {
      // Rozmowa jest głównym sterem procesu badawczego; granicy okien nie
      // ustalono.
      postacRozmowy: 'okno',
      granicaOkien: GRANICA_NIEPODANA,
      pamiecSesyjna: true,
    },
  ),
];
