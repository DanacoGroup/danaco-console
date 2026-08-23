import type { Window } from '../../../shared/contract';
import type { ObszarRoboczy } from '../powloka/obszar-roboczy';
import type { Kanal } from '../protokol/kanal';
import { zadajWejscieDoPrzestrzeni } from '../protokol/wejscie-do-przestrzeni';
import type { PozycjaModulu, Srodowisko } from '../powloka/srodowiska';
import { pokazKomunikat } from './komunikaty';
import {
  czySekcjaOrkiestracji,
  utworzWidokSekcji,
} from '../moduly/multitasking/panel-orkiestracji';
import { opisModulu, type WidokModulu } from './rejestr-modulow';

/**
 * Przestrzeń robocza modułu — co widać po wyborze pozycji bocznej nawigacji.
 * Przekłada wybór na treść obszaru roboczego i na komendę `workspace.enter`;
 * nie buduje ani jednego elementu, bo scenę i obszar dostaje gotowe.
 *
 * Scena rozmowy nie jest alternatywą dla widoku modułu: katalog rdzenia daje
 * każdemu modułowi okno rozmowy `<kod>.chat-window` obok okien operacyjnych,
 * a okno rozmowy jest oknem wiodącym modułu. Scena zostaje więc na planszy
 * zawsze, gdy pozycja ma moduł, a widok modułu staje obok niej. Stąd trzy
 * ścieżki:
 *  • pozycja z modułem i widokiem — plansza niesie widok modułu i scenę;
 *  • pozycja z modułem bez widoku — sama scena plus pasek uczciwości
 *    wymieniający okna operacyjne, których w tym module jeszcze nie ma;
 *  • pozycja bez modułu — sekcja panelu orkiestracji. Sekcja znana panelowi
 *    dostaje swój widok tak samo jak moduł; sekcja nieznana zostaje przy
 *    stanie pustym z własną nazwą.
 *
 * Panel orkiestracji jest wołany wprost, a nie przez rejestr modułów: rejestr
 * kluczuje po kodzie modułu z tabeli `modul` rdzenia, a sekcja panelu modułu
 * nie ma, bo środowisko MultitaskingAI nie udostępnia modułów w bocznej
 * nawigacji. Wpisanie sekcji do rejestru wprowadziłoby do niego byty, których
 * rdzeń za moduły nie uważa. Moduły idą rejestrem, sekcje panelem.
 *
 * `workspace.enter` idzie z identyfikatorem żywego okna rozmowy: okno wskazane
 * przestawia moduł i zachowuje historię, a dopiero jego brak zakłada okno nowe.
 * Okno założone po stronie rdzenia nie dostaje kanału modelu, więc pierwsze
 * `message.send` skończyłoby się odmową — gdy rdzeń odda okno bez kanału,
 * klient mówi o tym wprost, zamiast zgadywać.
 */
export interface ZaleznosciPrzestrzeni {
  /**
   * Kanał kontraktu — droga `workspace.enter` oraz kanał podawany widokom
   * modułów; moduł nie sięga po niego sam.
   */
  kanal: Kanal;
  /** Obszar roboczy powłoki. */
  obszar: ObszarRoboczy;
  /** Scena okien komunikacji — warstwa wspólna wszystkim modułom. */
  scena: HTMLElement;
  /** Identyfikator sesji rdzenia; pusty, dopóki uzgodnienie trwa. */
  idSesji(): string;
  /** Okno rozmowy do rekonfiguracji; puste zakłada okno nowe. */
  idOkna(): string | undefined;
}

export interface PrzestrzenModulu {
  /** Przestawia obszar roboczy na treść wybranej pozycji nawigacji. */
  wybierz(pozycja: PozycjaModulu, dane: Srodowisko): void;
  /** Ponawia wejście odłożone na czas uzgodnienia z rdzeniem. */
  ponow(): void;
}

