import type { SettingCategory } from '../../../shared/contract';
import { czyNazwaIkony, elementIkony } from '../ikony/ikony';

/**
 * Lewa kolumna okna konfiguracji — kategorie z katalogu. Ani jedna pozycja nie
 * jest zapisana w kodzie: wykaz przychodzi komendą `settings.category.list`
 * wraz z porządkiem, ikoną i zagnieżdżeniem, więc nowa kategoria to nowy
 * wiersz katalogu.
 */
export interface NawigacjaKategorii {
  /** Kolumna osadzana w ciele okna. */
  element: HTMLElement;
  /** Przebudowuje wykaz i utrzymuje wybór, o ile kategoria nadal istnieje. */
  odswiez(kategorie: readonly SettingCategory[]): void;
  /** Kategoria czynna; `null`, gdy katalog jest pusty. */
  wybrana(): SettingCategory | null;
  /** Subskrypcja wyboru kategorii. */
  naWybor(sluchacz: (kategoria: SettingCategory) => void): void;
}

export function utworzNawigacjeKategorii(): NawigacjaKategorii {
  const element = document.createElement('nav');
  element.className = 'dn-boczna dk-kategorie';
  element.setAttribute('aria-label', 'Kategorie konfiguracji');

  const sluchacze: Array<(kategoria: SettingCategory) => void> = [];
  let wykaz: readonly SettingCategory[] = [];
  let czynna: SettingCategory | null = null;

  function wybierz(kategoria: SettingCategory): void {
    czynna = kategoria;
    oznacz();
    for (const sluchacz of [...sluchacze]) sluchacz(kategoria);
  }

  function oznacz(): void {
    for (const pozycja of element.querySelectorAll<HTMLElement>('[data-kategoria]')) {
      const wybrana = pozycja.dataset.kategoria === czynna?.id;
      if (wybrana) pozycja.setAttribute('aria-current', 'page');
      else pozycja.removeAttribute('aria-current');
    }
  }

  function przebuduj(): void {
    const pozycje: HTMLElement[] = [];
    for (const kategoria of najwyzszyPoziom(wykaz)) {
      pozycje.push(przycisk(kategoria, false, wybierz));
      for (const dziecko of dzieci(wykaz, kategoria.id)) {
        pozycje.push(przycisk(dziecko, true, wybierz));
      }
    }
    element.replaceChildren(...pozycje);
    oznacz();
  }

  return {
    element,

    /**
     * Przebudowa wykazu nie ogłasza wyboru: formularzem kieruje warstwa,
     * która wykaz zamontowała.
     */
    odswiez(kategorie) {
      wykaz = kategorie;
      const nadal = kategorie.find((kategoria) => kategoria.id === czynna?.id) ?? null;
      czynna = nadal ?? kategorie[0] ?? null;
      przebuduj();
    },

    wybrana: () => czynna,

    naWybor: (sluchacz) => void sluchacze.push(sluchacz),
  };
}

/**
 * Kategorie najwyższego poziomu: bez rodzica albo z rodzicem spoza katalogu,
 * ponieważ niespójność katalogu nie może odebrać dostępu do ustawień.
 */
function najwyzszyPoziom(wykaz: readonly SettingCategory[]): SettingCategory[] {
  const znane = new Set(wykaz.map((kategoria) => kategoria.id));
  return wykaz.filter((kategoria) => {
    const rodzic = kategoria.parentId;
    return rodzic === undefined || rodzic === '' || !znane.has(rodzic);
  });
}

/**
 * Kategorie podrzędne wskazanej kategorii, w kolejności z katalogu, dobierane
 * po wypełnionym wskazaniu kategorii nadrzędnej.
 */
function dzieci(wykaz: readonly SettingCategory[], rodzic: string): SettingCategory[] {
  return wykaz.filter((kategoria) => kategoria.parentId === rodzic);
}

/**
 * Pozycja nawigacji; ikona nieznana zestawowi jest pomijana, a nie zastępowana
 * zamiennikiem, który sugerowałby znaczenie nieobecne w katalogu.
 */
function przycisk(
  kategoria: SettingCategory,
  wcieta: boolean,
  naWybor: (kategoria: SettingCategory) => void,
): HTMLElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = wcieta
    ? 'dn-boczna-pozycja dk-kategorie__pozycja dk-kategorie__pozycja--wcieta'
    : 'dn-boczna-pozycja dk-kategorie__pozycja';
  element.dataset.kategoria = kategoria.id;
  if (kategoria.description !== undefined) element.title = kategoria.description;

  const nazwa = document.createElement('span');
  nazwa.className = 'dk-kategorie__nazwa';
  nazwa.textContent = kategoria.name;

  const ikona = kategoria.icon;
  if (ikona !== undefined && czyNazwaIkony(ikona)) {
    element.append(elementIkony(ikona, { rozmiar: 16 }));
  }
  element.append(nazwa);
  element.addEventListener('click', () => naWybor(kategoria));

  return element;
}
