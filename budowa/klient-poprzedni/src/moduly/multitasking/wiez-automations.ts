import type { AutomationWorkflow } from '../../../../shared/contract';
import { poleWyboru, ustawPozycje } from '../../modele/kontrolki-formularza';
import { zdaniePuste, zObjasnieniem } from './powierzchnia-sekcji';

/**
 * Więź sekcji panelu orkiestracji z oknem modułu Automations. Sekcje Kolejki,
 * Orkiestracja, Harmonogram i Monitor nie mają własnych okien — pracują na oknach
 * modułu Automations po skonfigurowaniu powiązania.
 */
export interface WiezAutomations {
  /** Element podsekcji „Powiązanie z modułem Automations". */
  element: HTMLElement;
  /** Wskazana automatyka; pusty napis znaczy „powiązania nie ma". */
  uklad(): string;
  /** Wymienia wykaz automatyk po odczycie z rdzenia. */
  ustawWykaz(automatyki: readonly AutomationWorkflow[]): void;
  /** Zgłasza zmianę wskazania — sekcja przeładowuje wtedy swoją treść. */
  naZmiane(sluchacz: () => void): void;
}

export interface OpcjeWiezi {
  /** Nazwa okna Automations, na którym ta sekcja pracuje. */
  okno: string;
  /** Co sekcja robi na tym oknie — jedno zdanie dla Operatora. */
  rola: string;
}

export function utworzWiezAutomations(opcje: OpcjeWiezi): WiezAutomations {
  const sluchacze: Array<() => void> = [];

  const wybor = poleWyboru(
    { etykieta: 'Automatyka (układ) modułu Automations' },
    [{ wartosc: '', etykieta: 'Bez powiązania' }],
  );
  wybor.kontrolka.addEventListener('change', () => {
    for (const sluchacz of sluchacze) sluchacz();
  });

  const stan = document.createElement('p');
  stan.className = 'dn-pole-opis';

  const element = document.createElement('div');
  element.className = 'dm-orkiestracja__wiez';
  element.append(
    zdaniePuste(`Ta sekcja pracuje na oknie „${opcje.okno}" modułu Automations: ${opcje.rola}`),
    zObjasnieniem(
      wybor.element,
      'Wskazuje układ automatyki, na którym pracuje ta sekcja. Wybór nie uruchamia niczego sam — ustala adres, pod który jadą komendy zależności i harmonogramu. Integracja z Automations jest jawną i odwracalną decyzją Operatora, nie stanem domyślnym.',
    ),
    stan,
  );

  /** Zdanie o stanie powiązania; wypełnia akapit przy każdej zmianie wskazania. */
  function opiszStan(ile: number): void {
    if (ile === 0) {
      stan.textContent =
        'Rdzeń nie ma ani jednej automatyki, więc powiązania nie ma z czym założyć. Układ zakłada się w module Automations (okno Workflow Builder); po zapisaniu pojawi się na tej liście.';
      return;
    }
    stan.textContent =
      wybor.kontrolka.value === ''
        ? `Powiązanie NIE JEST skonfigurowane. Rdzeń zna ${ile} automatyk — wskaż jedną wyżej, żeby ta sekcja miała na czym pracować.`
        : `Powiązanie skonfigurowane: układ ${wybor.kontrolka.value}. Zdjęcie powiązania to wybór „Bez powiązania" — decyzja jest odwracalna.`;
  }

  sluchacze.push(() => opiszStan(wybor.kontrolka.length - 1));
  opiszStan(0);

  return {
    element,
    uklad: () => wybor.kontrolka.value,

    ustawWykaz(automatyki) {
      ustawPozycje(wybor.kontrolka, [
        { wartosc: '', etykieta: 'Bez powiązania' },
        ...automatyki.map((uklad) => ({
          wartosc: uklad.id,
          etykieta: uklad.enabled ? uklad.name : `${uklad.name} (wyłączona)`,
        })),
      ]);
      opiszStan(automatyki.length);
    },

    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),
  };
}
