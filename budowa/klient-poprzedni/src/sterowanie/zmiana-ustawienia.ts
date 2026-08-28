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

/** Wysyłka i odczyt ustawień z poziomu zasięgu okna komunikacji, adresowanych poziomem, bytem poziomu, osią i bytem osi. */
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
      // Zapis szerszy niż okno nie dotyczy okna, ale zmienia wartość obowiązującą, więc odświeża politykę.
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
