import { ResearchFindingStatus } from '../../../../shared/contract';
import { przycisk } from '../../modele/kontrolki-formularza';
import { KODY_OKIEN } from './kody-okien';
import type { StanBadania } from './stan-badania';

/**
 * Panel postępu badania liczy wyłącznie to, co rdzeń potwierdził i oddał w tym połączeniu,
 * i jawnie nazywa metryki, których dziś policzyć nie sposób.
 */
export interface PanelPostepu {
  element: HTMLElement;
  odswiez(): void;
}

/** Jeden licznik panelu opisuje nazwę pomiaru, wyliczoną wartość tekstową oraz kod okna, do którego kliknięcie prowadzi. */
interface Licznik {
  nazwa: string;
  wartosc: string;
  okno: string;
}

export function utworzPanelPostepu(
  stan: StanBadania,
  przejdz: (kodOkna: string) => void,
): PanelPostepu {
  const wykaz = document.createElement('div');
  wykaz.className = 'mr-postep__wykaz';

  const zastrzezenie = document.createElement('p');
  zastrzezenie.className = 'dn-pole-opis mr-postep__zastrzezenie';
  zastrzezenie.textContent =
    'Liczniki pokazują to, co rdzeń potwierdził w tym połączeniu. Pokrycia pytań badawczych ' +
    'i podziału źródeł na przeczytane panel jeszcze nie liczy: obie zdolności mają komendy ' +
    'w kontrakcie, ale rdzeń nie ma dla nich uchwytu, a stanu lektury nie niesie samo źródło, ' +
    'które rdzeń oddaje. Procent policzony z danych, których nie ma, byłby zmyśleniem.';

  const naglowek = document.createElement('p');
  naglowek.className = 'dn-pole-etykieta';
  naglowek.textContent = 'Postęp badania';

  const element = document.createElement('div');
  element.className = 'dn-karta mr-postep';
  element.setAttribute('aria-label', 'Postęp badania');
  element.append(naglowek, wykaz, zastrzezenie);

  return {
    element,

    odswiez() {
      wykaz.replaceChildren(...liczniki(stan).map((licznik) => pozycja(licznik, przejdz)));
    },
  };
}

/** Komplet liczników panelu, policzony ze stanu badania przekazanego przez pamięć modułu, gotowy do wyrenderowania jako pozycje. */
function liczniki(stan: StanBadania): Licznik[] {
  const zrodla = stan.zrodla();
  const ustalenia = stan.ustalenia();
  const raport = stan.raport();
  const sekcje = raport?.sections ?? [];

  // Ustalenia w dokumencie liczone po zbiorze `findingIds` sekcji, nie po ich liczbie.
  const wRaporcie = new Set<string>();
  for (const sekcja of sekcje) {
    for (const ustalenie of sekcja.findingIds ?? []) wRaporcie.add(ustalenie);
  }
  const objete = ustalenia.filter((ustalenie) => wRaporcie.has(ustalenie.id)).length;

  return [
    { nazwa: 'Źródła', wartosc: String(zrodla.length), okno: KODY_OKIEN.zrodla },
    { nazwa: 'Ustalenia', wartosc: String(ustalenia.length), okno: KODY_OKIEN.ustalenia },
    {
      nazwa: 'Ustalenia rozstrzygnięte',
      wartosc: `${String(rozstrzygniete(stan))} z ${String(ustalenia.length)}`,
      okno: KODY_OKIEN.ustalenia,
    },
    { nazwa: 'Sekcje raportu', wartosc: String(sekcje.length), okno: KODY_OKIEN.raport },
    {
      nazwa: 'Pokrycie ustaleń w raporcie',
      wartosc: pokrycie(objete, ustalenia.length),
      okno: KODY_OKIEN.raport,
    },
  ];
}

/** Liczba ustaleń w stanie rozstrzygniętym, jedynym rozróżnieniu statusu ustalenia, jakie kontrakt dziś przewiduje. */
function rozstrzygniete(stan: StanBadania): number {
  return stan
    .ustalenia()
    .filter((ustalenie) => ustalenie.status === ResearchFindingStatus.Resolved).length;
}

/**
 * Pokrycie ustaleń wyrażone ułamkiem, nie samym procentem.
 *
 * Procent bez mianownika kłamie przy małych liczbach: „100%" z jednego ustalenia
 * na jedno wygląda jak raport kompletny. Ułamek mówi, o ilu rzeczach mowa.
 */
function pokrycie(objete: number, wszystkie: number): string {
  if (wszystkie === 0) return 'brak ustaleń do objęcia';
  const procent = Math.round((objete / wszystkie) * 100);
  return `${String(objete)} z ${String(wszystkie)} (${String(procent)}%)`;
}

/** Jeden licznik przedstawiony jako kontrolka, której kliknięcie przechodzi do okna wypełniającego dany licznik danymi. */
function pozycja(licznik: Licznik, przejdz: (kodOkna: string) => void): HTMLElement {
  const element = document.createElement('div');
  element.className = 'mr-postep__pozycja';

  const nazwa = document.createElement('span');
  nazwa.className = 'mr-postep__nazwa';
  nazwa.textContent = licznik.nazwa;

  const kontrolka = przycisk(licznik.wartosc, 'dn-btn dn-btn--sm dn-btn--duch mr-postep__wartosc');
  kontrolka.setAttribute('aria-label', `${licznik.nazwa}: ${licznik.wartosc} — przejdź do okna`);
  kontrolka.addEventListener('click', () => przejdz(licznik.okno));

  element.append(nazwa, kontrolka);
  return element;
}
