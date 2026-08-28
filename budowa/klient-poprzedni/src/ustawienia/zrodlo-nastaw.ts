import {
  Command,
  EventType,
  type ConfigEntry,
  type ConfigScope,
  type ConfigSetResponse,
  type SettingDefinition,
} from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import { wywolaj } from '../protokol/wywolanie';

/** Odczyt i zapis pojedynczej nastawy katalogu rdzenia jest wspólną drogą sekcji okna ustawień, osobną od źródła wartości okna konfiguracji, bo kształt pytania jest inny. */
/** Nastawa odczytana z rdzenia: definicja katalogu i wartość bieżąca. */
export interface OdczytNastawy {
  /** Pozycja katalogu; `undefined` znaczy „rdzeń tego klucza nie zna". */
  definicja?: SettingDefinition;
  /** Wpis wartości; `undefined` znaczy „nie zapisano, obowiązuje domyślna". */
  wpis?: ConfigEntry;
  /** Zdanie odmowy rdzenia, gdy odczyt się nie udał; puste znaczy „udał się". */
  odmowa: string;
}

/** Poziom, na którym wolno zapisać klucz, wprost z definicji katalogu: definicja bez ani jednego poziomu oddaje brak, zamiast zgadywać. */
export function poziomZapisu(definicja: SettingDefinition | undefined): ConfigScope | undefined {
  return definicja?.allowedScopes[0];
}

/** Odczytuje definicję katalogu i wartość jednego klucza, dwoma wywołaniami po kolei, bo poziom odczytu musi pochodzić z definicji. */
export async function odczytajNastawe(kanal: Kanal, klucz: string): Promise<OdczytNastawy> {
  const katalog = await wywolaj(kanal, Command.SettingsDefinitionList, {
    key: klucz,
    includeDisabled: true,
  });

  if (!katalog.udany || katalog.wynik === undefined) {
    return { odmowa: zdanieOdmowy(katalog.blad?.message) };
  }

  const definicja = katalog.wynik.definitions.find((pozycja) => pozycja.key === klucz);
  if (definicja === undefined) {
    return {
      odmowa:
        `Katalog ustawień rdzenia nie zna klucza „${klucz}" — nie ma czego pokazać ` +
        'ani czym sterować.',
    };
  }

  const poziom = poziomZapisu(definicja);
  if (poziom === undefined) {
    return {
      definicja,
      odmowa:
        `Definicja klucza „${klucz}" nie wskazuje ani jednego dozwolonego poziomu ` +
        'zasięgu — nie wiadomo, skąd czytać ani gdzie zapisać.',
    };
  }

  const wartosc = await wywolaj(kanal, Command.ConfigGet, { key: klucz, scope: poziom });
  if (!wartosc.udany || wartosc.wynik === undefined) {
    return { definicja, odmowa: zdanieOdmowy(wartosc.blad?.message) };
  }

  const wpis = wartosc.wynik.entries.find((pozycja) => pozycja.key === klucz);
  return wpis === undefined ? { definicja, odmowa: '' } : { definicja, wpis, odmowa: '' };
}

/** Zapisuje wartość klucza na poziomie wskazanym przez definicję katalogu, jedną komendą kontraktu do rdzenia. */
export function zapiszNastawe(
  kanal: Kanal,
  klucz: string,
  wartosc: unknown,
  poziom: ConfigScope,
): Promise<Wynik<ConfigSetResponse>> {
  return wywolaj(kanal, Command.ConfigSet, { key: klucz, value: wartosc, scope: poziom });
}

/** Subskrypcja zmiany jednego klucza jest drogą na żywo z rdzenia do widoku; słuchacz dostaje wartość po zmianie albo brak, gdy wpis skasowano. */
export function naZmianeKlucza(
  kanal: Kanal,
  klucz: string,
  sluchacz: (wartosc: unknown | undefined, wpis: ConfigEntry) => void,
): Odsubskrybuj {
  return kanal.naZdarzenie(EventType.ConfigChanged, (tresc) => {
    if (tresc.entry.key !== klucz) return;
    sluchacz(tresc.change === 'deleted' ? undefined : tresc.entry.value, tresc.entry);
  });
}

/** Zdanie odmowy rdzenia albo nazwanie milczenia rdzenia jednym stałym zdaniem zastępczym, nigdy pustym napisem. */
function zdanieOdmowy(wiadomosc: string | undefined): string {
  return wiadomosc === undefined || wiadomosc === ''
    ? 'Rdzeń odmówił odczytu bez podania powodu.'
    : wiadomosc;
}
