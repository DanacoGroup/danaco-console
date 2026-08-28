import type { DesignPrompt } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { KLASY_DYMKA } from './dymek-objasnienia';

/**
 * Suwaki parametrów generowania obejmują kreatywność, ziarno oraz liczbę
 * wariantów. Każdy z parametrów ma w kontrakcie zakres zamknięty albo naturalny,
 * a wartość bieżąca stoi w odczycie umieszczonym obok suwaka.
 */
export interface SuwakiPromptu {
  element: HTMLElement;
  /** Trzy parametry liczbowe gotowe do włożenia w prompt. */
  wartosci(): Pick<DesignPrompt, 'creativity' | 'seed' | 'variants'>;
  /** Nanosi parametry promptu na suwaki — droga powrotu z historii. */
  nanies(prompt: DesignPrompt): void;
}

/**
 * Opis suwaka wiąże klucz pola kontraktu z etykietą, zakresem, krokiem, wartością
 * wyjściową oraz treścią dymka objaśnienia. Wykaz opisów rozstrzyga skład
 * i kolejność suwaków w interfejsie.
 */
interface OpisSuwaka {
  klucz: 'creativity' | 'seed' | 'variants';
  etykieta: string;
  min: number;
  max: number;
  krok: number;
  wyjsciowa: number;
  objasnienie: string;
}

const SUWAKI: readonly OpisSuwaka[] = [
  {
    klucz: 'creativity',
    etykieta: 'Kreatywność wobec wierności',
    min: 0,
    max: 1,
    krok: 0.05,
    wyjsciowa: 0.5,
    objasnienie:
      'Pole creativity promptu, od 0 do 1. Niżej znaczy wierniej wobec promptu i obrazu referencyjnego.',
  },
  {
    klucz: 'seed',
    etykieta: 'Ziarno generowania',
    min: 0,
    max: 999999,
    krok: 1,
    wyjsciowa: 0,
    objasnienie: 'Pole seed. Ta sama wartość przy tym samym promptcie powtarza wynik; zero znaczy losowe.',
  },
  {
    klucz: 'variants',
    etykieta: 'Liczba wariantów (generowanie wsadowe)',
    min: 1,
    max: 8,
    krok: 1,
    wyjsciowa: 1,
    objasnienie: 'Pole variants. Wartość powyżej jedynki jest generowaniem wsadowym jednym zleceniem.',
  },
];

export function utworzSuwaki(): SuwakiPromptu {
  const kontrolki = new Map<OpisSuwaka['klucz'], HTMLInputElement>();

  const element = document.createElement('div');
  element.className = 'md-suwaki';

  for (const opis of SUWAKI) {
    const { wiersz, kontrolka } = zbudujSuwak(opis);
    kontrolki.set(opis.klucz, kontrolka);
    element.append(wiersz);
  }

  function liczba(klucz: OpisSuwaka['klucz']): number {
    return Number(kontrolki.get(klucz)?.value ?? 0);
  }

  return {
    element,

    wartosci() {
      const ziarno = liczba('seed');
      return {
        creativity: liczba('creativity'),
        variants: liczba('variants'),
        ...(ziarno > 0 ? { seed: ziarno } : {}),
      };
    },

    nanies(prompt) {
      for (const opis of SUWAKI) {
        const kontrolka = kontrolki.get(opis.klucz);
        if (kontrolka === undefined) continue;
        const wartosc = prompt[opis.klucz];
        kontrolka.value = String(typeof wartosc === 'number' ? wartosc : opis.wyjsciowa);
        kontrolka.dispatchEvent(new Event('input'));
      }
    },
  };
}

/**
 * Buduje wiersz suwaka wraz z etykietą, dymkiem objaśnienia, kontrolką zakresu
 * oraz odczytem wartości bieżącej. Zwraca wiersz do osadzenia i samą kontrolkę,
 * potrzebną przy odczycie oraz przy nanoszeniu wartości promptu.
 */
function zbudujSuwak(opis: OpisSuwaka): { wiersz: HTMLElement; kontrolka: HTMLInputElement } {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'range';
  kontrolka.className = 'dn-suwak';
  kontrolka.min = String(opis.min);
  kontrolka.max = String(opis.max);
  kontrolka.step = String(opis.krok);
  kontrolka.value = String(opis.wyjsciowa);
  kontrolka.setAttribute('aria-label', opis.etykieta);

  const odczyt = document.createElement('span');
  odczyt.className = 'dn-plakietka';
  odczyt.textContent = String(opis.wyjsciowa);

  /** Wypełnienie toru: żeton `--dn-suwak-pozycja` czytany przez `komponenty/suwak.css`. */
  const ustawWypelnienie = (): void => {
    const zakres = opis.max - opis.min;
    const udzial = zakres === 0 ? 0 : (Number(kontrolka.value) - opis.min) / zakres;
    kontrolka.style.setProperty('--dn-suwak-pozycja', `${Math.round(udzial * 100)}%`);
  };

  kontrolka.addEventListener('input', () => {
    odczyt.textContent = kontrolka.value;
    ustawWypelnienie();
  });
  ustawWypelnienie();

  const etykieta = document.createElement('span');
  etykieta.className = 'dn-pole-etykieta';
  etykieta.textContent = opis.etykieta;

  const naglowek = document.createElement('span');
  naglowek.className = 'md-pole__rzad';
  naglowek.append(etykieta, utworzDymekObjasnienia(opis.objasnienie, KLASY_DYMKA));

  const wiersz = document.createElement('div');
  wiersz.className = 'md-suwaki__wiersz';
  wiersz.dataset['parametr'] = opis.klucz;
  wiersz.append(naglowek, kontrolka, odczyt);
  return { wiersz, kontrolka };
}
