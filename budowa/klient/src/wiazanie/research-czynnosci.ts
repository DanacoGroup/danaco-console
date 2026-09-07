// Wspólne rusztowanie czynności okna badania. Nasłuch stoi na korzeniu karty,
// bo wykazy paneli są wymieniane w całości i wiersz nie przeżywa odświeżenia.
import type { ErrorInfo } from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { oglos } from './ogloszenie.ts';

export const NAGLOWEK_BADAN = 'Badania';

/** Co czynność musi znać: kanał, korzeń karty, bieżące okno i odświeżenie. */
export interface KontekstBadania {
  kanal: Kanal;
  korzen: Element;
  idOkna: () => string;
  odswiez: () => Promise<void>;
}

/** Przycisk czynności dokładany do wiersza wykazu. */
export interface CzynnoscWiersza {
  kod: string;
  etykieta: string;
}

const CZAS_UZBROJENIA = 6000;
const UZBROJONE = new Set<string>();

export function nasluchCzynnosci(
  korzen: Element,
  atrybut: string,
  wykonaj: (kod: string, przycisk: HTMLElement) => Promise<void>,
): Odsubskrybuj {
  const sterowanie = new AbortController();
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przycisk = cel.closest<HTMLElement>(`[${atrybut}]`);
    if (przycisk === null) return;
    zdarzenie.stopPropagation();
    void wykonaj(przycisk.getAttribute(atrybut) ?? '', przycisk);
  }, { signal: sterowanie.signal, capture: true });
  return () => {
    sterowanie.abort();
  };
}

export function wskazanie(przycisk: HTMLElement, atrybut: string): string {
  return przycisk.closest<HTMLElement>(`[${atrybut}]`)?.getAttribute(atrybut) ?? '';
}

export function zapytaj(pytanie: string, domyslna = ''): string {
  return (globalThis.prompt?.(pytanie, domyslna) ?? '').trim();
}

/* Okno nie prowadzi listy wyboru, a identyfikatory rdzenia nic Operatorowi nie
   mówią, więc pozycję wskazuje się numerem z wypisanego spisu. */
export function wybierz(pytanie: string, etykiety: readonly string[]): number {
  if (etykiety.length === 0) return -1;
  const spis = etykiety.map((etykieta, numer) => `${String(numer + 1)}. ${etykieta}`).join('\n');
  const numer = Number.parseInt(zapytaj(`${pytanie}\n${spis}`, '1'), 10);
  if (Number.isNaN(numer) || numer < 1 || numer > etykiety.length) return -1;
  return numer - 1;
}

/* Czynność nieodwracalna wymaga dwóch naciśnięć: pierwsze uzbraja i ostrzega,
   drugie wykonuje. Uzbrojenie wygasa samo, żeby powrót po czasie nie usuwał. */
export function potwierdzone(klucz: string, ostrzezenie: string): boolean {
  if (UZBROJONE.has(klucz)) {
    UZBROJONE.delete(klucz);
    return true;
  }
  UZBROJONE.add(klucz);
  globalThis.setTimeout(() => UZBROJONE.delete(klucz), CZAS_UZBROJENIA);
  oglos(NAGLOWEK_BADAN, ostrzezenie, 'ostrzezenie');
  return false;
}

export function odmowa(blad: ErrorInfo | undefined, zdanie: string): void {
  oglos(NAGLOWEK_BADAN, blad?.message ?? zdanie, 'ostrzezenie');
}

export function powiedz(tresc: string): void {
  oglos(NAGLOWEK_BADAN, tresc);
}

export function brakOkna(idOkna: string): boolean {
  if (idOkna !== '') return false;
  odmowa(undefined, 'Okno badania jeszcze nie powstało, więc czynność nie ma do czego się odnieść.');
  return true;
}

/* Kod bez gałęzi wykonawczej meldowałby sukces bez działania, więc każde
   rozgałęzienie czynności kończy się tym nazwaniem. */
export function nieznanaCzynnosc(kod: string): void {
  odmowa(undefined, `Czynność „${kod}" nie ma pokrycia w kontrakcie, więc nic nie zostało wykonane.`);
}

export function wierszCzynnosci(
  cialo: Element,
  tytul: string,
  podpis: string,
  atrybutKlucza: string,
  klucz: string,
  atrybutCzynnosci: string,
  czynnosci: readonly CzynnoscWiersza[],
): HTMLElement {
  const dokument = cialo.ownerDocument;
  const wiersz = dokument.createElement('div');
  wiersz.className = 'dn-wykaz-modulu-poz';
  wiersz.setAttribute(atrybutKlucza, klucz);
  const nazwa = dokument.createElement('span');
  nazwa.textContent = tytul;
  wiersz.append(nazwa);
  if (podpis !== '') {
    const meta = dokument.createElement('span');
    meta.className = 'dn-meta';
    meta.textContent = podpis;
    wiersz.append(meta);
  }
  for (const czynnosc of czynnosci) {
    const przycisk = dokument.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-btn dn-btn--zarys dn-btn--sm';
    przycisk.setAttribute(atrybutCzynnosci, czynnosc.kod);
    przycisk.textContent = czynnosc.etykieta;
    wiersz.append(przycisk);
  }
  return wiersz;
}

/* Okno nie ma gdzie zostawić pliku — powłoka nie prowadzi wskazania katalogu —
   więc treść wydania idzie do schowka wraz z nazwą podaną przez rdzeń. */
export async function doSchowka(tresc: string, nazwa: string): Promise<void> {
  try {
    await navigator.clipboard.writeText(tresc);
    powiedz(`Treść jest w schowku; nazwa z rdzenia to „${nazwa}".`);
  } catch {
    odmowa(undefined, 'Przeglądarka nie dała dostępu do schowka, więc treść nie została przeniesiona.');
  }
}
