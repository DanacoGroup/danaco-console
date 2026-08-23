import type { BrowserNote, BrowserSource } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { poleWielowierszowe } from '../../modele/kontrolki-formularza';
import { KLASY_DYMKA, OBJASNIENIA } from './etykiety-browser';
import { utworzSterWyboru } from './ster-wyboru';
import type { TrescNotatki } from './czynnosci-notatek';

/**
 * Formularz notatki Notes Panel — treść, cytat, powiązane źródło i moduł
 * docelowy przekazania.
 *
 * Jedna odpowiedzialność: pola notatki wraz z ich stanem błędu. Wykaz notatek
 * i rozmowa z rdzeniem są osobno.
 *
 * Błąd pola nie kasuje formularza: alert liniowy staje pod polem, pole pozostaje
 * edytowalne, a wpisana treść zostaje na miejscu.
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
   * Wypełnia formularz treścią notatki (pozycja „Edytuj”).
   *
   * Oddaje `false`, gdy źródła notatki nie ma już w wykazie okna — ster wraca
   * wtedy do „bez powiązania”, a okno ma o czym powiedzieć. Ciche przełknięcie
   * tej różnicy zamieniłoby zapis w notatkę o innym powiązaniu niż pierwowzór.
   */
  wypelnij(notatka: BrowserNote): boolean;
  /** Wstawia zaznaczony fragment strony jako cytat. */
  ustawCytat(fragment: string): void;
  /** Odświeża wykaz źródeł do powiązania. */
  ustawZrodla(zrodla: readonly BrowserSource[]): void;
}

/** Dwa moduły docelowe pozycji „→ Wyślij do Research/Library". */
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

/** Pozycja steru źródeł, gdy notatka nie ma być z niczym wiązana. */
const BEZ_POWIAZANIA = { wartosc: '', etykieta: 'bez powiązania' };

export function utworzFormularzNotatki(): FormularzNotatki {
  const tresc = poleWielowierszowe({ etykieta: 'Treść notatki' }, 4);
  const cytat = poleWielowierszowe({ etykieta: 'Cytowany fragment strony' }, 2);
  // Dwa stery zamiast dwóch natywnych list wyboru: nastawa stoi na uchwycie,
  // a mechanizm rozwijania jest jeden dla całego produktu
  // (`komponenty/menu-drzewo.ts`). Podpis zostaje, bo oba stery stoją w rzędzie
  // pól formularza, gdzie sama wartość nie mówi, czego dotyczy.
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
      // Alert znika razem ze znacznikiem niepoprawności: pole zdjęte z alertu,
      // a wciąż oznaczone `aria-invalid`, czyta się jako błędne bez powodu.
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
