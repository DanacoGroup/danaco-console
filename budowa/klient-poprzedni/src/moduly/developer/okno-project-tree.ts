import {
  TreeNodeKind,
  type DeveloperTreeGetRequest,
  type DeveloperTreeNode,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pole,
  poleLiczbowe,
  przelacznik,
  przyciskAkcji,
  przyciskBezKomendy,
  wiersz,
} from '../../modele/kontrolki-formularza-braki';
import { powodBezKomendy } from './braki-kontraktu';
import type { StanDevelopera } from './stan-developer';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import type { ZrodloDeveloper } from './zrodlo-developer';
import './okno-project-tree.css';

/**
 * Project Tree — okno pomocnicze modułu Developer: nawigacja po plikach
 * i otwarcie pliku w Code Editor.
 *
 * Otwarcie pliku to wyłącznie `stan.wskazPlik` — Code Editor sam nasłuchuje
 * i sam woła `developer.file.open`; wołanie jej też stąd czytałoby plik dwa
 * razy.
 *
 * `developer.tree.get` oddaje listę płaską — hierarchię składa `zbudujDrzewo`
 * z `parentPath`. Żaden węzeł nie znika po cichu: samowskazanie rodzica, cykl
 * wzajemny, zduplikowana ścieżka i rodzic nieobecny w wykazie trafiają do
 * korzenia z jawną adnotacją (`wykryjCykle`). Gdyby po złożeniu nie zostało
 * nic, `rysujDrzewoProjektu` nazywa to stanem błędu, nie stanem „treść”.
 *
 * Pole „Katalog” niesie żądanie Operatora, `root` z odpowiedzi — miejsce, od
 * którego rdzeń naprawdę czytał. Okno pokazuje jedno i drugie (`opiszKorzen`).
 *
 * Menu kontekstowe pod prawym przyciskiem myszy jest tu jedynym takim
 * przypadkiem w platformie — gdzie indziej „menu kontekstowe” znaczy lewy klik
 * ikony `⋯`. Stąd nasłuch `contextmenu` z `preventDefault`.
 *
 * Jedenaście z trzynastu pozycji menu nie ma pokrycia w kontrakcie — stoją
 * jako `przyciskBezKomendy` z powodem. „Skopiuj ścieżkę” (schowek) i „Odśwież”
 * (odczyt) są zrobione naprawdę.
 */
export interface OknoProjectTree {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka nasłuch `stan.naZmiane` założony przy montażu okna. */
  zamknij(): void;
}

export function utworzOknoProjectTree(
  zrodlo: ZrodloDeveloper,
  stan: StanDevelopera,
): OknoProjectTree {
  const rama = utworzRameOkna({
    tytul: 'Project Tree',
    rola: 'pomocnicze',
    kod: 'project-tree',
    przeznaczenie:
      'Nawigacja po plikach; otwarcie pliku w Code Editor. Odzwierciedla aktualny stan repozytorium.',
    modul: 'Developer',
    przedrostek: 'mdev',
  });
  const tresc = utworzStanTresci();
  const powierzchnia = zlozPowierzchnieDrzewa(rama, tresc.element);
  const menu = utworzMenuKontekstowe();
  const rozwiniete = new Set<string>();
  let wezly: readonly DeveloperTreeNode[] = [];

  function odczytaj(): void {
    tresc.ladowanie('Odczyt drzewa projektu…');
    const zadanie = zadanieDrzewa(stan.okno(), powierzchnia);
    void zrodlo.drzewo(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Rdzeń nie oddał drzewa projektu.', wynik.blad);
        return;
      }
      stan.ustawKorzen(wynik.wynik.root);
      wezly = wynik.wynik.nodes;
      rysuj();
    });
  }

  function rysuj(): void {
    opiszKorzen(powierzchnia.korzen, stan.korzen());
    rysujDrzewoProjektu(wezly, tresc, {
      stan,
      rozwiniete,
      fraza: powierzchnia.szukaj.value.trim(),
      naPrzelacz: rysuj,
      naMenu: (event, wezel) => menu.pokaz(event.clientX, event.clientY, wezel, tresc, odczytaj),
    });
  }

  powierzchnia.odswiezPrzycisk.addEventListener('click', odczytaj);
  // Zawężenie przerysowuje wyłącznie widok — nic nie jedzie do rdzenia.
  powierzchnia.szukaj.addEventListener('input', rysuj);
  powierzchnia.zwinPrzycisk.addEventListener('click', () => {
    rozwiniete.clear();
    rysuj();
  });
  const odsubskrybuj = stan.naZmiane(rysuj);

  // Jeden punkt wejścia dla odczytu — złożenie modułu woła `odswiez()` samo,
  // tak jak `moduly/automations/indeks.ts` woła je raz dla każdego okna po
  // złożeniu. Wywołanie tutaj powielałoby `developer.tree.get` przy montażu.
  return { element: rama.element, odswiez: odczytaj, zamknij: odsubskrybuj };
}

