import manifest from '../../package.json' with { type: 'json' };
import { nowyIdentyfikator } from './identyfikator.ts';

/**
 * Tożsamość klienta przedstawiana rdzeniowi w powitaniu połączenia.
 *
 * Wersja pochodzi z manifestu pakietu, nie z literału w kodzie — produkt ma
 * jedną wersję i nie ma jej gdzie rozejść.
 */
export interface TozsamoscKlienta {
  /** Identyfikator instancji klienta, nadawany na czas uruchomienia. */
  id: string;
  /** Wersja klienta z manifestu pakietu. */
  wersja: string;
}

export function tozsamoscKlienta(): TozsamoscKlienta {
  return { id: nowyIdentyfikator('klient'), wersja: manifest.version };
}
