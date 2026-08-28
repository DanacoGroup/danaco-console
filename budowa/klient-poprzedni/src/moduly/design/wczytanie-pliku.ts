import { DesignAssetKind } from '../../../../shared/contract';

/**
 * Odczytanie pliku wskazanego do wgrania prowadzi jedną drogą wskazanie z pola
 * pliku oraz upuszczenie na płytę: oba kończą się obiektem `File` i tym samym
 * żądaniem. Treść czyta klient, ponieważ rdzeń stoi na innej maszynie.
 */
export interface WczytanyPlik {
  /** Nazwa pliku z dysku — trafia do pola `name` żądania. */
  nazwa: string;
  /** Rodzaj zasobu rozpoznany po typie MIME pliku. */
  rodzaj: DesignAssetKind;
  /** Treść pliku w zapisie base64, BEZ przedrostka `data:`. */
  trescBase64: string;
  /** Format pliku w postaci rozszerzenia; pusty, gdy nie da się go ustalić. */
  format: string;
  /** Szerokość w pikselach; zero znaczy „nie zmierzono". */
  szerokosc: number;
  /** Wysokość w pikselach; zero znaczy „nie zmierzono". */
  wysokosc: number;
  /** Rozmiar pliku w bajtach — do zdania o wgraniu, nie do żądania. */
  bajtow: number;
}

/**
 * Rodzaj zasobu z typu MIME pliku: SVG jest wektorem, reszta obrazów rastrem.
 * Trzecia wartość wyliczenia, `Composition`, powstaje w rdzeniu z kompozycji
 * Design Board i plikiem nigdy nie jest, więc nie ma jej w tym rozpoznaniu.
 */
export function rodzajZTypu(typMime: string, nazwa: string): DesignAssetKind {
  if (typMime === 'image/svg+xml' || nazwa.toLowerCase().endsWith('.svg')) {
    return DesignAssetKind.Vector;
  }
  return DesignAssetKind.Image;
}

/**
 * Format pliku ustala się z rozszerzenia nazwy, a gdy nazwa rozszerzenia nie
 * niesie — z podtypu typu MIME. Format pozostaje pusty, gdy żadne z tych dwóch
 * źródeł go nie podaje.
 */
export function formatPliku(typMime: string, nazwa: string): string {
  const kropka = nazwa.lastIndexOf('.');
  if (kropka > 0 && kropka < nazwa.length - 1) return nazwa.slice(kropka + 1).toLowerCase();
  const ukosnik = typMime.indexOf('/');
  return ukosnik === -1 ? '' : typMime.slice(ukosnik + 1).toLowerCase();
}

/**
 * Czyta plik w całości i oddaje wszystko, czego potrzebuje żądanie wgrania.
 * Obietnica jest odrzucana wyłącznie wtedy, gdy przeglądarka nie oddała treści
 * pliku; nieudany pomiar wymiarów jej nie przerywa.
 */
export async function wczytajPlik(plik: File): Promise<WczytanyPlik> {
  const trescBase64 = await odczytajBase64(plik);
  const rodzaj = rodzajZTypu(plik.type, plik.name);
  const wymiary =
    rodzaj === DesignAssetKind.Image
      ? await zmierzWymiary(plik.type, trescBase64)
      : { szerokosc: 0, wysokosc: 0 };
  return {
    nazwa: plik.name,
    rodzaj,
    trescBase64,
    format: formatPliku(plik.type, plik.name),
    szerokosc: wymiary.szerokosc,
    wysokosc: wymiary.wysokosc,
    bajtow: plik.size,
  };
}

/**
 * Bajty pliku w zapisie base64 odczytuje `readAsDataURL`, ponieważ oddaje ten
 * zapis bez ręcznego przepisywania bajtów. Przedrostek `data:…;base64,` zostaje
 * zdjęty, ponieważ kontrakt oczekuje samego zapisu, a nie adresu URI.
 */
function odczytajBase64(plik: File): Promise<string> {
  return new Promise((rozstrzygnij, odrzuc) => {
    const czytnik = new FileReader();
    czytnik.addEventListener('load', () => {
      const wynik = typeof czytnik.result === 'string' ? czytnik.result : '';
      const przecinek = wynik.indexOf(',');
      if (przecinek === -1) {
        odrzuc(new Error('Przeglądarka nie oddała treści pliku w zapisie base64.'));
        return;
      }
      rozstrzygnij(wynik.slice(przecinek + 1));
    });
    czytnik.addEventListener('error', () =>
      odrzuc(new Error('Przeglądarka nie zdołała odczytać wskazanego pliku.')),
    );
    czytnik.readAsDataURL(plik);
  });
}

/**
 * Wymiary rastra mierzy się zdekodowaniem treści w przeglądarce, a wynikiem są
 * zera, gdy przeglądarka obrazu nie zdekodowała. Oba pola wymiaru są w kontrakcie
 * opcjonalne, więc zasób bez nich pozostaje poprawny.
 */
function zmierzWymiary(
  typMime: string,
  trescBase64: string,
): Promise<{ szerokosc: number; wysokosc: number }> {
  return new Promise((rozstrzygnij) => {
    const brak = { szerokosc: 0, wysokosc: 0 };
    // Środowisko bez dekodera obrazów nie jest usterką wgrania: pomiar odpada,
    // treść leci dalej.
    if (typeof Image !== 'function') {
      rozstrzygnij(brak);
      return;
    }
    const obraz = new Image();
    obraz.addEventListener('load', () =>
      rozstrzygnij({ szerokosc: obraz.naturalWidth, wysokosc: obraz.naturalHeight }),
    );
    obraz.addEventListener('error', () => rozstrzygnij(brak));
    obraz.src = `data:${typMime === '' ? 'application/octet-stream' : typMime};base64,${trescBase64}`;
  });
}