/** Kontrolki Operatora: ścieżka, głębokość, pozycje ukryte oraz przycisk odświeżenia. */
interface PowierzchniaDrzewa {
  sciezka: HTMLInputElement;
  glebokosc: HTMLInputElement;
  ukryte: HTMLInputElement;
  /**
   * Zawężenie po fragmencie nazwy — czynność WYŁĄCZNIE kliencka nad węzłami już
   * odczytanymi. Nie jedzie do rdzenia i nie jest tym samym co pole „Katalog”,
   * które zmienia zakres odczytu.
   */
  szukaj: HTMLInputElement;
  odswiezPrzycisk: HTMLButtonElement;
  zwinPrzycisk: HTMLButtonElement;
  /** Miejsce na korzeń oddany przez rdzeń — nie na to, co wpisał Operator. */
  korzen: HTMLElement;
}

/**
 * Pokazuje korzeń, od którego rdzeń naprawdę czytał drzewo.
 *
 * Osobny wiersz obok pola „Katalog”, bo pole niesie żądanie Operatora (puste
 * znaczy „czytaj od korzenia okna”), a odpowiedź rdzenia niesie `root` — i to
 * są dwie różne rzeczy. Bez tego wiersza pusty katalog roboczy i katalog nie
 * ten, co trzeba, wyglądają identycznie.
 */
function opiszKorzen(miejsce: HTMLElement, korzen: string): void {
  miejsce.dataset['wskazana'] = korzen === '' ? 'nie' : 'tak';
  miejsce.textContent =
    korzen === ''
      ? 'Rdzeń nie podał jeszcze korzenia — odśwież drzewo.'
      : `Rdzeń czytał od: ${korzen}`;
}

/** Składa pasek narzędzi i panel akcji ramy; ciało ramy dostaje stan treści. */
function zlozPowierzchnieDrzewa(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
): PowierzchniaDrzewa {
  const sciezka = pole('Katalog', 'domyślnie korzeń repozytorium');
  const glebokosc = poleLiczbowe('Głębokość drzewa', 'domyślnie bez ograniczenia');
  const ukryte = przelacznik('Uwzględnij pozycje ukryte');
  const szukaj = pole('Szybkie otwarcie po fragmencie nazwy', 'fragment nazwy pliku');
  rama.narzedzia.append(
    wiersz('Katalog', sciezka, { klasa: 'mdev-wiersz', objasnienie: 'Puste pole czyta od korzenia okna.' }),
    wiersz('Głębokość', glebokosc, { klasa: 'mdev-wiersz' }),
    wiersz('Pozycje ukryte', ukryte, { klasa: 'mdev-wiersz' }),
    wiersz('Szybkie otwarcie', szukaj, {
      klasa: 'mdev-wiersz',
      objasnienie:
        'Zawężenie klienckie nad węzłami już odczytanymi — do rdzenia nic nie jedzie. Katalog ' +
        'zostaje na scenie, gdy zawiera trafienie, żeby ścieżka pliku była widoczna.',
    }),
  );

  const odswiezPrzycisk = przyciskAkcji('Odśwież drzewo', 'dn-btn dn-btn--atrament');
  const zwinPrzycisk = przyciskAkcji('Zwiń wszystko');
  rama.akcje.append(odswiezPrzycisk, zwinPrzycisk);

  const korzen = document.createElement('p');
  korzen.className = 'mdev-sciezka';
  opiszKorzen(korzen, '');

  rama.cialo.append(korzen, stanTresci);
  return { sciezka, glebokosc, ukryte, szukaj, odswiezPrzycisk, zwinPrzycisk, korzen };
}

