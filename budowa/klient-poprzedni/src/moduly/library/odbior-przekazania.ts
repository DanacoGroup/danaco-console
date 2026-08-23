import { type Window } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { przycisk } from '../../modele/kontrolki-formularza';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { StanBiblioteki } from './stan-biblioteki';
import { KOD_MODULU, type ZrodloOtoczenia } from './zrodlo-otoczenia';

/**
 * Odbiór przekazania kontekstu z innego modułu.
 *
 * Przekazanie (`context.transfer`) zakłada w rdzeniu okno modułu Library
 * i rozgłasza wyłącznie `window.changed`; pliku w repozytorium nie zakłada,
 * więc `library.file.changed` po przekazaniu nie przychodzi. Bez tej
 * subskrypcji przybycie kompletu byłoby dla modułu nieme.
 *
 * Zdarzenie `window.changed` niesie okno, nie powód jego zmiany, a `moduleId`
 * biblioteki ma także okno, w imieniu którego moduł sam działa (przestawia je
 * `workspace.enter`). Przekazanie rozpoznają więc dwa warunki:
 *   1. okno jest inne niż okno modułu (`stan.idOkna()`), a okno modułu jest już
 *      znane — dopóki pasek kontekstu nie odczytał własnego okna, nie ma czego
 *      porównywać i moduł milczy;
 *   2. rdzeń ma dla tego okna komplet z dokumentami — okno, którego nikt nie
 *      przekazywał, oddaje komplet bez `documentIds`.
 *
 * Identyfikatory z kompletu są identyfikatorami dokumentów modułu nadawcy
 * i biblioteka może ich nie znać. Przycisk odświeża więc wykaz i mówi wprost,
 * czy rdzeń ma pod tym identyfikatorem plik: gdy ma — wskazuje go jako plik
 * czynny (podgląd, wersje i etykiety przestawiają się razem), gdy nie ma —
 * nazywa to brakiem, zamiast pokazać pustą pozycję.
 */
export interface OdbiorPrzekazania {
  element: HTMLElement;
  /** Odpina subskrypcję `window.changed`. */
  rozlacz(): void;
}

/** Zdanie o przybyciu kompletu — osobno od widoku, żeby dało się je sprawdzić. */
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

  /**
   * Znacznik trwającego odczytu — zamiast wygaszenia przycisku.
   *
   * Odczyt w toku nie jest powodem do odebrania kontrolki: przycisk zostaje
   * w pełni klikalny, a powtórne naciśnięcie mówi, co się właśnie dzieje,
   * zamiast wysyłać drugie żądanie o to samo.
   */
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
        // Odmowa odczytu treści to nie brak przekazania. Okno modułu Library
        // powstało bez udziału Operatora, więc mówimy i o nim, i o powodzie,
        // dla którego treść kompletu została nieodczytana.
        pokaz(
          `Rdzeń założył okno ${okno.id} modułu ${okno.moduleId}, ale treści kompletu nie ` +
            `oddał. ${opisOdmowy('Odczyt kompletu okna', komplet.blad?.code, komplet.blad?.message)}`,
          'bez-tresci',
          false,
        );
        return;
      }
      const przyniesione = komplet.wynik.context.documentIds ?? [];
      // Komplet bez dokumentów ma także okno, którego nikt nie przekazywał,
      // więc to nie jest przybycie i nie ma o nim zdania — inaczej moduł
      // ogłaszałby przekazanie po każdym cudzym oknie.
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
