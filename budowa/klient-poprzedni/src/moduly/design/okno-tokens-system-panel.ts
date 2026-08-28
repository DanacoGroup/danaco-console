import './tokeny.css';

import { przycisk, poleWyboru, pobierzPlik, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { utworzZakladkiSekcji } from '../../modele/zakladki-sekcji';
import { ZDARZENIE_MOTYWU, motywObowiazujacy, ustawMotyw, type Motyw } from '../../motyw/motyw';
import {
  POSTACIE_WYDANIA,
  zlozPrzewodnikStylu,
  zlozWydanie,
  type PostacWydania,
} from './eksport-tokenow';
import { OKNO_TOKENS_SYSTEM_PANEL } from './etykiety-designu';
import {
  liczbaPar,
  odczytajBarwe,
  poniejProgu,
  zapisWspolczynnika,
  zmierzPary,
  type PomiarPary,
} from './kontrast-wcag';
import { naglowekOkna } from './stan-okna';
import {
  RODZAJE_WIDZENIA,
  symuluj,
  zapisBarwy,
  type RodzajWidzenia,
} from './symulacja-widzenia';
import { DesignTokenKind, type DesignToken } from '../../../../shared/contract';
import { utworzRozwiniecie, type Rozwiniecie } from './warstwy-designu';
import { utworzZestawyZetonow, type ZestawyZetonow } from './zestawy-zetonow';
import type { StanDesignu } from './stan-designu';
import {
  liczbaBezDefinicji,
  odczytajZetony,
  selektoryZetonu,
  type GrupaZetonow,
  type Zeton,
} from './zetony-systemu';

/**
 * Tokens & System Panel jest piątym oknem operacyjnym modułu Design: definiuje
 * i wydaje system projektowy — żetony, motywy, reguły dostępności i przewodnik
 * stylu — mierząc wartości bezpośrednio na uruchomionym produkcie, zamiast je
 * opisywać.
 */
export interface OknoTokensSystemPanel {
  element: HTMLElement;
  odswiez(): void;
  /** Rozwija okno i prowadzi do niego ognisko — droga skrótu klawiszowego. */
  otworz(): void;
  /** Zdejmuje nasłuch zmiany motywu. */
  zamknij(): void;
}

export function utworzOknoTokensSystemPanel(stan: StanDesignu): OknoTokensSystemPanel {
  let grupy: readonly GrupaZetonow[] = odczytajZetony();
  let pomiary: readonly PomiarPary[] = zmierzPary();

  // ── Zakładka „Żetony" ──────────────────────────────────────────────────────
  const drzewo = document.createElement('div');
  drzewo.className = 'md-zetony__drzewo';

  const powiazania = document.createElement('div');
  powiazania.className = 'md-zetony__powiazania';

  const obszarZetonow = sekcja([
    zdanie(
      'Drzewo czyta wartości z motywu obowiązującego w tej chwili. Naciśnięcie żetonu ' +
        'pokazuje selektory arkuszy produktu, w których ten żeton występuje — to jest ' +
        'powiązanie żetonu z komponentami, zmierzone, nie spisane.',
    ),
    drzewo,
    powiazania,
  ]);

  // ── Zakładka „Motywy" ──────────────────────────────────────────────────────
  const stanMotywu = document.createElement('p');
  stanMotywu.className = 'md-zetony__motyw';

  const naJasny = przycisk('Podgląd w motywie jasnym', 'dn-btn dn-btn--sm dn-btn--zarys');
  const naCiemny = przycisk('Podgląd w motywie ciemnym', 'dn-btn dn-btn--sm dn-btn--zarys');
  naJasny.addEventListener('click', () => przelaczPodglad('light'));
  naCiemny.addEventListener('click', () => przelaczPodglad('dark'));

  const paskaMotywu = document.createElement('div');
  paskaMotywu.className = 'md-zetony__pasek';
  paskaMotywu.append(naJasny, naCiemny);

  const obszarMotywow = sekcja([
    stanMotywu,
    paskaMotywu,
    zdanie(
      'Oba motywy są równoprawne i mają własne wartości tych samych ról, więc każdy wymaga ' +
        'osobnego pomiaru i osobnego wydania. Przełącznik woła tę samą czynność motywu co ' +
        'przełącznik paska górnego — drugiego stanu motywu tu nie ma.',
    ),
    zdanie(
      'Wariantów marki poza dwoma motywami produktu nie ma. Zestawu nadpisań nie ma też gdzie ' +
        'zapisać: kontrakt nie zna bytu zestawu żetonów, więc trwały system projektowy jest ' +
        'dziś poza zasięgiem modułu.',
    ),
  ]);

  // ── Zakładka „Dostępność" ──────────────────────────────────────────────────
  const podsumowanieKontrastu = document.createElement('p');
  podsumowanieKontrastu.className = 'md-zetony__podsumowanie';

  const tabelaPar = document.createElement('table');
  tabelaPar.className = 'dn-tabela md-zetony__tabela';

  const wyborWidzenia = poleWyboru(
    {
      etykieta: 'Symulacja widzenia barw',
      opis:
        'Symulacja obejmuje próbki żetonów. Obrazów nie obejmuje — bajtów zasobu nie ma czym ' +
        'pobrać, bo pole odsyłacza jest ścieżką w systemie plików rdzenia.',
    },
    [
      { wartosc: '', etykieta: 'widzenie bez wady' },
      ...RODZAJE_WIDZENIA.map(([kod, nazwa]) => ({ wartosc: kod, etykieta: nazwa })),
    ],
  );
  wyborWidzenia.kontrolka.addEventListener('change', () => rysujProbki());

  const probki = document.createElement('div');
  probki.className = 'md-zetony__probki';

  const obszarDostepnosci = sekcja([
    podsumowanieKontrastu,
    tabelaPar,
    zdanie(
      'Wykaz progów nie podaje ROLI pary, a norma dopuszcza próg niższy dla tekstu dużego oraz ' +
        'dla obrysów i wskaźników skupienia. Para obrysu wychodzi więc poniżej progu, choć ' +
        'wobec właściwej reguły może być zgodna — to jest pomiar, nie werdykt.',
    ),
    wyborWidzenia.element,
    probki,
    zdanie(
      'Rachunek symulacji jest przybliżeniem przyjętym w narzędziach projektowych: pokazuje, ' +
        'które barwy systemu zlewają się w jedną, i nie służy orzeczeniu medycznemu.',
    ),
  ]);

  // ── Zakładka „Wydanie" ─────────────────────────────────────────────────────
  const postac = poleWyboru(
    {
      etykieta: 'Postać wydania żetonów',
      opis: 'Plik składa przeglądarka z wartości motywu obowiązującego.',
    },
    POSTACIE_WYDANIA.map(([kod, nazwa]) => ({ wartosc: kod, etykieta: nazwa })),
  );

  const wydaj = przycisk('Wydaj żetony do pliku', 'dn-btn dn-btn--sm dn-btn--atrament');
  const przewodnik = przycisk('Wydaj przewodnik stylu', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedzWydania = utworzWierszOdpowiedzi();

  const wczytaj = document.createElement('input');
  wczytaj.type = 'file';
  wczytaj.className = 'dn-pole-kontrolka md-zetony__wczytanie';
  wczytaj.accept = 'application/json,.json';
  wczytaj.id = 'md-zetony-import';

  const etykietaWczytania = document.createElement('label');
  etykietaWczytania.className = 'dn-pole-etykieta';
  etykietaWczytania.htmlFor = wczytaj.id;
  etykietaWczytania.textContent = 'Wczytaj zestaw żetonów do porównania';

  const roznica = document.createElement('pre');
  roznica.className = 'md-zetony__roznica';
  roznica.hidden = true;

  // Panel zestawów trwałych czyta żetony motywu w chwili naciśnięcia, nie z kopii sprzed otwarcia okna.
  const zestawy: ZestawyZetonow = utworzZestawyZetonow(stan, {
    zetonyMotywu: () => zetonyKontraktu(odczytajZetony()),
  });

  const paskaWydania = document.createElement('div');
  paskaWydania.className = 'md-zetony__pasek';
  paskaWydania.append(wydaj, przewodnik);

  const obszarWydania = sekcja([
    postac.element,
    paskaWydania,
    odpowiedzWydania.element,
    etykietaWczytania,
    wczytaj,
    roznica,
    zdanie(
      'Wczytany zestaw jest ZESTAWIANY z żetonami motywu, nie nadpisuje ich: motyw należy do ' +
        'powłoki, a moduł nie ma prawa przestawiać wyglądu całego produktu. Zapis trwały ' +
        'zestawu, wydanie go do kodu i wydanie przewodnika do modułu docelowego prowadzi ' +
        'panel poniżej — komendami rodziny design.tokenset.* i design.styleguide.publish.',
    ),
    zestawy.element,
  ]);

  wydaj.addEventListener('click', () => wydajZetony());
  przewodnik.addEventListener('click', () => wydajPrzewodnik());
  wczytaj.addEventListener('change', () => void porownajWczytany());

  // ── Złożenie okna ──────────────────────────────────────────────────────────
  const zakladki = utworzZakladkiSekcji([
    { kod: 'zetony', nazwa: 'Żetony', element: obszarZetonow },
    { kod: 'motywy', nazwa: 'Motywy', element: obszarMotywow },
    { kod: 'dostepnosc', nazwa: 'Dostępność', element: obszarDostepnosci },
    { kod: 'wydanie', nazwa: 'Wydanie', element: obszarWydania },
  ]);

  const rozwiniecie: Rozwiniecie = utworzRozwiniecie({
    warstwa: 3,
    nazwa: OKNO_TOKENS_SYSTEM_PANEL.nazwa,
    wyjasnienie:
      'Żetony, motywy, kontrast i wydanie systemu projektowego. Wszystkie wartości są mierzone ' +
      'na uruchomionym produkcie, żadna nie jest przepisana do modułu.',
    znacznik: '⋮',
  });
  rozwiniecie.tresc.append(
    naglowekOkna(OKNO_TOKENS_SYSTEM_PANEL.nazwa, OKNO_TOKENS_SYSTEM_PANEL.rola),
    zakladki.element,
  );

  const element = document.createElement('section');
  element.className = 'md-okno md-okno--zetony';
  element.dataset['okno'] = OKNO_TOKENS_SYSTEM_PANEL.kod;
  element.append(rozwiniecie.element);

  // Przełączenie motywu podglądu woła czynność motywu produktu, nie własną — jeden nośnik stanu.
  function przelaczPodglad(wybor: Motyw): void {
    ustawMotyw(wybor);
    odswiez();
  }

  function wydajZetony(): void {
    const wybrana = postac.kontrolka.value as PostacWydania;
    const opis = POSTACIE_WYDANIA.find(([kod]) => kod === wybrana);
    if (opis === undefined) {
      odpowiedzWydania.pokaz('Nie wskazano postaci wydania.', false);
      return;
    }
    const [, nazwa, plik] = opis;
    const tresc = zlozWydanie(grupy, wybrana, motywObowiazujacy());
    pobierzPlik(plik, tresc, 'text/plain');
    const wydanych = grupy.reduce(
      (suma, grupa) => suma + grupa.zetony.filter((zeton) => zeton.wartosc !== '').length,
      0,
    );
    odpowiedzWydania.pokaz(
      `Wydano ${String(wydanych)} żetonów w postaci „${nazwa}" do pliku ${plik}, z motywu ` +
        `${motywObowiazujacy()}. Drugi motyw wymaga osobnego wydania.`,
      true,
    );
  }

  function wydajPrzewodnik(): void {
    pobierzPlik(
      'przewodnik-systemu.html',
      zlozPrzewodnikStylu(grupy, motywObowiazujacy()),
      'text/html',
    );
    odpowiedzWydania.pokaz(
      'Przewodnik stylu złożony i oddany plikiem. Wydania go do modułów Library i Studio ' +
        'kontrakt nie zna — plik zostaje u Operatora.',
      true,
    );
  }

  // Zestawienie wczytanego zestawu z żetonami motywu, porównane po nazwach ról, nie po strukturze pliku.
  async function porownajWczytany(): Promise<void> {
    const plik = wczytaj.files?.[0];
    if (plik === undefined) {
      roznica.hidden = true;
      return;
    }
    let wczytane: Record<string, string>;
    try {
      wczytane = splaszcz(JSON.parse(await plik.text()) as unknown);
    } catch (blad) {
      roznica.hidden = false;
      roznica.textContent =
        `Pliku „${plik.name}" nie udało się odczytać jako zapisu JSON: ` +
        `${blad instanceof Error ? blad.message : 'przeglądarka nie podała przyczyny'}.`;
      return;
    }

    const wlasne = new Map(
      grupy.flatMap((grupa) => grupa.zetony.map((zeton) => [zeton.nazwa, zeton.wartosc] as const)),
    );
    const zgodne: string[] = [];
    const inne: string[] = [];
    const nieznane: string[] = [];

    for (const [nazwa, wartosc] of Object.entries(wczytane)) {
      const bez = nazwa.replace(/^--dn-/, '');
      const moja = wlasne.get(bez);
      if (moja === undefined) nieznane.push(bez);
      else if (moja === wartosc) zgodne.push(bez);
      else inne.push(`${bez}: motyw „${moja}", zestaw „${wartosc}"`);
    }

    roznica.hidden = false;
    roznica.textContent = [
      `Zestaw „${plik.name}" wobec motywu ${motywObowiazujacy()}:`,
      `zgodnych ról: ${String(zgodne.length)}`,
      `ról o innej wartości: ${String(inne.length)}`,
      ...inne.map((wiersz) => `  ${wiersz}`),
      `ról spoza systemu produktu: ${String(nieznane.length)}`,
      ...(nieznane.length === 0 ? [] : [`  ${nieznane.join(', ')}`]),
      '',
      'Zestawienie niczego nie zmieniło — motyw pozostał nietknięty.',
    ].join('\n');
  }

  function rysujDrzewo(): void {
    drzewo.replaceChildren(...grupy.map(sekcjaGrupy));
  }

  function sekcjaGrupy(grupa: GrupaZetonow): HTMLElement {
    const tytul = document.createElement('h5');
    tytul.className = 'md-zetony__grupa';
    tytul.textContent = grupa.nazwa;

    const wykaz = document.createElement('ul');
    wykaz.className = 'md-zetony__wykaz';
    wykaz.replaceChildren(...grupa.zetony.map(wierszZetonu));

    const element = document.createElement('section');
    element.className = 'md-zetony__sekcja';
    element.append(tytul, wykaz);
    return element;
  }

  function wierszZetonu(zeton: Zeton): HTMLElement {
    const probka = document.createElement('span');
    probka.className = 'md-zetony__probka';
    probka.dataset['rodzaj'] = zeton.rodzaj;
    if (zeton.rodzaj === 'barwa' && zeton.wartosc !== '') {
      probka.style.background = zeton.wartosc;
    }

    const nazwa = document.createElement('span');
    nazwa.className = 'md-zetony__nazwa';
    nazwa.textContent = zeton.nazwa;

    const wartosc = document.createElement('span');
    wartosc.className = 'md-zetony__wartosc';
    wartosc.dataset['podana'] = String(zeton.wartosc !== '');
    wartosc.textContent =
      zeton.wartosc === '' ? 'motyw nie definiuje tego żetonu' : zeton.wartosc;

    const kontrolka = document.createElement('button');
    kontrolka.type = 'button';
    kontrolka.className = 'md-zetony__wiersz';
    kontrolka.dataset['zeton'] = zeton.nazwa;
    kontrolka.append(probka, nazwa, wartosc);
    kontrolka.addEventListener('click', () => pokazPowiazania(zeton));

    const element = document.createElement('li');
    element.append(kontrolka);
    return element;
  }

  /** Selektory, w których żeton występuje — powiązanie z komponentami. */
  function pokazPowiazania(zeton: Zeton): void {
    const selektory = selektoryZetonu(zeton.nazwa);
    const tytul = document.createElement('h5');
    tytul.className = 'md-zetony__grupa';
    tytul.textContent = `Powiązania żetonu ${zeton.nazwa}`;

    const tresc = document.createElement('p');
    tresc.className = 'md-zetony__powiazania-tresc';
    tresc.textContent =
      selektory.length === 0
        ? 'Żaden wczytany arkusz nie sięga po ten żeton. To znaczy albo że żeton nie ma dziś ' +
          'użycia, albo że arkusz, który go używa, nie został jeszcze wczytany do dokumentu.'
        : `${String(selektory.length)} selektorów: ${selektory.join(', ')}`;

    // Zapis barwy w składowych — konwersja przestrzeni barw dla żetonu barwnego.
    powiazania.replaceChildren(tytul, tresc, ...opisBarwy(zeton));
  }

  function opisBarwy(zeton: Zeton): readonly HTMLElement[] {
    if (zeton.rodzaj !== 'barwa') return [];
    const barwa = odczytajBarwe(zeton.wartosc);
    if (barwa === null) return [];
    const element = document.createElement('p');
    element.className = 'md-zetony__powiazania-tresc';
    element.textContent =
      `Składowe: czerwona ${String(Math.round(barwa.r))}, zielona ${String(Math.round(barwa.g))}, ` +
      `niebieska ${String(Math.round(barwa.b))}.`;
    return [element];
  }

  function rysujKontrast(): void {
    const zlamane = poniejProgu(pomiary);
    podsumowanieKontrastu.textContent =
      `Zmierzono ${String(pomiary.length)} par z wykazu progów produktu (wykaz niesie ` +
      `${String(liczbaPar())}). Poniżej progu albo bez pomiaru: ${String(zlamane.length)}. ` +
      `Motyw mierzony: ${motywObowiazujacy()}.`;

    const naglowek = document.createElement('tr');
    for (const tytul of ['Para', 'Próg', 'Pomiar', 'Ocena']) {
      const komorka = document.createElement('th');
      komorka.textContent = tytul;
      naglowek.append(komorka);
    }
    tabelaPar.replaceChildren(naglowek, ...pomiary.map(wierszPomiaru));
  }

  function wierszPomiaru(pomiar: PomiarPary): HTMLElement {
    const wiersz = document.createElement('tr');
    wiersz.dataset['spelnia'] = String(pomiar.spelnia);

    for (const tresc of [
      pomiar.para,
      pomiar.prog.toFixed(2),
      zapisWspolczynnika(pomiar.wspolczynnik),
      // Stan nigdy samym kolorem: ocena stoi słowem, nie odcieniem wiersza.
      pomiar.powod !== '' ? pomiar.powod : pomiar.spelnia ? 'sięga progu' : 'poniżej progu',
    ]) {
      const komorka = document.createElement('td');
      komorka.textContent = tresc;
      wiersz.append(komorka);
    }
    return wiersz;
  }

  function rysujProbki(): void {
    const wybrany = wyborWidzenia.kontrolka.value as RodzajWidzenia | '';
    const barwne = grupy
      .flatMap((grupa) => grupa.zetony)
      .filter((zeton) => zeton.rodzaj === 'barwa' && zeton.wartosc !== '');

    probki.replaceChildren(
      ...barwne.map((zeton) => {
        const barwa = odczytajBarwe(zeton.wartosc);
        const element = document.createElement('span');
        element.className = 'md-zetony__kafel';
        element.textContent = zeton.nazwa;
        if (barwa !== null) {
          element.style.background =
            wybrany === '' ? zapisBarwy(barwa) : zapisBarwy(symuluj(barwa, wybrany));
        }
        return element;
      }),
    );
  }

  function odswiez(): void {
    grupy = odczytajZetony();
    pomiary = zmierzPary();
    const bez = liczbaBezDefinicji(grupy);
    stanMotywu.textContent =
      `Motyw obowiązujący: ${motywObowiazujacy()}. Ról w wykazie panelu: ` +
      `${String(grupy.reduce((suma, grupa) => suma + grupa.zetony.length, 0))}, ` +
      `z czego bez definicji w motywie: ${String(bez)}.`;
    rysujDrzewo();
    rysujKontrast();
    rysujProbki();
  }

  // Zmiana motywu — także wywołana z paska górnego — przelicza cały panel.
  const naZmianeMotywu = (): void => odswiez();
  window.addEventListener(ZDARZENIE_MOTYWU, naZmianeMotywu);

  odswiez();

  return {
    element,
    odswiez,
    otworz() {
      rozwiniecie.rozwin();
      zakladki.pokaz('zetony');
    },
    zamknij() {
      window.removeEventListener(ZDARZENIE_MOTYWU, naZmianeMotywu);
    },
  };
}

/**
 * Spłaszczenie wczytanego zestawu żetonów do par nazwa roli i wartość,
 * obsługujące zarówno zapis płaski, jak i zapis wedle wzorca W3C z wartością
 * pod kluczem $value.
 */
function splaszcz(zestaw: unknown, przedrostek = ''): Record<string, string> {
  const wynik: Record<string, string> = {};
  if (typeof zestaw !== 'object' || zestaw === null) return wynik;

  for (const [klucz, wartosc] of Object.entries(zestaw as Record<string, unknown>)) {
    if (klucz === '$value') {
      if (typeof wartosc === 'string' || typeof wartosc === 'number') {
        wynik[przedrostek] = String(wartosc);
      }
      continue;
    }
    const nazwa = przedrostek === '' ? klucz : `${przedrostek}-${klucz}`;
    if (typeof wartosc === 'string' || typeof wartosc === 'number') {
      wynik[nazwa] = String(wartosc);
      continue;
    }
    Object.assign(wynik, splaszcz(wartosc, nazwa));
  }
  return wynik;
}

/** Akapit objaśnienia stosowany dla jednego, spójnego brzmienia oprawy tekstowej w całym tym oknie panelu. */
function zdanie(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis md-zetony__zdanie';
  element.textContent = tresc;
  return element;
}

/** Obszar zakładki panelu Tokens & System — nośnik elementów potomnych, bez własnej wiedzy o niesionej treści. */
function sekcja(dzieci: readonly HTMLElement[]): HTMLElement {
  const element = document.createElement('div');
  element.className = 'md-zetony__obszar';
  element.append(...dzieci);
  return element;
}

/**
 * Przekłada żetony odczytane z motywu na kształt kontraktu DesignToken, bo
 * panel nazywa rodzaje po polsku, a kontrakt — po angielsku wedle W3C Design
 * Tokens; żeton, którego motyw nie definiuje, nie wchodzi do zestawu.
 */
function zetonyKontraktu(wykaz: readonly GrupaZetonow[]): readonly DesignToken[] {
  const zetony: DesignToken[] = [];
  for (const grupa of wykaz) {
    for (const zeton of grupa.zetony) {
      if (zeton.wartosc === '') continue;
      zetony.push({
        name: zeton.nazwa,
        kind: rodzajKontraktu(zeton.rodzaj),
        value: zeton.wartosc,
        description: grupa.nazwa,
      });
    }
  }
  return zetony;
}

/** Przekład rodzaju żetonu panelu żetonów systemu na odpowiadający mu rodzaj kontraktu DesignTokenKind. */
function rodzajKontraktu(rodzaj: Zeton['rodzaj']): DesignTokenKind {
  switch (rodzaj) {
    case 'barwa':
      return DesignTokenKind.Color;
    case 'krój':
      return DesignTokenKind.FontFamily;
    case 'czas':
      return DesignTokenKind.Duration;
    case 'liczba':
      return DesignTokenKind.FontWeight;
    default:
      // Miara jest domyślnym rodzajem: odstępy, promienie, stopnie pisma i wymiary to grupa najliczniejsza.
      return DesignTokenKind.Dimension;
  }
}