/** Zadanie odczytu drzewa z kontrolek Operatora; puste pola nie trafiają do żądania. */
function zadanieDrzewa(idOkna: string, powierzchnia: PowierzchniaDrzewa): DeveloperTreeGetRequest {
  const zadanie: DeveloperTreeGetRequest = { windowId: idOkna };
  const sciezka = powierzchnia.sciezka.value.trim();
  if (sciezka !== '') zadanie.path = sciezka;
  const glebokosc = Number.parseInt(powierzchnia.glebokosc.value, 10);
  if (Number.isInteger(glebokosc) && glebokosc > 0) zadanie.depth = glebokosc;
  if (powierzchnia.ukryte.checked) zadanie.includeHidden = true;
  return zadanie;
}

/** Węzeł drzewa hierarchicznego złożony z listy płaskiej po `parentPath`. */
interface WezelDrzewa {
  wezel: DeveloperTreeNode;
  dzieci: WezelDrzewa[];
  /** Powód, dla którego rdzeń nie pozwolił umieścić węzła u właściwego rodzica. */
  osierocony?: string;
}

/**
 * Składa listę płaską węzłów w hierarchię; katalogi przed plikami, alfabetycznie.
 *
 * Żaden węzeł nie znika. Węzły stoją w tablicy indeksowanej po pozycji, a nie
 * w mapie kluczowanej ścieżką, bo mapa gubi jeden z dwóch węzłów o tej
 * samej `path` (drugi nadpisuje pierwszy). Rodzic każdego węzła idzie przez
 * odrębną mapę ścieżka→pierwszy indeks, więc dwa węzły o tej samej ścieżce
 * zostają oba na scenie, a dzieci trafiają do pierwszego z nich. Samowskazanie
 * i cykl wzajemny (`wykryjCykle`) trafiają do korzenia z adnotacją, zamiast
 * znikać z drzewa razem z całym cyklem.
 */
function zbudujDrzewo(wezly: readonly DeveloperTreeNode[]): WezelDrzewa[] {
  const wpisy: WezelDrzewa[] = wezly.map((wezel) => ({ wezel, dzieci: [] }));
  const poSciezce = new Map<string, number>();
  wezly.forEach((wezel, indeks) => {
    if (!poSciezce.has(wezel.path)) poSciezce.set(wezel.path, indeks);
  });
  const rodzice: Array<number | null> = wezly.map((wezel) => {
    if (wezel.parentPath === undefined || wezel.parentPath === '') return null;
    return poSciezce.get(wezel.parentPath) ?? null;
  });
  // Rodzic wskazany, lecz nieobecny w wykazie: węzeł zostaje korzeniem, nie znika.
  wezly.forEach((wezel, indeks) => {
    if (wezel.parentPath !== undefined && wezel.parentPath !== '' && rodzice[indeks] === null) {
      wpisy[indeks].osierocony = `rdzeń podał rodzica „${wezel.parentPath}”, którego nie ma w wykazie`;
    }
  });
  const naCyklu = wykryjCykle(rodzice);
  rodzice.forEach((rodzic, indeks) => {
    if (!naCyklu[indeks]) return;
    wpisy[indeks].osierocony = rodzic === indeks
      ? 'rdzeń wskazał węzeł jako własnego rodzica'
      : 'rdzeń podał cykl rodziców — węzeł wraca po łańcuchu sam do siebie';
    rodzice[indeks] = null;
  });
  const korzenie: WezelDrzewa[] = [];
  rodzice.forEach((rodzic, indeks) => {
    if (rodzic !== null) wpisy[rodzic].dzieci.push(wpisy[indeks]);
    else korzenie.push(wpisy[indeks]);
  });
  posortujDrzewo(korzenie);
  return korzenie;
}

