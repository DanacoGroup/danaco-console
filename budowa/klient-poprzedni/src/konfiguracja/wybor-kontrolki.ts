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
 * Jedyne miejsce, w którym rodzaj wartości zamienia się w kontrolkę.
 *
 * Formularz okna konfiguracji jest generowany z katalogu, więc rozdzielenie ma
 * dokładnie jeden punkt: nowy rodzaj wartości w kontrakcie to nowy przypadek
 * tutaj i nowy budowniczy obok — nie nowy ekran i nie zmiana w formularzu ani
 * w panelu kategorii.
 *
 * Rodzaj nieznany nie gasi pola. Rdzeń nowszy od klienta może przysłać rodzaj,
 * którego ten klient nie zna; zamiast pustego miejsca staje wtedy pole tekstowe
 * z ostrzeżeniem, a wartość pozostaje odczytywalna i zapisywalna jako napis.
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

/** Pole tekstowe z ostrzeżeniem — odpowiedź na rodzaj wartości spoza kontraktu. */
function kontrolkaZastepcza(zaleznosci: ZaleznosciKontrolki): Kontrolka {
  const rodzaj = String(zaleznosci.definicja.valueType);
  console.warn('[konfiguracja] rodzaj wartości nieznany klientowi', rodzaj);

  return {
    ...utworzKontrolkeNapisu(zaleznosci),
    ostrzezenie: `Rodzaj wartości „${rodzaj}" nie jest znany tej wersji interfejsu — pole przyjmuje wartość jako tekst.`,
  };
}
