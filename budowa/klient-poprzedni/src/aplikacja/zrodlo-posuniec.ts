import {
  AssistantActionStatus,
  ChangeKind,
  EventType,
  MessageRole,
  type Window,
} from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import { utworzRozstrzyganieSprawcy } from './rozstrzyganie-sprawcy';
import { migawka, roznice, type MigawkaUstawien } from './ustawienia-okna-sledzone';
import type { Kanal } from '../protokol/kanal';

/**
 * Źródło posunięć — jedyne miejsce klienta odpowiadające na pytanie, czy dane
 * zdarzenie wywołało to połączenie, czy inne.
 *
 * Asystent steruje platformą przez osobne połączenie, a rdzeń rozgłasza każdą
 * zmianę do wszystkich połączeń konta (`transport/rozgloszenie.go`), więc ekran
 * dostaje komplet zdarzeń — także cudzych.
 *
 * Kontrakt nie niesie sprawcy: `WindowChangedEvent` niesie okno, a nie autora
 * zmiany, a `Message.role` mówi „user" niezależnie od tego, czy zdanie wpisano
 * ręcznie, czy przez MCP. Jedynym zdarzeniem z jawnym sprawcą jest
 * `session.focus.changed` niosące `clientId`. Resztę rozstrzyga rejestr
 * własnych odpowiedzi: zdarzenie o bycie spoza rejestru przyszło skądinąd.
 * Mechanizm rejestru i zwłoki mieszka w `rozstrzyganie-sprawcy.ts`; tutaj
 * zostaje to, co z werdyktu wynika — które posunięcia trafiają na pas, za czym
 * podąża nawigacja i które okno wchodzi na scenę.
 *
 * Źródło mówi „spoza tego połączenia", nigdy „asystent": drugie urządzenie
 * użytkownika i asystent są z tego miejsca nieodróżnialne.
 */

/** Rodzaj posunięcia — czym asystent ruszył w aplikacji. */
export type RodzajPosuniecia = 'okno' | 'modul' | 'prompt' | 'ognisko' | 'ustawienie';

/** Jedno posunięcie wykonane poza tym połączeniem. */
export interface Posuniecie {
  rodzaj: RodzajPosuniecia;
  /** Zdanie gotowe do postawienia na pasku. */
  opis: string;
  /**
   * Czy sprawca jest dowiedziony, czy wyprowadzony z braku własnej odpowiedzi.
   * `pewne` przysługuje wyłącznie zdarzeniu niosącemu `clientId`.
   */
  pewnosc: 'pewne' | 'domniemane';
  /** Okno, którego posunięcie dotyczy; puste, gdy posunięcie nie ma okna. */
  idOkna: string;
  /** Moduł po zmianie; pusty poza posunięciem `modul`. */
  kodModulu: string;
  chwila: number;
}

/** Stan pracy asystenta odczytany ze zleceń modułu Assistant. */
export interface PracaAsystenta {
  /** Nazwa zlecenia albo zdanie zastępcze, gdy rdzeń nazwy nie podał. */
  tytul: string;
  /** Etap bieżący i liczba etapów; oba zerowe, gdy rdzeń ich nie zna. */
  etap: number;
  etapow: number;
}

/** Port posunięć dla widoków powłoki. */
export interface ZrodloPosuniec {
  /** Posunięcie rozstrzygnięte jako pochodzące spoza tego połączenia. */
  naPosuniecie(sluchacz: (posuniecie: Posuniecie) => void): Odsubskrybuj;
  /**
   * Okno założone w tej sesji przez inne połączenie — do wprowadzenia na scenę.
   *
   * Okna zakładane przez sam interfejs tędy nie idą: wchodzą na scenę drogą
   * własnego zamówienia (`scena-sesji.ts`), a puszczone tu drugi raz stanęłyby
   * w drugim gnieździe.
   */
  naNoweOkno(sluchacz: (okno: Window) => void): Odsubskrybuj;
  /** Zmiana modułu okna zgłoszona przez rdzeń (`workspace.enter` gdziekolwiek). */
  naModulOkna(sluchacz: (okno: Window) => void): Odsubskrybuj;
  /** Praca asystenta w toku albo `null`, gdy nic nie biegnie. */
  naPrace(sluchacz: (praca: PracaAsystenta | null) => void): Odsubskrybuj;
  /** Ostatni znany stan pracy — widok montowany później nie zaczyna od pustki. */
  praca(): PracaAsystenta | null;
  rozlacz(): void;
}

