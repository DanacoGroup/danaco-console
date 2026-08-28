import { utworzRozwiniecie } from './warstwy-translate';

/** Skróty klawiszowe modułu Translate stanowią wykaz z opracowania, z których część moduł wiąże. */

/** Jedna pozycja wykazu skrótów niesie zapis klawiszy, opis działania, okno, do którego skrót należy, oraz stan wiązania w module. */
export interface Skrot {
  /** Zapis klawiszy tak, jak podaje go opracowanie. */
  readonly klawisze: string;
  /** Co skrót robi. */
  readonly dziala: string;
  /** Okno, do którego skrót należy. */
  readonly okno: string;
  /** Czy moduł wiąże ten skrót u siebie. */
  readonly wiazany: boolean;
  /** Uwaga do skrótu — obowiązkowa, gdy moduł go nie wiąże albo wiąże inaczej. */
  readonly uwaga?: string;
}

export const SKROTY: readonly Skrot[] = [
  {
    klawisze: 'Ctrl/Cmd + Shift + L',
    dziala: 'Dodanie nowego języka docelowego',
    okno: 'Translation Panels',
    wiazany: true,
  },
  {
    klawisze: 'Alt + ↑ / Alt + ↓',
    dziala: 'Przejście do poprzedniego lub następnego panelu wymagającego uwagi',
    okno: 'Translation Panels',
    wiazany: true,
    uwaga:
      'Opracowanie mówi o segmencie wymagającym rewizji. Kontrakt nie zna stanu pojedynczego ' +
      'segmentu — stan niesie cały panel — więc skrót prowadzi między panelami, których stan ' +
      'nie jest „gotowe".',
  },
  {
    klawisze: 'Ctrl/Cmd + G',
    dziala: 'Otwarcie okna Glossary Manager',
    okno: 'dowolne okno modułu',
    wiazany: true,
  },
  {
    klawisze: 'Ctrl/Cmd + M',
    dziala: 'Otwarcie okna Translation Memory Panel',
    okno: 'dowolne okno modułu',
    wiazany: true,
  },
  {
    klawisze: 'Ctrl/Cmd + Shift + Q',
    dziala: 'Otwarcie okna QA & Review Center',
    okno: 'dowolne okno modułu',
    wiazany: true,
  },
  {
    klawisze: 'Ctrl/Cmd + Shift + E',
    dziala: 'Wstrzymanie i wznowienie pętli wykonawczej',
    okno: 'Execution Loop Window',
    wiazany: false,
    uwaga: 'Okno wspólne platformy — moduł Translate go nie buduje i skrótu nie wiąże.',
  },
  {
    klawisze: 'Ctrl/Cmd + Shift + R',
    dziala: 'Otwarcie rejestru przebiegu pętli wykonawczej',
    okno: 'Execution Loop Window',
    wiazany: false,
    uwaga: 'Okno wspólne platformy — moduł Translate go nie buduje i skrótu nie wiąże.',
  },
  {
    klawisze: 'Ctrl/Cmd + K',
    dziala: 'Wyszukiwarka funkcji modułu',
    okno: 'dowolne okno modułu',
    wiazany: true,
    uwaga:
      'Ten sam skrót prowadzi na pasku górnym do pola poleceń. W module pierwszeństwo ma ' +
      'wyszukiwarka funkcji; poza modułem skrót działa jak dotąd.',
  },
];

/** Czynności wyzwalane skrótami klawiszowymi, które okna modułu dostarczają jako obiekt przy montażu, wołane po dopasowaniu zdarzenia klawiatury. */
export interface CzynnosciSkrotow {
  /** `Ctrl/Cmd + Shift + L` — prowadzi do formularza dodania języka. */
  dodajJezyk(): void;
  /** `Alt + ↑` — panel wymagający uwagi wcześniejszy w wykazie. */
  poprzedniDoUwagi(): void;
  /** `Alt + ↓` — panel wymagający uwagi następny w wykazie. */
  nastepnyDoUwagi(): void;
  /** `Ctrl/Cmd + G`. */
  otworzGlosariusz(): void;
  /** `Ctrl/Cmd + M`. */
  otworzPamiec(): void;
  /** `Ctrl/Cmd + Shift + Q`. */
  otworzJakosc(): void;
  /** `Ctrl/Cmd + K`. */
  otworzWyszukiwarke(): void;
}

/** Podpina skróty klawiszowe do elementu modułu i zwraca funkcję odpięcia, wołaną przy rozłączeniu modułu od dokumentu. */
export function podepnijSkroty(element: HTMLElement, czynnosci: CzynnosciSkrotow): () => void {
  function przyKlawiszu(zdarzenie: KeyboardEvent): void {
    const czynnosc = dopasuj(zdarzenie, czynnosci);
    if (czynnosc === null) return;
    // Zdarzenie zatrzymuje się na module, bo skrót obsłużono tu, inaczej wykonałby się w powłoce.
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
  const klawisz = zdarzenie.key.toLowerCase();

  if (zdarzenie.altKey && !polecenie && klawisz === 'arrowup') return czynnosci.poprzedniDoUwagi;
  if (zdarzenie.altKey && !polecenie && klawisz === 'arrowdown') return czynnosci.nastepnyDoUwagi;
  if (!polecenie) return null;

  if (zdarzenie.shiftKey) {
    if (klawisz === 'l') return czynnosci.dodajJezyk;
    if (klawisz === 'q') return czynnosci.otworzJakosc;
    return null;
  }
  if (klawisz === 'g') return czynnosci.otworzGlosariusz;
  if (klawisz === 'm') return czynnosci.otworzPamiec;
  if (klawisz === 'k') return czynnosci.otworzWyszukiwarke;
  return null;
}

/** Wykaz skrótów jako element warstwy czwartej stanowi drogę bez klawiatury do tej samej informacji, którą niosą skróty. */
export function utworzWykazSkrotow(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 4,
    nazwa: 'Skróty klawiszowe modułu',
    wyjasnienie:
      'Wykaz z załącznika opracowania wraz z tym, które skróty moduł wiąże u siebie.',
    znacznik: '☰',
  });

  const wykaz = document.createElement('ul');
  wykaz.className = 'mt-skroty';
  wykaz.replaceChildren(...SKROTY.map(wierszSkrotu));
  rozwiniecie.tresc.append(wykaz);
  return rozwiniecie.element;
}

function wierszSkrotu(skrot: Skrot): HTMLElement {
  const klawisze = document.createElement('span');
  klawisze.className = 'mt-skroty__klawisze';
  klawisze.textContent = skrot.klawisze;

  const dziala = document.createElement('span');
  dziala.className = 'mt-skroty__dziala';
  dziala.textContent = `${skrot.dziala} — ${skrot.okno}`;

  const element = document.createElement('li');
  element.className = 'mt-skroty__wiersz';
  element.dataset['wiazany'] = skrot.wiazany ? 'tak' : 'nie';
  element.append(klawisze, dziala);

  if (skrot.uwaga !== undefined) {
    const uwaga = document.createElement('p');
    uwaga.className = 'mt-skroty__uwaga';
    uwaga.textContent = skrot.uwaga;
    element.append(uwaga);
  }
  return element;
}
