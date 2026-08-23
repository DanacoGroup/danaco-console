import { SettingValueType, type SettingDefinition } from '../../../shared/contract';
import {
  bezpiecznyZapis,
  lista,
  napis,
  utworzZbiornikZamiarow,
  type Kontrolka,
  type ZaleznosciKontrolki,
} from './kontrolka';

/**
 * Kontrolki wartości tekstowych: napis, tekst wielowierszowy, ścieżka,
 * lista ścieżek, poświadczenie i wartość złożona w zapisie JSON.
 *
 * Zatwierdzenie następuje na zdarzeniu `change`, czyli po opuszczeniu pola —
 * nie po każdym znaku. Zapis co znak zasypałby rdzeń komendami `config.set`
 * i odbierał możliwość poprawienia wartości przed wysyłką.
 *
 * Wartość rodzaju `secret` nie wraca z rdzenia: pole pozostaje puste i mówi to
 * wprost. Kontrolka przyjmuje wartość nową, nie pokazuje wartości zapisanej.
 */

/** Kontrolka jednowierszowa: napis, ścieżka, poświadczenie. */
export function utworzKontrolkeNapisu(zaleznosci: ZaleznosciKontrolki): Kontrolka {
  const { definicja, identyfikator } = zaleznosci;
  const zbiornik = utworzZbiornikZamiarow();

  const pole = document.createElement('input');
  pole.id = identyfikator;
  pole.type = definicja.valueType === SettingValueType.Secret ? 'password' : 'text';
  pole.className = klasaPola(definicja);
  pole.placeholder = podpowiedz(definicja);
  if (definicja.pattern !== undefined && definicja.pattern !== '') pole.pattern = definicja.pattern;
  if (definicja.valueType === SettingValueType.Secret) pole.autocomplete = 'off';
  pole.addEventListener('change', () => zbiornik.zglos());

  const poswiadczenie = definicja.valueType === SettingValueType.Secret;

  return {
    element: pole,
    odczytaj: () => pole.value,
    ustaw: (wartosc) => {
      pole.value = poswiadczenie ? '' : napis(wartosc);
    },
    naZatwierdzenie: zbiornik.naZatwierdzenie,
    ostrzezenie: poswiadczenie
      ? 'Poświadczenie nie wraca z rdzenia — puste pole znaczy wartość niezmienioną.'
      : undefined,
  };
}

/** Kontrolka wielowierszowa dla tekstu swobodnego. */
export function utworzKontrolkeTekstu(zaleznosci: ZaleznosciKontrolki): Kontrolka {
  const obszar = obszarTekstu(zaleznosci);
  const zbiornik = utworzZbiornikZamiarow();
  obszar.addEventListener('change', () => zbiornik.zglos());

  return {
    element: obszar,
    odczytaj: () => obszar.value,
    ustaw: (wartosc) => {
      obszar.value = napis(wartosc);
    },
    naZatwierdzenie: zbiornik.naZatwierdzenie,
  };
}

/** Kontrolka listy ścieżek — jedna ścieżka w wierszu. */
export function utworzKontrolkeListySciezek(zaleznosci: ZaleznosciKontrolki): Kontrolka {
  const obszar = obszarTekstu(zaleznosci);
  obszar.classList.add('dk-kontrolka--mono');
  const zbiornik = utworzZbiornikZamiarow();
  obszar.addEventListener('change', () => zbiornik.zglos());

  return {
    element: obszar,
    odczytaj: () =>
      obszar.value
        .split('\n')
        .map((wiersz) => wiersz.trim())
        .filter((wiersz) => wiersz !== ''),
    ustaw: (wartosc) => {
      obszar.value = lista(wartosc).join('\n');
    },
    naZatwierdzenie: zbiornik.naZatwierdzenie,
    ostrzezenie: 'Jedna ścieżka w wierszu.',
  };
}

/**
 * Kontrolka wartości złożonej zapisanej w JSON.
 *
 * Zapis niepoprawny nie idzie do rdzenia: kontrolka pokazuje błąd i nie
 * zgłasza zamiaru. Pole pozostaje czynne — treść błędu nazywa powód.
 */
export function utworzKontrolkeJson(zaleznosci: ZaleznosciKontrolki): Kontrolka {
  const obszar = obszarTekstu(zaleznosci);
  obszar.classList.add('dk-kontrolka--mono');
  const zbiornik = utworzZbiornikZamiarow();

  const blad = document.createElement('p');
  blad.className = 'dn-pole-blad dk-kontrolka__blad';
  blad.hidden = true;

  const koszyk = document.createElement('div');
  koszyk.className = 'dk-kontrolka__koszyk';
  koszyk.append(obszar, blad);

  let ostatnia: unknown = undefined;

  obszar.addEventListener('change', () => {
    if (obszar.value.trim() === '') {
      ostatnia = undefined;
      pokazBlad(obszar, blad, '');
      zbiornik.zglos();
      return;
    }
    try {
      ostatnia = JSON.parse(obszar.value) as unknown;
      pokazBlad(obszar, blad, '');
      zbiornik.zglos();
    } catch (powod) {
      pokazBlad(obszar, blad, `Zapis nie jest poprawnym JSON: ${String(powod)}`);
    }
  });

  return {
    element: koszyk,
    odczytaj: () => ostatnia,
    ustaw: (wartosc) => {
      ostatnia = wartosc;
      obszar.value = wartosc === undefined ? '' : bezpiecznyZapis(wartosc);
      pokazBlad(obszar, blad, '');
    },
    naZatwierdzenie: zbiornik.naZatwierdzenie,
  };
}

/** Wspólny obszar tekstu obu kontrolek wielowierszowych. */
function obszarTekstu(zaleznosci: ZaleznosciKontrolki): HTMLTextAreaElement {
  const obszar = document.createElement('textarea');
  obszar.id = zaleznosci.identyfikator;
  obszar.className = klasaPola(zaleznosci.definicja);
  obszar.rows = 4;
  obszar.placeholder = podpowiedz(zaleznosci.definicja);
  return obszar;
}

/** Ścieżka i lista ścieżek dostają krój o stałej szerokości znaku. */
function klasaPola(definicja: SettingDefinition): string {
  const sciezkowa =
    definicja.valueType === SettingValueType.Path ||
    definicja.valueType === SettingValueType.PathList;
  return sciezkowa ? 'dn-pole-kontrolka dk-kontrolka--mono' : 'dn-pole-kontrolka';
}

/** Podpowiedź pustego pola pochodzi z katalogu, nie z kodu klienta. */
function podpowiedz(definicja: SettingDefinition): string {
  return definicja.placeholder ?? '';
}

/** Pokazuje albo chowa opis błędu wraz z oznaczeniem `aria-invalid`. */
function pokazBlad(obszar: HTMLTextAreaElement, opis: HTMLElement, tresc: string): void {
  opis.textContent = tresc;
  opis.hidden = tresc === '';
  obszar.setAttribute('aria-invalid', String(tresc !== ''));
}
