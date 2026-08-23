import {
  ConfigScope,
  MemoryEntryOrigin,
  type WorkspaceMemoryEntry,
} from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleTresci,
  przelacznik,
  przyciskAkcji,
  utworzWierszOdpowiedzi,
  wiersz,
  wybor,
  wykaz,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { BRAKI } from './braki-kontraktu';
import { NAZWY_ZASIEGOW, ODCZYTY, POZIOMY_PAMIECI, PUSTE } from './etykiety-assistant';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAssistant } from './stan-assistant';
import { wierszPamieci } from './wiersz-pamieci';
import type { ZrodloPamieci } from './zrodlo-pamieci';

/**
 * Zakładka faktów Memory & Context Manager — jawny edytor pamięci asystenta.
 *
 * Realizuje zasadę jawności modułu: pamięć jest w całości widoczna, edytowalna
 * i usuwalna. Cztery czynności kontraktu wystarczają: `memory.list` czyta,
 * `memory.set` zakłada i zmienia (także przypięcie i pochodzenie),
 * `memory.delete` kasuje.
 *
 * Zasięgiem odczytu jest karta sesji, nie projekt. Moduł Assistant pracuje
 * w karcie sesji środowiska TalkIn i nie ma pojęcia projektu — `memory.list`
 * przyjmuje `sessionId` właśnie na taki przypadek („karta sesji, gdy projekt
 * nie został wskazany"). Zapis idzie tą samą drogą, z zasięgiem wskazanym
 * jawnie przez Operatora.
 *
 * Potwierdzenie mówi to, co zapisał rdzeń, a nie to, co wysłało okno: wpis
 * zapisany na poziomie szerszym niż karta sesji bywa niewidoczny w wykazie
 * poniżej, więc samo odświeżenie listy niczego by nie potwierdzało.
 *
 * Reguł retencji, wygaszania (TTL) i znaczników wrażliwości okno nie udaje.
 * `WorkspaceMemoryEntry` ma osiem pól i żadne z nich nie niesie czasu życia ani
 * wrażliwości; brak nazywa przycisk, zamiast stawiać formularz, którego rdzeń
 * nie zapisze.
 */
export interface PanelFaktow {
  element: HTMLElement;
  /** Odczyt pamięci widocznej w zasięgu karty sesji. */
  wczytaj(): Promise<void>;
}

