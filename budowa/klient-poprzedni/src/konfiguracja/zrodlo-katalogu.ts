import {
  Command,
  type SettingCategory,
  type SettingDefinition,
  type SettingsCategoryListRequest,
  type SettingsDefinitionListRequest,
} from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { czyLiczba, czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Katalog okna konfiguracji pobierany z rdzenia. Formularz jest budowany z danych:
 * klient nie zna żadnego klucza, kategorii ani etykiety pola, więc dodanie ustawienia
 * jest nowym wierszem katalogu, a nie zmianą kodu interfejsu.
 */
export interface ZrodloKatalogu {
  /** `settings.category.list` — kategorie w kolejności wyświetlania. */
  kategorie(zadanie?: SettingsCategoryListRequest): Promise<SettingCategory[]>;
  /** `settings.definition.list` — pozycje katalogu wraz z metadanymi pola. */
  definicje(zadanie?: SettingsDefinitionListRequest): Promise<SettingDefinition[]>;
}

/**
 * Powiadomienie o niepowodzeniu odczytu. Widok potrzebuje go, bo katalog pusty
 * po odmowie rdzenia i katalog pusty na świeżej instalacji to dwa różne stany.
 */
export type NaNiepowodzenie = (powod: string) => void;

export function utworzZrodloKatalogu(
  kanal: Kanal,
  naNiepowodzenie: NaNiepowodzenie = () => undefined,
): ZrodloKatalogu {
  return {
    async kategorie(zadanie = {}) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.SettingsCategoryList, zadanie),
        Command.SettingsCategoryList,
        (tresc) => czyTablica(tresc.categories),
      );
      if (!wynik.udany) {
        ostrzez(Command.SettingsCategoryList, wynik.blad?.message);
        naNiepowodzenie(powodOdczytu(Command.SettingsCategoryList, wynik.blad?.message));
        return [];
      }
      return uporzadkuj(wynik.wynik?.categories ?? []);
    },

    async definicje(zadanie = {}) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.SettingsDefinitionList, zadanie),
        Command.SettingsDefinitionList,
        (tresc) => czyTablica(tresc.definitions) && czyLiczba(tresc.total),
      );
      if (!wynik.udany) {
        ostrzez(Command.SettingsDefinitionList, wynik.blad?.message);
        naNiepowodzenie(powodOdczytu(Command.SettingsDefinitionList, wynik.blad?.message));
        return [];
      }
      return uporzadkuj(wynik.wynik?.definitions ?? []);
    },
  };
}

/**
 * Pozycja katalogu niosąca własną kolejność wyświetlania. Kolejność jest polem
 * odpowiedzi rdzenia, więc porządek pól formularza rozstrzyga się w katalogu,
 * a nie w kodzie okna.
 */
interface Uporzadkowana {
  order: number;
}

/**
 * Porządek wyświetlania bierzemy z katalogu, nie z kolejności odpowiedzi.
 * Wiersz bez kolejności trafia na koniec, zamiast zniknąć.
 */
function uporzadkuj<T extends Uporzadkowana>(pozycje: readonly T[]): T[] {
  return [...pozycje].sort((pierwsza, druga) => kolejnosc(pierwsza) - kolejnosc(druga));
}

function kolejnosc(pozycja: Uporzadkowana): number {
  return Number.isFinite(pozycja.order) ? pozycja.order : Number.MAX_SAFE_INTEGER;
}

/**
 * Niepowodzenie odczytu katalogu zostaje w dzienniku — wraz z nazwą komendy i powodem
 * podanym przez rdzeń. Wpis dziennika jest jedynym śladem odczytu, który nie doszedł,
 * bo okno pokazuje wtedy pusty katalog wraz z osobnym zdaniem.
 */
function ostrzez(komenda: string, powod: string | undefined): void {
  console.warn('[konfiguracja] katalog nie dotarł', komenda, powod ?? '');
}

/**
 * Zdanie o niepowodzeniu odczytu pokazywane w oknie. Gdy rdzeń nie podał
 * powodu, zdanie wskazuje przynajmniej komendę, która nie odpowiedziała.
 */
export function powodOdczytu(komenda: string, powod: string | undefined): string {
  const tresc = powod !== undefined && powod !== '' ? powod : 'rdzeń nie podał powodu';
  return `Komenda ${komenda} nie odpowiedziała: ${tresc}`;
}
