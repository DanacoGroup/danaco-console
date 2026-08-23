import { adresPoczatkowy, opisAdresu } from './adres-ustawienia';
import type { PunktWidzenia } from './rozstrzygniecie';
import { utworzWyborAdresu } from './wybor-adresu';
import { OSIE_OD_NAJWEZSZEJ, ZASIEGI_OD_NAJWEZSZEGO } from './zasiegi';

/**
 * Pasek punktu widzenia — selektor zasięgu okna konfiguracji.
 *
 * Wartość ustawienia nie jest jedna: zależy od tego, dla kogo pytamy.
 * Ten sam klucz może mieć inną wartość globalnie, inną w oknie komunikacji
 * i jeszcze inną dla wskazanego modelu. Pasek ustala, **dla kogo** okno liczy
 * wartości obowiązujące; każde pole formularza przelicza wtedy swoje
 * pochodzenie od nowa.
 *
 * Zmiana punktu widzenia niczego nie zapisuje. To wyłącznie soczewka —
 * odczyt, nie polecenie.
 */
export interface PasekPunktuWidzenia {
  /** Pas nad ciałem okna. */
  element: HTMLElement;
  /** Punkt widzenia wskazany kontrolkami. */
  punkt(): PunktWidzenia;
  /** Subskrypcja zmiany punktu widzenia. */
  naZmiane(sluchacz: (punkt: PunktWidzenia) => void): void;
}

export function utworzPasekPunktuWidzenia(): PasekPunktuWidzenia {
  const wybor = utworzWyborAdresu(ZASIEGI_OD_NAJWEZSZEGO, OSIE_OD_NAJWEZSZEJ);
  wybor.ustaw(adresPoczatkowy());

  const opis = document.createElement('p');
  opis.className = 'dk-punkt__opis';

  const naglowek = document.createElement('div');
  naglowek.className = 'dk-punkt__naglowek';

  const tytul = document.createElement('h3');
  tytul.className = 'dk-punkt__tytul';
  tytul.textContent = 'Punkt widzenia';

  const wyjasnienie = document.createElement('p');
  wyjasnienie.className = 'dn-pole-opis';
  wyjasnienie.textContent =
    'Dla kogo okno liczy wartości obowiązujące. Zmiana punktu widzenia niczego nie zapisuje.';

  naglowek.append(tytul, wyjasnienie);

  const element = document.createElement('section');
  element.className = 'dk-punkt';
  element.append(naglowek, wybor.element, opis);

  const sluchacze: Array<(punkt: PunktWidzenia) => void> = [];

  function ubierz(): void {
    opis.textContent = `Wartości obowiązujące liczone dla: ${opisAdresu(wybor.adres())}`;
  }

  wybor.naZmiane(() => {
    ubierz();
    const punkt = wybor.adres();
    for (const sluchacz of [...sluchacze]) sluchacz(punkt);
  });

  ubierz();

  return {
    element,
    punkt: wybor.adres,
    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),
  };
}
