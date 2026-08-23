import { ResearchFindingStatus } from '../../../../shared/contract';
import { przycisk } from '../../modele/kontrolki-formularza';
import { KODY_OKIEN } from './kody-okien';
import type { StanBadania } from './stan-badania';

/**
 * Panel postępu badania — liczbowy przegląd stanu, policzony z tego, co rdzeń
 * potwierdził.
 *
 * Każdy licznik powstaje z pamięci modułu, a pamięć niesie wyłącznie odpowiedzi
 * rdzenia z tego połączenia. Żadna z tych liczb nie jest oszacowaniem: „źródeł
 * 14" znaczy czternaście źródeł, które rdzeń skatalogował i oddał.
 *
 * **Czego panel NIE liczy i dlaczego.** Opracowanie modułu (rozdz. 3.3)
 * przewiduje w tym miejscu pokrycie pytań badawczych źródłami oraz podział
 * źródeł na przeczytane i nieprzeczytane. Obie rzeczy mają już w kontrakcie
 * swoje komendy, a mimo to policzyć ich dziś nie sposób — z dwóch różnych
 * powodów, więc warto je rozróżnić:
 *
 * - **pokrycie pytań** czeka wyłącznie na uchwyt w rdzeniu; kształt odpowiedzi
 *   jest rozstrzygnięty i licznik dojdzie razem z uchwytem;
 * - **stan lektury** ma pole w ŻĄDANIU zapisu źródła, ale `ResearchSource`,
 *   które rdzeń oddaje, tego pola nie niesie. Wartość da się więc wysłać,
 *   a nie da się jej odczytać z powrotem — i tego uchwyt sam nie naprawi,
 *   dopóki byt źródła nie dostanie pola.
 *
 * Procent policzony z danych, których nie ma, byłby metryką zmyśloną, więc
 * panel go nie pokazuje i mówi, czego brakuje.
 *
 * Pokrycie ustaleń jest natomiast policzalne i policzone: sekcje raportu niosą
 * pole `findingIds`, więc „ile z zebranych ustaleń weszło do dokumentu" wychodzi
 * z porównania dwóch zbiorów, które rdzeń oddał.
 */
export interface PanelPostepu {
  element: HTMLElement;
  odswiez(): void;
}

/** Jeden licznik panelu: co liczy, ile wyszło i do którego okna prowadzi. */
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

/** Komplet liczników policzonych ze stanu badania. */
function liczniki(stan: StanBadania): Licznik[] {
  const zrodla = stan.zrodla();
  const ustalenia = stan.ustalenia();
  const raport = stan.raport();
  const sekcje = raport?.sections ?? [];

  // Ustalenia, które weszły do dokumentu — liczone po zbiorze `findingIds`
  // sekcji, a nie po ich liczbie: to samo ustalenie może zasilać dwie sekcje.
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

/** Ustalenia w stanie rozstrzygniętym — jedyne rozróżnienie stanu, jakie kontrakt ma. */
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

/** Jeden licznik jako kontrolka prowadząca do okna, które go wypełnia. */
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
