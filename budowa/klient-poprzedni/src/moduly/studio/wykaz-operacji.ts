import { Command, type Action } from '../../../../shared/contract';
import {
  utworzMenuDrzewo,
  type GalazMenu,
  type PozycjaMenu,
} from '../../komponenty/menu-drzewo';
import { KATEGORIE_OPERACJI } from './kategorie-operacji';
import { utworzPozycjeSpozaOperacji } from './pozycja-operacji';

/**
 * Wykaz operacji panelu narzędzi łączy rejestr komend rdzenia z wykazem dokumentacji
 * redakcyjnej Studia i pozwala wskazać jedną operację kontekstową.
 */
export interface WykazOperacji {
  element: HTMLElement;
  /** Identyfikator wybranej operacji; pusty, gdy żadnej nie wybrano. */
  wybrana(): string;
  /** Wstawia pozycje rejestru rdzenia obok wykazu dokumentacji. */
  ustawRejestr(akcje: readonly Action[]): void;
  /** Zwija menu i zdejmuje jego nasłuchy dokumentu. */
  zamknij(): void;
}

const OPIS_REJESTRU_PUSTEGO =
  'Komenda action.list nie oddała dla zasięgu modułu Studio ani jednego wiersza.';
const OPIS_DOKUMENTACJI =
  'Pozycje bez wiersza rejestru. Wybór jest czynny: uruchomienie wychodzi do rdzenia ' +
  'i wraca jego odpowiedzią.';

/** Napis widoczny na uchwycie wykazu operacji, dopóki użytkownik nie wskaże żadnej operacji kontekstowej z rejestru rdzenia ani z wykazu dokumentacji. */
const UCHWYT_BEZ_WYBORU = 'wskaż operację kontekstową';

/**
 * Zdanie o zawartości rejestru, z rozbiciem na wiersze zdatne do wysłania jako
 * operacja kontekstowa i pozostałe. Liczby powstają przy każdym odczycie, więc
 * zdanie nadąża za rejestrem bez zmiany w kodzie.
 */
function opisRejestru(wszystkie: number, kontekstowe: number): string {
  if (kontekstowe === wszystkie) {
    return `Pozycje przysłane komendą action.list dla zasięgu modułu Studio (${wszystkie}).`;
  }
  return (
    `Komenda action.list oddała wierszy zasięgu modułu Studio: ${wszystkie}. Komendę ` +
    `${Command.StudioContextualOp} wskazuje z nich: ${kontekstowe}. Pozostałe są akcjami innych ` +
    'komend i wybrać ich jako operacji kontekstowej nie można — rdzeń przyjąłby ich ' +
    'identyfikator bez sprawdzenia i oddał wynik wyglądający na udany.'
  );
}

/**
 * Nazwa grupy wykazu dokumentacji; podaje, ile pozycji rejestr już przejął.
 * Liczba powstaje przy każdym odczycie, więc nazwa nadąża za rejestrem sama.
 */
function nazwaGrupyDokumentacji(przejete: number): string {
  if (przejete === 0) return 'Siedem kategorii wykazu dokumentacji';
  return `Wykaz dokumentacji — bez wiersza rejestru (rejestr przejął już: ${przejete})`;
}

