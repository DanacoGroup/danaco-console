import { elementIkony, type NazwaIkony } from '../ikony/ikony';

// Jeden mechanizm rozwijanego drzewa dla całej aplikacji, wspólny dla paska, nagłówków i paneli.

/** Wspólne pola każdej pozycji drzewa menu: klucz oddawany przy wyborze, nazwa widoczna, opcjonalna nazwa krótka, opis i ikona. */
interface PozycjaWspolna {
  /** Klucz oddawany wołającemu przy wyborze; niepowtarzalny w obrębie drzewa. */
  klucz: string;
  /** Nazwa widoczna w menu. */
  nazwa: string;
  /** Druga, krótsza nazwa pozycji bez źródła; trafienie w nią liczy się wyżej niż w nazwę pełną. */
  nazwaKrotka?: string;
  /** Zdanie mówiące, co pozycja robi; stoi wprost pod nazwą, nie w dymku wywoływanym najechaniem. */
  opis?: string;
  /** Ikona przy nazwie; pominięta znaczy „bez ikony". */
  ikona?: NazwaIkony;
}

/** Gałąź menu — pozycja otwierająca kolejny poziom drzewa; niesie własne dzieci i rozwija się niezależnie od innych gałęzi. */
export interface GalazMenu extends PozycjaWspolna {
  rodzaj: 'galaz';
  dzieci: readonly PozycjaMenu[];
}

/** Liść wyboru jednokrotnego w menu — dokładnie jeden w całym drzewie ma pole wybrany ustawione na wartość prawda. */
export interface LiscMenu extends PozycjaWspolna {
  rodzaj: 'wybor';
  wybrany: boolean;
}

/** Nastawa dwustanowa — przełącznik siedzący wewnątrz drzewa menu, nie obok niego jako osobna kontrolka interfejsu. */
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

/** Stopka menu — droga do rejestru pozycji; podłączenie stopki do menu nie jest tym samym co rejestracja nowej pozycji w wykazie. */
export interface StopkaMenu {
  nazwa: string;
  opis?: string;
  ikona?: NazwaIkony;
  wykonaj(): void;
}

export interface OpcjeMenuDrzewa {
  /** Nazwa rodzajowa nastawy (Model, Wysiłek); nie trafia na ekran, idzie do etykiety dostępności. */
  nastawa: string;
  /** Ikona stojąca na uchwycie przed wartością; pominięta znaczy „bez ikony". */
  ikona?: NazwaIkony;
  /** Wołane po wyborze liścia albo przełącznika; klucz pozycji. */
  naWybor(klucz: string): void;
  /** Droga do rejestru na dole menu; pominięta znaczy „to menu jej nie ma". */
  stopka?: StopkaMenu;
  /** Od ilu liści menu stawia pole szukania; próg liczy liście, nie wszystkie pozycje wykazu. */
  progSzukania?: number;
  /** Kierunek rozwinięcia wykazu; nie jest preferencją wizualną, zależy od miejsca steru na ekranie. */
  kierunek?: 'dol' | 'gora';
  /** Menu bez własnego uchwytu, rozwijane i filtrowane z zewnątrz — tryb dla wykazu komend po ukośniku. */
  bezUchwytu?: boolean;
  /** Opis rysuje się tylko przy pozycji wyróżnionej zamiast przy każdej — dla wykazu na setki pozycji. */
  opisTylkoPrzyWyroznionej?: boolean;
}

export interface MenuDrzewo {
  /** Element montowany w pasku albo w nagłówku okna. */
  element: HTMLElement;
  /** Podaje mechanizmowi całą zawartość naraz: wartość uchwytu i drzewo; gałęzie przeżywają odświeżenie. */
  ustaw(wartosc: string, drzewo: readonly PozycjaMenu[]): void;
  /** Zwija menu wraz z gałęziami i zdejmuje nasłuchy dokumentu. Obowiązkowe. */
  zwin(): void;
  /** Rozwija wykaz bez klikania uchwytu — jedyna droga trybu bezUchwytu; ognisko zostaje u wołającego. */
  rozwin(): void;
  /** Podaje frazę filtrującą z zewnątrz; wchodzi do pola wewnętrznego, gdy pole stoi na ekranie. */
  ustawFraze(fraza: string): void;
  /** Przesuwa wyróżnienie o krok pozycji, nie ruszając ogniska; przy ukośniku ognisko zostaje w polu. */
  przesunWyroznienie(krok: number): void;
  /** Wybiera pozycję wyróżnioną; oddaje fałsz, gdy nie ma czego wybrać. Gałąź wyróżniona nie jest wyborem */
  wybierzWyrozniona(): boolean;
}

/** Zdanie o pustym wykazie drzewa — menu ogłasza brak pozycji zdaniem, zamiast po prostu przestać reagować na otwarcie. */
const ZDANIE_PUSTEGO = 'Wykaz jest pusty — nie ma tu jeszcze ani jednej pozycji do wyboru.';

