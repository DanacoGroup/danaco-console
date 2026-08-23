import { utworzRozwiniecie } from './warstwy-designu';

/**
 * Skróty klawiszowe modułu Design wraz z tym, które z nich moduł faktycznie
 * wiąże.
 *
 * Skąd biorą się kombinacje. Opracowanie modułu wymienia skrót klawiszowy jako
 * drogę równorzędną kliknięciu — przy adnotacjach na kanwie, przy powiększeniu
 * podglądu i przy całej warstwie czwartej — ale **konkretnych kombinacji nie
 * podaje**. Kombinacje są więc rozstrzygnięciem projektowym tej budowy, podjętym
 * wedle jednej reguły: bierzemy to, co ten produkt już związał w module
 * Translate, żeby Operator przechodzący między modułami nie uczył się dwóch
 * układów klawiatury. Czynności, których opracowanie skrótem nie opatruje,
 * skrótu tu nie dostają — wymyślanie ich byłoby dokładaniem funkcji.
 *
 * Nasłuch wisi na elemencie modułu, nie na dokumencie. Moduł znika z drzewa
 * przy zejściu ze sceny i nasłuch znika razem z nim; nasłuch dokumentu trzeba by
 * odpinać osobno, a pierwszy przeoczony byłby wyciekiem.
 *
 * Skutek uboczny tej decyzji jest nazwany, nie przemilczany: skrót działa, gdy
 * ognisko stoi wewnątrz modułu. Poza modułem klawisze należą do powłoki.
 *
 * Dwa skróty stoją w wykazie, a moduł ich nie wiąże. Powiększenie podglądu
 * należy do Preview Window, którego wytwórnia leży poza katalogiem tego modułu;
 * sterowanie pętlą należy do Execution Loop Window, okna wspólnego platformy.
 * Wykaz mówi to wprost, zamiast pomijać pozycje i sugerować, że skrótów nie ma.
 */

/** Jedna pozycja wykazu skrótów. */
export interface Skrot {
  /** Zapis klawiszy w notacji opracowania. */
  readonly klawisze: string;
  /** Co skrót robi. */
  readonly dziala: string;
  /** Okno, do którego skrót należy. */
  readonly okno: string;
  /** Czy moduł wiąże ten skrót u siebie. */
  readonly wiazany: boolean;
  /** Uwaga — obowiązkowa, gdy moduł skrótu nie wiąże albo wiąże go inaczej. */
  readonly uwaga?: string;
}

export const SKROTY: readonly Skrot[] = [
  {
    klawisze: 'Ctrl/Cmd + K',
    dziala: 'Wyszukiwarka funkcji modułu',
    okno: 'dowolne okno modułu',
    wiazany: true,
    uwaga:
      'Ten sam skrót prowadzi na pasku górnym do pola poleceń. W module pierwszeństwo ma ' +
      'wyszukiwarka funkcji; poza modułem skrót działa jak dotąd. Zbieżność jest realna ' +
      'i wymaga rozstrzygnięcia właściciela projektu.',
  },
  {
    klawisze: 'Ctrl/Cmd + Shift + A',
    dziala: 'Adnotacja warstwy zaznaczonej na kanwie',
    okno: 'Design Board',
    wiazany: true,
    uwaga:
      'Opracowanie wymienia dla adnotacji menu kebab i skrót klawiszowy. Skrót prowadzi ' +
      'ognisko do pola adnotacji pierwszej zaznaczonej warstwy; bez zaznaczenia mówi, ' +
      'że nie ma czego opisać.',
  },
  {
    klawisze: 'Ctrl/Cmd + Shift + T',
    dziala: 'Tokens & System Panel',
    okno: 'dowolne okno modułu',
    wiazany: true,
    uwaga:
      'Panel jest elementem warstwy trzeciej, wywoływanym menu kebab obszaru roboczego. ' +
      'Skrót jest drugą drogą do tego samego miejsca, nie osobną funkcją.',
  },
  {
    klawisze: 'Ctrl/Cmd + Shift + P',
    dziala: 'Powiększenie pełnoekranowe podglądu',
    okno: 'Preview Window',
    wiazany: false,
    uwaga:
      'Preview Window ma w katalogu rdzenia jedną definicję i dwa przypięcia — do Studia ' +
      'i do Designu — a jego wytwórnia leży w katalogu modułu Studio. Moduł Design skrótu ' +
      'nie wiąże, bo nie buduje okna, którego on dotyczy.',
  },
  {
    klawisze: 'Ctrl/Cmd + Shift + E',
    dziala: 'Wstrzymanie i wznowienie pętli wykonawczej',
    okno: 'Execution Loop Window',
    wiazany: false,
    uwaga: 'Okno wspólne platformy — moduł Design go nie buduje i skrótu nie wiąże.',
  },
];

