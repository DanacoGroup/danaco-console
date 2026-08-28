import type { ResearchReportSection } from '../../../../shared/contract';
import { poleTekstowe, poleWielowierszowe, przycisk } from '../../modele/kontrolki-formularza';
import { nowyIdentyfikator as identyfikatorKlienta } from '../../protokol/identyfikator';

/**
 * Redakcja sekcji raportu Report Buildera wymienia sekcje jako pola edycji wraz z ich
 * kolejnością, gotowe do żądania kontraktu.
 */
export interface SekcjeRaportu {
  element: HTMLElement;
  /** Dokłada pustą sekcję i przenosi do niej ognisko. */
  dodaj(): void;
  /** Sekcje wypełnione, w kolejności wierszy, gotowe do żądania; wiersz pusty w obu polach nie wychodzi. */
  zebrane(): ResearchReportSection[];
  /** Wymienia wykaz sekcji na sekcje raportu potwierdzonego przez rdzeń. */
  wczytaj(sekcje: readonly ResearchReportSection[]): void;
  /** Liczba wierszy redakcji, także pustych — miara widoku, nie żądania. */
  liczba(): number;
}

interface WierszSekcji {
  element: HTMLElement;
  identyfikator: string;
  tytul: HTMLInputElement;
  tresc: HTMLTextAreaElement;
  /** Ustalenia przypisane sekcji przez rdzeń; wracają w żądaniu bez zmiany. */
  ustalenia: readonly string[];
}

export function utworzSekcjeRaportu(): SekcjeRaportu {
  const element = document.createElement('div');
  element.className = 'mr-sekcje';
  let wiersze: WierszSekcji[] = [];

  function przerysuj(): void {
    element.replaceChildren(...wiersze.map((wiersz) => wiersz.element));
  }

  function usun(identyfikator: string): void {
    wiersze = wiersze.filter((wiersz) => wiersz.identyfikator !== identyfikator);
    przerysuj();
  }

  function dodajWiersz(sekcja: ResearchReportSection): WierszSekcji {
    const wiersz = utworzWiersz(sekcja, usun);
    wiersze = [...wiersze, wiersz];
    przerysuj();
    return wiersz;
  }

  return {
    element,

    dodaj() {
      const wiersz = dodajWiersz({ id: nowyIdentyfikator(), title: '' });
      wiersz.tytul.focus();
    },

    // Kolejność liczona od jednego, bo zero rdzeń traktuje jak brak wskazania
    // i podstawia numer wiersza.
    zebrane: () =>
      wiersze.filter(wypelniony).map((wiersz, kolejnosc) => ({
        id: wiersz.identyfikator,
        title: wiersz.tytul.value.trim(),
        content: wiersz.tresc.value,
        ...(wiersz.ustalenia.length > 0 ? { findingIds: [...wiersz.ustalenia] } : {}),
        order: kolejnosc + 1,
      })),

    wczytaj(sekcje) {
      wiersze = [];
      for (const sekcja of sekcje) dodajWiersz(sekcja);
      przerysuj();
    },

    liczba: () => wiersze.length,
  };
}

/**
 * Czy wiersz redakcji jest sekcją dokumentu.
 *
 * Sam nagłówek wystarczy, sama treść wystarczy; wiersz pusty w obu polach jest
 * miejscem na pisanie, nie częścią raportu.
 */
function wypelniony(wiersz: WierszSekcji): boolean {
  return wiersz.tytul.value.trim() !== '' || wiersz.tresc.value.trim() !== '';
}

/**
 * Identyfikator sekcji nadawany po stronie klienta, dopóki rdzeń nie odda swojego,
 * niepowtarzalny poza sesją okna dzięki wspólnemu generatorowi.
 */
function nowyIdentyfikator(): string {
  return identyfikatorKlienta('sekcja');
}

/** Jeden wiersz redakcji sekcji raportu: pole tytułu, pole treści oraz zdjęcie sekcji, z jakiej wiersz powstał, wraz z jej ustaleniami. */
function utworzWiersz(
  sekcja: ResearchReportSection,
  naUsuniecie: (identyfikator: string) => void,
): WierszSekcji {
  const tytul = poleTekstowe({ etykieta: 'Tytuł sekcji', podpowiedz: 'nagłówek części raportu' });
  tytul.kontrolka.value = sekcja.title;

  const tresc = poleWielowierszowe({ etykieta: 'Treść sekcji', podpowiedz: 'redakcja części raportu' }, 3);
  tresc.kontrolka.value = sekcja.content ?? '';

  const zdejmij = przycisk('Zdejmij sekcję', 'dn-btn dn-btn--sm dn-btn--duch');
  zdejmij.addEventListener('click', () => naUsuniecie(sekcja.id));

  const element = document.createElement('div');
  element.className = 'mr-sekcje__wiersz';
  element.dataset['sekcja'] = sekcja.id;
  element.append(tytul.element, tresc.element, zdejmij);

  return {
    element,
    identyfikator: sekcja.id,
    tytul: tytul.kontrolka,
    tresc: tresc.kontrolka,
    ustalenia: sekcja.findingIds ?? [],
  };
}
