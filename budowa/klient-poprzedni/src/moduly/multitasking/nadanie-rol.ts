import { WindowRole, type ErrorInfo, type Window } from '../../../../shared/contract';
import { pole, przyciskAkcji as przycisk, wiersz } from '../../modele/kontrolki-formularza';
import type { StanMultitaskingu } from './stan-multitaskingu';
import type { ZrodloOkien } from './zrodlo-okien';

// Nadawanie ról oknom sceny odpowiada za samo nadanie i za zdanie o skutku z odpowiedzi rdzenia.

/** Treść odmowy wraz z kodem kontraktu, złożona w jedno pełne zdanie gotowe do pokazania w treści okna. */
export function powod(blad?: ErrorInfo): string {
  if (blad === undefined) return 'Rdzeń nie podał przyczyny.';
  return `Powód: ${blad.message} (kod ${blad.code}).`;
}

/**
 * Czy rdzeń nadał rolę, o którą prosiło okno.
 *
 * Czysta funkcja danych: porównuje żądanie z tym, co wróciło. Przypięcie do
 * koordynatora obowiązuje wyłącznie wykonawcę — inna rola nie należy do pętli.
 */
export function czyRolaNadana(
  okno: { windowRole?: WindowRole; coordinatorWindowId?: string },
  rola: WindowRole,
  koordynator: Window | null,
): boolean {
  if (okno.windowRole !== rola) return false;
  if (rola !== WindowRole.Executor) return true;
  return koordynator !== null && okno.coordinatorWindowId === koordynator.id;
}

/** Zdanie o założeniu okna wzięte z odpowiedzi rdzenia, nie z żądania, bo zgoda nie jest dowodem skutku. */
export function zdanieOZalozeniu(
  okno: Window,
  rola: WindowRole,
  koordynator: Window | null,
): string {
  if (okno.windowRole !== rola) {
    return `Rdzeń założył okno ${okno.id}, ale nadał mu rolę ${okno.windowRole} zamiast ${rola} — obsada tej roli nie powstała.`;
  }
  if (rola === WindowRole.Executor && okno.coordinatorWindowId !== koordynator?.id) {
    return `Rdzeń założył okno wykonawcy ${okno.id}, ale nie przypiął go pod koordynatora ${koordynator?.id ?? 'brak'} (przypięcie: ${okno.coordinatorWindowId ?? 'brak'}).`;
  }
  return `Rdzeń założył okno roli ${okno.windowRole}: ${okno.id}.`;
}

/**
 * Przepięcie cudzego okna wykonawczego pod koordynatora tej obsady — `role.assign`.
 *
 * Oddaje samo powodzenie; odczyt okien po udanym przepięciu należy do panelu,
 * więc funkcja nie wywołuje go za niego.
 */
