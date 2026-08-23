import {
  ConfigScope,
  MemoryEntryOrigin,
  type WorkspaceMemoryEntry,
} from '../../../../shared/contract';
import {
  przelacznik,
  przyciskAkcji as przycisk,
  poleTresci,
  wiersz,
  wykaz,
  pozycjaWykazu,
} from '../../modele/kontrolki-formularza';
import { POZIOMY_ZASIEGU, type ParaPoziomu } from './poziomy-zasiegu';
import type { StanProjektu } from './stan-projektu';
import { utworzStanTresci } from './stany-okna';
import type { ZrodloWorkspace } from './zrodlo-workspace';

/**
 * Pamięć karty sesji — poziomy włączone, zapis pamięci i wpisy widoczne w sesji.
 *
 * Odpowiedzialność odrębna od pamięci projektu. `pamiec-kontekstu.ts` prowadzi
 * pamięć projektu rodziną `workspace.context.*`: wpisy jednego projektu, wspólne
 * wszystkim jego kartom. Ten panel prowadzi to, co widzi karta — rodziną
 * `memory.*`, w której zasięg jest polem żądania, a nie założeniem. Dwie rodziny
 * komend, dwa byty, dwa pliki.
 *
 * Panel woła `memory.list`, `memory.set` i `memory.toggle`. `memory.detach`
 * zostaje bez wołacza: znaczenie odpięcia nie jest w kontrakcie ustalone, więc
 * kontrolka wołająca tę komendę zmieniałaby zawartość bazy w sposób, którego
 * okno nie potrafi nazwać.
 *
 * Każde zdanie mówi to, co oddał rdzeń: poziomy po przestawieniu biorą się
 * z odpowiedzi `memory.toggle`, nie ze stanu przełączników — rdzeń, który
 * przyjmie wywołanie i odda inny zestaw poziomów, nie zrobił tego, o co proszono.
 */
export interface PanelPamieciSesji {
  element: HTMLElement;
  /** Odczytuje wpisy widoczne w karcie sesji. */
  odswiez(): void;
}

/**
 * Poziomy zasięgu do przestawienia — ten sam wykaz, co w pamięci projektu
 * i w Instructions Panel. Jedno źródło nazw stoi w `poziomy-zasiegu.ts`, żeby
 * ta sama nastawa nie nazywała się w trzech oknach na trzy sposoby.
 */
const POZIOMY: readonly ParaPoziomu[] = POZIOMY_ZASIEGU;

