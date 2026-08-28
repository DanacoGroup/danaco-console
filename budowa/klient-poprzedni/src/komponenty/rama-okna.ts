/**
 * Rama okna operacyjnego to jedna obudowa dla okien wszystkich modułów; nie
 * buduje treści i nie zna komend, tylko dostaje gotowe elementy od widoku
 * okna. Fazy okna nie należą do ramy, a wygląd pochodzi wyłącznie
 * z biblioteki i żetonów motywu.
 */

/** Rola, jaką okno pełni w module — wiodące, pomocnicze, monitor, kreator albo zarządca; rozstrzyga układ w pasie. */
export type RolaOkna = 'wiodące' | 'pomocnicze' | 'monitor' | 'kreator' | 'zarządca';

/** Waga plakietki stanu w nagłówku okna — licznik obiegów, stan kolejki albo tura w biegu; rozstrzyga barwę plakietki. */
export type WagaZnacznika = 'neutralna' | 'sukces' | 'ostrzezenie' | 'blad';

/**
 * Odmiany plakietki rozpisane na warianty jawne, a nie składane z napisu.
 * Biblioteka nie ma wariantu neutralnego — neutralna to plakietka bazowa.
 */
const ODMIANA_ZNACZNIKA: Record<Exclude<WagaZnacznika, 'neutralna'>, string> = {
  sukces: 'dn-plakietka--sukces',
  ostrzezenie: 'dn-plakietka--ostrzezenie',
  blad: 'dn-plakietka--blad',
};

export interface OpisRamyOkna {
  /** Nazwa okna operacyjnego — tytuł w nagłówku. */
  tytul: string;
  rola: RolaOkna;
  /** Kod okna z katalogu rdzenia (`okno_operacyjne.kod`) — trafia w `data-okno`. */
  kod?: string;
  /** Zdanie „po co to okno" pod nagłówkiem. Pominięte = brak akapitu. */
  przeznaczenie?: string;
  /** Nazwa modułu w etykiecie dostępności; bez niej etykieta niesie rolę. */
  modul?: string;
  /** Elementy doklejane do wiersza tytułu — na przykład dymek objaśnienia [?]. */
  dodatkiNaglowka?: readonly HTMLElement[];
  /** Przedrostek klas modułu (`'dt'`, `'mp'`…) — bywa uchwytem reguł własnych arkusza modułu. */
  przedrostek?: string;
  /** Okno przyjmuje ognisko z kodu (`element.focus()`) — pas z `tabIndex = -1`. */
  ogniskowalne?: boolean;
}

export interface RamaOkna {
  /** Sekcja osadzana w pasie układu modułu. */
  element: HTMLElement;
  /** Miejsce na sterowanie w wierszu tytułu. */
  pasek: HTMLElement;
  /** Panel akcji okna — czynności modułu wykonywane na tym oknie. */
  akcje: HTMLElement;
  /** Pasek narzędzi kontekstowych — przełączniki i pola zawężające. */
  narzedzia: HTMLElement;
  /** Miejsce treści okna wraz z jego stanem. */
  cialo: HTMLElement;
  /** Plakietka stanu w nagłówku. Pusty napis chowa plakietkę. */
  ustawZnacznik(tekst: string, waga?: WagaZnacznika): void;
}

/**
 * Gniazdo ramy: klasa biblioteczna zawsze, klasa modułu tylko gdy podany
 * przedrostek. Jedno miejsce składania nazw — inaczej przedrostek trzeba by
 * doklejać przy każdym gnieździe z osobna i przy którymś się zapomni.
 */
function gniazdo<K extends keyof HTMLElementTagNameMap>(
  znacznik: K,
  klasy: string,
  nazwa: string,
  przedrostek: string,
): HTMLElementTagNameMap[K] {
  const element = document.createElement(znacznik);
  element.className = klasy;
  if (przedrostek !== '') element.classList.add(`${przedrostek}-${nazwa}`);
  return element;
}

