import type { BrowserNote, BrowserSource } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { poleWielowierszowe } from '../../modele/kontrolki-formularza';
import { KLASY_DYMKA, OBJASNIENIA } from './etykiety-browser';
import { utworzSterWyboru } from './ster-wyboru';
import type { TrescNotatki } from './czynnosci-notatek';

/**
 * Formularz notatki Notes Panel — treść, cytat, powiązane źródło i moduł
 * docelowy przekazania. Jedna odpowiedzialność: pola notatki wraz z ich stanem
 * błędu; wykaz notatek i rozmowa z rdzeniem są osobno.
 */
export interface FormularzNotatki {
  element: HTMLElement;
  /** Treść formularza gotowa do wysłania. */
  odczytaj(): TrescNotatki;
  /** Moduł docelowy pozycji „→ Wyślij do Research/Library". */
  modulDocelowy(): string;
  /** Pokazuje albo zdejmuje alert liniowy pola treści. */
  pokazBlad(zdanie: string): void;
  /**
   * Wypełnia formularz treścią notatki; oddaje `false`, gdy źródła nie ma
   * już w wykazie okna.
   */
  wypelnij(notatka: BrowserNote): boolean;
  /** Wstawia zaznaczony fragment strony jako cytat. */
  ustawCytat(fragment: string): void;
  /** Odświeża wykaz źródeł do powiązania. */
  ustawZrodla(zrodla: readonly BrowserSource[]): void;
}

/**
 * Dwa moduły docelowe pozycji „→ Wyślij do Research/Library", jedyne, do
 * których formularz pozwala notatkę przekazać.
 */
const MODULY_DOCELOWE = [
  {
    wartosc: 'research',
    etykieta: 'Research — Sources Manager',
    opis: 'Notatki idą do wykazu materiałów badania.',
  },
  {
    wartosc: 'library',
    etykieta: 'Library — Library Explorer',
    opis: 'Notatki idą do biblioteki dokumentów sesji.',
  },
];

/**
 * Pozycja steru źródeł, gdy notatka nie ma być z niczym wiązana; jest zarazem
 * wartością wyjściową formularza.
 */
const BEZ_POWIAZANIA = { wartosc: '', etykieta: 'bez powiązania' };

export function utworzFormularzNotatki(): FormularzNotatki {
  const tresc = poleWielowierszowe({ etykieta: 'Treść notatki' }, 4);
  const cytat = poleWielowierszowe({ etykieta: 'Cytowany fragment strony' }, 2);
  // Dwa stery zamiast natywnych list wyboru: mechanizm rozwijania jest jeden
  // dla całego produktu.
  const zrodlo = utworzSterWyboru({
    nastawa: 'Powiązane źródło',
    pozycje: [BEZ_POWIAZANIA],
    klasa: 'dn-pole mb-notatki__ster',
    podpis: true,
  });
  const cel = utworzSterWyboru({
    nastawa: 'Moduł docelowy przekazania',
    pozycje: MODULY_DOCELOWE,
    klasa: 'dn-pole mb-notatki__ster',
    podpis: true,
  });

  const blad = document.createElement('p');
  blad.className = 'dn-pole-blad mb-notatki__blad';
  blad.hidden = true;

  const element = document.createElement('div');
  element.className = 'mb-notatki__formularz';
  element.append(
    tresc.element,
    blad,
    cytat.element,
    utworzDymekObjasnienia(OBJASNIENIA.cytat, KLASY_DYMKA),
    zrodlo.element,
    utworzDymekObjasnienia(OBJASNIENIA.powiazanie, KLASY_DYMKA),
    cel.element,
  );

  return {
    element,

    odczytaj: () => ({
      tresc: tresc.kontrolka.value.trim(),
      idZrodla: zrodlo.wartosc(),
      cytat: cytat.kontrolka.value.trim(),
    }),

    modulDocelowy: () => cel.wartosc(),

    pokazBlad(zdanie) {
      blad.textContent = zdanie;
      blad.hidden = zdanie === '';
      tresc.kontrolka.setAttribute('aria-invalid', zdanie === '' ? 'false' : 'true');
    },

    wypelnij(notatka) {
      tresc.kontrolka.value = notatka.content;
      cytat.kontrolka.value = notatka.quote ?? '';
      blad.textContent = '';
      blad.hidden = true;
      // Alert znika ze znacznikiem niepoprawności: samo `aria-invalid`
      // czytałoby się jako błąd bez powodu.
      tresc.kontrolka.setAttribute('aria-invalid', 'false');
      return zrodlo.ustawWartosc(notatka.sourceId ?? '');
    },

    ustawCytat(fragment) {
      cytat.kontrolka.value = fragment;
      // Ognisko idzie na treść: po wstawieniu cytatu notatkę trzeba dopisać.
      tresc.kontrolka.focus();
    },

    ustawZrodla(zrodla) {
      zrodlo.ustawPozycje([
        BEZ_POWIAZANIA,
        ...zrodla.map((wpis) => ({
          wartosc: wpis.id,
          etykieta: (wpis.title ?? '').trim() === '' ? wpis.url : (wpis.title ?? ''),
          opis: wpis.url,
        })),
      ]);
    },
  };
}
