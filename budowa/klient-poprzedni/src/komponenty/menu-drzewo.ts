import { elementIkony, type NazwaIkony } from '../ikony/ikony';

/**
 * Jeden mechanizm rozwijanego drzewa dla całej aplikacji.
 *
 * Stoi w bibliotece komponentów, a nie wewnątrz paska zlecenia, bo pasek jest
 * tylko jego największym odbiorcą: takie samo menu stawiają nagłówki okien,
 * panele i ekrany konfiguracji. Mechanizm zamknięty w pasku byłby dla nich
 * nieosiągalny i zostałby skopiowany, a kopie rozjeżdżają się w szczegółach —
 * haczyk wyboru pojawiałby się w jednym menu, a w drugim nie, przy jednakowym
 * wymogu.
 *
 * To jest mechanizm docelowy dla sterów nastawy: dopisanie obok niego kolejnego
 * mechanizmu rozwijania odtwarza dokładnie ten rozjazd, który on likwiduje.
 * Trwałym wyjątkiem zostaje `okna-rownolegle/menu-rozwijane.ts`, bo obsługuje
 * menu czynności (uchwyt jest ikoną, treść wchodzi jako gotowe elementy), a nie
 * ster niosący bieżącą wartość — to osobny wzorzec, nie rozjazd do scalenia.
 *
 * Mechanizm niesie sam sześć cech, których odbiorcy nie budują u siebie:
 *   1. zagnieżdżenie — `GalazMenu` wchodzi w `GalazMenu` bez ograniczenia głębokości,
 *   2. znacznik bieżącego wyboru — `aria-checked` na liściu i ślad na gałęzi,
 *      która ten liść niesie, żeby Operator widział wybór bez wchodzenia w gałąź,
 *   3. opis przy pozycji — `opis` jest polem każdej pozycji,
 *   4. przełącznik dwustanowy wewnątrz gałęzi — `PrzelacznikMenu`,
 *   5. grupowanie po źródle lub rodzinie — `GrupaMenu`,
 *   6. droga do rejestru na dole — `stopka`.
 * Treść wchodzi tu jako dane, a formę nadaje mechanizm; przyjmowanie gotowych
 * elementów od odbiorcy zniosłoby jednolitość tych sześciu cech.
 *
 * Nie wie, co jest w drzewie, i nie wykonuje wyboru — oddaje klucz pozycji
 * wołającemu. Nie zna kontraktu, komendy ani stanu okna.
 *
 * Żadna pozycja nie dostaje `disabled`. Drzewo puste nie jest błędem: uchwyt
 * otwiera się i mówi zdaniem, że wykaz jest pusty, zamiast przestać reagować.
 */

/** Wspólne pola każdej pozycji drzewa. */
interface PozycjaWspolna {
  /** Klucz oddawany wołającemu przy wyborze; niepowtarzalny w obrębie drzewa. */
  klucz: string;
  /** Nazwa widoczna w menu. */
  nazwa: string;
  /**
   * Druga nazwa tej samej pozycji — krótsza, bez źródła.
   *
   * W wykazie płaskim na setki pozycji nazwa pełna musi nieść źródło, bo inaczej
   * dwie komendy o tej samej nazwie z dwóch wtyczek są nie do rozróżnienia
   * („anthropic-skills:skill-creator" kontra „moje:skill-creator"). To samo
   * źródło czyni jednak wykaz nieczytelnym, gdy sto pozycji zaczyna się tym
   * samym przedrostkiem. Dlatego pozycja niesie obie nazwy: pełną ze źródłem
   * i skróconą do samej rzeczy.
   *
   * Pominięta znaczy „ta pozycja ma jedną nazwę".
   *
   * Szukanie obejmuje obie nazwy i opis; trafienie w skróconą liczy się wyżej
   * niż w pełną, bo to ono jest tym, czego Operator naprawdę szukał.
   */
  nazwaKrotka?: string;
  /**
   * Zdanie mówiące, co ta pozycja robi.
   *
   * Stoi wprost pod nazwą, a nie pod znakiem zapytania, bo w menu o kilkudziesięciu
   * liściach najechanie na każdy liść z osobna nie jest drogą. Dymek
   * (`komponenty/dymek.ts`) zostaje formą dla kontrolek stojących pojedynczo;
   * drugiego mechanizmu objaśniania tu nie ma.
   *
   * Pominięty znaczy „ta pozycja opisu nie ma".
   */
  opis?: string;
  /** Ikona przy nazwie; pominięta znaczy „bez ikony". */
  ikona?: NazwaIkony;
}

/** Gałąź — pozycja prowadząca do kolejnego poziomu. */
export interface GalazMenu extends PozycjaWspolna {
  rodzaj: 'galaz';
  dzieci: readonly PozycjaMenu[];
}

/** Liść wyboru jednokrotnego — dokładnie jeden w drzewie bywa `wybrany`. */
export interface LiscMenu extends PozycjaWspolna {
  rodzaj: 'wybor';
  wybrany: boolean;
}

/** Nastawa dwustanowa siedząca w menu, nie obok niego. */
export interface PrzelacznikMenu extends PozycjaWspolna {
  rodzaj: 'przelacznik';
  wlaczony: boolean;
}

/**
 * Grupa pozycji jednej rodziny albo jednego źródła.
 *
 * Nie jest poziomem drzewa: jej dzieci stoją na tym samym poziomie co ona,
 * pod wspólnym nagłówkiem. Poziom otwiera wyłącznie gałąź.
 */
