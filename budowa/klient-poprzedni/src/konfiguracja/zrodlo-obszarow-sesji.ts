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
 * oddaje obszary rozstrzygnięte przez rdzeń wraz z pochodzeniem, a
 * `config.session.set` zapisuje wskazane obszary pod wskazanym adresem. Wynik
 * obu czynności wraca opakowany.
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
      // Pusty wykaz obszarów znaczy w kontrakcie komplet, więc nie jest wysyłany.
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
 * Kształt odpowiedzi `config.effective.get`: pola konfiguracji, pochodzenia
 * i katalogu roboczego wychodzą zawsze, a możliwości wyłącznie na żądanie.
 * Sprawdzian pilnuje trzech pól obowiązkowych i ani jednego więcej.
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
 * Kształt odpowiedzi `config.session.set`: rdzeń oddaje stan poziomu po zapisie,
 * wykaz obszarów faktycznie na nim zapisanych oraz wpisy rezolwera. Sprawdzian
 * obejmuje dwa pola obowiązkowe.
 */
function czyKsztaltZapisu(tresc: ConfigSessionSetResponse): boolean {
  return czyObiekt(tresc.config as unknown) && czyTablica(tresc.storedAreas);
}
