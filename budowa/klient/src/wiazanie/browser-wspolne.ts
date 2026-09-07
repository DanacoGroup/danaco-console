// Rusztowanie czynności przeglądania wspólne wszystkim panelom modułu: pasy
// działań nad wykazami, jedno zdanie odmowy rdzenia, przeniesienie treści do
// schowka i dwukrotne naciśnięcie dla kroku, którego nie da się cofnąć.

import { oglos } from './ogloszenie.ts';

export const NAGLOWEK = 'Przeglądanie';

const CZAS_POTWIERDZENIA = 5000;

export function zalozPas(etykieta: string): HTMLElement {
  const pas = document.createElement('div');
  pas.className = 'dn-pas-dzialan';
  pas.setAttribute('aria-label', etykieta);
  return pas;
}

/* Pas stoi NAD ciałem panelu, bo ciało jest wymieniane przy każdym odświeżeniu
   wykazu i wszystko w nim postawione znikłoby razem z wykazem. */
export function pasNadCialem(cialo: Element | null, etykieta: string): HTMLElement | null {
  if (cialo === null || cialo.parentElement === null) return null;
  const pas = zalozPas(etykieta);
  cialo.parentElement.insertBefore(pas, cialo);
  return pas;
}

export function polePasa(
  pas: HTMLElement | null,
  podpis: string,
  wartosc = '',
): HTMLInputElement | null {
  if (pas === null) return null;
  const pole = document.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole-kontrolka';
  pole.placeholder = podpis;
  pole.value = wartosc;
  pole.setAttribute('aria-label', podpis);
  pas.append(pole);
  return pole;
}

export function wyborPasa(
  pas: HTMLElement | null,
  podpis: string,
  opcje: readonly (readonly [string, string])[],
): HTMLSelectElement | null {
  if (pas === null) return null;
  const wybor = document.createElement('select');
  wybor.className = 'dn-pole-kontrolka';
  wybor.setAttribute('aria-label', podpis);
  for (const [wartosc, etykieta] of opcje) {
    const pozycja = document.createElement('option');
    pozycja.value = wartosc;
    pozycja.textContent = etykieta;
    wybor.append(pozycja);
  }
  pas.append(wybor);
  return wybor;
}

export function dolozPrzycisk(
  gniazdo: Element | null,
  etykieta: string,
  cecha: string,
  wartosc = '',
): HTMLButtonElement | null {
  if (gniazdo === null) return null;
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--zarys dn-btn--sm';
  przycisk.dataset[cecha] = wartosc;
  przycisk.textContent = etykieta;
  gniazdo.append(przycisk);
  return przycisk;
}

export function odmowa(blad: { message?: string } | undefined, zdanie: string): void {
  oglos(NAGLOWEK, blad?.message ?? zdanie, 'ostrzezenie');
}

export function powiedz(zdanie: string): void {
  oglos(NAGLOWEK, zdanie);
}

/* Krok nieodwracalny idzie dwoma naciśnięciami: pierwsze zamienia napis na
   ostrzeżenie, drugie wykonuje. Bez drugiego przycisk wraca sam. */
export function potwierdzone(przycisk: HTMLElement, ostrzezenie: string): boolean {
  if (przycisk.dataset.potwierdzenie === 'czeka') {
    przycisk.textContent = przycisk.dataset.napisPrzed ?? przycisk.textContent;
    delete przycisk.dataset.potwierdzenie;
    delete przycisk.dataset.napisPrzed;
    return true;
  }
  przycisk.dataset.napisPrzed = przycisk.textContent ?? '';
  przycisk.dataset.potwierdzenie = 'czeka';
  przycisk.textContent = 'Potwierdź';
  oglos(NAGLOWEK, ostrzezenie, 'ostrzezenie');
  globalThis.setTimeout(() => {
    if (przycisk.dataset.potwierdzenie !== 'czeka') return;
    przycisk.textContent = przycisk.dataset.napisPrzed ?? '';
    delete przycisk.dataset.potwierdzenie;
    delete przycisk.dataset.napisPrzed;
  }, CZAS_POTWIERDZENIA);
  return false;
}

/* Okno nie ma gdzie zostawić pliku, więc treść oddana przez rdzeń idzie do
   schowka razem z nazwą, pod którą rdzeń ją prowadzi. */
export async function doSchowka(tresc: string, nazwa: string): Promise<void> {
  if (tresc === '') {
    oglos(NAGLOWEK, `Rdzeń oddał „${nazwa}” bez treści — nie ma czego przenieść.`, 'ostrzezenie');
    return;
  }
  try {
    await navigator.clipboard.writeText(tresc);
    oglos(NAGLOWEK, `Treść „${nazwa}” jest w schowku (${tresc.length} znaków).`);
  } catch {
    oglos(NAGLOWEK,
      'Przeglądarka nie dała dostępu do schowka, więc treść nie została przeniesiona.',
      'ostrzezenie');
  }
}

export function adresZPola(korzen: Element): string {
  const pole = korzen.querySelector<HTMLInputElement>('[data-adres-przegladania]');
  return (pole?.value ?? '').trim();
}

export function adresWymagany(korzen: Element): string {
  const adres = adresZPola(korzen);
  if (adres === '') {
    oglos(NAGLOWEK, 'Pasek adresu jest pusty — czynność nie ma na czym stanąć.', 'ostrzezenie');
  }
  return adres;
}

export function wypiszWynik(gniazdo: Element | null, naglowek: string, wiersze: string[]): void {
  if (gniazdo === null) return;
  gniazdo.replaceChildren();
  const tytul = document.createElement('div');
  tytul.className = 'dn-etyk-mono';
  tytul.textContent = naglowek;
  gniazdo.append(tytul);
  if (wiersze.length === 0) {
    const pusty = document.createElement('p');
    pusty.className = 'dn-pusty';
    pusty.textContent = 'Rdzeń nie oddał ani jednej pozycji.';
    gniazdo.append(pusty);
    return;
  }
  for (const wiersz of wiersze) {
    const wezel = document.createElement('div');
    wezel.className = 'dn-wykaz-modulu-poz';
    wezel.textContent = wiersz;
    gniazdo.append(wezel);
  }
}

export function naBase64(tekst: string): string {
  const bajty = new TextEncoder().encode(tekst);
  let zapis = '';
  for (const bajt of bajty) zapis += String.fromCharCode(bajt);
  return btoa(zapis);
}
