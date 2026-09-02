// Wspólne części okna Agents: stanowisko eksperta i pomocniki węzłów prototypu.

import type { Agent } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';

export interface StanowiskoEkspertow {
  kanal: Kanal;
  korzen: Element;
  idOkna: string;
  przy: AddEventListenerOptions;
  ekspert: () => Agent | null;
  postaw: (agent: Agent) => void;
  odswiez: () => void;
}

export const BEZ_EKSPERTA = 'Wskaż eksperta w bibliotece — rdzeń nie ma czego oddać.';
export const BEZ_POKRYCIA = 'Rdzeń nie oddaje danych dla tego panelu.';

export function tekst(wezel: Element | null, wartosc: string): void {
  if (wezel !== null) wezel.textContent = wartosc;
}

export function pole(korzen: ParentNode, wybor: string): HTMLInputElement | null {
  const wezel = korzen.querySelector(wybor);
  return wezel instanceof HTMLInputElement ? wezel : null;
}

export function obszar(korzen: ParentNode, wybor: string): HTMLTextAreaElement | null {
  const wezel = korzen.querySelector(wybor);
  return wezel instanceof HTMLTextAreaElement ? wezel : null;
}

// Wiersz wzorcowy prototypu niesie treść przykładową: zostaje wzorem, schodzi z dokumentu.
export function wzorem(rodzic: Element | null, wybor: string): HTMLElement | null {
  if (rodzic === null) return null;
  const wiersze = [...rodzic.querySelectorAll<HTMLElement>(wybor)];
  const pierwszy = wiersze[0];
  const wzor = pierwszy === undefined ? null : (pierwszy.cloneNode(true) as HTMLElement);
  for (const wiersz of wiersze) wiersz.remove();
  return wzor;
}

export function pustka(rodzic: Element | null, zdanie: string): void {
  if (rodzic === null) return;
  rodzic.querySelector('[data-pustka]')?.remove();
  const wiersz = document.createElement('p');
  wiersz.className = 'dn-tekst-ciagly dn-tekst-ciagly--drobny';
  wiersz.dataset.pustka = 'tak';
  wiersz.textContent = zdanie;
  rodzic.appendChild(wiersz);
}

export function zdejmijPustke(rodzic: Element | null): void {
  rodzic?.querySelector('[data-pustka]')?.remove();
}

// Pole formularza rozpoznaje się po podpisie: prototyp nie znakuje ich niczym innym.
export function poleZPodpisem(korzen: ParentNode, podpis: string): HTMLElement | null {
  for (const kandydat of korzen.querySelectorAll<HTMLElement>('.ab-pole')) {
    const etykieta = kandydat.querySelector('label');
    if ((etykieta?.textContent ?? '').trim().startsWith(podpis)) return kandydat;
  }
  return null;
}

export function przelacz(wezel: Element | null, wlaczony: boolean): void {
  wezel?.setAttribute('aria-checked', String(wlaczony));
}

export function czyWlaczony(wezel: Element | null): boolean {
  return wezel?.getAttribute('aria-checked') === 'true';
}

export function godzinaZapisu(znacznik: number): string {
  return new Date(znacznik).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
}

// Grupa wyborów prototypu ma stałą liczbę opcji; wykaz z rdzenia bywa dłuższy,
// więc pierwsza opcja zostaje wzorem, a reszta powstaje z jej klonów.
export function wypelnijGrupe(
  grupa: Element | null,
  pozycje: readonly { klucz: string; podpis: string; zaznaczona: boolean }[],
): void {
  if (grupa === null) return;
  const wzorMozeMoze = wzorem(grupa, '.ab-opcja');
  if (wzorMozeMoze === null) return;
  const wzorMoze = wzorMozeMoze;
  const wzor = wzorMoze;
  for (const pozycja of pozycje) {
    const opcja = wzor.cloneNode(true) as HTMLElement;
    opcja.dataset.klucz = pozycja.klucz;
    const zaznaczenie = opcja.querySelector('input');
    if (zaznaczenie instanceof HTMLInputElement) zaznaczenie.checked = pozycja.zaznaczona;
    for (const dziecko of [...opcja.childNodes]) {
      if (dziecko.nodeType === Node.TEXT_NODE) dziecko.remove();
    }
    opcja.appendChild(document.createTextNode(pozycja.podpis));
    grupa.appendChild(opcja);
  }
}

export function zaznaczoneKlucze(grupa: Element | null): string[] {
  if (grupa === null) return [];
  return [...grupa.querySelectorAll<HTMLElement>('.ab-opcja')]
    .filter((opcja) => opcja.querySelector('input')?.checked === true)
    .map((opcja) => opcja.dataset.klucz ?? '')
    .filter((klucz) => klucz !== '');
}

export function nazwaEksperta(agent: Agent): string {
  return agent.displayName ?? agent.name;
}
