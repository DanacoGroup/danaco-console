import { SettingValueType } from '../../../shared/contract';
import {
  utworzZbiornikZamiarow,
  type Kontrolka,
  type ZaleznosciKontrolki,
} from './kontrolka';

/**
 * Kontrolki wartości liczbowych, całkowitej i zmiennoprzecinkowej, biorą granice
 * oraz skok z katalogu ustawień (`minimum`, `maximum`, `step`), a nie z kodu
 * klienta; katalog bez granic daje pole bez granic.
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

/**
 * Liczba w postaci tekstu pola: wartość pusta oraz nieliczbowa dają pole puste,
 * a wartość skończoną funkcja zapisuje jej postacią dziesiętną, więc pole nigdy
 * nie pokazuje zapisu, którego nie da się odczytać z powrotem jako liczby.
 */
function tekstLiczby(wartosc: unknown): string {
  if (wartosc === null || wartosc === undefined || wartosc === '') return '';
  const liczba = typeof wartosc === 'number' ? wartosc : Number(wartosc);
  return Number.isFinite(liczba) ? String(liczba) : '';
}

/**
 * Zdanie o granicach drukowane pod polem powstaje z metadanych katalogu: obie
 * granice dają zakres, jedna — warunek jednostronny, a brak obu daje brak
 * zdania, ponieważ pole bez granic niczego czytającemu nie obiecuje.
 */
function opisGranic(dolna: number | undefined, gorna: number | undefined): string | undefined {
  if (dolna !== undefined && gorna !== undefined) return `Zakres od ${dolna} do ${gorna}.`;
  if (dolna !== undefined) return `Wartość nie mniejsza niż ${dolna}.`;
  if (gorna !== undefined) return `Wartość nie większa niż ${gorna}.`;
  return undefined;
}