export function utworzPrzestrzenModulu(zaleznosci: ZaleznosciPrzestrzeni): PrzestrzenModulu {
  const { kanal, obszar, scena, idSesji, idOkna } = zaleznosci;

  // Scena wchodzi na planszę raz i zostaje do końca życia powłoki: niesie żywe
  // gniazdo WebSocket, strumień odpowiedzi modelu i historię wpisów.
  obszar.osadzCzat(scena);

  /** Pozycja czekająca na sesję — wejście przed uzgodnieniem nie przepada. */
  let odlozona: PozycjaModulu | null = null;
  /**
   * Numer bieżącego wyboru pozycji. Odpowiedź na wybór porzucony nie ma prawa
   * wejść na planszę: przy trzech kliknięciach pod rząd widoczna ma zostać
   * pozycja trzecia, a nie ta, której rdzeń odpowiedział najpóźniej. Bez żetonu
   * odczyt zlecony przy pozycji pierwszej dopisywałby się do widoku już
   * zdjętego z planszy.
   */
  let zeton = 0;
  /**
   * Widoki modułów powstają raz na moduł i zostają. Odbudowa przy każdym
   * przejściu gubiłaby stan okien operacyjnych, a obszar roboczy i tak
   * przestawia widoczność, zamiast zdejmować widok z drzewa.
   *
   * Mapa jest całym cyklem życia widoków modułów: przestrzeń powstaje raz
   * (`widok-srodowiska.ts`), a router nie zdejmuje widoku opuszczonej trasy,
   * więc mapa żyje tyle, co karta przeglądarki, i nigdy się nie opróżnia.
   * `WidokModulu.zamknij` nie ma tu wołacza, bo nie ma chwili, w której byłby
   * prawdziwy; pełne uzasadnienie stoi przy tym polu w `rejestr-modulow.ts`.
   */
  const widoki = new Map<string, WidokModulu>();

  /** Widok modułu z rejestru albo `undefined`, gdy moduł nie ma jeszcze widoku. */
  function widokModulu(kod: string): WidokModulu | undefined {
    const gotowy = widoki.get(kod);
    if (gotowy !== undefined) return gotowy;
    const opis = opisModulu(kod);
    if (opis === undefined) return undefined;
    const widok = opis.utworzWidok(kanal);
    widoki.set(kod, widok);
    return widok;
  }

  /**
   * Widok sekcji panelu orkiestracji albo `undefined` dla pozycji, która sekcją
   * nie jest.
   *
   * Widoki sekcji mieszkają w tej samej mapie co widoki modułów, bo mają ten
   * sam cykl życia: powstają raz, zostają w drzewie i są wyłącznie przełączane
   * widocznością. Klucz sekcji („zespoly", „role"…) nie zderzy się z kodem
   * modułu, bo rdzeń takich kodów w tabeli `modul` nie ma.
   */
  function widokSekcji(klucz: string): WidokModulu | undefined {
    if (!czySekcjaOrkiestracji(klucz)) return undefined;
    const gotowy = widoki.get(klucz);
    if (gotowy !== undefined) return gotowy;
    const widok = utworzWidokSekcji(klucz, kanal);
    widoki.set(klucz, widok);
    return widok;
  }

  /**
   * Wejście do przestrzeni roboczej, a dopiero po nim odczyt widoku modułu.
   *
   * Kolejność jest tu treścią, nie stylem. `workspace.enter` przestawia okno
   * rozmowy na wybrany moduł po stronie rdzenia, a widok modułu szuka swojego
   * okna komendą `window.list` zawężoną do modułu — znajdzie je wyłącznie
   * wtedy, gdy przestawienie już się dokonało. Rdzeń prowadzi każde żądanie
   * osobnym biegiem i odpowiada w kolejności zależnej od czasu obsługi, więc
   * puszczenie obu komend naraz dawałoby widok modułu poprzedniego albo stan
   * pusty. Odmowa wejścia odczytu nie uruchamia: bez przestawienia okna
   * `window.list` opisałby stan cudzy, a prawdziwa przyczyna idzie
   * komunikatem.
   */
  async function wejdzIWczytaj(pozycja: PozycjaModulu, widok?: WidokModulu): Promise<void> {
    const moj = zeton;
    const sesja = idSesji();
    if (sesja === '') {
      odlozona = pozycja;
      return;
    }
    odlozona = null;

    // Sekcja panelu orkiestracji nie wchodzi w przestrzeń modułu:
    // `workspace.enter` przestawia okno rozmowy na moduł, a sekcja modułu nie
    // ma — żądanie z pustym `moduleId` skończyłoby się odmową bez czytelnej
    // przyczyny. Sekcja czyta więc wprost swoją treść.
    if (pozycja.modul === undefined) {
      if (widok !== undefined) await widok.wczytaj(sesja);
      return;
    }

    const okno = idOkna();
    const wynik = await zadajWejscieDoPrzestrzeni(kanal, {
      sessionId: sesja,
      moduleId: pozycja.modul,
      ...(okno === undefined ? {} : { windowId: okno }),
    });
    // Pozycja zmieniona w trakcie oczekiwania: to już nie jest odpowiedź na
    // pytanie, które stoi na planszy.
    if (moj !== zeton) return;

    if (!wynik.udany) {
      pokazKomunikat({
        tytul: `Rdzeń odmówił wejścia do modułu ${pozycja.nazwa}`,
        tresc: wynik.blad?.message ?? 'Odmowa bez treści.',
        waga: 'blad',
      });
      return;
    }
    ostrzezOBrakuKanalu(pozycja, wynik.wynik?.window);

    if (widok === undefined) return;
    await widok.wczytaj(sesja);
  }

  return {
    wybierz(pozycja, dane) {
      // Każdy wybór unieważnia odpowiedzi wyboru poprzedniego.
      zeton += 1;

      // Stan pusty trzymany w zgodzie z pozycją także wtedy, gdy jest ukryty —
      // powrót na sekcję bez modułu nie może pokazać poprzedniej treści.
      obszar.zapowiedz(pozycja, dane, pozycja.okna);

      if (pozycja.modul === undefined) {
        // Sekcja bez modułu odwołuje wejście odłożone: po uzgodnieniu z rdzeniem
        // nie ma wchodzić w moduł, z którego Operator już zszedł.
        odlozona = null;
        obszar.nota(null);

        // Sekcja znana panelowi orkiestracji ma swój widok i wchodzi na planszę
        // tak samo jak moduł — obok sceny okien komunikacji, nie zamiast niej.
        // Sekcja nieznana zostaje przy stanie pustym: nie ma dla niej widoku.
        const sekcja = widokSekcji(pozycja.klucz);
        if (sekcja === undefined) {
          obszar.oproznij();
          return;
        }
        obszar.pokaz(sekcja.element);
        void wejdzIWczytaj(pozycja, sekcja);
        return;
      }

      // Rejestr kluczuje po kodzie modułu (`modul.code` rdzenia), nie po jego
      // identyfikatorze: kod jest nazwą modułu w katalogu i to jego moduł
      // deklaruje u siebie, a identyfikator jest numerem wiersza i idzie
      // wyłącznie do `workspace.enter`. Podanie identyfikatora nie trafiłoby
      // w rejestr i każdy zbudowany moduł wypadałby w gałąź bez widoku.
      const widok = widokModulu(pozycja.klucz);
      if (widok === undefined) {
        // Moduł bez zbudowanego widoku: sama scena okien komunikacji, bo okno
        // rozmowy należy do katalogu każdego modułu, a pasek uczciwości mówi
        // wprost, czego jeszcze nie ma.
        obszar.nota(notaModulu(pozycja));
        obszar.pokaz(null);
        void wejdzIWczytaj(pozycja);
        return;
      }

      // Moduł zbudowany: widok modułu obok sceny, nie zamiast niej. Rozmowa
      // jest oknem wiodącym tego modułu i zostaje na planszy.
      obszar.nota(null);
      obszar.pokaz(widok.element);
      void wejdzIWczytaj(pozycja, widok);
    },

    ponow() {
      // Wejście odłożone niesie ze sobą odczyt widoku: bez niego moduł wybrany
      // przed uzgodnieniem z rdzeniem wszedłby w przestrzeń, ale nigdy nie
      // przeczytał swojej treści.
      const pozycja = odlozona;
      if (pozycja === null) return;
      // Odłożona bywa też sekcja panelu orkiestracji — jej widok stoi w tej
      // samej mapie, ale bierze się z panelu, nie z rejestru modułów. Sięgnięcie
      // tu wyłącznie po rejestr zostawiłoby sekcję bez odczytu treści.
      void wejdzIWczytaj(
        pozycja,
        pozycja.modul === undefined ? widokSekcji(pozycja.klucz) : widokModulu(pozycja.klucz),
      );
    },
  };
}

