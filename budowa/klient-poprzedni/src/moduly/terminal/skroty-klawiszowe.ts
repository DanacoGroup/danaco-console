/**
 * Skróty klawiszowe modułu Terminal — załącznik dokumentacji modułu przełożony
 * na wiązania klawiatury i na wykaz widoczny dla Operatora.
 *
 * Wykaz obejmuje wszystkie skróty załącznika, także te należące do okna rozmowy
 * i okna pętli wykonawczej, których to złożenie nie osadza — bo są to okna
 * wspólne platformie, składane przez scenę sesji. Pozycja takiego skrótu stoi
 * w wykazie z jawnym zdaniem o tym, gdzie skrót działa; usunięcie jej byłoby
 * ukryciem połowy załącznika, a wiązanie go tutaj przechwytywałoby klawisze
 * cudzemu oknu.
 *
 * Nasłuch stoi na węźle modułu, nie na dokumencie: skrót ma działać, gdy praca
 * dzieje się w Terminalu, a nie odbierać klawisze reszcie platformy.
 *
 * Część skrótów załącznika należy do przeglądarki i przeglądarka oddaje je
 * dopiero wtedy, gdy pozwala na to jej własna polityka. Wykaz mówi to wprost
 * przy każdej takiej pozycji zamiast obiecywać zachowanie, którego moduł nie
 * ma jak wymusić — każda z tych czynności ma drugą drogę w palecie poleceń.
 */

/** Czynność, do której skrót prowadzi. Nazwa słowna, nie kod. */
export type CzynnoscSkrotu =
  | 'nowa-karta'
  | 'zamknij-karte'
  | 'nastepna-karta'
  | 'podziel-widok-karty'
  | 'paleta-polecen'
  | 'szukaj-w-karcie'
  | 'szukaj-w-konsoli'
  | 'przestaw-kolejke';

/** Jedno wiązanie klawiatury. */
export interface WiazanieSkrotu {
  /** Wartość `KeyboardEvent.key` po sprowadzeniu do małych liter. */
  klawisz: string;
  /** Czy skrót wymaga klawisza modyfikującego wspólnego dla obu systemów. */
  sterujacy: boolean;
  shift: boolean;
  /** Zapis skrótu widoczny dla Operatora. */
  zapis: string;
  czynnosc: CzynnoscSkrotu;
  /** Okno, którego skrót dotyczy — nazwa z nagłówka ramy okna. */
  okno: string;
  /** Zdanie o działaniu i o tym, co może je przechwycić. */
  opis: string;
}

/** Skróty, które to złożenie wiąże i wykonuje. */
export const WIAZANIA: readonly WiazanieSkrotu[] = [
  {
    klawisz: 't',
    sterujacy: true,
    shift: false,
    zapis: 'Ctrl/Cmd + T',
    czynnosc: 'nowa-karta',
    okno: 'Terminal Tabs',
    opis: 'Otwiera nową kartę powłoki wybranego rodzaju. Przeglądarka używa tego skrótu do własnej karty i bywa, że go nie oddaje — ta sama czynność stoi w palecie poleceń.',
  },
  {
    klawisz: 'w',
    sterujacy: true,
    shift: false,
    zapis: 'Ctrl/Cmd + W',
    czynnosc: 'zamknij-karte',
    okno: 'Terminal Tabs',
    opis: 'Zamyka kartę bieżącą w widoku; procesy już uruchomione biegną dalej. Przeglądarka zwykle zatrzymuje ten skrót dla siebie — czynność stoi także w panelu akcji okna.',
  },
  {
    klawisz: 'tab',
    sterujacy: true,
    shift: false,
    zapis: 'Ctrl/Cmd + Tab',
    czynnosc: 'nastepna-karta',
    okno: 'Terminal Tabs',
    opis: 'Przestawia ognisko na kolejną kartę paska. Przeglądarka używa tego skrótu do przełączania własnych kart.',
  },
  {
    klawisz: '\\',
    sterujacy: true,
    shift: true,
    zapis: 'Ctrl/Cmd + Shift + \\',
    czynnosc: 'podziel-widok-karty',
    okno: 'Terminal Tabs',
    opis: 'Ustawia treść karty w dwóch kolumnach. To jest podział widoku, nie druga powłoka — kontrakt nie zna panelu wewnątrz karty.',
  },
  {
    klawisz: 'k',
    sterujacy: true,
    shift: false,
    zapis: 'Ctrl/Cmd + K',
    czynnosc: 'paleta-polecen',
    okno: 'Moduł Terminal',
    opis: 'Otwiera paletę poleceń: czynności wszystkich okien złożenia i katalog funkcji modułu.',
  },
  {
    klawisz: 'f',
    sterujacy: true,
    shift: false,
    zapis: 'Ctrl/Cmd + F',
    czynnosc: 'szukaj-w-karcie',
    okno: 'Output Console',
    opis: 'Stawia ognisko w polu wyszukiwania i zawęża strumień do karty bieżącej — odpowiednik wyszukiwania w obrębie karty.',
  },
  {
    klawisz: 'f',
    sterujacy: true,
    shift: true,
    zapis: 'Ctrl/Cmd + Shift + F',
    czynnosc: 'szukaj-w-konsoli',
    okno: 'Output Console',
    opis: 'Stawia ognisko w polu wyszukiwania nad strumieniem zbiorczym wszystkich kart.',
  },
  {
    klawisz: '.',
    sterujacy: true,
    shift: false,
    zapis: 'Ctrl/Cmd + .',
    czynnosc: 'przestaw-kolejke',
    okno: 'Task & Schedule',
    opis: 'Wstrzymuje albo wznawia pierwszą odczytaną kolejkę silnika pętli obsługującą to okno.',
  },
];

