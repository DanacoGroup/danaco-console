/** Motyw obsługuje wybór motywu użytkownika: odczyt trwałego zapisu, jego zapis, ustawienie atrybutu koloru na dokumencie oraz śledzenie preferencji systemu, przy czym same wartości barw pozostają w arkuszu stylów. */

export type Motyw = 'light' | 'dark';

/** Klucz zapisu, pod którym wybór motywu użytkownika jest przechowywany w pamięci trwałej przeglądarki. */
export const KLUCZ_ZAPISU = 'danaco-console.motyw';

/** Nazwa zdarzenia rozgłaszanego w oknie po każdej zmianie motywu faktycznie obowiązującego w dokumencie. */
export const ZDARZENIE_MOTYWU = 'danaco-motyw';

/** Treść zdarzenia zmiany motywu, niosąca zarówno motyw obowiązujący, jak i wybór dokonany przez użytkownika. */
export interface ZmianaMotywu {
  /** Motyw faktycznie obowiązujący po zmianie. */
  obowiazujacy: Motyw;
  /** Wybór użytkownika; null oznacza zdanie się na preferencję systemu. */
  wybor: Motyw | null;
}

const ZAPYTANIE_CIEMNY = '(prefers-color-scheme: dark)';

function czyMotyw(wartosc: unknown): wartosc is Motyw {
  return wartosc === 'light' || wartosc === 'dark';
}

/** Pamięć trwała przeglądarki użyta do zapisu wyboru motywu, albo null, gdy środowisko jej nie udostępnia. */
function pamiec(): Storage | null {
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

/** Odczytuje zapisany wybór użytkownika z pamięci trwałej; brak wyboru albo zapis nieczytelny dają wartość null. */
export function odczytajWybor(): Motyw | null {
  try {
    const zapis = pamiec()?.getItem(KLUCZ_ZAPISU);
    return czyMotyw(zapis) ? zapis : null;
  } catch {
    return null;
  }
}

/** Zapisuje wybór motywu dokonany przez użytkownika; wartość null kasuje zapis i oddaje decyzję preferencji systemu. */
export function zapiszWybor(wybor: Motyw | null): void {
  try {
    const magazyn = pamiec();
    if (magazyn === null) return;
    if (wybor === null) magazyn.removeItem(KLUCZ_ZAPISU);
    else magazyn.setItem(KLUCZ_ZAPISU, wybor);
  } catch {
    // Brak zapisu nie odbiera możliwości przełączenia motywu w bieżącej sesji.
  }
}

/** Preferencja systemu operacyjnego odczytana z zapytania o media; bez obsługi tego zapytania przyjmuje się motyw jasny. */
export function preferencjaSystemu(): Motyw {
  if (typeof window.matchMedia !== 'function') return 'light';
  return window.matchMedia(ZAPYTANIE_CIEMNY).matches ? 'dark' : 'light';
}

/** Motyw faktycznie obowiązujący w dokumencie: wybór użytkownika, a w jego braku preferencja systemu operacyjnego. */
export function motywObowiazujacy(): Motyw {
  const jawny = document.documentElement.getAttribute('data-theme');
  if (czyMotyw(jawny)) return jawny;
  return odczytajWybor() ?? preferencjaSystemu();
}

function rozglos(wybor: Motyw | null, obowiazujacy: Motyw): void {
  const tresc: ZmianaMotywu = { obowiazujacy, wybor };
  document.dispatchEvent(new CustomEvent<ZmianaMotywu>(ZDARZENIE_MOTYWU, { detail: tresc }));
}

/**
 * Ustawia atrybut data-theme na elemencie głównym dokumentu.
 * Wartość null usuwa atrybut — od tej chwili rozstrzyga preferencja systemu.
 * Zwraca motyw faktycznie obowiązujący po zmianie.
 */
export function zastosujMotyw(wybor: Motyw | null): Motyw {
  const korzen = document.documentElement;
  if (wybor === null) korzen.removeAttribute('data-theme');
  else korzen.setAttribute('data-theme', wybor);

  const obowiazujacy = wybor ?? preferencjaSystemu();
  rozglos(wybor, obowiazujacy);
  return obowiazujacy;
}

/** Ustawia motyw dokumentu i zapisuje go jednocześnie jako świadomy wybór użytkownika w pamięci trwałej. */
export function ustawMotyw(wybor: Motyw): Motyw {
  zapiszWybor(wybor);
  return zastosujMotyw(wybor);
}

/** Kasuje zapisany wybór użytkownika i przywraca dokument do stanu wyznaczanego przez preferencję systemu. */
export function przywrocPreferencjeSystemu(): Motyw {
  zapiszWybor(null);
  return zastosujMotyw(null);
}

/** Przełącza motyw dokumentu na przeciwny względem obecnie obowiązującego i zapisuje ten wybór jako decyzję użytkownika. */
export function przelaczMotyw(): Motyw {
  return ustawMotyw(motywObowiazujacy() === 'dark' ? 'light' : 'dark');
}

/**
 * Śledzi zmianę preferencji systemu. Dopóki użytkownik nie dokonał wyboru,
 * zmiana ustawienia systemowego rozgłaszana jest dalej — interfejs nadąża
 * za systemem bez przeładowania. Zwraca funkcję odpinającą nasłuch.
 */
export function sledzPreferencjeSystemu(): () => void {
  if (typeof window.matchMedia !== 'function') return () => {};

  const zapytanie = window.matchMedia(ZAPYTANIE_CIEMNY);
  const obsluga = (): void => {
    if (odczytajWybor() !== null) return;
    rozglos(null, preferencjaSystemu());
  };

  zapytanie.addEventListener('change', obsluga);
  return () => zapytanie.removeEventListener('change', obsluga);
}

/**
 * Uruchomienie warstwy motywu. Odtwarza zapisany wybór, a w jego braku
 * pozostawia dokument bez atrybutu, aby rozstrzygnęła preferencja systemu.
 * Zwraca motyw obowiązujący po uruchomieniu.
 */
export function uruchomMotyw(): Motyw {
  const obowiazujacy = zastosujMotyw(odczytajWybor());
  sledzPreferencjeSystemu();
  return obowiazujacy;
}
