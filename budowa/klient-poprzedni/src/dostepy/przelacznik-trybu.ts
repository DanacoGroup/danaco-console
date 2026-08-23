import { AccessMode, type AccessPoint } from '../../../shared/contract';
import { elementIkony } from '../ikony/ikony';
import { nazwaTrybu, TRYBY } from './nazwy-dostepow';
import { ostrzezenieZapisu } from './ostrzezenie-zapisu';

/**
 * Przełącznik trybu nadania: odczyt albo zapis, per punkt dostępu.
 *
 * Ostrzeżenie należy do przełącznika, nie jest dodatkiem obok niego. Zapis na
 * maszynie chronionej pokazuje pełne zdanie pod przełącznikiem, a nie
 * w podpowiedzi pod kursorem — skutek zapisu na produkcyjnej platformie LEX
 * jest nieodwracalny.
 *
 * Przełącznik nie blokuje zapisu; wybór trybu zostaje po stronie człowieka.
 */
export interface PrzelacznikTrybu {
  /** Element osadzany w karcie punktu albo wierszu nadania. */
  element: HTMLElement;
  /** Tryb wybrany w chwili odczytu. */
  tryb(): AccessMode;
  /** Nanosi tryb bez zgłaszania zmiany. */
  ustaw(tryb: AccessMode): void;
  /** Zgłasza wybór trybu dokonany w przełączniku. */
  naZmiane(sluchacz: (tryb: AccessMode) => void): void;
  /** Nanosi punkt po zmianie — ostrzeżenie liczone jest z jego nazwy maszyny. */
  ustawPunkt(punkt: AccessPoint): void;
}

export interface ZaleznosciPrzelacznika {
  /** Punkt, którego dotyczy tryb — z niego bierze się ostrzeżenie. */
  punkt: AccessPoint;
  /** Tryb początkowy. */
  tryb: AccessMode;
  /** Przedrostek identyfikatorów pól; dwa przełączniki nie dzielą nazwy grupy. */
  identyfikator: string;
}

export function utworzPrzelacznikTrybu(
  zaleznosci: ZaleznosciPrzelacznika,
): PrzelacznikTrybu {
  const { identyfikator } = zaleznosci;
  let punkt = zaleznosci.punkt;
  let wybrany = zaleznosci.tryb;

  const sluchacze: Array<(tryb: AccessMode) => void> = [];
  const pola = new Map<AccessMode, HTMLInputElement>();

  const grupa = document.createElement('div');
  grupa.className = 'dd-tryb__grupa';
  grupa.setAttribute('role', 'radiogroup');
  grupa.setAttribute('aria-label', 'Tryb nadania dostępu');

  for (const tryb of TRYBY) {
    const pole = document.createElement('input');
    pole.type = 'radio';
    pole.className = 'dn-radio dd-tryb__pole';
    pole.name = `${identyfikator}-tryb`;
    pole.id = `${identyfikator}-tryb-${tryb}`;
    pole.value = tryb;
    pole.checked = tryb === wybrany;
    pole.addEventListener('change', () => {
      if (!pole.checked) return;
      wybrany = tryb;
      odswiezOstrzezenie();
      for (const sluchacz of [...sluchacze]) sluchacz(tryb);
    });

    const etykieta = document.createElement('label');
    etykieta.className = 'dd-tryb__etykieta';
    etykieta.htmlFor = pole.id;
    etykieta.textContent = nazwaTrybu(tryb);

    pola.set(tryb, pole);
    grupa.append(pole, etykieta);
  }

  const ostrzezenie = document.createElement('p');
  ostrzezenie.className = 'dd-tryb__ostrzezenie';
  ostrzezenie.setAttribute('role', 'alert');
  ostrzezenie.hidden = true;

  const znak = elementIkony('ostrzezenie', { rozmiar: 16 });
  const tresc = document.createElement('span');
  ostrzezenie.append(znak, tresc);

  const element = document.createElement('div');
  element.className = 'dd-tryb';
  element.append(grupa, ostrzezenie);

  function odswiezOstrzezenie(): void {
    const zdanie = ostrzezenieZapisu(punkt, wybrany);
    tresc.textContent = zdanie;
    ostrzezenie.hidden = zdanie === '';
  }

  odswiezOstrzezenie();

  return {
    element,

    tryb: () => wybrany,

    ustaw(tryb) {
      wybrany = tryb;
      for (const [nazwa, pole] of pola) pole.checked = nazwa === tryb;
      odswiezOstrzezenie();
    },

    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),

    ustawPunkt(nowy) {
      punkt = nowy;
      odswiezOstrzezenie();
    },
  };
}
