import { SettingValueType, type SettingDefinition } from '../../../shared/contract';
import {
  bezpiecznyZapis,
  lista,
  napis,
  utworzZbiornikZamiarow,
  type Kontrolka,
  type ZaleznosciKontrolki,
} from './kontrolka';

/** Kontrolki wartości tekstowych: napis, tekst wielowierszowy, ścieżka, poświadczenie i zapis JSON. */

/** Kontrolka jednowierszowa dla wartości tekstowej: napis zwykły, ścieżka pliku albo poświadczenie ukryte w polu hasła. */
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

/** Kontrolka wielowierszowa dla tekstu swobodnego, oparta na wspólnym obszarze tekstu obu kontrolek wielowierszowych. */
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

/** Kontrolka listy ścieżek — jedna ścieżka w wierszu obszaru tekstu, rozdzielona znakiem końca każdego wiersza. */
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

/** Wspólny obszar tekstu obu kontrolek wielowierszowych — tekstu swobodnego oraz listy ścieżek w wierszach. */
function obszarTekstu(zaleznosci: ZaleznosciKontrolki): HTMLTextAreaElement {
  const obszar = document.createElement('textarea');
  obszar.id = zaleznosci.identyfikator;
  obszar.className = klasaPola(zaleznosci.definicja);
  obszar.rows = 4;
  obszar.placeholder = podpowiedz(zaleznosci.definicja);
  return obszar;
}

/** Ścieżka i lista ścieżek dostają krój pisma o stałej szerokości znaku, czytelny dla ciągów systemowych. */
function klasaPola(definicja: SettingDefinition): string {
  const sciezkowa =
    definicja.valueType === SettingValueType.Path ||
    definicja.valueType === SettingValueType.PathList;
  return sciezkowa ? 'dn-pole-kontrolka dk-kontrolka--mono' : 'dn-pole-kontrolka';
}

/** Podpowiedź pustego pola pochodzi z katalogu ustawień, nie jest wpisana wprost w kodzie warstwy klienta. */
function podpowiedz(definicja: SettingDefinition): string {
  return definicja.placeholder ?? '';
}

/** Pokazuje albo chowa opis błędu pola wraz z oznaczeniem dostępności `aria-invalid` na obszarze tekstu. */
function pokazBlad(obszar: HTMLTextAreaElement, opis: HTMLElement, tresc: string): void {
  opis.textContent = tresc;
  opis.hidden = tresc === '';
  obszar.setAttribute('aria-invalid', String(tresc !== ''));
}
