import './rozmowa.css';

import type { ZrodloAkcjiModulu } from '../okno-komunikacji/katalog-akcji';
import type { KontekstOkna } from '../okno-komunikacji/panel-kontekstu';
import { utworzPowierzchnieModulu } from '../okno-komunikacji/powierzchnia-modulu';
import { utworzListeWpisow } from './lista-wpisow';
import { utworzPoleWysylki } from './pole-wysylki';
import type { Rozmowa } from './rozmowa';
import type { WidokZapisu } from './widok-zapisu';
import { utworzWykazPoUkosniku } from './wykaz-po-ukosniku';
import type { DolozenieNarzedzia, ZrodloWykazuUkosnika } from './zrodlo-wykazu-ukosnika';

/**
 * Ustawienia widoku rozmowy, których w żaden sposób nie da się wyprowadzić z samej rozmowy
 * tego okna.
 */
export interface OpcjeWidokuRozmowy {
  /** Moduł okna w chwili złożenia widoku; pusty znaczy „jeszcze nieustalony". */
  modul?: string;
  /** Źródło katalogu akcji modułu (`action.list`); brak = bez panelu akcji. */
  katalogAkcji?: ZrodloAkcjiModulu;
  // Źródło wykazu po ukośniku; brak oznacza pole bez wykazu, jak w podglądzie rozmowy.
  zrodloWykazuUkosnika?: ZrodloWykazuUkosnika;
  /** Czym wykonać dołożenie narzędzia (`session.tool.attach`). */
  dolozenieNarzedzia?: DolozenieNarzedzia;
  /** Parametry wykonania okna pokazywane w panelu kontekstu. */
  kontekst?: KontekstOkna;
  // Gotowy pasek zlecenia — rząd sterów koperty; brak oznacza, że rzędu nie ma na ekranie.
  pasekZlecenia?: HTMLElement;
}

/**
 * Widok rozmowy — powierzchnia modułowa, historia wpisów oraz pole wypowiedzi samego
 * Operatora okna.
 */
export interface WidokRozmowy {
  /** Element montowany w oknie komunikacji. */
  element: HTMLElement;
  /** Ustawia ognisko na polu wpisywania. */
  ustawOgnisko(): void;
  /** Przestawia okno na wskazany moduł; historia wątku zostaje. */
  ustawModul(kod: string): void;
  /** Moduł, w którym okno pracuje w tej chwili. */
  modul(): string;
  // Przestawia widok transkryptu na jeden z czterech trybów; filtr nad tym, co okno już ma.
  ustawWidokZapisu(widok: WidokZapisu): void;
  /** Tryb widoku transkryptu, w którym okno pracuje w tej chwili. */
  widokZapisu(): WidokZapisu;
}

/**
 * Widok rozmowy jednego okna, składający powierzchnię modułową, listę wpisów i pole
 * wysyłki wiadomości.
 */
export function utworzWidokRozmowy(
  rozmowa: Rozmowa,
  opcje: OpcjeWidokuRozmowy = {},
): WidokRozmowy {
  const lista = utworzListeWpisow();

  // Wykaz po ukośniku powstaje tutaj, bo tu spotykają się źródło pozycji i droga dołożenia.
  const zrodloWykazu = opcje.zrodloWykazuUkosnika;
  const wykaz =
    zrodloWykazu === undefined
      ? undefined
      : utworzWykazPoUkosniku({
          zrodlo: zrodloWykazu,
          ...(opcje.dolozenieNarzedzia === undefined
            ? {}
            : { dolozenie: opcje.dolozenieNarzedzia }),
          naKomunikat: (zdanie) => rozmowa.zglosKomunikat(zdanie),
        });

  const pole = utworzPoleWysylki(
    {
      naWyslanie: (tekst) => rozmowa.wyslij(tekst),
      naPrzerwanie: () => rozmowa.przerwij(),
    },
    {
      ...(opcje.pasekZlecenia === undefined ? {} : { pasekZlecenia: opcje.pasekZlecenia }),
      ...(wykaz === undefined ? {} : { wykazUkosnika: wykaz }),
    },
  );

  const powierzchnia = utworzPowierzchnieModulu({
    naPolecenie: (tekst) => pole.wstaw(tekst),
    ...(opcje.katalogAkcji === undefined ? {} : { katalogAkcji: opcje.katalogAkcji }),
    ...(opcje.modul === undefined ? {} : { modul: opcje.modul }),
    kontekst: { ...opcje.kontekst, wpisyWatku: () => lista.liczba() },
  });

  const element = document.createElement('section');
  // Klasa własna, nie biblioteczna; ramę okna niesie już gniazdo układu równoległego.
  element.className = 'dc-rozmowa';

  // Pasek zlecenia stoi wewnątrz pola wysyłki, w rzędzie akcji, nie w powierzchni modułowej.
  element.append(powierzchnia.element, lista.element, pole.element);

  for (const wpis of rozmowa.wpisy()) lista.pokaz(wpis);
  pole.pokazStan(rozmowa.stan());
  powierzchnia.odswiezKontekst();

  // Wyczyszczenie rozmowy ulotnej zdejmuje pozycje wątku, nie okno, i ogłasza to osobno od
  // wpisów.
  rozmowa.naWyczyszczenie(() => {
    lista.wyczysc();
    powierzchnia.odswiezKontekst();
  });

  rozmowa.naWpis((wpis) => {
    lista.pokaz(wpis);
    // Licznik wątku w panelu kontekstu ma być prawdziwy w każdej chwili, nawet między
    // modułami.
    powierzchnia.odswiezKontekst();
  });
  rozmowa.naStan((stan) => pole.pokazStan(stan));

  // Prompt wpisany spoza tego połączenia idzie przez pole wypowiedzi, nie samą listą wpisów.
  rozmowa.naWypowiedzZZewnatrz((tresc) => pole.odbijWypowiedzZZewnatrz(tresc));

  return {
    element,
    ustawOgnisko: pole.ustawOgnisko,
    ustawModul: powierzchnia.ustawModul,
    modul: powierzchnia.modul,
    // Widok transkryptu mieszka w liście wpisów, nie w powierzchni modułowej, i przeżywa jej
    // zmianę.
    ustawWidokZapisu: lista.ustawWidokZapisu,
    widokZapisu: lista.widokZapisu,
  };
}
