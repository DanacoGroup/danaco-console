import {
  Command,
  ConfigScope,
  type Component,
  type ComponentAssignRequest,
  type ComponentUpdateRequest,
} from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Kanal } from '../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/** Zmiana komponentu własnego i jego przypisanie do poziomu zasięgu są dwiema czynnościami strefy drugiej strony głównej, wykonywanymi na komponencie już założonym. */
export interface PoziomPrzypisania {
  kod: ConfigScope;
  /** Nazwa poziomu widziana przez Operatora. */
  nazwa: string;
  /** Czy poziom żąda identyfikatora bytu; global go nie przyjmuje wcale. */
  zBytem: boolean;
  /** Skąd bierze się wykaz bytów tego poziomu; pusty napis znaczy „wpisywany". */
  zrodloBytow: 'srodowiska' | 'okna' | 'wpisywany';
}

export const POZIOMY_PRZYPISANIA: readonly PoziomPrzypisania[] = [
  { kod: ConfigScope.Global, nazwa: 'Poziom globalny', zBytem: false, zrodloBytow: 'wpisywany' },
  {
    kod: ConfigScope.Environment,
    nazwa: 'Środowisko',
    zBytem: true,
    zrodloBytow: 'srodowiska',
  },
  {
    kod: ConfigScope.Window,
    nazwa: 'Okno komunikacji',
    zBytem: true,
    zrodloBytow: 'okna',
  },
  { kod: ConfigScope.Project, nazwa: 'Projekt', zBytem: true, zrodloBytow: 'wpisywany' },
  { kod: ConfigScope.Session, nazwa: 'Karta sesji', zBytem: true, zrodloBytow: 'wpisywany' },
];

/** Poziomy, których rdzeń w przypisaniu nie przyjmuje, wraz z powodem przeznaczonym do pokazania w zdaniu widoku. */
export const POZIOMY_ODRZUCANE =
  'Poziomów modułu, pary modułów, roli i aplikacji rdzeń w przypisaniu nie przyjmuje — ' +
  'odpowiada odmową sprawdzenia żądania. Brak jest po stronie RDZENIA, w rozstrzygnięciu ' +
  'komendy, nie w tym oknie.';

/** Byt poziomu wraz z jego nazwą — to, z czym okno wiąże komponent podczas przypisania zasięgu komponentu. */
export interface BytPoziomu {
  id: string;
  nazwa: string;
}

/** Wynik czynności zmiany albo przypisania w postaci, którą widok pokazuje bez dopowiadania niczego od siebie. */
export interface WynikCzynnosci {
  udane: boolean;
  zdanie: string;
  komponent?: Component;
}

export interface ZmianaKomponentu {
  /** `component.list` — komponenty wraz z wyłączonymi, bo te też się zmienia. */
  wykaz(): Promise<{ komponenty: readonly Component[]; zdanie: string }>;
  /** `component.update` — zmienia wyłącznie pola podane. */
  zmien(zadanie: ComponentUpdateRequest): Promise<WynikCzynnosci>;
  /** `component.assign` — wiąże komponent z poziomem zasięgu. */
  przypisz(zadanie: ComponentAssignRequest): Promise<WynikCzynnosci>;
  /** `environment.list` — byty poziomu środowiska wraz z nazwami. */
  srodowiska(): Promise<readonly BytPoziomu[]>;
  /** `window.list` — byty poziomu okna komunikacji. */
  okna(): Promise<readonly BytPoziomu[]>;
}

/** Zmiana komponentu udostępnia wykaz, zmianę pól, przypisanie zasięgu oraz byty poziomów środowiska i okna komunikacji. */
export function utworzZmianeKomponentu(kanal: Kanal): ZmianaKomponentu {
  return {
    async wykaz() {
      // Wykaz z włączeniem wyłączonych: komponent wyłączony jest tym, który operator przychodzi włączyć.
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ComponentList, { includeDisabled: true }),
        Command.ComponentList,
        (tresc) => czyTablica(tresc.components),
      );
      if (!wynik.udany || wynik.wynik === undefined) {
        return { komponenty: [], zdanie: opisOdmowyBledu('Wykaz komponentów', wynik.blad) };
      }
      return { komponenty: wynik.wynik.components, zdanie: '' };
    },

    async zmien(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ComponentUpdate, zadanie),
        Command.ComponentUpdate,
        (tresc) => czyObiekt(tresc.component),
      );
      if (!wynik.udany || wynik.wynik === undefined) {
        return { udane: false, zdanie: opisOdmowyBledu('Zmiana komponentu', wynik.blad) };
      }
      const komponent = wynik.wynik.component;
      // Zdanie mówi o tym, co oddał rdzeń, a nie o tym, co poszło w żądaniu: nazwa bywa przycięta.
      return {
        udane: true,
        zdanie:
          `Komponent „${komponent.name}" zmieniony: ` +
          `${komponent.enabled ? 'czynny' : 'wyłączony'}, ostatnia zmiana ` +
          `${new Date(komponent.updatedAt).toLocaleString('pl-PL')}.`,
        komponent,
      };
    },

    async przypisz(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ComponentAssign, zadanie),
        Command.ComponentAssign,
        (tresc) => czyObiekt(tresc.component),
      );
      if (!wynik.udany || wynik.wynik === undefined) {
        return { udane: false, zdanie: opisOdmowyBledu('Przypisanie komponentu', wynik.blad) };
      }
      const tresc = wynik.wynik;
      // Wynik odmowny znaczy powtórzenie tego samego przypisania — nic nie doszło do skutku.
      return {
        udane: true,
        zdanie: tresc.assigned
          ? `Komponent „${tresc.component.name}" przypisany do poziomu ${zadanie.scope}` +
            `${zadanie.scopeId === undefined || zadanie.scopeId === '' ? '' : ` (${zadanie.scopeId})`}. ` +
            'Przypisanie zapamiętuje POZIOM, na którym komponent obowiązuje — nie przenosi bytu ' +
            'modułowego i nie nadaje uprawnień.'
          : `Nic się nie zmieniło: komponent „${tresc.component.name}" był już przypisany ` +
            'dokładnie tam. Rdzeń odpowiedział, że przypisanie nie doszło do skutku.',
        komponent: tresc.component,
      };
    },

    async srodowiska() {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.EnvironmentList, {}),
        Command.EnvironmentList,
        (tresc) => czyTablica(tresc.environments),
      );
      if (!wynik.udany || wynik.wynik === undefined) return [];
      return wynik.wynik.environments.map((srodowisko) => ({
        id: srodowisko.id,
        nazwa: `${srodowisko.name} (${srodowisko.code})`,
      }));
    },

    async okna() {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowList, {}),
        Command.WindowList,
        (tresc) => czyTablica(tresc.windows),
      );
      if (!wynik.udany || wynik.wynik === undefined) return [];
      // Okno komunikacji nie ma w kontrakcie nazwy własnej, więc nazywa się tym, co je odróżnia.
      return wynik.wynik.windows.map((okno) => ({
        id: okno.id,
        nazwa: `${okno.moduleId} · sesja ${okno.sessionId} · ${okno.id}`,
      }));
    },
  };
}

/** Poziom przypisania o wskazanym kodzie; kod spoza wykazu bierze poziom globalny jako wartość zastępczą. */
export function poziomOKodzie(kod: string): PoziomPrzypisania {
  return (
    POZIOMY_PRZYPISANIA.find((pozycja) => pozycja.kod === kod) ??
    (POZIOMY_PRZYPISANIA[0] as PoziomPrzypisania)
  );
}
