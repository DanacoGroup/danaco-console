import { SettingValueType } from '../../../shared/contract';
import {
  utworzZbiornikZamiarow,
  type Kontrolka,
  type ZaleznosciKontrolki,
} from './kontrolka';

/**
 * Kontrolki wartości liczbowych: całkowitej i zmiennoprzecinkowej.
 *
 * Granice i skok pochodzą z katalogu (`minimum`, `maximum`, `step`), nie
 * z kodu klienta. Katalog, który granic nie podaje, daje pole bez granic —
 * brak metadanej nie jest błędem i niczego nie blokuje.
 *
 * Pole puste znaczy „bez wartości", a nie zero: `odczytaj` oddaje wtedy
 * `null`, więc zapis czyści wartość zamiast wpisywać liczbę, której Operator
 * nie podał.
 */
export function utworzKontrolkeLiczby(zaleznosci: ZaleznosciKontrolki): Kontrolka {
  const { definicja, identyfikator } = zaleznosci;
  const calkowita = definicja.valueType === SettingValueType.Int;
  const zbiornik = utworzZbiornikZamiarow();

  const pole = document.createElement('input');
  pole.id = identyfikator;
  pole.type = 'number';
  pole.className = 'dn-pole-kontrolka dk-kontrolka--liczba';
  pole.inputMode = calkowita ? 'numeric' : 'decimal';
  pole.step = String(definicja.step ?? (calkowita ? 1 : 'any'));
  if (definicja.minimum !== undefined) pole.min = String(definicja.minimum);
  if (definicja.maximum !== undefined) pole.max = String(definicja.maximum);
  if (definicja.placeholder !== undefined) pole.placeholder = definicja.placeholder;
  pole.addEventListener('change', () => zbiornik.zglos());

  return {
    element: pole,

    odczytaj: () => {
      if (pole.value.trim() === '') return null;
      const liczba = Number(pole.value);
      if (!Number.isFinite(liczba)) return null;
      return calkowita ? Math.trunc(liczba) : liczba;
    },

    ustaw: (wartosc) => {
      pole.value = tekstLiczby(wartosc);
    },

    naZatwierdzenie: zbiornik.naZatwierdzenie,
    ostrzezenie: opisGranic(definicja.minimum, definicja.maximum),
  };
}

/** Liczba w postaci tekstu pola; wartość nieliczbowa daje pole puste. */
function tekstLiczby(wartosc: unknown): string {
  if (wartosc === null || wartosc === undefined || wartosc === '') return '';
  const liczba = typeof wartosc === 'number' ? wartosc : Number(wartosc);
  return Number.isFinite(liczba) ? String(liczba) : '';
}

/** Zdanie o granicach drukowane pod polem; brak granic daje brak zdania. */
function opisGranic(dolna: number | undefined, gorna: number | undefined): string | undefined {
  if (dolna !== undefined && gorna !== undefined) return `Zakres od ${dolna} do ${gorna}.`;
  if (dolna !== undefined) return `Wartość nie mniejsza niż ${dolna}.`;
  if (gorna !== undefined) return `Wartość nie większa niż ${gorna}.`;
  return undefined;
}
