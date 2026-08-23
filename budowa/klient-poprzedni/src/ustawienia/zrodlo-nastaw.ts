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

/**
 * Odczyt i zapis pojedynczej nastawy katalogu rdzenia — wspólna droga sekcji
 * Okna Ustawień.
 *
 * Osobne źródło obok `konfiguracja/zrodlo-wartosci.ts`, bo kształt pytania jest
 * inny. Tamto obsługuje Okno Konfiguracji: czyta poziom w całości (`config.get`
 * bez klucza), skleja łańcuch ośmiu zasięgów i rozstrzyga oś, bo tamto okno
 * rysuje każdą pozycję katalogu na każdym poziomie. Sekcja Ustawień pyta o jeden
 * klucz na jednym poziomie i potrzebuje przy tym definicji katalogu — opcje
 * wyboru przychodzą z katalogu, nie z wykazu zaszytego w kliencie. Komendy są te
 * same (`config.get`, `config.set`, `settings.definition.list`) i idą tym samym
 * `wywolaj` na tym samym kanale.
 *
 * Dwie właściwości rdzenia, na które ten plik jest przygotowany:
 *
 *  - zdarzenie `config.changed` o rodzaju `deleted` niesie w `entry` wartość
 *    zdjętą, nie nową; kto zastosuje `entry.value` bez patrzenia na `change`,
 *    przywróci właśnie skasowaną wartość — dlatego `naZmianeKlucza` oddaje przy
 *    usunięciu `undefined`, a nie treść wpisu;
 *  - rdzeń nie sprawdza poziomu zapisu wobec `allowedScopes` definicji
 *    i przyjmuje zapis na poziomie węższym, niż katalog dopuszcza; taki zapis
 *    wygrywa potem rozstrzyganie osi. Poziom bierzemy więc z `allowedScopes`
 *    definicji, nigdy z domysłu wołającego.
 */

/** Nastawa odczytana z rdzenia: definicja katalogu i wartość bieżąca. */
export interface OdczytNastawy {
  /** Pozycja katalogu; `undefined` znaczy „rdzeń tego klucza nie zna". */
  definicja?: SettingDefinition;
  /** Wpis wartości; `undefined` znaczy „nie zapisano, obowiązuje domyślna". */
  wpis?: ConfigEntry;
  /** Zdanie odmowy rdzenia, gdy odczyt się nie udał; puste znaczy „udał się". */
  odmowa: string;
}

/**
 * Poziom, na którym wolno zapisać klucz — wprost z definicji katalogu.
 *
 * Definicja niesie `allowedScopes` i to ona rozstrzyga; klient nie wybiera
 * poziomu za katalog, bo rdzeń takiego zapisu nie odrzuci i pomyłka byłaby
 * cicha. Definicja bez ani jednego poziomu oddaje `undefined` — wołający ma
 * wtedy powiedzieć, że nie wie, gdzie zapisać, zamiast zgadywać.
 */
export function poziomZapisu(definicja: SettingDefinition | undefined): ConfigScope | undefined {
  return definicja?.allowedScopes[0];
}

/**
 * Odczytuje definicję katalogu i wartość jednego klucza.
 *
 * Dwa wywołania, nie jedno, bo to dwie różne rzeczy: katalog mówi, co wolno
 * ustawić i z czego wybierać, a `config.get` mówi, co ustawiono. Idą po kolei,
 * bo poziom odczytu musi pochodzić z definicji.
 *
 * Odmowa nie wywraca odczytu: brak definicji albo brak wartości wraca jako pole
 * puste wraz ze zdaniem odmowy, a nie jako wyjątek.
 */
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

/** Zapisuje wartość klucza na poziomie wskazanym przez definicję katalogu. */
export function zapiszNastawe(
  kanal: Kanal,
  klucz: string,
  wartosc: unknown,
  poziom: ConfigScope,
): Promise<Wynik<ConfigSetResponse>> {
  return wywolaj(kanal, Command.ConfigSet, { key: klucz, value: wartosc, scope: poziom });
}

/**
 * Subskrypcja zmiany jednego klucza — droga na żywo z rdzenia do widoku.
 *
 * Słuchacz dostaje wartość po zmianie albo `undefined`, gdy wpis skasowano
 * (`config.reset`); przy usunięciu rdzeń niesie w `entry` wartość zdjętą.
 */
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

/** Zdanie odmowy rdzenia albo nazwanie milczenia; nigdy pustka. */
function zdanieOdmowy(wiadomosc: string | undefined): string {
  return wiadomosc === undefined || wiadomosc === ''
    ? 'Rdzeń odmówił odczytu bez podania powodu.'
    : wiadomosc;
}
