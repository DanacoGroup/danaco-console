/**
 * Strefa 1 — pasmo siedmiu kart okna roboczego. Znak modułu w narożniku,
 * wykaz kart `role="tablist"` i sterowanie pasma poza wykazem — rola
 * `tablist` nie przyjmuje innych dzieci. Przełączanie kart wiąże `montaz.ts`,
 * bo dotyczy też widoczności paneli w zaczepie, poza tym składnikiem.
 *
 * Zamknięcie karty pomocniczej i jej przywrócenie stoją tutaj, bo obie
 * czynności dotyczą wyłącznie widoku: karta znika z wykazu i wraca, a rdzeń
 * nie prowadzi wykazu kart okna roboczego. Karta robocza zamknięcia nie ma —
 * jest jedyna w oknie. Po zamknięciu karty bieżącej pasmo klika w kartę
 * sąsiednią, bo to montaż wiąże kliknięcie z odsłonięciem panelu i z jego
 * domontowaniem przy pierwszym wejściu. Tę samą czynność niesie przycisk
 * „Zamknij kartę" w belce okna pomocniczego — pasmo odnajduje kartę po
 * identyfikatorze panelu, w którym ten przycisk stoi.
 */

import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';
import { KARTA_EDITOR, karty, type DefinicjaKarty } from './definicje.ts';

/* Napisy pasma zebrane w jednym miejscu pliku: katalog `moduly/studio/tresci.ts`
   leży poza terenem tej zmiany, a tekst wpleciony w kod byłby nie do wydania
   tłumaczowi. {nazwa} podstawia nazwę karty. */
const TRESC = {
  zamknij: 'Zamknij kartę {nazwa}',
  przywrocPozycja: '+ {nazwa}',
  przywrocOpis: 'Przywróć kartę {nazwa}',
} as const;

/* Znak zamknięcia stoi tutaj, a nie w `ikony.ts`: ten plik leży poza terenem
   zmiany. Rysunek przeniesiony z prototypu bez zmian. */
const ZNAK_ZAMKNIECIA =
  '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" ' +
  'stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18M6 6l12 12"/></svg>';

function podstaw(wzor: string, nazwa: string): string {
  return wzor.replace('{nazwa}', nazwa);
}

