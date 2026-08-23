import type { AutomationSchedule, AutomationWorkflow } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  pole,
  pozycjaWykazu,
  przelacznik,
  przyciskAkcji,
  ustawPozycje,
  utworzWierszOdpowiedzi,
  wiersz,
  wybor,
  wykaz,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { chwila, ODCZYTY, PUSTE } from './etykiety-assistant';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { ZrodloNarzedzi } from './zrodlo-narzedzi';

/**
 * Zakładka rutyn Command & Tools Hub — automatyki asystenta wraz z ich
 * cyklicznością.
 *
 * Rutyna („poranny brief o 8:00 w dni robocze") jest w kontrakcie automatyką
 * z harmonogramem: `automation.workflow.list` mówi, jakie automatyki są,
 * `schedule.get` — kiedy biegną, a `automation.schedule.set` nadaje im
 * cykliczność. Osobnego bytu rutyny asystenta kontrakt nie ma i okno go nie
 * zakłada.
 *
 * Wykaz i harmonogramy czytane są jednym ruchem, bo bez pary tych odczytów
 * automatyka bez harmonogramu byłaby nie do odróżnienia od automatyki, której
 * harmonogramu okno jeszcze nie zna.
 *
 * Wyzwalaczy zdarzeniowych okno nie ustawia. `AutomationTrigger` niesie ich
 * cztery rodzaje, ale każdy wymaga wyrażenia właściwego dla swojego rodzaju
 * (adres wywołania zdalnego, ścieżka pliku, warunek na wyniku modelu) —
 * ich redakcja należy do modułu Automations, którego Workflow Builder jest
 * miejscem budowy procesów. Assistant, zgodnie z granicą tematyczną modułu,
 * inicjuje i nadzoruje pojedyncze zlecenia, a nie projektuje pełnych procesów.
 */
export interface PanelRutyn {
  element: HTMLElement;
  /** Odczyt automatyk i ich harmonogramów. */
  wczytaj(): Promise<void>;
}

