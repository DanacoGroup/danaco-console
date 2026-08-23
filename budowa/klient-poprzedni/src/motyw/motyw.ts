// ============================================================================
// DANACO CONSOLE — MOTYW · ODCZYT, ZAPIS I PRZEŁĄCZANIE
// ----------------------------------------------------------------------------
// Jedna odpowiedzialność: wybór motywu użytkownika — odczyt trwałego zapisu,
// zapis, ustawienie atrybutu data-theme oraz śledzenie preferencji systemu.
// Wartości barw nie występują w tym pliku; należą do arkuszy motyw.css.
//
// Zasady:
//   · brak zapisanego wyboru nie ustawia atrybutu — rozstrzyga prefers-color-scheme,
//     a zmiana preferencji systemu działa na żywo;
//   · żaden błąd pamięci trwałej nie zatrzymuje uruchomienia —
//     wybór degraduje się do preferencji systemu, nigdy do blokady;
//   · oba motywy są równoprawne, żaden nie jest wartością „domyślną" produktu.
// ============================================================================

export type Motyw = 'light' | 'dark';

/** Klucz zapisu wyboru w pamięci trwałej przeglądarki. */
export const KLUCZ_ZAPISU = 'danaco-console.motyw';

/** Nazwa zdarzenia rozgłaszanego po każdej zmianie obowiązującego motywu. */
export const ZDARZENIE_MOTYWU = 'danaco-motyw';

/** Treść zdarzenia zmiany motywu. */
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

/** Pamięć trwała albo null, gdy środowisko jej nie udostępnia. */
function pamiec(): Storage | null {
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

/** Odczytuje zapisany wybór użytkownika. Brak wyboru albo zapis nieczytelny → null. */
export function odczytajWybor(): Motyw | null {
  try {
    const zapis = pamiec()?.getItem(KLUCZ_ZAPISU);
    return czyMotyw(zapis) ? zapis : null;
  } catch {
    return null;
  }
}

/** Zapisuje wybór użytkownika; null kasuje zapis i oddaje decyzję systemowi. */
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

/** Preferencja systemu operacyjnego. Bez obsługi matchMedia przyjmujemy jasny. */
export function preferencjaSystemu(): Motyw {
  if (typeof window.matchMedia !== 'function') return 'light';
  return window.matchMedia(ZAPYTANIE_CIEMNY).matches ? 'dark' : 'light';
}

/** Motyw faktycznie obowiązujący: wybór użytkownika, a w jego braku system. */
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

/** Ustawia motyw i zapisuje go jako wybór użytkownika. */
export function ustawMotyw(wybor: Motyw): Motyw {
  zapiszWybor(wybor);
  return zastosujMotyw(wybor);
}

/** Kasuje wybór użytkownika i wraca do preferencji systemu. */
export function przywrocPreferencjeSystemu(): Motyw {
  zapiszWybor(null);
  return zastosujMotyw(null);
}

/** Przełącza motyw na przeciwny względem obowiązującego i zapisuje wybór. */
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