/**
 * Wykrywa węzły na cyklu wskazań rodzica (samowskazanie albo cykl wzajemny).
 * Barwienie DFS: węzeł „w toku”, który wraca sam do siebie po łańcuchu
 * rodziców, jest na cyklu razem z całym odcinkiem ścieżki od niego.
 */
function wykryjCykle(rodzice: ReadonlyArray<number | null>): boolean[] {
  const STAN_NOWY = 0;
  const STAN_W_TOKU = 1;
  const STAN_GOTOWY = 2;
  const stan = new Array<number>(rodzice.length).fill(STAN_NOWY);
  const naCyklu = new Array<boolean>(rodzice.length).fill(false);
  for (let start = 0; start < rodzice.length; start += 1) {
    if (stan[start] !== STAN_NOWY) continue;
    const sciezka: number[] = [];
    let biezacy: number | null = start;
    while (biezacy !== null && stan[biezacy] === STAN_NOWY) {
      stan[biezacy] = STAN_W_TOKU;
      sciezka.push(biezacy);
      biezacy = rodzice[biezacy] ?? null;
    }
    if (biezacy !== null && stan[biezacy] === STAN_W_TOKU) {
      const poczatekCyklu = sciezka.indexOf(biezacy);
      for (let i = poczatekCyklu; i < sciezka.length; i += 1) naCyklu[sciezka[i]] = true;
    }
    for (const wezel of sciezka) stan[wezel] = STAN_GOTOWY;
  }
  return naCyklu;
}

function posortujDrzewo(lista: WezelDrzewa[]): void {
  lista.sort((a, b) => {
    if (a.wezel.kind !== b.wezel.kind) return a.wezel.kind === TreeNodeKind.Directory ? -1 : 1;
    return a.wezel.name.localeCompare(b.wezel.name);
  });
  for (const wpis of lista) posortujDrzewo(wpis.dzieci);
}

/** Wywołania zwrotne widoku drzewa: stan wspólny, rozwinięcia lokalne, menu. */
interface ObslugaWezla {
  stan: StanDevelopera;
  rozwiniete: Set<string>;
  /** Fragment nazwy zawężający widok; pusty znaczy „bez zawężenia”. */
  fraza: string;
  naPrzelacz: () => void;
  naMenu: (event: MouseEvent, wezel: DeveloperTreeNode) => void;
}

/**
 * Zawęża drzewo do gałęzi zawierających trafienie.
 *
 * Katalog zostaje na scenie, gdy trafienie ma którykolwiek z jego potomków —
 * inaczej plik znaleziony trzy poziomy w głąb wisiałby bez ścieżki i nie dałoby
 * się powiedzieć, skąd pochodzi. Katalog trafiony własną nazwą zachowuje
 * komplet dzieci, bo szukającemu katalogu chodzi o jego zawartość.
 *
 * Zawężenie buduje nowe wpisy zamiast przycinać istniejące: `zbudujDrzewo`
 * biegnie od nowa przy każdym rysowaniu, ale adnotacja osierocenia ma przeżyć
 * zawężenie, więc jest przenoszona wprost.
 */
