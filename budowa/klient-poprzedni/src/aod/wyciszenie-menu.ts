import { NAZWY_WYCISZEN, WyciszenieCzasowe } from './progi-aod';
import { TrybObecnosci, type StanObecnosci } from './tryb-obecnosci';
import {
  KLASY_ROZPOZNAWANE,
  KlasaZdarzen,
  NAZWY_KLAS,
  ZDARZENIA_KLAS,
  ZakresKontekstu,
  kluczWyciszenia,
  zdanieWyciszenia,
} from './wyciszenie-aod';
import { zdanieGranicyWyciszenia } from './wyciszenie-braki-kontraktu';
import type { KontekstWyciszenia } from './wyciszenie-kontekst';

/**
 * Menu kebab (⋮) wyciszania — rozdz. 2.6 („Wyciszenie sugestii … Menu kebab
 * powierzchni interakcji"), 3.5 (pięć rodzajów wyciszenia) i 8.3 („grupowanie
 * logiczne akcji: menu kebab zbierające wyciszanie") opracowania
 * `docs/funkcje-globalne/always-on-display.md`.
 *
 * Jedno menu, dwa miejsca osadzenia: powierzchnia interakcji awatara (warstwa 1,
 * wyciszenie od ręki bez otwierania kolumny) i nagłówek powierzchni interakcji
 * (rozdz. 2.3). Dwie kopie tego samego menu czytają ten sam stan i rysują się
 * jego powiadomieniem, więc nigdy nie mówią dwóch różnych rzeczy.
 *
 * Menu niesie komplet rozdz. 3.5:
 *   • trzy czasy wyciszenia czasowego,
 *   • wyciszenie bieżącego modułu i bieżącej karty sesji (kontekstowe),
 *   • wyciszenie każdej z sześciu klas zdarzeń rozdz. 3.2,
 *   • tryb cichy,
 *   • podgląd wyciszeń czynnych — każde swoim zdaniem, z „do kiedy" — wraz ze
 *     zniesieniem jednym kliknięciem, osobno i wszystkich naraz.
 *
 * ŻADNA pozycja nie pyta o potwierdzenie i żadna nie jest wyszarzana — reguła
 * przyjęta w `cztery-stery.ts`. Pozycja, której nakładka dziś nie ma czym
 * wykonać (wyciszenie modułu, którego rdzeń nie wskazał), zostaje klikalna
 * i mówi, czego brakuje oraz po czyjej stronie.
 *
 * Menu jest osobnym wyzwalaczem wobec awatara: kliknięcie pojedyncze i podwójne
 * awatara zostają nietknięte, a znak `⋮` jest zwykłym przyciskiem — osiągalnym
 * klawiszem tabulacji, otwieranym `Enter` i `Spacja`, zamykanym `Esc`, z ruchem
 * strzałkami po pozycjach.
 */

export interface MenuWyciszenia {
  /** Element do wstawienia obok awatara albo w nagłówek powierzchni. */
  element: HTMLElement;
  /** Znak `⋮` — do nadania mu własnego umiejscowienia przez osadzającego. */
  znak: HTMLButtonElement;
  /** Przerysowuje menu ze stanu. */
  odswiez(teraz: number): void;
  otworz(): void;
  zamknij(): void;
  czyOtwarte(): boolean;
  rozlacz(): void;
}

export interface OpisMenuWyciszenia {
  /** Stan obecności wraz z wykazem wyciszeń. */
  stan: StanObecnosci;
  /** Zegar podawany z zewnątrz — reguła i menu są sprawdzalne bez zegara maszyny. */
  teraz: () => number;
  /**
   * Kontekst nakładki: bieżący moduł i bieżąca karta sesji.
   *
   * Pominięty znaczy „menu nie zna bieżącego bytu" — pozycje kontekstowe
   * zostają widoczne i mówią to wprost, zamiast zniknąć.
   */
  kontekst?: KontekstWyciszenia;
  /** Krótkie zameldowanie czynności poza menu, gdy osadzający ma gdzie je pokazać. */
  zamelduj?: (zdanie: string) => void;
  /** Wołane przy otwarciu menu — miejsce na odświeżenie kontekstu z rdzenia. */
  naOtwarcie?: () => void;
  /** Etykieta znaku dla technologii wspomagających. */
  etykietaZnaku?: string;
}

