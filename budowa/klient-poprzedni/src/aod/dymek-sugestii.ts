import { utworzCzteryStery, type OpisCzterechSterow } from './cztery-stery';
import { utworzAkapit } from './pola-wykazu';
import { NAZWY_RODZAJOW, NAZWY_WAG, POWSTANIE_RODZAJU } from './rodzaje-sugestii';
import type { WpisKolejkiDecyzji } from './kolejka-decyzji';
import { NAZWY_POWODOW, WagaDecyzji, opiszOdstep } from './rozpoznanie-decyzji';
import type { Queue } from '../../../shared/contract';

/**
 * Dymek kontekstowy sugestii — warstwa 2 funkcji, popover przy awatarze
 * (`docs/funkcje-globalne/always-on-display.md`, rozdz. 2.1, 2.6, 4).
 *
 * Dymek przekazuje sugestię BEZ PRZERYWANIA PRACY: otwiera się po lewej stronie
 * awatara, nie przesuwa kolumn obszaru roboczego, nie przyciemnia tła i nie
 * zabiera ogniska klawiatury siłą. Zamknięcie klawiszem `Esc` ani odejście
 * ogniska nie zmieniają statusu sugestii (rozdz. 2.6 i zał. A.1).
 *
 * Treść dymka mówi cztery rzeczy: co się stało (zdanie decyzji), jakiego jest
 * rodzaju (katalog rozdz. 4), jak głośno wchodzi (waga rozdz. 4.5) i czyja to
 * ocena — rdzenia czy nakładki (`WagaDecyzji`). Pod treścią stoją działania:
 * te, za którymi stoi komenda kontraktu, przyciskiem
 * (`cztery-stery.ts`), a cykl życia sugestii — „Odłóż" i „Odrzuć" — osobno,
 * bo należy do nakładki, nie do rdzenia.
 *
 * Przy wąskiej kolumnie obszaru roboczego dymek skraca się do jednego zdania
 * i działania „Rozwiń" otwierającego powierzchnię interakcji (rozdz. 2.5).
 */

export interface OpisDymkaSugestii {
  /** Wspólny opis sterów decyzji — ten sam, którym karmi się lista sugestii. */
  stery: OpisCzterechSterow;
  /** „Odłóż" — pozycja wraca do listy oczekujących, status pozostaje `nowa`. */
  naOdloz(wpis: WpisKolejkiDecyzji): void;
  /** „Odrzuć" — status `odrzucona`, sugestia nie wraca. */
  naOdrzuc(wpis: WpisKolejkiDecyzji): void;
  /** „Rozwiń" — otwiera powierzchnię interakcji na liście oczekujących. */
  naRozwin(): void;
  /** Zamknięcie dymka bez zmiany statusu — `Esc`, odejście ogniska, kliknięcie obok. */
  naZamkniecie(): void;
}

export interface DymekSugestii {
  element: HTMLElement;
  /** Pokazuje dymek z treścią wskazanej decyzji. */
  pokaz(wpis: WpisKolejkiDecyzji, kolejka?: Queue): void;
  /** Chowa dymek. Status sugestii pozostaje bez zmiany. */
  schowaj(): void;
  czyWidoczny(): boolean;
  /** Wąska kolumna — dymek skraca treść do jednego zdania (rozdz. 2.5). */
  ustawWaski(waski: boolean): void;
  /** Przesuwa dymek o szerokość otwartej kolumny. */
  ustawOdsuniecie(pikseli: number): void;
  rozlacz(): void;
}

