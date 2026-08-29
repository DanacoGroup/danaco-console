/**
 * Panel Plan — katalog treści własny tego panelu. Ten sam wzorzec co
 * `moduly/studio/tresci.ts`: jedyne miejsce z tekstem widocznym dla Operatora
 * w tym pliku. Stoi osobno od katalogu wspólnego modułu, żeby siedmiu
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

  /* Panel nie filtruje po stanie zadania, więc pustka znaczy dokładnie tyle,
     że rozkład nie ma ani jednego zadania — nie że żadne nie pasuje do wyboru. */
  pusty: {
    tytul: 'Brak zadań',
    opis: 'Rozkład nie zawiera zadań.',
  },

  odmowa: {
    okno: 'Nie udało się otworzyć okna modułu',
    sesja: 'Nie udało się założyć sesji',
    /* Jeden nagłówek na obie drogi: odmowa wywołania i odpowiedź bez rozkładu
       znaczą dla Operatora to samo — zlecenie nie zostało rozłożone. */
    rozklad: 'Nie udało się rozłożyć zlecenia na zadania',
    uruchomienie: 'Nie udało się uruchomić pętli',
    zatrzymanie: 'Nie udało się zatrzymać pętli',
  },
} as const;
