import { poleWyboru } from '../../modele/kontrolki-formularza';
import { oznaczWarstwe } from './warstwy-designu';

/**
 * Tryby narzędzia kanwy Design Board — znacznik `Tryb: … ▼` warstwy drugiej.
 *
 * Opracowanie wymienia siedem trybów: zaznaczanie, pióro, kształt, tekst,
 * pędzel-maska, retusz i ramka UI. Wykaz podaje wszystkie siedem, bo ukrycie
 * czterech, których ta budowa nie wykonuje, przedstawiłoby moduł jako mniejszy,
 * niż zaprojektowano.
 *
 * Żaden tryb nie jest wygaszony — wybór zawsze się udaje, a pod wykazem staje
 * zdanie mówiące, co tryb w tej budowie robi albo czego mu brakuje. To jest
 * różnica między „nie da się wybrać" a „wybrałeś i wiesz, co masz": pierwsze
 * niczego nie tłumaczy, drugie tłumaczy wszystko.
 *
 * Dwa tryby zmieniają zachowanie kanwy. Zaznaczanie jest stanem wyjściowym
 * i prowadzi wskazywanie oraz przeciąganie warstw. Ramka interfejsu kieruje
 * pracę do presetów ramek przybornika, bo ramka jest w tej budowie warstwą
 * o zadanych wymiarach, a nie osobnym bytem.
 *
 * Pięć pozostałych wymaga bytów, których kompozycja nie zna: ścieżki z węzłami,
 * kształtu innego niż prostokąt, tekstu, maski przezroczystości i pędzla
 * działającego na pikselach zasobu. Kompozycja niesie warstwę o położeniu,
 * rozmiarze, kolejności, blokadzie i adnotacji — i tyle jedzie do rdzenia.
 */

/** Tryb narzędzia kanwy; ta sama wartość idzie w `data-tryb` kanwy. */
export type TrybKanwy =
  | 'zaznaczanie'
  | 'pioro'
  | 'ksztalt'
  | 'tekst'
  | 'pedzel-maska'
  | 'retusz'
  | 'ramka-ui';

/** Opis trybu: nazwa dla Operatora i zdanie o tym, co tryb w tej budowie robi. */
interface OpisTrybu {
  readonly kod: TrybKanwy;
  readonly nazwa: string;
  /** Czy tryb zmienia zachowanie kanwy w tej budowie. */
  readonly czynny: boolean;
  readonly zdanie: string;
}

const TRYBY: readonly OpisTrybu[] = [
  {
    kod: 'zaznaczanie',
    nazwa: 'zaznaczanie',
    czynny: true,
    zdanie:
      'Wskazywanie i przeciąganie warstw na kanwie. Klawisz Ctrl albo Cmd przy naciśnięciu ' +
      'dokłada warstwę do zaznaczenia zamiast zastępować je; warstwa zablokowana nie rusza się ' +
      'z miejsca.',
  },
  {
    kod: 'pioro',
    nazwa: 'pióro (ścieżki Béziera)',
    czynny: false,
    zdanie:
      'Rysunek ścieżkami wymaga bytu ścieżki z węzłami. Warstwa kompozycji niesie położenie, ' +
      'rozmiar, kolejność, blokadę i adnotację — węzłów nie ma gdzie zapisać ani czym wysłać ' +
      'do rdzenia.',
  },
  {
    kod: 'ksztalt',
    nazwa: 'kształt',
    czynny: false,
    zdanie:
      'Elementy pomocnicze przybornika są prostokątnymi warstwami bez zasobu. Elipsy, ' +
      'wielokąta ani gwiazdy warstwa nie wyraża, bo nie ma pola opisującego kształt.',
  },
  {
    kod: 'tekst',
    nazwa: 'tekst',
    czynny: false,
    zdanie:
      'Adnotacja warstwy jest opisem warstwy, nie tekstem kompozycji: nie ma kroju, stopnia ' +
      'ani pozycji własnej. Bytu tekstu kompozycja nie zna.',
  },
  {
    kod: 'pedzel-maska',
    nazwa: 'pędzel-maska',
    czynny: false,
    zdanie:
      'Maska jest osobnym obrazem wiązanym z warstwą. Ani wiązania, ani obrazu maski kontrakt ' +
      'nie niesie, a bajtów zasobu przeglądarka i tak nie ma po czym pobrać.',
  },
  {
    kod: 'retusz',
    nazwa: 'retusz',
    czynny: false,
    zdanie:
      'Poprawka obrazu w kontrakcie obejmuje CAŁY obraz jednym natężeniem. Pędzla ani obszaru ' +
      'nie ma czym wskazać, więc retusz miejscowy nie ma drogi.',
  },
  {
    kod: 'ramka-ui',
    nazwa: 'ramka interfejsu',
    czynny: true,
    zdanie:
      'Ramka jest warstwą o zadanych wymiarach. Zakłada się ją presetem urządzenia ' +
      'w przyborniku kanwy, a wymiary jadą do rdzenia polami szerokości i wysokości warstwy.',
  },
];

export interface TrybyKanwy {
  element: HTMLElement;
  /** Tryb wybrany — kanwa znakuje się nim w `data-tryb`. */
  tryb(): TrybKanwy;
}

/**
 * Znacznik trybu kanwy.
 *
 * @param naZmiane wołane po każdym wyborze; okno przepisuje tryb na kanwę.
 */
export function utworzTrybyKanwy(naZmiane: (tryb: TrybKanwy) => void): TrybyKanwy {
  const wybor = poleWyboru(
    {
      etykieta: 'Tryb narzędzia kanwy',
      opis:
        'Siedem trybów z opracowania modułu. Dwa zmieniają zachowanie kanwy, pięć nazywa ' +
        'to, czego kompozycja nie wyraża — żaden nie jest zablokowany.',
    },
    TRYBY.map((opis) => ({
      wartosc: opis.kod,
      etykieta: opis.czynny ? opis.nazwa : `${opis.nazwa} — bez drogi w tej budowie`,
    })),
  );

  const zdanie = document.createElement('p');
  zdanie.className = 'dn-pole-opis md-tryby__zdanie';

  const element = document.createElement('div');
  element.className = 'md-tryby';
  oznaczWarstwe(element, 2);
  element.append(wybor.element, zdanie);

  function biezacy(): OpisTrybu {
    // Wartość pochodzi z wykazu złożonego z kodów trybów, więc pierwsza pozycja
    // jako zapas opisuje stan faktyczny kontrolki, a nie założenie o niej.
    return TRYBY.find((opis) => opis.kod === wybor.kontrolka.value) ?? (TRYBY[0] as OpisTrybu);
  }

  function przepisz(): void {
    const opis = biezacy();
    zdanie.textContent = opis.zdanie;
    element.dataset['czynny'] = String(opis.czynny);
    naZmiane(opis.kod);
  }

  wybor.kontrolka.addEventListener('change', przepisz);
  przepisz();

  return { element, tryb: () => biezacy().kod };
}
