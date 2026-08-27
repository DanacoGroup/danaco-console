import { ComponentKind } from '../../../shared/contract';
import type { NazwaIkony } from '../ikony/ikony';

/** Cztery komponenty własne jako pozycje strefy drugiej strony głównej niosą wyłącznie treść kafli, bez elementów i bez stylu, z kodem zamkniętym w czterech rodzajach kontraktu. */
export type KodKomponentu = ComponentKind;

/** Metadane komponentu widoczne na kaflu pokazują wyłącznie to, co rdzeń oddaje, bo typ komponentu nie niesie ani autora, ani ikony własnej, ani zasięgu. */
export interface MetadaneKomponentu {
  /** Czy komponent jest czynny (`Component.enabled`). */
  czynny: boolean;
  /** Czas ostatniej zmiany w milisekundach epoki (`Component.updatedAt`). */
  zmieniony: number;
  /** Byt magazynu modułowego, który komponent reprezentuje (`Component.targetId`). */
  cel?: string;
}

export interface PozycjaKomponentu {
  /** Kod niesiony przez kafel jest napisem, nie wyliczeniem: kod spoza wyliczenia to normalna odpowiedź. */
  kod: string;
  /** Nazwa kanoniczna komponentu — etykieta kafla. */
  nazwa: string;
  /** Wezwanie do działania: co Operator zrobi po naciśnięciu. */
  wezwanie: string;
  /** Ikona kafla, renderowana w skali mniejszej niż godło środowiska. */
  ikona: NazwaIkony;
  /** Metadane komponentu zbudowanego, wyłącznie dla kafli personalizowanych. */
  metadane?: MetadaneKomponentu;
  /** Identyfikator komponentu operatora, wyłącznie dla kafli personalizowanych. */
  komponent?: string;
}

/** Treść kafla bez kodu — kod dokłada wykaz z kontraktu, tak samo jak dla pozycji środowisk strony głównej. */
type TrescKomponentu = Omit<PozycjaKomponentu, 'kod'>;

/** Rekord wymusza komplet: rodzaj dołożony do wyliczenia bez treści kafla jest błędem kompilacji, nie kaflem, który po cichu nie powstał. */
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

/** Wykaz zastany strefy drugiej niesie treść kafli, dopóki rdzeń nie odpowie; kolejność jest kolejnością kontraktu, klient jej nie układa po swojemu. */
export const POZYCJE_KOMPONENTOW: readonly PozycjaKomponentu[] = Object.values(
  ComponentKind,
).map((kod) => ({ kod, ...TRESCI[kod] }));

/** Ikona rodzaju komponentu, którą kafel personalizowany dziedziczy zawsze po swoim rodzaju bazowym w wykazie. */
export function ikonaRodzaju(kod: KodKomponentu): NazwaIkony {
  return TRESCI[kod].ikona;
}
