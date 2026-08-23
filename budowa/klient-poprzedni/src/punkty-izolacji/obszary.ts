/**
 * Rejestr sześciu obszarów okna Punktów Izolacji wraz z kontraktem wpięcia.
 *
 * Podział na obszary idzie za maszynerią rdzenia
 * (`server/internal/konfig/definicje_izolacji.go`, `osie.go`, `poziomy.go`,
 * `polityka_efektywna.go`; egzekutor w `server/internal/session/izolacja*.go`
 * i `server/internal/core/izolacja*.go`), która rozstrzyga punkt izolacji na
 * dwóch prostopadłych osiach:
 *
 *   - oś rozstrzygania (`konfig.Os`, `osie.go`): konto → model → platforma —
 *     dla czego wartość obowiązuje;
 *   - poziom zasięgu (`konfig.Poziom`, `poziomy.go`): okno → rola → sesja →
 *     projekt → para modułów → moduł → środowisko → globalny — jak wąsko
 *     wartość obowiązuje. Poziom rozstrzyga pierwszy, oś dopiero w jego ramach
 *     (`osieOdNajwezszej` w `osie.go`).
 *
 * Na te dwie osie nakłada się przedmiot izolacji — jedenaście kluczy
 * z `definicje_izolacji.go`, w dwóch grupach:
 *
 *   - kontekst (trzy klucze, wartość `odrebna`/`wspoldzielona`, domyślnie
 *     odrębna): historia wymiany, pamięć długoterminowa, bieżący stan
 *     roboczy (pliki, projekt, załączniki, zmienne);
 *   - zakres techniczny (osiem kluczy, wartość `wlaczony`/`wylaczony`,
 *     domyślnie wyłączony — stan wyjściowy platformy to pełna swoboda
 *     operacyjna): katalog roboczy sesji, środowisko procesu, katalog danych
 *     modelu, dostęp sieciowy, odczyt/zapis plików, konto i token, model
 *     procesu, serwer wykonania. Egzekutor (`session/izolacja.go`) bierze
 *     wartość włączoną i odrzuca wykonanie, które by ją naruszyło; wartość
 *     wyłączona niczego nie ogranicza.
 *
 * `PolitykaEfektywna` (`polityka_efektywna.go`) składa te ustawienia w jeden
 * rozstrzygnięty podgląd na dany kontekst: dla każdego klucza wynik i poziom,
 * z którego pochodzi.
 *
 * Osobnym pojęciem maszynerii jest profil: nazwany zestaw trzech przełączników
 * kontekstu i ośmiu zakresów technicznych, niezależny od przypisania — ten sam
 * profil bywa przypisany kilku sesjom, rolom, projektom i oknom.
 *
 * Stąd sześć obszarów tego okna, rozłożonych na trzy panele wymagane przez
 * rozdz. 6.1 Modelu konfiguracji — selektor zasięgu, macierz izolacji, profil
 * i podgląd polityki efektywnej:
 *
 *   panel „selektor zasięgu" (lewy)
 *   1. `poziomy`     — poziom zasięgu (okno … globalny), na którym reguła obowiązuje.
 *   2. `osie`        — oś rozstrzygania (konto / model / platforma) w ramach poziomu.
 *
 *   panel „macierz izolacji" (środkowy)
 *   3. `kontekst`    — trzy punkty izolacji kontekstu (historia, pamięć, stan roboczy).
 *   4. `techniczny`  — osiem punktów izolacji zakresu technicznego (pliki, sieć, konto…).
 *
 *   panel „profil i podgląd" (prawy)
 *   5. `profile`     — nazwane zestawy jedenastu punktów: zapis, wczytanie,
 *                       przypisanie do poziomu, usunięcie.
 *   6. `efektywna`   — podgląd polityki efektywnej: rozstrzygnięcie i pochodzenie
 *                       każdego z jedenastu punktów dla wybranego kontekstu.
 *
 * Kontrakt wpięcia: każdy plik obszaru eksportuje dokładnie jedną funkcję
 * wytwórczą o kształcie
 *
 *   export function utworzObszar(zaleznosci: ZaleznosciObszaru): ObszarIzolacji
 *
 * i nic więcej publicznego poza własnymi typami pomocniczymi. `ObszarIzolacji`
 * i `ZaleznosciObszaru` są zdefiniowane tutaj — obszar ich nie definiuje
 * ponownie. Rama (`okno-punktow-izolacji.ts`) montuje `element` każdego obszaru
 * w jego panelu i woła `odswiez()` przy otwarciu okna oraz po każdej zmianie
 * zasięgu, warstwy i po zdarzeniu rdzenia; `zamknij()` przy demontażu okna,
 * jeśli obszar go zwrócił.
 *
 * Wszystkie dwanaście komend `isolation.*` stoi w `shared/contract.ts`
 * (`isolation.scope.list`, `context.get/set`, `technical.get/set`,
 * `profile.save/list/load/assign/delete`, `layer.set`, `policy.preview`),
 * a rdzeń ma je zarejestrowane (`core/kompozycja.go`). Obszary wołają je
 * naprawdę; stan „rdzeń nie niesie" zostaje wyłącznie tam, gdzie czynności
 * brakuje w kontrakcie — dziś jest to zapis wartości na wskazanej osi
 * rozstrzygania (obszar `osie`).
 */
