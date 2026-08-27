import type { LibraryFile, LibraryVersion } from '../../../../shared/contract';
import {
  pozycjaWykazu,
  przelacznik,
  przyciskAkcji as przycisk,
  wykaz,
} from '../../modele/kontrolki-formularza';

/**
 * Jedna pozycja wykazu zasobów projektu wraz z zaznaczeniem do czynności
 * zbiorowych i skokiem do wersji pliku.
 *
 * Osobny plik od okna, bo to inna odpowiedzialność: okno prowadzi odczyt,
 * wgranie i pracę zbiorową, pozycja rysuje jeden plik.
 */
export interface CzynnosciPliku {
  zaznaczony: boolean;
  przelacz(zaznacz: boolean): void;
  pokazWersje(): void;
}

export function pozycjaPliku(wpis: LibraryFile, czynnosci: CzynnosciPliku): HTMLElement {
  const opis = [
    wpis.path ?? wpis.name,
    wpis.mimeType ?? 'rodzaj nierozpoznany',
    `${Math.max(1, Math.round((wpis.sizeBytes ?? 0) / 1024))} kB`,
    new Date(wpis.updatedAt).toLocaleString('pl-PL'),
  ].join(' · ');
  const { element, akcje } = pozycjaWykazu(wpis.name, opis, 'dw');

  const zaznaczenie = przelacznik(`Zaznacz ${wpis.name}`);
  zaznaczenie.checked = czynnosci.zaznaczony;
  zaznaczenie.addEventListener('change', () => czynnosci.przelacz(zaznaczenie.checked));

  const wersje = przycisk('Wersje', 'dn-btn dn-btn--zarys');
  wersje.addEventListener('click', () => czynnosci.pokazWersje());

  element.prepend(zaznaczenie);
  akcje.append(wersje);
  return element;
}

/** Buduje wykaz wersji jednego pliku biblioteki projektu na podstawie odpowiedzi komendy rdzenia wymieniającej wersje. */
export function wykazWersji(nazwaPliku: string, wersje: readonly LibraryVersion[]): HTMLElement {
  const lista = wykaz(`Wersje pliku ${nazwaPliku}`, 'dw-wykaz');
  for (const wersja of wersje) {
    lista.append(
      pozycjaWykazu(
        wersja.label ?? wersja.id,
        `${wersja.author ?? 'sprawca nieznany'} · ${new Date(wersja.createdAt).toLocaleString('pl-PL')}`,
        'dw',
      ).element,
    );
  }
  return lista;
}

/** Zamienia binarną zawartość pliku na postać tekstową base64 zgodną z odpowiednim polem kontraktu wgrywania. */
export function naBase64(zawartosc: ArrayBuffer): string {
  const bajty = new Uint8Array(zawartosc);
  let tekst = '';
  for (const bajt of bajty) tekst += String.fromCharCode(bajt);
  return btoa(tekst);
}
