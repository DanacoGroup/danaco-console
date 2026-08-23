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

/** Ustawienia widoku, których nie da się wyprowadzić z samej rozmowy. */
export interface OpcjeWidokuRozmowy {
  /** Moduł okna w chwili złożenia widoku; pusty znaczy „jeszcze nieustalony". */
  modul?: string;
  /** Źródło katalogu akcji modułu (`action.list`); brak = bez panelu akcji. */
  katalogAkcji?: ZrodloAkcjiModulu;
  /**
   * Źródło wykazu po ukośniku (`tools.catalog.list`); brak = pole bez wykazu.
   *
   * Pole jest opcjonalne, tak samo jak `pasekZlecenia`: wykaz czyta rdzeń
   * i dokłada narzędzia do sesji, więc bez kanału nie ma go z czego złożyć,
   * a stanowisko podglądu rozmowy (`polaczenie-podgladu.ts`) kanału komend nie
   * prowadzi i ma dalej działać.
   */
  zrodloWykazuUkosnika?: ZrodloWykazuUkosnika;
  /** Czym wykonać dołożenie narzędzia (`session.tool.attach`). */
  dolozenieNarzedzia?: DolozenieNarzedzia;
  /** Parametry wykonania okna pokazywane w panelu kontekstu. */
  kontekst?: KontekstOkna;
  /**
   * Gotowy pasek zlecenia — rząd sterów koperty; brak = rzędu nie ma.
   *
   * Pole jest opcjonalne: stery paska potrzebują kompletu sterowania okna
   * (migawka, subskrypcja, `window.update`, `config.set`), a ten powstaje
   * dopiero w powłoce — stanowisko podglądu rozmowy (`podglad-rozmowy.ts`) go
   * nie ma i ma dalej działać. Widok przepuszcza element nietknięty: nie zna ani
   * nastaw, ani kontraktu.
   *
   * Pasek niesie komplet sterów koperty — przełącznik urządzenia, katalog
   * roboczy, model, wysiłek i tryb zatwierdzania — jako jedną całość.
   */
  pasekZlecenia?: HTMLElement;
}

/** Widok rozmowy — powierzchnia modułowa, historia i pole wypowiedzi. */
export interface WidokRozmowy {
  /** Element montowany w oknie komunikacji. */
  element: HTMLElement;
  /** Ustawia ognisko na polu wpisywania. */
  ustawOgnisko(): void;
  /** Przestawia okno na wskazany moduł; historia wątku zostaje. */
  ustawModul(kod: string): void;
  /** Moduł, w którym okno pracuje w tej chwili. */
  modul(): string;
  /**
   * Przestawia widok transkryptu na jeden z czterech trybów.
   *
   * Filtr nad tym, co okno już ma. Nie woła rdzenia, nie czyści historii i nie
   * zmienia modułu — przełączenie jest natychmiastowe i nigdy nie odmawia.
   */
  ustawWidokZapisu(widok: WidokZapisu): void;
  /** Tryb widoku transkryptu, w którym okno pracuje w tej chwili. */
  widokZapisu(): WidokZapisu;
}

/**
 * Widok rozmowy jednego okna.
 *
 * Plik wyłącznie składa: powierzchnia modułowa na górze, lista wpisów pośrodku,
 * pole wysyłki na dole. Nie zna kontraktu, nie buduje koperty i nie wie, czym
 * jest fragment strumienia.
 *
 * Historia nie znika przy zmianie modułu: lista wpisów powstaje raz, a
 * `ustawModul` sięga wyłącznie do powierzchni modułowej. Chat Window
 * przestawia wygląd, możliwości, narzędzia i kontekst, a wątek zachowuje.
 *
 * Wpisy zastane w rozmowie są rysowane przy złożeniu widoku, żeby podłączenie
 * widoku po rozpoczęciu tury nie gubiło tego, co już przyszło.
 */
