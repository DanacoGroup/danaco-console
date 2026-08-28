/**
 * Ster nastawy modułu Browser obsadza wspólny mechanizm menu-drzewa: rozwijanie,
 * filtrowanie, haczyk i wędrówka klawiszami należą do mechanizmu, a ten plik podaje
 * mu pozycje wyboru i przechowuje wartość nastawy w jednym polu.
 */
import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';

/**
 * Jedna pozycja wyboru w sterze: wartość zapisywana w nastawie, napis widoczny dla
 * Operatora oraz zdanie opisu wyświetlane przy pozycji w rozwiniętym wykazie.
 */
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
  /** Ustawia nastawę bez udziału Operatora; oddaje `false` przy wartości spoza wykazu. */
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
      // Wybór przeżywa wymianę wykazu tylko wtedy, gdy jego pozycja nadal w nim stoi.
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