export interface GrupaMenu {
  rodzaj: 'grupa';
  nazwa: string;
  dzieci: readonly PozycjaMenu[];
}

export type PozycjaMenu = GalazMenu | LiscMenu | PrzelacznikMenu | GrupaMenu;

/** Stopka menu — droga do rejestru; podłączenie nie jest tym samym co rejestracja. */
export interface StopkaMenu {
  nazwa: string;
  opis?: string;
  ikona?: NazwaIkony;
  wykonaj(): void;
}

export interface OpcjeMenuDrzewa {
  /**
   * Nazwa rodzajowa nastawy — „Model", „Wysiłek", „Urządzenie".
   *
   * Nie trafia na ekran, bo na uchwycie stoi wartość, a nie nazwa nastawy.
   * Idzie do `aria-label` uchwytu i listy: czytnik ekranu musi wiedzieć, czego
   * dotyczy wartość, której nazwa sama tego nie mówi („Opus 5" nie niesie słowa
   * „model").
   */
  nastawa: string;
  /** Ikona stojąca na uchwycie przed wartością; pominięta znaczy „bez ikony". */
  ikona?: NazwaIkony;
  /** Wołane po wyborze liścia albo przełącznika; klucz pozycji. */
  naWybor(klucz: string): void;
  /** Droga do rejestru na dole menu; pominięta znaczy „to menu jej nie ma". */
  stopka?: StopkaMenu;
  /**
   * Od ilu liści menu stawia pole szukania. Domyślnie 12.
   *
   * Przy dużym wykazie samo drzewo nie wystarczy: droga do liścia na czwartym
   * poziomie jest dłuższa niż cierpliwość Operatora, więc od pewnej skali
   * wyszukiwanie jest obowiązkowe. Próg liczy liście, nie pozycje — drzewo
   * o trzech gałęziach i dwóch liściach filtra nie potrzebuje, drzewo o dwóch
   * gałęziach i stu ekspertach potrzebuje go natychmiast.
   *
   * Próg, a nie „zawsze": pole szukania nad wykazem pięciu pozycji zabiera
   * wiersz i nie skraca ani jednego ruchu.
   */
  progSzukania?: number;
  /**
   * W którą stronę rozwija się wykaz. Domyślnie `'dol'`.
   *
   * Nie jest preferencją wizualną. Ster stojący w pasku u góry okna ma pod sobą
   * całą wysokość sceny i rozwija się w dół. Ster przy polu wpisywania stoi
   * u dołu okna i wykaz rozwinięty w dół nie ma dokąd pójść — wyszedłby poza
   * krawędź sceny.
   */
  kierunek?: 'dol' | 'gora';
  /**
   * Menu bez własnego uchwytu — rozwijane i filtrowane z zewnątrz.
   *
   * Tryb dla wykazu komend po ukośniku: wykaz ma się pojawić natychmiast po
   * wpisaniu `/`, a filtrem ma być to samo pole wpisywania, nie osobne okienko.
   * Menu z własnym uchwytem i własnym polem szukania stawiałoby drugie pole obok
   * pierwszego i tego wymogu nie spełnia.
   *
   * W tym trybie uchwyt nie wchodzi do dokumentu, wewnętrzne pole szukania nie
   * powstaje wcale, a wołający steruje mechanizmem przez `rozwin`, `ustawFraze`,
   * `przesunWyroznienie` i `wybierzWyrozniona`. Ognisko zostaje w jego polu —
   * stąd wyróżnienie jest wirtualne (`aria-activedescendant`), a nie ogniskiem.
   */
  bezUchwytu?: boolean;
  /**
   * Opis rysuje się tylko przy pozycji wyróżnionej, zamiast przy każdej.
   *
   * Domyślnie opis stoi przy każdej pozycji i tak zostaje przy drzewie
   * o kilkudziesięciu liściach. Wykaz płaski na setki pozycji to inna skala:
   * setka opisów naraz nie jest objaśnieniem, tylko ścianą tekstu, przez którą
   * nie widać już samych nazw. Wtedy opis należy się jednej pozycji — tej, na
   * którą Operator właśnie patrzy.
   */
  opisTylkoPrzyWyroznionej?: boolean;
}

