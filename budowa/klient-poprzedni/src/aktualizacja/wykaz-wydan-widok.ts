import {
  ADRES_WYDAN,
  pobierzWykazWydan,
  type OdczytWykazu,
  type Wydanie,
} from './wykaz-wydan';

/** Siedem kolumn chronologii wydań, w tej samej kolejności co na stronie „Pobierz": wersja, data, system, plik, rozmiar, suma kontrolna i opis zmian. */
const KOLUMNY = [
  'Wersja',
  'Data',
  'System',
  'Plik',
  'Rozmiar',
  'Suma SHA-256',
  'Co się zmieniło',
] as const;

/** Zdanie wypisywane w miejscu sumy kontrolnej wydania bez sumy, brzmiące tak samo jak na stronie „Pobierz" w tabeli chronologii wydań. */
const BRAK_SUMY = 'BRAK — aplikacja takiego wydania NIE ZAŁOŻY';

/** Wynik budowy widoku wykazu wydań: element gotowy do osadzenia w oknie oraz metoda odświeżająca jego treść z kanału. */
export interface WykazWydanWidok {
  /** Element do osadzenia przez przywołującego. Powstaje pusty i milczący. */
  element: HTMLElement;
  /** Pyta kanał i przerysowuje widok tym, co zastał. Nie rzuca. */
  odswiez(): Promise<void>;
}

/** Parametry budowy widoku wykazu wydań: adres kanału, z którego czytana jest chronologia, oraz wersja aktualnie zainstalowana. */
export interface UstawieniaWykazu {
  /** Adres kanału wydań. Domyślnie ten sam, o który pyta baner. */
  adres?: string;
  /** Wersja zainstalowana — wypisywana nad chronologią, gdy podana. */
  wersjaBiezaca?: string;
}

/**
 * Buduje widok chronologii wydań.
 *
 * Odczyt zaczyna się dopiero od `odswiez()`: samo zbudowanie widoku nie rusza
 * w sieć, bo widok bywa zbudowany na zapas i nigdy nieprzywołany.
 */
