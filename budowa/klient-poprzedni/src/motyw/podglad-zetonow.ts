// ============================================================================
// DANACO CONSOLE — PODGLĄD ŻETONÓW  (RUSZTOWANIE, POZA PAKIETEM PRODUKTU)
// ----------------------------------------------------------------------------
// Wykłada na jeden ekran wszystkie żetony motywu — barwy, typografię,
// przestrzeń, ruch — w obu motywach naraz, żeby zmianę wartości dało się
// zobaczyć, a nie tylko odczytać z pliku.
//
// Strona nie dowodzi, że produkt działa: nie ma tu rdzenia, kanału kontraktu,
// sesji ani okien — sam katalog wartości z `zetony.json`. Poprawny wygląd
// żetonu tutaj nie znaczy, że jakikolwiek ekran produktu go używa. Do pakietu
// plik nie wchodzi: `main.ts` go nie wciąga, `vite build` czyta wyłącznie
// `index.html`.
//
// Uruchomienie:  cd budowa/client && npm run dev
//                → http://localhost:5173/src/motyw/podglad-zetonow.html
// ============================================================================

import zetony from './zetony.json';
import { elementIkony } from '../ikony/ikony';
import { ZNACZENIA_STANOW, type ZnaczenieStanu } from './znaczenia-stanow';
import {
  motywObowiazujacy,
  przelaczMotyw,
  przywrocPreferencjeSystemu,
  uruchomMotyw,
  ustawMotyw,
  ZDARZENIE_MOTYWU,
} from './motyw';

type Slownik = Record<string, string>;

function element(id: string): HTMLElement | null {
  return document.getElementById(id);
}

/**
 * Kafel z próbką barwy i nazwą żetonu.
 * Próbkę maluje żeton — `var(nazwa)`, nie wartość przepisana z pliku. Podpis
 * pozostaje wartością oczekiwaną, więc rozjazd między warstwą CSS a zetony.json
 * widać gołym okiem: kafel pokazuje jedno, podpis drugie.
 */
function probka(nazwa: string, oczekiwana: string): HTMLElement {
  const kafel = document.createElement('div');
  kafel.className = 'probka';

  const pole = document.createElement('i');
  pole.style.background = `var(${nazwa})`;
  kafel.append(pole);

  const podpis = document.createElement('span');
  podpis.textContent = `${nazwa} · ${oczekiwana}`;
  kafel.append(podpis);

  return kafel;
}

function wypelnijProbkami(id: string, slowniki: Slownik[]): void {
  const cel = element(id);
  if (cel === null) return;
  cel.replaceChildren(
    ...slowniki.flatMap((slownik) =>
      Object.entries(slownik).map(([nazwa, wartosc]) => probka(nazwa, wartosc)),
    ),
  );
}

/** Barwy prymitywne — skale surowe, nie do użycia wprost w komponentach. */
function pokazPrymitywy(): void {
  const { szarosc, sygnal, stany } = zetony.prymitywy;
  wypelnijProbkami('prymitywy', [
    szarosc,
    sygnal,
    stany.zielen,
    stany.bursztyn,
    stany.czerwien,
  ]);
}

/** Rama kokpitu — stała w obu motywach; dowód, że nie przełącza się z motywem. */
function pokazRame(): void {
  const { _uwaga: _pominiete, ...zetonyRamy } = zetony['rama-kokpitu'];
  wypelnijProbkami('rama', [zetonyRamy]);
}

/** Gradienty — wyłącznie ilustracyjne, nigdy tło przycisku, karty ani sekcji. */
function pokazGradienty(): void {
  const { _zakres: _pominiete, ...zetonyGradientow } = zetony.gradienty;
  wypelnijProbkami('gradienty', [zetonyGradientow]);
}

/**
 * Żetony semantyczne aktywnego motywu. Wartości czytamy z dokumentu, a nie
 * z pliku — dowodzi, że przełączenie motywu zmienia wartości, nie reguły.
 */
function pokazSemantyczne(): void {
  const cel = element('semantyczne');
  if (cel === null) return;
  const styl = getComputedStyle(document.documentElement);
  const nazwy = Object.keys(zetony.semantyczne.jasny);
  cel.replaceChildren(
    ...nazwy.map((nazwa) => probka(nazwa, styl.getPropertyValue(nazwa).trim())),
  );
}

/**
 * Plakietka stanu — barwa zawsze z ikoną i etykietą, nigdy sama.
 *
 * Barwy bierze rodzina z `stany.css`; ikonę i napis — katalog znaczeń
 * (`znaczenia-stanow.ts`). Rodzina `neutralna` nie ma własnych żetonów: to
 * plakietka bazowa, a znaczenie niesie w niej wyłącznie znak.
 *
 * Sekcja odpowiada na pytanie, czy Operator odróżni te stany, nie widząc barw
 * — dlatego pokazuje stany, a nie samą paletę.
 */
