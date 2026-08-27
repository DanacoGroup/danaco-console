import type { SettingDefinition } from '../../../shared/contract';
import {
  odczytajWybor,
  zapiszWybor,
  zastosujMotyw,
  ZDARZENIE_MOTYWU,
  type Motyw,
  type ZmianaMotywu,
} from '../motyw/motyw';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import {
  naZmianeKlucza,
  odczytajNastawe,
  poziomZapisu,
  zapiszNastawe,
} from './zrodlo-nastaw';

/** Most motywu jest jedynym właścicielem nastawy motywu po stronie klienta: zwija dwa miejsca zapisu, katalog rdzenia i pamięć przeglądarki, w jedną nastawę, w której prawdą jest rdzeń. */
export const KLUCZ_MOTYWU = 'personalizacja.motyw';

/** Wartość nastawy motywu: pusta znaczy preferencję systemu, w przeciwnym razie niesie nazwę motywu jasnego albo ciemnego. */
export type WyborMotywu = '' | Motyw;

export interface MostMotywu {
  /** Wartość bieżąca nastawy — to, co ma stanąć na uchwycie steru. */
  wybor(): WyborMotywu;
  /** Definicja katalogu: etykiety i opcje wyboru; `undefined` przed odczytem. */
  definicja(): SettingDefinition | undefined;
  /** Zdanie odmowy z ostatniego odczytu albo zapisu; puste znaczy „bez odmowy". */
  odmowa(): string;
  /** Odczytuje katalog i wartość z rdzenia, po czym stosuje motyw. */
  odczytaj(): Promise<void>;
  /** Zapisuje wybór do rdzenia. Oddaje zdanie odmowy albo puste przy powodzeniu. */
  ustaw(wybor: WyborMotywu): Promise<string>;
  /** Subskrypcja zmiany — dla każdego miejsca pokazującego tę nastawę. */
  naZmiane(sluchacz: (wybor: WyborMotywu) => void): () => void;
  /** Odpina nasłuch kanału i nasłuch dokumentu. */
  rozlacz(): void;
}

/** Most zbudowany dla danego kanału połączenia — jeden na kanał, zgodnie z zasadą opisaną w nagłówku pliku. */
let most: MostMotywu | null = null;
let kanalMostu: Kanal | null = null;

/** Podpina most motywu do kanału i oddaje go; wołanie z tym samym kanałem oddaje most już stojący, bez drugiego odczytu. */
export function podepnijMostMotywu(kanal: Kanal): MostMotywu {
  if (most !== null && kanalMostu === kanal) return most;
  if (most !== null) most.rozlacz();
  most = zbudujMostMotywu(kanal);
  kanalMostu = kanal;
  void most.odczytaj();
  return most;
}

/** Zdejmuje most bieżący — droga wyjścia dla sprawdzianów jednostkowych oraz dla rozłączenia całego okna. */
export function odepnijMostMotywu(): void {
  if (most !== null) most.rozlacz();
  most = null;
  kanalMostu = null;
}

