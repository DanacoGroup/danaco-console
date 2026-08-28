import {
  AutomationDependencyKind,
  type AutomationDependency,
} from '../../../../shared/contract';
import {
  poleTekstowe,
  poleWyboru,
  pozycjaWykazu,
  przyciskAkcji as przycisk,
  wykaz,
} from '../../modele/kontrolki-formularza';
import type { ZrodloSekcjiPaneli } from '../../powloka/zrodlo-sekcji-paneli';
import { PRZEDROSTEK } from './kontrolki';
import { utworzPowierzchnieSekcji, zdaniePuste, zObjasnieniem } from './powierzchnia-sekcji';
import { utworzWiezAutomations } from './wiez-automations';
import type { ZrodloNadzoru } from './zrodlo-nadzoru';

/**
 * Sekcja Orkiestracja steruje trzema rodzajami zależności między krokami:
 * sekwencyjną, równoległą i warunkową, na wskazanym układzie prowadzonym przez
 * Orchestrator modułu Automations.
 */
export interface SekcjaOrkiestracji {
  element: HTMLElement;
  odswiez(): void;
  ustawOkno(idOkna: string): void;
}

export interface OpcjeSekcjiOrkiestracji {
  nadzor: ZrodloNadzoru;
  sekcje: ZrodloSekcjiPaneli;
}

/** Trzy rodzaje zależności wraz z tym, co dokładnie znaczą dla biegu pracy w wykonywanym układzie automatyki. */
const RODZAJE: ReadonlyArray<[AutomationDependencyKind, string, string]> = [
  [AutomationDependencyKind.Sequential, 'Sekwencyjna', 'Krok następny rusza dopiero po zakończeniu poprzednika.'],
  [AutomationDependencyKind.Parallel, 'Równoległa', 'Oba kroki biegną obok siebie; żaden nie czeka na drugi.'],
  [AutomationDependencyKind.Conditional, 'Warunkowa', 'Krok następny rusza wyłącznie po spełnieniu warunku podanego niżej.'],
];