export interface MenuDrzewo {
  /** Element montowany w pasku albo w nagłówku okna. */
  element: HTMLElement;
  /**
   * Podaje mechanizmowi całą zawartość naraz: wartość na uchwyt i drzewo.
   *
   * `wartosc` jest napisem widocznym na uchwycie — bieżącą nastawą. Drzewo
   * przerysowuje się w całości, bo wykaz bywa czytany z rdzenia i pozycje
   * przychodzą oraz znikają; rozwinięte gałęzie przeżywają przerysowanie po
   * kluczu, żeby odświeżenie katalogu nie zwijało menu pod ręką Operatora.
   */
  ustaw(wartosc: string, drzewo: readonly PozycjaMenu[]): void;
  /** Zwija menu wraz z gałęziami i zdejmuje nasłuchy dokumentu. Obowiązkowe. */
  zwin(): void;
  /**
   * Rozwija wykaz bez klikania uchwytu — jedyna droga dla trybu `bezUchwytu`.
   *
   * Ognisko zostaje tam, gdzie było. Mechanizm nie zabiera go wołającemu, bo
   * przy obsadzie ukośnikiem to pole wpisywania jest miejscem, w którym Operator
   * pisze dalej.
   */
  rozwin(): void;
  /**
   * Podaje frazę filtrującą z zewnątrz — z pola wpisywania wołającego.
   *
   * Gdy wewnętrzne pole szukania stoi na ekranie, ta fraza w nie wchodzi, żeby
   * Operator widział, czym wykaz jest przycięty. Gdy nie stoi (tryb `bezUchwytu`
   * albo wykaz poniżej progu szukania) — fraza i tak przycina wykaz. Próg
   * rozstrzyga, czy mechanizm stawia pole, a nie czy w ogóle umie filtrować.
   */
  ustawFraze(fraza: string): void;
  /**
   * Przesuwa wyróżnienie o `krok` pozycji, nie ruszając ogniska.
   *
   * Dla obsady ukośnikiem strzałki naciska się w polu wpisywania, więc ognisko
   * musi w nim zostać; przenoszenie go na pozycję wyrwałoby Operatorowi klawiaturę
   * spod palców w połowie pisania. Wykaz zawija się na obu końcach.
   */
  przesunWyroznienie(krok: number): void;
  /**
   * Wybiera pozycję wyróżnioną — Enter działa natychmiast, bez dodatkowego ruchu.
   *
   * Oddaje `false`, gdy nie ma czego wybrać (wykaz pusty albo przycięty do zera),
   * i wtedy wołający wie, że Enter ma zrobić swoje zwykłe zadanie zamiast niczego.
   * Wyróżniona gałąź nie jest wyborem — Enter na niej otwiera poziom.
   */
  wybierzWyrozniona(): boolean;
}

/** Zdanie o pustym wykazie — menu mówi o braku, zamiast przestać reagować. */
const ZDANIE_PUSTEGO = 'Wykaz jest pusty — nie ma tu jeszcze ani jednej pozycji do wyboru.';

/** Od ilu liści drzewo dostaje pole szukania, gdy wołający nie rozstrzygnie. */
const PROG_SZUKANIA = 12;

/* ---------------------------------------------------------------------------
   Ocena trafienia — z czego bierze się porządek wykazu przy czynnej frazie
   ---------------------------------------------------------------------------
   Kolejność wejścia jest właściwa dla drzewa ułożonego ręcznie i bezużyteczna
   dla wykazu płaskiego na setki pozycji: Operator wpisuje trzy litery i szuka
   swojej pozycji na czele wykazu, nie w kolejności rejestru. Stąd ocena.

   Trzy progi po sto, żeby trafienie w lepszym polu zawsze biło trafienie
   w gorszym, choćby tamto stało na samym początku napisu. Odsunięcie trafienia
   od początku odejmuje w obrębie progu i jest przycięte do 99, żeby nigdy nie
   przelało się do progu niższego.
   --------------------------------------------------------------------------- */
const OCENA_SKROT = 300;
const OCENA_NAZWA = 200;
const OCENA_OPIS = 100;
const OCENA_ODSUNIECIE_MAX = 99;
/** Ocena „ta pozycja nie pasuje". Wszystko od zera w górę jest trafieniem. */
const OCENA_BRAK = -1;

/** Ocena jednego pola; `szukane` przychodzi już małymi literami i przycięte. */
function ocenaPola(tekst: string | undefined, baza: number, szukane: string): number {
  if (tekst === undefined || tekst === '') return OCENA_BRAK;
  const gdzie = tekst.toLocaleLowerCase('pl-PL').indexOf(szukane);
  if (gdzie < 0) return OCENA_BRAK;
  return baza - Math.min(gdzie, OCENA_ODSUNIECIE_MAX);
}

/** Najlepsze trafienie pozycji: nazwa skrócona bije pełną, pełna bije opis. */
function ocenaPozycji(pozycja: PozycjaWspolna, szukane: string): number {
  return Math.max(
    ocenaPola(pozycja.nazwaKrotka, OCENA_SKROT, szukane),
    ocenaPola(pozycja.nazwa, OCENA_NAZWA, szukane),
    ocenaPola(pozycja.opis, OCENA_OPIS, szukane),
  );
}

/**
 * Najlepsze trafienie w całym poddrzewie — ocena gałęzi bierze się z dzieci.
 *
 * Gałąź, pod którą siedzi trafienie idealne, ma stanąć wyżej niż gałąź
 * z trafieniem bylejakim. Bez tego gałęzie ustawiłyby się kolejnością wejścia
 * i porządek wedle trafności kończyłby się na pierwszym poziomie.
 */
function najlepszaOcena(pozycje: readonly PozycjaMenu[], szukane: string): number {
  let najlepsza = OCENA_BRAK;
  for (const p of pozycje) {
    // Grupa nie ma własnego trafienia — jest nagłówkiem, nie pozycją do wyboru.
    if (p.rodzaj !== 'grupa') najlepsza = Math.max(najlepsza, ocenaPozycji(p, szukane));
    if (p.rodzaj === 'galaz' || p.rodzaj === 'grupa') {
      najlepsza = Math.max(najlepsza, najlepszaOcena(p.dzieci, szukane));
    }
  }
  return najlepsza;
}

/**
 * Wstawia napis do elementu, wytłuszczając każde wystąpienie frazy.
 *
 * Każde, nie pierwsze: przy nazwie „skill-creator-skill" Operator szukający
 * „skill" widzi, że trafił dwa razy, i wie, czy fraza jest dość wyostrzona.
 * Bez frazy wchodzi jeden węzeł tekstowy.
 */
