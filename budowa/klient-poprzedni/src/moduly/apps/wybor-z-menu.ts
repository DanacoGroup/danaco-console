import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';

/**
 * Obsada bibliotecznego `komponenty/menu-drzewo.ts` dla wyborów lokalnych okna.
 *
 * Rozwijanie, znacznik wyboru, opisy pozycji, wędrówka strzałkami, pole
 * szukania i zdanie o pustym wykazie należą do `komponenty/menu-drzewo.ts`.
 * Ten plik podaje mechanizmowi dane, a formy nie odtwarza.
 *
 * `okno-komunikacji/ster-nastawy.ts` tu nie wystarcza: tamta obudowa wysyła
 * klucz do rdzenia, czeka na potwierdzenie i wyróżnienie bierze z migawki
 * stanu potwierdzonego. Tutaj żaden wybór nie jedzie do rdzenia sam z siebie —
 * środowisko wdrożenia, poziom wpisu dziennika czy priorytet błędu są wejściem
 * formularza, czytanym dopiero przy naciśnięciu przycisku komendy. Nie ma więc
 * czego potwierdzać ani skąd wziąć stanu potwierdzonego, a dorabianie
 * `ster-nastawy.ts` trybu bez rdzenia dałoby dwie prawdy o tym, co ster robi
 * z kluczem.
 *
 * Menu oddaje klucz, nie `.value`, więc wybór trzyma ten plik, a nie element
 * DOM. Zachowanie odwzorowuje natywne `<select>` w dwóch rzeczach:
 *   — pusty wybór jest wartością, nie brakiem (pozycja „bez wskazania" niesie
 *     klucz pusty i tak jedzie do rdzenia — pole nieustawione),
 *   — wymiana pozycji utrzymuje wybór, o ile nadal istnieje; gdy zniknął,
 *     wybór wraca na pozycję pierwszą, tak jak robi to przeglądarka.
 *
 * Wykaz pusty nie odbiera uchwytowi klikalności — mechanizm otwiera się
 * i mówi zdaniem, że nie ma tu jeszcze czego wybrać.
 */

/** Jedna pozycja wyboru — para klucz–nazwa wraz z opcjonalnym zdaniem opisu. */
export interface PozycjaWyboruMenu {
  /** Klucz oddawany oknu; pusty jest poprawną wartością („bez wskazania"). */
  wartosc: string;
  /** Nazwa widoczna w wykazie i na uchwycie po wybraniu. */
  etykieta: string;
  /** Zdanie mówiące, co ta pozycja znaczy; pominięte znaczy „bez opisu". */
  opis?: string;
}

export interface WyborZMenu {
  /** Element montowany w formularzu albo w pasku narzędzi okna. */
  element: HTMLElement;
  /** Klucz pozycji wybranej; pusty łańcuch, gdy wybrano pozycję o pustym kluczu. */
  wartosc(): string;
  /**
   * Ustawia wybór z zewnątrz — z odpowiedzi rdzenia albo z odczytu.
   *
   * Klucz spoza wykazu nie jest wybierany i nie trafia na uchwyt: okno nie ma
   * prawa pokazać jako wybranej pozycji, której w wykazie nie ma. Oddaje wtedy
   * `false`, żeby wołający mógł to powiedzieć Operatorowi zamiast przemilczeć.
   */
  ustawWartosc(klucz: string): boolean;
  /** Wymienia pozycje; wybór zostaje, o ile nadal istnieje (patrz nagłówek). */
  ustawPozycje(pozycje: readonly PozycjaWyboruMenu[]): void;
  /**
   * Zgłasza słuchacza wyboru dokonanego przez Operatora — i wyłącznie takiego.
   *
   * Wymiana pozycji z kodu (`ustawPozycje`) ani ustawienie z odpowiedzi rdzenia
   * (`ustawWartosc`) słuchacza nie wołają, tak samo jak natywny `<select>` nie
   * wysyła `change` przy zmianie z kodu. Gdyby wołały, odświeżenie wykazu
   * w oknie filtra zleciłoby odczyt, który sam kończy się odświeżeniem wykazu.
   *
   * Rejestracja jest osobną czynnością, a nie polem konstruktora, bo okna
   * modułu Diagnostics składają powierzchnię najpierw, a podpinają obsługę
   * dopiero wtedy, gdy zna ona całą powierzchnię.
   */
  naZmiane(sluchacz: (klucz: string) => void): void;
}

