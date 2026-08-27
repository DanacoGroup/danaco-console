/**
 * Rysunki okna Execution Monitor: panel metryk, oś czasu przebiegów, dziennik
 * i szuflada szczegółów błędu. Wszystkie powstają z przebiegów oddanych przez
 * rdzeń, więc plik nie zna źródła ani stanu okna i daje się sprawdzić bez kanału.
 */

import { AutomationExecutionStatus, type AutomationExecution } from '../../../../shared/contract';
import { pozycjaWykazu, wykaz } from '../../modele/kontrolki-formularza';
import { zapisWpisu, type WpisDziennika } from './dziennik-przebiegow';
import { nazwaStanu } from './maszyna-stanow';
import { czasTrwania, miaryPrzebiegow } from './przeglad-przebiegow';

/** Wysokość paska osi czasu w jednostkach rysunku SVG, wspólna dla wszystkich pasków wykazu przebiegów. */
const WYSOKOSC_PASKA = 12;
const ODSTEP_PASKA = 4;
const SZEROKOSC_OSI = 480;
const PRZESTRZEN_SVG = 'http://www.w3.org/2000/svg';

/**
 * Panel metryk niezawodności. Koszt liczy się z pól przebiegu: sumuje liczbę
 * tokenów i wywołań modelu tam, gdzie rdzeń je podał. Pola są nieobowiązkowe,
 * więc panel podaje też liczbę przebiegów, których suma dotyczy.
 */
export function panelMetryk(przebiegi: readonly AutomationExecution[]): HTMLElement {
  const miary = miaryPrzebiegow(przebiegi);
  const lista = wykaz('Metryki przebiegów automatyki', 'da-wykaz');
  const wiersze: ReadonlyArray<[string, string]> = [
    ['Przebiegów w wykazie', String(miary.liczba)],
    ['Zakończonych powodzeniem', String(miary.udane)],
    ['Zakończonych błędem', String(miary.nieudane)],
    ['W toku', String(miary.wToku)],
    [
      'Wskaźnik powodzenia',
      miary.wskaznikPowodzenia === null
        ? 'brak przebiegów zamkniętych, więc nie ma z czego policzyć'
        : `${miary.wskaznikPowodzenia}% przebiegów zamkniętych`,
    ],
    [
      'Średni czas trwania',
      miary.sredniCzasMs === null ? 'brak przebiegów zamkniętych' : czasTrwania(miary.sredniCzasMs),
    ],
    [
      'Opóźnienie p95',
      miary.opoznienieP95Ms === null
        ? 'brak przebiegów zamkniętych'
        : czasTrwania(miary.opoznienieP95Ms),
    ],
    ['Koszt: tokeny', opisKosztu(przebiegi, (przebieg) => przebieg.tokensUsed)],
    ['Koszt: wywołania modelu', opisKosztu(przebiegi, (przebieg) => przebieg.modelCalls)],
  ];
  for (const [nazwa, wartosc] of wiersze) {
    lista.append(pozycjaWykazu(nazwa, wartosc, 'da').element);
  }
  return lista;
}

/** Suma pola kosztu wraz z liczbą przebiegów, które faktycznie je podały, ponieważ pole bywa nieobowiązkowe i puste. */
function opisKosztu(
  przebiegi: readonly AutomationExecution[],
  wartosc: (przebieg: AutomationExecution) => number | undefined,
): string {
  const podane = przebiegi.map(wartosc).filter((liczba): liczba is number => liczba !== undefined);
  if (podane.length === 0) return 'żaden przebieg wykazu nie podał tej wartości';
  const suma = podane.reduce((razem, liczba) => razem + liczba, 0);
  if (podane.length === przebiegi.length) return `${suma} we wszystkich ${przebiegi.length} przebiegach`;
  return `${suma} w ${podane.length} z ${przebiegi.length} przebiegów, które tę wartość podały`;
}

/**
 * Oś czasu przebiegów: jeden pasek na przebieg, długość według czasu trwania.
 * Przebieg trwający nie ma czasu zakończenia, więc jego pasek sięga chwili
 * bieżącej. Skala jest wspólna dla wykazu, żeby paski dało się porównać wzrokiem.
 */
