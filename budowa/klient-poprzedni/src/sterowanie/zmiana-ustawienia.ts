import {
  ChangeKind,
  Command,
  ConfigAxis,
  EventType,
  type ConfigEntry,
} from '../../../shared/contract';
import {
  adresWpisu,
  tenSamAdres,
  type AdresUstawienia,
} from '../konfiguracja/adres-ustawienia';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import { KLUCZE_KOMPLETU, ZASIEG_OKNA } from './klucze-ustawien';
import { niepowodzenie, potwierdzenie, type OdbiorcaKomunikatu } from './komunikat-zmiany';
import type { StanSterowania } from './stan-sterowania';
import { utworzOdczytObowiazujacej } from './wartosc-obowiazujaca';

/**
 * Wysyłka i odczyt ustawień z poziomu zasięgu okna komunikacji.
 *
 * Trzy sterowania — host wykonania, kanał zapasowy, nakład rozumowania — nie
 * mają pola w treści `window.update`, więc idą komendą `config.set` na poziomie
 * `window`, z identyfikatorem okna jako bytem poziomu. Odczyt przy otwarciu
 * kompletu (`config.get`) sprawia, że po przeładowaniu ustawienia wracają
 * z bazy, a nie z pamięci przeglądarki.
 *
 * Adres ustawienia ma cztery człony, nie dwa: przestrzeń konfiguracji ma dwa
 * prostopadłe wymiary — poziom zasięgu z bytem poziomu oraz oś rozstrzygania
 * z bytem osi. Sprawdzanie samego poziomu przyjmowałoby za swój zapis
 * adresowany do innego modelu albo innego konta, na przykład
 * `naklad_rozumowania` na osi `model/…`. Adres składa
 * `konfiguracja/adres-ustawienia.ts`, jedyny w kliencie przepis na adres tej
 * przestrzeni.
 *
 * Wpis konfiguracji innego okna, innego poziomu albo innej osi jest pomijany;
 * to adresowanie, a nie wstrzymanie zmiany.
 *
 * Zdarzenie `config.changed` niesie rodzaj zmiany: `created` i `updated`
 * nanoszą wartość, `deleted` ją zdejmuje. Rdzeń rozgłasza usunięcie z wpisem
 * o starej wartości (`server/internal/core/handlers_config.go`), więc czytanie
 * samego wpisu bez rodzaju nanosiłoby skasowaną wartość jak świeży zapis,
 * a powrót do wartości domyślnej nie docierałby do żadnej powierzchni.
 */
export interface ZmianaUstawienia {
  /** Wczytuje ustawienia poziomu okna z rdzenia wraz z wartością obowiązującą. */
  wczytaj(): void;
  /** Zapisuje pojedyncze ustawienie poziomu okna. */
  zapisz(nazwa: string, klucz: string, wartosc: unknown): void;
  /** Odłącza subskrypcję zdarzeń. */
  rozlacz(): void;
}

export function utworzZmianeUstawienia(
  kanal: Kanal,
  stan: StanSterowania,
  zglos: OdbiorcaKomunikatu,
): ZmianaUstawienia {
  const obowiazujaca = utworzOdczytObowiazujacej(kanal, stan.idOkna);

  /** Adres tego okna: poziom okna, oś platformy (oś pominięta znaczy platform). */
  function adresOkna(): AdresUstawienia {
    return {
      zasieg: ZASIEG_OKNA,
      bytZasiegu: stan.idOkna(),
      os: ConfigAxis.Platform,
      bytOsi: '',
    };
  }

  /** Czy wpis konfiguracji leży dokładnie pod adresem tego okna. */
  function dotyczyOkna(wpis: ConfigEntry): boolean {
    return tenSamAdres(adresWpisu(wpis), adresOkna());
  }

  /** Odświeża wartość obowiązującą; wynik wchodzi do stanu, gdy rdzeń odpowie. */
  function odswiezObowiazujaca(): void {
    obowiazujaca.wczytaj(stan.przyjmijObowiazujace);
  }

  const odsubskrybuj: Odsubskrybuj = kanal.naZdarzenie(
    EventType.ConfigChanged,
    (tresc) => {
      // Zapis na poziomie szerszym niż okno nie dotyczy poziomu okna, ale
      // zmienia wartość obowiązującą, dlatego rodzina kluczy kompletu odświeża
      // politykę niezależnie od adresu wpisu.
      if (KLUCZE_KOMPLETU.includes(tresc.entry.key)) odswiezObowiazujaca();
      if (!dotyczyOkna(tresc.entry)) return;
      if (tresc.change === ChangeKind.Deleted) {
        stan.przyjmijUsuniecie(tresc.entry.key);
        return;
      }
      stan.przyjmijUstawienie(tresc.entry.key, tresc.entry.value);
    },
  );

  return {
    wczytaj() {
      kanal.wyslij(
        Command.ConfigGet,
        { scope: ZASIEG_OKNA, scopeId: stan.idOkna() },
        (wynik) => {
          for (const wpis of wynik.wynik?.entries ?? []) {
            if (dotyczyOkna(wpis)) stan.przyjmijUstawienie(wpis.key, wpis.value);
          }
        },
      );
      odswiezObowiazujaca();
    },

    zapisz(nazwa, klucz, wartosc) {
      kanal.wyslij(
        Command.ConfigSet,
        { key: klucz, value: wartosc, scope: ZASIEG_OKNA, scopeId: stan.idOkna() },
        (wynik) => {
          const wpis = wynik.wynik?.entry;
          if (wynik.udany && wpis !== undefined) {
            stan.przyjmijUstawienie(wpis.key, wpis.value);
            zglos(potwierdzenie(nazwa));
            return;
          }
          stan.odswiez();
          zglos(niepowodzenie(nazwa, wynik.blad));
        },
      );
    },

    rozlacz: odsubskrybuj,
  };
}