/**
 * Czy kod okna wskazuje okno rozmowy — w obu postaciach, jakie niesie rdzeń.
 *
 * Plik `migracja_030_rejestr_okien_operacyjnych.sql` zakłada kod bez
 * przedrostka (`chat-window`); postać przedrostkowana (`studio.chat-window`)
 * pojawia się tam, gdzie moduł powiela okno u siebie. Warunek pytający
 * wyłącznie o końcówkę `.chat-window` pominąłby postać bezprzedrostkową
 * i pasek uczciwości wymieniałby działające okno rozmowy jako niezbudowane.
 */
function czyOknoRozmowy(kod: string): boolean {
  return kod === 'chat-window' || kod.endsWith('.chat-window');
}

/** Zdanie paska uczciwości: co w tym module działa, a czego jeszcze nie ma. */
function notaModulu(pozycja: PozycjaModulu): string {
  const brakujace = pozycja.okna.filter((kod) => !czyOknoRozmowy(kod));
  if (brakujace.length === 0) {
    return `${pozycja.nazwa}: rdzeń nie podaje dla tego modułu okien operacyjnych poza oknem rozmowy.`;
  }
  return `${pozycja.nazwa}: działa okno rozmowy. Okna operacyjne tego modułu (${brakujace.join(', ')}) nie zostały jeszcze zbudowane.`;
}

/** Ostrzega, gdy rdzeń odda okno bez kanału modelu: pierwsza wiadomość w takim
 * oknie skończy się odmową. */
function ostrzezOBrakuKanalu(pozycja: PozycjaModulu, okno: Window | undefined): void {
  if (okno === undefined || okno.id === '' || okno.modelChannelId !== '') return;
  pokazKomunikat({
    tytul: `Okno rozmowy modułu ${pozycja.nazwa} bez kanału modelu`,
    tresc: 'Rdzeń oddał okno przestrzeni roboczej bez kanału modelu — wysłanie wiadomości w tym oknie skończy się odmową.',
    waga: 'ostrz',
  });
}
