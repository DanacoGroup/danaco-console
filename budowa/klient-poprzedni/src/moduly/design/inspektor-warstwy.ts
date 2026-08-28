import type { DesignBoardLayer } from '../../../../shared/contract';
import { poleTekstowe } from '../../modele/kontrolki-formularza';
import { oznaczWarstwe } from './warstwy-designu';
import { BOK_WYJSCIOWY } from './zapis-kompozycji';
import type { StanKompozycji } from './stan-kompozycji';

/**
 * Inspektor właściwości warstwy zaznaczonej — jednocześnie edytor liczb dla projektanta i handoff,
 * wydający te same cztery liczby jako gotowy zapis reguł stylu do przepisania bez mierzenia na oko.
 */
export interface InspektorWarstwy {
  element: HTMLElement;
  odswiez(): void;
  /** Prowadzi ognisko do pola adnotacji — droga skrótu klawiszowego. */
  ogniskujAdnotacje(): boolean;
}

/** Pola geometrii warstwy zaznaczonej wraz z ich kluczem identyfikującym pole w kontrakcie danych rdzenia. */
const POLA_GEOMETRII = [
  ['x', 'Położenie poziome'],
  ['y', 'Położenie pionowe'],
  ['width', 'Szerokość'],
  ['height', 'Wysokość'],
] as const;

type KluczGeometrii = (typeof POLA_GEOMETRII)[number][0];

/** Właściwości bez pola w kontrakcie: warstwa kompozycji ich nie niesie, więc wykaz nazywa ten brak wprost. */
const BEZ_POLA_W_KONTRAKCIE: readonly (readonly [string, string])[] = [
  ['Wypełnienie', 'Warstwa nie ma pola barwy wypełnienia — kanwa rysuje ją oprawą modułu.'],
  ['Obrys', 'Ani grubości, ani barwy obrysu warstwa nie wyraża.'],
  ['Efekty', 'Cienia, poświaty i rozmycia warstwa nie niesie.'],
  [
    'Więzy responsywne',
    'Reguły kotwiczenia przy zmianie rozmiaru ramki nie ma czym zapisać przy warstwie.',
  ],
  [
    'Auto-layout',
    'Odstępu i wyściółki warstwa nie niesie, więc rozmieszczenie jest jednorazowe, ' +
      'a nie trwałą właściwością.',
  ],
];