export function utworzDymekSugestii(opis: OpisDymkaSugestii): DymekSugestii {
  const element = document.createElement('section');
  element.className = 'ao-dymek';
  element.hidden = true;
  element.setAttribute('role', 'dialog');
  element.setAttribute('aria-label', 'Dymek kontekstowy sugestii Always On Display');

  const naglowek = document.createElement('header');
  naglowek.className = 'ao-dymek__glowa';

  const rodzaj = document.createElement('span');
  rodzaj.className = 'dn-plakietka ao-dymek__rodzaj';

  const zamknij = document.createElement('button');
  zamknij.type = 'button';
  zamknij.className = 'dn-btn-ikona ao-dymek__zamknij';
  zamknij.textContent = '×';
  zamknij.setAttribute('aria-label', 'Zamknij dymek bez zmiany statusu sugestii');
  zamknij.addEventListener('click', () => opis.naZamkniecie());

  naglowek.append(rodzaj, zamknij);

  const tresc = document.createElement('div');
  tresc.className = 'ao-dymek__tresc';

  element.append(naglowek, tresc);

  /** Wpis pokazywany w tej chwili; `null`, gdy dymek jest zamknięty. */
  let biezacy: WpisKolejkiDecyzji | null = null;

  /**
   * Kolejka rdzenia dopasowana do pokazywanego wpisu.
   *
   * Trzymana obok wpisu, bo przerysowanie po zwężeniu kolumny nie ma jej skąd
   * wziąć po raz drugi — a bez niej ster wstrzymania stracił by drogę wyjścia,
   * której Operator przed chwilą używał.
   */
  let biezacaKolejka: Queue | undefined;

  /** Czy kolumna obszaru roboczego jest wąska — rozstrzyga długość treści. */
  let waski = false;

  function naKlawisz(zdarzenie: KeyboardEvent): void {
    if (element.hidden) return;
    if (zdarzenie.key !== 'Escape') return;
    // `Esc` zamyka dymek BEZ zmiany statusu sugestii (zał. A.1).
    opis.naZamkniecie();
  }

  document.addEventListener('keydown', naKlawisz);

  function narysuj(wpis: WpisKolejkiDecyzji, kolejka?: Queue): void {
    const { decyzja } = wpis;

    rodzaj.textContent = `${NAZWY_RODZAJOW[decyzja.rodzaj]} · ${NAZWY_WAG[decyzja.wagaUjawnienia]}`;
    element.dataset['rodzaj'] = decyzja.rodzaj;
    element.dataset['waga'] = decyzja.wagaUjawnienia;

    const zdanie = document.createElement('p');
    zdanie.className = 'ao-dymek__zdanie';
    zdanie.textContent = decyzja.zdanie;

    const czesci: HTMLElement[] = [zdanie];

    if (waski) {
      // Wąska kolumna: jedno zdanie i droga do powierzchni interakcji.
      czesci.push(utworzRozwin(opis.naRozwin));
      tresc.replaceChildren(...czesci);
      return;
    }

    const podpis = document.createElement('p');
    podpis.className = 'ao-dymek__podpis';
    podpis.textContent =
      `${NAZWY_POWODOW[decyzja.powod]} · czeka ${opiszOdstep(wpis.czekaMs)} · ` +
      (decyzja.waga === WagaDecyzji.Pewna
        ? 'stan rdzenia mówi to wprost'
        : 'OCENA NAKŁADKI, nie orzeczenie rdzenia');
    czesci.push(podpis);

    czesci.push(
      utworzAkapit('ao-granica', `Kiedy taka sugestia powstaje: ${POWSTANIE_RODZAJU[decyzja.rodzaj]}`),
    );

    czesci.push(utworzCzteryStery(opis.stery, { decyzja, kolejka }));

    czesci.push(utworzCyklZycia(wpis, opis));
    czesci.push(utworzRozwin(opis.naRozwin));

    tresc.replaceChildren(...czesci);
  }

  return {
    element,

    pokaz(wpis, kolejka) {
      biezacy = wpis;
      biezacaKolejka = kolejka;
      narysuj(wpis, kolejka);
      element.hidden = false;
    },

    schowaj() {
      element.hidden = true;
      biezacy = null;
      biezacaKolejka = undefined;
      tresc.replaceChildren();
    },

    czyWidoczny: () => !element.hidden,

    ustawWaski(nowy) {
      if (nowy === waski) return;
      waski = nowy;
      element.dataset['waski'] = String(waski);
      if (biezacy !== null) narysuj(biezacy, biezacaKolejka);
    },

    ustawOdsuniecie(pikseli) {
      element.style.setProperty('--ao-odsuniecie', `${Math.max(0, pikseli)}px`);
    },

    rozlacz() {
      document.removeEventListener('keydown', naKlawisz);
      element.remove();
    },
  };
}

/**
 * Cykl życia sugestii — „Odłóż" i „Odrzuć" (rozdz. 9.4).
 *
 * Oba działania należą do nakładki, nie do rdzenia: kontrakt nie ma statusu
 * sugestii ani komendy go zmieniającej. Skutek jest miejscowy i tak jest
 * nazwany, żeby nikt nie wziął go za zapis w bazie rdzenia.
 */
function utworzCyklZycia(wpis: WpisKolejkiDecyzji, opis: OpisDymkaSugestii): HTMLElement {
  const pas = document.createElement('div');
  pas.className = 'ao-dymek__cykl';

  const odloz = document.createElement('button');
  odloz.type = 'button';
  odloz.className = 'dn-btn dn-btn--zarys ao-dymek__przycisk';
  odloz.textContent = 'Odłóż';
  odloz.title = 'Pozycja wraca do listy oczekujących, dymek się zamyka, status pozostaje „nowa".';
  odloz.addEventListener('click', () => opis.naOdloz(wpis));

  const odrzuc = document.createElement('button');
  odrzuc.type = 'button';
  odrzuc.className = 'dn-btn dn-btn--zarys ao-dymek__przycisk';
  odrzuc.textContent = 'Odrzuć';
  odrzuc.title = 'Status „odrzucona"; ta sugestia nie wraca do listy oczekujących.';
  // Bez pytania „czy na pewno" — decyzja Operatora jest decyzją.
  odrzuc.addEventListener('click', () => opis.naOdrzuc(wpis));

  const granica = utworzAkapit(
    'ao-granica',
    'Status sugestii żyje w nakładce: kontrakt nie ma pola statusu ani komendy, która by go ' +
      'zapisała w rdzeniu. Odłożenie i odrzucenie obowiązują do końca życia tej warstwy.',
  );

  pas.append(odloz, odrzuc, granica);
  return pas;
}

/** Działanie „Rozwiń" — jedyne wyjście z dymka do powierzchni interakcji. */
function utworzRozwin(naRozwin: () => void): HTMLElement {
  const rozwin = document.createElement('button');
  rozwin.type = 'button';
  rozwin.className = 'dn-btn dn-btn--zarys ao-dymek__rozwin';
  rozwin.textContent = 'Rozwiń';
  rozwin.title = 'Otwiera powierzchnię interakcji z listą oczekujących sugestii.';
  rozwin.addEventListener('click', naRozwin);
  return rozwin;
}
