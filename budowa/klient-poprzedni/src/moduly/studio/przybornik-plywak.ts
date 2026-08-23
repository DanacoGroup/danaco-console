import { nazwaOperacji } from './kategorie-operacji';
import {
  SKALA_SUWAKA,
  WIELKOSCI_CIAGLE,
  opiszNastawe,
  type NastawySuwakow,
} from './suwaki-koncepcyjne';
import { TrybOperacji } from './przybornik-uzycie';

/**
 * Pływak kontekstowy przy zaznaczeniu — narzędzia ukryte, a nie stały panel.
 *
 * ── Rozstrzygnięcie Właściciela ─────────────────────────────────────────────
 * Katalog operacji nie może zjadać stałej kolumny powierzchni. Pływak pojawia
 * się w chwili zaznaczenia fragmentu i znika po jego zdjęciu; gdy Operator
 * z operacji nie korzysta, nie zajmuje ani jednego punktu ekranu.
 *
 * ── Trzy warstwy pływaka ────────────────────────────────────────────────────
 *   1. **czynności na wierzchu** — kilka najczęstszych, wzięte z policzonego
 *      użycia (`przybornik-uzycie.ts`), z przypięciem własnym przy każdej;
 *   2. **uchwyt pełnego katalogu** — reszta czynności pod jednym uchwytem
 *      (`przybornik-katalog.ts`), z grupami, opisami i szukaniem po nazwie;
 *   3. **suwaki wielkości ciągłych** — objętość, ton, rejestr, poziom szczegółu
 *      i stopień dopracowania jako sterowanie CIĄGŁE, nie przycisk dający skok
 *      o nieznanej wielkości.
 *
 * ── Nie zasłania zaznaczenia ────────────────────────────────────────────────
 * Wymaganie wprost: pływak staje **nad** zaznaczeniem, a gdy nad nim nie ma
 * miejsca — **pod** nim. Rozstrzyga to zmierzona wysokość pływaka i odległość
 * kursora od górnej krawędzi powierzchni, nie stała wpisana w kod.
 *
 * ── Stały panel zostaje, ale jako wybór ─────────────────────────────────────
 * Przełącznik trybu stoi na pływaku, bo pływak jest zawsze pod ręką: Operator
 * może wrócić do stałego panelu bocznego jednym naciśnięciem, a nastawa jest
 * pamiętana w rdzeniu. Wybór jest jawny i odwracalny — tak stanowi zasada ogólna
 * zlecenia.
 */

/** Czynności pływaka zlecane oknu. */
export interface CzynnosciPlywaka {
  /** Uruchamia operację o wskazanym identyfikatorze. */
  naOperacje(idAkcji: string): void;
  /** Przypina czynność do wierzchu pływaka albo zdejmuje przypięcie. */
  naPrzypiecie(idAkcji: string): void;
  /** Przestawia wielkość ciągłą. */
  naSuwak(kod: string, wartosc: number): void;
  /** Przestawia tryb wykazu operacji: narzędzia ukryte albo stały panel. */
  naTryb(tryb: TrybOperacji): void;
}

/** Pływak wraz z jego sterowaniem. */
export interface PlywakOperacji {
  element: HTMLElement;
  /** Przerysowuje czynności na wierzchu wraz z ich przypięciem. */
  ustawCzynnosci(
    naWierzchu: readonly string[],
    przypieta: (idAkcji: string) => boolean,
    opisKolejnosci: string,
  ): void;
  /** Przestawia suwaki na podane nastawy. */
  ustawSuwaki(nastawy: NastawySuwakow): void;
  /** Zapisuje tryb wykazu, żeby przełącznik mówił prawdę. */
  ustawTryb(tryb: TrybOperacji): void;
  /**
   * Ustawia pływak wobec kursora tak, żeby nie zasłonił zaznaczenia.
   *
   * `wysokoscWiersza` jest wysokością wiersza treści — pod nią pływak stanie,
   * gdy nad kursorem nie ma dla niego miejsca.
   */
  ustawPolozenie(polozenie: { x: number; y: number }, wysokoscWiersza: number): void;
  /** Gdzie pływak ostatecznie stanął — do zdania paska stanu i sprawdzianu. */
  strona(): 'nad' | 'pod';
}

/** Ile czynności staje na wierzchu pływaka. */
export const CZYNNOSCI_NA_WIERZCHU = 4;

/** Odstęp pływaka od wiersza treści w punktach. */
const ODSTEP_OD_TRESCI = 8;