export function utworzSekcjeOrkiestracji(
  opcje: OpcjeSekcjiOrkiestracji,
): SekcjaOrkiestracji {
  const { nadzor, sekcje } = opcje;

  const skad = poleTekstowe({ etykieta: 'Krok poprzedzający', podpowiedz: 'identyfikator kroku' });
  const dokad = poleTekstowe({ etykieta: 'Krok następujący', podpowiedz: 'identyfikator kroku' });
  const rodzaj = poleWyboru(
    { etykieta: 'Rodzaj zależności' },
    RODZAJE.map(([wartosc, etykieta]) => ({ wartosc, etykieta })),
  );
  const warunek = poleTekstowe({
    etykieta: 'Warunek przejścia',
    podpowiedz: 'wyrażenie sprawdzane przed krokiem następnym',
  });
  const zapisz = przycisk('Zapisz zależność', 'dn-btn dn-btn--sm');

  const formularz = document.createElement('div');
  formularz.className = 'dm-orkiestracja__pasek';
  formularz.append(
    zObjasnieniem(
      skad.element,
      'Krok, po którym ma ruszyć następny. Identyfikator bierze się z definicji układu w module Automations — zależność wiąże kroki układu, nie okna ról.',
    ),
    zObjasnieniem(
      dokad.element,
      'Krok, którego dotyczy zależność. Zapis nadpisuje zależność istniejącą między tą samą parą kroków, zamiast dokładać drugą.',
    ),
    zObjasnieniem(
      rodzaj.element,
      'Rozstrzyga, czy krok następny czeka na poprzednik (sekwencyjna), biegnie obok niego (równoległa), czy rusza po spełnieniu warunku (warunkowa). Rodzaj zmienia kolejność pracy całego układu.',
    ),
    zObjasnieniem(
      warunek.element,
      'Wyrażenie sprawdzane przed uruchomieniem kroku następnego. Ma znaczenie wyłącznie przy zależności warunkowej; przy pozostałych rdzeń je pomija.',
    ),
    zapisz,
  );

  const lista = wykaz('Zależności układu', 'dm-wykaz');
  const ocena = document.createElement('div');
  ocena.className = 'dm-orkiestracja__ocena';

  const sprawdz = przycisk('Sprawdź układ', 'dn-btn dn-btn--sm dn-btn--zarys');
  sprawdz.title =
    'Pyta rdzeń, czy układ zależności jest wykonalny, i pokazuje ścieżkę krytyczną — kroki, których opóźnienie opóźnia całość.';
  sprawdz.setAttribute('aria-description', sprawdz.title);
  ocena.append(sprawdz);

  const wiez = utworzWiezAutomations({
    okno: 'Orchestrator',
    rola: 'zależności sekwencyjne, równoległe i warunkowe między krokami układu.',
  });

  const powierzchnia = utworzPowierzchnieSekcji({
    klucz: 'orkiestracja',
    tytul: 'Orkiestracja',
    zakres:
      'Zależności między krokami układu: co po czym, co obok czego i co dopiero po spełnieniu warunku.',
    sekcje,
    podsekcje: [
      { id: 'wiez', tytul: 'Powiązanie z modułem Automations', tresc: wiez.element },
      { id: 'nowa', tytul: 'Nowa zależność', tresc: formularz },
      { id: 'zaleznosci', tytul: 'Zależności układu', tresc: lista },
      { id: 'ocena', tytul: 'Ocena układu', tresc: ocena },
    ],
  });

  zapisz.addEventListener('click', () => {
    void zapiszZaleznosc();
  });
  sprawdz.addEventListener('click', () => {
    void sprawdzUklad();
  });

  /** Zapis zależności na wskazanym układzie; brak układu nazywa powód odmowy. */
  async function zapiszZaleznosc(): Promise<void> {
    const uklad = wiez.uklad();
    if (uklad === '') {
      powierzchnia.meldunek(
        'Zależność należy do UKŁADU, a układ nie jest wskazany. Wybierz automatykę w podsekcji powiązania — bez niej żądanie nie ma adresu.',
        false,
      );
      return;
    }
    const od = skad.kontrolka.value.trim();
    const do_ = dokad.kontrolka.value.trim();
    if (od === '' || do_ === '') {
      powierzchnia.meldunek('Zależność wiąże DWA kroki — oba identyfikatory są polami obowiązkowymi kontraktu.', false);
      return;
    }
    const tresc = warunek.kontrolka.value.trim();
    const zaleznosc: AutomationDependency = {
      fromStepId: od,
      toStepId: do_,
      kind: rodzaj.kontrolka.value as AutomationDependencyKind,
      ...(tresc === '' ? {} : { condition: tresc }),
    };
    const wynik = await nadzor.zapiszZaleznosc({ workflowId: uklad, dependency: zaleznosc });
    if (!wynik.udany || wynik.wynik === undefined) {
      powierzchnia.meldunek(`Rdzeń odmówił zapisu zależności. ${powod(wynik.blad?.message, wynik.blad?.code)}`, false);
      return;
    }
    // Rdzeń oddaje układ po zapisie wraz z oceną — na ekran idzie ta odpowiedź.
    pokazZaleznosci(wynik.wynik.dependencies);
    powierzchnia.meldunek(
      wynik.wynik.valid
        ? `Rdzeń zapisał zależność ${od} → ${do_}; układ pozostaje wykonalny (${wynik.wynik.dependencies.length} zależności).`
        : `Rdzeń zapisał zależność ${od} → ${do_}, ale układ ma zastrzeżenia: ${(wynik.wynik.issues ?? []).join('; ') || 'rdzeń nie wypisał treści'}.`,
      wynik.wynik.valid,
    );
  }

  /** Zdjęcie zależności — para kroków, nie identyfikator wpisu. */
  async function zdejmij(zaleznosc: AutomationDependency): Promise<void> {
    const uklad = wiez.uklad();
    if (uklad === '') return;
    const wynik = await nadzor.zdejmijZaleznosc({
      workflowId: uklad,
      fromStepId: zaleznosc.fromStepId,
      toStepId: zaleznosc.toStepId,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      powierzchnia.meldunek(`Rdzeń odmówił zdjęcia zależności. ${powod(wynik.blad?.message, wynik.blad?.code)}`, false);
      return;
    }
    pokazZaleznosci(wynik.wynik.dependencies);
    powierzchnia.meldunek(
      wynik.wynik.removed
        ? `Zależność ${zaleznosc.fromStepId} → ${zaleznosc.toStepId} zdjęta; zostaje ${wynik.wynik.dependencies.length}.`
        : `Rdzeń przyjął żądanie, ale zależności ${zaleznosc.fromStepId} → ${zaleznosc.toStepId} nie zdjął — w układzie jej nie było.`,
      wynik.wynik.removed,
    );
  }

  /** Ocena układu wraz ze ścieżką krytyczną prosto z rdzenia. */
  async function sprawdzUklad(): Promise<void> {
    const uklad = wiez.uklad();
    if (uklad === '') {
      powierzchnia.meldunek('Sprawdzenie dotyczy układu, a układ nie jest wskazany.', false);
      return;
    }
    const wynik = await nadzor.sprawdzUklad(uklad);
    if (!wynik.udany || wynik.wynik === undefined) {
      powierzchnia.meldunek(`Rdzeń odmówił sprawdzenia układu. ${powod(wynik.blad?.message, wynik.blad?.code)}`, false);
      return;
    }
    const sciezka = wynik.wynik.criticalPathStepIds ?? [];
    powierzchnia.meldunek(
      wynik.wynik.valid
        ? `Układ wykonalny. Ścieżka krytyczna: ${sciezka.length === 0 ? 'rdzeń jej nie wyliczył' : sciezka.join(' → ')}.`
        : `Układ ma zastrzeżenia: ${(wynik.wynik.issues ?? []).join('; ') || 'rdzeń nie wypisał treści'}.`,
      wynik.wynik.valid,
    );
  }

  /** Wykaz zależności; pustka jest tu stanem poprawnym. */
  function pokazZaleznosci(zaleznosci: readonly AutomationDependency[]): void {
    if (zaleznosci.length === 0) {
      lista.replaceChildren(
        zdaniePuste(
          'Ten układ nie ma jeszcze ani jednej zależności — jego kroki biegną niezależnie. To jest stan poprawny, nie brak.',
        ),
      );
      return;
    }
    lista.replaceChildren(
      ...zaleznosci.map((zaleznosc) => {
        const { element, akcje } = pozycjaWykazu(
          `${zaleznosc.fromStepId} → ${zaleznosc.toStepId}`,
          `${zaleznosc.kind}${zaleznosc.condition === undefined ? '' : ` · warunek: ${zaleznosc.condition}`}`,
          PRZEDROSTEK,
        );
        const zdejmijJa = przycisk('Zdejmij', 'dn-btn dn-btn--sm dn-btn--zarys');
        zdejmijJa.addEventListener('click', () => {
          void zdejmij(zaleznosc);
        });
        akcje.append(zdejmijJa);
        return element;
      }),
    );
  }

  /** Odczyt: wykaz automatyk pod więź, a po wskazaniu układu — zależności. */
  function odczytaj(): void {
    powierzchnia.tresci.ladowanie('Odczyt układów i zależności…');
    void (async () => {
      const automatyki = await nadzor.automatyki({});
      wiez.ustawWykaz(automatyki.udany ? (automatyki.wynik ?? []) : []);

      const uklad = wiez.uklad();
      if (uklad === '') {
        powierzchnia.tresci.pusto('');
        lista.replaceChildren(
          zdaniePuste(
            'Zależności czyta się z układu, a układ nie jest wskazany. Wybierz automatykę w podsekcji powiązania — wykaz wejdzie tutaj.',
          ),
        );
        return;
      }
      const wynik = await nadzor.zaleznosci({ workflowId: uklad });
      if (!wynik.udany || wynik.wynik === undefined) {
        powierzchnia.tresci.blad('Odczyt zależności układu odmówiony.', wynik.blad);
        return;
      }
      powierzchnia.tresci.pusto('');
      pokazZaleznosci(wynik.wynik);
    })();
  }

  wiez.naZmiane(() => odczytaj());

  return {
    element: powierzchnia.element,
    odswiez: odczytaj,
    ustawOkno: powierzchnia.ustawOkno,
  };
}

/** Treść odmowy wraz z kodem kontraktu, złożona w jedno pełne zdanie gotowe do pokazania w meldunku sekcji. */
function powod(zdanie: string | undefined, kod: string | undefined): string {
  return `Powód: ${zdanie ?? 'rdzeń nie podał przyczyny'} (kod ${kod ?? 'brak'}).`;
}
