/**
 * Panel Plan — katalog treści własny tego panelu. Ten sam wzorzec co
 * `moduly/studio/tresci.ts`: jedyne miejsce z tekstem widocznym dla Operatora
 * w tym panelu. Stoi osobno od katalogu wspólnego modułu, żeby siedmiu
 * wykonawców paneli okna roboczego nie pisało jednego pliku równolegle.
 */

export const planTresci = {
  panel: {
    tytul: 'Plan',
    /* Znacznik belki niesie stosunek zadań domkniętych do wszystkich, więc
       czytnik ekranu potrzebuje zdania, którego sam zapis „2 / 4" nie daje. */
    etykietaZnacznika: 'Zadania wykonane',
  },

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

  zadanie: {
    etykietaDzialan: 'Czynności wskazanego zadania',
    wykonane: 'Oznacz wykonane',
    pomin: 'Pomiń',
    ponow: 'Ponów',
    przestawianie: 'Przestawianie…',
    wykonawca: 'Wykonawca:',
    powod: 'Powód:',
    wynik: 'Wynik:',
  },

  /* Dokument bywa rozłożony więcej niż raz; wybór stoi tylko wtedy, gdy
     rozkładów jest kilka — przy jednym nie ma czego wybierać. */
  rozklady: {
    etykieta: 'Rozkłady dokumentu',
    bezZlecenia: 'Rozkład bez zlecenia',
  },

  nowe: {
    etykieta: 'Nowe zlecenie',
    zastepcza: 'Napisz, co ma się stać z dokumentem…',
    rozloz: 'Rozłóż na zadania',
    rozkladanie: 'Rozkładanie…',
    brakDokumentu: 'W tym oknie nie ma jeszcze dokumentu, więc nie ma czego rozłożyć na zadania.',
  },

  lancuchy: {
    etykieta: 'Łańcuchy operacji',
    uruchom: 'Uruchom',
    uruchamianie: 'Uruchamianie…',
    brak: 'Żaden łańcuch operacji nie jest zapisany.',
    /** Poprzedza liczbę kroków przyjętych do wykonania zaraz po uruchomieniu. */
    przyjeto: 'Kroki przekazane do wykonania:',
    jednostkaKroki: { jedna: 'krok', kilka: 'kroki', wiele: 'kroków' },
  },

  ladowanie: 'Wczytywanie rozkładu…',

  /* Rozkładu nie ma i odmowa to dwa różne stany: pierwszy znaczy, że zlecenia
     jeszcze nikt nie rozłożył, drugi — że odpowiedź w ogóle nie przyszła. */
  brakRozkladu: {
    tytul: 'Brak rozkładu',
    opis: 'Żadne zlecenie tego dokumentu nie zostało jeszcze rozłożone na zadania.',
  },

  /* Panel nie zawęża wykazu do stanu zadania, więc pustka znaczy dokładnie
     tyle, że rozkład nie ma ani jednego zadania — nie że żadne nie pasuje. */
  pusty: {
    tytul: 'Brak zadań',
    opis: 'Rozkład nie zawiera zadań.',
  },

  odmowa: {
    okno: 'Nie udało się otworzyć okna modułu',
    sesja: 'Nie udało się założyć sesji',
    /* Jeden nagłówek na obie drogi: odmowa wywołania i odpowiedź bez rozkładu
       znaczą dla Operatora to samo — zlecenie nie zostało rozłożone. */
    rozklad: 'Nie udało się wczytać rozkładu zlecenia',
    rozlozenie: 'Nie udało się rozłożyć zlecenia na zadania',
    uruchomienie: 'Nie udało się uruchomić pętli',
    zatrzymanie: 'Nie udało się zatrzymać pętli',
    zadanie: 'Nie udało się przestawić zadania',
    wykazLancuchow: 'Nie udało się wczytać łańcuchów operacji',
    lancuch: 'Nie udało się uruchomić łańcucha operacji',
  },
} as const;

/** Postęp łańcucha zdaniem, a nie samą parą liczb, którą trzeba sobie objaśnić. */
export function opisPostepuLancucha(krok: number, krokow: number): string {
  return `Krok ${krok} z ${krokow}`;
}