export function utworzPlywakOperacji(
  czynnosci: CzynnosciPlywaka,
  katalog: HTMLElement,
): PlywakOperacji {
  let gdzieStanal: 'nad' | 'pod' = 'nad';
  let trybBiezacy: TrybOperacji = TrybOperacji.Ukryte;

  const szybkie = document.createElement('div');
  szybkie.className = 'ms-plywak__szybkie';
  szybkie.setAttribute('aria-label', 'Czynności najczęstsze na tym fragmencie');

  const podstawa = document.createElement('p');
  podstawa.className = 'dn-pole-opis ms-plywak__podstawa';

  const suwaki = new Map<string, HTMLInputElement>();
  const zdaniaSuwakow = new Map<string, HTMLElement>();

  const rzadSuwakow = document.createElement('div');
  rzadSuwakow.className = 'ms-plywak__suwaki';
  for (const wielkosc of WIELKOSCI_CIAGLE) {
    const etykieta = document.createElement('label');
    etykieta.className = 'ms-suwak__etykieta';
    etykieta.textContent = wielkosc.nazwa;

    const kontrolka = document.createElement('input');
    kontrolka.type = 'range';
    kontrolka.className = 'ms-suwak__kontrolka';
    kontrolka.min = String(SKALA_SUWAKA.dol);
    kontrolka.max = String(SKALA_SUWAKA.gora);
    kontrolka.step = String(SKALA_SUWAKA.krok);
    kontrolka.value = String(wielkosc.neutralna);
    kontrolka.dataset['suwak'] = wielkosc.kod;
    kontrolka.setAttribute(
      'aria-label',
      `${wielkosc.nazwa}: od „${wielkosc.koniecDolny}" do „${wielkosc.koniecGorny}"`,
    );
    kontrolka.title = wielkosc.opis;

    const zdanie = document.createElement('p');
    zdanie.className = 'dn-pole-opis ms-suwak__zdanie';
    zdanie.textContent = opiszNastawe(wielkosc, wielkosc.neutralna);

    kontrolka.addEventListener('input', () => {
      const wartosc = Number(kontrolka.value);
      zdanie.textContent = opiszNastawe(wielkosc, wartosc);
      czynnosci.naSuwak(wielkosc.kod, wartosc);
    });

    const uruchom = document.createElement('button');
    uruchom.type = 'button';
    uruchom.className = 'dn-btn dn-btn--sm dn-btn--duch';
    uruchom.textContent = 'Zastosuj';
    // Klucz zbioru `dataset` jest w postaci wielbłądziej; myślnik w nazwie
    // klucza jest błędem składni, a nie nazwą atrybutu. Atrybut wychodzi z tego
    // jako `data-suwak-uruchom`.
    uruchom.dataset['suwakUruchom'] = wielkosc.kod;
    uruchom.addEventListener('click', () => czynnosci.naOperacje(wielkosc.idAkcji));

    const pozycja = document.createElement('div');
    pozycja.className = 'ms-suwak';
    pozycja.append(etykieta, kontrolka, uruchom, zdanie);
    rzadSuwakow.append(pozycja);

    suwaki.set(wielkosc.kod, kontrolka);
    zdaniaSuwakow.set(wielkosc.kod, zdanie);
  }

  const przelacznikTrybu = document.createElement('button');
  przelacznikTrybu.type = 'button';
  przelacznikTrybu.className = 'dn-btn dn-btn--sm dn-btn--duch';
  przelacznikTrybu.dataset['czynnosc'] = 'tryb-operacji';
  przelacznikTrybu.addEventListener('click', () => {
    czynnosci.naTryb(
      trybBiezacy === TrybOperacji.Panel ? TrybOperacji.Ukryte : TrybOperacji.Panel,
    );
  });

  const uchwytKatalogu = document.createElement('div');
  uchwytKatalogu.className = 'ms-plywak__katalog';
  uchwytKatalogu.append(katalog);

  const element = document.createElement('div');
  element.className = 'ms-plywak';
  element.setAttribute('aria-label', 'Narzędzia ukryte — czynności na zaznaczonym fragmencie');
  element.append(szybkie, uchwytKatalogu, rzadSuwakow, podstawa, przelacznikTrybu);

  function ustawTryb(tryb: TrybOperacji): void {
    trybBiezacy = tryb;
    element.dataset['tryb'] = tryb;
    przelacznikTrybu.textContent =
      tryb === TrybOperacji.Panel
        ? 'Schowaj stały panel operacji'
        : 'Pokaż stały panel operacji';
    przelacznikTrybu.title =
      tryb === TrybOperacji.Panel
        ? 'Stały panel boczny stoi otwarty i zajmuje kolumnę. Naciśnięcie wraca do narzędzi ' +
          'ukrytych; nastawa jedzie do rdzenia komendą config.set i jest pamiętana.'
        : 'Katalog operacji stoi dziś jako narzędzia ukryte: ten pływak i uchwyt katalogu. ' +
          'Naciśnięcie otwiera stały panel boczny — jest to tryb do wyboru, nie postać domyślna.';
  }

  ustawTryb(TrybOperacji.Ukryte);

  return {
    element,

    ustawCzynnosci(naWierzchu, przypieta, opisKolejnosci) {
      szybkie.replaceChildren(
        ...naWierzchu.map((idAkcji) => przybornikCzynnoscSzybka(idAkcji, przypieta(idAkcji), czynnosci)),
      );
      podstawa.textContent = opisKolejnosci;
    },

    ustawSuwaki(nastawy) {
      for (const wielkosc of WIELKOSCI_CIAGLE) {
        const kontrolka = suwaki.get(wielkosc.kod);
        const zdanie = zdaniaSuwakow.get(wielkosc.kod);
        const wartosc = nastawy[wielkosc.kod] ?? wielkosc.neutralna;
        if (kontrolka !== undefined) kontrolka.value = String(wartosc);
        if (zdanie !== undefined) zdanie.textContent = opiszNastawe(wielkosc, wartosc);
      }
    },

    ustawTryb,

    ustawPolozenie(polozenie, wysokoscWiersza) {
      element.style.left = `${Math.max(0, polozenie.x)}px`;
      // Wysokość pływaka mierzy się na elemencie, a nie zakłada: liczba czynności
      // na wierzchu i liczba suwaków mogą się zmienić, a wpisana stała przestałaby
      // wtedy mówić prawdę i pływak zasłoniłby tekst.
      const wysokosc = element.offsetHeight;
      const zmiescSieNad = polozenie.y - wysokosc - ODSTEP_OD_TRESCI >= 0;
      gdzieStanal = zmiescSieNad ? 'nad' : 'pod';
      element.dataset['strona'] = gdzieStanal;
      element.style.top = zmiescSieNad
        ? `${polozenie.y - wysokosc - ODSTEP_OD_TRESCI}px`
        : `${polozenie.y + Math.max(wysokoscWiersza, ODSTEP_OD_TRESCI)}px`;
    },

    strona: () => gdzieStanal,
  };
}

