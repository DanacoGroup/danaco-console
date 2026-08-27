import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { utworzSekcjePanelu, type SekcjePanelu } from '../../powloka/sekcje-panelu';
import type { OpisSekcji } from '../../powloka/uklad-sekcji';
import type { ZrodloSekcjiPaneli } from '../../powloka/zrodlo-sekcji-paneli';
import { PRZEDROSTEK } from './kontrolki';
import { utworzStanTresci, type StanTresci } from './stany-okna';

/**
 * Wspólna powierzchnia sekcji panelu orkiestracji: nagłówek z opisem zakresu,
 * wewnętrzny układ podsekcji i jedno miejsce meldunku pod spodem. Kolejność
 * podsekcji prowadzi wzorzec powłoki, a ten plik go woła i nic nie przestawia.
 */
export interface PowierzchniaSekcji {
  /** Element sekcji montowany w obszarze roboczym. */
  element: HTMLElement;
  /** Stan treści sekcji — jedno miejsce meldunku dla wszystkich podsekcji. */
  tresci: StanTresci;
  /** Meldunek: powodzenie i odmowa idą tą samą drogą. */
  meldunek(zdanie: string, udane: boolean): void;
  /** Wskazanie okna, do którego przypięty jest układ podsekcji. */
  ustawOkno(idOkna: string): void;
  /** Wewnętrzny układ podsekcji — do odświeżenia po zmianie okna. */
  uklad: SekcjePanelu;
}

export interface OpcjePowierzchni {
  /** Klucz sekcji panelu orkiestracji — trafia w `panelId` i w `data-sekcja`. */
  klucz: string;
  /** Nagłówek sekcji widoczny dla Operatora. */
  tytul: string;
  /** Jedno zdanie: czym ta sekcja steruje. */
  zakres: string;
  /** Podsekcje w kolejności domyślnej; układ z rdzenia może ją zmienić. */
  podsekcje: readonly OpisSekcji[];
  /** Wzorzec powłoki — źródło `panel.sections.get/set`. */
  sekcje: ZrodloSekcjiPaneli;
}

export function utworzPowierzchnieSekcji(opcje: OpcjePowierzchni): PowierzchniaSekcji {
  const tresci = utworzStanTresci();
  const meldunek = (zdanie: string, udane: boolean): void =>
    tresci.potwierdzenie(zdanie, udane);

  const tytul = document.createElement('h2');
  tytul.className = 'dm-orkiestracja__tytul';
  tytul.textContent = opcje.tytul;

  const zakres = document.createElement('p');
  zakres.className = 'dn-pole-opis';
  zakres.textContent = opcje.zakres;

  // Identyfikatorem panelu jest klucz sekcji, ten sam, którym wskazuje ją nawigacja.
  const uklad = utworzSekcjePanelu({
    zrodlo: opcje.sekcje,
    panelId: `orkiestracja-${opcje.klucz}`,
    przedrostek: PRZEDROSTEK,
    meldunek,
    sekcje: opcje.podsekcje,
  });

  const element = document.createElement('section');
  element.className = 'dm-orkiestracja';
  element.dataset['sekcja'] = opcje.klucz;
  element.setAttribute('aria-label', `Panel orkiestracji — ${opcje.tytul}`);
  element.append(tytul, zakres, uklad.element, tresci.element);

  return {
    element,
    tresci,
    meldunek,
    uklad,
    ustawOkno: (idOkna) => uklad.ustawOkno(idOkna),
  };
}

/**
 * Akapit stanu pustego — wspólny kształt meldunku dla wszystkich sekcji panelu.
 * Jeden kształt zamiast pięciu sprawia, że sekcja bez zawartości czyta się tak
 * samo niezależnie od tego, która to sekcja.
 */
export function zdaniePuste(tresc: string): HTMLElement {
  const akapit = document.createElement('p');
  akapit.className = 'dn-pole-opis';
  akapit.textContent = tresc;
  return akapit;
}

/**
 * Element konfiguracji wraz z objaśnieniem `[?]`.
 *
 * Objaśnienie idzie w `aria-label` znaku, więc dociera także tam, gdzie dymek
 * się nie pokazuje — czytnik ekranu przeczyta je bez najechania kursorem.
 */
export function zObjasnieniem(kontrolka: HTMLElement, objasnienie: string): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dm-orkiestracja__ustawienie';
  element.append(kontrolka, utworzDymekObjasnienia(objasnienie));
  return element;
}
