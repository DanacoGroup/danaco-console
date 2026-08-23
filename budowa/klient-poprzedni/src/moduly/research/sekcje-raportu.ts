import type { ResearchReportSection } from '../../../../shared/contract';
import { poleTekstowe, poleWielowierszowe, przycisk } from '../../modele/kontrolki-formularza';
import { nowyIdentyfikator as identyfikatorKlienta } from '../../protokol/identyfikator';

/**
 * Redakcja sekcji raportu Report Buildera: wykaz sekcji jako pola edycji wraz
 * z ich kolejnością.
 *
 * Sekcje jadą w polu `sections` komendy `research.report.build`, więc redakcja
 * nie potrzebuje osobnej komendy; kolejność bierze się z porządku wierszy
 * i wchodzi w pole `order`. Sekcja bez tytułu nie blokuje zapisu — kreator
 * wysyła to, co ma, a brak nazywa komunikatem przy polu.
 *
 * Wiersz pusty w obu polach nie jest sekcją i tu wypada. Podanie `sections`
 * przestawia rdzeń na gałąź zapisu wprost: nie woła modelu i nie zamienia
 * zaznaczonych ustaleń na sekcje. Pusta tablica też jest podanymi sekcjami,
 * więc warstwa żądania pomija pole `sections` w całości, gdy nie zostało nic —
 * inaczej kreator, który dokłada pusty wiersz przy otwarciu, składałby raport
 * z jednej sekcji bez tytułu i bez treści zamiast streszczenia ustaleń.
 *
 * Wiersz wczytany z raportu niesie dalej swoje `findingIds`, więc powtórne
 * złożenie nie zrywa wiązania ustalenie↔sekcja.
 */
export interface SekcjeRaportu {
  element: HTMLElement;
  /** Dokłada pustą sekcję i przenosi do niej ognisko. */
  dodaj(): void;
  /**
   * Sekcje wypełnione, w kolejności wierszy, gotowe do żądania kontraktu.
   * Wiersz bez tytułu i bez treści nie wychodzi.
   */
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
 * Identyfikator sekcji nadawany po stronie klienta, dopóki rdzeń nie odda
 * swojego.
 *
 * Musi być niepowtarzalny poza sesją okna, nie tylko w niej: kolumna
 * `sekcja_raportu_badania.identyfikator_zewnetrzny` jest unikalna globalnie,
 * a nie w obrębie raportu, więc kod powtórzony przy nowym raporcie kończy się
 * odmową rdzenia. Ten sam kod przy tym samym raporcie jest w porządku — rdzeń
 * poprawia sekcję zastaną. Człon losowy bierzemy ze wspólnego
 * `protokol/identyfikator`, żeby nie hodować drugiej reguły nadawania
 * identyfikatorów w kliencie.
 */
function nowyIdentyfikator(): string {
  return identyfikatorKlienta('sekcja');
}

/** Jeden wiersz redakcji: tytuł, treść i zdjęcie sekcji z raportu. */
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