export async function przypnijPodKoordynatora(
  zrodlo: ZrodloOkien,
  stan: StanMultitaskingu,
  okno: Window,
  potwierdz: (zdanie: string, udane: boolean) => void,
): Promise<boolean> {
  const koordynator = stan.obsada().koordynator;
  if (koordynator === null) {
    potwierdz('Nie ma koordynatora, pod którego dałoby się przypiąć wykonawcę.', false);
    return false;
  }
  const wynik = await zrodlo.nadajRole({
    windowId: okno.id,
    role: WindowRole.Executor,
    coordinatorWindowId: koordynator.id,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    potwierdz(`Rdzeń odmówił nadania roli wykonawcy. ${powod(wynik.blad)}`, false);
    return false;
  }
  const nadane = wynik.wynik;
  if (!czyRolaNadana({ windowRole: nadane.role, ...(nadane.coordinatorWindowId === undefined ? {} : { coordinatorWindowId: nadane.coordinatorWindowId }) }, WindowRole.Executor, koordynator)) {
    potwierdz(
      `Rdzeń przyjął nadanie roli oknu ${nadane.windowId}, ale oddał rolę ${nadane.role} i przypięcie ${nadane.coordinatorWindowId ?? 'brak'} — przepięcie się nie odbyło.`,
      false,
    );
    return false;
  }
  potwierdz(
    `Rdzeń nadał oknu ${nadane.windowId} rolę wykonawcy pod koordynatorem ${nadane.coordinatorWindowId ?? ''}.`,
    true,
  );
  return true;
}

/**
 * Funkcja dopina więź wykonawcy z koordynatorem drugim żądaniem, bo założenie
 * okna z tą więzią kończy się powodzeniem, ale rdzeń jej nie utrwala, więc
 * obsada zostaje niepełna.
 */
export async function dopnijWiezWykonawcy(
  zrodlo: ZrodloOkien,
  stan: StanMultitaskingu,
  okno: Window,
): Promise<{ udane: boolean; zdanie: string }> {
  let odpowiedz = '';
  const udane = await przypnijPodKoordynatora(zrodlo, stan, okno, (zdanie) => {
    odpowiedz = zdanie;
  });
  return {
    udane,
    zdanie: udane
      ? `Więź z koordynatorem dopięta DRUGIM żądaniem (role.assign), bo window.create jej nie utrwalił. ${odpowiedz}`
      : `Więź z koordynatorem NIE została dopięta — window.create jej nie utrwalił, a naprawcze role.assign też się nie powiodło. ${odpowiedz}`,
  };
}

/** Pasek wcielenia — jedyny wołacz komendy zmiany roli w kliencie, z blokadą, gdy obsada nie ma koordynatora. */
export interface PasekWcielenia {
  element: HTMLElement;
  /** Przerysowuje blokadę po zmianie obsady. */
  odswiez(): void;
}

/**
 * Pasek nadania wcielenia roli koordynatora: wcielenie jest wolnym napisem, bo
 * kontrakt nie definiuje katalogu wcieleń, a bez koordynatora przycisk jest
 * zablokowany.
 */
export function utworzPasekWcielenia(
  zrodlo: ZrodloOkien,
  stan: StanMultitaskingu,
  potwierdz: (zdanie: string, udane: boolean) => void,
): PasekWcielenia {
  const wpis = pole('Wcielenie roli', 'na przykład Validator albo Security Auditor');
  const nadaj = przycisk('Nadaj wcielenie', 'dn-btn dn-btn--sm');

  const element = document.createElement('div');
  element.className = 'dm-obsada__wcielenie';
  element.append(
    wiersz('Wcielenie koordynatora', wpis, {
      klasa: 'dm-wiersz',
      objasnienie:
        'Wcielenie zawęża sposób pracy roli. Katalogu wcieleń nie rozstrzygnięto — kontrakt niesie napis.',
    }),
    nadaj,
  );

  function odswiez(): void {
    const koordynator = stan.obsada().koordynator;
    nadaj.disabled = koordynator === null;
    if (koordynator === null) {
      const zdanie = 'Obsada nie ma koordynatora — nie ma komu nadać wcielenia.';
      nadaj.title = zdanie;
      nadaj.setAttribute('aria-description', zdanie);
      nadaj.dataset['powodBlokady'] = 'tak';
      return;
    }
    nadaj.removeAttribute('title');
    nadaj.removeAttribute('aria-description');
    delete nadaj.dataset['powodBlokady'];
  }

  nadaj.addEventListener('click', () => {
    const koordynator = stan.obsada().koordynator;
    if (koordynator === null) {
      potwierdz('Obsada nie ma koordynatora — nie ma komu nadać wcielenia.', false);
      return;
    }
    const wcielenie = wpis.value.trim();
    if (wcielenie === '') {
      potwierdz('Wcielenie puste nie zostanie nadane — wpisz je najpierw.', false);
      return;
    }
    void zrodlo
      .zmienRole({ windowId: koordynator.id, persona: wcielenie })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          potwierdz(`Rdzeń odmówił zmiany wcielenia. ${powod(wynik.blad)}`, false);
          return;
        }
        // Zdanie mówi, co oddał rdzeń: przyjęcie innego wcielenia nie jest nadaniem tego, o które proszono.
        const oddane = wynik.wynik.persona ?? '';
        if (oddane !== wcielenie) {
          potwierdz(
            `Rdzeń przyjął wywołanie, ale oddał wcielenie ${oddane === '' ? 'puste' : oddane} zamiast ${wcielenie} — wcielenie nie zostało nadane.`,
            false,
          );
          return;
        }
        potwierdz(
          `Rdzeń nadał oknu ${wynik.wynik.windowId} wcielenie ${oddane} przy roli ${wynik.wynik.role}.`,
          true,
        );
      });
  });

  odswiez();
  return { element, odswiez };
}