export function utworzWykazOperacji(naWybor: () => void): WykazOperacji {
  let zRejestru: readonly Action[] = [];
  let wybrana = '';

  const menu = utworzMenuDrzewo({
    nastawa: 'Operacja kontekstowa',
    naWybor(klucz) {
      wybrana = klucz;
      przerysuj();
      naWybor();
    },
  });

  const zdanieRejestru = document.createElement('p');
  zdanieRejestru.className = 'dn-pole-opis ms-wykaz__zdanie';

  const spozaOperacji = document.createElement('div');
  spozaOperacji.className = 'ms-wykaz__grupa';

  const element = document.createElement('div');
  element.className = 'ms-wykaz';
  element.append(menu.element, zdanieRejestru, spozaOperacji);

  /** Nazwa pozycji wybranej — z rejestru albo z wykazu dokumentacji. */
  function nazwaWybranej(): string | undefined {
    const zRejestruNazwa = zRejestru.find((akcja) => akcja.id === wybrana)?.name;
    if (zRejestruNazwa !== undefined) return zRejestruNazwa;
    for (const kategoria of KATEGORIE_OPERACJI) {
      const operacja = kategoria.operacje.find((pozycja) => pozycja.id === wybrana);
      if (operacja !== undefined) return `${operacja.nazwa} · ${kategoria.nazwa}`;
    }
    return undefined;
  }

  /** Gałęzie kategorii z pozycjami, których rejestr rdzenia jeszcze nie niesie. */
  function galezieDokumentacji(wRejestrze: ReadonlySet<string>): GalazMenu[] {
    const galezie: GalazMenu[] = [];
    for (const kategoria of KATEGORIE_OPERACJI) {
      const czekajace = kategoria.operacje.filter((operacja) => !wRejestrze.has(operacja.id));
      // Kategoria w całości pokryta rejestrem znika z gałęzią: pusta gałąź sugerowałaby czekający wiersz.
      if (czekajace.length === 0) continue;
      galezie.push({
        rodzaj: 'galaz',
        klucz: `kategoria:${kategoria.kod}`,
        nazwa: kategoria.nazwa,
        dzieci: czekajace.map((operacja) => ({
          rodzaj: 'wybor',
          klucz: operacja.id,
          nazwa: operacja.nazwa,
          nazwaKrotka: operacja.id,
          wybrany: operacja.id === wybrana,
        })),
      });
    }
    return galezie;
  }

  function przerysuj(): void {
    const kontekstowe = zRejestru.filter((akcja) => akcja.command === Command.StudioContextualOp);
    const wRejestrze = new Set(kontekstowe.map((akcja) => akcja.id));
    const galezie = galezieDokumentacji(wRejestrze);

    const drzewo: PozycjaMenu[] = [];
    if (kontekstowe.length > 0) {
      drzewo.push({
        rodzaj: 'grupa',
        nazwa: 'Rejestr akcji rdzenia',
        dzieci: kontekstowe.map((akcja) => ({
          rodzaj: 'wybor',
          klucz: akcja.id,
          nazwa: akcja.name,
          nazwaKrotka: akcja.id,
          ...(akcja.description === undefined ? {} : { opis: akcja.description }),
          wybrany: akcja.id === wybrana,
        })),
      });
    }
    if (galezie.length > 0) {
      drzewo.push({
        rodzaj: 'grupa',
        nazwa: nazwaGrupyDokumentacji(kontekstowe.length),
        dzieci: galezie,
      });
    }
    menu.ustaw(nazwaWybranej() ?? UCHWYT_BEZ_WYBORU, drzewo);

    zdanieRejestru.textContent =
      zRejestru.length === 0
        ? OPIS_REJESTRU_PUSTEGO
        : `${opisRejestru(zRejestru.length, kontekstowe.length)} ${OPIS_DOKUMENTACJI}`;

    spozaOperacji.replaceChildren(
      ...zRejestru
        .filter((akcja) => akcja.command !== Command.StudioContextualOp)
        .map((akcja) => utworzPozycjeSpozaOperacji(akcja.id, akcja.name, akcja.command)),
    );
  }

  przerysuj();

  return {
    element,
    wybrana: () => wybrana,

    ustawRejestr(akcje) {
      zRejestru = akcje;
      // Wybór, który rejestr już nie niesie, znika z nią: uchwyt nie pokazuje operacji spoza wykazu.
      if (wybrana !== '' && nazwaWybranej() === undefined) wybrana = '';
      przerysuj();
    },

    zamknij: () => menu.zwin(),
  };
}