function zawezDrzewo(lista: readonly WezelDrzewa[], fraza: string): WezelDrzewa[] {
  const szukane = fraza.toLocaleLowerCase('pl-PL');
  const wynik: WezelDrzewa[] = [];
  for (const wpis of lista) {
    const trafiony = wpis.wezel.name.toLocaleLowerCase('pl-PL').includes(szukane);
    const dzieci = trafiony ? wpis.dzieci : zawezDrzewo(wpis.dzieci, fraza);
    if (!trafiony && dzieci.length === 0) continue;
    const kopia: WezelDrzewa = { wezel: wpis.wezel, dzieci };
    if (wpis.osierocony !== undefined) kopia.osierocony = wpis.osierocony;
    wynik.push(kopia);
  }
  return wynik;
}

/** Rysuje drzewo albo stan pustki — pustka bywa poprawna. */
function rysujDrzewoProjektu(
  wezly: readonly DeveloperTreeNode[],
  tresc: StanTresci,
  obsluga: ObslugaWezla,
): void {
  if (wezly.length === 0) {
    // Zdanie mówi, co oddał rdzeń, i nie orzeka przy okazji, że katalog
    // roboczy okna jest repozytorium — katalog sesji świeżego okna nie ma
    // `.git`, a `developer.git.action` odpowiada wtedy „not a git repository”.
    tresc.pusto('Rdzeń nie oddał ani jednego węzła dla tego katalogu.');
    return;
  }
  const wszystkie = zbudujDrzewo(wezly);
  // Węzły są, ale żadnego nie dało się umieścić w hierarchii — to nie jest
  // „treść” (miejsce byłoby puste bez wyjaśnienia), to nazwany błąd.
  if (wszystkie.length === 0) {
    tresc.blad('Rdzeń oddał węzły drzewa, których nie dało się złożyć w hierarchię.');
    return;
  }
  const korzenie = obsluga.fraza === '' ? wszystkie : zawezDrzewo(wszystkie, obsluga.fraza);
  // Pustka po zawężeniu to nie pustka katalogu — dwa różne zdania o tym samym
  // odczycie, więc i dwa różne komunikaty.
  if (korzenie.length === 0) {
    tresc.pusto(
      `Żadna pozycja odczytanego drzewa nie zawiera w nazwie „${obsluga.fraza}”. ` +
        'Zawężenie obejmuje wyłącznie węzły oddane przez rdzeń dla wskazanego katalogu ' +
        'i głębokości — nie przeszukuje repozytorium.',
    );
    return;
  }
  const lista = document.createElement('ul');
  lista.className = 'mdev-drzewo';
  lista.setAttribute('aria-label', 'Drzewo projektu');
  for (const wpis of korzenie) lista.append(rysujWezel(wpis, obsluga));
  tresc.tresc().append(lista);
}