/** Nagłówek okna: tytuł, plakietka roli, dodatki przy tytule, znacznik stanu oraz pasek sterowania kontekstowego. */
function utworzNaglowek(
  opis: OpisRamyOkna,
  przedrostek: string,
): { element: HTMLElement; pasek: HTMLElement; znacznik: HTMLElement } {
  const tytul = gniazdo('h3', 'dn-karta-tytul dn-okno__tytul', 'okno__tytul', przedrostek);
  tytul.textContent = opis.tytul;

  const rola = document.createElement('span');
  rola.className = 'dn-plakietka dn-plakietka--rola';
  rola.textContent = opis.rola;

  const znacznik = gniazdo('span', 'dn-plakietka dn-okno__znacznik', 'okno__znacznik', przedrostek);
  znacznik.hidden = true;

  const pasek = gniazdo('div', 'dn-okno__pasek', 'okno__pasek', przedrostek);

  const element = gniazdo('header', 'dn-karta-naglowek', 'okno__naglowek', przedrostek);
  element.append(tytul, rola, ...(opis.dodatkiNaglowka ?? []), znacznik, pasek);
  return { element, pasek, znacznik };
}

/**
 * Pas kontrolek pod nagłówkiem — panel akcji albo narzędzia kontekstowe.
 * Klasa idzie napisem stałym, nie składana z członów: kontrola pokrycia klas
 * CSS czyta napisy, a nazwa złożona w czasie działania jest dla niej ślepa.
 */
function utworzPas(
  klasa: string,
  nazwa: string,
  etykieta: string,
  przedrostek: string,
): HTMLElement {
  const element = gniazdo('div', klasa, nazwa, przedrostek);
  element.setAttribute('aria-label', etykieta);
  return element;
}

export function utworzRameOkna(opis: OpisRamyOkna): RamaOkna {
  const przedrostek = opis.przedrostek ?? '';
  const naglowek = utworzNaglowek(opis, przedrostek);
  const akcje = utworzPas(
    'dn-okno__akcje',
    'okno__akcje',
    `Panel akcji okna ${opis.tytul}`,
    przedrostek,
  );
  const narzedzia = utworzPas(
    'dn-okno__narzedzia',
    'okno__narzedzia',
    `Narzędzia kontekstowe okna ${opis.tytul}`,
    przedrostek,
  );
  const cialo = gniazdo('div', 'dn-karta-cialo dn-okno__cialo', 'okno__cialo', przedrostek);

  const element = gniazdo('section', 'dn-karta dn-okno', 'okno', przedrostek);
  element.dataset['rola'] = opis.rola;
  if (opis.kod !== undefined) element.dataset['okno'] = opis.kod;
  if (opis.ogniskowalne === true) element.tabIndex = -1;
  element.setAttribute(
    'aria-label',
    opis.modul !== undefined
      ? `${opis.tytul} — okno modułu ${opis.modul}`
      : `${opis.tytul} — okno ${opis.rola}`,
  );
  const opisOkna = akapitPrzeznaczenia(opis, przedrostek);
  element.append(naglowek.element, ...opisOkna, akcje, narzedzia, cialo);

  return {
    element,
    pasek: naglowek.pasek,
    akcje,
    narzedzia,
    cialo,
    ustawZnacznik: (tekst, waga = 'neutralna') => ustawZnacznik(naglowek.znacznik, tekst, waga),
  };
}

/** Akapit „po co to okno" pod nagłówkiem — pusty wykaz elementów, gdy moduł opisu przeznaczenia nie podał. */
function akapitPrzeznaczenia(opis: OpisRamyOkna, przedrostek: string): HTMLElement[] {
  if (opis.przeznaczenie === undefined || opis.przeznaczenie === '') return [];
  const element = gniazdo('p', 'dn-pole-opis dn-okno__opis', 'okno__opis', przedrostek);
  element.textContent = opis.przeznaczenie;
  return [element];
}

function ustawZnacznik(znacznik: HTMLElement, tekst: string, waga: WagaZnacznika): void {
  znacznik.classList.remove(
    'dn-plakietka--sukces',
    'dn-plakietka--ostrzezenie',
    'dn-plakietka--blad',
  );
  if (waga !== 'neutralna') znacznik.classList.add(ODMIANA_ZNACZNIKA[waga]);
  znacznik.textContent = tekst;
  znacznik.hidden = tekst === '';
}
