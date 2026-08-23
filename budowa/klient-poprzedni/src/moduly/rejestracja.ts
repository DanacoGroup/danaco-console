import { Command } from '../../../shared/contract';
import type { WidokModulu } from '../aplikacja/rejestr-modulow';
import type { FazaOkna } from '../komponenty/faza-okna';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Kanal } from '../protokol/kanal';

// Typy WidokModulu i OpisModulu należą do powłoki — to ona określa, czego
// wymaga od modułu. Ten plik ich nie powtarza, tylko przekazuje dalej, żeby
// moduły nie musiały sięgać do katalogu aplikacji po kontrakt, którego używają.
// Import typu nie tworzy zależności czasu wykonania, więc krąg nie powstaje.
export type { OpisModulu, WidokModulu } from '../aplikacja/rejestr-modulow';

/**
 * Przejście z modułu montującego się samodzielnie na widok powłoki.
 *
 * Część modułów montuje się do gospodarza podanego z zewnątrz
 * (`zamontujX(gospodarz, …)`), a nie oddaje własnego elementu. Przejście daje
 * im gospodarza własnego i zwraca go jako element — zachowanie modułu zostaje
 * bez zmian, a powłoka dostaje kształt `WidokModulu`, którego wymaga.
 *
 * `wczytaj` jest tu bezczynne: te moduły odczytują rdzeń już przy montażu,
 * więc powtórny odczyt przy wejściu do przestrzeni byłby drugim zapytaniem
 * o to samo.
 */
export function widokZMontazu(
  zamontuj: (gospodarz: HTMLElement, kanal: Kanal) => { element: HTMLElement; zamknij(): void },
  kanal: Kanal,
): WidokModulu {
  const gospodarz = document.createElement('div');
  const zamontowane = zamontuj(gospodarz, kanal);
  return {
    element: gospodarz,
    wczytaj: async () => {},
    zamknij: () => zamontowane.zamknij(),
  };
}

/**
 * Przejście dla modułu, który pracuje w oknie, nie w sesji.
 *
 * `WidokModulu.wczytaj` dostaje identyfikator sesji, a moduł pracujący
 * w konkretnym oknie komunikacji nie ma bez niego czego otworzyć. Dopóki umowa
 * nie niesie okna, moduł odracza montaż do chwili wczytania i sam pyta rdzeń
 * o okna sesji.
 *
 * Każde wyjście bez montażu nazywa swój stan zdaniem w obszarze roboczym.
 * Ciche `return` zostawiałoby obszar o zerowej treści, nie do odróżnienia od
 * sesji bez okien i od usterki widoku.
 *
 * Okno wybiera się po `Window.moduleId` (pole ustawiane przez rdzeń z kodu
 * modułu w `server/internal/core/przeklad_nawigacja.go`), a nie po pierwszej
 * pozycji wykazu: sesja z oknem cudzego modułu na początku dałaby modułowi
 * okno nie jego, a komendy pojechałyby do rdzenia z obcym `windowId`.
 *
 * `kodModulu` jest opcjonalny ze względu na granicę własności: wywołania,
 * które kodu nie podają, zachowują wybór pierwszego okna sesji. Jest to jawny
 * dług — takie wywołania należy uzupełnić o kod modułu.
 */
export function widokZOknaSesji(
  kanal: Kanal,
  zamontuj: (gospodarz: HTMLElement, kanal: Kanal, okno: string) => { zamknij(): void },
  kodModulu?: string,
): WidokModulu {
  const gospodarz = document.createElement('div');
  let zamontowane: { zamknij(): void } | null = null;

  return {
    element: gospodarz,
    async wczytaj(idSesji: string) {
      if (zamontowane !== null) return;
      powiedz(gospodarz, 'ladowanie', 'Odczyt okien sesji…');
      kanal.wyslij(Command.WindowList, { sessionId: idSesji }, (wynik) => {
        if (!wynik.udany) {
          powiedz(
            gospodarz,
            'blad',
            opisOdmowyBledu('Odczyt okien sesji nie udał się', wynik.blad),
          );
          return;
        }
        const okna = wynik.wynik?.windows ?? [];
        const okno =
          kodModulu === undefined ? okna[0] : okna.find((w) => w.moduleId === kodModulu);
        if (okno === undefined) {
          powiedz(gospodarz, 'puste', zdanieBrakuOkna(kodModulu, okna.length));
          return;
        }
        gospodarz.replaceChildren();
        zamontowane = zamontuj(gospodarz, kanal, okno.id);
      });
    },
    zamknij: () => zamontowane?.zamknij(),
  };
}

/**
 * Stan modułu przed montażem, powiedziany wprost.
 *
 * Bez klasy `dn-*`: to nie jest okno operacyjne ani element pakietu
 * design, tylko jedno zdanie w obszarze roboczym na czas przed montażem albo
 * zamiast niego. `data-faza` bierze wartości ze wspólnego słownika
 * `komponenty/faza-okna`, żeby sprawdzian i powłoka rozpoznawały ten stan tak
 * samo jak stany okien.
 */
function powiedz(gospodarz: HTMLElement, faza: FazaOkna, zdanie: string): void {
  const akapit = document.createElement('p');
  akapit.dataset['faza'] = faza;
  akapit.textContent = zdanie;
  gospodarz.replaceChildren(akapit);
}

/** Czego dokładnie brakuje: okna tego modułu czy okien w sesji w ogóle. */
function zdanieBrakuOkna(kodModulu: string | undefined, ileOkien: number): string {
  if (ileOkien === 0) {
    return 'Ta sesja nie ma jeszcze ani jednego okna — moduł otworzy się razem z pierwszym oknem.';
  }
  if (kodModulu === undefined) {
    return 'Rdzeń oddał okna sesji, ale żadnego nie dało się użyć do montażu tego modułu.';
  }
  return `Ta sesja ma ${ileOkien} okno/okien, ale żadne nie należy do modułu ${kodModulu} — moduł otworzy się razem ze swoim oknem.`;
}

/**
 * Przejście dla modułu, który potrzebuje sesji już przy montażu.
 *
 * Odmiana `widokZOknaSesji`: tamten pyta rdzeń o okna, ten wystarcza sobie samą
 * sesją. Montaż jest odroczony do wczytania, bo w chwili tworzenia widoku sesja
 * bywa jeszcze nieznana — powłoka odkłada wtedy wejście i ponawia je po otwarciu
 * okna przez rdzeń.
 */
export function widokZSesji(
  kanal: Kanal,
  zamontuj: (gospodarz: HTMLElement, kanal: Kanal, sesja: string) => { zamknij(): void },
): WidokModulu {
  const gospodarz = document.createElement('div');
  let zamontowane: { zamknij(): void } | null = null;

  return {
    element: gospodarz,
    async wczytaj(idSesji: string) {
      if (zamontowane !== null || idSesji === '') return;
      zamontowane = zamontuj(gospodarz, kanal, idSesji);
    },
    zamknij: () => zamontowane?.zamknij(),
  };
}