/** Od ilu liści drzewo menu dostaje własne pole szukania, gdy wołający nie poda progu szukania we własnych opcjach. */
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
/** Ocena oznaczająca, że pozycja nie pasuje do szukanej frazy; każda wartość od zera w górę liczy się jako trafienie. */
const OCENA_BRAK = -1;

/** Ocena dopasowania jednego pola tekstowego pozycji; parametr szukane przychodzi już małymi literami i przycięty z brzegów. */
function ocenaPola(tekst: string | undefined, baza: number, szukane: string): number {
  if (tekst === undefined || tekst === '') return OCENA_BRAK;
  const gdzie = tekst.toLocaleLowerCase('pl-PL').indexOf(szukane);
  if (gdzie < 0) return OCENA_BRAK;
  return baza - Math.min(gdzie, OCENA_ODSUNIECIE_MAX);
}

/** Najlepsze trafienie pojedynczej pozycji spośród jej pól: nazwa skrócona bije pełną nazwę, a pełna nazwa bije opis. */
function ocenaPozycji(pozycja: PozycjaWspolna, szukane: string): number {
  return Math.max(
    ocenaPola(pozycja.nazwaKrotka, OCENA_SKROT, szukane),
    ocenaPola(pozycja.nazwa, OCENA_NAZWA, szukane),
    ocenaPola(pozycja.opis, OCENA_OPIS, szukane),
  );
}

/**
 * Najlepsze trafienie w całym poddrzewie — ocena gałęzi bierze się z dzieci.
 * Gałąź z trafieniem idealnym ma stanąć wyżej niż gałąź z trafieniem
 * bylejakim, inaczej porządek trafności kończyłby się na pierwszym poziomie.
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
 * Wstawia napis do elementu, wytłuszczając każde wystąpienie frazy, nie
 * tylko pierwsze: Operator widzi wtedy liczbę trafień w nazwie. Bez frazy
 * wchodzi jeden węzeł tekstowy.
 */
function wstawTrafienia(cel: HTMLElement, napis: string, szukane: string): void {
  if (szukane === '') {
    cel.textContent = napis;
    return;
  }
  const male = napis.toLocaleLowerCase('pl-PL');
  // Zmiana wielkości liter bywa zmianą długości; wtedy napis zostaje w całości, bez wytłuszczenia.
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
  // Trafienia nie ma, gdy pozycja weszła do wykazu przez inne pole — nazwa zostaje zwykłym napisem.
  if (czesci.length === 0) {
    cel.textContent = napis;
    return;
  }
  if (od < napis.length) czesci.push(document.createTextNode(napis.slice(od)));
  cel.replaceChildren(...czesci);
}

/** Licznik egzemplarzy mechanizmu menu drzewa — daje każdemu wystąpieniu niepowtarzalny przedrostek identyfikatora pozycji. */
let licznikMenu = 0;

