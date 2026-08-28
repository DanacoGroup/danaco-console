import { przycisk } from '../../modele/kontrolki-formularza';
import { utworzRozwiniecie } from './warstwy-translate';

/**
 * Karta oceny jakości językowej modelu MQM/DQF na warstwie czwartej QA & Review Center.
 */

/**
 * Kategoria błędu górnego poziomu metryki jakości MQM/DQF, przypisana ocenianemu fragmentowi przekładu.
 */
export interface KategoriaOceny {
  /** Nazwa kategorii. */
  readonly nazwa: string;
  /** Czego kategoria dotyczy. */
  readonly opis: string;
}

export const KATEGORIE_OCENY: readonly KategoriaOceny[] = [
  { nazwa: 'Dokładność', opis: 'Zgodność znaczenia przekładu z tekstem źródłowym.' },
  { nazwa: 'Poprawność językowa', opis: 'Gramatyka, ortografia, interpunkcja języka docelowego.' },
  { nazwa: 'Terminologia', opis: 'Zgodność z bazą terminologiczną i spójność odpowiedników.' },
  { nazwa: 'Styl', opis: 'Zgodność rejestru i tonu z przeznaczeniem materiału.' },
  { nazwa: 'Konwencje lokalne', opis: 'Zapis liczb, dat, walut, adresów i typografii rynku.' },
  { nazwa: 'Zgodność z rzeczywistością', opis: 'Prawdziwość treści wobec realiów rynku docelowego.' },
  { nazwa: 'Postać i skład', opis: 'Znaczniki, symbole zastępcze, układ i długość treści.' },
];

/**
 * Dotkliwość błędu wraz z wagą metryki, jaką ta dotkliwość wnosi do liczbowego wyniku oceny tej karty.
 */
export interface Dotkliwosc {
  readonly nazwa: string;
  readonly waga: number;
}

export const DOTKLIWOSCI: readonly Dotkliwosc[] = [
  { nazwa: 'obojętna', waga: 0 },
  { nazwa: 'drobna', waga: 1 },
  { nazwa: 'poważna', waga: 5 },
  { nazwa: 'krytyczna', waga: 25 },
];

/**
 * Jedna liczba karty oceny: kategoria błędu, jego dotkliwość i liczba stwierdzonych błędów tego rodzaju.
 */
export interface WpisOceny {
  readonly kategoria: string;
  readonly dotkliwosc: Dotkliwosc;
  readonly liczba: number;
}

/**
 * Wynik przeliczenia karty oceny jakości językowej na liczbę punktów modelu MQM/DQF dla tego przekładu.
 */
export interface OcenaLqa {
  /** Suma ważona kar — liczba błędów przemnożona przez wagę dotkliwości. */
  readonly sumaWazona: number;
  /** Liczba stwierdzonych błędów bez ważenia. */
  readonly liczbaBledow: number;
  /** Liczba słów materiału, po której wzór dzieli sumę kar. */
  readonly liczbaSlow: number;
  /** Wynik metryki albo `null`, gdy materiał nie ma ani jednego słowa. */
  readonly wynik: number | null;
}

/**
 * Liczy słowa ocenianego tekstu jako odcinki rozdzielone białymi znakami w treści całego przekładu tekstu.
 */
export function liczSlowa(tekst: string): number {
  const przyciety = tekst.trim();
  if (przyciety === '') return 0;
  return przyciety.split(/\s+/u).length;
}

/**
 * Przelicza kartę wzorem metryki.
 *
 * Wynik ujemny schodzi do zera: metryka nie zna jakości gorszej niż zero,
 * a liczba ujemna sugerowałaby skalę, której nie ma.
 */
export function policzOcene(
  wpisy: readonly WpisOceny[],
  liczbaSlow: number,
): OcenaLqa {
  const sumaWazona = wpisy.reduce((suma, wpis) => suma + wpis.liczba * wpis.dotkliwosc.waga, 0);
  const liczbaBledow = wpisy.reduce((suma, wpis) => suma + wpis.liczba, 0);
  if (liczbaSlow <= 0) {
    return { sumaWazona, liczbaBledow, liczbaSlow: 0, wynik: null };
  }
  const wynik = 100 - (sumaWazona / liczbaSlow) * 100;
  return { sumaWazona, liczbaBledow, liczbaSlow, wynik: Math.max(0, Math.round(wynik * 100) / 100) };
}

/**
 * Zdanie o wyniku oceny: liczby punktów wraz z tym, czego rdzeń o samym tekście przekładu nie może wiedzieć.
 */
export function zdanieOceny(ocena: OcenaLqa): string {
  const podstawa =
    `Błędów stwierdzonych: ${String(ocena.liczbaBledow)}, suma ważona kar: ` +
    `${String(ocena.sumaWazona)}, słów materiału: ${String(ocena.liczbaSlow)}.`;
  if (ocena.wynik === null) {
    return (
      `${podstawa} Wyniku nie ma z czego policzyć: wzór metryki dzieli sumę kar przez liczbę ` +
      'słów, a tekst źródłowy jest pusty.'
    );
  }
  return (
    `${podstawa} Wynik metryki MQM/DQF: ${String(ocena.wynik)} na 100. ` +
    'Ocena powstaje i zostaje w oknie — kontrakt nie ma komendy jej zapisu ani wydania raportu.'
  );
}

/**
 * Karta oceny jakości językowej przekładu jako element warstwy czwartej QA & Review Center tej budowy.
 */
