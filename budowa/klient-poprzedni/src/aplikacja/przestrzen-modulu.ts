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

/** Przestrzeń robocza modułu — co widać po wyborze pozycji bocznej nawigacji, przełożone na obszar roboczy. */
export interface ZaleznosciPrzestrzeni {
  // Kanał kontraktu — droga `workspace.enter` oraz kanał podawany widokom modułów.
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

  // Scena wchodzi na planszę raz i zostaje do końca życia powłoki, niesie gniazdo WebSocket.
  obszar.osadzCzat(scena);

  /** Pozycja czekająca na sesję — wejście przed uzgodnieniem nie przepada. */
  let odlozona: PozycjaModulu | null = null;
  // Numer bieżącego wyboru pozycji — odpowiedź na wybór porzucony nie ma prawa wejść na planszę.
  let zeton = 0;
  // Widoki modułów powstają raz na moduł i zostają — mapa żyje tyle, co karta przeglądarki.
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

  // Widok sekcji panelu orkiestracji, z tej samej mapy co widoki modułów — ten sam cykl życia.
  function widokSekcji(klucz: string): WidokModulu | undefined {
    if (!czySekcjaOrkiestracji(klucz)) return undefined;
    const gotowy = widoki.get(klucz);
    if (gotowy !== undefined) return gotowy;
    const widok = utworzWidokSekcji(klucz, kanal);
    widoki.set(klucz, widok);
    return widok;
  }

  // Wejście do przestrzeni roboczej, a dopiero po nim odczyt widoku modułu — kolejność jest tu treścią.
  async function wejdzIWczytaj(pozycja: PozycjaModulu, widok?: WidokModulu): Promise<void> {
    const moj = zeton;
    const sesja = idSesji();
    if (sesja === '') {
      odlozona = pozycja;
      return;
    }
    odlozona = null;

    // Sekcja panelu orkiestracji nie wchodzi w przestrzeń modułu — czyta wprost swoją treść.
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
    // Pozycja zmieniona w trakcie oczekiwania — to już nie odpowiedź na pytanie z planszy.
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

      // Stan pusty trzymany w zgodzie z pozycją nawet ukryty — powrót nie może pokazać poprzedniej treści.
      obszar.zapowiedz(pozycja, dane, pozycja.okna);

      if (pozycja.modul === undefined) {
        // Sekcja bez modułu odwołuje wejście odłożone — nie wchodzi w moduł już opuszczony.
        odlozona = null;
        obszar.nota(null);

        // Sekcja znana panelowi ma swój widok i wchodzi na planszę obok sceny; nieznana zostaje pusta.
        const sekcja = widokSekcji(pozycja.klucz);
        if (sekcja === undefined) {
          obszar.oproznij();
          return;
        }
        obszar.pokaz(sekcja.element);
        void wejdzIWczytaj(pozycja, sekcja);
        return;
      }

      // Rejestr kluczuje po kodzie modułu (`modul.code` rdzenia), nie po identyfikatorze wiersza.
      const widok = widokModulu(pozycja.klucz);
      if (widok === undefined) {
        // Moduł bez zbudowanego widoku: sama scena, a pasek uczciwości mówi, czego jeszcze nie ma.
        obszar.nota(notaModulu(pozycja));
        obszar.pokaz(null);
        void wejdzIWczytaj(pozycja);
        return;
      }

      // Moduł zbudowany: widok modułu obok sceny, nie zamiast niej — rozmowa jest oknem wiodącym.
      obszar.nota(null);
      obszar.pokaz(widok.element);
      void wejdzIWczytaj(pozycja, widok);
    },

    ponow() {
      // Wejście odłożone niesie ze sobą odczyt widoku, inaczej moduł wszedłby w przestrzeń bez treści.
      const pozycja = odlozona;
      if (pozycja === null) return;
      // Odłożona bywa też sekcja panelu — jej widok bierze się z panelu, nie z rejestru modułów.
      void wejdzIWczytaj(
        pozycja,
        pozycja.modul === undefined ? widokSekcji(pozycja.klucz) : widokModulu(pozycja.klucz),
      );
    },
  };
}

/** Czy kod okna wskazuje okno rozmowy — w obu postaciach, przedrostkowanej i bez przedrostka, jakie niesie rdzeń. */
function czyOknoRozmowy(kod: string): boolean {
  return kod === 'chat-window' || kod.endsWith('.chat-window');
}

/** Zdanie paska uczciwości: co w tym module działa, a jakich okien operacyjnych jeszcze dziś w nim nie ma. */
function notaModulu(pozycja: PozycjaModulu): string {
  const brakujace = pozycja.okna.filter((kod) => !czyOknoRozmowy(kod));
  if (brakujace.length === 0) {
    return `${pozycja.nazwa}: rdzeń nie podaje dla tego modułu okien operacyjnych poza oknem rozmowy.`;
  }
  return `${pozycja.nazwa}: działa okno rozmowy. Okna operacyjne tego modułu (${brakujace.join(', ')}) nie zostały jeszcze zbudowane.`;
}

/** Ostrzega, gdy rdzeń odda okno bez kanału modelu — pierwsza wiadomość w takim oknie skończy się odmową. */
function ostrzezOBrakuKanalu(pozycja: PozycjaModulu, okno: Window | undefined): void {
  if (okno === undefined || okno.id === '' || okno.modelChannelId !== '') return;
  pokazKomunikat({
    tytul: `Okno rozmowy modułu ${pozycja.nazwa} bez kanału modelu`,
    tresc: 'Rdzeń oddał okno przestrzeni roboczej bez kanału modelu — wysłanie wiadomości w tym oknie skończy się odmową.',
    waga: 'ostrz',
  });
}