import type { Kanal } from '../protokol/kanal';
import { utworzObszar as utworzObszarEfektywna } from './obszar-efektywna';
import { utworzObszar as utworzObszarKontekst } from './obszar-kontekst';
import { utworzObszar as utworzObszarOsie } from './obszar-osie';
import { utworzObszar as utworzObszarPoziomy } from './obszar-poziomy';
import { utworzObszar as utworzObszarProfile } from './obszar-profile';
import { utworzObszar as utworzObszarTechniczny } from './obszar-techniczny';
import type { StanWarstwy } from './stan-warstwy';
import type { StanZasiegu } from './stan-zasiegu';

/** Kod jednego z sześciu obszarów okna. */
export type KodObszaruIzolacji = 'kontekst' | 'techniczny' | 'profile' | 'osie' | 'poziomy' | 'efektywna';

/** Panel okna, w którym obszar stoi — trzy kolumny rozdz. 6.1 Modelu konfiguracji. */
export type PanelIzolacji = 'zasieg' | 'macierz' | 'profil';

/** Zależności dostępne każdemu obszarowi przy budowie. */
export interface ZaleznosciObszaru {
  /** Kanał do rdzenia — obszar woła nim komendy `isolation.*` i zakłada własne subskrypcje. */
  kanal: Kanal;
  /**
   * Warstwa izolacji czynna w oknie (`default` / `session`), wspólna dla
   * wszystkich obszarów. Wybiera ją pas narzędzi ramy (`sterowanie-warstwa.ts`),
   * bo dotyczy każdego odczytu i zapisu, nie jednego panelu.
   */
  warstwa: StanWarstwy;
  /**
   * Zasięg czynny okna — poziom i byt, na których reguła ma obowiązywać.
   * Wybiera go panel lewy (`panel-zasiegu.ts`); macierz izolacji, przypisanie
   * profilu i podgląd polityki czytają stąd, zamiast trzymać własny poziom.
   */
  zasieg: StanZasiegu;
}

/** Kontrakt obszaru — kształt jednolity dla wszystkich plików obszarów. */
export interface ObszarIzolacji {
  /** Element montowany w ciele okna, gdy obszar jest czynny. */
  element: HTMLElement;
  /** Odświeża treść obszaru; wołane przy każdym pokazaniu. */
  odswiez(): void;
  /** Odłącza subskrypcje obszaru. Pominięte, gdy obszar żadnych nie zakłada. */
  zamknij?(): void;
}

/** Wpis rejestru: kod, etykieta widoczna nad obszarem, panel i wytwórnia obszaru. */
export interface WpisObszaru {
  kod: KodObszaruIzolacji;
  nazwa: string;
  opis: string;
  /** Kolumna okna, w której obszar stoi. */
  panel: PanelIzolacji;
  utworz(zaleznosci: ZaleznosciObszaru): ObszarIzolacji;
}

/**
 * Rejestr sześciu obszarów rozłożonych na trzy panele okna, zgodnie z rozdz. 6.1
 * Modelu konfiguracji: selektor zasięgu po lewej, macierz izolacji pośrodku,
 * profil i podgląd polityki efektywnej po prawej.
 *
 * Kolejność wewnątrz panelu idzie za porządkiem rozstrzygania maszynerii:
 * w panelu zasięgu najpierw poziom (rozstrzyga pierwszy), pod nim oś; w macierzy
 * najpierw kontekst, potem zakres techniczny; w panelu prawym najpierw profile —
 * bo profil jest nazwanym zestawem tego, co stoi w macierzy — a pod nimi
 * podgląd wypadkowy.
 */
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

/** Trzy panele okna w kolejności kolumn, wraz z nagłówkiem każdej z nich. */
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
