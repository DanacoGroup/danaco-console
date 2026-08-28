import {
  ConfigScope,
  MemoryEntryOrigin,
  type WorkspaceMemoryEntry,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  poleTresci,
  przelacznikWidoku,
  przestaw,
  przyciskAkcji as przycisk,
  wiersz,
  wybor,
  wykaz,
} from '../../modele/kontrolki-formularza';
import { eksportPamieci, pozycjaPamieci } from './pamiec-pozycja';
import { poziomyZPierwszym, type ParaPoziomu } from './poziomy-zasiegu';
import { utworzPanelPamieciSesji } from './pamiec-sesji';
import type { StanProjektu } from './stan-projektu';
import { utworzStanTresci } from './stany-okna';
import type { ZrodloWorkspace } from './zrodlo-workspace';

/**
 * Context Memory to okno zarządcy modułu Workspace: zapis ustalenia, przegląd pamięci i edycja
 * wpisu odrębnej dla każdego projektu.
 */
export interface OknoPamieci {
  element: HTMLElement;
  odswiez(): void;
}

/**
 * Zasięgi współdzielenia wpisu tworzą ten sam wykaz, co w Instructions Panel i w panelu pamięci
 * karty sesji, z projektem na czele jako domyślnym zasięgiem.
 */
const ZASIEGI: readonly ParaPoziomu[] = poziomyZPierwszym(ConfigScope.Project);

/**
 * Pochodzenia w filtrze wykazu. Pusta wartość znaczy „bez zawężania” i stoi
 * pierwsza, bo wykaz nieprzefiltrowany jest stanem wyjściowym okna.
 */
const POCHODZENIA: ReadonlyArray<readonly [string, string]> = [
  ['', 'pochodzenie: wszystkie'],
  [MemoryEntryOrigin.Operator, 'pochodzenie: Operator'],
  [MemoryEntryOrigin.Model, 'pochodzenie: propozycje modelu'],
];

