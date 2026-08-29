/**
 * Strefa 2 — wstążka narzędziowa okna roboczego. Leży wewnątrz bryły, w jednej
 * barwie z kartą bieżącą, na pełną szerokość okna — odrębna od wstążki
 * aplikacji, która niesie wyłącznie narzędzia poziomu aplikacji (rama).
 *
 * Przyciski układu przestawiają atrybuty, które czyta `stanowisko.css`:
 * `data-czaty` na obszarze i `data-uklad` na strefie roboczej. Przyciski wersji
 * i operacji nie wołają rdzenia same — przenoszą na kartę, na której praca
 * naprawdę stoi, bo to tam panel trzyma dokument okna i obsługę odmowy.
 * `Wyślij do Library` stoi nieczynny: kontrakt nie niesie komendy odkładającej
 * dokument Studia do Library, a przycisk czynny bez komendy byłby atrapą.
 * Tabliczka sesji niesie tytuł sesji bieżącej, gdy rdzeń go podał;
 * w przeciwnym razie nazwany stan pusty, nie wymyślona nazwa.
 */

import type { Session } from '../../../../../shared/contract.ts';
import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

export interface WlasciwosciWstazki {
  /** Karty sesji odtworzone przez rdzeń — pierwsza nazywa tabliczkę sesji. */
  sesje: Session[];
  /** Przenosi pasmo na kartę o podanym kodzie — tak, jak zrobiłby to klik w nią. */
  naKarte(kod: string): void;
}

function znak(rysunek: NazwaZnaku): SVGElement {
  const wezel = zeZnacznika(ikony[rysunek]);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function przyciskIkony(rysunek: NazwaZnaku, etykieta: string, atrybuty?: Record<string, string>): HTMLElement {
  return el(
    'button',
    { klasa: 'dn-nrz-btn dn-etykietka', type: 'button', 'data-etykietka': etykieta, 'aria-label': etykieta, ...atrybuty },
    [znak(rysunek)],
  );
}

function grupa(etykieta: string, dzieci: HTMLElement[]): HTMLElement {
  return el('span', { klasa: 'st-wstazka-grupa', role: 'group', 'aria-label': etykieta }, dzieci);
}

/* Przełącznik atrybutu na węźle wskazanym z wnętrza bryły: wstążka powstaje
   przed bryłą, więc węzły odnajduje dopiero w chwili kliknięcia. */
function przelacz(przycisk: HTMLElement, wybor: string, atrybut: string, wartosc: string): void {
  przycisk.addEventListener('click', () => {
    const bryla = przycisk.closest('.st-okno-robocze');
    const cel = bryla?.querySelector<HTMLElement>(wybor);
    if (!cel) return;
    const wlaczone = cel.getAttribute(atrybut) === wartosc;
    if (wlaczone) cel.removeAttribute(atrybut);
    else cel.setAttribute(atrybut, wartosc);
    przycisk.setAttribute('aria-pressed', wlaczone ? 'false' : 'true');
  });
}

export function wstazkaOkna(w: WlasciwosciWstazki): HTMLElement {
  const nazwaSesji = w.sesje[0]?.title ?? tekst('wstazka.bezSesji');

  const oknoKomunikacji = przyciskIkony('dymek', tekst('wstazka.oknoKomunikacji'), { 'aria-pressed': 'true' });
  przelacz(oknoKomunikacji, '.sta-obszar', 'data-czaty', 'ukryte');

  const podzialPionowy = przyciskIkony('podzial', tekst('wstazka.podzialPionowy'));
  przelacz(podzialPionowy, '.sta-robocza', 'data-uklad', 'siatka');

  const zapiszWersje = przyciskIkony('zapisz', tekst('wstazka.zapiszWersje'));
  zapiszWersje.addEventListener('click', () => w.naKarte('repo'));

  const porownajWersje = przyciskIkony('diff', tekst('wstazka.porownajWersje'));
  porownajWersje.addEventListener('click', () => w.naKarte('diff'));

  const podgladWydruku = przyciskIkony('oko', tekst('wstazka.podgladWydruku'));
  podgladWydruku.addEventListener('click', () => w.naKarte('preview'));

  const uruchomOperacje = el(
    'button',
    { klasa: 'dn-btn dn-btn--atrament dn-btn--sm', type: 'button' },
    [znak('klucz'), el('span', { tekst: tekst('wstazka.uruchomOperacje') })],
  );
  uruchomOperacje.addEventListener('click', () => w.naKarte('tools'));

  const wyslijDoLibrary = el(
    'button',
    {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      disabled: 'disabled',
      title: tekst('wstazka.brakKomendyLibrary'),
    },
    [znak('wyslij'), el('span', { tekst: tekst('wstazka.wyslijDoLibrary') })],
  );

  return el('div', { klasa: 'st-wstazka', role: 'group', 'aria-label': tekst('wstazka.etykieta') }, [
    el('span', { klasa: 'st-wstazka-sesja' }, [znak('olowek'), el('span', { tekst: nazwaSesji })]),
    grupa(tekst('wstazka.ukladEtykieta'), [oknoKomunikacji, podzialPionowy]),
    grupa(tekst('wstazka.wersjeEtykieta'), [zapiszWersje, porownajWersje, podgladWydruku]),
    el(
      'span',
      { klasa: 'st-wstazka-grupa st-wstazka-grupa--drugorzedna', role: 'group', 'aria-label': tekst('wstazka.operacjeEtykieta') },
      [uruchomOperacje, wyslijDoLibrary],
    ),
    el('span', { klasa: 'st-wstazka-odstep' }),
    przyciskIkony('dostosuj', tekst('wstazka.dostosuj')),
  ]);
}