/** Czynności wyzwalane skrótami — okna dostarczają je przy montażu modułu. */
export interface CzynnosciSkrotow {
  /** `Ctrl/Cmd + K`. */
  otworzWyszukiwarke(): void;
  /** `Ctrl/Cmd + Shift + A`. */
  opiszWarstwe(): void;
  /** `Ctrl/Cmd + Shift + T`. */
  otworzZetony(): void;
}

/** Podpina skróty do elementu modułu; zwraca odpięcie wołane przy rozłączeniu. */
export function podepnijSkroty(element: HTMLElement, czynnosci: CzynnosciSkrotow): () => void {
  function przyKlawiszu(zdarzenie: KeyboardEvent): void {
    const czynnosc = dopasuj(zdarzenie, czynnosci);
    if (czynnosc === null) return;
    // Zdarzenie zatrzymuje się na module, bo skrót został tu obsłużony. Bez tego
    // ta sama kombinacja wykonałaby się drugi raz na poziomie powłoki.
    zdarzenie.preventDefault();
    zdarzenie.stopPropagation();
    czynnosc();
  }

  element.addEventListener('keydown', przyKlawiszu);
  return () => element.removeEventListener('keydown', przyKlawiszu);
}

/**
 * Dopasowanie zdarzenia do czynności.
 *
 * `Ctrl` i `Cmd` są tu równoważne, tak jak w zapisie skrótów opracowania:
 * na komputerach Apple modyfikatorem polecenia jest klawisz `Meta`.
 */
function dopasuj(zdarzenie: KeyboardEvent, czynnosci: CzynnosciSkrotow): (() => void) | null {
  const polecenie = zdarzenie.ctrlKey || zdarzenie.metaKey;
  if (!polecenie) return null;
  const klawisz = zdarzenie.key.toLowerCase();

  if (zdarzenie.shiftKey) {
    if (klawisz === 'a') return czynnosci.opiszWarstwe;
    if (klawisz === 't') return czynnosci.otworzZetony;
    return null;
  }
  if (klawisz === 'k') return czynnosci.otworzWyszukiwarke;
  return null;
}

/** Wykaz skrótów jako element warstwy czwartej — droga bez klawiatury. */
export function utworzWykazSkrotow(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 4,
    nazwa: 'Skróty klawiszowe modułu',
    wyjasnienie:
      'Kombinacje są rozstrzygnięciem tej budowy — opracowanie wskazuje skrót jako drogę, ' +
      'nie podając klawiszy. Wykaz mówi też, których skrótów moduł nie wiąże i dlaczego.',
    znacznik: '☰',
  });

  const wykaz = document.createElement('ul');
  wykaz.className = 'md-skroty';
  wykaz.replaceChildren(...SKROTY.map(wierszSkrotu));
  rozwiniecie.tresc.append(wykaz);
  return rozwiniecie.element;
}

function wierszSkrotu(skrot: Skrot): HTMLElement {
  const klawisze = document.createElement('span');
  klawisze.className = 'md-skroty__klawisze';
  klawisze.textContent = skrot.klawisze;

  const dziala = document.createElement('span');
  dziala.className = 'md-skroty__dziala';
  dziala.textContent = `${skrot.dziala} — ${skrot.okno}`;

  const element = document.createElement('li');
  element.className = 'md-skroty__wiersz';
  element.dataset['wiazany'] = skrot.wiazany ? 'tak' : 'nie';
  element.append(klawisze, dziala);

  if (skrot.uwaga !== undefined) {
    const uwaga = document.createElement('p');
    uwaga.className = 'md-skroty__uwaga';
    uwaga.textContent = skrot.uwaga;
    element.append(uwaga);
  }
  return element;
}