export function utworzWidokRozmowy(
  rozmowa: Rozmowa,
  opcje: OpcjeWidokuRozmowy = {},
): WidokRozmowy {
  const lista = utworzListeWpisow();

  // Wykaz po ukośniku powstaje tutaj, bo tutaj spotykają się jego trzy potrzeby:
  // źródło pozycji (z powłoki), droga dołożenia (z powłoki) i wątek rozmowy,
  // w którym ma stanąć zdanie o dołożeniu. Pole wypowiedzi dostaje go gotowego
  // i steruje nim ruchami klawiatury — kontraktu nie zna ani ono, ani ten plik.
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
  // Klasa własna, nie biblioteczna. Powód stoi przy regule `.dc-rozmowa`
  // w `rozmowa.css`: `.dn-rozmowa` niesie ramę okna, a ramę tego okna niesie
  // już gniazdo układu równoległego.
  element.className = 'dc-rozmowa';

  // Pasek zlecenia stoi wewnątrz pola wysyłki, w rzędzie akcji pod obszarem
  // tekstu; widok podaje go tam jako gotowy element. Do powierzchni modułowej
  // nie wchodzi z trzech powodów:
  //   1. Jego miejsce jest przy polu wypowiedzi, a powierzchnia modułowa stoi
  //      nad zapisem, czyli po drugiej stronie historii wątku.
  //   2. Powierzchnia ma `max-height: 33%` i własne przewijanie
  //      (`powierzchnia-modulu.css`) — ster ustawiający kopertę mógłby w niej
  //      odjechać poza widok.
  //   3. Powierzchnia jest tym, co okno przestawia przy zmianie modułu,
  //      a koperta zlecenia należy do okna, nie do modułu.
  element.append(powierzchnia.element, lista.element, pole.element);

  for (const wpis of rozmowa.wpisy()) lista.pokaz(wpis);
  pole.pokazStan(rozmowa.stan());
  powierzchnia.odswiezKontekst();

  // Wyczyszczenie rozmowy ulotnej zdejmuje pozycje, nie okno. Warstwa rozmowy
  // ogłasza je osobno od wpisów, bo „nie ma już tych wpisów" jest zdarzeniem
  // innego rodzaju niż „jest nowy wpis"; zaraz po nim przychodzą zdania
  // o powodzie, więc lista nie zostaje pusta dłużej, niż trzeba.
  rozmowa.naWyczyszczenie(() => {
    lista.wyczysc();
    powierzchnia.odswiezKontekst();
  });

  rozmowa.naWpis((wpis) => {
    lista.pokaz(wpis);
    // Licznik wątku w panelu kontekstu ma być prawdziwy w każdej chwili, także
    // między zmianami modułu — inaczej dowód zachowania historii byłby pozorny.
    powierzchnia.odswiezKontekst();
  });
  rozmowa.naStan((stan) => pole.pokazStan(stan));

  // Prompt wpisany spoza tego połączenia idzie przez pole wypowiedzi, nie samą
  // listą wpisów: tekst pojawia się w polu wpisywania okna docelowego, więc
  // widać, skąd się wziął, zamiast zastanego faktu w wątku. Wpis wchodzi
  // do wątku równocześnie i bez zwłoki, bo odwrócenie kolejności postawiłoby
  // odpowiedź modelu przed pytaniem, na które odpowiada.
  rozmowa.naWypowiedzZZewnatrz((tresc) => pole.odbijWypowiedzZZewnatrz(tresc));

  return {
    element,
    ustawOgnisko: pole.ustawOgnisko,
    ustawModul: powierzchnia.ustawModul,
    modul: powierzchnia.modul,
    // Widok transkryptu mieszka w liście, nie w powierzchni modułowej. Lista
    // przeżywa zmianę modułu, więc wybrany tryb też ją przeżywa — inaczej
    // `workspace.enter` po cichu cofałby ustawienie.
    ustawWidokZapisu: lista.ustawWidokZapisu,
    widokZapisu: lista.widokZapisu,
  };
}
