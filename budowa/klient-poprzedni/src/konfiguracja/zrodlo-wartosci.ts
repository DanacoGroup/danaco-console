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
 * i przywrócenie pojedynczego klucza pod wskazanym adresem.
 *
 * Pola osi są w kontrakcie nieobowiązkowe — puste znaczy `platform`
 * — więc adres platformy wychodzi bez nich.
 *
 * Odczyt idzie po poziomach łańcucha, nie po polityce efektywnej. `config.get`
 * bez poziomu oddaje samą wartość obowiązującą, bez pozostałych zapisów, a okno
 * pokazuje również, z którego poziomu wartość pochodzi i jakie zapisy stoją
 * obok. Łańcuch punktu widzenia liczy klient (`rozstrzygniecie.ts`) z wpisów
 * globalnych oraz wpisów zapisanych dokładnie na wskazanym poziomie i bycie,
 * więc odczyt pobiera surowe wpisy tych poziomów osobnymi zapytaniami
 * z podanym `scope`.
 *
 * Surowy odczyt zawęża się do granulacji poziomu (zasięg i byt poziomu), bez
 * osi: łańcuch bierze wszystkie zapisy poziomu niezależnie od osi, a którą oś
 * przyjąć rozstrzyga klient względem punktu widzenia.
 *
 * Nazwy komend i zdarzeń pochodzą wyłącznie ze stałych kontraktu.
 */
export interface ZrodloWartosci {
  /**
   * `config.get` — surowe wpisy tworzące łańcuch dziedziczenia punktu
   * widzenia: poziom globalny oraz, gdy punkt nie jest globalny, poziom nim
   * wskazany. Wynik zasila rozstrzyganie pochodzenia po stronie klienta.
   *
   * Pominięty punkt znaczy widok globalny — tyle wystarcza czytelnikowi
   * przypiętemu do poziomu globalnego (sekcja katalogu roboczego).
   */
  wpisy(punkt?: AdresUstawienia): Promise<ConfigEntry[]>;
  /** `config.set` — zapis wartości pod wskazanym adresem. */
  zapisz(klucz: string, wartosc: unknown, adres: AdresUstawienia): Promise<Wynik<unknown>>;
  /** `config.reset` — usunięcie zapisu spod adresu; wraca wartość szersza. */
  przywroc(klucz: string, adres: AdresUstawienia): Promise<Wynik<unknown>>;
  /**
   * Subskrypcja `config.changed` — zapis albo przywrócenie z innego okna
   * lub urządzenia.
   *
   * Słuchacz dostaje wpis wraz z rodzajem zmiany. Bez rodzaju nie da się
   * odróżnić zapisu od usunięcia, bo rdzeń rozgłasza przywrócenie wartości
   * domyślnej wpisem niosącym starą wartość (`core/handlers_config.go`,
   * `core/adapter_ustawienia.go` funkcja `Przywroc`); czytelnik patrzący na sam
   * wpis wstawiłby skasowany zapis z powrotem do wykazu.
   */
  naZmiane(sluchacz: (wpis: ConfigEntry, rodzaj: ChangeKind) => void): Odsubskrybuj;
}

export function utworzZrodloWartosci(
  kanal: Kanal,
  naNiepowodzenie: NaNiepowodzenie = () => undefined,
): ZrodloWartosci {
  // Surowe wpisy jednego poziomu: podanie `scope` zdejmuje z rdzenia liczenie
  // polityki efektywnej i zwraca wpisy zapisane dokładnie na tym poziomie.
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
 *
 * Łańcuchowi potrzebne są wszystkie zapisy poziomu niezależnie od osi — oś
 * rozstrzyga później klient — więc surowy odczyt nie niesie osi, choćby punkt
 * widzenia ją wskazywał. Zasięg globalny wychodzi bez bytu.
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
 *
 * Eksportowany, bo ten sam czteroczłonowy adres niosą także `config.session.set`
 * i `config.effective.get`; wszystkie trzy komendy używają jednego tłumaczenia
 * `AdresUstawienia` na pola kontraktu.
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

/** Czy wpis leży dokładnie pod wskazanym adresem i dotyczy wskazanego klucza. */
export function wpisSpodAdresu(
  wpis: ConfigEntry,
  klucz: string,
  adres: AdresUstawienia,
): boolean {
  return wpis.key === klucz && tenSamAdres(adresWpisu(wpis), adres);
}