export function utworzMenuDrzewo(opcje: OpcjeMenuDrzewa): MenuDrzewo {
  let otwarte = false;
  /** Klucze gałęzi rozwiniętych — przeżywają przerysowanie drzewa. */
  const rozwiniete = new Set<string>();
  let drzewo: readonly PozycjaMenu[] = [];
  const bezUchwytu = opcje.bezUchwytu === true;
  /** Jedyne źródło prawdy o frazie — wewnętrzne pole tylko je odzwierciedla przy każdym uderzeniu. */
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

  // Na uchwycie stoi wartość, nie nazwa nastawy — tym ster różni się od menu czynności.
  const wartosc = document.createElement('span');
  wartosc.className = 'dn-drzewo__wartosc';
  uchwyt.append(wartosc);
  uchwyt.append(elementIkony('grot-dol', { rozmiar: 12, klasa: 'dn-drzewo__grot' }));

  const lista = document.createElement('div');
  lista.className = 'dn-drzewo__lista';
  lista.setAttribute('role', 'menu');
  lista.setAttribute('aria-label', opcje.nastawa);
  lista.hidden = true;

  // Pole szukania stoi poza rolą menu, nad wykazem — strzałki mają chodzić po pozycjach, nie po polu.
  const fraza = document.createElement('input');
  fraza.type = 'search';
  fraza.className = 'dn-drzewo__szukanie';
  fraza.placeholder = 'Szukaj…';
  fraza.setAttribute('aria-label', `Szukanie w wykazie: ${opcje.nastawa}`);
  fraza.hidden = true;
  // Strzałka w dół z pola przenosi ognisko na pierwsze trafienie, bez sięgania po mysz.
  fraza.addEventListener('input', () => {
    frazaBiezaca = fraza.value;
    odrysuj();
  });
  fraza.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'ArrowDown' && zdarzenie.key !== 'Enter') return;
    zdarzenie.preventDefault();
    wedrowne()[0]?.focus();
  });

  // W trybie bez uchwytu ani uchwyt, ani pole szukania nie wchodzą do dokumentu; wyzwala je ukośnik.
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

  // Wędrówka ognisk i sterowanie poziomami stoją na całym menu — ognisko bywa też na uchwycie.
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
    // Prawo otwiera gałąź, lewo ją zamyka — dwa klawisze na dwa kierunki, tak samo na każdej głębokości.
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

  // Ognisko i wyróżnienie wskazują tę samą pozycję; z uchwytem strzałki przenoszą ognisko naprawdę.
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

  /** Pozycje wykazu w porządku rysowania, bez pytania o offsetParent; stopka odpada, bo nie jest wyborem. */
  function pozycjeWykazu(): HTMLElement[] {
    const wybor =
      '[role="menuitem"]:not(.dn-drzewo__pozycja--stopka), ' +
      '[role="menuitemcheckbox"], [role="menuitemradio"]';
    return [...lista.querySelectorAll<HTMLElement>(wybor)];
  }

  /** Przenosi wyróżnienie; null znaczy brak czego wyróżnić. Wyróżnienie jest wirtualne, nie rusza ognisk */
  function ustawWyrozniona(cel: HTMLElement | null): void {
    if (wyrozniona !== null && wyrozniona !== cel) delete wyrozniona.dataset['wyrozniona'];
    wyrozniona = cel;
    if (cel === null) {
      lista.removeAttribute('aria-activedescendant');
    } else {
      // Identyfikator nadaje się dopiero tutaj i tylko raz: wskazuje na niego etykieta dostępności listy.
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

  /** Drzewo przycięte frazą — gałąź zostaje, gdy cokolwiek pod nią pasuje; trafiona rozwija się sama. */
  function przytnij(pozycje: readonly PozycjaMenu[], szukane: string): PozycjaMenu[] {
    const wynik: { pozycja: PozycjaMenu; ocena: number }[] = [];
    for (const p of pozycje) {
      if (p.rodzaj === 'galaz' || p.rodzaj === 'grupa') {
        const dzieci = przytnij(p.dzieci, szukane);
        const sama = p.rodzaj === 'galaz' ? ocenaPozycji(p, szukane) : OCENA_BRAK;
        if (dzieci.length === 0 && sama < 0) continue;
        if (p.rodzaj === 'galaz') rozwiniete.add(p.klucz);
        // Gałąź trafiona własną nazwą zachowuje całe poddrzewo: Operator szuka gałęzi, nie liścia pod nią.
        const zawartosc = dzieci.length === 0 ? p.dzieci : dzieci;
        const ocena = Math.max(sama, najlepszaOcena(zawartosc, szukane));
        wynik.push({ pozycja: { ...p, dzieci: zawartosc }, ocena });
        continue;
      }
      const ocena = ocenaPozycji(p, szukane);
      if (ocena >= 0) wynik.push({ pozycja: p, ocena });
    }
    // Porządek wedle trafności obowiązuje tylko przy frazie; sortowanie w JavaScripcie jest stabilne.
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

    // Pierwsza pozycja jest wyróżniona od razu, żeby Enter działał natychmiast przy jednym trafieniu.
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
      // Ślad wyboru na gałęzi — menu pokazuje wybór, zamiast kazać go szukać wejściem w każdą gałąź.
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
    // Znak wyboru stoi w każdym liściu, nie tylko wybranym — miejsce zajęte na stałe nie przesuwa wierszy.
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

  /** Szkielet wiersza: haczyk, ikona, nazwa i opis — jeden rytm na wszystkie rodzaje pozycji. */
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

    // Miejsce na haczyk jest stałe, także gdy go nie widać — inaczej wiersze przeskakiwałyby w bok.
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

    // Druga nazwa stoi w tym samym wierszu co pierwsza, nie pod nią — pod nazwą jest miejsce opisu.
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
    // Stopka nie ma klucza, bo nie jest wyborem — nie wraca do naWybor, tylko woła własne wykonaj.
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
      // Czytnik dostaje jedno zdanie: czego nastawa dotyczy i co w niej stoi, nie samą gołą wartość.
      uchwyt.setAttribute('aria-label', `${opcje.nastawa}: ${nowa}`);
      uchwyt.title = `${opcje.nastawa}: ${nowa}`;
      drzewo = noweDrzewo;
      fraza.hidden = bezUchwytu || liczbaLisci(noweDrzewo) < (opcje.progSzukania ?? PROG_SZUKANIA);
      // Zniknięcie pola kasuje frazę, ale tylko tam, gdzie pole w ogóle bywa — nie w trybie bez uchwytu.
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
      // Pole, gdy stoi na ekranie, ma pokazywać to, czym wykaz jest przycięty.
      if (!fraza.hidden) fraza.value = nowa;
      odrysuj();
    },

    przesunWyroznienie(krok) {
      const wykaz = pozycjeWykazu();
      if (wykaz.length === 0) return;
      const teraz = wyrozniona === null ? -1 : wykaz.indexOf(wyrozniona);
      // Wykaz zawija się na obu końcach, tak samo jak wędrówka ogniskiem.
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
      // Gałąź nie jest wyborem — Enter na niej otwiera poziom, tak samo jak klik.
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
