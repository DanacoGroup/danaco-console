import { ChunkKind, type LoopState } from '../../../../shared/contract';
import { RodzajFragmentu } from '../../rozmowa/indeks';
import { wierszOpisu } from './kontrolki';
import type { StanMultitaskingu } from './stan-multitaskingu';

/**
 * Pełny strumień wykonawcy widziany przez koordynatora wraz
 * z licznikiem obiegów biegu naprawczego.
 *
 * Coordinator Chat widzi cztery rzeczy, nie samą odpowiedź: tok rozumowania,
 * wywołania narzędzi, wyniki i pliki. Rodzaje fragmentów pochodzą z kontraktu
 * przez katalog warstwy rozmowy (`RodzajFragmentu`), więc moduł nie przepisuje
 * ich po swojemu.
 *
 * Gdy bieg stoi, wiersz licznika niesie powód zatrzymania wprost
 * z `LoopState.stopReason` — zatrzymanie nie wygląda wtedy jak cisza.
 */
export interface WidokStrumienia {
  element: HTMLElement;
  odswiez(): void;
}

/** Nazwy rodzajów fragmentu widoczne w nagłówku wpisu. */
const NAZWY_RODZAJOW: Readonly<Record<string, string>> = {
  [RodzajFragmentu.Tekst]: 'wynik',
  [RodzajFragmentu.Rozumowanie]: 'tok rozumowania',
  [RodzajFragmentu.WywolanieNarzedzia]: 'wywołanie narzędzia',
  [RodzajFragmentu.WynikNarzedzia]: 'wynik narzędzia',
  [RodzajFragmentu.Obraz]: 'plik obrazowy',
  [RodzajFragmentu.Dzwiek]: 'plik dźwiękowy',
  [RodzajFragmentu.Blad]: 'błąd tury',
  [RodzajFragmentu.Prowenancja]: 'prowenancja wywołania',
  [RodzajFragmentu.Konto]: 'konto kanału',
};

/** Ile ostatnich wpisów strumienia rysuje widok koordynatora. */
const WIDOCZNYCH = 40;

export function utworzWidokStrumienia(stan: StanMultitaskingu): WidokStrumienia {
  const licznik = document.createElement('div');
  licznik.className = 'dm-licznik';
  licznik.setAttribute('aria-label', 'Licznik obiegów biegu naprawczego');

  const strumien = document.createElement('div');
  strumien.className = 'dm-strumien';
  strumien.setAttribute('aria-label', 'Strumień wykonawców');

  const element = document.createElement('div');
  element.className = 'dm-podglad';
  element.append(licznik, strumien);

  function odswiez(): void {
    rysujLicznik(licznik, stan.bieg());
    rysujStrumien(strumien, stan);
  }

  odswiez();
  return { element, odswiez };
}

/** Licznik obiegów z `window.state.get`; brak biegu nie jest usterką. */
function rysujLicznik(gospodarz: HTMLElement, bieg: LoopState | null): void {
  gospodarz.replaceChildren();
  if (bieg === null) {
    gospodarz.append(
      wierszOpisu('Bieg naprawczy', 'nie rozpoczęty — rdzeń nie zliczył ani jednego obiegu'),
    );
    return;
  }
  gospodarz.dataset['zatrzymany'] = String(bieg.stopped);
  gospodarz.append(
    wierszOpisu('Obiegów', String(bieg.loops)),
    wierszOpisu('Bez postępu', `${bieg.loopsWithoutProgress} z ${bieg.threshold}`),
    wierszOpisu('Stan biegu', bieg.stopped ? 'zatrzymany' : 'w toku'),
  );
  if (bieg.stopped) {
    gospodarz.append(wierszOpisu('Powód zatrzymania', bieg.stopReason ?? 'rdzeń nie podał powodu'));
  }
  if (bieg.lastExecutorWindowId !== undefined) {
    gospodarz.append(
      wierszOpisu('Ostatnia tura', `${bieg.lastExecutorWindowId} — ${bieg.lastTurnReason ?? 'bez powodu'}`),
    );
  }
}

/** Wpisy strumienia obu wykonawców złożone w jeden ciąg czasu. */
function rysujStrumien(gospodarz: HTMLElement, stan: StanMultitaskingu): void {
  const wykonawcy = stan.obsada().wykonawcy;
  const wpisy = wykonawcy
    .flatMap((okno) => stan.strumien().wpisy(okno.id))
    .sort((a, b) => a.chwila - b.chwila)
    .slice(-WIDOCZNYCH);

  gospodarz.replaceChildren();
  if (wpisy.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dm-stan';
    pusto.dataset['stan'] = 'pusto';
    pusto.textContent =
      wykonawcy.length === 0
        ? 'Obsada bez wykonawcy — koordynator nie ma czyjego strumienia oglądać.'
        : 'Wykonawcy nie przysłali jeszcze ani jednego fragmentu.';
    gospodarz.append(pusto);
    return;
  }

  for (const wpis of wpisy) {
    const element = document.createElement('div');
    element.className = 'dm-fragment';
    element.dataset['rodzaj'] = wpis.rodzaj;
    element.dataset['okno'] = wpis.okno;
    // Pliki wyróżnione osobno — stoją obok toku rozumowania, wywołań i wyników
    // jako czwarta rzecz widoczna dla koordynatora.
    element.dataset['plik'] = String(czyPlik(wpis.rodzaj));

    const naglowek = document.createElement('span');
    naglowek.className = 'dm-fragment__rodzaj';
    naglowek.textContent = NAZWY_RODZAJOW[wpis.rodzaj] ?? wpis.rodzaj;

    const tresc = document.createElement('pre');
    tresc.className = 'dm-fragment__tresc';
    tresc.textContent = wpis.tresc;

    element.append(naglowek, tresc);
    gospodarz.append(element);
  }
}

/** Czy rodzaj fragmentu niesie plik — obraz albo dźwięk zwrócony przez narzędzie. */
function czyPlik(rodzaj: ChunkKind): boolean {
  return rodzaj === ChunkKind.Image || rodzaj === ChunkKind.Audio;
}
