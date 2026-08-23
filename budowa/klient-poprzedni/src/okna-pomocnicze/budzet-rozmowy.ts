import { LICZBA_MIN } from '../okna-rownolegle/identyfikatory';
import { liczbaOkienRozmowy, type ProfilModulu } from '../okno-komunikacji/profil-modulu';
import { profilModulu } from '../okno-komunikacji/rejestr-profilow';

/**
 * Ile okien rozmowy ma moduł — odczyt profilu, nie własny mechanizm.
 *
 * Budżet stoi przy oknach pomocniczych, bo pas pomocniczych podaje obie połowy
 * naraz: ile okien rozmowy moduł prowadzi i gdzie leży jego ciężar. Pas
 * milczący o rozmowie mówiłby tylko połowę.
 *
 * Liczbę daje funkcja `liczbaOkienRozmowy`, nie pole `granicaOkien`: pole
 * niesie granicę, funkcja — prawo do otwarcia. Dla modułu bez rozmowy pole
 * mówi 1, a funkcja 0, więc odczyt samego pola obiecałby Operatorowi czat
 * w module, który czatu nie prowadzi.
 *
 * Liczba początkowa pochodzi z `LICZBA_MIN` (`okna-rownolegle/identyfikatory.ts`),
 * czyli stamtąd, skąd bierze ją silnik okien równoległych. Własna jedynka
 * rozjechałaby się z silnikiem przy pierwszej jego zmianie.
 *
 * Granicy zapisanej w profilu nie da się tu odróżnić od stałej
 * `GRANICA_NIEPODANA` — obie są w czasie działania czwórką. Zdanie podaje więc
 * samą liczbę i nic o jej pochodzeniu nie orzeka.
 */
export interface BudzetRozmowy {
  /** Ile okien rozmowy stoi na scenie na starcie. */
  naStart: number;
  /** Do ilu Operator może rozbudować; 0 znaczy „ten moduł rozmowy nie prowadzi". */
  granica: number;
  /** Zdanie dla Operatora — pełne, nie skrót. */
  zdanie: string;
}

/** Budżet okien rozmowy modułu; moduł spoza rejestru dostaje profil wspólny. */
export function budzetRozmowy(kodModulu: string): BudzetRozmowy {
  const profil = profilModulu(kodModulu);
  const granica = liczbaOkienRozmowy(profil);
  const naStart = granica === 0 ? 0 : Math.min(LICZBA_MIN, granica);
  return { naStart, granica, zdanie: zdanieBudzetu(profil, naStart, granica) };
}

function zdanieBudzetu(profil: ProfilModulu, naStart: number, granica: number): string {
  if (granica === 0) {
    // Zdanie nazywa moduł i mówi, co stoi zamiast rozmowy — pusty pasek
    // czytałby się jak usterka.
    return (
      `Moduł ${profil.nazwa} nie prowadzi rozmowy — ${profil.przeznaczenie.toLowerCase()}. ` +
      'Okno komunikacji nie otworzy się tu wcale i nie jest to awaria.'
    );
  }
  const podstawa =
    naStart === granica
      ? `Okno rozmowy: ${granica}, bez rozbudowy.`
      : `Okno rozmowy: ${naStart} na start, rozbudowa do ${granica}.`;
  if (profil.postacRozmowy === 'dymek-glosowy') {
    return `${podstawa} Postać: pływający awatar z oknem dymkowym, nie zwykłe okno czatu.`;
  }
  return `${podstawa} Ciężar modułu leży w oknach pomocniczych poniżej.`;
}