export function utworzInspektorWarstwy(stan: StanKompozycji): InspektorWarstwy {
  const tytul = document.createElement('h4');
  tytul.className = 'md-panel__tytul';
  tytul.textContent = 'Inspektor właściwości';

  const zakres = document.createElement('p');
  zakres.className = 'dn-pole-opis md-panel__opis';

  const kontrolki = new Map<KluczGeometrii, HTMLInputElement>();
  const siatka = document.createElement('div');
  siatka.className = 'md-inspektor__siatka';

  for (const [klucz, etykieta] of POLA_GEOMETRII) {
    const pole = poleTekstowe({ etykieta, podpowiedz: '0' });
    pole.kontrolka.inputMode = 'numeric';
    pole.kontrolka.addEventListener('change', () => nanies());
    kontrolki.set(klucz, pole.kontrolka);
    siatka.append(pole.element);
  }

  const adnotacja = poleTekstowe({
    etykieta: 'Adnotacja warstwy',
    podpowiedz: 'co ta warstwa przedstawia',
    opis: 'Pole adnotacji warstwy — jedyne pole opisowe, które jedzie do rdzenia z układem.',
  });
  adnotacja.kontrolka.addEventListener('change', () => {
    const warstwa = wskazana();
    if (warstwa === null) return;
    stan.ustawAdnotacje(warstwa.id, adnotacja.kontrolka.value);
  });

  const zapisStylu = document.createElement('pre');
  zapisStylu.className = 'md-inspektor__styl';

  const oZapisie = document.createElement('p');
  oZapisie.className = 'dn-pole-opis md-inspektor__o-zapisie';
  oZapisie.textContent =
    'Zapis reguł stylu złożony z geometrii warstwy — to jest handoff tej budowy. Barw ' +
    'i typografii w nim nie ma, bo warstwa ich nie niesie; wykaz pod spodem mówi, czego ' +
    'dokładnie brakuje.';

  const braki = document.createElement('dl');
  braki.className = 'md-inspektor__braki';
  braki.replaceChildren(...BEZ_POLA_W_KONTRAKCIE.flatMap(([nazwa, zdanie]) => wierszBraku(nazwa, zdanie)));

  const element = document.createElement('section');
  element.className = 'md-panel md-inspektor';
  oznaczWarstwe(element, 2);
  element.append(tytul, zakres, siatka, adnotacja.element, zapisStylu, oZapisie, braki);

  /** Warstwa opisywana: pierwsza zaznaczona; `null`, gdy nie ma zaznaczenia. */
  function wskazana(): DesignBoardLayer | null {
    const zaznaczone = stan.zaznaczone();
    if (zaznaczone.length === 0) return null;
    return stan.warstwy().find((warstwa) => warstwa.id === zaznaczone[0]) ?? null;
  }

  /** Przepisuje liczby z pól na warstwę; wartość nieliczbowa zostawia wymiar bez zmiany. */
  function nanies(): void {
    const warstwa = wskazana();
    if (warstwa === null) return;
    stan.ustawGeometrie(warstwa.id, {
      x: liczbaPola('x', warstwa.x ?? 0),
      y: liczbaPola('y', warstwa.y ?? 0),
      width: liczbaPola('width', warstwa.width ?? BOK_WYJSCIOWY),
      height: liczbaPola('height', warstwa.height ?? BOK_WYJSCIOWY),
    });
  }

  function liczbaPola(klucz: KluczGeometrii, zapasowa: number): number {
    const wpisana = Number.parseFloat(kontrolki.get(klucz)?.value ?? '');
    return Number.isFinite(wpisana) ? wpisana : zapasowa;
  }

  function odswiez(): void {
    const warstwa = wskazana();
    const zaznaczonych = stan.zaznaczone().length;

    if (warstwa === null) {
      zakres.textContent =
        'Żadna warstwa nie jest zaznaczona — inspektor opisuje jedną warstwę, więc nie ma ' +
        'czego opisać. Wskaż warstwę na kanwie albo w panelu warstw.';
      for (const kontrolka of kontrolki.values()) kontrolka.value = '';
      adnotacja.kontrolka.value = '';
      zapisStylu.textContent = '';
      return;
    }

    zakres.textContent =
      zaznaczonych === 1
        ? `Warstwa ${warstwa.id}${warstwa.assetId === undefined ? ' — element pomocniczy' : ' — zasób'}.`
        : `Zaznaczonych warstw: ${String(zaznaczonych)}. Inspektor opisuje pierwszą z nich ` +
          `(${warstwa.id}) — cztery liczby opisują jeden prostokąt, nie zbiór.`;

    ustawPole('x', warstwa.x ?? 0);
    ustawPole('y', warstwa.y ?? 0);
    ustawPole('width', warstwa.width ?? BOK_WYJSCIOWY);
    ustawPole('height', warstwa.height ?? BOK_WYJSCIOWY);
    adnotacja.kontrolka.value = warstwa.note ?? '';
    zapisStylu.textContent = zapisRegulStylu(warstwa);
  }

  function ustawPole(klucz: KluczGeometrii, wartosc: number): void {
    const kontrolka = kontrolki.get(klucz);
    if (kontrolka === undefined) return;
    // Liczba zaokrąglana do części setnych, bo przeciąganie wskaźnikiem daje ułamki zbyt długie dla kodu.
    kontrolka.value = String(Math.round(wartosc * 100) / 100);
  }

  return {
    element,
    odswiez,

    ogniskujAdnotacje() {
      if (wskazana() === null) return false;
      adnotacja.kontrolka.focus();
      adnotacja.kontrolka.select();
      return true;
    },
  };
}

/**
 * Reguły stylu warstwy — wydanie geometrii do przepisania w kodzie.
 *
 * Selektor bierze się z kodu warstwy, bo to jedyne oznaczenie, które warstwa
 * niesie i które wraca z rdzenia po zapisie kompozycji.
 */
export function zapisRegulStylu(warstwa: DesignBoardLayer): string {
  const wiersze = [
    'position: absolute;',
    `left: ${String(warstwa.x ?? 0)}px;`,
    `top: ${String(warstwa.y ?? 0)}px;`,
    `width: ${String(warstwa.width ?? BOK_WYJSCIOWY)}px;`,
    `height: ${String(warstwa.height ?? BOK_WYJSCIOWY)}px;`,
    `z-index: ${String(warstwa.order ?? 0)};`,
  ];
  return `[data-warstwa="${warstwa.id}"] {\n  ${wiersze.join('\n  ')}\n}`;
}

function wierszBraku(nazwa: string, zdanie: string): readonly HTMLElement[] {
  const klucz = document.createElement('dt');
  klucz.className = 'md-inspektor__brak-nazwa';
  klucz.textContent = nazwa;

  const wartosc = document.createElement('dd');
  wartosc.className = 'md-inspektor__brak-zdanie';
  wartosc.textContent = zdanie;
  return [klucz, wartosc];
}
