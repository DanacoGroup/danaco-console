import './przeciaganie.css';

/**
 * Zmiana kolejności pozycji w pasie to jedna mechanika dla wszystkich pasów, oparta na przechwyconym wskaźniku i klawiaturze zamiast HTML5 Drag and Drop, trzymająca wiązanie nasłuchów, które trzeba odłączyć po użyciu.
 */
export interface WiezPrzeciagania {
  /** Odłącza nasłuchy i sprząta kreskę upuszczenia. */
  rozlacz(): void;
  /** Czy gest trwa w tej chwili. */
  wRuchu(): boolean;
}

/** Ustawienia wzorca przeciągania; pojemnik, pozycje oraz skutek przestawienia są wymagane do działania mechaniki. */
export interface OpcjePrzeciagania {
  /** Element obejmujący pozycje pasa — na nim wiszą nasłuchy. */
  pojemnik: HTMLElement;
  /** Pozycje pasa w kolejności wyświetlania; czytane przy każdym geście. */
  pozycje(): readonly HTMLElement[];
  /** Skutek przestawienia: pozycja z jednego miejsca staje na drugim, wołany raz na koniec gestu. */
  przyPrzestawieniu(z: number, na: number): void;
  /** Nazwa pozycji dla czytnika ekranu; domyślnie jej tekst. */
  nazwaPozycji?(pozycja: HTMLElement): string;
  /** Nazwa pasa w komunikatach — „Karty sesji". */
  nazwaPasa?: string;
}

/**
 * Ile pikseli musi przejechać wskaźnik, zanim naciśnięcie stanie się gestem.
 *
 * Bez tego progu każde kliknięcie karty byłoby przeciągnięciem o zero pikseli,
 * a wybór karty przestałby działać.
 */
const PROG_GESTU = 5;

