import './uchwyt.css';

/**
 * Uchwyt ręcznego ustawiania szerokości — pionowa krecha między dwiema
 * kolumnami, rozsuwana wskaźnikiem albo klawiaturą.
 *
 * Szerokość ustawia się ręcznie, każdemu oknu wewnętrznemu osobno. Nie ma
 * automatycznego zwężania i nie ma chowania pod zakładkę — nic nie znika
 * z ekranu bez ruchu ręki.
 *
 * Obsługa klawiaturą jest obowiązkowa: uchwyt osiągalny wyłącznie wskaźnikiem
 * odcinałby od jedynego sposobu ustawienia szerokości. Stąd `role="separator"`,
 * `tabIndex = 0`, strzałki, Home/End i pełny komplet `aria-value*`.
 *
 * Uchwyt nie zna paneli ani rozmowy. Oddaje wołającemu jedną liczbę w pikselach
 * i nic o niej nie zakłada; przycięcie do minimów robi `szerokosci-gniazda.ts`.
 *
 * Uchwyt stoi po lewej stronie sterowanej kolumny — tak układa gniazdo
 * `kolumnyGniazda`: rozmowa, uchwyt, panele. Ruch w lewo poszerza kolumnę, ruch
 * w prawo ją zwęża; strzałki idą tą samą logiką, żeby ręka i klawiatura nie
 * mówiły dwóch różnych rzeczy.
 */
export interface UchwytSzerokosci {
  element: HTMLElement;
  /** Ustawia wartość pokazywaną przez uchwyt (bez wołania `naZmiane`). */
  ustawWartosc(px: number): void;
}

/** Zależności uchwytu. */
export interface OpcjeUchwytu {
  /** Etykieta dostępności — co ten uchwyt rozsuwa. */
  etykieta: string;
  wartosc: number;
  minimum: number;
  /** Górna granica liczona w chwili ruchu — scena zmienia szerokość. */
  maksimum(): number;
  naZmiane(px: number): void;
}

/**
 * Krok jednego naciśnięcia strzałki, w pikselach.
 *
 * Cztery jednostki skali 4 px. Krok jednopikselowy kazałby trzymać strzałkę
 * kilkaset razy, żeby przejść przez kolumnę; krok stupikselowy nie pozwalałby
 * dojść do wartości, którą wskaźnik trafia bez wysiłku.
 */
const KROK_PX = 16;

/** Uchwyt rozsuwający sąsiadującą kolumnę; wartość oddaje w pikselach. */
export function utworzUchwytSzerokosci(opcje: OpcjeUchwytu): UchwytSzerokosci {
  const element = document.createElement('div');
  element.className = 'dn-uchwyt';
  element.setAttribute('role', 'separator');
  element.setAttribute('aria-orientation', 'vertical');
  element.setAttribute('aria-label', opcje.etykieta);
  element.tabIndex = 0;

  let wartosc = opcje.wartosc;

  /** Wartość sprowadzona do granic bieżących; górna liczona w chwili pytania. */
  function wGranicach(px: number): number {
    if (!Number.isFinite(px)) return opcje.minimum;
    const gora = Math.max(opcje.minimum, opcje.maksimum());
    return Math.round(Math.min(gora, Math.max(opcje.minimum, px)));
  }

  /** Zapisuje wartość i opowiada o niej czytnikowi ekranu. */
  function pokaz(px: number): void {
    wartosc = px;
    element.setAttribute('aria-valuenow', String(px));
    element.setAttribute('aria-valuemin', String(opcje.minimum));
    element.setAttribute('aria-valuemax', String(Math.max(opcje.minimum, opcje.maksimum())));
  }

  /** Zmiana wychodząca na zewnątrz — wołana wyłącznie przy ruchu uchwytem. */
  function przestaw(px: number): void {
    const nowa = wGranicach(px);
    if (nowa === wartosc) {
      // Granica osiągnięta: wartość ta sama, ale `aria-valuemax` mogła się
      // zmienić razem ze sceną, więc opis i tak idzie na nowo.
      pokaz(nowa);
      return;
    }
    pokaz(nowa);
    opcje.naZmiane(nowa);
  }

  // --- WSKAŹNIK ------------------------------------------------------------
  // Przechwycenie wskaźnika na `pointerdown`, żeby ruch poza uchwytem nadal do
  // niego trafiał — tak samo jak w `moduly/design/plansza-kompozycji.ts`.
  let poczatek: { x: number; wartosc: number } | null = null;

  element.addEventListener('pointerdown', (zdarzenie) => {
    poczatek = { x: zdarzenie.clientX, wartosc };
    element.setPointerCapture(zdarzenie.pointerId);
    element.dataset['przeciagany'] = 'tak';
    // Bez tego przeciąganie zaznacza tekst obu kolumn zamiast rozsuwać.
    zdarzenie.preventDefault();
  });

  element.addEventListener('pointermove', (zdarzenie) => {
    if (poczatek === null) return;
    przestaw(poczatek.wartosc - (zdarzenie.clientX - poczatek.x));
  });

  function koniecRuchu(): void {
    poczatek = null;
    delete element.dataset['przeciagany'];
  }

  element.addEventListener('pointerup', koniecRuchu);
  element.addEventListener('pointercancel', koniecRuchu);

  // --- KLAWIATURA ----------------------------------------------------------
  element.addEventListener('keydown', (zdarzenie) => {
    switch (zdarzenie.key) {
      case 'ArrowLeft':
        przestaw(wartosc + KROK_PX);
        break;
      case 'ArrowRight':
        przestaw(wartosc - KROK_PX);
        break;
      case 'Home':
        przestaw(opcje.minimum);
        break;
      case 'End':
        przestaw(opcje.maksimum());
        break;
      default:
        return;
    }
    // Strzałka na uchwycie rozsuwa kolumnę, nie przewija toru okien.
    zdarzenie.preventDefault();
  });

  pokaz(wGranicach(opcje.wartosc));

  return {
    element,
    ustawWartosc(px: number): void {
      pokaz(wGranicach(px));
    },
  };
}