export function utworzOknoPamieci(
  zrodlo: ZrodloWorkspace,
  stan: StanProjektu,
): OknoPamieci {
  const rama = utworzRameOkna({
    tytul: 'Context Memory',
    kod: 'context-memory',
    rola: 'zarządca',
    przeznaczenie: 'Ustalenia projektu wraz z ich zasięgiem, przypięciem i pochodzeniem.',
    przedrostek: 'dw',
  });
  const tresc = utworzStanTresci();

  const edytor = poleTresci('Treść ustalenia', 4);
  const zasieg = wybor('Zasięg współdzielenia wpisu', ZASIEGI.map(([wartosc, opis]) => [wartosc, opis]));
  const przypiecie = document.createElement('input');
  przypiecie.type = 'checkbox';
  przypiecie.className = 'dn-przelacznik';
  const wspolne = document.createElement('input');
  wspolne.type = 'checkbox';
  wspolne.className = 'dn-przelacznik';

  const dodaj = przycisk('+ Dodaj wpis', 'dn-btn dn-btn--atrament');
  const odswiez = przycisk('Odśwież pamięć');
  const eksport = przycisk('Eksport pamięci');
  const usunEdytowany = przycisk('Usuń wpis');
  usunEdytowany.addEventListener('click', () => usunZEdytora());

  rama.akcje.append(dodaj, odswiez, eksport, usunEdytowany);

  // Filtr zawęża wykaz już odczytany, stoi w pasku narzędzi, nie w panelu, i nic nie wysyła rdzeniowi.
  const fraza = pole('Fraza w treści wpisu', 'fragment ustalenia');
  const pochodzenie = wybor('Pochodzenie wpisu', POCHODZENIA);
  const tylkoPrzypiete = przelacznikWidoku('Tylko przypięte', false);
  rama.narzedzia.append(fraza, pochodzenie, tylkoPrzypiete);

  rama.cialo.append(
    wiersz('Ustalenie', edytor, { klasa: 'dw-wiersz', objasnienie: 'Wpis wchodzi do kontekstu pracy modelu w tym projekcie.' }),
    wiersz('Zasięg', zasieg, { klasa: 'dw-wiersz', objasnienie: 'Domyślnie projekt — pamięć odrębna; zasięg szerszy czyni z wpisu ustalenie wspólne.' }),
    wiersz('Przypnij wpis', przypiecie, { klasa: 'dw-wiersz', objasnienie: 'Przypięte stoją na początku wykazu.' }),
    wiersz('Pokaż ustalenia wspólne', wspolne, { klasa: 'dw-wiersz', objasnienie: 'Dokłada wpisy innych projektów zapisane na poziomie szerszym.' }),
    tresc.element,
  );

  // Pamięć karty sesji stoi w tym oknie, ale prowadzi ją osobny panel i osobna rodzina komend.
  const pamiecKarty = utworzPanelPamieciSesji(zrodlo, stan);
  rama.cialo.append(pamiecKarty.element);

  /** Wpis w trakcie edycji; pusty znaczy „zapis zakłada wpis nowy”. */
  let zmieniany: WorkspaceMemoryEntry | null = null;
  let ostatnie: WorkspaceMemoryEntry[] = [];
  /** Stan przełącznika „Tylko przypięte” — nośnikiem jest przycisk, nie pole. */
  let filtrPrzypiete = false;

  function zapisz(zadanie: {
    tresc: string;
    zasieg: ConfigScope;
    przypiety: boolean;
    pochodzenie: MemoryEntryOrigin;
    wpis?: string;
  }): void {
    const projekt = stan.projekt();
    if (projekt === '') {
      tresc.blad('Zapis bez projektu nie ma zakresu — wskaż projekt na pulpicie.');
      return;
    }
    if (zadanie.tresc.trim() === '') {
      tresc.potwierdzenie('Wpis bez treści nie zostanie zapisany.', false);
      return;
    }
    void zrodlo
      .zapiszWpisPamieci({
        projectId: projekt,
        content: zadanie.tresc,
        entryId: zadanie.wpis ?? '',
        pinned: zadanie.przypiety,
        origin: zadanie.pochodzenie,
        scope: zadanie.zasieg,
      })
      .then((wynik) => {
        if (!wynik.udany) {
          tresc.blad('Rdzeń nie przyjął wpisu pamięci.', wynik.blad);
          return;
        }
        edytor.value = '';
        zmieniany = null;
        odczytaj();
        // Potwierdzenie mówi to, co zapisał rdzeń, nie to, co wysłało okno; odświeżenie samo nie potwierdza.
        const zapisany = wynik.wynik?.entry;
        if (zapisany === undefined) {
          tresc.potwierdzenie(
            'Rdzeń przyjął wywołanie, ale nie oddał zapisanego wpisu — nie wiadomo, co zapisał.',
            false,
          );
          return;
        }
        tresc.potwierdzenie(
          `Rdzeń zapisał wpis ${zapisany.id} — zasięg ${zapisany.scope}, ` +
            `pochodzenie ${zapisany.origin}, ${zapisany.pinned === true ? 'przypięty' : 'nieprzypięty'}.`,
          true,
        );
      });
  }

  /** Usuwa wpis komendą usunięcia; potwierdzenie mówi to, co oddał rdzeń, nie to, co wysłało okno. */
  function usun(wpis: WorkspaceMemoryEntry): void {
    void zrodlo.usunWpisPamieci({ entryId: wpis.id }).then((wynik) => {
      if (!wynik.udany) {
        tresc.blad('Rdzeń nie usunął wpisu pamięci.', wynik.blad);
        return;
      }
      if (wynik.wynik?.deleted !== true) {
        tresc.potwierdzenie('Rdzeń przyjął wywołanie, ale nie potwierdził usunięcia wpisu.', false);
        return;
      }
      if (zmieniany !== null && zmieniany.id === wpis.id) {
        zmieniany = null;
        edytor.value = '';
      }
      odczytaj();
      tresc.potwierdzenie(`Rdzeń usunął wpis ${wpis.id} z pamięci projektu.`, true);
    });
  }

  /** „Usuń wpis” z nagłówka dotyczy wpisu wczytanego do edytora („Edytuj”). */
  function usunZEdytora(): void {
    if (zmieniany === null) {
      tresc.potwierdzenie(
        'Wczytaj wpis do edytora („Edytuj” przy pozycji), zanim usuniesz go tym przyciskiem — ' +
          'albo użyj „Usuń” bezpośrednio przy wybranej pozycji.',
        false,
      );
      return;
    }
    usun(zmieniany);
  }

  function pozycja(wpis: WorkspaceMemoryEntry): HTMLElement {
    return pozycjaPamieci(
      wpis,
      {
        zapisz,
        usun,
        wczytajDoEdytora(zrodlowy) {
          zmieniany = zrodlowy;
          edytor.value = zrodlowy.content;
          zasieg.value = zrodlowy.scope;
          przypiecie.checked = zrodlowy.pinned === true;
          tresc.potwierdzenie('Wpis wczytany do edytora — „+ Dodaj wpis” zapisze zmianę.', true);
        },
        dolaczDoEdytora(zrodlowy) {
          edytor.value =
            edytor.value === '' ? zrodlowy.content : `${edytor.value}
${zrodlowy.content}`;
          // Zdanie mówi o skutku scalenia, a powodu, dla którego wpis źródłowy zostaje, nie powtarza.
          tresc.potwierdzenie(
            'Treść dołączona do edytora. Zapis scali ją w jeden wpis; wpis źródłowy zostaje ' +
              'w pamięci projektu, dopóki nie usuniesz go przyciskiem „Usuń”.',
            true,
          );
        },
      },
    );
  }

  /** Czy wpis przechodzi przez nastawę filtra w pasku narzędzi. */
  function przechodziFiltr(wpis: WorkspaceMemoryEntry): boolean {
    const szukana = fraza.value.trim().toLocaleLowerCase('pl-PL');
    if (szukana !== '' && !wpis.content.toLocaleLowerCase('pl-PL').includes(szukana)) return false;
    if (pochodzenie.value !== '' && wpis.origin !== pochodzenie.value) return false;
    return !(filtrPrzypiete && wpis.pinned !== true);
  }

  /** Odczyt i rysowanie są rozdzielone, bo filtr zawęża materiał już posiadany bez pytania rdzenia. */
  function rysuj(): void {
    if (ostatnie.length === 0) {
      rama.ustawZnacznik('');
      tresc.pusto('Projekt nie ma jeszcze ustaleń. Pierwszy wpis możesz dodać powyżej.');
      return;
    }
    const widoczne = ostatnie.filter(przechodziFiltr);
    rama.ustawZnacznik(
      widoczne.length === ostatnie.length
        ? `wpisy: ${ostatnie.length}`
        : `wpisy: ${widoczne.length} z ${ostatnie.length}`,
    );
    if (widoczne.length === 0) {
      tresc.pusto('Żaden wpis pamięci projektu nie odpowiada nastawie filtra w pasku narzędzi.');
      return;
    }
    const lista = wykaz('Wpisy pamięci projektu', 'dw-wykaz');
    for (const wpis of widoczne) lista.append(pozycja(wpis));
    tresc.tresc().append(lista);
  }

  function odczytaj(): void {
    const projekt = stan.projekt();
    if (projekt === '') {
      ostatnie = [];
      rama.ustawZnacznik('');
      tresc.pusto('Wskaż projekt na pulpicie, aby zobaczyć jego pamięć.');
      return;
    }
    tresc.ladowanie('Odczyt pamięci projektu…');
    void zrodlo
      .wpisyPamieci({ projectId: projekt, includeShared: wspolne.checked })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          // Plakietka gaśnie razem z wykazem: licznik sprzed odmowy mówiłby o pamięci, której nie odczytano.
          ostatnie = [];
          rama.ustawZnacznik('');
          tresc.blad('Rdzeń nie oddał pamięci projektu.', wynik.blad);
          return;
        }
        ostatnie = wynik.wynik;
        rysuj();
      });
  }

  dodaj.addEventListener('click', () =>
    zapisz({
      tresc: edytor.value,
      zasieg: zasieg.value as ConfigScope,
      przypiety: przypiecie.checked,
      pochodzenie: MemoryEntryOrigin.Operator,
      ...(zmieniany === null ? {} : { wpis: zmieniany.id }),
    }),
  );
  odswiez.addEventListener('click', odczytaj);
  // Zasięg wchodzi do zapytania, więc jego zmiana wymaga odczytu; filtr wystarcza mu przerysowanie.
  wspolne.addEventListener('change', odczytaj);
  fraza.addEventListener('input', rysuj);
  pochodzenie.addEventListener('change', rysuj);
  tylkoPrzypiete.addEventListener('click', () => {
    filtrPrzypiete = przestaw(tylkoPrzypiete);
    rysuj();
  });
  eksport.addEventListener('click', () => {
    if (ostatnie.length === 0) {
      tresc.potwierdzenie('Pamięć jest pusta — nie ma czego wyeksportować.', false);
      return;
    }
    pobierzPlik(`${stan.projekt()}-pamiec.md`, eksportPamieci(ostatnie), 'text/markdown');
    tresc.potwierdzenie('Pamięć projektu pobrana jako plik Markdown.', true);
  });

  // Zmiana projektu, także zdarzeniem z innego okna, przeładowuje wykaz, by pamięć nie była nieaktualna.
  stan.naZmiane(odczytaj);

  return { element: rama.element, odswiez: odczytaj };
}