export function zwiazPrzeciaganieKolejnosci(opcje: OpcjePrzeciagania): WiezPrzeciagania {
  const { pojemnik } = opcje;
  const nazwaPasa = opcje.nazwaPasa ?? 'pasa';
  const nazwa = opcje.nazwaPozycji ?? ((el: HTMLElement) => (el.textContent ?? '').trim());

  const kreska = document.createElement('div');
  kreska.className = 'dnp-przenoszenie__kreska';
  kreska.setAttribute('aria-hidden', 'true');
  kreska.hidden = true;

  // Zdanie dla czytnika ekranu: gest wzrokowy musi mieć odpowiednik słowny.
  const mowa = document.createElement('span');
  mowa.className = 'dn-sr-only';
  mowa.setAttribute('role', 'status');
  mowa.setAttribute('aria-live', 'polite');

  pojemnik.append(kreska, mowa);
  // Znacznik dla arkusza: kursor „chwytu" pojawia się wyłącznie na pasie,
  // który naprawdę prowadzi gest.
  pojemnik.dataset.przenoszenie = 'tak';

  let zrodlo: number | null = null;
  let chwytak: HTMLElement | null = null;
  let startX = 0;
  let ruszyl = false;
  let cel: number | null = null;

  /** Pozycja, w której leży dany węzeł — albo -1. */
  function miejsceWezla(wezel: EventTarget | null): number {
    if (!(wezel instanceof Node)) return -1;
    return opcje.pozycje().findIndex((pozycja) => pozycja === wezel || pozycja.contains(wezel));
  }

  function przyNacisnieciu(zdarzenie: PointerEvent): void {
    // Wyłącznie przycisk główny i poza czynnościami wewnątrz pozycji — te nie mają być przechwytywane.
    if (zdarzenie.button !== 0) return;
    if (zdarzenie.target instanceof Element && zdarzenie.target.closest('button') !== null) return;

    const miejsce = miejsceWezla(zdarzenie.target);
    if (miejsce < 0) return;

    zrodlo = miejsce;
    chwytak = opcje.pozycje()[miejsce] ?? null;
    startX = zdarzenie.clientX;
    ruszyl = false;
    cel = null;
  }

  function przyRuchu(zdarzenie: PointerEvent): void {
    if (zrodlo === null || chwytak === null) return;

    if (!ruszyl) {
      if (Math.abs(zdarzenie.clientX - startX) < PROG_GESTU) return;
      ruszyl = true;
      chwytak.dataset.przenoszona = 'tak';
      chwytak.setAttribute('aria-grabbed', 'true');
      // Przechwycenie wskaźnika utrzymuje gest, gdy kursor wyjedzie poza pas.
      if (typeof chwytak.setPointerCapture === 'function') {
        chwytak.setPointerCapture(zdarzenie.pointerId);
      }
      powiedz(`Przenoszenie: ${nazwa(chwytak)}. Escape anuluje.`);
    }

    const pozycje = opcje.pozycje();
    cel = miejsceUpuszczenia(srodkiPozycji(pozycje), zdarzenie.clientX);
    pokazKreske(pozycje, cel);
    zdarzenie.preventDefault();
  }

  function przyPuszczeniu(): void {
    if (zrodlo === null) {
      posprzataj();
      return;
    }
    const z = zrodlo;
    const na = cel;
    const bylRuch = ruszyl;
    posprzataj();

    if (!bylRuch || na === null) return;
    // Przestawienie „na to samo miejsce" nie jest zmianą i nie ma o niej mówić.
    const docelowe = na > z ? na - 1 : na;
    if (docelowe === z) return;
    opcje.przyPrzestawieniu(z, docelowe);
    powiedz(`Przestawiono na miejsce ${docelowe + 1} z ${opcje.pozycje().length}.`);
  }

  function przyKlawiszu(zdarzenie: KeyboardEvent): void {
    if (zdarzenie.key === 'Escape' && zrodlo !== null) {
      // Gest przerwany nie przestawia niczego.
      zdarzenie.preventDefault();
      posprzataj();
      powiedz('Przenoszenie anulowane.');
      return;
    }

    // Równoważnik klawiaturowy gestu idzie tą samą drogą, więc obie dają tę samą kolejność.
    if (!zdarzenie.ctrlKey) return;
    const krok = zdarzenie.key === 'ArrowRight' ? 1 : zdarzenie.key === 'ArrowLeft' ? -1 : 0;
    if (krok === 0) return;

    const pozycje = opcje.pozycje();
    const miejsce = miejsceWezla(zdarzenie.target);
    if (miejsce < 0) return;
    const docelowe = miejsce + krok;
    if (docelowe < 0 || docelowe >= pozycje.length) return;

    zdarzenie.preventDefault();
    // Zatrzymanie tu jest konieczne, bo wędrówka pasa czyta te same strzałki bez tego wyboru.
    zdarzenie.stopPropagation();
    opcje.przyPrzestawieniu(miejsce, docelowe);
    powiedz(
      `${nazwa(pozycje[miejsce] ?? pozycje[0]!)}: miejsce ${docelowe + 1} z ${pozycje.length}.`,
    );
    // Ognisko wraca na przestawioną pozycję, żeby dało się przesuwać dalej.
    opcje.pozycje()[docelowe]?.focus();
  }

  function pokazKreske(pozycje: readonly HTMLElement[], miejsce: number): void {
    const przed = pozycje[miejsce] ?? null;
    kreska.hidden = false;
    if (przed === null) pojemnik.append(kreska);
    else pojemnik.insertBefore(kreska, przed);
  }

  function posprzataj(): void {
    if (chwytak !== null) {
      delete chwytak.dataset.przenoszona;
      chwytak.removeAttribute('aria-grabbed');
    }
    kreska.hidden = true;
    zrodlo = null;
    chwytak = null;
    ruszyl = false;
    cel = null;
  }

  function powiedz(zdanie: string): void {
    mowa.textContent = `${nazwaPasa}: ${zdanie}`;
  }

  pojemnik.addEventListener('pointerdown', przyNacisnieciu);
  pojemnik.addEventListener('pointermove', przyRuchu);
  pojemnik.addEventListener('pointerup', przyPuszczeniu);
  pojemnik.addEventListener('pointercancel', posprzataj);
  // Klawiatura na etapie przechwytywania trafia tu wcześniej niż wędrówka pasa.
  pojemnik.addEventListener('keydown', przyKlawiszu, true);

  return {
    wRuchu: () => ruszyl,
    rozlacz() {
      pojemnik.removeEventListener('pointerdown', przyNacisnieciu);
      pojemnik.removeEventListener('pointermove', przyRuchu);
      pojemnik.removeEventListener('pointerup', przyPuszczeniu);
      pojemnik.removeEventListener('pointercancel', posprzataj);
      pojemnik.removeEventListener('keydown', przyKlawiszu, true);
      posprzataj();
      delete pojemnik.dataset.przenoszenie;
      kreska.remove();
      mowa.remove();
    },
  };
}

/** Środki poziome wszystkich pozycji pasa — podstawa rachuby miejsca, w którym nastąpi upuszczenie przeciąganej pozycji. */
function srodkiPozycji(pozycje: readonly HTMLElement[]): number[] {
  return pozycje.map((pozycja) => {
    const rama = pozycja.getBoundingClientRect();
    return rama.left + rama.width / 2;
  });
}

/**
 * Miejsce, przed którym stanie przenoszona pozycja.
 *
 * Wynik jest z zakresu 0…n, gdzie n znaczy „na koniec". Funkcja jest czysta,
 * żeby dało się ją sprawdzić bez rozkładu strony.
 */
export function miejsceUpuszczenia(srodki: readonly number[], x: number): number {
  const miejsce = srodki.findIndex((srodek) => x < srodek);
  return miejsce < 0 ? srodki.length : miejsce;
}

/**
 * Przestawienie w tablicy — jedna rachuba dla gestu, klawiatury i nakładki.
 *
 * Osobna funkcja, bo tę samą kolejność liczą trzy miejsca; trzy kopie
 * rozjechałyby się przy pierwszej poprawce.
 */
export function przestawWTablicy<T>(wpisy: readonly T[], z: number, na: number): T[] {
  const wynik = [...wpisy];
  if (z < 0 || z >= wynik.length) return wynik;
  const [wyjeta] = wynik.splice(z, 1);
  if (wyjeta === undefined) return wynik;
  wynik.splice(Math.max(0, Math.min(na, wynik.length)), 0, wyjeta);
  return wynik;
}
