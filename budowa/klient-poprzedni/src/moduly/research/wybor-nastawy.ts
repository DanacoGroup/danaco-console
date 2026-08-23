import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';

/**
 * Obsada bibliotecznego `komponenty/menu-drzewo.ts` dla trzech nastaw modułu
 * Research: rodzaju źródła, oceny wiarygodności i formatu eksportu.
 *
 * Ster nastawy niesie wartość bieżącą na uchwycie — „PDF", nie „Format
 * wyjściowy" — bo to ona jest odpowiedzią na pytanie „co jest teraz
 * ustawione". Natywny `<select>` tego nie robi: pokazuje wartość dopiero po
 * rozwinięciu.
 *
 * Plik nie jest drugim mechanizmem rozwijania: rozwijanie, znacznik wyboru,
 * opisy pozycji, wędrówka strzałkami, pole szukania po progu i zdanie o pustym
 * wykazie należą do biblioteki. Ten plik podaje mechanizmowi dane.
 *
 * Ta sama obsada stoi też w `moduly/apps/wybor-z-menu.ts` (wraz z wymianą
 * pozycji w locie) i w `moduly/browser/ster-wyboru.ts`. Miejscem docelowym
 * jest `komponenty/`; do czasu przeniesienia moduł nie sięga po cudzą obsadę
 * przez granicę modułu, bo import między modułami wiąże je mocniej niż
 * powtórzenie kilkudziesięciu wierszy.
 *
 * Wykaz jest stały — to jedyna różnica wobec obsady Apps. Trzy nastawy
 * Research pochodzą z wyliczeń kontraktu (`ResearchSourceKind`,
 * `ResearchCredibility`, `ExportFormat`), a nie z odpowiedzi rdzenia, więc
 * wymiany pozycji w locie tu nie ma.
 *
 * Uchwyt jest klikalny zawsze, także zanim Operator cokolwiek wybrał.
 */

/** Jedna pozycja wyboru — para klucz–nazwa wraz z opcjonalnym zdaniem opisu. */
export interface PozycjaNastawy {
  /** Wartość jadąca do rdzenia — napis wyliczenia kontraktu. */
  wartosc: string;
  /** Nazwa widoczna w wykazie i na uchwycie po wybraniu. */
  etykieta: string;
  /** Zdanie mówiące, co ta pozycja znaczy; pominięte znaczy „bez opisu". */
  opis?: string;
}

export interface WyborNastawy {
  /** Element montowany w formularzu okna. */
  element: HTMLElement;
  /** Wartość bieżąca — ta sama, którą niesie uchwyt. */
  wartosc(): string;
  /**
   * Ustawia wybór z zewnątrz. Klucz spoza wykazu nie jest wybierany i oddaje
   * `false`, żeby okno nie pokazało jako wybranej pozycji, której w wykazie
   * nie ma.
   */
  ustawWartosc(klucz: string): boolean;
}

/**
 * @param nastawa nazwa rodzajowa wyboru — idzie do `aria-label`, nie na ekran
 * @param pozycje wykaz stały; pierwsza pozycja jest wyborem początkowym, tak
 *   jak w natywnej liście wyboru
 * @param naZmiane wołane wyłącznie po wyborze Operatora — nie po `ustawWartosc`,
 *   tak jak natywny `<select>` nie wysyła `change` przy zmianie z kodu
 */
export function utworzWyborNastawy(
  nastawa: string,
  pozycje: readonly PozycjaNastawy[],
  naZmiane?: (wartosc: string) => void,
): WyborNastawy {
  let wybrany = pozycje[0]?.wartosc ?? '';

  const menu = utworzMenuDrzewo({
    nastawa,
    naWybor: (klucz) => {
      wybrany = klucz;
      odrysuj();
      naZmiane?.(klucz);
    },
  });

  function odrysuj(): void {
    const biezaca = pozycje.find((pozycja) => pozycja.wartosc === wybrany);
    const drzewo: PozycjaMenu[] = pozycje.map((pozycja) => ({
      rodzaj: 'wybor',
      klucz: pozycja.wartosc,
      nazwa: pozycja.etykieta,
      ...(pozycja.opis === undefined ? {} : { opis: pozycja.opis }),
      wybrany: pozycja.wartosc === wybrany,
    }));
    // Uchwyt niesie wartość bieżącą; wykaz pusty mówi to zdaniem, a nie pustym
    // przyciskiem bez znaczenia.
    menu.ustaw(biezaca?.etykieta ?? 'brak pozycji do wyboru', drzewo);
  }

  odrysuj();

  return {
    element: menu.element,
    wartosc: () => wybrany,

    ustawWartosc(klucz) {
      if (!pozycje.some((pozycja) => pozycja.wartosc === klucz)) return false;
      wybrany = klucz;
      odrysuj();
      return true;
    },
  };
}

/**
 * Wiersz formularza: podpis nad sterem — kształt `dn-pole` biblioteki.
 *
 * Nie jest `<label>`, bo uchwyt menu jest przyciskiem, a przycisk nie jest
 * elementem etykietowalnym: `<label>` owinięta wokół niego nie przeniosłaby ani
 * kliknięcia, ani ogniska, więc udawałaby wiązanie, którego nie ma. Nazwę
 * nastawy niesie `aria-label` uchwytu, stawiany przez mechanizm z pola
 * `nastawa` — podpis jest tu dla oka, nie dla czytnika ekranu.
 */
export function wierszNastawy(etykieta: string, wybor: WyborNastawy): HTMLElement {
  const podpis = document.createElement('span');
  podpis.className = 'dn-pole-etykieta';
  podpis.textContent = etykieta;

  const element = document.createElement('div');
  element.className = 'dn-pole mr-nastawa';
  element.append(podpis, wybor.element);
  return element;
}
