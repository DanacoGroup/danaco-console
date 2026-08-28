/**
 * Wykaz mikrofonów widocznych dla przeglądarki odczytuje `navigator.mediaDevices.enumerateDevices()` i podaje same wejścia dźwięku, bez proszenia o zgodę i bez otwierania strumienia.
 */

/** Pojedyncze wejście dźwięku widziane przez przeglądarkę, niosące identyfikator urządzenia oraz nazwę gotową na listę wyboru. */
export interface UrzadzenieDzwieku {
  /** `deviceId` z przeglądarki — to samo, co przyjmuje `getUserMedia`. */
  id: string;
  /** Nazwa na listę wyboru — nigdy pusta, w razie braku zastąpiona przedrostkiem `NAZWA_ZASTEPCZA`. */
  nazwa: string;
  /** Wejście, którego przeglądarka użyje bez wskazania wprost. */
  domyslne: boolean;
}

/**
 * Wynik odczytu wykazu niesie obok listy urządzeń pole `powod`, rozróżniające brak mikrofonu od odmowy przeglądarki i wyjaśniające nazwy zastępcze.
 */
export interface WykazUrzadzen {
  urzadzenia: readonly UrzadzenieDzwieku[];
  powod: string;
}

/** Przedrostek nazwy zastępczej stosowanej dla urządzeń dźwiękowych o ukrytej etykiecie przed pierwszą udzieloną zgodą. */
const NAZWA_ZASTEPCZA = 'Mikrofon';

/**
 * Zdanie na brak `navigator.mediaDevices` wskazuje, że wykazu nie ma, że przeglądarka nie udostępnia dostępu do sprzętu, oraz jak to zmienić.
 */
const POWOD_BRAK_DOSTEPU =
  'Wykaz mikrofonów jest niedostępny. Ta przeglądarka nie udostępnia dostępu do ' +
  'sprzętu dźwiękowego — dzieje się tak w starszych oknach osadzonych i pod adresem ' +
  'bez szyfrowania. Otwórz konsolę pod adresem „https” albo w nowszym oknie przeglądarki.';

/**
 * Zdanie na ukryte etykiety wyjaśnia, że nazwy urządzeń są zastępcze, dopóki przeglądarka nie otrzyma pierwszej zgody na nagrywanie.
 */
const POWOD_UKRYTE_NAZWY =
  'Nazwy mikrofonów są zastępcze. Przeglądarka ukrywa etykiety sprzętu, dopóki nie ' +
  'przyznasz dostępu do mikrofonu. Nagraj raz i zgódź się na dostęp — po tym wykaz ' +
  'pokaże prawdziwe nazwy urządzeń.';

/** Zdanie zwracane, gdy wykaz wejść dźwięku jest pusty mimo poprawnie wykonanego odczytu przez przeglądarkę. */
const POWOD_BRAK_URZADZEN =
  'Nie widać żadnego mikrofonu. System nie zgłasza podłączonego wejścia dźwięku — ' +
  'sprzęt bywa odłączony albo zajęty przez inny program. Wepnij mikrofon lub zamknij ' +
  'program, który go trzyma, a wykaz odświeży się sam.';

/** Dostęp do interfejsu sprzętu dźwiękowego przeglądarki, równy `null`, gdy przeglądarka go nie udostępnia. */
function sprzet(): MediaDevices | null {
  // Dostęp sięga przez `globalThis`, nie `navigator`, bo `mediaDevices` bywa `undefined` w starym oknie.
  const nawigator = (globalThis as { navigator?: Navigator }).navigator;
  return nawigator?.mediaDevices ?? null;
}

/**
 * Składanie wykazu z surowej odpowiedzi przeglądarki zastępuje pustą etykietę numerem porządkowym, gdy zgoda na mikrofon nie została jeszcze udzielona.
 */
function zlozWykaz(wejscia: readonly MediaDeviceInfo[]): WykazUrzadzen {
  if (wejscia.length === 0) {
    return { urzadzenia: [], powod: POWOD_BRAK_URZADZEN };
  }

  // Wpis `default` to umowa na wejście systemowe — bez niego domyślnym zostaje pierwszy wpis wykazu.
  const maWpisDomyslny = wejscia.some((wejscie) => wejscie.deviceId === 'default');
  let ukryte = 0;

  const urzadzenia = wejscia.map((wejscie, numer) => {
    const etykieta = (wejscie.label ?? '').trim();
    if (etykieta === '') ukryte += 1;
    return {
      id: wejscie.deviceId,
      nazwa: etykieta === '' ? `${NAZWA_ZASTEPCZA} ${numer + 1}` : etykieta,
      domyslne: maWpisDomyslny ? wejscie.deviceId === 'default' : numer === 0,
    };
  });

  return { urzadzenia, powod: ukryte > 0 ? POWOD_UKRYTE_NAZWY : '' };
}

/**
 * Odczyt wykazu wejść dźwięku nie kończy się wyjątkiem przy odmowie — zwraca pusty wykaz wraz ze zdaniem wyjaśniającym w polu `powod`.
 */
export async function odczytajUrzadzenia(): Promise<WykazUrzadzen> {
  const urzadzeniaMediow = sprzet();
  if (urzadzeniaMediow === null || typeof urzadzeniaMediow.enumerateDevices !== 'function') {
    return { urzadzenia: [], powod: POWOD_BRAK_DOSTEPU };
  }

  try {
    const wszystkie = await urzadzeniaMediow.enumerateDevices();
    return zlozWykaz(wszystkie.filter((wejscie) => wejscie.kind === 'audioinput'));
  } catch (blad) {
    // Odmowa wyliczenia urządzeń przez politykę ramki zwraca pusty wykaz ze zdaniem wyjaśniającym.
    console.error('[dyktowanie] odczyt urządzeń odmówiony', blad);
    return { urzadzenia: [], powod: POWOD_BRAK_DOSTEPU };
  }
}

/**
 * Subskrypcja zmian sprzętu budzi wywołującego przy zmianie urządzeń audio, nie odczytując wykazu samodzielnie, a przy braku wsparcia zwraca odwołanie bez działania.
 */
export function naZmianeUrzadzen(sluchacz: () => void): () => void {
  const urzadzeniaMediow = sprzet();
  if (urzadzeniaMediow === null || typeof urzadzeniaMediow.addEventListener !== 'function') {
    return () => {};
  }

  const reakcja = (): void => sluchacz();
  urzadzeniaMediow.addEventListener('devicechange', reakcja);
  return () => urzadzeniaMediow.removeEventListener('devicechange', reakcja);
}