function zbudujMostMotywu(kanal: Kanal): MostMotywu {
  /** Wartość potwierdzona przez rdzeń; przed pierwszym odczytem — zapis lokalny. */
  let wybor: WyborMotywu = odczytajWybor() ?? '';
  let definicja: SettingDefinition | undefined;
  let odmowa = '';
  const sluchacze: Array<(wybor: WyborMotywu) => void> = [];

  // Zapora pętli zwrotnej: bez niej zastosowanie wartości z rdzenia wróciłoby jako zmiana operatora.
  let stosujeZRdzenia = false;

  /** Odczyt w toku — żeby zapis nie wyprzedził poznania poziomu zasięgu. */
  let wToku: Promise<void> | null = null;

  function oglos(): void {
    for (const sluchacz of [...sluchacze]) sluchacz(wybor);
  }

  // Stosuje wartość: atrybut dokumentu, odbicie lokalne, ogłoszenie słuchaczom jedną drogą.
  function zastosuj(nowy: WyborMotywu): void {
    wybor = nowy;
    stosujeZRdzenia = true;
    try {
      zastosujMotyw(nowy === '' ? null : nowy);
    } finally {
      stosujeZRdzenia = false;
    }
    // Odbicie prawdy rdzenia, żeby następne uruchomienie nie migało cudzym motywem.
    zapiszWybor(nowy === '' ? null : nowy);
    oglos();
  }

  // Wartość z rdzenia sprowadzona do trzech stanów; nierozpoznane degraduje się do preferencji systemu.
  function jakoWybor(wartosc: unknown): WyborMotywu {
    return wartosc === 'light' || wartosc === 'dark' ? wartosc : '';
  }

  const odsubskrybuj: Odsubskrybuj = naZmianeKlucza(kanal, KLUCZ_MOTYWU, (wartosc) => {
    // Brak wartości znaczy zerowanie: wpis zdjęty, obowiązuje wartość domyślna katalogu.
    zastosuj(wartosc === undefined ? jakoWybor(definicja?.defaultValue) : jakoWybor(wartosc));
  });

  /** Kliknięcie przełącznika na pasku — drugi ster tej samej nastawy motywu produktu. */
  function odZdarzeniaMotywu(zdarzenie: Event): void {
    if (stosujeZRdzenia) return;
    const zmiana = (zdarzenie as CustomEvent<ZmianaMotywu>).detail;
    const zadany: WyborMotywu = zmiana.wybor ?? '';
    // Zmiana preferencji systemu przy braku wyboru operatora też idzie tym zdarzeniem.
    if (zadany === wybor) return;
    void ustaw(zadany);
  }

  document.addEventListener(ZDARZENIE_MOTYWU, odZdarzeniaMotywu);

  async function odczytaj(): Promise<void> {
    const odczyt = odczytajNastawe(kanal, KLUCZ_MOTYWU).then((wynik) => {
      definicja = wynik.definicja;
      odmowa = wynik.odmowa;
      zastosuj(jakoWybor(wynik.wpis?.value ?? wynik.definicja?.defaultValue));
    });
    wToku = odczyt;
    await odczyt;
    wToku = null;
  }

  async function ustaw(nowy: WyborMotywu): Promise<string> {
    // Poziom zapisu bierze się z definicji katalogu, nigdy z domysłu — rdzeń niedozwolonego nie odrzuca.
    if (definicja === undefined && wToku !== null) await wToku;
    const poziom = poziomZapisu(definicja);
    if (poziom === undefined) {
      odmowa =
        'Nie ma gdzie zapisać motywu: katalog ustawień rdzenia nie oddał definicji ' +
        `klucza „${KLUCZ_MOTYWU}" ani jego dozwolonego poziomu zasięgu.`;
      oglos();
      return odmowa;
    }

    // Wartość stosuje się od razu: operator ma zobaczyć skutek kliknięcia natychmiast, bez czekania.
    zastosuj(nowy);

    const wynik = await zapiszNastawe(kanal, KLUCZ_MOTYWU, nowy, poziom);
    odmowa = wynik.udany ? '' : (wynik.blad?.message ?? 'Rdzeń odmówił zapisu bez podania powodu.');
    if (odmowa !== '') oglos();
    return odmowa;
  }

  return {
    wybor: () => wybor,
    definicja: () => definicja,
    odmowa: () => odmowa,
    odczytaj,
    ustaw,

    naZmiane(sluchacz) {
      sluchacze.push(sluchacz);
      return () => {
        const miejsce = sluchacze.indexOf(sluchacz);
        if (miejsce >= 0) sluchacze.splice(miejsce, 1);
      };
    },

    rozlacz() {
      odsubskrybuj();
      document.removeEventListener(ZDARZENIE_MOTYWU, odZdarzeniaMotywu);
      sluchacze.length = 0;
    },
  };
}
