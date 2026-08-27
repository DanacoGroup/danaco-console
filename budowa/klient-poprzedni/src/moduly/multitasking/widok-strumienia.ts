import { ChunkKind, type LoopState } from '../../../../shared/contract';
import { RodzajFragmentu } from '../../rozmowa/indeks';
import { wierszOpisu } from './kontrolki';
import type { StanMultitaskingu } from './stan-multitaskingu';

/**
 * Pełny strumień wykonawców widziany przez okno koordynatora: tok rozumowania,
 * wywołania narzędzi, wyniki i pliki wraz z licznikiem obiegów biegu naprawczego.
 * Rodzaje fragmentów pochodzą z kontraktu przez katalog warstwy rozmowy.
 */
export interface WidokStrumienia {
  element: HTMLElement;
  odswiez(): void;
}

/**
 * Nazwy rodzajów fragmentu widoczne w nagłówku wpisu. Odwzorowanie sprowadza
 * wartości kontraktu na określenia polszczyzny zawodowej, a rodzaj nieujęty
 * w odwzorowaniu trafia do nagłówka w postaci surowej.
 */
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

/**
 * Ile ostatnich wpisów strumienia rysuje widok koordynatora. Ograniczenie
 * powstrzymuje rozrost drzewa dokumentu przy długim biegu, a wpisy starsze
 * pozostają zachowane w stanie modułu.
 */
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

/**
 * Rysuje licznik obiegów biegu naprawczego odczytany komendą `window.state.get`.
 * Brak biegu nie jest usterką i oddaje wiersz o biegu nierozpoczętym, a bieg
 * zatrzymany niesie powód zatrzymania z pola `LoopState.stopReason`.
 */
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

/**
 * Składa wpisy strumienia obu wykonawców w jeden ciąg uporządkowany chwilą
 * powstania i rysuje ostatnie z nich. Każdy wpis niesie znacznik okna
 * źródłowego, więc pochodzenie fragmentu pozostaje czytelne po złączeniu.
 */
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
    // Pliki wyróżnione osobno, obok toku rozumowania, wywołań i wyników.
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

/**
 * Rozstrzyga, czy rodzaj fragmentu niesie plik. Plikiem jest obraz albo dźwięk
 * zwrócony przez narzędzie, a rozpoznanie opiera się wyłącznie na wartościach
 * wyliczenia `ChunkKind` pochodzącego z kontraktu.
 */
function czyPlik(rodzaj: ChunkKind): boolean {
  return rodzaj === ChunkKind.Image || rodzaj === ChunkKind.Audio;
}