function wstawTrafienia(cel: HTMLElement, napis: string, szukane: string): void {
  if (szukane === '') {
    cel.textContent = napis;
    return;
  }
  const male = napis.toLocaleLowerCase('pl-PL');
  // Zmiana wielkości liter bywa zmianą długości (np. „İ" schodzi na dwa znaki).
  // Wtedy indeksy z wersji małej nie pasują do oryginału i cięcie rozerwałoby
  // napis w złym miejscu, pokazując przekłamaną nazwę. Napis zostaje w całości,
  // bez wytłuszczenia.
  if (male.length !== napis.length) {
    cel.textContent = napis;
    return;
  }

  const czesci: Node[] = [];
  let od = 0;
  for (;;) {
    const gdzie = male.indexOf(szukane, od);
    if (gdzie < 0) break;
    if (gdzie > od) czesci.push(document.createTextNode(napis.slice(od, gdzie)));
    const traf = document.createElement('strong');
    traf.className = 'dn-drzewo__traf';
    traf.textContent = napis.slice(gdzie, gdzie + szukane.length);
    czesci.push(traf);
    od = gdzie + szukane.length;
  }
  // Trafienia nie ma, gdy pozycja weszła do wykazu przez inne pole — nazwa
  // zostaje wtedy zwykłym napisem, bo wytłuszczać nie ma czego.
  if (czesci.length === 0) {
    cel.textContent = napis;
    return;
  }
  if (od < napis.length) czesci.push(document.createTextNode(napis.slice(od)));
  cel.replaceChildren(...czesci);
}

/** Licznik egzemplarzy — daje niepowtarzalny przedrostek `id` pozycji. */
let licznikMenu = 0;

