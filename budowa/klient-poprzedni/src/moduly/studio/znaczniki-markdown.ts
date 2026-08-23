/**
 * Narzędzia znacznikowe edytora Studio — przekształcenia tekstu, nie wywołania rdzenia.
 *
 * Narzędzia nie pytają rdzenia, bo markdown jest tekstem: nagłówek to `#` na
 * początku wiersza, pogrubienie to dwie gwiazdki wokół fragmentu.
 * Przekształcenie zapisuje się w buforze edytora w całości, a treść bufora
 * idzie do rdzenia dopiero komendą `studio.document.save`, tak samo jak każdy
 * znak wpisany z klawiatury.
 *
 * PDF i DOCX to inna rzecz: tam przekształcenie jest konwersją do formatu
 * binarnego, klient jej nie wykona i kontraktu na nią nie ma, więc przycisk
 * odmawia z powodem (`pasek-narzedzi-tekstu.ts`).
 *
 * Plik nie zna DOM. Wejściem jest treść i zakres, wyjściem nowa treść i nowy
 * zakres — dzięki temu przekształcenia sprawdza się bez stawiania okna, a pasek
 * narzędzi nie zna reguł składni.
 */

/**
 * Sposób, w jaki narzędzie dotyka tekstu.
 *
 * - `otoczenie` — obejmuje zaznaczenie parą znaczników (pogrubienie, kod);
 * - `wiersz`    — stawia przedrostek na początku każdego zaznaczonego wiersza
 *                 (nagłówek, lista, cytat);
 * - `numeracja` — jak `wiersz`, ale przedrostek rośnie z numerem wiersza;
 * - `blok`      — dokłada gotowy blok w osobnych wierszach (tabela, linia).
 */
export type RodzajZnacznika = 'otoczenie' | 'wiersz' | 'numeracja' | 'blok';

/** Jedno narzędzie znacznikowe paska edytora. */
export interface NarzedzieTekstu {
  /** Kod pozycji; trafia do `data-narzedzie` przycisku. */
  kod: string;
  /** Etykieta przycisku. */
  nazwa: string;
  /** Zdanie objaśnienia — czym narzędzie jest i co robi z zaznaczeniem. */
  opis: string;
  rodzaj: RodzajZnacznika;
  /** Znacznik otwierający albo przedrostek wiersza. */
  przed: string;
  /** Znacznik zamykający; puste dla `wiersz`, `numeracja` i `blok`. */
  po: string;
  /** Treść wstawiana, gdy zaznaczenia nie ma; dla `blok` — sam blok. */
  zastepnik: string;
}

/** Treść wraz z zakresem zaznaczenia — wejście i wyjście przekształcenia. */
export interface ZakresTekstu {
  tresc: string;
  poczatek: number;
  koniec: number;
}

function narzedzieTekstu(
  kod: string,
  nazwa: string,
  opis: string,
  rodzaj: RodzajZnacznika,
  przed: string,
  po: string,
  zastepnik: string,
): NarzedzieTekstu {
  return { kod, nazwa, opis, rodzaj, przed, po, zastepnik };
}

/**
 * Zestaw narzędzi znacznikowych — szesnaście pozycji w czterech grupach.
 *
 * Stoją tu wyłącznie przekształcenia będące składnią tekstu; formaty binarne
 * wiersza Studio Editor (PDF, DOCX) wymagają konwersji i klient ich nie wykona.
 */
export const NARZEDZIA_TEKSTU: readonly NarzedzieTekstu[] = [
  narzedzieTekstu('naglowek-1', 'H1', 'Nagłówek pierwszego stopnia — „# " na początku wiersza.', 'wiersz', '# ', '', 'Nagłówek'),
  narzedzieTekstu('naglowek-2', 'H2', 'Nagłówek drugiego stopnia — „## " na początku wiersza.', 'wiersz', '## ', '', 'Nagłówek'),
  narzedzieTekstu('naglowek-3', 'H3', 'Nagłówek trzeciego stopnia — „### " na początku wiersza.', 'wiersz', '### ', '', 'Nagłówek'),
  narzedzieTekstu('pogrubienie', 'B', 'Pogrubienie — zaznaczenie w podwójnych gwiazdkach.', 'otoczenie', '**', '**', 'pogrubienie'),
  narzedzieTekstu('kursywa', 'I', 'Kursywa — zaznaczenie w pojedynczych gwiazdkach.', 'otoczenie', '*', '*', 'kursywa'),
  narzedzieTekstu('przekreslenie', 'S', 'Przekreślenie — zaznaczenie w podwójnych tyldach.', 'otoczenie', '~~', '~~', 'przekreślenie'),
  narzedzieTekstu('kod-wiersz', 'Kod', 'Kod w wierszu — zaznaczenie w grawisach.', 'otoczenie', '`', '`', 'kod'),
  narzedzieTekstu('kod-blok', 'Blok kodu', 'Blok kodu — zaznaczenie między liniami grawisów.', 'otoczenie', '```\n', '\n```', 'kod'),
  narzedzieTekstu('lista', 'Lista', 'Lista punktowana — „- " przed każdym zaznaczonym wierszem.', 'wiersz', '- ', '', 'pozycja'),
  narzedzieTekstu('lista-numerowana', 'Lista 1.', 'Lista numerowana — numer rosnący przed każdym wierszem.', 'numeracja', '. ', '', 'pozycja'),
  narzedzieTekstu('zadania', 'Zadania', 'Lista zadań — „- [ ] " przed każdym zaznaczonym wierszem.', 'wiersz', '- [ ] ', '', 'zadanie'),
  narzedzieTekstu('cytat', 'Cytat', 'Cytat blokowy — „> " przed każdym zaznaczonym wierszem.', 'wiersz', '> ', '', 'cytat'),
  narzedzieTekstu('odnosnik', 'Odnośnik', 'Odnośnik — zaznaczenie staje się treścią, adres dopisz w nawiasie.', 'otoczenie', '[', '](adres)', 'treść odnośnika'),
  narzedzieTekstu('obraz', 'Obraz', 'Obraz — zaznaczenie staje się opisem, adres dopisz w nawiasie.', 'otoczenie', '![', '](adres)', 'opis obrazu'),
  narzedzieTekstu('tabela', 'Tabela', 'Zrąb tabeli o dwóch kolumnach, dokładany w osobnych wierszach.', 'blok', '', '', '| Kolumna A | Kolumna B |\n| --- | --- |\n|  |  |'),
  narzedzieTekstu('linia', 'Linia', 'Linia pozioma — trzy myślniki w osobnym wierszu.', 'blok', '', '', '---'),
];

