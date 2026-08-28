import { Command } from '../../../shared/contract';
import type { WidokModulu } from '../aplikacja/rejestr-modulow';
import type { FazaOkna } from '../komponenty/faza-okna';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Kanal } from '../protokol/kanal';

// Typy WidokModulu i OpisModulu należą do powłoki — to ona określa, czego
// wymaga od modułu, ten plik je przekazuje dalej.
export type { OpisModulu, WidokModulu } from '../aplikacja/rejestr-modulow';

/**
 * Przejście z modułu montującego się samodzielnie na widok powłoki, z
 * gospodarzem własnym zwracanym jako element.
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
 * Przejście dla modułu, który pracuje w oknie, nie w sesji: odracza montaż do
 * chwili, gdy okno tego modułu jest znane rdzeniowi.
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
 * Stan modułu przed montażem, powiedziany wprost jednym zdaniem w obszarze
 * roboczym, bez klasy okna operacyjnego.
 */
function powiedz(gospodarz: HTMLElement, faza: FazaOkna, zdanie: string): void {
  const akapit = document.createElement('p');
  akapit.dataset['faza'] = faza;
  akapit.textContent = zdanie;
  gospodarz.replaceChildren(akapit);
}

/** Czego dokładnie brakuje: okna tego modułu czy okien w sesji w ogóle, nazwane zdaniem dla operatora okna. */
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
 * Przejście dla modułu, który potrzebuje sesji już przy montażu, z montażem
 * odroczonym do jej wczytania.
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
