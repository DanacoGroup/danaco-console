import { Command, type DesignAsset } from '../../../../shared/contract';

/**
 * Wspólny opis zasobu oddanego przez `design.asset.generate` — używany przez
 * Prompt Buildera i Preview Window, żeby oba mówiły o zasobie tak samo.
 *
 * Rdzeń odkłada bajty obrazu w magazynie pod sumą sha256 i zakłada wiersz zasobu
 * z `uri`, `format` i wymiarami odczytanymi z nagłówka pliku
 * (`adapter_modul_design_generowanie.go`). Pole `name` jest opcjonalne i bywa
 * odpowiedzią modelu tekstowego, więc okno traktuje je jako słowo modelu,
 * a nie tytuł zasobu.
 */

/**
 * Słowo modelu — treść pola `name` bez upiększania.
 *
 * Pusty wynik znaczy, że rdzeń nie podał nazwy w ogóle; to inny stan niż
 * nazwa pusta w znakach i okno musi je rozróżniać.
 */
export function slowoModelu(zasob: DesignAsset): string {
  return (zasob.name ?? '').trim();
}

/** Czy rdzeń zna treść obrazu tego zasobu — pole `uri` wypełnione. */
export function czyRdzenZnaObraz(zasob: DesignAsset): boolean {
  return (zasob.uri ?? '').trim() !== '';
}

/**
 * Zdanie o tym, czym jest nazwa zasobu — pokazywane zawsze, niezależnie od jej
 * treści, bo warunkowe pokazywanie byłoby oceną treści modelu.
 */
export function zdanieOSlowieModelu(zasob: DesignAsset): string {
  if (slowoModelu(zasob) === '') {
    return 'Rdzeń nie podał nazwy tego zasobu — pole name kontraktu jest opcjonalne.';
  }
  return (
    'Nazwa zasobu jest treścią pola name zapisaną w rdzeniu, przepisaną co do znaku. ' +
    'Okno jej nie ocenia i nie podmienia na własną — nie wie, kto ją nadał.'
  );
}

/**
 * Zdanie o zasobie, przy którym rdzeń nie podał `uri`.
 *
 * Generowanie zakłada wiersz dopiero po utrwaleniu bajtów w magazynie, a wgranie
 * — po zapisaniu pliku, więc pusty `uri` znaczy zasób założony drogą, która pola
 * nie wypełniła. Okno nie zgaduje którą — mówi, czego nie ma.
 */
export function zdanieOBrakuObrazu(): string {
  return (
    'Rdzeń nie podał dla tego zasobu pola uri, więc nie wskazał żadnej treści obrazu. ' +
    'Generowanie zakłada wiersz dopiero po utrwaleniu bajtów w magazynie, a wgranie — po ' +
    'zapisaniu pliku, więc zasób bez uri pochodzi sprzed tej zasady albo powstał inaczej. ' +
    'Podgląd pokazuje to, co rdzeń o zasobie wie, i nie dorabia obrazu, którego nie wskazano.'
  );
}

/**
 * Dlaczego treść obrazu nie wyświetli się w przeglądarce, choć rdzeń ją ma.
 *
 * Pole `uri` zasobu jest ścieżką w systemie plików rdzenia: magazyn treści
 * oddaje `filepath.Join(katalog, suma[:2], suma)` i ta ścieżka wchodzi do
 * wiersza jako odwołanie (`core/adapter_modul_library_magazyn.go`, `Zapisz`).
 * Przeglądarka ścieżki dyskowej nie otworzy, więc wstawiona w `<img src>` daje
 * zdarzenie `error`, którego przyczyną nie jest sama przeglądarka.
 *
 * Komenda oddająca treść zasobu jest już w kontrakcie i obejmuje cały magazyn,
 * nie sam obszar Designu — zasób Design, plik Library i dokument Studia leżą
 * w jednym repozytorium. Brakuje jej uchwytu w rdzeniu, więc okno nadal nie ma
 * skąd wziąć bajtów; to jest jednak brak obsługi, a nie brak drogi.
 */
export function zdanieOSciezceRdzenia(adres: string): string {
  return (
    `Rdzeń podał uri „${adres}", ale to jest ŚCIEŻKA W JEGO SYSTEMIE PLIKÓW, nie adres ` +
    'do pobrania. Magazyn treści zapisuje bajty pod sumą kontrolną i oddaje ścieżkę pliku, ' +
    'a przeglądarka takiego adresu nie otworzy — także wtedy, gdy obraz powstał poprawnie. ' +
    `Droga po bajty JEST już opisana w kontrakcie (${Command.DesignAssetContentGet}) i jest ` +
    'jedna dla całego magazynu, więc obsłuży zarówno zasób Designu, jak i plik Library. ' +
    'Czeka na dobudowę uchwytu w rdzeniu — to jest brak obsługi, nie usterka przeglądarki ' +
    'ani tego okna.'
  );
}
