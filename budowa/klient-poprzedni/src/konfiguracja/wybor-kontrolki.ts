import { SettingValueType } from '../../../shared/contract';
import type { Kontrolka, ZaleznosciKontrolki } from './kontrolka';
import {
  utworzKontrolkeLiczby,
} from './kontrolki-liczbowe';
import {
  utworzKontrolkeJson,
  utworzKontrolkeListySciezek,
  utworzKontrolkeNapisu,
  utworzKontrolkeTekstu,
} from './kontrolki-tekstowe';
import {
  utworzKontrolkeListy,
  utworzKontrolkeListyWielokrotnej,
  utworzKontrolkePrzelacznika,
} from './kontrolki-wyboru';

/**
 * Jedyne miejsce, w którym rodzaj wartości z kontraktu zamienia się w kontrolkę
 * formularza konfiguracji. Nowy rodzaj wartości to nowy przypadek tutaj i nowy
 * budowniczy obok, bez zmiany formularza i panelu kategorii.
 */
export function utworzKontrolke(zaleznosci: ZaleznosciKontrolki): Kontrolka {
  switch (zaleznosci.definicja.valueType) {
    case SettingValueType.String:
    case SettingValueType.Path:
    case SettingValueType.Secret:
      return utworzKontrolkeNapisu(zaleznosci);

    case SettingValueType.Text:
      return utworzKontrolkeTekstu(zaleznosci);

    case SettingValueType.Json:
      return utworzKontrolkeJson(zaleznosci);

    case SettingValueType.PathList:
      return utworzKontrolkeListySciezek(zaleznosci);

    case SettingValueType.Int:
    case SettingValueType.Float:
      return utworzKontrolkeLiczby(zaleznosci);

    case SettingValueType.Bool:
      return utworzKontrolkePrzelacznika(zaleznosci);

    case SettingValueType.Enum:
      return utworzKontrolkeListy(zaleznosci);

    case SettingValueType.EnumList:
      return utworzKontrolkeListyWielokrotnej(zaleznosci);

    default:
      return kontrolkaZastepcza(zaleznosci);
  }
}

/**
 * Pole tekstowe z ostrzeżeniem, stawiane w odpowiedzi na rodzaj wartości spoza
 * kontraktu znanego tej wersji interfejsu; wartość pozostaje odczytywalna
 * i zapisywalna jako napis, a rodzaj trafia do dziennika przeglądarki.
 */
function kontrolkaZastepcza(zaleznosci: ZaleznosciKontrolki): Kontrolka {
  const rodzaj = String(zaleznosci.definicja.valueType);
  console.warn('[konfiguracja] rodzaj wartości nieznany klientowi', rodzaj);

  return {
    ...utworzKontrolkeNapisu(zaleznosci),
    ostrzezenie: `Rodzaj wartości „${rodzaj}" nie jest znany tej wersji interfejsu — pole przyjmuje wartość jako tekst.`,
  };
}