/** Pojedynczy węzeł drzewa: nazwa klikalna, opis, ewentualne dzieci rozwinięte. */
function rysujWezel(wpis: WezelDrzewa, obsluga: ObslugaWezla): HTMLLIElement {
  const jestKatalogiem = wpis.wezel.kind === TreeNodeKind.Directory;
  const li = document.createElement('li');
  li.className = 'mdev-wezel';
  li.dataset['rodzaj'] = wpis.wezel.kind;
  li.dataset['wskazany'] = obsluga.stan.sciezka() === wpis.wezel.path ? 'tak' : 'nie';
  // Wskazany to nie to samo co wczytany: wskazanie biegnie natychmiast po
  // kliknięciu, a plik pojawia się w edytorze dopiero z odpowiedzią rdzenia —
  // która potrafi nie przyjść wcale (odmowa) albo przyjść dla innej ścieżki.
  const wczytany = obsluga.stan.plik()?.path === wpis.wezel.path;
  li.dataset['wczytany'] = wczytany ? 'tak' : 'nie';
  if (wpis.osierocony !== undefined) li.dataset['osierocony'] = 'tak';
  // Zawężenie rozwija gałęzie samo: trafienie ukryte w zwiniętym katalogu nie
  // byłoby trafieniem widocznym, a Operator wpisał frazę właśnie po to, żeby je
  // zobaczyć. Rozwinięcia własne Operatora zostają nietknięte i wracają
  // z chwilą wyczyszczenia pola.
  const rozwiniety = obsluga.fraza !== '' || obsluga.rozwiniete.has(wpis.wezel.path);

  const nazwa = document.createElement('button');
  nazwa.type = 'button';
  nazwa.className = 'mdev-wezel__nazwa';
  const strzalka = jestKatalogiem ? (rozwiniety ? '▾ ' : '▸ ') : '';
  nazwa.textContent = strzalka + wpis.wezel.name;
  nazwa.addEventListener('click', () => {
    if (jestKatalogiem) {
      // Katalog nie otwiera się jak plik — przełącza wyłącznie rozwinięcie lokalne.
      if (obsluga.rozwiniete.has(wpis.wezel.path)) obsluga.rozwiniete.delete(wpis.wezel.path);
      else obsluga.rozwiniete.add(wpis.wezel.path);
      obsluga.naPrzelacz();
    } else {
      obsluga.stan.wskazPlik(wpis.wezel.path);
    }
  });
  nazwa.addEventListener('contextmenu', (event) => {
    event.preventDefault();
    obsluga.naMenu(event, wpis.wezel);
  });

  const opis = document.createElement('span');
  opis.className = 'mdev-wezel__opis';
  opis.textContent = opisWezla(wpis, wczytany);
  li.append(nazwa, opis);
  if (jestKatalogiem && rozwiniety && wpis.dzieci.length > 0) {
    const dzieci = document.createElement('ul');
    dzieci.className = 'mdev-drzewo';
    for (const dziecko of wpis.dzieci) dzieci.append(rysujWezel(dziecko, obsluga));
    li.append(dzieci);
  }
  return li;
}

/** Zdanie opisowe węzła: rozmiar, czas zmiany, a na końcu powód osierocenia wprost. */
function opisWezla(wpis: WezelDrzewa, wczytany: boolean): string {
  const wezel = wpis.wezel;
  const czesci: string[] = [];
  if (wczytany) czesci.push('wczytany w edytorze');
  if (wezel.kind === TreeNodeKind.File && wezel.sizeBytes !== undefined) {
    czesci.push(`${wezel.sizeBytes} B`);
  }
  if (wezel.updatedAt !== undefined) {
    czesci.push(new Date(wezel.updatedAt).toLocaleString('pl-PL'));
  }
  if (wpis.osierocony !== undefined) czesci.push(`osierocony: ${wpis.osierocony}`);
  return czesci.join(' · ');
}

/** Menu kontekstowe prawego przycisku — jedyny udokumentowany przypadek platformy. */
interface MenuKontekstowe {
  pokaz(x: number, y: number, wezel: DeveloperTreeNode, tresc: StanTresci, naOdswiez: () => void): void;
}

function utworzMenuKontekstowe(): MenuKontekstowe {
  let biezace: HTMLElement | null = null;
  function zamknij(): void {
    biezace?.remove();
    biezace = null;
    document.removeEventListener('click', zamknij);
    document.removeEventListener('keydown', naEscape);
  }
  function naEscape(event: KeyboardEvent): void {
    if (event.key === 'Escape') zamknij();
  }
  return {
    pokaz(x, y, wezel, tresc, naOdswiez) {
      zamknij();
      const menu = document.createElement('ul');
      menu.className = 'mpt-menu';
      menu.setAttribute('role', 'menu');
      menu.style.left = `${x}px`;
      menu.style.top = `${y}px`;
      const dodaj = (element: HTMLElement): void => {
        const pozycja = document.createElement('li');
        pozycja.className = 'mpt-menu__pozycja';
        pozycja.append(element);
        menu.append(pozycja);
      };
      for (const [etykieta, powod] of brakiKomendMenu()) dodaj(przyciskBezKomendy(etykieta, powod));
      dodaj(zbudujCzynnoscMenu('Skopiuj ścieżkę', () => {
        void skopiujSciezke(wezel.path, tresc);
        zamknij();
      }));
      dodaj(zbudujCzynnoscMenu('Odśwież', () => {
        naOdswiez();
        zamknij();
      }));
      document.body.append(menu);
      biezace = menu;
      window.setTimeout(() => document.addEventListener('click', zamknij), 0);
      document.addEventListener('keydown', naEscape);
    },
  };
}