export function osCzasuPrzebiegow(
  przebiegi: readonly AutomationExecution[],
  teraz: number = Date.now(),
): SVGElement {
  const uporzadkowane = [...przebiegi].sort((pierwszy, drugi) => pierwszy.startedAt - drugi.startedAt);
  const czasy = uporzadkowane.map((przebieg) => (przebieg.finishedAt ?? teraz) - przebieg.startedAt);
  const najdluzszy = czasy.reduce((ile, czas) => Math.max(ile, czas), 0);
  const wysokosc = Math.max(
    WYSOKOSC_PASKA,
    uporzadkowane.length * (WYSOKOSC_PASKA + ODSTEP_PASKA),
  );

  const rysunek = document.createElementNS(PRZESTRZEN_SVG, 'svg');
  rysunek.setAttribute('class', 'da-os-czasu');
  rysunek.setAttribute('viewBox', `0 0 ${SZEROKOSC_OSI} ${wysokosc}`);
  rysunek.setAttribute('role', 'img');
  rysunek.setAttribute(
    'aria-label',
    `Oś czasu ${uporzadkowane.length} przebiegów; najdłuższy trwał ${czasTrwania(najdluzszy)}`,
  );

  uporzadkowane.forEach((przebieg, numer) => {
    const czas = czasy[numer] ?? 0;
    const szerokosc = najdluzszy === 0 ? SZEROKOSC_OSI : Math.max(2, (czas / najdluzszy) * SZEROKOSC_OSI);
    const pasek = document.createElementNS(PRZESTRZEN_SVG, 'rect');
    pasek.setAttribute('x', '0');
    pasek.setAttribute('y', String(numer * (WYSOKOSC_PASKA + ODSTEP_PASKA)));
    pasek.setAttribute('width', String(szerokosc));
    pasek.setAttribute('height', String(WYSOKOSC_PASKA));
    pasek.setAttribute('rx', '2');
    pasek.dataset['stanPrzebiegu'] = przebieg.status;

    const podpis = document.createElementNS(PRZESTRZEN_SVG, 'title');
    podpis.textContent =
      `${przebieg.id} — ${nazwaStanu(przebieg.status)}, ` +
      `${new Date(przebieg.startedAt).toLocaleString('pl-PL')}, ${czasTrwania(czas)}`;
    pasek.append(podpis);
    rysunek.append(pasek);
  });

  return rysunek;
}

/**
 * Szuflada szczegółów błędu — przebiegi zakończone błędem wraz z powodem,
 * krokiem załamania i numerem obiegu; wykaz pusty oddaje wartość pustą
 * zamiast elementu szuflady.
 */
export function szufladaBledow(przebiegi: readonly AutomationExecution[]): HTMLElement | null {
  const nieudane = przebiegi.filter(
    (przebieg) => przebieg.status === AutomationExecutionStatus.Failed,
  );
  if (nieudane.length === 0) return null;
  const szuflada = document.createElement('details');
  szuflada.className = 'da-szuflada';
  const naglowek = document.createElement('summary');
  naglowek.textContent = `Szczegóły błędów (${nieudane.length})`;
  const lista = wykaz('Przebiegi zakończone błędem', 'da-wykaz');
  for (const przebieg of nieudane) {
    // Krok załamania jest dokładniejszy niż etap bieżący, więc ma pierwszeństwo w opisie.
    const etap =
      przebieg.failedStepId !== undefined && przebieg.failedStepId !== ''
        ? `krok ${przebieg.failedStepId}`
        : przebieg.currentStep === undefined
          ? 'rdzeń nie wskazał kroku'
          : `etap ${przebieg.currentStep}`;
    const pozycja = pozycjaWykazu(
      przebieg.id,
      `${przebieg.errorMessage ?? 'rdzeń nie podał powodu'} — ${etap}, obieg ${przebieg.attempt ?? 0}`,
      'da',
    );
    pozycja.element.dataset['stanPrzebiegu'] = przebieg.status;
    lista.append(pozycja.element);
  }
  szuflada.append(naglowek, lista);
  return szuflada;
}

/** Wykaz wierszy dziennika przebiegów w postaci elementu strony; wykaz pusty oddaje wartość pustą zamiast elementu. */
export function wykazDziennika(wpisy: readonly WpisDziennika[]): HTMLElement | null {
  if (wpisy.length === 0) return null;
  const lista = wykaz('Dziennik przebiegów', 'da-wykaz');
  for (const wpis of wpisy) {
    const pozycja = pozycjaWykazu(
      new Date(wpis.czas).toLocaleTimeString('pl-PL'),
      wpis.etap === '' ? wpis.tresc : `${wpis.etap}: ${wpis.tresc}`,
      'da',
    );
    pozycja.element.dataset['stanPrzebiegu'] = wpis.stan;
    lista.append(pozycja.element);
  }
  return lista;
}

/** Treść pliku dziennika przebiegów — jeden wiersz tekstu na każdy wpis, w kolejności jego przyjęcia do dziennika. */
export function zapisDziennika(wpisy: readonly WpisDziennika[]): string {
  return `${wpisy.map(zapisWpisu).join('\n')}\n`;
}

/**
 * Zestawienie przebiegów w zapisie rozdzielanym średnikiem.
 *
 * Średnik, nie przecinek: polskie arkusze kalkulacyjne czytają przecinek jako
 * znak dziesiętny, więc plik z przecinkiem rozjeżdża się w kolumnach.
 */
export function zestawieniePrzebiegow(przebiegi: readonly AutomationExecution[]): string {
  const wiersze = [
    ['przebieg', 'automatyka', 'stan', 'etap', 'etapow', 'obieg', 'rozpoczeto', 'zakonczono', 'powod'].join(';'),
  ];
  for (const przebieg of przebiegi) {
    wiersze.push(
      [
        przebieg.id,
        przebieg.workflowId,
        przebieg.status,
        String(przebieg.currentStep ?? ''),
        String(przebieg.totalSteps ?? ''),
        String(przebieg.attempt ?? ''),
        new Date(przebieg.startedAt).toISOString(),
        przebieg.finishedAt === undefined ? '' : new Date(przebieg.finishedAt).toISOString(),
        (przebieg.errorMessage ?? '').replace(/;/g, ','),
      ].join(';'),
    );
  }
  return `${wiersze.join('\n')}\n`;
}
