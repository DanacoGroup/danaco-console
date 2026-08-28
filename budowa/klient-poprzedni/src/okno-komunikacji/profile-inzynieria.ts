import { narzedzie, profil, type ProfilModulu } from './profil-modulu';

/**
 * Profile pięciu modułów pracy inżynierskiej — Terminal, Developer, Diagnostics, Apps i Agents — gdzie Terminal zostaje bez własnej postaci rozmowy, będąc oknem pomocniczym wewnątrz pozostałych modułów.
 */
export const PROFILE_INZYNIERII: readonly ProfilModulu[] = [
  profil(
    'terminal',
    'Terminal',
    'Powłoki, polecenia i procesy urządzenia',
    ['Terminal Tabs', 'Output Console', 'Process Monitor'],
    [
      narzedzie(
        'terminal.zbuduj',
        'Zbuduj polecenie',
        'Polecenie powłoki wyprowadzone z opisu naturalnego (terminal.command.exec)',
        'Zbuduj polecenie powłoki z opisu i wstaw je do karty terminala bez uruchamiania: ',
      ),
      narzedzie(
        'terminal.uruchom',
        'Uruchom teraz',
        'Wykonanie polecenia przez model z wpisem „inicjator: AI"',
        'Uruchom w bieżącej karcie terminala polecenie i pokaż jego wyjście: ',
      ),
      narzedzie(
        'terminal.wyjscie',
        'Wyjaśnij wyjście',
        'Rozbiór fragmentu wyjścia z Output Console (terminal.output.stream)',
        'Wyjaśnij ostatnie wyjście z Output Console i wskaż przyczynę kodu wyjścia.',
      ),
      narzedzie(
        'terminal.procesy',
        'Procesy',
        'Rejestr procesów rdzenia z rozróżnieniem inicjatora (terminal.process.list)',
        'Pokaż procesy tej sesji wraz z inicjatorem, zużyciem i kodem wyjścia.',
      ),
    ],
  ),
  profil(
    'developer',
    'Developer',
    'Kod, budowanie, repozytorium Git',
    ['Code Editor', 'Project Tree', 'Git Panel', 'Build Output'],
    [
      narzedzie(
        'developer.generuj',
        'Generuj',
        'Nowy fragment kodu z opisu albo kontekstu repozytorium',
        'Wygeneruj kod w kontekście otwartego pliku i pokaż różnicę przed wstawieniem: ',
      ),
      narzedzie(
        'developer.refaktor',
        'Refaktoryzuj',
        'Zmiana struktury bez zmiany zachowania',
        'Zrefaktoryzuj wskazany fragment bez zmiany zachowania i uzasadnij każdą zmianę.',
      ),
      narzedzie(
        'developer.test',
        'Napisz test',
        'Sprawdzian pokrywający wskazane zachowanie (developer.build.run)',
        'Napisz sprawdzian pokrywający wskazane zachowanie i uruchom budowanie.',
      ),
      narzedzie(
        'developer.git',
        'Git',
        'Działanie na repozytorium powiązanym z sesją (developer.git.action)',
        'Wykonaj działanie Git na repozytorium sesji i podaj stan przed i po: ',
      ),
    ],
    {
      // Nacisk modułu leży na oknach pomocniczych, więc cztery okna kontekstu są obowiązkowe, nie rozmowa.
      postacRozmowy: 'okno',
      granicaOkien: 4,
      pamiecSesyjna: true,
      oknaObowiazkowe: ['Code Editor', 'Project Tree', 'Git Panel', 'Build Output'],
    },
  ),
  profil(
    'diagnostics',
    'Diagnostics',
    'Stan systemu, logi i analiza błędów',
    ['Diagnostics Center', 'Logs Viewer', 'Errors Panel', 'Recommendations Panel'],
    [
      narzedzie(
        'diagnostics.analiza',
        'Uruchom analizę',
        'Analiza zagregowanego stanu systemu (diagnostics.analyze.run)',
        'Uruchom analizę diagnostyczną dla wskazanego zakresu czasu i podsumuj wynik.',
      ),
      narzedzie(
        'diagnostics.logi',
        'Przeszukaj logi',
        'Zapytanie do strumienia logów na żywo (diagnostics.log.query)',
        'Przeszukaj logi wzorcem i pokaż wpisy pasujące wraz z ich źródłem: ',
      ),
      narzedzie(
        'diagnostics.bledy',
        'Błędy',
        'Grupowanie błędów po odcisku i przegląd stanu (diagnostics.error.list)',
        'Pokaż błędy zgrupowane po odcisku, od najczęstszego, wraz z ich kontekstem.',
      ),
      narzedzie(
        'diagnostics.poprawka',
        'Poproś o poprawkę',
        'Rekomendacja poprawki z Recommendations Panel (diagnostics.recommendation.list)',
        'Zaproponuj poprawkę dla wskazanego błędu i podaj skutek jej zastosowania.',
      ),
    ],
    {
      // Rozmowa jest częścią analizy, ale rolę główną pełnią logi — stąd trzy okna obowiązkowe.
      postacRozmowy: 'okno',
      granicaOkien: 4,
      pamiecSesyjna: true,
      oknaObowiazkowe: ['Logs Viewer', 'Errors Panel', 'Recommendations Panel'],
    },
  ),
  profil(
    'apps',
    'Apps',
    'Budowa produktu od architektury po wdrożenie',
    [
      'Product Builder',
      'Architecture Designer',
      'Frontend Workspace',
      'Backend Workspace',
      'Deployment Panel',
    ],
    [
      narzedzie(
        'apps.architektura',
        'Architektura',
        'Definicja komponentów i zależności rozwiązania (apps.architecture.define)',
        'Zdefiniuj komponenty rozwiązania i ich zależności w Architecture Designer: ',
      ),
      narzedzie(
        'apps.warstwa',
        'Warstwa produktu',
        'Praca w Frontend albo Backend Workspace (apps.workspace.update)',
        'Wykonaj pracę w warstwie produktu (frontend albo backend) — zakres: ',
      ),
      narzedzie(
        'apps.wdrozenie',
        'Wdrożenie',
        'Uruchomienie wdrożenia wraz ze strategią i rollbackiem (apps.deployment.run)',
        'Uruchom wdrożenie na wskazane środowisko i podaj strategię oraz drogę powrotu.',
      ),
      narzedzie(
        'apps.wykonawca',
        'Wybór wykonawcy',
        'Praca własna albo podział Executor 1 / Executor 2',
        'Podziel to zadanie między wykonawców i wskaż, kto bierze frontend, a kto backend.',
      ),
    ],
    {
      // Pełna sesja z multitaskingiem do czterech okien: frontend, backend, architektura albo wdrożenie.
      postacRozmowy: 'okno',
      granicaOkien: 4,
      pamiecSesyjna: true,
    },
  ),
  profil(
    'agents',
    'Agents',
    'Agenci: tożsamość, model, uprawnienia, konektory',
    [
      'Agent Builder',
      'Model Configuration',
      'Skills Manager',
      'Connectors Manager',
      'Permissions Center',
    ],
    [
      narzedzie(
        'agents.ekspert',
        'Nowy ekspert',
        'Utworzenie agenta wraz z tożsamością i instrukcjami (agent.create)',
        'Zbuduj eksperta: nazwa, przeznaczenie i instrukcje systemowe — ',
      ),
      narzedzie(
        'agents.model',
        'Model i kanał',
        'Model bazowy oraz kanał wywołania: Code CLI, Agent SDK, API (agent.model.set)',
        'Ustaw modelowi tego eksperta model bazowy i kanał wywołania: ',
      ),
      narzedzie(
        'agents.skille',
        'Skille i konektory',
        'Przypisanie skilla i konektora MCP (agent.skill.add, agent.connector.add)',
        'Przypisz ekspertowi skille i konektory potrzebne do jego pracy: ',
      ),
      narzedzie(
        'agents.uprawnienia',
        'Uprawnienia',
        'Cztery grupy zakresu; stan wyjściowy to pełny dostęp',
        'Pokaż uprawnienia eksperta w czterech grupach zakresu i wskaż zmiany do wykonania.',
      ),
    ],
    {
      // Jedno okno pełni funkcję testową, bez trwałej pamięci — zmiana agenta rozpoczyna nowy kontekst.
      postacRozmowy: 'okno',
      granicaOkien: 1,
      pamiecSesyjna: false,
    },
  ),
];
