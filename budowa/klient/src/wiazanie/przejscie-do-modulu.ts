// Jedna droga wejścia z przedsionka do okna modułu dla czterech środowisk.
import { Command } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

// Zapora domyka przejście, gdy zdarzenie końca animacji nie przyjdzie.
const ZAPORA_RUCHU_MS = 600;

function czyRuchZniesiony(): boolean {
  return globalThis.matchMedia('(prefers-reduced-motion: reduce)').matches;
}

function plotnoCentrum(): HTMLElement | null {
  return document.querySelector<HTMLElement>('.cd-plotno');
}

function kafelModulu(kodModulu: string): HTMLElement | null {
  for (const kafel of document.querySelectorAll<HTMLElement>('.pd-kafel[data-modul]')) {
    if ((kafel.dataset.modul ?? '').toLowerCase() === kodModulu) return kafel;
  }
  return null;
}

/* Centrum otwiera moduł wyłącznie z kliknięcia przechwyconego na dokumencie,
   więc wejście bez widocznego kafla podstawia własny węzeł z kodem. */
function otworzModulWCentrum(kodModulu: string, plotno: HTMLElement): void {
  const kafel = kafelModulu(kodModulu);
  if (kafel !== null) {
    kafel.click();
    return;
  }
  const nosnik = document.createElement('div');
  nosnik.className = 'pd-kafel';
  nosnik.dataset.modul = kodModulu;
  nosnik.hidden = true;
  plotno.appendChild(nosnik);
  nosnik.click();
  nosnik.remove();
}

function zagrajRuchWidoku(widok: HTMLElement, stanPrzejscia: string): void {
  widok.dataset.przejscie = stanPrzejscia;
  const domknij = (): void => {
    widok.removeEventListener('animationend', domknij);
    if (widok.dataset.przejscie === stanPrzejscia) delete widok.dataset.przejscie;
  };
  widok.addEventListener('animationend', domknij);
  globalThis.setTimeout(domknij, ZAPORA_RUCHU_MS);
}

/* Wejście z przedsionka do okna modułu: komenda rdzenia i otwarcie widoku idą
   od razu, ruch gra równolegle. Pusty identyfikator znaczy wejście bez sesji. */
export function przejdzDoModulu(kanal: Kanal, kodModulu: string, idSesji = ''): void {
  const kod = kodModulu.trim().toLowerCase();
  if (kod === '') return;
  const plotno = plotnoCentrum();
  if (plotno === null) return;
  // Wiązanie sesji nie wstrzymuje obrazu: nie rozstrzyga o widoku płótna.
  if (idSesji !== '') {
    void wywolaj(kanal, Command.SessionBind, {
      sessionId: idSesji,
      clientId: tozsamoscKlienta().id,
    });
  }
  const schodzacy = plotno.querySelector<HTMLElement>('.cd-tresc--przedsionek:not([hidden])');
  otworzModulWCentrum(kod, plotno);
  if (czyRuchZniesiony()) return;
  if (schodzacy !== null && schodzacy.hidden) zagrajRuchWidoku(schodzacy, 'wychodzi');
  const wchodzacy = plotno.querySelector<HTMLElement>('.cd-tresc:not([hidden])');
  if (wchodzacy !== null) zagrajRuchWidoku(wchodzacy, 'wchodzi');
}
