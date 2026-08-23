import { QueueAction, type Queue } from '../../../../shared/contract';
import {
  poleTekstowe,
  poleWyboru,
  pozycjaWykazu,
  przyciskAkcji as przycisk,
  wykaz,
} from '../../modele/kontrolki-formularza';
import type { ZrodloSekcjiPaneli } from '../../powloka/zrodlo-sekcji-paneli';
import { PRZEDROSTEK, wierszOpisu } from './kontrolki';
import { utworzPowierzchnieSekcji, zdaniePuste, zObjasnieniem } from './powierzchnia-sekcji';
import { utworzWiezAutomations } from './wiez-automations';
import type { ZrodloBiegu } from './zrodlo-biegu';
import type { ZrodloNadzoru } from './zrodlo-nadzoru';

/**
 * Sekcja KOLEJKI panelu orkiestracji — sterowanie silnikiem kolejek rdzenia.
 *
 * `queue.create`, `queue.action` i `queue.list` obsługują pętlę sesyjną, moduł
 * Automations i to środowisko; sekcja nimi steruje, zamiast budować własnego
 * wykonawcę zleceń. Okno robocze jest cudze — po skonfigurowaniu powiązania
 * kolejkę prowadzi Queue Manager modułu Automations.
 *
 * `QueueAction` kontraktu niesie `start · pause · resume · stop · retry · clear`.
 * Czynności spoza tego zbioru (`enqueue`, `dequeue`, `delay`, `split`/`merge`,
 * `route`/`branch`/`condition`) sekcja wypisuje wraz z brakującym kształtem
 * komendy, zamiast składać je z czynności istniejących: przycisk „split" oparty
 * na dwóch `retry` robiłby co innego, niż mówi jego napis.
 *
 * `Queue` niesie sesję, okna, stan i licznik obiegów — pola zasięgu nie ma.
 * Zasięg jest więc wykazem opisowym wraz z kształtem brakującego pola, a nie
 * kontrolką bez skutku.
 */
export interface SekcjaKolejek {
  element: HTMLElement;
  odswiez(): void;
  ustawOkno(idOkna: string): void;
  rozlacz(): void;
}

export interface OpcjeSekcjiKolejek {
  bieg: ZrodloBiegu;
  nadzor: ZrodloNadzoru;
  sekcje: ZrodloSekcjiPaneli;
  /** Karta sesji — `queue.create` i `queue.list` adresują właśnie ją. */
  sesja(): string;
}

/** Zasięgi kolejki wraz z ich znaczeniem; kontrakt pola zasięgu nie niesie. */
const ZASIEGI: ReadonlyArray<[string, string]> = [
  ['globalna', 'jedna kolejka na całą instalację'],
  ['lokalna', 'kolejka jednej karty sesji'],
  ['dla modelu', 'kolejka wspólna oknom jednego kanału modelu'],
  ['dla agenta', 'kolejka wspólna oknom jednego eksperta'],
  ['dla projektu', 'kolejka wspólna oknom jednego projektu'],
];

/** Czynności bez odpowiednika w `QueueAction` wraz z kształtem brakującej komendy. */
const BRAKI_AKCJI: ReadonlyArray<[string, string]> = [
  ['enqueue', 'queue.action { queueId, action: "enqueue", item: QueueItem }'],
  ['dequeue', 'queue.action { queueId, action: "dequeue", itemId }'],
  ['delay', 'queue.action { queueId, action: "delay", itemId, delayMs }'],
  ['split / merge', 'queue.action { queueId, action: "split" | "merge", itemId, parts?: number }'],
  ['route / branch / condition', 'queue.action { queueId, action: "route" | "branch" | "condition", itemId, target?, condition? }'],
];

/** Kształt pola zasięgu, którego `Queue` kontraktu nie ma. */
const BRAK_ZASIEGU =
  'Queue { …, scope: "global" | "local" | "model" | "agent" | "project", scopeId?: string } oraz queue.create przyjmujące te dwa pola.';

