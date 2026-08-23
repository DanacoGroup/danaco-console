import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';

/**
 * Ster nastawy modułu Browser — obsada wspólnego menu-drzewa.
 *
 * Trzy nastawy modułu — rodzaj wyodrębnienia w pasku dolnym oraz powiązane
 * źródło i moduł docelowy w formularzu notatki — korzystają z jednego
 * mechanizmu rozwijania z `komponenty/menu-drzewo.ts`. Rozwijanie,
 * filtrowanie, haczyk i wędrówka klawiszami należą do tego mechanizmu; ten
 * plik go wyłącznie obsadza.
 *
 * Wartość nastawy mieszka w jednym polu `wybrana`. Napis na uchwycie i haczyk
 * przy pozycji biorą się z niej przy każdym przerysowaniu, a `wartosc()`
 * oddaje to samo pole.
 *
 * Podpis nad sterem jest blokiem, nie `<label>`. Etykieta bez `for` związałaby
 * się z pierwszym potomkiem dającym się etykietować — czyli z uchwytem menu —
 * i wtedy kliknięcie w podpis otwierałoby wykaz, a nazwa dostępna uchwytu
 * konkurowałaby z `aria-label`, które mechanizm składa sam z nazwy nastawy
 * i wartości bieżącej.
 *
 * Podpis staje na ekranie tam, gdzie ster sąsiaduje z polami formularza: sama
 * wartość („bez powiązania") nie mówi, czego dotyczy. Przy pasku dolnym nazwa
 * nastawy idzie wyłącznie do `aria-label` uchwytu.
 */

/** Jedna pozycja wyboru: wartość nastawy, napis dla Operatora, zdanie opisu. */
export interface PozycjaSteru {
  wartosc: string;
  etykieta: string;
  /** Zdanie mówiące, co ta pozycja robi; pominięte znaczy „ta pozycja opisu nie ma". */
  opis?: string;
}

export interface SterWyboru {
  /** Blok osadzany w pasku albo w formularzu — z podpisem, gdy podpis zamówiono. */
  element: HTMLElement;
  /** Wartość nastawy potwierdzona wyborem; pusty napis, gdy wykaz jest pusty. */
  wartosc(): string;
  /** Wymienia pozycje, utrzymując wybór, o ile nadal istnieje. */
  ustawPozycje(pozycje: readonly PozycjaSteru[]): void;
  /**
   * Ustawia nastawę bez udziału Operatora i oddaje informację, czy ją przyjęto.
   *
   * Wartość spoza wykazu nie zostaje na uchwycie — ster wraca do pozycji
   * pierwszej i oddaje `false`. Milczące przyjęcie wartości, której w wykazie
   * nie ma, byłoby uchwytem pokazującym nastawę nie do wybrania;
   * wołający dostaje więc odpowiedź i sam rozstrzyga, czy to powiedzieć
   * Operatorowi.
   */
  ustawWartosc(wartosc: string): boolean;
}

export interface OpcjeSteruWyboru {
  /** Nazwa rodzajowa nastawy — „Rodzaj wyodrębnienia", „Powiązane źródło". */
  nastawa: string;
  /** Pozycje początkowe; pierwsza z nich jest nastawą, dopóki Operator nie wybierze. */
  pozycje: readonly PozycjaSteru[];
  /** Klasa bloku — układ należy do arkusza modułu, nie do tego pliku. */
  klasa: string;
  /** Czy nazwa nastawy ma stanąć na ekranie nad sterem. */
  podpis: boolean;
  /** Wołane po każdym wyborze Operatora; pominięte znaczy „nikt tego nie słucha". */
  naZmiane?(wartosc: string): void;
}

export function utworzSterWyboru(opcje: OpcjeSteruWyboru): SterWyboru {
  let pozycje: readonly PozycjaSteru[] = opcje.pozycje;
  let wybrana = pozycje[0]?.wartosc ?? '';

  const menu = utworzMenuDrzewo({
    nastawa: opcje.nastawa,
    naWybor: (klucz) => {
      wybrana = klucz;
      odrysuj();
      opcje.naZmiane?.(klucz);
    },
  });

  const element = document.createElement('div');
  element.className = opcje.klasa;
  element.dataset['nastawa'] = opcje.nastawa;

  if (opcje.podpis) {
    const napis = document.createElement('span');
    napis.className = 'dn-pole-etykieta';
    napis.textContent = opcje.nastawa;
    element.append(napis);
  }
  element.append(menu.element);

  /** Napis uchwytu: wartość bieżąca, a gdy wykazu nie ma — powiedziane wprost. */
  function napisUchwytu(): string {
    const wpis = pozycje.find((pozycja) => pozycja.wartosc === wybrana);
    if (wpis !== undefined) return wpis.etykieta;
    return pozycje.length === 0 ? 'wykaz pusty' : 'nie wybrano';
  }

  function drzewo(): PozycjaMenu[] {
    return pozycje.map((pozycja) => ({
      rodzaj: 'wybor',
      klucz: pozycja.wartosc,
      nazwa: pozycja.etykieta,
      wybrany: pozycja.wartosc === wybrana,
      ...(pozycja.opis === undefined ? {} : { opis: pozycja.opis }),
    }));
  }

  function odrysuj(): void {
    menu.ustaw(napisUchwytu(), drzewo());
  }

  odrysuj();

  return {
    element,
    wartosc: () => wybrana,

    ustawPozycje(nowe) {
      pozycje = nowe;
      // Wybór przeżywa wymianę wykazu, jeśli jego pozycja nadal w nim jest.
      // Gdy nie jest, nastawą staje się pierwsza pozycja nowego wykazu — ster
      // nie zostaje przy wartości spoza wykazu, bo uchwyt pokazywałby wtedy
      // nastawę nie do wybrania.
      if (!nowe.some((pozycja) => pozycja.wartosc === wybrana)) {
        wybrana = nowe[0]?.wartosc ?? '';
      }
      odrysuj();
    },

    ustawWartosc(nowa) {
      const jest = pozycje.some((pozycja) => pozycja.wartosc === nowa);
      wybrana = jest ? nowa : (pozycje[0]?.wartosc ?? '');
      odrysuj();
      return jest;
    },
  };
}
