import type { SettingDefinition, SettingOption } from '../../../shared/contract';
import { utworzMenuDrzewo, type PozycjaMenu } from '../komponenty/menu-drzewo';

/**
 * Wiersz nastawy — jedna nastawa produktu jako ster, nie jako wyświetlacz.
 *
 * Etykieta uchwytu niesie wartość bieżącą, a nie napis rodzajowy: stoi na niej
 * „Ciemny", nie „Motyw" ani „Wybierz…". Nazwa rodzajowa idzie do etykiety
 * wiersza po lewej i do `aria-label` uchwytu, żeby czytnik ekranu wiedział,
 * czego dotyczy wartość, której nazwa sama tego nie mówi.
 *
 * Wybór nie stoi rozwinięty: sekcja pokazuje po jednym wierszu na nastawę,
 * a opcje rozwijają się dopiero pod kliknięciem. Sześć sekcji rozwiniętych
 * naraz byłoby sześcioma płachtami, a nie oknem ustawień.
 *
 * Uchwyt, wykaz, opisy pozycji, haczyk przy wybranej, zwijanie kliknięciem obok
 * i obsługa klawiatury należą do `komponenty/menu-drzewo.ts`; ten plik obsadza
 * ten mechanizm danymi i nie odtwarza go u siebie. Bliźniaczą obsadę ma pasek
 * zlecenia (`okno-komunikacji/ster-nastawy.ts`) — nie da się jej zaimportować,
 * bo niesie klasy arkusza tamtego okna (`dc-ster-zlecenia`,
 * `pasek-zlecenia.css`), którego to okno nie wczytuje.
 *
 * Wykaz wyboru buduje się z `SettingDefinition.options`, wraz z etykietami
 * i opisami. Klient nie zna ani jednej wartości dopuszczalnej z góry — gdy
 * katalog dołoży czwarty motyw, wiersz pokaże go bez zmiany tego pliku.
 */
export interface WierszNastawy {
  /** Element montowany w sekcji. */
  element: HTMLElement;
  /** Podaje wartość bieżącą i przerysowuje uchwyt oraz wykaz. */
  ustaw(wartosc: string): void;
  /** Zdanie pod wierszem: odmowa rdzenia albo cisza (treść pusta chowa). */
  zdanie(tresc: string): void;
}

export interface OpcjeWiersza {
  /** Definicja katalogu — źródło nazwy, opisu i wykazu wyboru. */
  definicja: SettingDefinition;
  /** Wartość bieżąca w postaci, w jakiej stoi w katalogu (`''` bywa wartością). */
  wartosc: string;
  /** Wykonanie wyboru; oddaje zdanie odmowy albo puste przy powodzeniu. */
  wykonaj(wartosc: string): Promise<string>;
}

export function utworzWierszNastawy(opcje: OpcjeWiersza): WierszNastawy {
  const { definicja } = opcje;

  const element = document.createElement('div');
  element.className = 'du-wiersz';
  element.dataset['klucz'] = definicja.key;

  const opis = document.createElement('span');
  opis.className = 'du-wiersz__nazwa';
  opis.textContent = definicja.name;

  const zdanieOdmowy = document.createElement('p');
  zdanieOdmowy.className = 'du-wiersz__zdanie';
  zdanieOdmowy.setAttribute('role', 'alert');
  zdanieOdmowy.hidden = true;

  let biezaca = opcje.wartosc;

  const menu = utworzMenuDrzewo({
    nastawa: definicja.name,
    naWybor: (klucz) => void wybierz(klucz),
  });

  const uchwyt = document.createElement('div');
  uchwyt.className = 'du-wiersz__ster';
  uchwyt.append(menu.element);

  element.append(opis, uchwyt, zdanieOdmowy);

  if (definicja.description !== undefined && definicja.description !== '') {
    const objasnienie = document.createElement('p');
    objasnienie.className = 'du-wiersz__opis';
    objasnienie.textContent = definicja.description;
    element.append(objasnienie);
  }

  function powiedz(tresc: string): void {
    zdanieOdmowy.textContent = tresc;
    zdanieOdmowy.hidden = tresc === '';
  }

  /**
   * Klucz pozycji menu nie może być wartością pustą.
   *
   * Katalog rdzenia ma opcję o wartości pustej — dla motywu jest nią
   * „Preferencja systemu" i to trzeci pełnoprawny stan nastawy, nie brak
   * wyboru. Mechanizm menu rozdaje pozycje po kluczu, więc pusty klucz byłby
   * pozycją bez tożsamości. Stąd przedrostek zdejmowany przy wyborze.
   */
  const PRZEDROSTEK = 'wartosc:';

  function drzewo(): PozycjaMenu[] {
    return (definicja.options ?? []).map((wybor: SettingOption) => ({
      rodzaj: 'wybor' as const,
      klucz: `${PRZEDROSTEK}${wybor.value}`,
      nazwa: wybor.label,
      ...(wybor.description === undefined || wybor.description === ''
        ? {}
        : { opis: wybor.description }),
      wybrany: wybor.value === biezaca,
    }));
  }

  /**
   * Napis na uchwycie: etykieta opcji z katalogu, nie wartość surowa.
   *
   * Wartość spoza katalogu wychodzi dosłownie, zamiast podmiany na pierwszą
   * opcję z brzegu — Operator ma zobaczyć, że w bazie stoi coś, czego katalog
   * nie zna.
   */
  function napisUchwytu(): string {
    const opcja = (definicja.options ?? []).find((wybor) => wybor.value === biezaca);
    if (opcja !== undefined) return opcja.label;
    return biezaca === ''
      ? 'bez wskazania'
      : `${biezaca} — wartość spoza katalogu ustawień`;
  }

  async function wybierz(klucz: string): Promise<void> {
    const wartosc = klucz.startsWith(PRZEDROSTEK) ? klucz.slice(PRZEDROSTEK.length) : klucz;
    powiedz('');
    // Uchwyt nie traci klikalności na czas wysyłki: `aria-busy` mówi o pracy,
    // niczego nie odbierając.
    element.setAttribute('aria-busy', 'true');
    menu.zwin();
    try {
      const odmowa = await opcje.wykonaj(wartosc);
      if (odmowa !== '') powiedz(odmowa);
    } finally {
      element.removeAttribute('aria-busy');
    }
  }

  function ustaw(wartosc: string): void {
    biezaca = wartosc;
    menu.ustaw(napisUchwytu(), drzewo());
  }

  ustaw(biezaca);

  return { element, ustaw, zdanie: powiedz };
}
