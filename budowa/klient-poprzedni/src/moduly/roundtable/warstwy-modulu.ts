import {
  bilansFunkcji,
  funkcjeOkna,
  opisStanuFunkcji,
  type BilansFunkcji,
  type FunkcjaKatalogu,
  type OknoModulu,
  type WarstwaWidocznosci,
} from './katalog-funkcji';

/** Warstwy widoczności modułu Roundtable — rysowanie tego, czego okno nie pokazuje w spoczynku: nazwa warstwy, jak czyta ją Operator. */
const NAZWY_WARSTW: Record<WarstwaWidocznosci, string> = {
  1: 'warstwa 1 — widoczne bez interakcji',
  2: 'warstwa 2 — przycisk, przełącznik albo znacznik kontekstowy',
  3: 'warstwa 3 — zestaw akcji okna',
  4: 'warstwa 4 — nastawa poza stanem spoczynku okna',
};

/** Rodzaj warstwy zwiniętej — wartość atrybutu danych dla arkusza, jedna z trzech nośników rozwijanych tego pliku. */
export type RodzajWarstwy = 'zestaw-akcji' | 'wykaz-funkcji' | 'diagnostyka';

/**
 * Element rozwijany warstwy — wspólna obudowa zestawu akcji, wykazu funkcji i diagnostyki,
 * z wnętrzem oddawanym wywołującemu.
 */
export function utworzWarstweTresci(
  podpis: string,
  waga: RodzajWarstwy,
): { element: HTMLDetailsElement; wnetrze: HTMLElement } {
  return utworzRozwijane(podpis, waga);
}

function utworzRozwijane(podpis: string, waga: RodzajWarstwy): {
  element: HTMLDetailsElement;
  wnetrze: HTMLElement;
} {
  const element = document.createElement('details');
  element.className = 'dr-warstwa';
  element.dataset['warstwa'] = waga;

  const naglowek = document.createElement('summary');
  naglowek.className = 'dr-warstwa__podpis';
  naglowek.textContent = podpis;

  const wnetrze = document.createElement('div');
  wnetrze.className = 'dr-warstwa__wnetrze';

  element.append(naglowek, wnetrze);
  return { element, wnetrze };
}

/**
 * Zestaw akcji okna — warstwa trzecia, gotowe przyciski od okna, bo warstwa nie zna ani jednej
 * komendy i niczego nie wywołuje.
 */
export function utworzZestawAkcji(podpis: string, akcje: readonly HTMLElement[]): HTMLElement {
  const { element, wnetrze } = utworzRozwijane(podpis, 'zestaw-akcji');
  wnetrze.classList.add('dr-warstwa__akcje');
  wnetrze.append(...akcje);
  return element;
}

/**
 * Wykaz narzędzi okna wraz ze stanem ich wykonania — warstwa czwarta, jedyna droga poznania
 * granicy działania modułu.
 */
export function utworzWykazFunkcji(okno: OknoModulu): HTMLElement {
  const funkcje = funkcjeOkna(okno);
  const bilans = bilansFunkcji(funkcje);
  const { element, wnetrze } = utworzRozwijane(
    `Narzędzia tego okna wg opracowania modułu — ${bilans.wszystkich}, w tym bez zbudowanej obsługi: ${bilans.bezObslugi}`,
    'wykaz-funkcji',
  );

  const zdanie = document.createElement('p');
  zdanie.className = 'dr-warstwa__zdanie';
  zdanie.textContent = zdanieBilansu(bilans);
  wnetrze.append(zdanie);

  for (const warstwa of [1, 2, 3, 4] as const) {
    const pozycje = funkcje.filter((funkcja) => funkcja.warstwa === warstwa);
    if (pozycje.length === 0) continue;
    wnetrze.append(sekcjaWarstwy(warstwa, pozycje));
  }
  return element;
}

/**
 * Zdanie bilansu — liczby, nie ocena dojrzałości: obsługa niezbudowana i brak pokrycia
 * w kontrakcie to dwie różne rzeczy.
 */
function zdanieBilansu(bilans: BilansFunkcji): string {
  const czlony = [
    `wykonane w całości: ${bilans.wykonanych}`,
    `wykonane częściowo: ${bilans.czesciowych}`,
    `bez zbudowanej obsługi: ${bilans.bezObslugi}`,
  ];
  if (bilans.bezPokrycia > 0) czlony.push(`bez pokrycia w kontrakcie: ${bilans.bezPokrycia}`);
  if (bilans.pozaModulem > 0) czlony.push(`poza tym modułem: ${bilans.pozaModulem}`);
  return (
    `Opracowanie modułu wymienia dla tego okna ${bilans.wszystkich} narzędzi — ${czlony.join(' · ')}. ` +
    'Stan mówi, co okno dziś robi, a nie co da się zrobić: komenda obecna w kontrakcie bez zbudowanej obsługi liczy się jako niewykonana.'
  );
}

/** Jedna warstwa wykazu wraz z pozycjami, podpisana nazwą warstwy i liczbą pozycji, które do niej należą. */
function sekcjaWarstwy(
  warstwa: WarstwaWidocznosci,
  pozycje: readonly FunkcjaKatalogu[],
): HTMLElement {
  const sekcja = document.createElement('div');
  sekcja.className = 'dr-warstwa__sekcja';

  const podpis = document.createElement('p');
  podpis.className = 'dr-warstwa__nazwa';
  podpis.textContent = `${NAZWY_WARSTW[warstwa]} (${pozycje.length})`;

  const wykaz = document.createElement('ul');
  wykaz.className = 'dr-funkcje';
  wykaz.setAttribute('aria-label', NAZWY_WARSTW[warstwa]);
  for (const funkcja of pozycje) wykaz.append(pozycjaFunkcji(funkcja));

  sekcja.append(podpis, wykaz);
  return sekcja;
}

/** Pozycja wykazu: nazwa narzędzia, stan słowem i zdanie o granicy wykonania, wraz z atrybutem stanu dla arkusza. */
function pozycjaFunkcji(funkcja: FunkcjaKatalogu): HTMLElement {
  const pozycja = document.createElement('li');
  pozycja.className = 'dr-funkcja';
  // Stan idzie atrybutem i słowem naraz — sama barwa nie niesie stanu nikomu, kto barw nie rozróżnia.
  pozycja.dataset['stanFunkcji'] = funkcja.stan;

  const nazwa = document.createElement('strong');
  nazwa.className = 'dr-funkcja__nazwa';
  nazwa.textContent = funkcja.nazwa;

  const stan = document.createElement('span');
  stan.className = 'dr-funkcja__stan';
  stan.textContent = opisStanuFunkcji(funkcja.stan);

  const zdanie = document.createElement('span');
  zdanie.className = 'dr-funkcja__zdanie';
  zdanie.textContent = funkcja.zdanie;

  pozycja.append(nazwa, stan, zdanie);
  return pozycja;
}
