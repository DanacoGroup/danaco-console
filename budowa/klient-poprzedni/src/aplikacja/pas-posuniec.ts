import './pas-posuniec.css';

import type { Posuniecie, PracaAsystenta, ZrodloPosuniec } from './zrodlo-posuniec';

/**
 * Pas posunięć asystenta — wskaźnik pracy i wykaz ruchów na jednym pasku przy dolnej
 * krawędzi. Pas mieści wiersz stanu i wiersz ostatniego posunięcia, stoi poza torem
 * pracy, nie przechwytuje wskaźnika poza własnym przyciskiem i nie zabiera ogniska.
 */
export interface PasPosuniec {
  /** Element montowany w korzeniu powłoki środowiska. */
  element: HTMLElement;
  /** Stawia na pasie zaproszenie do przejścia, którego widok nie wykonał sam. */
  zaproponujPrzejscie(opis: string, wykonaj: () => void): void;
  rozlacz(): void;
}

/**
 * Ile posunięć pas trzyma w pamięci rozwiniętego wykazu. Wykaz jest podręczny, a nie
 * archiwalny, więc pojemność jest skończona i najstarsze wpisy schodzą, zamiast rosnąć
 * przez cały czas pracy okna.
 */
const POJEMNOSC_WYKAZU = 8;

export function utworzPasPosuniec(zrodlo: ZrodloPosuniec): PasPosuniec {
  const element = document.createElement('aside');
  element.className = 'dc-pas-asystenta';
  element.dataset['pasAsystenta'] = 'tak';
  // Nazwa idzie atrybutem `aria-label`; pas jest obszarem uzupełniającym, nie sekcją.
  element.setAttribute('aria-label', 'Posunięcia asystenta');

  const kropka = document.createElement('span');
  kropka.className = 'dc-pas-asystenta__kropka';
  kropka.setAttribute('aria-hidden', 'true');

  const stan = document.createElement('p');
  stan.className = 'dc-pas-asystenta__stan';
  stan.setAttribute('role', 'status');

  const ostatnie = document.createElement('p');
  ostatnie.className = 'dc-pas-asystenta__ostatnie';
  ostatnie.setAttribute('aria-live', 'polite');

  const przejscie = document.createElement('button');
  przejscie.type = 'button';
  przejscie.className = 'dc-pas-asystenta__przejscie';
  przejscie.hidden = true;

  const wykaz = document.createElement('ol');
  wykaz.className = 'dc-pas-asystenta__wykaz';

  element.append(kropka, stan, ostatnie, przejscie, wykaz);

  /** Nanosi na pas stan pracy asystenta; `null` znaczy brak zlecenia. */
  function opiszPrace(praca: PracaAsystenta | null): void {
    element.dataset['praca'] = praca === null ? 'nie' : 'tak';
    if (praca === null) {
      stan.textContent = 'Asystent nie prowadzi zlecenia';
      return;
    }
    const etapy =
      praca.etapow > 0 ? ` · etap ${praca.etap} z ${praca.etapow}` : '';
    stan.textContent = `Asystent pracuje: ${praca.tytul}${etapy}`;
  }

  function dopisz(posuniecie: Posuniecie): void {
    // Pewność jedzie razem z opisem, bo bez niej domysł czytałoby się jak fakt.
    const zrodloRuchu =
      posuniecie.pewnosc === 'pewne'
        ? 'inne połączenie tego konta'
        : 'spoza tego połączenia';
    const zdanie = `${posuniecie.opis} — ${zrodloRuchu}`;

    ostatnie.textContent = zdanie;
    element.dataset['posuniecie'] = posuniecie.rodzaj;

    const wiersz = document.createElement('li');
    wiersz.className = 'dc-pas-asystenta__wiersz';
    wiersz.textContent = zdanie;
    wykaz.prepend(wiersz);
    while (wykaz.children.length > POJEMNOSC_WYKAZU) wykaz.lastElementChild?.remove();
  }

  const odsubskrybowania = [
    zrodlo.naPrace(opiszPrace),
    zrodlo.naPosuniecie(dopisz),
  ];
  opiszPrace(zrodlo.praca());

  let wykonajPrzejscie: (() => void) | null = null;
  przejscie.addEventListener('click', () => {
    const wykonaj = wykonajPrzejscie;
    wykonajPrzejscie = null;
    przejscie.hidden = true;
    wykonaj?.();
  });

  return {
    element,

    zaproponujPrzejscie(opis, wykonaj) {
      przejscie.textContent = opis;
      przejscie.hidden = false;
      wykonajPrzejscie = wykonaj;
    },

    rozlacz() {
      for (const odsubskrybuj of odsubskrybowania.splice(0)) odsubskrybuj();
      element.remove();
    },
  };
}