function znak(rysunek: NazwaZnaku, klasa?: string): SVGElement {
  const wezel = zeZnacznika(ikony[rysunek]);
  if (klasa !== undefined) wezel.setAttribute('class', klasa);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function przyciskIkony(rysunek: NazwaZnaku, etykieta: string, atrybuty?: Record<string, string>): HTMLElement {
  return el(
    'button',
    { klasa: 'dn-btn-ikona dn-etykietka', type: 'button', 'data-etykietka': etykieta, 'aria-label': etykieta, ...atrybuty },
    [znak(rysunek)],
  );
}

export function pasmoKart(): HTMLElement {
  const przywroc = el('span', { klasa: 'st-pasmo-przywroc', 'data-przywroc-pas': true });

  /* Karta odnajdywana po identyfikatorze swojego panelu: przycisk w belce okna
     pomocniczego wie tylko, w którym panelu stoi. */
  const wgPanelu = new Map<string, { wezel: HTMLElement; nazwa: string }>();

  function widoczneKarty(): HTMLElement[] {
    return Array.from(lista.querySelectorAll<HTMLElement>('.dn-karta[data-karta]')).filter((k) => !k.hidden);
  }

  /* Karta odłożona nie może zostać wybrana: wykaz musi mieć dokładnie jedną
     kartę bieżącą także wtedy, gdy zamknięto tę, która nią była. */
  function zamknij(wezel: HTMLElement, nazwa: string): void {
    if (wezel.hidden) return;
    const wykaz = widoczneKarty();
    const miejsce = wykaz.indexOf(wezel);
    const nastepna = wykaz[miejsce + 1] ?? wykaz[miejsce - 1];
    if (nastepna === undefined) return;

    const byla = wezel.getAttribute('aria-selected') === 'true';
    wezel.setAttribute('aria-selected', 'false');
    wezel.tabIndex = -1;
    wezel.hidden = true;
    przywroc.appendChild(przyciskPrzywrocenia(wezel, nazwa));
    if (byla) nastepna.click();
  }

  function przyciskPrzywrocenia(wezel: HTMLElement, nazwa: string): HTMLElement {
    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--zarys dn-btn--sm',
      type: 'button',
      tekst: podstaw(TRESC.przywrocPozycja, nazwa),
      'aria-label': podstaw(TRESC.przywrocOpis, nazwa),
    });
    przycisk.addEventListener('click', () => {
      wezel.hidden = false;
      przycisk.remove();
      wezel.click();
    });
    return przycisk;
  }

  function karta(k: DefinicjaKarty): HTMLElement {
    const robocza = k.kod === KARTA_EDITOR;
    const wezel = el(
      'div',
      {
        klasa: robocza ? 'dn-karta-widoku dn-karta dn-karta--robocza' : 'dn-karta-widoku dn-karta',
        role: 'tab',
        id: `karta-${k.kod}`,
        'aria-controls': `panel-${k.kod}`,
        'aria-selected': robocza ? 'true' : 'false',
        tabindex: robocza ? '0' : '-1',
        'data-karta': k.kod,
      },
      [znak(k.ikona, 'dn-karta-widoku-ikona'), el('span', { klasa: 'dn-karta-widoku-nazwa', tekst: k.nazwa })],
    );
    if (robocza) return wezel;

    const zamkniecie = el(
      'span',
      {
        klasa: 'dn-karta-widoku-zamknij dn-etykietka',
        'data-etykietka': podstaw(TRESC.zamknij, k.nazwa),
        'data-karta-zamknij': true,
        'aria-hidden': 'true',
      },
      [zeZnacznika(ZNAK_ZAMKNIECIA)],
    );
    /* Zatrzymanie wędrówki zdarzenia: bez niego ten sam klik trafiłby jeszcze
       do montażu okna, który odsłoniłby panel karty właśnie zamykanej. */
    zamkniecie.addEventListener('click', (zdarzenie) => {
      zdarzenie.stopPropagation();
      zamknij(wezel, k.nazwa);
    });
    wezel.appendChild(zamkniecie);
    wgPanelu.set(`panel-${k.kod}`, { wezel, nazwa: k.nazwa });

    /* Znak zamknięcia jest wyłącznie wzrokowy (`aria-hidden`), więc bez klawisza
       Delete czynność nie byłaby osiągalna z klawiatury ani dla czytnika. */
    wezel.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key !== 'Delete') return;
      zdarzenie.preventDefault();
      zamknij(wezel, k.nazwa);
    });
    return wezel;
  }

  const lista = el(
    'div',
    { klasa: 'st-karty', role: 'tablist', 'aria-label': tekst('pasmo.etykietaKart') },
    karty.map(karta),
  );

  const sterowanie = el('span', { klasa: 'st-pasmo-sterowanie' }, [
    przywroc,
    przyciskIkony('plus', tekst('pasmo.nowe')),
    przyciskIkony('maksymalizuj', tekst('pasmo.maksymalizuj'), { 'aria-pressed': 'false' }),
  ]);

  const pasmo = el('div', { klasa: 'dn-karty-pasmo st-pasmo', 'data-gestosc': 'ciasna' }, [
    el('span', { klasa: 'st-pasmo-znak', 'aria-hidden': 'true' }, [znak('olowek')]),
    lista,
    sterowanie,
  ]);

  /* Przycisk „Zamknij kartę" w belce okna pomocniczego zamyka tę samą kartę co
     znak w paśmie. Nasłuch siada na bryle okna roboczego, bo panele powstają
     poza tym składnikiem; bryła powstaje po paśmie, stąd odłożenie o obieg. */
  queueMicrotask(() => {
    const bryla = pasmo.closest('.st-okno-robocze');
    bryla?.addEventListener('click', (zdarzenie) => {
      const cel = zdarzenie.target;
      if (!(cel instanceof Element)) return;
      const zamkniecie = cel.closest('[data-karta-zamknij]');
      if (zamkniecie === null || pasmo.contains(zamkniecie)) return;
      const okno = zamkniecie.closest('.sta-okno');
      const wpis = okno === null ? undefined : wgPanelu.get(okno.id);
      if (wpis === undefined) return;
      zamknij(wpis.wezel, wpis.nazwa);
    });
  });

  return pasmo;
}
