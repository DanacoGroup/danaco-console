import { DiffHunkKind, type StudioDiffHunk } from '../../../../shared/contract';

/** Liczby opisujące różnicę dwóch stron porównania: liczba fragmentów dodanych, usuniętych i zmienionych, oraz słowa i znaki dodane i usunięte. */
export interface StatystykaRoznicy {
  fragmenty: number;
  dodane: number;
  usuniete: number;
  zmienione: number;
  /** Słowa dopisane z fragmentów dodanych i strony po dla zmienionych. */
  slowaDodane: number;
  /** Słowa zdjęte z fragmentów usuniętych i strony przed dla zmienionych. */
  slowaUsuniete: number;
  znakiDodane: number;
  znakiUsuniete: number;
}

/** Wartość filtru wykazu różnic w panelu porównania; wskazuje jeden rodzaj fragmentu różnicy albo brak filtrowania i pokazanie wszystkich rodzajów. */
export type FiltrRoznicy = 'wszystkie' | DiffHunkKind;

/** Liczy słowa w napisie rozdzielone dowolną sekwencją białych znaków; napis pusty albo złożony z samych odstępów nie ma ani jednego słowa. */
function slowa(tresc: string): number {
  return tresc.split(/\s+/u).filter((slowo) => slowo !== '').length;
}

/**
 * Liczy statystykę z fragmentów oddanych przez rdzeń; fragment zmieniony liczy
 * się do obu stron naraz, a fragment kontekstowy nie liczy się do żadnej.
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

/** Zawęża wykaz fragmentów różnicy do jednego wskazanego rodzaju; wartość `wszystkie` niczego nie odsiewa i oddaje wykaz w całości. */
export function przefiltrujRoznice(
  fragmenty: readonly StudioDiffHunk[],
  filtr: FiltrRoznicy,
): readonly StudioDiffHunk[] {
  if (filtr === 'wszystkie') return fragmenty;
  return fragmenty.filter((fragment) => fragment.kind === filtr);
}

/** Buduje zdanie paska statystyki różnicy z liczby fragmentów, słów i znaków dodanych oraz usuniętych; jest jedynym miejscem składania tych liczb w napis. */
export function opiszStatystyke(statystyka: StatystykaRoznicy): string {
  if (statystyka.fragmenty === 0) return 'Rdzeń nie oddał ani jednego fragmentu różnicy.';
  return (
    `Fragmentów ${statystyka.fragmenty}: dodanych ${statystyka.dodane}, ` +
    `usuniętych ${statystyka.usuniete}, zmienionych ${statystyka.zmienione}. ` +
    `Słów +${statystyka.slowaDodane} / −${statystyka.slowaUsuniete} · ` +
    `znaków +${statystyka.znakiDodane} / −${statystyka.znakiUsuniete}.`
  );
}
