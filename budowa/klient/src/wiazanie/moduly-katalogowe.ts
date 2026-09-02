// Moduły prowadzone katalogiem operacji: okno z prototypu, a w nim generyczny
// katalog komend rodziny. Dla rodzin bez własnego, formularzowego okna to jest
// pełna droga od interfejsu do rdzenia — Operator wybiera operację i ją uruchamia.

import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { zwiazKatalogModulu } from './katalog-modulu.ts';

const WIAZANIA = new Map<string, Odsubskrybuj>();

function korzenKarty(wskazanie: Element | string): Element | null {
  if (typeof wskazanie !== 'string') return wskazanie;
  return document.querySelector(`.cd-tresc--modul[data-karta="${wskazanie}"]`);
}

function zwiaz(
  kanal: Kanal,
  idOkna: string,
  wskazanie: Element | string,
  rodzina: string,
  nazwa: string,
): boolean {
  const korzen = korzenKarty(wskazanie);
  if (korzen === null) return false;
  const idKarty = korzen.getAttribute('data-karta') ?? '';
  if (idKarty === '') return false;
  if (WIAZANIA.has(idKarty)) return false;
  const odlacz = zwiazKatalogModulu(kanal, idOkna, korzen, rodzina, nazwa);
  if (odlacz !== null) WIAZANIA.set(idKarty, odlacz);
  return odlacz !== null;
}

function zwolnij(idKarty: string): void {
  const odlacz = WIAZANIA.get(idKarty);
  if (odlacz === undefined) return;
  odlacz();
  WIAZANIA.delete(idKarty);
}

export function zwiazBadania(k: Kanal, o: string, w: Element | string): boolean {
  return zwiaz(k, o, w, 'research', 'Badania');
}
export function zwiazTlumaczenie(k: Kanal, o: string, w: Element | string): boolean {
  return zwiaz(k, o, w, 'translate', 'Tłumaczenie');
}
export function zwiazDevelopera(k: Kanal, o: string, w: Element | string): boolean {
  return zwiaz(k, o, w, 'developer', 'Developer');
}
export function zwiazDebate(k: Kanal, o: string, w: Element | string): boolean {
  return zwiaz(k, o, w, 'roundtable', 'Roundtable');
}
export function zwiazAplikacje(k: Kanal, o: string, w: Element | string): boolean {
  return zwiaz(k, o, w, 'apps', 'Aplikacje');
}
export function zwiazAutomatyzacje(k: Kanal, o: string, w: Element | string): boolean {
  return zwiaz(k, o, w, 'automation', 'Automatyzacje');
}
export function zwiazAsystenta(k: Kanal, o: string, w: Element | string): boolean {
  return zwiaz(k, o, w, 'assistant', 'Asystent');
}
export function zwiazDiagnostyke(k: Kanal, o: string, w: Element | string): boolean {
  return zwiaz(k, o, w, 'diagnostics', 'Diagnostyka');
}

export const zwolnijKatalogModulu = zwolnij;