/** Pozycja załącznika, której to złożenie nie wiąże, wraz z powodem. */
export interface SkrotPozaZlozeniem {
  zapis: string;
  okno: string;
  opis: string;
}

/** Skróty załącznika należące do okien wspólnych platformie. */
export const SKROTY_POZA_ZLOZENIEM: readonly SkrotPozaZlozeniem[] = [
  {
    zapis: 'Ctrl/Cmd + Enter',
    okno: 'Chat Window',
    opis: 'Wysłanie polecenia do modelu. Okno rozmowy jest bytem sesji wspólnym wszystkim modułom i nie wchodzi do tego złożenia.',
  },
  {
    zapis: 'Ctrl/Cmd + Shift + L',
    okno: 'Execution Loop Window',
    opis: 'Otwarcie kolumny pętli wykonawczej. Okno pętli składa scena sesji, nie moduł.',
  },
  {
    zapis: 'Ctrl/Cmd + .',
    okno: 'Execution Loop Window',
    opis: 'Wstrzymanie i wznowienie przebiegu pętli. W tym złożeniu ten sam skrót przestawia kolejkę silnika pętli obsługującą okno modułu.',
  },
  {
    zapis: '/',
    okno: 'Chat Window',
    opis: 'Wywołanie menu operacji z pola polecenia okna rozmowy.',
  },
  {
    zapis: 'Strzałka w górę / w dół w polu polecenia',
    okno: 'Terminal Tabs',
    opis: 'Wędrówka po historii poleceń karty. Kontrakt nie oddaje historii poleceń; moduł pamięta wyłącznie ostatnie polecenie karty, po które sięga ponowienie.',
  },
];

/** Nasłuch skrótów założony na węźle modułu. */
export interface Skroty {
  /** Zdejmuje nasłuch. Wołane przy zamykaniu modułu. */
  zamknij(): void;
}

/**
 * Wiąże skróty z czynnościami modułu.
 *
 * @param element węzeł modułu — skrót działa, gdy ognisko jest wewnątrz niego.
 * @param czynnosci wykonania po nazwie czynności; nazwa bez wykonania jest
 *   pomijana, bo skrót bez skutku byłby obietnicą bez pokrycia.
 */
export function zwiazSkroty(
  element: HTMLElement,
  czynnosci: Readonly<Partial<Record<CzynnoscSkrotu, () => void>>>,
): Skroty {
  function naKlawisz(zdarzenie: KeyboardEvent): void {
    // Klawisz modyfikujący jest inny na obu systemach i oba są tu równoprawne:
    // zapis „Ctrl/Cmd" załącznika znaczy dokładnie tyle.
    const sterujacy = zdarzenie.ctrlKey || zdarzenie.metaKey;
    const klawisz = zdarzenie.key.toLowerCase();
    for (const wiazanie of WIAZANIA) {
      if (wiazanie.klawisz !== klawisz) continue;
      if (wiazanie.sterujacy !== sterujacy) continue;
      if (wiazanie.shift !== zdarzenie.shiftKey) continue;
      const wykonaj = czynnosci[wiazanie.czynnosc];
      if (wykonaj === undefined) continue;
      zdarzenie.preventDefault();
      wykonaj();
      return;
    }
  }

  element.addEventListener('keydown', naKlawisz);
  return { zamknij: () => element.removeEventListener('keydown', naKlawisz) };
}

/**
 * Wykaz skrótów widoczny dla Operatora — cały załącznik, z rozróżnieniem
 * pozycji wiązanych w tym złożeniu i należących do okien wspólnych.
 */
export function wykazSkrotow(): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-skroty';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-skroty__podpis';
  podpis.textContent = 'Skróty klawiszowe modułu';
  blok.append(podpis);

  const lista = document.createElement('dl');
  lista.className = 'dt-skroty__wykaz';
  for (const wiazanie of WIAZANIA) {
    const zapis = document.createElement('dt');
    zapis.textContent = wiazanie.zapis;
    const opis = document.createElement('dd');
    opis.textContent = `${wiazanie.okno} — ${wiazanie.opis}`;
    lista.append(zapis, opis);
  }
  for (const pozycja of SKROTY_POZA_ZLOZENIEM) {
    const zapis = document.createElement('dt');
    zapis.dataset['poza'] = 'zlozenie';
    zapis.textContent = pozycja.zapis;
    const opis = document.createElement('dd');
    opis.textContent = `${pozycja.okno} — ${pozycja.opis}`;
    lista.append(zapis, opis);
  }
  blok.append(lista);
  return blok;
}
