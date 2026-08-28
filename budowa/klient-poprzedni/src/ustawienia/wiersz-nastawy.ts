import type { SettingDefinition, SettingOption } from '../../../shared/contract';
import { utworzMenuDrzewo, type PozycjaMenu } from '../komponenty/menu-drzewo';

/** Wiersz nastawy przedstawia jedną nastawę produktu jako ster, nie jako wyświetlacz: etykieta uchwytu niesie wartość bieżącą, a wykaz wyboru buduje się z definicji katalogu rdzenia. */
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

  // Klucz pozycji menu nie może być pusty, bo pusty klucz byłby pozycją bez tożsamości w menu.
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

  // Napis na uchwycie bierze etykietę opcji z katalogu; wartość spoza katalogu wychodzi dosłownie.
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
    // Uchwyt nie traci klikalności na czas wysyłki: znacznik zajętości mówi o pracy, nic nie odbierając.
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