/** Granice wierszy obejmujących zaznaczenie — narzędzia wierszowe pracują na całych wierszach. */
function graniceWierszy(zakres: ZakresTekstu): { od: number; do: number } {
  const { tresc, poczatek } = zakres;
  // Zaznaczenie kończące się na złamaniu wiersza nie obejmuje wiersza
  // następnego: bez tego cofnięcia przedrostek trafiałby do wiersza, którego
  // Operator nie zaznaczył.
  const koniec =
    zakres.koniec > zakres.poczatek && tresc[zakres.koniec - 1] === '\n'
      ? zakres.koniec - 1
      : zakres.koniec;
  const od = tresc.lastIndexOf('\n', poczatek - 1) + 1;
  const nastepne = tresc.indexOf('\n', koniec);
  return { od, do: nastepne === -1 ? tresc.length : nastepne };
}

/** Przedrostek wiersza o wskazanym numerze; numeracja liczy od 1. */
function przedrostek(narzedzie: NarzedzieTekstu, numer: number): string {
  return narzedzie.rodzaj === 'numeracja' ? `${numer + 1}${narzedzie.przed}` : narzedzie.przed;
}

/**
 * Narzędzie wierszowe przełącza, a nie dokłada w nieskończoność.
 *
 * Gdy wszystkie zaznaczone wiersze przedrostek już mają, kolejne naciśnięcie go
 * zdejmuje — inaczej przycisk nie ma stanu odwrotnego, a dwa naciśnięcia
 * „Cytat" zostawiłyby `> > `.
 */
function przestawWiersze(zakres: ZakresTekstu, narzedzie: NarzedzieTekstu): ZakresTekstu {
  const { od, do: az } = graniceWierszy(zakres);
  const wiersze = zakres.tresc.slice(od, az).split('\n');
  const majaJuz = wiersze.every((w, i) => w.startsWith(przedrostek(narzedzie, i)));
  const nowe = wiersze
    .map((w, i) => {
      const znak = przedrostek(narzedzie, i);
      if (majaJuz) return w.slice(znak.length);
      return `${znak}${w === '' && wiersze.length === 1 ? narzedzie.zastepnik : w}`;
    })
    .join('\n');
  return {
    tresc: `${zakres.tresc.slice(0, od)}${nowe}${zakres.tresc.slice(az)}`,
    poczatek: od,
    koniec: od + nowe.length,
  };
}

/** Narzędzie obejmujące: obejmuje zaznaczenie znacznikami albo je zdejmuje. */
function przestawOtoczenie(zakres: ZakresTekstu, narzedzie: NarzedzieTekstu): ZakresTekstu {
  const { tresc, poczatek, koniec } = zakres;
  const wybrany = tresc.slice(poczatek, koniec);
  const { przed, po } = narzedzie;
  const objety =
    wybrany.length >= przed.length + po.length &&
    wybrany.startsWith(przed) &&
    wybrany.endsWith(po);
  if (objety) {
    const rdzen = wybrany.slice(przed.length, wybrany.length - po.length);
    return {
      tresc: `${tresc.slice(0, poczatek)}${rdzen}${tresc.slice(koniec)}`,
      poczatek,
      koniec: poczatek + rdzen.length,
    };
  }
  const rdzen = wybrany === '' ? narzedzie.zastepnik : wybrany;
  const zlozony = `${przed}${rdzen}${po}`;
  return {
    tresc: `${tresc.slice(0, poczatek)}${zlozony}${tresc.slice(koniec)}`,
    // Ognisko wraca na treść, nie na znaczniki: po wstawieniu zastępnika
    // Operator pisze od razu w miejsce, które ma nadpisać.
    poczatek: poczatek + przed.length,
    koniec: poczatek + przed.length + rdzen.length,
  };
}

/** Narzędzie blokowe: dokłada blok w osobnych wierszach za zaznaczeniem. */
function dolozBlok(zakres: ZakresTekstu, narzedzie: NarzedzieTekstu): ZakresTekstu {
  const { tresc } = zakres;
  const { do: az } = graniceWierszy(zakres);
  const wstawka = `\n${narzedzie.zastepnik}`;
  return {
    tresc: `${tresc.slice(0, az)}${wstawka}${tresc.slice(az)}`,
    poczatek: az + 1,
    koniec: az + wstawka.length,
  };
}

/** Wykonuje narzędzie na treści i oddaje treść nową wraz z zakresem do zaznaczenia. */
export function zastosujNarzedzie(
  zakres: ZakresTekstu,
  narzedzie: NarzedzieTekstu,
): ZakresTekstu {
  if (narzedzie.rodzaj === 'otoczenie') return przestawOtoczenie(zakres, narzedzie);
  if (narzedzie.rodzaj === 'blok') return dolozBlok(zakres, narzedzie);
  return przestawWiersze(zakres, narzedzie);
}
