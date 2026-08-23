import { DiffHunkKind, type StudioDiffHunk } from '../../../../shared/contract';

/**
 * Statystyka i filtr różnicy — rachunek na tym, co rdzeń już oddał.
 *
 * Opracowanie żąda od Diff/Grep Panelu paska statystyki („+48 słów · −12 słów ·
 * 3 zmiany") oraz filtrowania różnic po rodzaju. Kontrakt nie niesie komendy
 * liczącej statystykę i nie musi: `studio.diff.compare` oddaje fragmenty wraz
 * z treścią przed i po, więc liczby wynikają z odpowiedzi, którą panel już ma.
 * Drugie wywołanie po te same liczby byłoby pytaniem o coś, co leży na stole.
 *
 * Plik nie zna DOM: wejściem są fragmenty kontraktu, wyjściem liczby i wykaz
 * przefiltrowany. Dzięki temu rachunek sprawdza się bez stawiania okna.
 */

/** Liczby opisujące różnicę dwóch stron porównania. */
export interface StatystykaRoznicy {
  fragmenty: number;
  dodane: number;
  usuniete: number;
  zmienione: number;
  /** Słowa dopisane — z fragmentów dodanych i ze strony „po" fragmentów zmienionych. */
  slowaDodane: number;
  /** Słowa zdjęte — z fragmentów usuniętych i ze strony „przed" fragmentów zmienionych. */
  slowaUsuniete: number;
  znakiDodane: number;
  znakiUsuniete: number;
}

/** Wartość filtru wykazu różnic. */
export type FiltrRoznicy = 'wszystkie' | DiffHunkKind;

/** Liczba słów w napisie; pusty napis nie ma ani jednego. */
function slowa(tresc: string): number {
  return tresc.split(/\s+/u).filter((slowo) => slowo !== '').length;
}

/**
 * Liczy statystykę z fragmentów oddanych przez rdzeń.
 *
 * Fragment zmieniony liczy się do obu stron naraz — jego treść „przed" jest
 * ubytkiem, a „po" przyrostem. Liczenie go tylko raz zaniżałoby obie liczby
 * i pasek mówiłby o mniejszej zmianie, niż zaszła.
 *
 * Fragment kontekstowy nie liczy się do żadnej strony: rdzeń podaje go dla
 * czytelności, a nie jako zmianę.
 */
export function policzRoznice(fragmenty: readonly StudioDiffHunk[]): StatystykaRoznicy {
  const wynik: StatystykaRoznicy = {
    fragmenty: fragmenty.length,
    dodane: 0,
    usuniete: 0,
    zmienione: 0,
    slowaDodane: 0,
    slowaUsuniete: 0,
    znakiDodane: 0,
    znakiUsuniete: 0,
  };

  for (const fragment of fragmenty) {
    const przed = fragment.before ?? '';
    const po = fragment.after ?? '';
    if (fragment.kind === DiffHunkKind.Added) wynik.dodane += 1;
    if (fragment.kind === DiffHunkKind.Removed) wynik.usuniete += 1;
    if (fragment.kind === DiffHunkKind.Changed) wynik.zmienione += 1;
    if (fragment.kind === DiffHunkKind.Context) continue;

    if (fragment.kind !== DiffHunkKind.Removed) {
      wynik.slowaDodane += slowa(po);
      wynik.znakiDodane += po.length;
    }
    if (fragment.kind !== DiffHunkKind.Added) {
      wynik.slowaUsuniete += slowa(przed);
      wynik.znakiUsuniete += przed.length;
    }
  }
  return wynik;
}

/** Zawęża wykaz fragmentów do jednego rodzaju; `wszystkie` niczego nie odsiewa. */
export function przefiltrujRoznice(
  fragmenty: readonly StudioDiffHunk[],
  filtr: FiltrRoznicy,
): readonly StudioDiffHunk[] {
  if (filtr === 'wszystkie') return fragmenty;
  return fragmenty.filter((fragment) => fragment.kind === filtr);
}

/** Zdanie paska statystyki — jedno miejsce składania tych liczb w napis. */
export function opiszStatystyke(statystyka: StatystykaRoznicy): string {
  if (statystyka.fragmenty === 0) return 'Rdzeń nie oddał ani jednego fragmentu różnicy.';
  return (
    `Fragmentów ${statystyka.fragmenty}: dodanych ${statystyka.dodane}, ` +
    `usuniętych ${statystyka.usuniete}, zmienionych ${statystyka.zmienione}. ` +
    `Słów +${statystyka.slowaDodane} / −${statystyka.slowaUsuniete} · ` +
    `znaków +${statystyka.znakiDodane} / −${statystyka.znakiUsuniete}.`
  );
}