export function utworzMenuWyciszenia(opis: OpisMenuWyciszenia): MenuWyciszenia {
  const { stan, teraz } = opis;
  const wyciszenia = stan.wyciszenia;

  const element = document.createElement('div');
  element.className = 'ao-kebab ao-wyciszenie';

  const znak = document.createElement('button');
  znak.type = 'button';
  znak.className = 'dn-btn-ikona ao-kebab__znak';
  znak.textContent = '⋮';
  znak.setAttribute(
    'aria-label',
    opis.etykietaZnaku ?? 'Menu wyciszania Always On Display — wyciszenie od ręki',
  );
  znak.setAttribute('aria-haspopup', 'menu');
  znak.setAttribute('aria-expanded', 'false');

  const lista = document.createElement('div');
  lista.className = 'ao-kebab__lista';
  lista.hidden = true;
  lista.setAttribute('role', 'menu');
  lista.setAttribute('aria-label', 'Wyciszanie Always On Display');

  /** Stan wyciszenia opisany zdaniem — pierwsza rzecz w menu, zawsze prawdziwa. */
  const opisStanu = document.createElement('p');
  opisStanu.className = 'ao-kebab__stan';
  opisStanu.setAttribute('aria-live', 'polite');
  lista.append(opisStanu);

  /** Podgląd wyciszeń czynnych — każde swoim zdaniem, ze zniesieniem obok. */
  const czynne = document.createElement('div');
  czynne.className = 'ao-wyciszenie__czynne';
  lista.append(czynne);

  const zniesWszystkie = document.createElement('button');
  zniesWszystkie.type = 'button';
  zniesWszystkie.className = 'ao-kebab__pozycja';
  zniesWszystkie.setAttribute('role', 'menuitem');
  zniesWszystkie.textContent = 'Znieś wszystkie wyciszenia';
  zniesWszystkie.addEventListener('click', () => {
    wyciszenia.zniesWszystkie();
    zamelduj('Wszystkie wyciszenia zniesione. Sugestie ujawniają się zgodnie z progami.');
  });
  lista.append(zniesWszystkie);

  // --- wyciszenie czasowe -------------------------------------------------

  lista.append(naglowek('Wyciszenie czasowe'));

  for (const czas of Object.values(WyciszenieCzasowe)) {
    const pozycja = document.createElement('button');
    pozycja.type = 'button';
    pozycja.className = 'ao-kebab__pozycja';
    pozycja.setAttribute('role', 'menuitem');
    pozycja.textContent = NAZWY_WYCISZEN[czas];
    pozycja.title =
      'Sugestie gromadzą się w liście oczekujących, dymek się nie otwiera, plakietka pozostaje ' +
      'ukryta. Punkt decyzyjny wstrzymujący proces ujawnia się mimo tego — plakietką.';
    pozycja.addEventListener('click', () => {
      wyciszenia.wyciszCzasem(czas, teraz());
      const wlaczone = wyciszenia.czasowe(teraz());
      zamelduj(
        wlaczone === null
          ? 'Wyciszenie czasowe nie zostało zapisane — powtórz czynność.'
          : zdanieWyciszenia(wlaczone),
      );
    });
    lista.append(pozycja);
  }

  // --- wyciszenie kontekstowe ---------------------------------------------

  lista.append(naglowek('Wyciszenie kontekstowe'));

  const pozycjaModulu = document.createElement('button');
  pozycjaModulu.type = 'button';
  pozycjaModulu.className = 'ao-kebab__pozycja';
  pozycjaModulu.setAttribute('role', 'menuitem');
  pozycjaModulu.addEventListener('click', () => przelaczKontekst(ZakresKontekstu.Modul));
  lista.append(pozycjaModulu);

  const pozycjaKarty = document.createElement('button');
  pozycjaKarty.type = 'button';
  pozycjaKarty.className = 'ao-kebab__pozycja';
  pozycjaKarty.setAttribute('role', 'menuitem');
  pozycjaKarty.addEventListener('click', () => przelaczKontekst(ZakresKontekstu.KartaSesji));
  lista.append(pozycjaKarty);

  const opisKontekstu = document.createElement('p');
  opisKontekstu.className = 'ao-kebab__stan';
  opisKontekstu.textContent =
    'Sugestie wskazanego zakresu nie ujawniają się; pozostałe zachowują pełne działanie.';
  lista.append(opisKontekstu);

  // --- wyciszenie klasy zdarzeń -------------------------------------------

  lista.append(naglowek('Wyciszenie klasy zdarzeń'));

  const pozycjeKlas = new Map<KlasaZdarzen, HTMLButtonElement>();
  for (const klasa of Object.values(KlasaZdarzen)) {
    const pozycja = document.createElement('button');
    pozycja.type = 'button';
    pozycja.className = 'ao-kebab__pozycja';
    pozycja.setAttribute('role', 'menuitem');
    pozycja.title = KLASY_ROZPOZNAWANE.has(klasa)
      ? ZDARZENIA_KLAS[klasa]
      : `${ZDARZENIA_KLAS[klasa]} Nakładka nie widzi dziś sygnału tej klasy — kontrakt go nie ` +
        'niesie. Wyciszenie zapisze się i zadziała, gdy sygnał wejdzie.';
    pozycja.addEventListener('click', () => {
      wyciszenia.przelaczKlase(klasa);
      zamelduj(
        wyciszenia.czyKlasaWyciszona(klasa)
          ? zdanieWyciszenia({ rodzaj: 'klasa-zdarzen', klasa })
          : `Klasa „${NAZWY_KLAS[klasa]}" znów tworzy sugestie.`,
      );
    });
    pozycjeKlas.set(klasa, pozycja);
    lista.append(pozycja);
  }

  // --- tryb cichy ---------------------------------------------------------

  lista.append(naglowek('Tryb cichy'));

  const pozycjaCiszy = document.createElement('button');
  pozycjaCiszy.type = 'button';
  pozycjaCiszy.className = 'ao-kebab__pozycja';
  pozycjaCiszy.setAttribute('role', 'menuitem');
  pozycjaCiszy.title =
    'Cała platforma: funkcja nie ujawnia sugestii samoczynnie i nie stosuje syntezy mowy; lista ' +
    'oczekujących pozostaje dostępna po otwarciu powierzchni interakcji.';
  pozycjaCiszy.addEventListener('click', () => {
    stan.przelaczTrybCichy();
    zamelduj(
      stan.tryb() === TrybObecnosci.Cichy
        ? 'Tryb cichy: sugestie gromadzą się bez dymka, synteza mowy wyłączona. Punkt decyzyjny ' +
            'wstrzymujący proces ujawnia się nadal — plakietką.'
        : 'Tryb pełny: sugestie ujawniają się zgodnie z progami, synteza mowy czynna.',
    );
  });
  lista.append(pozycjaCiszy);

  // --- zdania zamykające --------------------------------------------------

  const wyjatek = document.createElement('p');
  wyjatek.className = 'ao-kebab__stan';
  wyjatek.textContent =
    'Wyjątek wagi krytycznej: punkt decyzyjny pętli wykonawczej wstrzymujący proces ujawnia się ' +
    'mimo każdego wyciszenia — plakietką, bez dymka; w trybie cichym bez syntezy mowy.';
  lista.append(wyjatek);

  const granica = document.createElement('p');
  granica.className = 'ao-kebab__stan ao-wyciszenie__granica';
  granica.textContent = zdanieGranicyWyciszenia();
  lista.append(granica);

  element.append(znak, lista);

  // --- czynności ----------------------------------------------------------

  function naglowek(tresc: string): HTMLElement {
    const element = document.createElement('p');
    element.className = 'ao-wyciszenie__naglowek';
    element.textContent = tresc;
    return element;
  }

  /** Zameldowanie czynności: w menu zawsze, poza menu — gdy osadzający ma gdzie. */
  function zamelduj(zdanie: string): void {
    opisStanu.textContent = zdanie;
    opis.zamelduj?.(zdanie);
    odswiezCzynne(teraz());
  }

  function przelaczKontekst(zakres: ZakresKontekstu): void {
    const byt =
      zakres === ZakresKontekstu.Modul
        ? (opis.kontekst?.biezacyModul() ?? null)
        : (opis.kontekst?.biezacaKarta() ?? null);

    if (byt === null) {
      // Pozycja nie jest wyszarzana i nie milczy: mówi, czego brakuje.
      zamelduj(
        zakres === ZakresKontekstu.Modul
          ? (opis.kontekst?.powodBrakuModulu() ??
              'Menu nie dostało kontekstu nakładki, więc nie zna bieżącego modułu. Wyciszenie ' +
                'czasowe i wyciszenie klasy zdarzeń działają bez niego.')
          : (opis.kontekst?.powodBrakuKarty() ??
              'Menu nie dostało kontekstu nakładki, więc nie zna bieżącej karty sesji.'),
      );
      return;
    }

    wyciszenia.przelaczKontekst(zakres, byt.id, byt.nazwa);
    zamelduj(
      wyciszenia.czyKontekstWyciszony(zakres, byt.id)
        ? zdanieWyciszenia({ rodzaj: 'kontekstowe', zakres, wartosc: byt.id, nazwa: byt.nazwa })
        : `Wyciszenie zniesione — sugestie zakresu „${byt.nazwa}" ujawniają się znów.`,
    );
  }

  /** Podgląd wyciszeń czynnych wraz ze zniesieniem każdego jednym kliknięciem. */
  function odswiezCzynne(chwila: number): void {
    const wykaz = wyciszenia.czynne(chwila);
    czynne.replaceChildren();

    if (wykaz.length === 0) {
      const puste = document.createElement('p');
      puste.className = 'ao-kebab__stan';
      puste.textContent = 'Wyciszeń czynnych: brak. Sugestie ujawniają się zgodnie z progami.';
      czynne.append(puste);
      return;
    }

    for (const wyciszenie of wykaz) {
      const wiersz = document.createElement('div');
      wiersz.className = 'ao-wyciszenie__wiersz';

      const zdanie = document.createElement('p');
      zdanie.className = 'ao-kebab__stan';
      zdanie.textContent = zdanieWyciszenia(wyciszenie);

      const znies = document.createElement('button');
      znies.type = 'button';
      znies.className = 'ao-kebab__pozycja ao-wyciszenie__znies';
      znies.setAttribute('role', 'menuitem');
      znies.textContent = 'Znieś';
      const klucz = kluczWyciszenia(wyciszenie);
      znies.addEventListener('click', () => {
        wyciszenia.znies(klucz);
        zamelduj('Wyciszenie zniesione jednym kliknięciem — bez pytania o potwierdzenie.');
      });

      wiersz.append(zdanie, znies);
      czynne.append(wiersz);
    }
  }

  function odswiez(chwila: number): void {
    const wykaz = wyciszenia.czynne(chwila);
    odswiezCzynne(chwila);

    if (opisStanu.textContent === '' || opisStanu.textContent === null) {
      opisStanu.textContent =
        wykaz.length === 0
          ? 'Wyciszenie: brak. Sugestie ujawniają się zgodnie z progami opracowania.'
          : `Wyciszeń czynnych: ${wykaz.length}. Każde znosi się jednym kliknięciem.`;
    }

    const modul = opis.kontekst?.biezacyModul() ?? null;
    pozycjaModulu.textContent =
      modul === null
        ? 'Wycisz bieżący moduł'
        : wyciszenia.czyKontekstWyciszony(ZakresKontekstu.Modul, modul.id)
          ? `● Znieś wyciszenie modułu „${modul.nazwa}"`
          : `Wycisz bieżący moduł „${modul.nazwa}"`;
    pozycjaModulu.title =
      modul === null
        ? (opis.kontekst?.powodBrakuModulu() ??
          'Menu nie dostało kontekstu nakładki, więc nie zna bieżącego modułu.')
        : 'Sugestie tego modułu nie ujawniają się; pozostałe zachowują pełne działanie.';

    const karta = opis.kontekst?.biezacaKarta() ?? null;
    pozycjaKarty.textContent =
      karta === null
        ? 'Wycisz bieżącą kartę sesji'
        : wyciszenia.czyKontekstWyciszony(ZakresKontekstu.KartaSesji, karta.id)
          ? `● Znieś wyciszenie karty sesji „${karta.nazwa}"`
          : `Wycisz bieżącą kartę sesji „${karta.nazwa}"`;
    pozycjaKarty.title =
      karta === null
        ? (opis.kontekst?.powodBrakuKarty() ??
          'Menu nie dostało kontekstu nakładki, więc nie zna bieżącej karty sesji.')
        : 'Sugestie tej karty sesji nie ujawniają się; pozostałe zachowują pełne działanie.';

    for (const [klasa, pozycja] of pozycjeKlas) {
      const wyciszona = wyciszenia.czyKlasaWyciszona(klasa);
      // Stan nie idzie samą barwą: wyciszona klasa niesie znak wyboru i zmienia
      // czasownik pozycji na przeciwny.
      pozycja.dataset['wyciszone'] = String(wyciszona);
      pozycja.textContent = wyciszona
        ? `● Przywróć klasę „${NAZWY_KLAS[klasa]}"`
        : `Wycisz klasę „${NAZWY_KLAS[klasa]}"`;
    }

    const cichy = stan.tryb() === TrybObecnosci.Cichy;
    pozycjaCiszy.dataset['wyciszone'] = String(cichy);
    pozycjaCiszy.textContent = cichy ? '● Wyjdź z trybu cichego' : 'Przejdź w tryb cichy';

    znak.dataset['wyciszone'] = String(wykaz.length > 0 || cichy);
    znak.title =
      wykaz.length === 0 && !cichy
        ? 'Wyciszanie Always On Display — trzy czasy, moduł, karta sesji, klasa zdarzeń, tryb cichy.'
        : wykaz.map(zdanieWyciszenia).join(' ') + (cichy ? ' Tryb cichy czynny.' : '');
  }

  function otworz(): void {
    if (!lista.hidden) return;
    lista.hidden = false;
    znak.setAttribute('aria-expanded', 'true');
    opis.naOtwarcie?.();
    odswiez(teraz());
  }

  function zamknij(): void {
    if (lista.hidden) return;
    lista.hidden = true;
    znak.setAttribute('aria-expanded', 'false');
  }

  znak.addEventListener('click', () => (lista.hidden ? otworz() : zamknij()));

  /** Kliknięcie poza menu je zamyka — menu progresywne znika, gdy nie jest potrzebne. */
  function naKlikniecieDokumentu(zdarzenie: MouseEvent): void {
    if (lista.hidden) return;
    const cel = zdarzenie.target;
    if (cel instanceof Node && element.contains(cel)) return;
    zamknij();
  }

  /**
   * Klawiatura menu: `Esc` zamyka i wraca na znak, strzałki chodzą po pozycjach.
   *
   * Menu jest dostępne z klawiatury bez skrótu własnego — znak stoi w kolejności
   * tabulacji obok awatara. Skrótu nie wymyślamy: załącznik A.1 opracowania
   * wymienia pięć skrótów funkcji i menu kebab nie jest jednym z nich.
   */
  function naKlawisz(zdarzenie: KeyboardEvent): void {
    if (lista.hidden) return;
    const cel = zdarzenie.target;
    if (!(cel instanceof Node) || !element.contains(cel)) return;

    if (zdarzenie.key === 'Escape') {
      zamknij();
      znak.focus();
      zdarzenie.preventDefault();
      return;
    }

    if (zdarzenie.key !== 'ArrowDown' && zdarzenie.key !== 'ArrowUp') return;

    const pozycje = [...lista.querySelectorAll<HTMLButtonElement>('button')];
    if (pozycje.length === 0) return;
    const biezaca = pozycje.findIndex((pozycja) => pozycja === document.activeElement);
    const krok = zdarzenie.key === 'ArrowDown' ? 1 : -1;
    const nastepna = pozycje[(biezaca + krok + pozycje.length) % pozycje.length];
    nastepna?.focus();
    zdarzenie.preventDefault();
  }

  document.addEventListener('click', naKlikniecieDokumentu);
  element.addEventListener('keydown', naKlawisz);

  const przestanObserwowac = stan.obserwuj(() => odswiez(teraz()));

  odswiez(teraz());

  return {
    element,
    znak,
    odswiez,
    otworz,
    zamknij,
    czyOtwarte: () => !lista.hidden,

    rozlacz() {
      document.removeEventListener('click', naKlikniecieDokumentu);
      element.removeEventListener('keydown', naKlawisz);
      przestanObserwowac();
      element.remove();
    },
  };
}
