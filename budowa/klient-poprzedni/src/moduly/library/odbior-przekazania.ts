import { type Window } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { przycisk } from '../../modele/kontrolki-formularza';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { StanBiblioteki } from './stan-biblioteki';
import { KOD_MODULU, type ZrodloOtoczenia } from './zrodlo-otoczenia';

/**
 * Odbiór przekazania kontekstu z innego modułu subskrybuje zmianę okna, bo
 * przekazanie zakłada okno bez pliku w repozytorium, i rozpoznaje przybycie po
 * innym oknie z kompletem dokumentów.
 */
export interface OdbiorPrzekazania {
  element: HTMLElement;
  /** Odpina subskrypcję `window.changed`. */
  rozlacz(): void;
}

/** Zdanie o przybyciu kompletu — osobno od widoku, żeby dało się je sprawdzić niezależnie od stanu przycisku. */
export function zdanieOPrzybyciu(okno: Window, dokumenty: readonly string[]): string {
  const nazwa = okno.title === undefined || okno.title === '' ? okno.id : `${okno.title} (${okno.id})`;
  return (
    `Rdzeń założył okno ${nazwa} modułu ${okno.moduleId} — przekazanie kontekstu z innego ` +
    `modułu. Komplet niesie dokumenty (${dokumenty.length}): ${dokumenty.join(', ')}. ` +
    'Przekazanie nie zakłada pliku w repozytorium — odśwież wykaz, żeby sprawdzić, czy ' +
    'biblioteka zna ten dokument.'
  );
}

export function utworzOdbiorPrzekazania(
  stan: StanBiblioteki,
  otoczenie: ZrodloOtoczenia,
): OdbiorPrzekazania {
  const zdanie = document.createElement('p');
  zdanie.className = 'ml-odbior__zdanie';

  const wskaz = przycisk('Odśwież wykaz i wskaż dokument', 'dn-btn dn-btn--sm dn-btn--zarys');
  wskaz.dataset['czynnosc'] = 'odbior-wskaz';
  wskaz.hidden = true;

  const element = document.createElement('div');
  element.className = 'ml-odbior';
  element.dataset['odbior'] = 'spoczynek';
  element.hidden = true;
  element.append(zdanie, wskaz);

  /** Dokumenty ostatniego kompletu — wejście przycisku, nie ozdoba zdania. */
  let dokumenty: readonly string[] = [];

  function pokaz(tresc: string, stanOdbioru: string, zPrzyciskiem: boolean): void {
    zdanie.textContent = tresc;
    element.dataset['odbior'] = stanOdbioru;
    element.hidden = false;
    wskaz.hidden = !zPrzyciskiem;
  }

  // Odczyt w toku nie jest powodem do odebrania kontrolki: przycisk zostaje w pełni klikalny.
  let odczytTrwa = false;

  wskaz.addEventListener('click', () => {
    void (async () => {
      if (odczytTrwa) {
        pokaz(
          'Odczyt wykazu już trwa — poczekaj na odpowiedź rdzenia; drugie żądanie ' +
            'pytałoby o to samo.',
          'przybylo',
          true,
        );
        return;
      }
      odczytTrwa = true;
      await stan.odczytaj('', '');
      odczytTrwa = false;
      const znane = stan.pliki().filter((plik) => dokumenty.includes(plik.id));
      const pierwszy = znane[0];
      if (pierwszy === undefined) {
        pokaz(
          `Wykaz odświeżony — rdzeń nie ma w bibliotece pliku o żadnym z identyfikatorów ` +
            `kompletu (${dokumenty.join(', ')}). Przekazanie przenosi kontekst okna, nie ` +
            'zakłada zasobu; dokument mieszka nadal w module, który go przysłał.',
          'nieznany',
          false,
        );
        return;
      }
      stan.wskaz(pierwszy.id);
      pokaz(
        `Wykaz odświeżony — plik „${pierwszy.name}" (${pierwszy.id}) z kompletu jest ` +
          `w bibliotece i stoi jako plik czynny. Rozpoznanych pozycji kompletu: ` +
          `${znane.length} z ${dokumenty.length}.`,
        'wskazany',
        false,
      );
    })();
  });

  const odsubskrybuj: Odsubskrybuj = otoczenie.naZmianeOkna((tresc) => {
    const okno = tresc.window;
    if (okno.moduleId !== KOD_MODULU) return;
    // Okno modułu nieznane albo to właśnie ono — nie ma o czym mówić.
    if (stan.idOkna() === '' || okno.id === stan.idOkna()) return;
    void (async () => {
      const komplet = await otoczenie.komplet(okno.id);
      if (!komplet.udany || komplet.wynik === undefined) {
        // Odmowa odczytu treści to nie brak przekazania: mówimy o oknie i o powodzie nieodczytania kompletu.
        pokaz(
          `Rdzeń założył okno ${okno.id} modułu ${okno.moduleId}, ale treści kompletu nie ` +
            `oddał. ${opisOdmowy('Odczyt kompletu okna', komplet.blad?.code, komplet.blad?.message)}`,
          'bez-tresci',
          false,
        );
        return;
      }
      const przyniesione = komplet.wynik.context.documentIds ?? [];
      // Komplet bez dokumentów ma też okno nieprzekazywane, więc nie ma o nim zdania przybycia.
      if (przyniesione.length === 0) return;
      dokumenty = przyniesione;
      pokaz(zdanieOPrzybyciu(okno, przyniesione), 'przybylo', true);
    })();
  });

  return {
    element,
    rozlacz: () => odsubskrybuj(),
  };
}