/** Sześć czynności kontraktu wraz z tym, co robią. */
const AKCJE: ReadonlyArray<[QueueAction, string, string]> = [
  [QueueAction.Start, 'Uruchom', 'Wprawia kolejkę w bieg: pozycje zaczynają schodzić do okien, które kolejka obsługuje.'],
  [QueueAction.Pause, 'Wstrzymaj', 'Zatrzymuje pobieranie kolejnych pozycji; pozycja w biegu dochodzi do końca.'],
  [QueueAction.Resume, 'Wznów', 'Wraca do pobierania pozycji od miejsca wstrzymania; nic nie jest wykonywane po raz drugi.'],
  [QueueAction.Stop, 'Zatrzymaj', 'Kończy bieg kolejki decyzją Operatora. Przycisk zatrzymania jest zawsze czynny.'],
  [QueueAction.Retry, 'Powtórz krok', 'Ponawia krok bieżący i podnosi licznik obiegów naprawczych; liczba obiegów jest nieograniczona.'],
  [QueueAction.Clear, 'Opróżnij', 'Zdejmuje pozycje oczekujące. Kolejka zostaje, jej zawartość nie.'],
];

export function utworzSekcjeKolejek(opcje: OpcjeSekcjiKolejek): SekcjaKolejek {
  const { bieg, nadzor, sekcje, sesja } = opcje;

  const nazwa = poleTekstowe({
    etykieta: 'Nazwa kolejki',
    podpowiedz: 'np. Pętla ciągła wykonawców',
  });
  const zaloz = przycisk('Załóż kolejkę', 'dn-btn dn-btn--sm');

  const pasek = document.createElement('div');
  pasek.className = 'dm-orkiestracja__pasek';
  pasek.append(
    zObjasnieniem(
      nazwa.element,
      'Nazwa kolejki widoczna w wykazie i w oknie Queue Manager. Kolejka powstaje na tej karcie sesji i od razu jest prawdziwym bytem rdzenia — nie zapisem w widoku.',
    ),
    zaloz,
  );

  const lista = wykaz('Kolejki tej karty sesji', 'dm-wykaz');
  const zasiegi = wykaz('Pięć zasięgów kolejki', 'dm-wykaz');
  const braki = wykaz('Czynności bez pokrycia w kontrakcie', 'dm-wykaz');

  const wiez = utworzWiezAutomations({
    okno: 'Queue Manager',
    rola: 'obieg zadań, wstrzymywanie, ponawianie i kierowanie warunkowe pozycji.',
  });

  const powiazanie = poleWyboru({ etykieta: 'Kolejka do powiązania' }, []);
  const powiaz = przycisk('Powiąż z układem', 'dn-btn dn-btn--sm dn-btn--zarys');
  const paskiWiezi = document.createElement('div');
  paskiWiezi.className = 'dm-orkiestracja__pasek';
  paskiWiezi.append(
    zObjasnieniem(
      powiazanie.element,
      'Kolejka, którą wiążemy ze wskazanym układem automatyki. Po powiązaniu Queue Manager modułu Automations prowadzi tę kolejkę jako część pracy ciągłej; powiązanie jest odwracalne.',
    ),
    powiaz,
  );
  wiez.element.append(paskiWiezi);

  const powierzchnia = utworzPowierzchnieSekcji({
    klucz: 'kolejki',
    tytul: 'Kolejki',
    zakres:
      'Silnik kolejek pętli sesyjnej — ten sam, którym jedzie moduł Automations. Zakładanie, sterowanie i wiązanie kolejek tej karty sesji.',
    sekcje,
    podsekcje: [
      { id: 'zaloz', tytul: 'Nowa kolejka', tresc: pasek },
      { id: 'kolejki', tytul: 'Kolejki tej karty sesji', tresc: lista },
      { id: 'wiez', tytul: 'Powiązanie z modułem Automations', tresc: wiez.element },
      { id: 'zasiegi', tytul: 'Pięć zasięgów kolejki', tresc: zasiegi },
      { id: 'braki', tytul: 'Czynności bez pokrycia w kontrakcie', tresc: braki },
    ],
  });

  zasiegi.replaceChildren(
    ...ZASIEGI.map(([miano, znaczenie]) => pozycjaWykazu(miano, znaczenie, PRZEDROSTEK).element),
    zdaniePuste(
      `Zasięg jest dziś opisem, nie polem: kontrakt go nie niesie. Kształt brakujący: ${BRAK_ZASIEGU}`,
    ),
  );

  braki.replaceChildren(
    ...BRAKI_AKCJI.map(([miano, ksztalt]) => pozycjaWykazu(miano, ksztalt, PRZEDROSTEK).element),
    zdaniePuste(
      'Te czynności NIE MAJĄ przycisków, bo rdzeń nie ma ich uchwytów. Przycisk składający je z czynności istniejących robiłby co innego, niż mówi jego napis.',
    ),
  );

  zaloz.addEventListener('click', () => {
    void zalozKolejke();
  });
  powiaz.addEventListener('click', () => {
    void powiazZUkladem();
  });

  /** Zakłada kolejkę na tej karcie sesji; nazwa pusta idzie bez pola. */
  async function zalozKolejke(): Promise<void> {
    const karta = sesja();
    if (karta === '') {
      powierzchnia.meldunek(
        'Kolejka należy do karty sesji, a sesja nie jest jeszcze uzgodniona z rdzeniem — nie ma czego adresować.',
        false,
      );
      return;
    }
    const miano = nazwa.kontrolka.value.trim();
    const wynik = await bieg.zalozKolejke({
      sessionId: karta,
      ...(miano === '' ? {} : { name: miano }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      powierzchnia.meldunek(`Rdzeń odmówił założenia kolejki. ${powod(wynik.blad?.message, wynik.blad?.code)}`, false);
      return;
    }
    powierzchnia.meldunek(
      `Rdzeń założył kolejkę ${wynik.wynik.id} w stanie „${wynik.wynik.status}".`,
      true,
    );
    nazwa.kontrolka.value = '';
    odczytaj();
  }

  /** Wiąże wybraną kolejkę ze wskazanym układem automatyki. */
  async function powiazZUkladem(): Promise<void> {
    const kolejka = powiazanie.kontrolka.value;
    const uklad = wiez.uklad();
    if (kolejka === '' || uklad === '') {
      powierzchnia.meldunek(
        'Powiązanie potrzebuje obu stron: kolejki i układu automatyki. Wskaż je wyżej — bez nich żądanie nie ma dokąd pojechać.',
        false,
      );
      return;
    }
    const wynik = await nadzor.powiazKolejke({ queueId: kolejka, workflowId: uklad });
    if (!wynik.udany || wynik.wynik === undefined) {
      powierzchnia.meldunek(`Rdzeń odmówił powiązania kolejki. ${powod(wynik.blad?.message, wynik.blad?.code)}`, false);
      return;
    }
    powierzchnia.meldunek(
      `Rdzeń powiązał kolejkę ${wynik.wynik.id} z układem ${uklad}; Queue Manager prowadzi ją odtąd jako część tej automatyki.`,
      true,
    );
    odczytaj();
  }

  /** Czynność na kolejce; stan po niej bierze się z odpowiedzi rdzenia. */
  async function steruj(kolejka: Queue, akcja: QueueAction, miano: string): Promise<void> {
    const wynik = await bieg.sterujKolejka({ queueId: kolejka.id, action: akcja });
    if (!wynik.udany || wynik.wynik === undefined) {
      powierzchnia.meldunek(
        `Rdzeń odmówił czynności „${miano}" na kolejce ${kolejka.id}. ${powod(wynik.blad?.message, wynik.blad?.code)}`,
        false,
      );
      return;
    }
    powierzchnia.meldunek(
      `Czynność „${miano}" wykonana: kolejka ${wynik.wynik.id} jest w stanie „${wynik.wynik.status}", obieg ${wynik.wynik.cycle ?? 0}.`,
      true,
    );
    odczytaj();
  }

  /** Odczyt kolejek karty sesji wraz z wykazem automatyk pod więź. */
  function odczytaj(): void {
    powierzchnia.tresci.ladowanie('Odczyt kolejek karty sesji…');
    void (async () => {
      const automatyki = await nadzor.automatyki({});
      wiez.ustawWykaz(automatyki.udany ? (automatyki.wynik ?? []) : []);

      const karta = sesja();
      const wynik = await bieg.wykazKolejek(karta === '' ? {} : { sessionId: karta });
      if (!wynik.udany || wynik.wynik === undefined) {
        powierzchnia.tresci.blad('Odczyt kolejek odmówiony.', wynik.blad);
        return;
      }
      powierzchnia.tresci.pusto('');
      const kolejki = wynik.wynik;
      powiazanie.kontrolka.replaceChildren(
        ...kolejki.map((kolejka) => {
          const pozycja = document.createElement('option');
          pozycja.value = kolejka.id;
          pozycja.textContent = kolejka.name ?? kolejka.id;
          return pozycja;
        }),
      );
      if (kolejki.length === 0) {
        lista.replaceChildren(
          zdaniePuste(
            'Ta karta sesji nie ma jeszcze ani jednej kolejki. Pustka jest tu stanem poprawnym — kolejkę zakłada Operator albo koordynator, gdy pętla rusza.',
          ),
        );
        return;
      }
      lista.replaceChildren(...kolejki.map(pozycjaKolejki));
    })();
  }

  /** Jedna kolejka wraz z sześcioma czynnościami kontraktu. */
  function pozycjaKolejki(kolejka: Queue): HTMLElement {
    const { element, akcje } = pozycjaWykazu(
      kolejka.name ?? kolejka.id,
      `stan ${kolejka.status} · obieg ${kolejka.cycle ?? 0} · okna ${kolejka.windowIds?.length ?? 0}`,
      PRZEDROSTEK,
    );
    for (const [akcja, miano, objasnienie] of AKCJE) {
      const kontrolka = przycisk(miano, 'dn-btn dn-btn--sm dn-btn--zarys');
      // Objaśnienie idzie na przycisk trzema drogami, tak jak przy kontrolkach
      // z powodem blokady: wskaźnik, czytnik ekranu i sprawdzian mają je widzieć.
      kontrolka.title = objasnienie;
      kontrolka.setAttribute('aria-description', objasnienie);
      kontrolka.addEventListener('click', () => {
        void steruj(kolejka, akcja, miano);
      });
      akcje.append(kontrolka);
    }
    element.append(wierszOpisu('Założona', new Date(kolejka.createdAt).toLocaleString('pl-PL')));
    return element;
  }

  wiez.naZmiane(() => odczytaj());
  // `queue.changed` niesie stan kolejki na żywo: kolejka rusza i staje bez
  // udziału tego widoku, więc bez subskrypcji wykaz zastygałby na chwili odczytu.
  const odsubskrybuj = bieg.naKolejke(() => odczytaj());

  return {
    element: powierzchnia.element,
    odswiez: odczytaj,
    ustawOkno: powierzchnia.ustawOkno,
    rozlacz: odsubskrybuj,
  };
}

/** Treść odmowy wraz z kodem kontraktu. */
function powod(zdanie: string | undefined, kod: string | undefined): string {
  return `Powód: ${zdanie ?? 'rdzeń nie podał przyczyny'} (kod ${kod ?? 'brak'}).`;
}