function plakietkaStanu(stan: ZnaczenieStanu): HTMLElement {
  const znak = ZNACZENIA_STANOW[stan];
  const plakietka = document.createElement('span');
  plakietka.className = 'plakietka';

  if (znak.rodzina === 'neutralna') {
    plakietka.style.color = 'var(--dn-tekst-2)';
    plakietka.style.backgroundColor = 'var(--dn-powierzchnia-2)';
    plakietka.style.borderColor = 'var(--dn-obrys)';
  } else {
    plakietka.style.color = `var(--dn-${znak.rodzina}-tekst)`;
    plakietka.style.backgroundColor = `var(--dn-${znak.rodzina}-tlo)`;
    plakietka.style.borderColor = `var(--dn-${znak.rodzina}-obrys)`;
  }

  const napis = document.createElement('span');
  napis.textContent = znak.etykieta;
  plakietka.append(elementIkony(znak.ikona, { rozmiar: 14 }), napis);

  // Wskaźnika pracy trwającej (`.dn-spinner`) ta strona nie pokazuje: mieszka
  // w bibliotece komponentów, a rusztowanie wciąga wyłącznie `motyw.css`.
  // Odtworzenie go tutaj lokalnie byłoby drugim zapisem tego samego bytu.
  return plakietka;
}

/**
 * Wszystkie stany znaczące naraz — w tym dwa, których rdzeń jeszcze nie zna.
 *
 * „Wstrzymany" i „Do weryfikacji" stoją na tej samej barwie i różnią się
 * wyłącznie ikoną oraz napisem — sekcja stawia je obok siebie, żeby to było
 * widać bez dokładania piątej barwy.
 */
function pokazStany(): void {
  const cel = element('stany');
  if (cel === null) return;
  const stany = Object.keys(ZNACZENIA_STANOW) as ZnaczenieStanu[];
  cel.replaceChildren(...stany.map(plakietkaStanu));
}

function wiersz(tresc: string): HTMLElement {
  const pozycja = document.createElement('li');
  pozycja.textContent = tresc;
  return pozycja;
}

function pokazTypografie(): void {
  const cel = element('typografia');
  if (cel === null) return;
  const stopnie = Object.entries(zetony.typografia.stopnie as Slownik);
  cel.replaceChildren(
    ...stopnie.map(([nazwa, oczekiwana]) => {
      const pozycja = wiersz(`${oczekiwana} — Danaco Console · ${nazwa}`);
      // Stopień bierze się z żetonu, nie z liczby w pliku: gęstość przestronna
      // (data-gestosc) przestawia --dn-fs-base i wiersz musi za tym nadążyć.
      pozycja.style.fontSize = `var(${nazwa})`;
      pozycja.style.lineHeight = 'var(--dn-lh-ciasny)';
      return pozycja;
    }),
  );
}

function pokazOdstepy(): void {
  const cel = element('przestrzen');
  if (cel === null) return;
  cel.replaceChildren(
    ...Object.entries(zetony.przestrzen.odstepy as Slownik).map(([nazwa, oczekiwana]) => {
      const pozycja = wiersz(`${nazwa} · ${oczekiwana}`);
      const pasek = document.createElement('div');
      pasek.className = 'odstep';
      pasek.style.width = `var(${nazwa})`;
      pozycja.append(pasek);
      return pozycja;
    }),
  );
}

function pokazPromienieICienie(): void {
  const promienie = element('promienie');
  const cienie = element('cienie');
  const kafel = (etykieta: string): HTMLElement => {
    const pole = document.createElement('div');
    pole.className = 'przycisk';
    pole.textContent = etykieta;
    return pole;
  };

  promienie?.replaceChildren(
    ...Object.entries(zetony.przestrzen.promienie as Slownik).map(([nazwa, oczekiwana]) => {
      const pole = kafel(`${nazwa} · ${oczekiwana}`);
      pole.style.borderRadius = `var(${nazwa})`;
      return pole;
    }),
  );

  cienie?.replaceChildren(
    ...Object.keys(zetony.cienie.jasny).map((nazwa) => {
      const pole = kafel(nazwa);
      pole.style.boxShadow = `var(${nazwa})`;
      return pole;
    }),
  );
}

/** Wykaz wartości skalarnych grupy — wymiary, siatka, warstwy. */
function pokazWykaz(id: string, grupa: Record<string, unknown>): void {
  const cel = element(id);
  if (cel === null) return;
  cel.replaceChildren(
    ...Object.entries(grupa)
      .filter(([nazwa, wartosc]) => {
        const skalar = typeof wartosc === 'string' || typeof wartosc === 'number';
        return skalar && nazwa.startsWith('--dn-');
      })
      .map(([nazwa, wartosc]) => wiersz(`${nazwa} · ${String(wartosc)}`)),
  );
}

function pokazWymiaryISiatke(): void {
  pokazWykaz('wymiary', zetony.wymiary);
  pokazWykaz('siatka', zetony['siatka-i-lamanie']);
  pokazWykaz('warstwy', zetony.warstwy);
}

function odswiez(): void {
  pokazSemantyczne();
  const stan = element('stan-motywu');
  if (stan !== null) stan.textContent = `motyw: ${motywObowiazujacy()}`;
}

function podepnijPrzelacznik(): void {
  element('motyw-jasny')?.addEventListener('click', () => ustawMotyw('light'));
  element('motyw-ciemny')?.addEventListener('click', () => ustawMotyw('dark'));
  element('motyw-system')?.addEventListener('click', () => przywrocPreferencjeSystemu());
  element('motyw-przelacz')?.addEventListener('click', () => przelaczMotyw());
  document.addEventListener(ZDARZENIE_MOTYWU, odswiez);
}

uruchomMotyw();
pokazPrymitywy();
pokazRame();
pokazGradienty();
pokazStany();
pokazTypografie();
pokazOdstepy();
pokazPromienieICienie();
pokazWymiaryISiatke();
podepnijPrzelacznik();
odswiez();