export function utworzPanelFaktow(stan: StanAssistant, zrodlo: ZrodloPamieci): PanelFaktow {
  const okno: StanOkna = utworzStanOkna();
  const odpowiedz: WierszOdpowiedzi = utworzWierszOdpowiedzi();

  const edytor = poleTresci('Treść ustalenia', 3, 'co asystent ma zapamiętać');
  const zasieg = wybor(
    'Zasięg współdzielenia wpisu',
    POZIOMY_PAMIECI.map((poziom) => [poziom, `Zasięg wpisu: ${NAZWY_ZASIEGOW[poziom]}`] as const),
  );
  zasieg.value = ConfigScope.Session;
  const przypiecie = przelacznik('Przypnij wpis');
  const wspolne = przelacznik('Pokaż ustalenia z poziomów szerszych');

  const lista = wykaz('Wpisy pamięci asystenta', 'ma-wykaz');
  okno.tresc.append(lista, odpowiedz.element);

  /** Wpis wczytany do edytora; pusty napis znaczy „zapis zakłada wpis nowy". */
  let zmieniany = '';

  const formularz = document.createElement('div');
  formularz.className = 'ma-formularz';
  formularz.append(
    wiersz('Ustalenie', edytor, {
      klasa: 'ma-wiersz',
      objasnienie:
        'Treść wchodzi do pamięci asystenta komendą memory.set. Pusty edytor zakłada wpis ' +
        'nowy; „Edytuj" przy pozycji wczytuje wpis istniejący i zapis go zmienia.',
    }),
    wiersz('Zasięg', zasieg, {
      klasa: 'ma-wiersz',
      objasnienie:
        'Poziom, na którym wpis jest współdzielony. Karta sesji trzyma ustalenie przy tej ' +
        'rozmowie; poziom szerszy czyni z niego ustalenie wspólne.',
    }),
    wiersz('Przypnij wpis', przypiecie, {
      klasa: 'ma-wiersz',
      objasnienie: 'Przypięte stoją na początku wykazu oddawanego przez rdzeń.',
    }),
    wiersz('Pokaż ustalenia wspólne', wspolne, {
      klasa: 'ma-wiersz',
      objasnienie: 'Dokłada do wykazu wpisy zapisane na poziomach szerszych niż karta sesji.',
    }),
    przyciski(),
  );

  const element = document.createElement('div');
  element.className = 'ma-obszar';
  element.dataset['obszar'] = 'fakty';
  element.append(formularz, okno.element);

  function przyciski(): HTMLElement {
    const zapisz = przyciskAkcji('Zapisz ustalenie', 'dn-btn dn-btn--sm dn-btn--atrament');
    zapisz.addEventListener('click', () =>
      zapiszWpis({
        tresc: edytor.value,
        zasieg: zasieg.value as ConfigScope,
        przypiety: przypiecie.checked,
        pochodzenie: MemoryEntryOrigin.Operator,
        wpis: zmieniany,
      }),
    );

    const odswiez = przyciskAkcji('Odczytaj pamięć', 'dn-btn dn-btn--sm dn-btn--zarys');
    odswiez.addEventListener('click', () => void wczytaj());

    const rzad = document.createElement('div');
    rzad.className = 'ma-formularz__przyciski';
    rzad.append(zapisz, odswiez);
    return rzad;
  }

  function zapiszWpis(zamowienie: {
    tresc: string;
    zasieg: ConfigScope;
    przypiety: boolean;
    pochodzenie: MemoryEntryOrigin;
    wpis: string;
  }): void {
    if (zamowienie.tresc.trim() === '') {
      // Pusta treść jest brakiem w polu formularza, nie odmową rdzenia: wykaz
      // zostaje nietknięty, a zdanie stoi przy formularzu.
      odpowiedz.pokaz('Wpisz treść ustalenia — rdzeń odmówi zapisu pustego.', false);
      return;
    }
    void zrodlo
      .zapisz({
        content: zamowienie.tresc,
        entryId: zamowienie.wpis,
        scope: zamowienie.zasieg,
        pinned: zamowienie.przypiety,
        origin: zamowienie.pochodzenie,
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          odpowiedz.pokaz(
            opisOdmowy('Zapis ustalenia', wynik.blad?.code, wynik.blad?.message),
            false,
          );
          return;
        }
        const zapisany = wynik.wynik.entry;
        edytor.value = '';
        zmieniany = '';
        odpowiedz.pokaz(
          `Rdzeń zapisał wpis ${zapisany.id} — zasięg ${NAZWY_ZASIEGOW[zapisany.scope]}, ` +
            `${zapisany.pinned === true ? 'przypięty' : 'nieprzypięty'}.`,
          true,
        );
        void wczytaj();
      });
  }

  function usunWpis(wpis: WorkspaceMemoryEntry): void {
    void zrodlo.usun(wpis.id).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        odpowiedz.pokaz(
          opisOdmowy('Usunięcie wpisu pamięci', wynik.blad?.code, wynik.blad?.message),
          false,
        );
        return;
      }
      if (!wynik.wynik.deleted) {
        // Rdzeń przyjął wywołanie i odpowiedział „nie usunąłem". Zdanie mówi to,
        // co oddał, a nie to, o co okno prosiło.
        odpowiedz.pokaz('Rdzeń przyjął wywołanie, ale nie potwierdził usunięcia wpisu.', false);
        return;
      }
      if (zmieniany === wpis.id) {
        zmieniany = '';
        edytor.value = '';
      }
      odpowiedz.pokaz(`Rdzeń usunął wpis ${wpis.id} z pamięci asystenta.`, true);
      void wczytaj();
    });
  }

  function pozycja(wpis: WorkspaceMemoryEntry): HTMLElement {
    return wierszPamieci(wpis, {
      zapisz: (zamowienie) => zapiszWpis(zamowienie),
      usun: usunWpis,
      wczytajDoEdytora(zrodlowy) {
        zmieniany = zrodlowy.id;
        edytor.value = zrodlowy.content;
        zasieg.value = zrodlowy.scope;
        przypiecie.checked = zrodlowy.pinned === true;
        edytor.focus();
        odpowiedz.pokaz('Wpis wczytany do edytora — „Zapisz ustalenie" zmieni go w rdzeniu.', true);
      },
    });
  }

  async function wczytaj(): Promise<void> {
    const sesja = stan.idSesji();
    if (sesja === '') {
      okno.blad(BRAKI.brakSesji);
      return;
    }
    okno.ladowanie(ODCZYTY.pamiec);
    const wynik = await zrodlo.wpisy({ sessionId: sesja, includeShared: wspolne.checked });
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Odczyt pamięci asystenta', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    const wpisy = wynik.wynik.entries;
    lista.replaceChildren(...wpisy.map(pozycja));
    if (wpisy.length === 0) {
      okno.puste(PUSTE.pamiec);
      return;
    }
    okno.gotowe();
  }

  wspolne.addEventListener('change', () => void wczytaj());

  okno.puste(PUSTE.pamiecSpoczynek);
  return { element, wczytaj };
}
