/** Węzły karty WorkSpace z prototypu i wspólne wypełniacze tej karty. */

import type { Wynik } from '../protokol/kanal.ts';
import { oglos } from './ogloszenie.ts';

export interface WezlyWorkspace {
  korzen: Element;
  pulpit: HTMLElement;
  instrukcje: HTMLElement | null;
  pamiec: HTMLElement | null;
  biblioteka: HTMLElement | null;
  agenci: HTMLElement | null;
  wTle: HTMLElement | null;
  plan: HTMLElement | null;
  artefakty: HTMLElement | null;
  terminal: HTMLElement | null;
}

export function zbierzWezly(korzen: Element): WezlyWorkspace | null {
  const pulpitMoze = korzen.querySelector<HTMLElement>('#panel-dashboard');
  if (pulpitMoze === null) return null;
  const pulpit = pulpitMoze;
  return {
    korzen,
    pulpit,
    instrukcje: korzen.querySelector<HTMLElement>('#panel-instructions'),
    pamiec: korzen.querySelector<HTMLElement>('#panel-context'),
    biblioteka: korzen.querySelector<HTMLElement>('#panel-library'),
    agenci: korzen.querySelector<HTMLElement>('#panel-agents'),
    wTle: korzen.querySelector<HTMLElement>('#panel-zadania'),
    plan: korzen.querySelector<HTMLElement>('#panel-plan'),
    artefakty: korzen.querySelector<HTMLElement>('#panel-artefakty'),
    terminal: korzen.querySelector<HTMLElement>('#panel-terminal'),
  };
}

export function cialoPanelu(panel: HTMLElement | null): HTMLElement | null {
  return panel?.querySelector<HTMLElement>('.sta-okno-tresc') ?? null;
}

export function plakietkaPanelu(panel: HTMLElement | null): HTMLElement | null {
  return panel?.querySelector<HTMLElement>('.sta-okno-belka .dn-plakietka') ?? null;
}

export function wpisz(wezel: Element | null | undefined, wartosc: string): void {
  if (wezel === null || wezel === undefined) return;
  wezel.textContent = wartosc;
}

/** Pierwszy wiersz listy staje się wzorem, lista zostaje pusta. */
export function wzorWiersza(lista: Element | null, selektor: string): HTMLElement | null {
  if (lista === null) return null;
  const pierwszy = lista.querySelector<HTMLElement>(selektor);
  const wzor = pierwszy === null ? null : (pierwszy.cloneNode(true) as HTMLElement);
  lista.replaceChildren();
  return wzor;
}

export function kopiaWzoru(wzor: HTMLElement | null): HTMLElement | null {
  return wzor === null ? null : (wzor.cloneNode(true) as HTMLElement);
}

/** Jedno zdanie zamiast treści, gdy rdzeń nie ma czym wypełnić węzła. */
export function niegotowe(gniazdo: Element | null, zdanie: string): void {
  if (gniazdo === null) return;
  const nota = document.createElement('p');
  nota.className = 'dn-meta';
  nota.textContent = zdanie;
  gniazdo.replaceChildren(nota);
}

export function ustawPostep(postep: Element | null, procent: number): void {
  if (postep === null) return;
  const zaokraglony = Math.max(0, Math.min(100, Math.round(procent)));
  postep.setAttribute('data-postep-do', String(zaokraglony));
  const wartosc = postep.querySelector<HTMLElement>('.dn-postep-wartosc');
  if (wartosc !== null) wartosc.style.width = `${zaokraglony}%`;
}

export function udzialProcentowy(czesc: number | undefined, calosc: number | undefined): number {
  if (czesc === undefined || calosc === undefined || calosc <= 0) return 0;
  return (czesc / calosc) * 100;
}

export function przyjmij<T>(tytul: string, wynik: Wynik<T>): T | null {
  if (wynik.udany && wynik.wynik !== undefined) return wynik.wynik;
  oglos(tytul, wynik.blad?.message ?? 'Rdzeń nie zwrócił odpowiedzi.', 'blad');
  return null;
}

export function usun(wezel: Element | null | undefined): void {
  wezel?.remove();
}