/** Przycisk czynności menu, realnej — zamyka menu po wykonaniu. */
function zbudujCzynnoscMenu(etykieta: string, naKlik: () => void): HTMLButtonElement {
  const przycisk = przyciskAkcji(etykieta, 'dn-btn dn-btn--zarys');
  przycisk.addEventListener('click', naKlik);
  return przycisk;
}

/**
 * Jedenaście pozycji menu bez pokrycia w kontrakcie.
 *
 * Powody składa `powodBezKomendy` z wykazu komend odczytanego w czasie
 * działania, a każda pozycja nazywa komendę, która musiałaby powstać — po to,
 * żeby przycisk i komenda dołożona później mówiły o czynności tym samym słowem.
 *
 * Wykaz powstaje wytwórnią, a nie stałą modułu: `powodBezKomendy` czyta wykaz
 * komend, więc stała liczona przy wczytaniu pliku zamroziłaby zdanie na stan
 * z chwili wczytania modułu.
 */
function brakiKomendMenu(): ReadonlyArray<readonly [string, string]> {
  const jedenPlikNaOkno =
    'Code Editor trzyma jeden plik na okno; wiele kart wymagałoby ponadto stanu wielu otwartych ' +
    'plików, którego kontrakt nie niesie.';
  return [
    [
      'Nowy plik',
      powodBezKomendy(
        'developer.tree.get jest wyłącznie do odczytu; założenie pliku wymagałoby komendy ' +
          'developer.file.create.',
      ),
    ],
    [
      'Nowy folder',
      powodBezKomendy(
        'Założenie katalogu wymagałoby komendy developer.file.create z rodzajem węzła „directory”.',
      ),
    ],
    ['Zmień nazwę', powodBezKomendy('Zmiana nazwy węzła wymagałaby komendy developer.file.rename.')],
    [
      'Usuń',
      powodBezKomendy(
        'Usunięcie węzła wymagałoby komendy developer.file.delete oddającej znacznik cofnięcia — ' +
          'opracowanie żąda usuwania od razu z „Cofnij” w powiadomieniu, nie modala-bramki.',
      ),
    ],
    [
      'Duplikuj',
      powodBezKomendy(
        'Duplikat wymagałby komendy developer.file.create z treścią pliku źródłowego.',
      ),
    ],
    ['Wytnij', powodBezKomendy('Przeniesienie węzła wymagałoby komendy developer.file.move.')],
    [
      'Kopiuj',
      powodBezKomendy('Kopia węzła wymagałaby komendy developer.file.create z treścią źródła.'),
    ],
    ['Wklej', powodBezKomendy('Wklejenie wymagałoby komendy developer.file.move albo file.create.')],
    ['Otwórz w nowej karcie', jedenPlikNaOkno],
    ['Otwórz obok', jedenPlikNaOkno],
    [
      'Pokaż w eksploratorze systemowym',
      'Przeglądarka nie ma dostępu do eksploratora plików Operatora. To granica środowiska ' +
        'uruchomienia klienta, nie brak komendy — żadna komenda rdzenia tego nie zmieni.',
    ],
  ];
}

/** Jedyna czynność menu czysto kliencka: schowek Operatora, żadnej komendy rdzenia. */
async function skopiujSciezke(sciezka: string, tresc: StanTresci): Promise<void> {
  try {
    await navigator.clipboard.writeText(sciezka);
    tresc.potwierdzenie(`Ścieżka skopiowana do schowka: ${sciezka}`, true);
  } catch {
    tresc.potwierdzenie('Nie udało się skopiować ścieżki — schowek niedostępny.', false);
  }
}
