import { ComponentKind } from '../../../shared/contract';
import type { NazwaIkony } from '../ikony/ikony';

/**
 * Cztery komponenty własne jako pozycje strefy drugiej strony głównej.
 *
 * Jedna odpowiedzialność: treść kafli komponentów. Bez elementów, bez stylu.
 *
 * Różnica wobec strefy pierwszej: karta środowiska prowadzi do przestrzeni
 * pracy, a kafel komponentu do zbudowania rzeczy, która w tej przestrzeni potem
 * pracuje. Różnicę niesie krój — nagłówkowy dla wejścia do środowiska, bazowy
 * półgruby dla zbudowania komponentu.
 *
 * Kody pochodzą z kontraktu: `ComponentKind` jest zamkniętym zbiorem czterech
 * rodzajów komponentu własnego tej strefy. Rodzaj modułu i rodzaj komponentu
 * własnego to dwa różne byty, więc kod spoza kontraktu jest tu błędem
 * kompilacji, nie brakującym kaflem.
 */

/** Kod komponentu własnego przekazywany w zdarzeniu wyboru. */
export type KodKomponentu = ComponentKind;

/**
 * Metadane komponentu widoczne na kaflu.
 *
 * Kafel pokazuje wyłącznie to, co rdzeń oddaje. Typ `Component` nie niesie ani
 * autora, ani ikony własnej, ani zasięgu — więc kafel o nich nie mówi. Czas
 * przychodzi w milisekundach epoki; sformatowanie należy do widoku, bo rdzeń
 * nie zna strefy czasowej Operatora.
 */
export interface MetadaneKomponentu {
  /** Czy komponent jest czynny (`Component.enabled`). */
  czynny: boolean;
  /** Czas ostatniej zmiany w milisekundach epoki (`Component.updatedAt`). */
  zmieniony: number;
  /** Byt magazynu modułowego, który komponent reprezentuje (`Component.targetId`). */
  cel?: string;
}

export interface PozycjaKomponentu {
  /**
   * Kod niesiony przez kafel — `modul.kod` z rdzenia albo rodzaj komponentu.
   *
   * Typ jest napisem, nie wyliczeniem, z rozmysłu. O tym, co nastawia się na
   * Stronie głównej, rozstrzyga kolumna `modul.konfigurowany_na_stronie_glownej`
   * wystawiana w kontrakcie jako `Module.configuredOnHome` — kolejny moduł
   * tak oznaczony ma dostać kafel bez zmiany w kliencie. Kod spoza wyliczenia
   * jest tu więc normalną odpowiedzią rdzenia, a nie błędem.
   *
   * Zamkniętego zbioru `ComponentKind` pilnują nadal `TRESCI`
   * i `RODZAJE_DO_ZALOZENIA`, bo `component.create` przyjmuje wyłącznie rodzaj
   * z wyliczenia. Kafel i zakładanie to dwie różne rzeczy.
   */
  kod: string;
  /** Nazwa kanoniczna komponentu — etykieta kafla. */
  nazwa: string;
  /** Wezwanie do działania: co Operator zrobi po naciśnięciu. */
  wezwanie: string;
  /** Ikona kafla, renderowana w skali mniejszej niż godło środowiska. */
  ikona: NazwaIkony;
  /**
   * Metadane komponentu zbudowanego — wyłącznie dla kafli personalizowanych.
   *
   * Kafel rodzaju ich nie ma i mieć nie może: rodzaj nie jest bytem w bazie,
   * więc nie ma daty założenia ani stanu czynności. Pole nieobecne znaczy
   * „kafel rodzaju", nie „metadanych nie odczytano".
   */
  metadane?: MetadaneKomponentu;
  /**
   * Identyfikator komponentu Operatora — wyłącznie dla kafli personalizowanych.
   *
   * Kafel rodzaju (jeden z czterech stałych) prowadzi do zbudowania komponentu
   * i tego pola nie ma. Kafel personalizowany wskazuje komponent już
   * zbudowany, więc niesie jego identyfikator — po nim rozpoznaje się jeden
   * kafel spośród wielu tego samego rodzaju.
   */
  komponent?: string;
}

/** Treść kafla bez kodu — kod dokłada wykaz z kontraktu. */
type TrescKomponentu = Omit<PozycjaKomponentu, 'kod'>;

/**
 * Rekord wymusza komplet: rodzaj dołożony do `ComponentKind` bez treści kafla
 * jest błędem kompilacji, a nie kaflem, który po cichu nie powstał.
 */
const TRESCI: Record<KodKomponentu, TrescKomponentu> = {
  automations: {
    nazwa: 'Automations',
    wezwanie: 'Zbuduj automatykę',
    ikona: 'odswiez',
  },
  agents: {
    nazwa: 'Agents',
    wezwanie: 'Skonfiguruj agenta',
    ikona: 'tarcza',
  },
  workspace: {
    nazwa: 'Workspace',
    wezwanie: 'Załóż projekt',
    ikona: 'folder',
  },
  assistant: {
    nazwa: 'Assistant',
    wezwanie: 'Ustaw profil asystenta',
    ikona: 'uzytkownik',
  },
};

/**
 * Wykaz zastany strefy drugiej — treść kafli, dopóki rdzeń nie odpowie.
 *
 * To wartość początkowa, nie źródło prawdy; ten sam wzorzec co
 * `POZYCJE_SRODOWISK` w `pozycje-srodowisk.ts`. Pusty ekran przed pierwszą
 * odpowiedzią byłby gorszy niż cztery kafle, które i tak zostaną przerysowane.
 * Po odpowiedzi `module.list` rozstrzyga rdzeń —
 * `pozycjeModulowStrefyDrugiej` w `pozycje-modulow.ts`.
 *
 * Kolejność jest kolejnością kontraktu — klient nie układa jej po swojemu.
 * Nazwy i wezwania zostają miejscowe wyłącznie dla tej jednej klatki:
 * `ComponentKind` niesie kod rodzaju, a nie napis na kaflu.
 */
export const POZYCJE_KOMPONENTOW: readonly PozycjaKomponentu[] = Object.values(
  ComponentKind,
).map((kod) => ({ kod, ...TRESCI[kod] }));

/** Ikona rodzaju — kafel personalizowany dziedziczy ją po swoim rodzaju. */
export function ikonaRodzaju(kod: KodKomponentu): NazwaIkony {
  return TRESCI[kod].ikona;
}
