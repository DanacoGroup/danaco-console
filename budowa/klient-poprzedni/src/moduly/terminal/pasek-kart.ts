import { elementIkony } from '../../ikony/ikony';
import type { StanTerminala } from './stan-terminala';

/**
 * Pasek kart okna Terminal Tabs — przełączanie między powłokami.
 *
 * Pasek pokazuje komplet kart, wskazuje ognisko i oddaje dwie czynności: wybór
 * i zamknięcie. Karty przypięte stoją przed zwykłymi.
 *
 * Krzyżyk jest zwykłym przyciskiem, a przypięcia pilnuje `zamknijKarte` w oknie
 * wiodącym — dzięki temu odmowa ma jedno miejsce i jeden powód widoczny
 * w stanie treści, zamiast bramy w pasku i obejścia w panelu akcji.
 */
export interface CzynnosciPaska {
  wybierz(idKarty: string): void;
  zamknij(idKarty: string): void;
}

export interface PasekKart {
  element: HTMLElement;
  odswiez(): void;
}

export function pasekKart(stan: StanTerminala, czynnosci: CzynnosciPaska): PasekKart {
  const element = document.createElement('div');
  element.className = 'dt-pasek-kart';
  element.setAttribute('role', 'tablist');
  element.setAttribute('aria-label', 'Karty powłok okna');

  function odswiez(): void {
    const biezaca = stan.kartaBiezaca();
    const uporzadkowane = [...stan.karty()].sort((a, b) => {
      const pierwsza = stan.czyPrzypieta(a.id) ? 0 : 1;
      const druga = stan.czyPrzypieta(b.id) ? 0 : 1;
      return pierwsza - druga;
    });

    element.replaceChildren();
    if (uporzadkowane.length === 0) {
      const pusto = document.createElement('span');
      pusto.className = 'dn-pole-opis';
      pusto.textContent = 'Brak kart.';
      element.append(pusto);
      return;
    }

    for (const karta of uporzadkowane) {
      const zakladka = document.createElement('button');
      zakladka.type = 'button';
      zakladka.className = 'dt-karta';
      zakladka.setAttribute('role', 'tab');
      zakladka.dataset['karta'] = karta.id;
      zakladka.dataset['ognisko'] = String(karta.id === biezaca?.id);
      zakladka.setAttribute('aria-selected', String(karta.id === biezaca?.id));
      zakladka.textContent = karta.title ?? karta.shell;
      // Przypięcie niesie znak z zestawu ikon platformy wraz z etykietą dla
      // czytnika ekranu — stan karty nie jest tu przekazywany samym kształtem.
      if (stan.czyPrzypieta(karta.id)) {
        zakladka.prepend(
          elementIkony('spinacz', { rozmiar: 14, etykieta: 'karta przypięta' }),
        );
      }
      zakladka.addEventListener('click', () => czynnosci.wybierz(karta.id));

      const zamkniecie = document.createElement('button');
      zamkniecie.type = 'button';
      zamkniecie.className = 'dt-karta__zamknij';
      zamkniecie.textContent = '×';
      zamkniecie.setAttribute('aria-label', `Zamknij kartę ${karta.title ?? karta.shell}`);
      zamkniecie.title = stan.czyPrzypieta(karta.id)
        ? 'Karta przypięta — zamknięcie odmówi, dopóki jej nie odepniesz.'
        : `Zamyka kartę ${karta.title ?? karta.shell} w widoku; procesy biegną dalej.`;
      zamkniecie.addEventListener('click', (zdarzenie) => {
        zdarzenie.stopPropagation();
        czynnosci.zamknij(karta.id);
      });

      const opakowanie = document.createElement('span');
      opakowanie.className = 'dt-karta__opakowanie';
      opakowanie.append(zakladka, zamkniecie);
      element.append(opakowanie);
    }
  }

  odswiez();
  return { element, odswiez };
}
