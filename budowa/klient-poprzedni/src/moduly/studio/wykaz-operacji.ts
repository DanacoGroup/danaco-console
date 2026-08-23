import { Command, type Action } from '../../../../shared/contract';
import {
  utworzMenuDrzewo,
  type GalazMenu,
  type PozycjaMenu,
} from '../../komponenty/menu-drzewo';
import { KATEGORIE_OPERACJI } from './kategorie-operacji';
import { utworzPozycjeSpozaOperacji } from './pozycja-operacji';

/**
 * Wykaz operacji Tools Panel — rejestr rdzenia i wykaz dokumentacji obok siebie.
 *
 * Pozycje mają dwa pochodzenia i panel je rozróżnia. Rejestr akcji rdzenia jest
 * źródłem właściwym: nowa operacja ma być wierszem tabeli `akcja`, nie zmianą
 * w kodzie. Dopóki rejestr nie niesie operacji redakcyjnych Studia, brakujące
 * pozycje daje wykaz dokumentacji, a przy każdej pozycji stoi jej pochodzenie.
 *
 * Wybór jest jeden na cały panel, bo `studio.contextual.op` przyjmuje dokładnie
 * jeden `actionId`. Pozycja bez wiersza rejestru pozostaje wybieralna:
 * uruchomienie wychodzi do rdzenia i wraca jego odpowiedzią.
 *
 * Mechanizm rozwijania pochodzi z `komponenty/menu-drzewo.ts`: gałąź to
 * kategoria, liść to operacja, grupa to pochodzenie. Menu wnosi pole szukania
 * po przekroczeniu progu liczby liści, znacznik wyboru widoczny na gałęzi bez
 * wchodzenia w nią oraz wędrówkę strzałkami; ten plik żadnej z tych rzeczy nie
 * odtwarza. Etykieta uchwytu niesie nazwę wybranej operacji, nie nazwę rodzajową.
 *
 * Wiersz rejestru zdejmuje pozycję z wykazu dokumentacji. Oba wykazy mówią
 * o tych samych identyfikatorach `studio.<kategoria>.<kod>`, więc bez tej zapory
 * operacja pokazywałaby się dwa razy. Wykaz dokumentacji kurczy się sam, w miarę
 * jak rejestr się zapełnia, i znika w całości, gdy rejestr obejmie komplet.
 *
 * Wiersze rejestru wskazujące inną komendę stoją poza menu. Menu z biblioteki
 * nie zna pozycji niewybieralnych, a wybieralne te wiersze być nie mogą: rdzeń
 * nie sprawdza `actionId` względem katalogu komendy, więc wysłanie takiego
 * identyfikatora wraca odpowiedzią wyglądającą na udaną. Zostają widoczne obok
 * menu jako wiersze rejestru o innej komendzie.
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

/** Napis na uchwycie, dopóki operacja nie została wskazana. */
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
      // Kategoria w całości pokryta rejestrem znika razem z gałęzią: pusta
      // gałąź mówiłaby, że coś tu czeka na wiersz, choć nie czeka nic.
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
      // Wybór wskazujący pozycję, której rejestr już nie niesie, znika razem
      // z nią: uchwyt nie pokazuje operacji spoza wykazu.
      if (wybrana !== '' && nazwaWybranej() === undefined) wybrana = '';
      przerysuj();
    },

    zamknij: () => menu.zwin(),
  };
}