export function utworzMenuDrzewo(opcje: OpcjeMenuDrzewa): MenuDrzewo {
  let otwarte = false;
  /** Klucze gałęzi rozwiniętych — przeżywają przerysowanie drzewa. */
  const rozwiniete = new Set<string>();
  let drzewo: readonly PozycjaMenu[] = [];
  const bezUchwytu = opcje.bezUchwytu === true;
  /**
   * Jedyne źródło prawdy o frazie — wewnętrzne pole tylko je odzwierciedla.
   *
   * Przy dwóch źródłach (pole wewnętrzne i fraza z zewnątrz) trzeba by
   * rozstrzygać pierwszeństwo, a każde rozstrzygnięcie zawodzi w drugą stronę.
   * Pole, gdy stoi, wpisuje tutaj przy każdym uderzeniu w klawisz; wykaz czyta
   * stąd, nie z pola.
   */
  let frazaBiezaca = '';
  /** Fraza użyta przy ostatnim odrysowaniu — z niej bierze się wytłuszczenie. */
  let frazaCzynna = '';
  /** Pozycja wyróżniona; przy obsadzie z zewnątrz zastępuje ognisko. */
  let wyrozniona: HTMLElement | null = null;
  const przedrostekId = `dn-drzewo-${++licznikMenu}-`;
  let licznikPozycji = 0;

  const element = document.createElement('div');
  element.className = 'dn-drzewo';
  element.dataset['kierunek'] = opcje.kierunek ?? 'dol';
  if (bezUchwytu) element.dataset['bezUchwytu'] = 'tak';

  const uchwyt = document.createElement('button');
  uchwyt.type = 'button';
  uchwyt.className = 'dn-drzewo__uchwyt';
  uchwyt.setAttribute('aria-haspopup', 'menu');
  uchwyt.setAttribute('aria-expanded', 'false');

  if (opcje.ikona !== undefined) {
    uchwyt.append(elementIkony(opcje.ikona, { rozmiar: 14, klasa: 'dn-drzewo__ikona' }));
  }

  // Na uchwycie stoi wartość, nie nazwa nastawy — tym ster różni się od
  // wyświetlacza i od menu czynności z `menu-rozwijane.ts`.
  const wartosc = document.createElement('span');
  wartosc.className = 'dn-drzewo__wartosc';
  uchwyt.append(wartosc);
  uchwyt.append(elementIkony('grot-dol', { rozmiar: 12, klasa: 'dn-drzewo__grot' }));

  const lista = document.createElement('div');
  lista.className = 'dn-drzewo__lista';
  lista.setAttribute('role', 'menu');
  lista.setAttribute('aria-label', opcje.nastawa);
  lista.hidden = true;

  // Pole szukania stoi poza `role="menu"`, nad wykazem. Wewnątrz menu byłoby
  // dla czytnika ekranu pozycją menu, którą nie jest, i wchodziłoby w wędrówkę
  // strzałkami — a strzałki mają chodzić po pozycjach, nie po polu tekstowym.
  const fraza = document.createElement('input');
  fraza.type = 'search';
  fraza.className = 'dn-drzewo__szukanie';
  fraza.placeholder = 'Szukaj…';
  fraza.setAttribute('aria-label', `Szukanie w wykazie: ${opcje.nastawa}`);
  fraza.hidden = true;
  // Strzałka w dół z pola przenosi ognisko na pierwsze trafienie — bez tego
  // filtrowanie kończyłoby się koniecznością sięgnięcia po mysz.
  fraza.addEventListener('input', () => {
    frazaBiezaca = fraza.value;
    odrysuj();
  });
  fraza.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'ArrowDown' && zdarzenie.key !== 'Enter') return;
    zdarzenie.preventDefault();
    wedrowne()[0]?.focus();
  });

  // W trybie bez uchwytu ani uchwyt, ani wewnętrzne pole szukania nie wchodzą
  // do dokumentu: wyzwalaczem jest ukośnik w polu wpisywania wołającego, a filtrem
  // to samo pole. Drugie pole szukania obok pierwszego byłoby dublowaniem.
  if (bezUchwytu) element.append(lista);
  else element.append(uchwyt, fraza, lista);

  uchwyt.addEventListener('click', () => ustawOtwarte(!otwarte));
  uchwyt.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'ArrowDown' && zdarzenie.key !== 'ArrowUp') return;
    zdarzenie.preventDefault();
    ustawOtwarte(true);
    const wykaz = wedrowne();
    const cel = zdarzenie.key === 'ArrowDown' ? wykaz[0] : wykaz[wykaz.length - 1];
    cel?.focus();
  });

  // Wędrówka ognisk i sterowanie poziomami stoją na całym menu, nie na liście:
  // ognisko bywa na uchwycie, na gałęzi i na liściu podgałęzi.
  element.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Escape') {
      if (!otwarte) return;
      zdarzenie.stopPropagation();
      zwinWszystko();
      uchwyt.focus();
      return;
    }
    if (zdarzenie.key === 'ArrowDown' || zdarzenie.key === 'ArrowUp') {
      if (document.activeElement === uchwyt) return;
      zdarzenie.preventDefault();
      przesun(zdarzenie.key === 'ArrowDown' ? 1 : -1);
      return;
    }
    // Prawo otwiera gałąź, lewo ją zamyka — dwa klawisze na dwa kierunki
    // poziomu, tak samo na każdej głębokości.
    const cel = document.activeElement;
    if (!(cel instanceof HTMLElement) || cel.dataset['galaz'] === undefined) return;
    const klucz = cel.dataset['galaz'];
    if (zdarzenie.key === 'ArrowRight' && !rozwiniete.has(klucz)) {
      zdarzenie.preventDefault();
      rozwiniete.add(klucz);
      odrysuj();
      ogniskoNaKlucz(klucz);
      return;
    }
    if (zdarzenie.key === 'ArrowLeft' && rozwiniete.has(klucz)) {
      zdarzenie.preventDefault();
      rozwiniete.delete(klucz);
      odrysuj();
      ogniskoNaKlucz(klucz);
    }
  });

  // Ognisko i wyróżnienie wskazują tę samą pozycję. Przy obsadzie z uchwytem
  // strzałki przenoszą ognisko naprawdę; wyróżnienie pozostawione na pierwszej
  // pozycji dałoby dwa sprzeczne wskazania naraz.
  lista.addEventListener('focusin', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof HTMLElement)) return;
    if (cel.getAttribute('role') === null) return;
    if (cel.classList.contains('dn-drzewo__pozycja--stopka')) return;
    ustawWyrozniona(cel);
  });

  /** Pozycje, po których chodzi ognisko — wyłącznie widoczne w tej chwili. */
  function wedrowne(): HTMLElement[] {
    const wybor = '[role="menuitem"], [role="menuitemcheckbox"], [role="menuitemradio"]';
    return [...lista.querySelectorAll<HTMLElement>(wybor)].filter((poz) => poz.offsetParent !== null);
  }

  /**
   * Pozycje wykazu w porządku rysowania — bez pytania o `offsetParent`.
   *
   * `wedrowne()` odsiewa po `offsetParent`, bo ognisko ma chodzić tylko po tym,
   * co Operator widzi. Wyróżnienie nie może tak pytać: przy obsadzie z zewnątrz
   * menu bywa wyróżniane, zanim układ policzy geometrię, a zwinięte gałęzie i tak
   * nie mają dzieci w dokumencie. Stopka odpada, bo nie jest wyborem.
   */
  function pozycjeWykazu(): HTMLElement[] {
    const wybor =
      '[role="menuitem"]:not(.dn-drzewo__pozycja--stopka), ' +
      '[role="menuitemcheckbox"], [role="menuitemradio"]';
    return [...lista.querySelectorAll<HTMLElement>(wybor)];
  }

  /**
   * Przenosi wyróżnienie; `null` znaczy „nie ma czego wyróżnić".
   *
   * Wyróżnienie jest wirtualne — nie rusza ogniska. Przy obsadzie ukośnikiem
   * ognisko siedzi w polu wpisywania wołającego i ma tam zostać, więc czytnik
   * ekranu dowiaduje się o bieżącej pozycji z `aria-activedescendant`, a nie
   * z przeskoku ogniska.
   */
  function ustawWyrozniona(cel: HTMLElement | null): void {
    if (wyrozniona !== null && wyrozniona !== cel) delete wyrozniona.dataset['wyrozniona'];
    wyrozniona = cel;
    if (cel === null) {
      lista.removeAttribute('aria-activedescendant');
    } else {
      // `id` nadaje się dopiero tutaj i tylko raz: `aria-activedescendant`
      // wskazuje identyfikatorem, więc pozycja wyróżniona musi go mieć.
      if (cel.id === '') cel.id = `${przedrostekId}${++licznikPozycji}`;
      cel.dataset['wyrozniona'] = 'tak';
      lista.setAttribute('aria-activedescendant', cel.id);
    }
    zastosujOpisy();
  }

  /** Chowa opisy poza wyróżnioną — wyłącznie gdy wołający o to poprosił. */
  function zastosujOpisy(): void {
    if (opcje.opisTylkoPrzyWyroznionej !== true) return;
    for (const opis of lista.querySelectorAll<HTMLElement>('.dn-drzewo__opis')) {
      opis.hidden = wyrozniona === null || !wyrozniona.contains(opis);
    }
  }

  function przesun(krok: number): void {
    const wykaz = wedrowne();
    if (wykaz.length === 0) return;
    const teraz = wykaz.indexOf(document.activeElement as HTMLElement);
    wykaz[(teraz + krok + wykaz.length) % wykaz.length]?.focus();
  }

  function ogniskoNaKlucz(klucz: string): void {
    lista.querySelector<HTMLElement>(`[data-galaz="${CSS.escape(klucz)}"]`)?.focus();
  }

  function poKlikuPoza(zdarzenie: Event): void {
    const cel = zdarzenie.target;
    if (cel instanceof Node && element.contains(cel)) return;
    zwinWszystko();
  }

  function ustawOtwarte(nowe: boolean): void {
    if (nowe === otwarte) return;
    otwarte = nowe;
    lista.hidden = !otwarte;
    uchwyt.setAttribute('aria-expanded', String(otwarte));
    element.dataset['otwarte'] = otwarte ? 'tak' : 'nie';
    if (otwarte) document.addEventListener('pointerdown', poKlikuPoza, true);
    else document.removeEventListener('pointerdown', poKlikuPoza, true);
  }

  function zwinWszystko(): void {
    rozwiniete.clear();
    odrysuj();
    ustawOtwarte(false);
  }

  /** Wybór liścia albo przełącznika zwija menu i oddaje klucz wołającemu. */
  function wybierz(klucz: string): void {
    zwinWszystko();
    uchwyt.focus();
    opcje.naWybor(klucz);
  }

  /** Liczba liści w poddrzewie — próg szukania liczy je, nie pozycje. */
  function liczbaLisci(pozycje: readonly PozycjaMenu[]): number {
    return pozycje.reduce((suma, p) => {
      if (p.rodzaj === 'galaz' || p.rodzaj === 'grupa') return suma + liczbaLisci(p.dzieci);
      return suma + 1;
    }, 0);
  }

  /**
   * Drzewo przycięte frazą — gałąź zostaje, jeśli cokolwiek pod nią pasuje.
   *
   * Filtrowanie nie łamie ujawniania stopniowego, tylko je skraca: fraza jest
   * drogą na skróty do liścia, który i tak istnieje. Gałąź z trafieniem
   * rozwija się sama, bo Operator szukający po nazwie nie ma jak wiedzieć,
   * w której gałęzi jego trafienie siedzi.
   */
  function przytnij(pozycje: readonly PozycjaMenu[], szukane: string): PozycjaMenu[] {
    const wynik: { pozycja: PozycjaMenu; ocena: number }[] = [];
    for (const p of pozycje) {
      if (p.rodzaj === 'galaz' || p.rodzaj === 'grupa') {
        const dzieci = przytnij(p.dzieci, szukane);
        const sama = p.rodzaj === 'galaz' ? ocenaPozycji(p, szukane) : OCENA_BRAK;
        if (dzieci.length === 0 && sama < 0) continue;
        if (p.rodzaj === 'galaz') rozwiniete.add(p.klucz);
        // Gałąź trafiona własną nazwą zachowuje całe poddrzewo: skoro Operator trafił
        // w jej nazwę, szuka gałęzi, a nie jednego liścia pod nią.
        const zawartosc = dzieci.length === 0 ? p.dzieci : dzieci;
        const ocena = Math.max(sama, najlepszaOcena(zawartosc, szukane));
        wynik.push({ pozycja: { ...p, dzieci: zawartosc }, ocena });
        continue;
      }
      const ocena = ocenaPozycji(p, szukane);
      if (ocena >= 0) wynik.push({ pozycja: p, ocena });
    }
    // Porządek wedle trafności obowiązuje wyłącznie przy czynnej frazie:
    // `przytnij` bez frazy nie jest wołane wcale, więc drzewo bez szukania
    // zostaje w porządku podanym przez wołającego. Sortowanie w JavaScripcie
    // jest stabilne (ES2019), więc remis oceny oddaje kolejność wejścia sam
    // z siebie i nie trzeba go rozstrzygać drugim kluczem.
    wynik.sort((a, b) => b.ocena - a.ocena);
    return wynik.map((w) => w.pozycja);
  }

  /** Fraza filtrująca — z pola wewnętrznego albo z zewnątrz, jedno źródło. */
  function biezacaFraza(): string {
    return frazaBiezaca.trim().toLocaleLowerCase('pl-PL');
  }

  function odrysuj(): void {
    const szukane = biezacaFraza();
    // Musi stanąć przed budową węzłów: `wiersz()` czyta stąd frazę do wytłuszczenia.
    frazaCzynna = szukane;
    const widoczne = szukane === '' ? drzewo : przytnij(drzewo, szukane);

    let dzieci: HTMLElement[];
    if (drzewo.length === 0) dzieci = [zdaniePustego()];
    else if (widoczne.length === 0) dzieci = [zdanieBrakuTrafien(frazaBiezaca.trim())];
    else dzieci = widoczne.flatMap((p) => wezel(p, 0));

    if (opcje.stopka !== undefined) dzieci.push(kreska(), stopka(opcje.stopka));
    lista.replaceChildren(...dzieci);

    // Pierwsza pozycja jest wyróżniona od razu, żeby Enter działał natychmiast.
    // Bez tego Operator z jednym trafieniem musiałby jeszcze sięgnąć po strzałkę
    // albo mysz.
    ustawWyrozniona(pozycjeWykazu()[0] ?? null);
  }

  /** Jedna pozycja drzewa; `poziom` steruje wcięciem, nie zachowaniem. */
  function wezel(pozycja: PozycjaMenu, poziom: number): HTMLElement[] {
    if (pozycja.rodzaj === 'grupa') {
      const naglowek = document.createElement('p');
      naglowek.className = 'dn-drzewo__grupa';
      naglowek.textContent = pozycja.nazwa;
      naglowek.style.setProperty('--dn-drzewo-poziom', String(poziom));
      return [naglowek, ...pozycja.dzieci.flatMap((d) => wezel(d, poziom))];
    }

    if (pozycja.rodzaj === 'galaz') {
      const { guzik } = wiersz(pozycja, poziom, 'menuitem');
      guzik.dataset['galaz'] = pozycja.klucz;
      guzik.dataset['klucz'] = pozycja.klucz;
      const otwarta = rozwiniete.has(pozycja.klucz);
      guzik.setAttribute('aria-haspopup', 'menu');
      guzik.setAttribute('aria-expanded', String(otwarta));
      // Ślad wyboru na gałęzi. Bez niego Operator musiałby wejść w każdą gałąź,
      // żeby się dowiedzieć, w której siedzi jego wybór — menu pokazuje wybór,
      // zamiast kazać go szukać.
      if (niesieWybor(pozycja.dzieci)) guzik.dataset['znacznik'] = 'tak';
      guzik.append(elementIkony('grot-prawo', { rozmiar: 12, klasa: 'dn-drzewo__grot-galezi' }));
      guzik.addEventListener('click', () => {
        if (otwarta) rozwiniete.delete(pozycja.klucz);
        else rozwiniete.add(pozycja.klucz);
        odrysuj();
        ogniskoNaKlucz(pozycja.klucz);
      });
      if (!otwarta) return [guzik];
      const podlista = document.createElement('div');
      podlista.className = 'dn-drzewo__podlista';
      podlista.setAttribute('role', 'menu');
      podlista.setAttribute('aria-label', pozycja.nazwa);
      podlista.append(...pozycja.dzieci.flatMap((d) => wezel(d, poziom + 1)));
      return [guzik, podlista];
    }

    const rola = pozycja.rodzaj === 'wybor' ? 'menuitemradio' : 'menuitemcheckbox';
    const { guzik, haczyk } = wiersz(pozycja, poziom, rola);
    guzik.dataset['klucz'] = pozycja.klucz;
    const czynna = pozycja.rodzaj === 'wybor' ? pozycja.wybrany : pozycja.wlaczony;
    guzik.setAttribute('aria-checked', String(czynna));
    // Znak wyboru stoi w każdym liściu, nie tylko w wybranym — miejsce zajęte na
    // stałe nie pozwala wierszom skakać w bok przy przełączeniu nastawy. Gałąź
    // znaku nie dostaje: jej ślad wyboru jest kropką arkusza, a nie ikoną.
    haczyk.append(elementIkony('ptaszek', { rozmiar: 14 }));
    guzik.addEventListener('click', () => wybierz(pozycja.klucz));
    return [guzik];
  }

  /** Czy w poddrzewie siedzi liść wybrany albo przełącznik włączony. */
  function niesieWybor(dzieci: readonly PozycjaMenu[]): boolean {
    return dzieci.some((d) => {
      if (d.rodzaj === 'wybor') return d.wybrany;
      if (d.rodzaj === 'przelacznik') return d.wlaczony;
      return niesieWybor(d.dzieci);
    });
  }

  /**
   * Szkielet wiersza: haczyk · ikona · nazwa i opis. Jeden rytm na wszystkie.
   *
   * Haczyk wraca osobno, bo tylko liść stawia w nim znak wyboru, a znaku tego
   * nie da się wstawić arkuszem: jest ikoną zestawu, a arkusz zestawu nie zna.
   */
  function wiersz(
    pozycja: PozycjaWspolna,
    poziom: number,
    rola: string,
  ): { guzik: HTMLButtonElement; haczyk: HTMLElement } {
    const guzik = document.createElement('button');
    guzik.type = 'button';
    guzik.className = 'dn-drzewo__pozycja';
    guzik.setAttribute('role', rola);
    guzik.style.setProperty('--dn-drzewo-poziom', String(poziom));

    // Miejsce na haczyk jest stałe, także gdy haczyka nie widać: inaczej wiersze
    // przeskakiwałyby w bok przy zmianie wyboru. Widoczność rozstrzyga arkusz po
    // `aria-checked`, więc czytnik nie czyta znaku drugi raz.
    const haczyk = document.createElement('span');
    haczyk.className = 'dn-drzewo__haczyk';
    haczyk.setAttribute('aria-hidden', 'true');
    guzik.append(haczyk);

    if (pozycja.ikona !== undefined) {
      guzik.append(elementIkony(pozycja.ikona, { rozmiar: 14, klasa: 'dn-drzewo__ikona' }));
    }

    const napisy = document.createElement('span');
    napisy.className = 'dn-drzewo__napisy';

    const nazwa = document.createElement('span');
    nazwa.className = 'dn-drzewo__nazwa';
    wstawTrafienia(nazwa, pozycja.nazwa, frazaCzynna);
    napisy.append(nazwa);

    // Druga nazwa stoi w tym samym wierszu co pierwsza, nie pod nią: pod nazwą
    // jest miejsce opisu, a skrót nie jest opisem — jest tą samą nazwą krócej.
    // Nawiasy niesie arkusz (`::before`/`::after`), żeby wytłuszczenie cięło
    // sam napis, a nie znaki, których w danych nie było.
    if (pozycja.nazwaKrotka !== undefined && pozycja.nazwaKrotka !== '') {
      const skrot = document.createElement('span');
      skrot.className = 'dn-drzewo__skrot';
      wstawTrafienia(skrot, pozycja.nazwaKrotka, frazaCzynna);
      nazwa.append(document.createTextNode(' '), skrot);
    }

    if (pozycja.opis !== undefined && pozycja.opis !== '') {
      const opis = document.createElement('span');
      opis.className = 'dn-drzewo__opis';
      opis.textContent = pozycja.opis;
      napisy.append(opis);
    }

    guzik.append(napisy);
    return { guzik, haczyk };
  }

  function kreska(): HTMLElement {
    const linia = document.createElement('hr');
    linia.className = 'dn-drzewo__kreska';
    return linia;
  }

  function stopka(dane: StopkaMenu): HTMLElement {
    // Stopka nie ma klucza, bo nie jest wyborem — nie wraca do `naWybor`, tylko
    // woła swoje `wykonaj`. Wiersz dostaje więc klucz pusty i to jest prawda
    // o niej, a nie obejście typu.
    const { guzik } = wiersz({ ...dane, klucz: '' }, 0, 'menuitem');
    guzik.classList.add('dn-drzewo__pozycja--stopka');
    guzik.addEventListener('click', () => {
      zwinWszystko();
      dane.wykonaj();
    });
    return guzik;
  }

  function zdaniePustego(): HTMLElement {
    const zdanie = document.createElement('p');
    zdanie.className = 'dn-drzewo__pusty';
    zdanie.textContent = ZDANIE_PUSTEGO;
    return zdanie;
  }

  /** Fraza bez trafień nie kasuje wykazu ani nie blokuje menu. */
  function zdanieBrakuTrafien(wpisana: string): HTMLElement {
    const zdanie = document.createElement('p');
    zdanie.className = 'dn-drzewo__pusty';
    zdanie.textContent =
      `Fraza „${wpisana}" nie pasuje do żadnej pozycji tego wykazu. ` +
      'Wyczyść pole szukania, żeby zobaczyć całe drzewo.';
    return zdanie;
  }

  element.dataset['otwarte'] = 'nie';
  odrysuj();

  return {
    element,

    ustaw(nowa, noweDrzewo) {
      wartosc.textContent = nowa;
      // Czytnik dostaje jedno zdanie: czego nastawa dotyczy i co w niej stoi.
      // Sama wartość („Opus 5") nie niesie słowa „model" i bez tego byłaby
      // dla czytającego uchem napisem bez przynależności.
      uchwyt.setAttribute('aria-label', `${opcje.nastawa}: ${nowa}`);
      uchwyt.title = `${opcje.nastawa}: ${nowa}`;
      drzewo = noweDrzewo;
      fraza.hidden = bezUchwytu || liczbaLisci(noweDrzewo) < (opcje.progSzukania ?? PROG_SZUKANIA);
      // Zniknięcie pola kasuje frazę — ale tylko tam, gdzie pole w ogóle bywa.
      // W trybie bez uchwytu pole nie stoi nigdy, a fraza przychodzi z zewnątrz
      // i skasowanie jej tutaj wycierałoby Operatorowi to, co właśnie wpisał.
      if (fraza.hidden && !bezUchwytu) {
        fraza.value = '';
        frazaBiezaca = '';
      }
      odrysuj();
    },

    zwin() {
      document.removeEventListener('pointerdown', poKlikuPoza, true);
      rozwiniete.clear();
      otwarte = false;
      lista.hidden = true;
      uchwyt.setAttribute('aria-expanded', 'false');
      element.dataset['otwarte'] = 'nie';
      odrysuj();
    },

    rozwin() {
      ustawOtwarte(true);
    },

    ustawFraze(nowa) {
      frazaBiezaca = nowa;
      // Pole, gdy stoi na ekranie, ma pokazywać to, czym wykaz jest przycięty —
      // inaczej Operator widziałby wykaz przycięty frazą, której nigdzie nie widać.
      if (!fraza.hidden) fraza.value = nowa;
      odrysuj();
    },

    przesunWyroznienie(krok) {
      const wykaz = pozycjeWykazu();
      if (wykaz.length === 0) return;
      const teraz = wyrozniona === null ? -1 : wykaz.indexOf(wyrozniona);
      // Wykaz zawija się na obu końcach, tak samo jak wędrówka ogniskiem — dół
      // z ostatniej pozycji wraca na pierwszą, a nie zatrzymuje się bez słowa.
      const nastepna =
        teraz < 0
          ? krok > 0
            ? 0
            : wykaz.length - 1
          : (teraz + krok + wykaz.length) % wykaz.length;
      ustawWyrozniona(wykaz[nastepna] ?? null);
    },

    wybierzWyrozniona() {
      if (wyrozniona === null) return false;
      // Gałąź nie jest wyborem — Enter na niej otwiera poziom. Klik zamiast
      // powtórzenia rozwijania, żeby droga była dokładnie ta sama co myszą.
      if (wyrozniona.dataset['galaz'] !== undefined) {
        wyrozniona.click();
        return true;
      }
      const klucz = wyrozniona.dataset['klucz'];
      if (klucz === undefined || klucz === '') return false;
      wybierz(klucz);
      return true;
    },
  };
}