/** Napis na uchwycie, gdy wykaz jest pusty i nie ma czego pokazać. */
const UCHWYT_PUSTEGO = 'brak pozycji do wyboru';

/**
 * @param nastawa rodzajowa nazwa wyboru — idzie do `aria-label`, nie na ekran.
 * @param poczatkowe wykaz startowy; pusty znaczy „wykaz wchodzi później".
 */
export function utworzWyborZMenu(
  nastawa: string,
  poczatkowe: readonly PozycjaWyboruMenu[] = [],
): WyborZMenu {
  let pozycje: readonly PozycjaWyboruMenu[] = [];
  let wybrany = '';
  let jestWybor = false;
  const sluchacze = new Set<(klucz: string) => void>();

  const menu = utworzMenuDrzewo({
    nastawa,
    naWybor: (klucz) => {
      wybrany = klucz;
      jestWybor = true;
      odrysuj();
      for (const sluchacz of [...sluchacze]) sluchacz(klucz);
    },
  });

  /** Pozycja wybrana albo `undefined`, gdy wykaz jej nie niesie. */
  function wybrana(): PozycjaWyboruMenu | undefined {
    if (!jestWybor) return undefined;
    return pozycje.find((pozycja) => pozycja.wartosc === wybrany);
  }

  function odrysuj(): void {
    const biezaca = wybrana();
    const drzewo: PozycjaMenu[] = pozycje.map((pozycja) => ({
      rodzaj: 'wybor',
      klucz: pozycja.wartosc,
      nazwa: pozycja.etykieta,
      ...(pozycja.opis === undefined ? {} : { opis: pozycja.opis }),
      wybrany: biezaca !== undefined && pozycja.wartosc === biezaca.wartosc,
    }));
    menu.ustaw(biezaca?.etykieta ?? UCHWYT_PUSTEGO, drzewo);
  }

  function ustawPozycje(nowe: readonly PozycjaWyboruMenu[]): void {
    pozycje = [...nowe];
    // Wybór utrzymany, o ile przetrwał wymianę; inaczej pierwsza pozycja,
    // tak jak w natywnym `<select>` po `replaceChildren`.
    if (!jestWybor || !pozycje.some((pozycja) => pozycja.wartosc === wybrany)) {
      const pierwsza = pozycje[0];
      wybrany = pierwsza?.wartosc ?? '';
      jestWybor = pierwsza !== undefined;
    }
    odrysuj();
  }

  ustawPozycje(poczatkowe);

  return {
    element: menu.element,

    wartosc: () => (wybrana() === undefined ? '' : wybrany),

    ustawWartosc(klucz) {
      if (!pozycje.some((pozycja) => pozycja.wartosc === klucz)) return false;
      wybrany = klucz;
      jestWybor = true;
      odrysuj();
      return true;
    },

    ustawPozycje,

    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
    },
  };
}

/**
 * Wiersz formularza: podpis nad wyborem — kształt `dn-pole` biblioteki.
 *
 * Nie jest `<label>`, bo uchwyt menu jest przyciskiem, a przycisk nie jest
 * elementem etykietowalnym: `<label>` owinięta wokół niego nie przenosi ani
 * kliknięcia, ani ogniska, więc udawałaby wiązanie, którego nie ma. Nazwę
 * nastawy niesie `aria-label` uchwytu — mechanizm stawia go sam z pola
 * `nastawa`, więc podpis jest tu dla oka, nie dla czytnika ekranu.
 */
export function wierszWyboru(etykieta: string, wybor: WyborZMenu): HTMLElement {
  const podpis = document.createElement('span');
  podpis.className = 'dn-pole-etykieta';
  podpis.textContent = etykieta;

  const element = document.createElement('div');
  element.className = 'dn-pole dm-pole mp-pole-wyboru';
  element.append(podpis, wybor.element);
  return element;
}
