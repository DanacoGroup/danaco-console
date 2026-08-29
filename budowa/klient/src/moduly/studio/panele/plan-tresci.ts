/**
 * Panel Plan — katalog treści własny tego panelu. Ten sam wzorzec co
 * `moduly/studio/tresci.ts`: jedyne miejsce z tekstem widocznym dla Operatora
 * w tym plik. Stoi osobno od katalogu wspólnego modułu, żeby siedmiu
 * wykonawców paneli okna roboczego nie pisało jednego pliku równolegle.
 */

export const planTresci = {
  /** Poprzedza zlecenie Operatora własnymi słowami, wzięte wprost z rozkładu. */
  zlecenie: 'Zlecenie:',
  /** Poprzedza słowne miano stanu pętli wykonawczej rozkładu. */
  stanPetli: 'Stan pętli:',

  stanRozkladuSlowem: {
    draft: 'ułożony',
    running: 'w biegu',
    paused: 'wstrzymany',
    done: 'zakończony',
    stopped: 'zatrzymany',
  },

  /** Znak i miano stanu zadania, kluczowane wartością `StudioTaskState`. */
  stanZadania: {
    pending: { glif: '○', etykieta: 'Oczekuje' },
    running: { glif: '▶', etykieta: 'W biegu' },
    done: { glif: '✓', etykieta: 'Wykonane' },
    failed: { glif: '✕', etykieta: 'Zakończone błędem' },
    blocked: { glif: '‖', etykieta: 'Wstrzymane' },
    skipped: { glif: '⊘', etykieta: 'Pominięte' },
  },

  akcje: {
    uruchom: 'Uruchom pętlę',
    zatrzymaj: 'Zatrzymaj pętlę',
    uruchamianie: 'Uruchamianie…',
    zatrzymywanie: 'Zatrzymywanie…',
  },

  ladowanie: 'Wczytywanie rozkładu…',

  pusty: {
    tytul: 'Rdzeń nie podał żadnego zadania',
    opis: 'Rozkład nie niesie dziś ani jednego zadania w tym stanie.',
  },

  odmowa: {
    okno: 'Rdzeń nie założył okna modułu',
    sesja: 'Rdzeń nie założył sesji',
    rozklad: 'Rdzeń odmówił rozkładu zlecenia',
    brakRozkladu: 'Rdzeń nie podał rozkładu',
    uruchomienie: 'Rdzeń odmówił uruchomienia pętli',
    zatrzymanie: 'Rdzeń odmówił zatrzymania pętli',
  },
} as const;