export interface KartaLqa {
  element: HTMLElement;
  /** Rozwija kartę z zewnątrz. */
  rozwin(): void;
}

export function utworzKarteLqa(liczbaSlow: () => number): KartaLqa {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 4,
    nazwa: 'Karta oceny LQA (MQM/DQF)',
    wyjasnienie:
      'Ocena jakości przekładu wedle kategorii błędu i wag dotkliwości metryki MQM/DQF.',
    znacznik: '☰',
  });

  const pola = new Map<string, HTMLInputElement>();
  const tabela = document.createElement('table');
  tabela.className = 'dn-tabela mt-lqa';
  tabela.append(naglowekTabeli(), cialoTabeli(pola));

  const zdanie = document.createElement('p');
  zdanie.className = 'mt-lqa__wynik';

  const raport = document.createElement('pre');
  raport.className = 'mt-lqa__raport';

  const przelicz = przycisk('Przelicz ocenę', 'dn-btn dn-btn--sm dn-btn--atrament');
  przelicz.addEventListener('click', () => {
    const wpisy = odczytajWpisy(pola);
    const ocena = policzOcene(wpisy, liczbaSlow());
    zdanie.textContent = zdanieOceny(ocena);
    raport.textContent = zlozRaport(wpisy, ocena);
  });

  const pasek = document.createElement('div');
  pasek.className = 'mt-pasek';
  pasek.append(przelicz);

  rozwiniecie.tresc.append(tabela, pasek, zdanie, raport);
  return { element: rozwiniecie.element, rozwin: rozwiniecie.rozwin };
}

function naglowekTabeli(): HTMLElement {
  const wiersz = document.createElement('tr');
  wiersz.append(komorka('th', 'Kategoria błędu'));
  for (const dotkliwosc of DOTKLIWOSCI) {
    wiersz.append(komorka('th', `${dotkliwosc.nazwa} (waga ${String(dotkliwosc.waga)})`));
  }
  const naglowek = document.createElement('thead');
  naglowek.append(wiersz);
  return naglowek;
}

function cialoTabeli(pola: Map<string, HTMLInputElement>): HTMLElement {
  const cialo = document.createElement('tbody');
  for (const kategoria of KATEGORIE_OCENY) {
    const wiersz = document.createElement('tr');
    const nazwa = komorka('td', kategoria.nazwa);
    nazwa.title = kategoria.opis;
    wiersz.append(nazwa);
    for (const dotkliwosc of DOTKLIWOSCI) {
      const komorkaLiczby = document.createElement('td');
      const pole = poleLiczby(`${kategoria.nazwa}, dotkliwość ${dotkliwosc.nazwa}`);
      pola.set(klucz(kategoria.nazwa, dotkliwosc.nazwa), pole);
      komorkaLiczby.append(pole);
      wiersz.append(komorkaLiczby);
    }
    cialo.append(wiersz);
  }
  return cialo;
}

function klucz(kategoria: string, dotkliwosc: string): string {
  return `${kategoria}|${dotkliwosc}`;
}

function komorka(rodzaj: 'th' | 'td', tresc: string): HTMLTableCellElement {
  const element = document.createElement(rodzaj);
  element.textContent = tresc;
  return element;
}

function poleLiczby(nazwa: string): HTMLInputElement {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'number';
  kontrolka.min = '0';
  kontrolka.step = '1';
  kontrolka.value = '0';
  kontrolka.className = 'dn-pole-kontrolka mt-lqa__liczba';
  kontrolka.setAttribute('aria-label', nazwa);
  return kontrolka;
}

function odczytajWpisy(pola: Map<string, HTMLInputElement>): readonly WpisOceny[] {
  const wpisy: WpisOceny[] = [];
  for (const kategoria of KATEGORIE_OCENY) {
    for (const dotkliwosc of DOTKLIWOSCI) {
      const pole = pola.get(klucz(kategoria.nazwa, dotkliwosc.nazwa));
      const liczba = Number(pole?.value ?? '0');
      // Wartość nieliczbowa albo ujemna schodzi do zera przy przeliczaniu pola karty.
      wpisy.push({
        kategoria: kategoria.nazwa,
        dotkliwosc,
        liczba: Number.isFinite(liczba) && liczba > 0 ? Math.trunc(liczba) : 0,
      });
    }
  }
  return wpisy;
}

/**
 * Raport oceny jakości jako tekst gotowy do przeniesienia, bo kontrakt nie ma osobnej komendy jego eksportu.
 */
function zlozRaport(wpisy: readonly WpisOceny[], ocena: OcenaLqa): string {
  const wiersze = wpisy
    .filter((wpis) => wpis.liczba > 0)
    .map(
      (wpis) =>
        `${wpis.kategoria} · dotkliwość ${wpis.dotkliwosc.nazwa} · błędów ${String(wpis.liczba)} ` +
        `· kara ${String(wpis.liczba * wpis.dotkliwosc.waga)}`,
    );
  if (wiersze.length === 0) {
    return 'Karta oceny LQA (MQM/DQF): nie wpisano ani jednego błędu, więc raport nie ma pozycji.';
  }
  return [
    'Karta oceny LQA (MQM/DQF)',
    ...wiersze,
    `Suma ważona kar: ${String(ocena.sumaWazona)}`,
    `Słów materiału: ${String(ocena.liczbaSlow)}`,
    ocena.wynik === null
      ? 'Wynik: nie policzono — materiał nie ma ani jednego słowa.'
      : `Wynik: ${String(ocena.wynik)} na 100`,
  ].join('\n');
}