export function utworzWykazWydanWidok(ustawienia: UstawieniaWykazu = {}): WykazWydanWidok {
  const adres = ustawienia.adres ?? ADRES_WYDAN;
  const wersjaBiezaca = ustawienia.wersjaBiezaca ?? '';

  const element = document.createElement('div');
  element.className = 'da-wykaz';

  /** Ostatnie żądanie wygrywa — starsza odpowiedź nie przerysuje nowszej. */
  let obieg = 0;

  function czysc(): void {
    element.textContent = '';
  }

  function zdanie(tresc: string, klasa = ''): HTMLParagraphElement {
    const akapit = document.createElement('p');
    akapit.className = klasa === '' ? 'da-wykaz__zdanie' : `da-wykaz__zdanie ${klasa}`;
    akapit.textContent = tresc;
    return akapit;
  }

  function naglowekChronologii(): HTMLElement {
    const wiersz = document.createElement('p');
    wiersz.className = 'da-wykaz__zrodlo';
    wiersz.textContent =
      wersjaBiezaca === ''
        ? `Chronologia wydań z kanału ${adres}.`
        : `Chronologia wydań z kanału ${adres}. Pracujesz na wersji ${wersjaBiezaca}.`;
    return wiersz;
  }

  function komorka(tresc: string, klasa = ''): HTMLTableCellElement {
    const pole = document.createElement('td');
    if (klasa !== '') pole.className = klasa;
    pole.textContent = tresc;
    return pole;
  }

  function wiersz(wydanie: Wydanie): HTMLTableRowElement {
    const rzad = document.createElement('tr');

    const wersja = document.createElement('td');
    const mocno = document.createElement('strong');
    mocno.textContent = wydanie.wersja;
    wersja.append(mocno);
    rzad.append(wersja);

    rzad.append(komorka(wydanie.data));
    rzad.append(komorka(wydanie.system ?? '—'));

    // Wpis bez pliku jest pozycją historyczną, a nie przeoczeniem — pole nazywa to wprost.
    const plik = wydanie.plik
      ? komorka(wydanie.nazwaPliku ?? wydanie.plik, 'dn-dane')
      : komorka('— wpis historyczny, bez pliku do pobrania', 'dn-dane');
    if (wydanie.plik) plik.title = wydanie.plik;
    rzad.append(plik);

    rzad.append(komorka(wydanie.rozmiar ?? '—', 'dn-dane'));

    // Brak sumy jest przeszkodą, nie pustym polem — komórka nazywa to tym samym zdaniem co witryna.
    const suma = document.createElement('td');
    suma.className = 'dn-dane';
    const zapis = document.createElement('code');
    zapis.className = wydanie.suma ? 'da-wykaz__suma' : 'da-wykaz__suma da-wykaz__suma--brak';
    zapis.textContent = wydanie.suma ?? BRAK_SUMY;
    suma.append(zapis);
    rzad.append(suma);

    rzad.append(komorka(wydanie.zmiany ?? '—', 'dn-dane'));
    return rzad;
  }

  function chronologia(wydania: Wydanie[]): HTMLElement {
    const przewijak = document.createElement('div');
    przewijak.className = 'da-wykaz__przewijak';

    const tabela = document.createElement('table');
    tabela.className = 'dn-tabela';

    const glowa = document.createElement('thead');
    const rzadGlowy = document.createElement('tr');
    for (const nazwa of KOLUMNY) {
      const naglowek = document.createElement('th');
      naglowek.scope = 'col';
      naglowek.textContent = nazwa;
      rzadGlowy.append(naglowek);
    }
    glowa.append(rzadGlowy);

    const cialo = document.createElement('tbody');
    for (const wydanie of wydania) cialo.append(wiersz(wydanie));

    tabela.append(glowa, cialo);
    przewijak.append(tabela);
    return przewijak;
  }

  /** Zamienia stan kanału na zdanie widoku; „brak wydań" pada wyłącznie dla stanu `pusty`. */
  function rysuj(odczyt: OdczytWykazu): void {
    czysc();
    switch (odczyt.stan) {
      case 'wykaz':
        element.append(naglowekChronologii(), chronologia(odczyt.wydania));
        return;
      case 'pusty':
        element.append(
          zdanie(
            'Kanał wydań odpowiedział i mówi wprost: nie ma jeszcze żadnego wydania. ' +
              'Chronologia pokaże wydanie w tej samej chwili, w której ono powstanie — ' +
              'a do tego czasu mówi o tym wprost, zamiast pokazywać pustą tabelę.',
          ),
        );
        return;
      case 'brak-wydania-pod-adresem':
        element.append(
          zdanie(
            `Pod adresem kanału wydań (${adres}) nie ma wykazu — serwer odpowiedział ` +
              `kodem ${odczyt.status}. To NIE jest kłopot z Twoim łączem: żądanie doszło ` +
              'i doczekało się odpowiedzi. Brakuje samego pliku wykazu po stronie kanału.',
            'da-wykaz__zdanie--odmowa',
          ),
        );
        return;
      case 'odpowiedz-serwera':
        element.append(
          zdanie(
            `Serwer kanału wydań odpowiedział kodem ${odczyt.status}. Kanał stoi, ale ma ` +
              'kłopot po swojej stronie — chronologii nie znam i nie zgaduję jej. ' +
              'Przywołaj wykaz jeszcze raz za jakiś czas.',
            'da-wykaz__zdanie--odmowa',
          ),
        );
        return;
      case 'brak-lacznosci':
        element.append(
          zdanie(
            `Nie udało się zapytać kanału wydań (${adres}) — żądanie nie doszło do nikogo. ` +
              'To brak łączności albo adres, który się nie rozwiązuje. ' +
              'NIE ZNACZY TO, ŻE WYDAŃ NIE MA — znaczy, że nie wiem, jakie są.',
            'da-wykaz__zdanie--odmowa',
          ),
        );
        return;
      case 'odpowiedz-nieczytelna':
        element.append(
          zdanie(
            `Spod adresu ${adres} coś przyszło, ale nie jest to wykaz wydań. ` +
              'Chronologii nie pokazuję, bo nie umiem tej odpowiedzi przeczytać — ' +
              'zmyślona lista byłaby gorsza niż jej brak.',
            'da-wykaz__zdanie--odmowa',
          ),
        );
        return;
      default: {
        // Wyczerpanie stanów pilnuje kompilator — dopisanie siódmego stanu zapali się tutaj.
        const nieznany: never = odczyt;
        return nieznany;
      }
    }
  }

  return {
    element,
    async odswiez(): Promise<void> {
      const moj = (obieg += 1);
      czysc();
      element.append(zdanie('Pytam kanał wydań…'));
      const odczyt = await pobierzWykazWydan(adres);
      // Odpowiedź na żądanie już nieaktualne nie przerysowuje widoku — czeka się na nowsze.
      if (moj !== obieg) return;
      rysuj(odczyt);
    },
  };
}