export function utworzPanelPamieciSesji(
  zrodlo: ZrodloWorkspace,
  stan: StanProjektu,
): PanelPamieciSesji {
  const tresc = utworzStanTresci();
  const przelaczniki = new Map<ConfigScope, HTMLInputElement>();

  const siatka = document.createElement('div');
  siatka.className = 'dw-poziomy';
  for (const [poziom, opis] of POZIOMY) {
    const kontrolka = przelacznik(`Poziom pamięci: ${opis}`);
    przelaczniki.set(poziom, kontrolka);
    siatka.append(wiersz(opis, kontrolka, { klasa: 'dw-wiersz' }));
  }

  const zapisCzynny = przelacznik('Zapis pamięci czynny');
  const przestaw = przycisk('Przestaw pamięć karty', 'dn-btn dn-btn--atrament');
  const odczyt = przycisk('Odśwież wpisy karty');
  const ustalenie = poleTresci('Ustalenie karty sesji', 3);
  const zapisz = przycisk('Zapisz ustalenie karty');

  const element = document.createElement('section');
  element.className = 'dw-pamiec-sesji';
  element.setAttribute('aria-label', 'Pamięć karty sesji');
  element.append(
    siatka,
    wiersz('Zapis pamięci', zapisCzynny, {
      klasa: 'dw-wiersz',
      objasnienie: 'Wyłączony zapis zostawia pamięć do czytania; rdzeń nie dopisze nic nowego.',
    }),
    wiersz('Ustalenie karty', ustalenie, {
      klasa: 'dw-wiersz',
      objasnienie: 'Wpis zapisany na poziomie karty sesji, nie w pamięci projektu.',
    }),
    paskiem(przestaw, odczyt, zapisz),
    tresc.element,
  );

  /**
   * Warunek wstępny wszystkich trzech komend panelu: bez karty sesji nie ma
   * czego wysłać. Zwraca `true` i wypisuje powód, gdy karty jeszcze nie ma.
   */
  function bezKarty(): boolean {
    if (stan.sesja() !== '') return false;
    tresc.pusto(
      'Moduł nie zna jeszcze karty sesji — pamięć karty da się przestawić dopiero po jej wczytaniu przez powłokę.',
    );
    return true;
  }

  /** Poziomy zaznaczone przełącznikami, w kolejności wykazu. */
  function zaznaczone(): ConfigScope[] {
    return POZIOMY.filter(([poziom]) => przelaczniki.get(poziom)?.checked === true).map(
      ([poziom]) => poziom,
    );
  }

  przestaw.addEventListener('click', () => {
    if (bezKarty()) return;
    const poziomy = zaznaczone();
    if (poziomy.length === 0) {
      // Pusta lista poziomów zostawia stan bez zmian, więc jej wysłanie dałoby
      // potwierdzenie udanego wywołania, po którym nic się nie zmieniło.
      tresc.potwierdzenie(
        'Pusty zestaw poziomów zostawia stan bez zmian — zaznacz przynajmniej jeden poziom.',
        false,
      );
      return;
    }
    tresc.ladowanie('Przestawianie poziomów pamięci karty…');
    void zrodlo
      .przestawPamiec({ sessionId: stan.sesja(), levels: poziomy, writeEnabled: zapisCzynny.checked })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Rdzeń nie przestawił pamięci karty.', wynik.blad);
          return;
        }
        const oddane = wynik.wynik;
        przyjmijPoziomy(oddane.levels);
        zapisCzynny.checked = oddane.writeEnabled;
        tresc.potwierdzenie(
          `Rdzeń ustawił poziomy: ${oddane.levels.join(', ') || 'brak'} — zapis pamięci ${oddane.writeEnabled ? 'czynny' : 'wstrzymany'}.`,
          true,
        );
      });
  });

  zapisz.addEventListener('click', () => {
    if (bezKarty()) return;
    if (ustalenie.value.trim() === '') {
      tresc.potwierdzenie('Ustalenie bez treści nie zostanie zapisane.', false);
      return;
    }
    void zrodlo
      .zapiszPamiec({
        content: ustalenie.value,
        scope: ConfigScope.Session,
        scopeId: stan.sesja(),
        origin: MemoryEntryOrigin.Operator,
        ...(stan.projekt() === '' ? {} : { projectId: stan.projekt() }),
      })
      .then((wynik) => {
        if (!wynik.udany) {
          tresc.blad('Rdzeń nie przyjął ustalenia karty.', wynik.blad);
          return;
        }
        const zapisany = wynik.wynik?.entry;
        if (zapisany === undefined) {
          tresc.potwierdzenie(
            'Rdzeń przyjął wywołanie, ale nie oddał zapisanego wpisu — nie wiadomo, co zapisał.',
            false,
          );
          return;
        }
        ustalenie.value = '';
        odczytaj();
        tresc.potwierdzenie(
          `Rdzeń zapisał wpis ${zapisany.id} w zasięgu ${zapisany.scope}, pochodzenie ${zapisany.origin}.`,
          true,
        );
      });
  });

  /** Zaznacza przełączniki poziomów wedle zestawu oddanego przez rdzeń. */
  function przyjmijPoziomy(poziomy: readonly ConfigScope[]): void {
    const wlaczone = new Set<string>(poziomy);
    for (const [poziom, kontrolka] of przelaczniki) kontrolka.checked = wlaczone.has(poziom);
  }

  function odczytaj(): void {
    if (bezKarty()) return;
    tresc.ladowanie('Odczyt pamięci widocznej w karcie sesji…');
    void zrodlo
      .pamiecSesji({
        sessionId: stan.sesja(),
        includeShared: true,
        ...(stan.projekt() === '' ? {} : { projectId: stan.projekt() }),
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Rdzeń nie oddał pamięci karty sesji.', wynik.blad);
          return;
        }
        if (wynik.wynik.length === 0) {
          tresc.pusto('Karta sesji nie widzi jeszcze żadnego ustalenia.');
          return;
        }
        const lista = wykaz('Wpisy widoczne w karcie sesji', 'dw-wykaz');
        for (const wpis of wynik.wynik) lista.append(pozycja(wpis));
        tresc.tresc().append(lista);
      });
  }

  odczyt.addEventListener('click', odczytaj);
  stan.naZmiane(odczytaj);

  return { element, odswiez: odczytaj };
}

/** Pozycja wykazu: treść wpisu wraz z jego zasięgiem i pochodzeniem. */
function pozycja(wpis: WorkspaceMemoryEntry): HTMLElement {
  const { element } = pozycjaWykazu(
    wpis.content,
    `zasięg ${wpis.scope} · pochodzenie ${wpis.origin}${wpis.pinned === true ? ' · przypięty' : ''}`,
    'dw',
  );
  return element;
}

/** Pasek kontrolek panelu — samo złożenie, bez reguły w środku. */
function paskiem(...kontrolki: readonly HTMLElement[]): HTMLElement {
  const pasek = document.createElement('div');
  pasek.className = 'dw-pasek';
  pasek.append(...kontrolki);
  return pasek;
}
