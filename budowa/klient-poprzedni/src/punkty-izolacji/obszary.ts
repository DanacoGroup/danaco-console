// Rejestr sześciu obszarów okna Punktów Izolacji wraz z pełnym kontraktem wpięcia każdego obszaru w to okno.
import type { Kanal } from '../protokol/kanal';
import { utworzObszar as utworzObszarEfektywna } from './obszar-efektywna';
import { utworzObszar as utworzObszarKontekst } from './obszar-kontekst';
import { utworzObszar as utworzObszarOsie } from './obszar-osie';
import { utworzObszar as utworzObszarPoziomy } from './obszar-poziomy';
import { utworzObszar as utworzObszarProfile } from './obszar-profile';
import { utworzObszar as utworzObszarTechniczny } from './obszar-techniczny';
import type { StanWarstwy } from './stan-warstwy';
import type { StanZasiegu } from './stan-zasiegu';

/** Kod jednego z sześciu obszarów okna Punktów Izolacji, jednoznacznie wskazujący dany obszar w rejestrze. */
export type KodObszaruIzolacji = 'kontekst' | 'techniczny' | 'profile' | 'osie' | 'poziomy' | 'efektywna';

/** Panel okna, w którym obszar stoi — trzy kolumny widoku: zasięg, macierz izolacji oraz profil obszaru. */
export type PanelIzolacji = 'zasieg' | 'macierz' | 'profil';

/** Zależności dostępne każdemu obszarowi przy budowie: kanał do rdzenia, warstwa oraz zasięg czynny okna. */
export interface ZaleznosciObszaru {
  /** Kanał do rdzenia — obszar woła nim komendy `isolation.*` i zakłada własne subskrypcje. */
  kanal: Kanal;
  // Warstwa izolacji czynna w oknie, wspólna dla wszystkich obszarów, wybierana pasem narzędzi ramy.
  warstwa: StanWarstwy;
  // Zasięg czynny okna — poziom i byt, na których reguła ma obowiązywać, wybierany panelem lewym.
  zasieg: StanZasiegu;
}

/** Kontrakt obszaru — kształt jednolity dla wszystkich plików obszarów tego okna Punktów Izolacji klienta. */
export interface ObszarIzolacji {
  /** Element montowany w ciele okna, gdy obszar jest czynny. */
  element: HTMLElement;
  /** Odświeża treść obszaru; wołane przy każdym pokazaniu. */
  odswiez(): void;
  /** Odłącza subskrypcje obszaru. Pominięte, gdy obszar żadnych nie zakłada. */
  zamknij?(): void;
}

/** Wpis rejestru: kod obszaru, etykieta widoczna nad obszarem, panel okna oraz jego własna wytwórnia obszaru. */
export interface WpisObszaru {
  kod: KodObszaruIzolacji;
  nazwa: string;
  opis: string;
  /** Kolumna okna, w której obszar stoi. */
  panel: PanelIzolacji;
  utworz(zaleznosci: ZaleznosciObszaru): ObszarIzolacji;
}

/** Rejestr sześciu obszarów rozłożonych na trzy panele okna: zasięg po lewej, macierz pośrodku, profil po prawej. */
export const REJESTR_OBSZAROW: readonly WpisObszaru[] = [
  {
    kod: 'poziomy',
    nazwa: 'Poziom zasięgu',
    opis: 'Poziom zapisu: okno → rola → sesja → projekt → para modułów → moduł → środowisko → globalny.',
    panel: 'zasieg',
    utworz: utworzObszarPoziomy,
  },
  {
    kod: 'osie',
    nazwa: 'Oś rozstrzygania',
    opis: 'Adresat zapisu w ramach poziomu: konto, model albo platforma.',
    panel: 'zasieg',
    utworz: utworzObszarOsie,
  },
  {
    kod: 'kontekst',
    nazwa: 'Izolacja kontekstu',
    opis: 'Historia wymiany, pamięć długoterminowa i bieżący stan roboczy — odrębne albo współdzielone.',
    panel: 'macierz',
    utworz: utworzObszarKontekst,
  },
  {
    kod: 'techniczny',
    nazwa: 'Izolacja techniczna',
    opis: 'Osiem punktów izolacji technicznej: katalog roboczy, środowisko, sieć, pliki, konto, model, serwer.',
    panel: 'macierz',
    utworz: utworzObszarTechniczny,
  },
  {
    kod: 'profile',
    nazwa: 'Profile',
    opis: 'Nazwane zestawy jedenastu punktów: zapis, wczytanie do podglądu, przypisanie do poziomu, usunięcie.',
    panel: 'profil',
    utworz: utworzObszarProfile,
  },
  {
    kod: 'efektywna',
    nazwa: 'Polityka efektywna',
    opis: 'Podgląd rozstrzygnięcia jedenastu punktów izolacji wraz z pochodzeniem każdej wartości.',
    panel: 'profil',
    utworz: utworzObszarEfektywna,
  },
];

/** Trzy panele okna Punktów Izolacji w kolejności kolumn, wraz z nagłówkiem i opisem każdej z nich osobno. */
export const PANELE_IZOLACJI: ReadonlyArray<{
  kod: PanelIzolacji;
  tytul: string;
  opis: string;
}> = [
  {
    kod: 'zasieg',
    tytul: 'Selektor zasięgu',
    opis: 'Poziom, na którym reguła izolacji ma obowiązywać, i oś, dla której obowiązuje w jego ramach.',
  },
  {
    kod: 'macierz',
    tytul: 'Macierz izolacji',
    opis:
      'Jedenaście punktów na wybranym zasięgu: trzy przełączniki kontekstu (współdzielone albo ' +
      'odrębne) i osiem zakresów technicznych (włączony albo wyłączony). Okno nigdy nie wymusza ' +
      'izolacji — udostępnia ją jako możliwość.',
  },
  {
    kod: 'profil',
    tytul: 'Profil i podgląd',
    opis: 'Zapis, wczytanie i przypisanie profilu izolacji oraz podgląd polityki efektywnej po dziedziczeniu.',
  },
];
