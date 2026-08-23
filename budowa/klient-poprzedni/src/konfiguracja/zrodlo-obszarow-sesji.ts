import {
  Command,
  type ConfigEffectiveGetResponse,
  type ConfigSessionSetResponse,
  type SessionConfig,
  type SessionConfigArea,
} from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';
import type { AdresUstawienia } from './adres-ustawienia';
import { czescAdresu } from './zrodlo-wartosci';

/**
 * Dwie komendy jednolitego modelu konfiguracji sesji: `config.effective.get`
 * i `config.session.set`.
 *
 * `config.get` oddaje pojedyncze klucze katalogu ustawień — płaskie wpisy
 * `ConfigEntry` ze wszystkich poziomów naraz — więc dziedziczenie jednego klucza
 * rozstrzyga klient (`rozstrzygniecie.ts`). `config.effective.get` oddaje obszary
 * konfiguracji sesji już rozstrzygnięte przez rdzeń, wraz z pochodzeniem obszaru
 * (pole `source`) i z rozejściem katalogu roboczego. Klient tego rachunku wykonać
 * nie może: nie zna pełnej ścieżki bytów ani rejestrów spoza rodziny `config.*`
 * (rejestr kont, rejestr kanałów, katalog tożsamości, nadania dostępu), z których
 * treść obszaru jest składana.
 *
 * Wynik obu czynności wraca opakowany, więc okno odróżnia obszar bez zapisu
 * od nieudanego zapytania.
 */
export interface ZrodloObszarowSesji {
  /** `config.effective.get` — obszary rozstrzygnięte wraz z pochodzeniem. */
  obowiazujaca(
    punkt: AdresUstawienia,
    obszary?: readonly SessionConfigArea[],
  ): Promise<Wynik<ConfigEffectiveGetResponse>>;
  /** `config.session.set` — zapis wskazanych obszarów pod wskazanym adresem. */
  zapiszObszary(
    punkt: AdresUstawienia,
    obszary: readonly SessionConfigArea[],
    konfiguracja: SessionConfig,
  ): Promise<Wynik<ConfigSessionSetResponse>>;
}

export function utworzZrodloObszarowSesji(kanal: Kanal): ZrodloObszarowSesji {
  return {
    async obowiazujaca(punkt, obszary) {
      // Poziom i byt idą tym samym tłumaczeniem adresu, co `config.set`.
      // Pusty wykaz obszarów znaczy w kontrakcie komplet, więc go nie wysyłamy.
      const zadanie = {
        ...czescAdresu(punkt),
        ...(obszary === undefined || obszary.length === 0 ? {} : { areas: [...obszary] }),
      };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ConfigEffectiveGet, zadanie),
        Command.ConfigEffectiveGet,
        czyKsztaltObowiazujacej,
      );
    },

    async zapiszObszary(punkt, obszary, konfiguracja) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ConfigSessionSet, {
          areas: [...obszary],
          config: konfiguracja,
          ...czescAdresu(punkt),
        }),
        Command.ConfigSessionSet,
        czyKsztaltZapisu,
      );
    },
  };
}

/**
 * Kształt odpowiedzi `config.effective.get`, tak jak produkuje ją
 * `server/internal/core/sesja_konfiguracja_skladanie.go`: `config`, `origins`
 * i `workingDirectory` wychodzą zawsze, `capabilities` wyłącznie na żądanie,
 * a `unsupportedFields` tą ścieżką nie wychodzi. Sprawdzian pilnuje trzech pól
 * obowiązkowych i ani jednego więcej — pole nieobowiązkowe w sprawdzianie
 * zamieniłoby zdrową odpowiedź w fałszywą odmowę.
 */
function czyKsztaltObowiazujacej(tresc: ConfigEffectiveGetResponse): boolean {
  const obowiazujaca = tresc.effective as unknown;
  if (!czyObiekt(obowiazujaca)) return false;
  return (
    czyObiekt(obowiazujaca['config']) &&
    czyTablica(obowiazujaca['origins']) &&
    czyObiekt(obowiazujaca['workingDirectory'])
  );
}

/**
 * Kształt odpowiedzi `config.session.set`: rdzeń oddaje stan poziomu po zapisie
 * (`config`), wykaz obszarów faktycznie na nim zapisanych (`storedAreas`)
 * i wpisy rezolwera. `unsupportedFields` tą ścieżką nie wraca — okno pokazuje
 * je, gdy przyjdą, i nie wymaga ich do uznania zapisu za udany.
 */
function czyKsztaltZapisu(tresc: ConfigSessionSetResponse): boolean {
  return czyObiekt(tresc.config as unknown) && czyTablica(tresc.storedAreas);
}
