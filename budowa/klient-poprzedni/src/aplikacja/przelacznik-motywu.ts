import { elementIkony } from '../ikony/ikony';
import type { NazwaIkony } from '../ikony/ikony';
import {
  motywObowiazujacy,
  odczytajWybor,
  przelaczMotyw,
  przywrocPreferencjeSystemu,
  ZDARZENIE_MOTYWU,
  type Motyw,
  type ZmianaMotywu,
} from '../motyw/motyw';

/**
 * Przełącznik motywu na pasku górnym. Niesie dwie kontrolki: przełączenie między
 * motywem jasnym i ciemnym oraz powrót do preferencji systemu.
 */
export interface PrzelacznikMotywu {
  /** Grupa kontrolek montowana w akcjach paska. */
  element: HTMLElement;
  /** Doprowadza wygląd kontrolek do stanu bieżącego. */
  odswiez(): void;
}

/**
 * Przełączenie motywu jasnego i ciemnego oraz powrót do preferencji systemu.
 * Oba motywy są równoprawne, kontrolka nie wskazuje żadnego jako domyślnego.
 * Wartości motywu ani żadnej barwy ten plik nie zna, zostaje w nim sama obsługa
 * kontrolek.
 */
export function utworzPrzelacznikMotywu(): PrzelacznikMotywu {
  const element = document.createElement('div');
  element.className = 'dn-powloka__akcje';

  const przelacz = przyciskIkonowy();
  const system = przyciskIkonowy();

  element.append(przelacz, system);

  function odswiez(): void {
    const obowiazujacy = motywObowiazujacy();
    ubierzPrzelacznik(przelacz, obowiazujacy);
    ubierzPowrot(system, odczytajWybor() === null);
  }

  przelacz.addEventListener('click', () => {
    przelaczMotyw();
    odswiez();
  });

  system.addEventListener('click', () => {
    przywrocPreferencjeSystemu();
    odswiez();
  });

  // Zmiana preferencji systemu dochodzi tą samą drogą co przełączenie ręczne,
  // bez przeładowania.
  document.addEventListener(ZDARZENIE_MOTYWU, (zdarzenie: Event) => {
    const zmiana = (zdarzenie as CustomEvent<ZmianaMotywu>).detail;
    ubierzPrzelacznik(przelacz, zmiana.obowiazujacy);
    ubierzPowrot(system, zmiana.wybor === null);
  });

  odswiez();

  return { element, odswiez };
}

/**
 * Pusty przycisk ikonowy w wariancie przeznaczonym na ramę kokpitu. Ikonę i opis
 * dokładają funkcje ubierające, więc przycisk powstaje raz i nie jest wymieniany.
 */
function przyciskIkonowy(): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn-ikona dn-btn-ikona--na-ramie';
  return przycisk;
}

/**
 * Przełącznik pokazuje motyw, w który przejdzie po naciśnięciu — słońce przy
 * motywie ciemnym, księżyc przy jasnym.
 */
function ubierzPrzelacznik(przycisk: HTMLButtonElement, obowiazujacy: Motyw): void {
  const ciemny = obowiazujacy === 'dark';
  const opis = ciemny ? 'Włącz motyw jasny' : 'Włącz motyw ciemny';
  ubierz(przycisk, ciemny ? 'slonce' : 'ksiezyc', opis);
  przycisk.dataset.motyw = obowiazujacy;
}

/**
 * Powrót do preferencji systemu; stan wciśnięcia mówi, czy już obowiązuje.
 * Wciśnięty znaczy brak wyboru własnego, czyli motyw idący za ustawieniem
 * systemu.
 */
function ubierzPowrot(przycisk: HTMLButtonElement, wedlugSystemu: boolean): void {
  const opis = wedlugSystemu
    ? 'Motyw zgodny z systemem'
    : 'Wróć do motywu zgodnego z systemem';
  ubierz(przycisk, 'odswiez', opis);
  przycisk.setAttribute('aria-pressed', String(wedlugSystemu));
}

/**
 * Wymienia ikonę i opis przycisku, zachowując jego tożsamość w dokumencie.
 * Podmiana potomków zostawia nasłuch zdarzeń nietknięty, więc przycisk działa
 * dalej.
 */
function ubierz(przycisk: HTMLButtonElement, ikona: NazwaIkony, opis: string): void {
  przycisk.replaceChildren(elementIkony(ikona, { rozmiar: 18 }));
  przycisk.setAttribute('aria-label', opis);
  przycisk.title = opis;
}