export function utworzPanelRutyn(zrodlo: ZrodloNarzedzi): PanelRutyn {
  const okno: StanOkna = utworzStanOkna();
  const odpowiedz: WierszOdpowiedzi = utworzWierszOdpowiedzi();

  const automatyka = wybor('Automatyka objęta harmonogramem', [['', 'Automatyka: nie wskazano']]);
  const cyklicznosc = pole('Cykliczność w notacji cron', 'na przykład: 0 8 * * 1-5');
  const strefa = pole('Strefa czasowa harmonogramu', 'na przykład: Europe/Warsaw');
  const obowiazuje = przelacznik('Harmonogram obowiązuje');
  obowiazuje.checked = true;

  const lista = wykaz('Automatyki wraz z harmonogramami', 'ma-wykaz');
  okno.tresc.append(lista);
  okno.puste(PUSTE.rutynySpoczynek);

  const element = document.createElement('div');
  element.className = 'ma-obszar';
  element.dataset['obszar'] = 'rutyny';
  element.append(
    wiersz('Automatyka', automatyka, {
      klasa: 'ma-wiersz',
      objasnienie:
        'Wykaz pochodzi z automation.workflow.list. Makro wysłane z zakładki obok staje ' +
        'tu po odczycie.',
    }),
    wiersz('Cykliczność', cyklicznosc, {
      klasa: 'ma-wiersz',
      objasnienie:
        'Notacja cron rdzenia. Puste pole zostawia automatykę bez cykliczności — biegnie ' +
        'wtedy wyłącznie na żądanie.',
    }),
    wiersz('Strefa czasowa', strefa, {
      klasa: 'ma-wiersz',
      objasnienie: 'Puste pole zostawia strefę wartości domyślnej rdzenia.',
    }),
    wiersz('Obowiązuje', obowiazuje, {
      klasa: 'ma-wiersz',
      objasnienie: 'Wstrzymanie zostawia harmonogram zapisany, lecz nie wyzwala przebiegów.',
    }),
    przyciski(),
    okno.element,
    odpowiedz.element,
  );

  function przyciski(): HTMLElement {
    const czytaj = przyciskAkcji('Odczytaj rutyny', 'dn-btn dn-btn--sm dn-btn--zarys');
    czytaj.addEventListener('click', () => void wczytaj());

    const zapisz = przyciskAkcji('Zapisz harmonogram rutyny', 'dn-btn dn-btn--sm dn-btn--atrament');
    zapisz.addEventListener('click', () => void zapiszHarmonogram());

    const rzad = document.createElement('div');
    rzad.className = 'ma-formularz__przyciski';
    rzad.append(czytaj, zapisz);
    return rzad;
  }

  async function zapiszHarmonogram(): Promise<void> {
    if (automatyka.value === '') {
      odpowiedz.pokaz(
        'Wskaż automatykę — automation.schedule.set wymaga pola workflowId, więc bez niej ' +
          'harmonogram nie ma do czego przylgnąć.',
        false,
      );
      return;
    }
    odpowiedz.pokaz(ODCZYTY.harmonogram, true);
    const wynik = await zrodlo.ustawHarmonogram({
      workflowId: automatyka.value,
      enabled: obowiazuje.checked,
      ...(cyklicznosc.value.trim() === '' ? {} : { cron: cyklicznosc.value.trim() }),
      ...(strefa.value.trim() === '' ? {} : { timeZone: strefa.value.trim() }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Zapis harmonogramu rutyny', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    odpowiedz.pokaz(`Rdzeń zapisał harmonogram: ${opisHarmonogramu(wynik.wynik.schedule)}`, true);
    await wczytaj();
  }

  async function wczytaj(): Promise<void> {
    okno.ladowanie(ODCZYTY.rutyny);
    // Oba odczyty idą równolegle: automatyki i harmonogramy to dwie różne
    // komendy, a żadna nie warunkuje drugiej.
    const [wykazAutomatyk, wykazHarmonogramow] = await Promise.all([
      zrodlo.automatyki(),
      zrodlo.harmonogramy(),
    ]);
    if (!wykazAutomatyk.udany || wykazAutomatyk.wynik === undefined) {
      okno.blad(
        opisOdmowy('Odczyt automatyk', wykazAutomatyk.blad?.code, wykazAutomatyk.blad?.message),
      );
      return;
    }
    // Odmowa odczytu harmonogramów nie zabiera wykazu automatyk: bez nich
    // wiersze mówią „rdzeń nie oddał harmonogramu", zamiast znikać.
    const harmonogramy = wykazHarmonogramow.wynik?.schedules ?? [];
    if (!wykazHarmonogramow.udany) {
      odpowiedz.pokaz(
        opisOdmowy(
          'Odczyt harmonogramów',
          wykazHarmonogramow.blad?.code,
          wykazHarmonogramow.blad?.message,
        ),
        false,
      );
    }

    const automatyki = wykazAutomatyk.wynik.workflows;
    ustawPozycje(automatyka, [
      { wartosc: '', etykieta: 'Automatyka: nie wskazano' },
      ...automatyki.map((pozycja) => ({
        wartosc: pozycja.id,
        etykieta: `Automatyka: ${pozycja.name}`,
      })),
    ]);
    lista.replaceChildren(
      ...automatyki.map((pozycja) =>
        wierszRutyny(
          pozycja,
          harmonogramy.filter((harmonogram) => harmonogram.workflowId === pozycja.id),
          () => wczytajDoFormularza(pozycja, harmonogramy),
        ),
      ),
    );
    if (automatyki.length === 0) {
      okno.puste(PUSTE.rutyny);
      return;
    }
    okno.gotowe();
  }

  /** „Edytuj" wczytuje automatykę i jej pierwszy harmonogram do formularza. */
  function wczytajDoFormularza(
    pozycja: AutomationWorkflow,
    harmonogramy: readonly AutomationSchedule[],
  ): void {
    automatyka.value = pozycja.id;
    const harmonogram = harmonogramy.find((wpis) => wpis.workflowId === pozycja.id);
    cyklicznosc.value = harmonogram?.cron ?? '';
    strefa.value = harmonogram?.timeZone ?? '';
    obowiazuje.checked = harmonogram?.enabled ?? true;
    odpowiedz.pokaz(
      harmonogram === undefined
        ? `Automatyka „${pozycja.name}" wczytana do formularza. Rdzeń nie ma dla niej ` +
            'harmonogramu — zapis go założy.'
        : `Automatyka „${pozycja.name}" wczytana wraz z harmonogramem.`,
      true,
    );
  }

  return { element, wczytaj };
}

/** Jeden wiersz wykazu rutyn: automatyka wraz z jej harmonogramami. */
function wierszRutyny(
  pozycja: AutomationWorkflow,
  harmonogramy: readonly AutomationSchedule[],
  naEdycje: () => void,
): HTMLElement {
  const czesci = [
    pozycja.description ?? '',
    pozycja.enabled ? 'automatyka czynna' : 'automatyka wyłączona',
    `kroki: ${String(pozycja.steps?.length ?? 0)}`,
    harmonogramy.length === 0
      ? 'bez harmonogramu — biegnie wyłącznie na żądanie'
      : harmonogramy.map(opisHarmonogramu).join(' · '),
  ];

  const { element, akcje } = pozycjaWykazu(
    pozycja.name,
    czesci.filter((czesc) => czesc !== '').join(' · '),
    'ma',
  );
  element.dataset['automatyka'] = pozycja.id;

  const edytuj = przyciskAkcji('Wczytaj do formularza', 'dn-btn dn-btn--sm dn-btn--zarys');
  edytuj.addEventListener('click', naEdycje);
  akcje.append(edytuj);
  return element;
}

/** Harmonogram słowami: cykliczność, obowiązywanie i najbliższy przebieg. */
function opisHarmonogramu(harmonogram: AutomationSchedule): string {
  const czesci = [
    harmonogram.cron !== undefined && harmonogram.cron !== ''
      ? `cykliczność: ${harmonogram.cron}`
      : 'bez cykliczności',
    harmonogram.enabled ? 'obowiązuje' : 'wstrzymany',
  ];
  if (harmonogram.timeZone !== undefined && harmonogram.timeZone !== '') {
    czesci.push(`strefa: ${harmonogram.timeZone}`);
  }
  if (harmonogram.nextRunAt !== undefined) {
    czesci.push(`najbliższy przebieg: ${chwila(harmonogram.nextRunAt)}`);
  }
  return czesci.join(', ');
}