export function utworzZrodloPosuniec(kanal: Kanal, idKlienta: string): ZrodloPosuniec {
  const posuniecia = utworzMagistrale<Posuniecie>();
  const noweOkna = utworzMagistrale<Window>();
  const zmianyModulu = utworzMagistrale<Window>();
  const prace = utworzMagistrale<PracaAsystenta | null>();
  const odsubskrybowania: Odsubskrybuj[] = [];

  /**
   * Ostatnio widziane ustawienia okna: moduł, kanał modelu, agent, rola,
   * zasięg wykonania, tryb uprawnień.
   *
   * Zdarzenie `window.changed` niesie okno po zmianie i nie mówi, co się w nim
   * zmieniło. Bez tej migawki pas meldowałby „okno zmienione" przy każdym
   * dotknięciu zamiast nazwać przestawione ustawienie.
   */
  const migawkaOkna = new Map<string, MigawkaUstawien>();
  let biezacaPraca: PracaAsystenta | null = null;

  // Całe rozstrzyganie „czyja to czynność" mieszka osobno: mechanizm jest
  // niezależny od tego, co z werdyktu wynika.
  const sprawca = utworzRozstrzyganieSprawcy();
  const { zapamietajWlasne, poZwloce, poZwloceOkno } = sprawca;

  function oglos(posuniecie: Omit<Posuniecie, 'chwila'>): void {
    posuniecia.oglos({ ...posuniecie, chwila: Date.now() });
  }

  // Zasiew migawki idzie tą samą drogą, co rejestr własności: odpowiedzi na
  // własne komendy niosą okno w stanie bieżącym i są jedynym źródłem stanu
  // wyjściowego, jakie ta warstwa ma. Bez zasiewu pierwsza cudza zmiana okna
  // przepada, bo nie ma z czym jej porównać.
  //
  // Zasiew, nie nadpisanie: stan z odpowiedzi wchodzi wyłącznie tam, gdzie
  // migawki jeszcze nie ma. Nadpisywanie kasowałoby stan zapisany przez
  // zdarzenie, które tę odpowiedź wyprzedziło, i różnica kolejnej zmiany
  // liczyłaby się od stanu nieaktualnego.
  odsubskrybowania.push(
    kanal.naDowolny((koperta) =>
      zapamietajWlasne(koperta, (okno) => {
        if (!migawkaOkna.has(okno.id)) migawkaOkna.set(okno.id, migawka(okno));
      }),
    ),
  );

  // Okna. Wszystko, co dotyczy okna, idzie przez rozstrzygnięcie o sprawcy —
  // także wprowadzenie okna na scenę. Okno zakładane przez sam interfejs
  // (`workspace.enter` przed uzgodnieniem) również przychodzi zdarzeniem
  // `created`, i to wcześniej niż odpowiedź, która je zamawiała; bez zwłoki
  // scena wprowadziłaby je drugi raz, do wolnego gniazda. Zwłoka niczego nie
  // wstrzymuje po stronie rdzenia.
  odsubskrybowania.push(
    kanal.naZdarzenie(EventType.WindowChanged, ({ change, window: okno }) => {
      if (okno === undefined || okno.id === '') return;

      if (change === ChangeKind.Created) {
        migawkaOkna.set(okno.id, migawka(okno));
        poZwloceOkno(okno, () => {
          noweOkna.oglos(okno);
          oglos({
            rodzaj: 'okno',
            opis: `Otwarte okno „${okno.title ?? 'bez nazwy'}" w module ${okno.moduleId}`,
            pewnosc: 'domniemane',
            idOkna: okno.id,
            kodModulu: okno.moduleId,
          });
        });
        return;
      }

      // Zdarzenie niesie okno po zmianie i nie mówi, co się zmieniło —
      // porównanie z migawką poprzednią jest jedyną drogą do nazwania
      // przestawionego ustawienia zamiast meldunku „okno zostało dotknięte".
      const poprzednia = migawkaOkna.get(okno.id);
      const biezaca = migawka(okno);
      migawkaOkna.set(okno.id, biezaca);
      if (poprzednia === undefined) return;

      const zmianaModulu = poprzednia.modul !== biezaca.modul;
      const zmianyUstawien = roznice(poprzednia, biezaca);
      if (!zmianaModulu && zmianyUstawien.length === 0) return;

      poZwloceOkno(okno, () => {
        if (zmianaModulu) {
          // Podążanie nawigacji rusza wyłącznie za zmianą cudzą; za własną nie
          // ma czego podążać, bo ekran już tam jest.
          zmianyModulu.oglos(okno);
          oglos({
            rodzaj: 'modul',
            opis: `Przejście do modułu ${biezaca.modul}`,
            pewnosc: 'domniemane',
            idOkna: okno.id,
            kodModulu: biezaca.modul,
          });
        }
        // Ustawienia meldowane są pojedynczo i po nazwie: zbiorcze „ustawienia
        // okna zmienione" nie mówi, co zostało przestawione i na co.
        for (const zmiana of zmianyUstawien) {
          oglos({
            rodzaj: 'ustawienie',
            opis: zmiana,
            pewnosc: 'domniemane',
            idOkna: okno.id,
            kodModulu: biezaca.modul,
          });
        }
      });
    }),
  );

  // Prompty. Pasek melduje sam prompt, nigdy odpowiedzi modelu: posunięcie
  // asystenta i praca modelu docelowego to dwie różne rzeczy, a zlanie ich
  // w jeden strumień zaciera, kto co zrobił.
  odsubskrybowania.push(
    kanal.naZdarzenie(EventType.MessageChanged, ({ message }) => {
      if (message === undefined || message.role !== MessageRole.User) return;
      poZwloce(message.id, () =>
        oglos({
          rodzaj: 'prompt',
          opis: `Wpisany prompt: „${skroc(message.content)}"`,
          pewnosc: 'domniemane',
          idOkna: message.windowId,
          kodModulu: '',
        }),
      );
    }),
  );

  // Ognisko — jedyne zdarzenie z dowiedzionym sprawcą.
  odsubskrybowania.push(
    kanal.naZdarzenie(EventType.SessionFocusChanged, (zmiana) => {
      if (zmiana.clientId === idKlienta || zmiana.clientId === '') return;
      oglos({
        rodzaj: 'ognisko',
        opis:
          zmiana.windowId === undefined || zmiana.windowId === ''
            ? 'Ognisko przeniesione na inną kartę sesji'
            : 'Ognisko przeniesione na inne okno',
        // `clientId` w treści zdarzenia jest dowodem, nie domysłem: to nie jest
        // ten klient, więc czynność na pewno wyszła z innego połączenia.
        pewnosc: 'pewne',
        idOkna: zmiana.windowId ?? '',
        kodModulu: '',
      });
    }),
  );

  // Praca asystenta. Jedyne zdarzenie kontraktu mówiące wprost, że asystent
  // pracuje, i mówi to o swoich zleceniach, nie o posunięciach w cudzych oknach.
  // Pasek trzyma więc dwie warstwy osobno: stan pracy bierze się stąd, a wykaz
  // posunięć ze zdarzeń obsługiwanych wyżej.
  odsubskrybowania.push(
    kanal.naZdarzenie(EventType.AssistantActionChanged, ({ action }) => {
      if (action === undefined) return;
      if (action.status !== AssistantActionStatus.Running) {
        // Zlecenie domknięte gasi wskaźnik wyłącznie wtedy, gdy to o nim pasek
        // mówił. Inaczej koniec zlecenia bocznego zgasiłby pracę wciąż trwającą.
        if (biezacaPraca?.tytul === tytulZlecenia(action.title)) ustawPrace(null);
        return;
      }
      ustawPrace({
        tytul: tytulZlecenia(action.title),
        etap: action.currentStep ?? 0,
        etapow: action.totalSteps ?? 0,
      });
    }),
  );

  function ustawPrace(praca: PracaAsystenta | null): void {
    biezacaPraca = praca;
    prace.oglos(praca);
  }

  return {
    naPosuniecie: (sluchacz) => posuniecia.subskrybuj(sluchacz),
    naNoweOkno: (sluchacz) => noweOkna.subskrybuj(sluchacz),
    naModulOkna: (sluchacz) => zmianyModulu.subskrybuj(sluchacz),
    naPrace: (sluchacz) => prace.subskrybuj(sluchacz),
    praca: () => biezacaPraca,
    rozlacz() {
      sprawca.rozlacz();
      for (const odsubskrybuj of odsubskrybowania.splice(0)) odsubskrybuj();
    },
  };
}

/** Nazwa zlecenia albo zdanie zastępcze — pasek nie pokazuje pustego miejsca. */
function tytulZlecenia(tytul: string | undefined): string {
  const nazwa = (tytul ?? '').trim();
  return nazwa.length > 0 ? nazwa : 'zlecenie bez nazwy';
}

/** Skrót promptu na pasek; pasek ma jeden wiersz, nie akapit. */
function skroc(tresc: string): string {
  const jednym = tresc.replace(/\s+/g, ' ').trim();
  return jednym.length <= 80 ? jednym : `${jednym.slice(0, 79)}…`;
}
