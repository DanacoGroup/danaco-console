/**
 * Wspólny opis zasobu oddanego przez `design.asset.generate` — używany przez Prompt
 * Buildera i Preview Window, żeby oba mówiły o zasobie tak samo. Pole `name` jest
 * opcjonalne, więc okno traktuje je jako słowo modelu, a nie tytuł zasobu.
 */

import { Command, type DesignAsset } from '../../../../shared/contract';


/**
 * Słowo modelu — treść pola `name` bez upiększania.
 *
 * Pusty wynik znaczy, że rdzeń nie podał nazwy w ogóle; to inny stan niż
 * nazwa pusta w znakach i okno musi je rozróżniać.
 */
export function slowoModelu(zasob: DesignAsset): string {
  return (zasob.name ?? '').trim();
}

/**
 * Czy rdzeń zna treść obrazu tego zasobu — pole `uri` wypełnione. Pole puste znaczy
 * zasób założony drogą, która treści nie wskazała, i okno mówi o tym wprost zamiast
 * pokazywać podgląd bez pokrycia.
 */
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
 * Zdanie o zasobie, przy którym rdzeń nie podał pola `uri`. Okno nie zgaduje, którą
 * drogą taki zasób powstał — mówi, czego nie ma, bo podgląd pokazuje wyłącznie to,
 * co rdzeń o zasobie wie.
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
 * Dlaczego treść obrazu nie wyświetli się w przeglądarce, choć rdzeń ją ma. Pole
 * `uri` zasobu jest ścieżką w systemie plików rdzenia, a przeglądarka ścieżki dyskowej
 * nie otworzy, więc wstawiona w `<img src>` daje zdarzenie błędu.
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
