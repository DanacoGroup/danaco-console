import {
  Command,
  EventType,
  ConfigScope,
  type ChangeKind,
  type ConfigAxis,
  type ConfigEntry,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';
import {
  adresPoczatkowy,
  adresWpisu,
  tenSamAdres,
  type AdresUstawienia,
} from './adres-ustawienia';
import { powodOdczytu, type NaNiepowodzenie } from './zrodlo-katalogu';

/**
 * Wartości ustawień: odczyt łańcucha dziedziczenia punktu widzenia oraz zapis
 * i przywrócenie klucza pod wskazanym adresem. Odczyt idzie po poziomach
 * łańcucha, nie po polityce efektywnej. Nazwy komend pochodzą wyłącznie ze
 * stałych kontraktu.
 */
export interface ZrodloWartosci {
  /** `config.get` — surowe wpisy tworzące łańcuch dziedziczenia dla punktu widzenia i poziomu. */
  wpisy(punkt?: AdresUstawienia): Promise<ConfigEntry[]>;
  /** `config.set` — zapis wartości pod wskazanym adresem. */
  zapisz(klucz: string, wartosc: unknown, adres: AdresUstawienia): Promise<Wynik<unknown>>;
  /** `config.reset` — usunięcie zapisu spod adresu; wraca wartość szersza. */
  przywroc(klucz: string, adres: AdresUstawienia): Promise<Wynik<unknown>>;
  /** Subskrypcja `config.changed`: wpis niesie rodzaj zmiany — zapis albo przywrócenie z innego okna. */
  naZmiane(sluchacz: (wpis: ConfigEntry, rodzaj: ChangeKind) => void): Odsubskrybuj;
}

export function utworzZrodloWartosci(
  kanal: Kanal,
  naNiepowodzenie: NaNiepowodzenie = () => undefined,
): ZrodloWartosci {
  // Surowe wpisy jednego poziomu: `scope` zdejmuje z rdzenia liczenie polityki efektywnej.
  const pobierzPoziom = async (adres: AdresUstawienia): Promise<ConfigEntry[]> => {
    const wynik = sprawdzKsztalt(
      await wywolaj(kanal, Command.ConfigGet, poziomZadania(adres)),
      Command.ConfigGet,
      (tresc) => czyTablica(tresc.entries),
    );
    if (!wynik.udany) {
      console.warn('[konfiguracja] wpisy nie dotarły', wynik.blad?.message ?? '');
      naNiepowodzenie(powodOdczytu(Command.ConfigGet, wynik.blad?.message));
      return [];
    }
    return wynik.wynik?.entries ?? [];
  };

  return {
    async wpisy(punkt = adresPoczatkowy()) {
      const globalne = await pobierzPoziom(adresPoczatkowy());
      if (punkt.zasieg === ConfigScope.Global) return globalne;
      const poziomowe = await pobierzPoziom(punkt);
      return [...globalne, ...poziomowe];
    },

    zapisz(klucz, wartosc, adres) {
      return wywolaj(kanal, Command.ConfigSet, {
        key: klucz,
        value: wartosc,
        ...czescAdresu(adres),
      });
    },

    przywroc(klucz, adres) {
      return wywolaj(kanal, Command.ConfigReset, {
        key: klucz,
        ...czescAdresu(adres),
      });
    },

    naZmiane(sluchacz) {
      return kanal.naZdarzenie(EventType.ConfigChanged, (tresc) =>
        sluchacz(tresc.entry, tresc.change),
      );
    },
  };
}

/**
 * Treść `config.get` zawężona do poziomu: sam zasięg i byt poziomu, bez osi.
 * Łańcuchowi potrzebne są wszystkie zapisy poziomu niezależnie od osi — oś
 * rozstrzyga później klient. Zasięg globalny wychodzi bez bytu.
 */
export function poziomZadania(adres: AdresUstawienia): {
  scope: ConfigScope;
  scopeId?: string;
} {
  const zadanie: { scope: ConfigScope; scopeId?: string } = { scope: adres.zasieg };
  if (adres.bytZasiegu !== '') zadanie.scopeId = adres.bytZasiegu;
  return zadanie;
}

/**
 * Adres w kształcie treści komendy; człony puste zostają pominięte.
 * Eksportowany, bo ten sam czteroczłonowy adres niosą także
 * `config.session.set` i `config.effective.get`.
 */
export function czescAdresu(adres: AdresUstawienia): {
  scope: ConfigScope;
  scopeId?: string;
  axis?: ConfigAxis;
  axisId?: string;
} {
  const czesc: { scope: ConfigScope; scopeId?: string; axis?: ConfigAxis; axisId?: string } = {
    scope: adres.zasieg,
  };
  if (adres.bytZasiegu !== '') czesc.scopeId = adres.bytZasiegu;
  if (adres.bytOsi !== '') {
    czesc.axis = adres.os;
    czesc.axisId = adres.bytOsi;
  }
  return czesc;
}

/** Czy wpis leży dokładnie pod wskazanym adresem konfiguracji i dotyczy podanego wprost klucza ustawienia. */
export function wpisSpodAdresu(
  wpis: ConfigEntry,
  klucz: string,
  adres: AdresUstawienia,
): boolean {
  return wpis.key === klucz && tenSamAdres(adresWpisu(wpis), adres);
}