/** Jedna czynność na wierzchu pływaka wraz z jej przypięciem. */
function przybornikCzynnoscSzybka(
  idAkcji: string,
  przypieta: boolean,
  czynnosci: CzynnosciPlywaka,
): HTMLElement {
  const nazwa = nazwaOperacji(idAkcji);

  const uruchom = document.createElement('button');
  uruchom.type = 'button';
  uruchom.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  uruchom.textContent = nazwa === undefined ? idAkcji : nazwa.split(' · ')[0] ?? idAkcji;
  uruchom.dataset['operacja'] = idAkcji;
  uruchom.title =
    (nazwa ?? idAkcji) +
    ' — idzie komendą studio.contextual.op na zaznaczonym fragmencie. Wynik wchodzi jako zmiana ' +
    'śledzona autora „model".';
  uruchom.addEventListener('click', () => czynnosci.naOperacje(idAkcji));

  const przypnij = document.createElement('button');
  przypnij.type = 'button';
  przypnij.className = 'dn-btn dn-btn--sm dn-btn--duch ms-plywak__przypiecie';
  przypnij.textContent = przypieta ? 'odepnij' : 'przypnij';
  przypnij.dataset['przypiecie'] = idAkcji;
  przypnij.setAttribute(
    'aria-label',
    przypieta
      ? `Zdejmij przypięcie czynności ${nazwa ?? idAkcji}`
      : `Przypnij czynność ${nazwa ?? idAkcji} na wierzchu pływaka`,
  );
  przypnij.title =
    'Przypięta czynność stoi na wierzchu niezależnie od liczby użyć. Nastawa jedzie do rdzenia ' +
    'i jest pamiętana.';
  przypnij.addEventListener('click', () => czynnosci.naPrzypiecie(idAkcji));

  const pozycja = document.createElement('div');
  pozycja.className = 'ms-plywak__czynnosc';
  pozycja.dataset['przypieta'] = przypieta ? 'tak' : 'nie';
  pozycja.append(uruchom, przypnij);
  return pozycja;
}
