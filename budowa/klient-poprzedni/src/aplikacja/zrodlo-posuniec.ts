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

/** Rodzaj posunięcia — czym dokładnie asystent ruszył w aplikacji: okno, moduł, prompt, ognisko albo ustawienie. */
export type RodzajPosuniecia = 'okno' | 'modul' | 'prompt' | 'ognisko' | 'ustawienie';

/** Jedno posunięcie wykonane poza tym połączeniem, gotowe do pokazania na pasku posunięć asystenta klienta. */
export interface Posuniecie {
  rodzaj: RodzajPosuniecia;
  /** Zdanie gotowe do postawienia na pasku. */
  opis: string;
  // Czy sprawca dowiedziony, czy wyprowadzony z braku odpowiedzi; `pewne` wymaga `clientId`.
  pewnosc: 'pewne' | 'domniemane';
  /** Okno, którego posunięcie dotyczy; puste, gdy posunięcie nie ma okna. */
  idOkna: string;
  /** Moduł po zmianie; pusty poza posunięciem `modul`. */
  kodModulu: string;
  chwila: number;
}

/** Stan pracy asystenta odczytany ze zleceń modułu Assistant — tytuł zlecenia oraz jego bieżący etap pracy. */
export interface PracaAsystenta {
  /** Nazwa zlecenia albo zdanie zastępcze, gdy rdzeń nazwy nie podał. */
  tytul: string;
  /** Etap bieżący i liczba etapów; oba zerowe, gdy rdzeń ich nie zna. */
  etap: number;
  etapow: number;
}

/** Port posunięć dla widoków powłoki: subskrypcje posunięć, nowych okien, modułu i pracy samego asystenta. */
export interface ZrodloPosuniec {
  /** Posunięcie rozstrzygnięte jako pochodzące spoza tego połączenia. */
  naPosuniecie(sluchacz: (posuniecie: Posuniecie) => void): Odsubskrybuj;
  // Okno założone w tej sesji przez inne połączenie — własne okna tędy nie idą.
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

  // Ostatnio widziane ustawienia okna — bez migawki pas meldowałby ogólne „okno zmienione".
  const migawkaOkna = new Map<string, MigawkaUstawien>();
  let biezacaPraca: PracaAsystenta | null = null;

  // Całe rozstrzyganie „czyja to czynność" mieszka osobno, niezależnie od tego, co z werdyktu wynika.
  const sprawca = utworzRozstrzyganieSprawcy();
  const { zapamietajWlasne, poZwloce, poZwloceOkno } = sprawca;

  function oglos(posuniecie: Omit<Posuniecie, 'chwila'>): void {
    posuniecia.oglos({ ...posuniecie, chwila: Date.now() });
  }

  // Zasiew migawki idzie tą samą drogą, co rejestr własności — bez niego pierwsza cudza zmiana przepada.

  // Zasiew, nie nadpisanie — stan z odpowiedzi wchodzi wyłącznie tam, gdzie migawki jeszcze nie ma.
  odsubskrybowania.push(
    kanal.naDowolny((koperta) =>
      zapamietajWlasne(koperta, (okno) => {
        if (!migawkaOkna.has(okno.id)) migawkaOkna.set(okno.id, migawka(okno));
      }),
    ),
  );

  // Okna — wszystko, co ich dotyczy, idzie przez rozstrzygnięcie o sprawcy, także wprowadzenie na scenę.
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

      // Zdarzenie niesie okno po zmianie — porównanie z migawką jest jedyną drogą do nazwania ustawienia.
      const poprzednia = migawkaOkna.get(okno.id);
      const biezaca = migawka(okno);
      migawkaOkna.set(okno.id, biezaca);
      if (poprzednia === undefined) return;

      const zmianaModulu = poprzednia.modul !== biezaca.modul;
      const zmianyUstawien = roznice(poprzednia, biezaca);
      if (!zmianaModulu && zmianyUstawien.length === 0) return;

      poZwloceOkno(okno, () => {
        if (zmianaModulu) {
          // Podążanie nawigacji rusza wyłącznie za zmianą cudzą — za własną nie ma czego podążać.
          zmianyModulu.oglos(okno);
          oglos({
            rodzaj: 'modul',
            opis: `Przejście do modułu ${biezaca.modul}`,
            pewnosc: 'domniemane',
            idOkna: okno.id,
            kodModulu: biezaca.modul,
          });
        }
        // Ustawienia meldowane są pojedynczo i po nazwie, nie zbiorczym „ustawienia okna zmienione".
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

  // Prompty — pasek melduje sam prompt, nigdy odpowiedzi modelu, bo to dwie różne rzeczy.
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
        // `clientId` w treści zdarzenia jest dowodem, nie domysłem, że wyszła z innego połączenia.
        pewnosc: 'pewne',
        idOkna: zmiana.windowId ?? '',
        kodModulu: '',
      });
    }),
  );

  // Praca asystenta — jedyne zdarzenie kontraktu mówiące wprost, że asystent pracuje.
  odsubskrybowania.push(
    kanal.naZdarzenie(EventType.AssistantActionChanged, ({ action }) => {
      if (action === undefined) return;
      if (action.status !== AssistantActionStatus.Running) {
        // Zlecenie domknięte gasi wskaźnik wyłącznie wtedy, gdy to o nim pasek mówił.
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

/** Nazwa zlecenia albo zdanie zastępcze, gdy rdzeń nazwy nie podał — pasek nie pokazuje pustego miejsca. */
function tytulZlecenia(tytul: string | undefined): string {
  const nazwa = (tytul ?? '').trim();
  return nazwa.length > 0 ? nazwa : 'zlecenie bez nazwy';
}

/** Skrót promptu na pasek posunięć — pasek ma jeden wiersz, nie akapit, dłuższa treść kończy się wielokropkiem. */
function skroc(tresc: string): string {
  const jednym = tresc.replace(/\s+/g, ' ').trim();
  return jednym.length <= 80 ? jednym : `${jednym.slice(0, 79)}…`;
}
