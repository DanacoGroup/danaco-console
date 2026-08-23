/**
 * Wykaz mikrofonów widocznych dla przeglądarki.
 *
 * Jedna odpowiedzialność: odczyt `navigator.mediaDevices.enumerateDevices()`
 * i podanie samych wejść dźwięku w kształcie nadającym się na listę wyboru.
 * Plik nie prosi o zgodę, nie otwiera strumienia i nie tworzy DOM-u — pytanie
 * „czym nagrywać” jest osobne od pytania „nagrywaj”.
 */

/** Pojedyncze wejście dźwięku widziane przez przeglądarkę. */
export interface UrzadzenieDzwieku {
  /** `deviceId` z przeglądarki — to samo, co przyjmuje `getUserMedia`. */
  id: string;
  /** Nazwa na listę wyboru; nigdy pusta (patrz `NAZWA_ZASTEPCZA`). */
  nazwa: string;
  /** Wejście, którego przeglądarka użyje bez wskazania wprost. */
  domyslne: boolean;
}

/**
 * Wynik odczytu wykazu.
 *
 * Pusty wykaz nie tłumaczy się sam. Lista bez pozycji znaczy „nie ma
 * mikrofonu” albo „przeglądarka nie chce ich pokazać” — dla Operatora to dwie
 * różne sytuacje z dwoma różnymi wyjściami. Dlatego `powod` idzie obok wykazu
 * i bywa niepusty także wtedy, gdy urządzenia są (bo mają nazwy zastępcze).
 * Pusty `powod` znaczy: wykaz jest kompletny i nie ma czego dopowiadać.
 */
export interface WykazUrzadzen {
  urzadzenia: readonly UrzadzenieDzwieku[];
  powod: string;
}

/** Przedrostek nazwy zastępczej dla urządzeń o ukrytej etykiecie. */
const NAZWA_ZASTEPCZA = 'Mikrofon';

/**
 * Zdanie na brak `navigator.mediaDevices` — trzyczęściowe: co (nie ma wykazu),
 * dlaczego (przeglądarka nie udostępnia dostępu do sprzętu) i czym Operator to
 * zmieni (adres `https` albo nowsze okno).
 */
const POWOD_BRAK_DOSTEPU =
  'Wykaz mikrofonów jest niedostępny. Ta przeglądarka nie udostępnia dostępu do ' +
  'sprzętu dźwiękowego — dzieje się tak w starszych oknach osadzonych i pod adresem ' +
  'bez szyfrowania. Otwórz konsolę pod adresem „https” albo w nowszym oknie przeglądarki.';

/**
 * Zdanie na ukryte etykiety. Też trzyczęściowe: nazwy są zastępcze, bo
 * przeglądarka ukrywa je do pierwszej zgody, a Operator odsłoni je nagrywając
 * raz i przyznając dostęp do mikrofonu.
 */
const POWOD_UKRYTE_NAZWY =
  'Nazwy mikrofonów są zastępcze. Przeglądarka ukrywa etykiety sprzętu, dopóki nie ' +
  'przyznasz dostępu do mikrofonu. Nagraj raz i zgódź się na dostęp — po tym wykaz ' +
  'pokaże prawdziwe nazwy urządzeń.';

/** Zdanie na wykaz pusty mimo działającego odczytu. */
const POWOD_BRAK_URZADZEN =
  'Nie widać żadnego mikrofonu. System nie zgłasza podłączonego wejścia dźwięku — ' +
  'sprzęt bywa odłączony albo zajęty przez inny program. Wepnij mikrofon lub zamknij ' +
  'program, który go trzyma, a wykaz odświeży się sam.';

/** Dostęp do sprzętu dźwiękowego; `null` gdy przeglądarka go nie ma. */
function sprzet(): MediaDevices | null {
  // Sięgamy przez `globalThis`, a nie przez `navigator` wprost, bo w starym
  // oknie osadzonym `navigator.mediaDevices` bywa `undefined` mimo typu, który
  // obiecuje obiekt. Typ tego nie wyłapie — sprawdzenie musi być na wykonaniu.
  const nawigator = (globalThis as { navigator?: Navigator }).navigator;
  return nawigator?.mediaDevices ?? null;
}

/**
 * Składa wykaz z surowej odpowiedzi przeglądarki.
 *
 * Przed pierwszą zgodą `label` każdego urządzenia jest pustym napisem —
 * przeglądarka nie zdradza, jaki sprzęt stoi przy maszynie. Wykaz pustych
 * wierszy jest gorszy niż brak wykazu, więc pustą etykietę zastępuje numer
 * porządkowy, a `powod` mówi, skąd te nazwy się wzięły. Sprzętu nie zgadujemy:
 * „Mikrofon 2” jest przyznaniem się do niewiedzy, nie nazwą.
 */
function zlozWykaz(wejscia: readonly MediaDeviceInfo[]): WykazUrzadzen {
  if (wejscia.length === 0) {
    return { urzadzenia: [], powod: POWOD_BRAK_URZADZEN };
  }

  // Wpis o identyfikatorze `default` jest umową przeglądarek na „wejście
  // wybrane w systemie”. Gdy go nie ma, domyślnym zostaje pierwszy z wykazu —
  // bo taką kolejnością przeglądarka odpowiada i taką przyjmie `getUserMedia`
  // bez wskazania urządzenia.
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
 * Odczyt wykazu wejść dźwięku.
 *
 * Odmowa odczytu nie wychodzi wyjątkiem — wychodzi pustym wykazem ze zdaniem
 * w `powod`. Wywołujący rysuje listę wyboru i nie ma dokąd rzucić błędu, a
 * zdanie i tak musi trafić przed oczy Operatora.
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
    // Przeglądarki potrafią odmówić samego wyliczenia (polityka uprawnień
    // ramki). Zdanie nazywa i to, i drogę wyjścia — pusty wykaz bez słowa
    // czytałby się jako „nie masz mikrofonu”.
    console.error('[dyktowanie] odczyt urządzeń odmówiony', blad);
    return { urzadzenia: [], powod: POWOD_BRAK_DOSTEPU };
  }
}

/**
 * Subskrypcja zmian sprzętu.
 *
 * Operator wpina słuchawki w połowie dnia i wykaz sprzed godziny przestaje być
 * prawdą. Subskrypcja nie odczytuje wykazu sama — tylko budzi wywołującego,
 * bo to on wie, czy lista wyboru w ogóle stoi na ekranie. Gdy przeglądarka nie
 * ma `mediaDevices`, zwracamy odwołanie, które nic nie robi: brak zdarzeń jest
 * stanem, nie błędem wywołania.
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
