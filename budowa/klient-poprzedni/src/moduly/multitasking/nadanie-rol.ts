import { WindowRole, type ErrorInfo, type Window } from '../../../../shared/contract';
import { pole, przyciskAkcji as przycisk, wiersz } from '../../modele/kontrolki-formularza';
import type { StanMultitaskingu } from './stan-multitaskingu';
import type { ZrodloOkien } from './zrodlo-okien';

/**
 * Nadawanie ról oknom sceny MultitaskingAI. Panel obsady odpowiada za scenę —
 * zakłada okna, pokazuje skład i blokuje drugi egzemplarz roli; ten plik
 * odpowiada za samo nadanie roli i za zdanie o jego skutku.
 *
 * Rolę nadaje `role.assign`, nie `window.update`: ta druga komenda przy roli
 * spoza kontraktu odpowiada `status: ok`, a oddaje okno o roli `standalone` bez
 * `coordinatorWindowId`. Wcielenie niesie wyłącznie `role.update`; pasek niżej
 * jest jego jedynym wołaczem. `window.update` zostaje przy tytule i katalogach.
 *
 * Każde zdanie tego pliku powstaje z pól odpowiedzi, nigdy z treści żądania —
 * zgoda rdzenia nie jest dowodem skutku.
 */

/** Treść odmowy wraz z kodem kontraktu. */
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

/** Zdanie o założeniu okna wzięte z odpowiedzi rdzenia, nie z żądania. */
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
 * Dopięcie więzi wykonawcy z koordynatorem drugim żądaniem — jawne obejście
 * usterki rdzenia.
 *
 * `window.create` z polem `coordinatorWindowId` kończy się `ok`, ale więzi nie
 * utrwala: `dolozWiezi` w `budowa/server/internal/core/adapter_okna.go` nadpisuje
 * więź z pamięci wartością z bazy, a wiersz okna zakładany przez `utrwalZalozone`
 * jest w tym miejscu pusty. Obsada zostaje wtedy na „Wykonawcy 0 z 2", a okna
 * wykonawców nie mają adresata polecenia.
 *
 * Funkcja powtarza więc `role.assign` zaraz po założeniu okna i mówi w zdaniu
 * wprost, że więź poszła drugim żądaniem. Wołający sprawdza obsadę odczytaną
 * `window.list` już po założeniu okna, więc gdy rdzeń więź utrwali, drugie
 * żądanie nie idzie wcale.
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

/** Pasek wcielenia — jedyny wołacz `role.update` w kliencie. */
export interface PasekWcielenia {
  element: HTMLElement;
  /** Przerysowuje blokadę po zmianie obsady. */
  odswiez(): void;
}

/**
 * Pasek nadania wcielenia roli koordynatora (`role.update`).
 *
 * Wcielenie jest wolnym napisem, bo kontrakt nie definiuje katalogu wcieleń
 * (komentarz pola `persona`) — pasek nie buduje listy do wyboru, tylko przyjmuje
 * napis i powtarza to, co rdzeń oddał. Bez koordynatora przycisk jest zablokowany
 * i niesie powód blokady.
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
        // Zdanie mówi, co oddał rdzeń: przyjęcie wywołania z innym wcieleniem
        // (albo bez wcielenia) nie jest nadaniem tego, o które proszono.
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
