import { AssistantActionControl, type AssistantAction } from '../../../../shared/contract';
import { utworzNaglowekOkna } from '../../komponenty/naglowek-okna';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  przycisk,
  utworzWierszOdpowiedzi,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { BRAKI, zglosBrak } from './braki-kontraktu';
import { NAZWY_STANOW, ODCZYTY, PUSTE, SPOZA_OKNA } from './etykiety-assistant';
import { utworzFiltrZlecen, type FiltrZlecen } from './filtr-zlecen';
import { opisSkutku } from './skutek-sterowania';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAssistant } from './stan-assistant';
import { wierszZlecenia } from './wiersz-zlecenia';

/**
 * Kod okna w katalogu okien operacyjnych rdzenia, pole `okno_operacyjne.kod`,
 * po którym rdzeń rozpoznaje to okno w zleceniach i w zapisie czynności.
 */
export const KOD_OKNA = 'actions-monitor';

/**
 * Actions Monitor to okno monitorujące modułu Assistant: przegląd statusu
 * zlecenia wieloetapowego oraz wstrzymanie albo anulowanie działania. Panel
 * akcji i priorytet mieszczą się w jednej komendzie stanu zlecenia.
 */
export interface OknoActionsMonitor {
  element: HTMLElement;
  odswiez(): void;
}

/**
 * Nagłówki kolumn tabeli zleceń; stoją niezależnie od zawartości, więc tabela
 * pozostaje czytelna także przy pustym wykazie i przy odczycie w toku.
 */
const KOLUMNY = ['Zlecenie', 'Stan', 'Etap', 'Priorytet', 'Panel akcji'] as const;

export function utworzOknoActionsMonitor(stan: StanAssistant): OknoActionsMonitor {
  const okno: StanOkna = utworzStanOkna();
  const odpowiedz: WierszOdpowiedzi = utworzWierszOdpowiedzi();
  const filtr: FiltrZlecen = utworzFiltrZlecen(() => odswiez());

  const wlasne = utworzTabeleZlecen();
  const spoza = utworzTabeleZlecen();

  const sekcjaSpoza = document.createElement('section');
  sekcjaSpoza.className = 'ma-spoza';
  sekcjaSpoza.hidden = true;
  sekcjaSpoza.append(tytulSpoza(), zasiegSpoza(), spoza.tabela);

  const szczegoly = document.createElement('pre');
  szczegoly.className = 'ma-zlecenia__szczegoly';
  szczegoly.hidden = true;

  okno.tresc.append(wlasne.tabela, sekcjaSpoza, odpowiedz.element, szczegoly);

  const element = document.createElement('section');
  element.className = 'ma-okno ma-okno--monitor';
  element.dataset['okno'] = KOD_OKNA;
  element.append(
    utworzNaglowekOkna({
      tytul: 'Actions Monitor',
      rola: 'monitor · zlecenia wieloetapowe asystenta na żywo',
      kontrolki: [filtr.element],
    }),
    okno.element,
    stopka(),
  );

  async function steruj(
    idZlecenia: string,
    sterowanie: AssistantActionControl,
    priorytet?: number,
  ): Promise<void> {
    odpowiedz.pokaz('Rdzeń przyjmuje sterowanie zleceniem…', true);
    const wynik = await stan.zrodlo.steruj({ idZlecenia, sterowanie, priorytet });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Sterowanie zleceniem', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    stan.wchlon(wynik.wynik.actions);
    // Ocena dostaje zamówienie w całości, bo porównuje je ze zleceniem zwróconym.
    const skutek = opisSkutku(wynik.wynik.actions, idZlecenia, sterowanie, priorytet);
    odpowiedz.pokaz(skutek.zdanie, skutek.udany);
  }

  function pokaz(zlecenie: AssistantAction, wynik: boolean): void {
    szczegoly.hidden = false;
    szczegoly.textContent = wynik ? opisWyniku(zlecenie) : opisZlecenia(zlecenie);
    void stan.wybierz(zlecenie.id);
  }

  /** Jeden wiersz — ten sam skład i te same sterowania w obu tabelach. */
  const rysuj = (zlecenie: AssistantAction): HTMLTableRowElement =>
    wierszZlecenia(
      zlecenie,
      (id, sterowanie, priorytet) => void steruj(id, sterowanie, priorytet),
      pokaz,
    );

  function odswiez(): void {
    // Filtr zawęża oba wykazy tym samym warunkiem, bo stany zleceń są wspólne.
    const zlecenia = stan.zlecenia().filter((zlecenie) => filtr.przepusc(zlecenie));
    const obce = stan.zleceniaObce().filter((zlecenie) => filtr.przepusc(zlecenie));
    wlasne.cialo.replaceChildren(...zlecenia.map(rysuj));
    spoza.cialo.replaceChildren(...obce.map(rysuj));
    // Sekcja bez wierszy nie ma stanu własnego; pustkę wykazu nazywa okno.
    sekcjaSpoza.hidden = obce.length === 0;
    naniesFaze(okno, stan, zlecenia.length + obce.length, filtr);
  }

  odswiez();
  return { element, odswiez };
}

/**
 * Stopka okna z bramą potwierdzeń zleceń wymagających decyzji Operatora,
 * dla której kontrakt nie niesie osobnej komendy przekazania.
 */
function stopka(): HTMLElement {
  const brama = przycisk('Brama potwierdzeń akcji', 'dn-btn dn-btn--sm dn-btn--duch');
  brama.dataset['brak'] = 'brama-potwierdzen';
  brama.addEventListener('click', () => zglosBrak('Brama potwierdzeń', BRAKI.bramaPotwierdzen));

  const element = document.createElement('div');
  element.className = 'ma-okno__stopka';
  element.append(brama);
  return element;
}

/**
 * Nanosi fazę odczytu na tabelę zleceń: ciało tabeli znika na czas pytania
 * do rdzenia, a nagłówek kolumn zostaje, żeby układ okna nie skakał.
 */
function naniesFaze(okno: StanOkna, stan: StanAssistant, ile: number, filtr: FiltrZlecen): void {
  const faza = stan.faza();
  if (faza === 'blad') {
    okno.blad(stan.powod());
    return;
  }
  if (faza === 'ladowanie') {
    okno.ladowanie(ODCZYTY.zlecenia);
    return;
  }
  if (ile === 0) {
    // Pustka przed pytaniem, po odpowiedzi i po zawężeniu to trzy różne zdania.
    if (stan.pytanoOZlecenia() && filtr.zawezony()) {
      okno.puste(
        `Wykaz jest zawężony do zleceń o statusie: ${filtr.nazwa()}. ` +
          'Żadne zlecenie nie ma dziś tego statusu — przestaw filtr na „wszystkie", ' +
          'aby zobaczyć wykaz w całości.',
      );
      return;
    }
    okno.puste(stan.pytanoOZlecenia() ? PUSTE.zlecenia : PUSTE.monitorSpoczynek);
    return;
  }
  okno.gotowe();
}

/**
 * Buduje tabelę zleceń wraz z jej ciałem; ta sama budowa obsługuje wykaz
 * zleceń własnych okna oraz wykaz zleceń pochodzących spoza niego.
 */
function utworzTabeleZlecen(): { tabela: HTMLTableElement; cialo: HTMLTableSectionElement } {
  const cialo = document.createElement('tbody');
  const tabela = document.createElement('table');
  tabela.className = 'dn-tabela ma-zlecenia';
  tabela.append(naglowekTabeli(), cialo);
  return { tabela, cialo };
}

/**
 * Tytuł sekcji zleceń pochodzących spoza okna modułu, o stopień niższy niż
 * nazwa samego okna, żeby porządek nagłówków strony pozostał zachowany.
 */
function tytulSpoza(): HTMLElement {
  const element = document.createElement('h4');
  element.className = 'dn-karta-tytul ma-spoza__tytul';
  element.textContent = SPOZA_OKNA.tytul;
  return element;
}

/**
 * Zdanie o zasięgu wykazu zleceń, nazywające zarówno to, co wykaz obejmuje,
 * jak i to, czego w nim nie ma, żeby Operator nie brał braku za pustkę.
 */
function zasiegSpoza(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis ma-spoza__zasieg';
  element.textContent = SPOZA_OKNA.zasieg;
  return element;
}

function naglowekTabeli(): HTMLTableSectionElement {
  const wiersz = document.createElement('tr');
  for (const nazwa of KOLUMNY) {
    const komorka = document.createElement('th');
    komorka.scope = 'col';
    komorka.textContent = nazwa;
    wiersz.append(komorka);
  }
  const naglowek = document.createElement('thead');
  naglowek.append(wiersz);
  return naglowek;
}

/**
 * Rozwinięcie szczegółów zlecenia: podaje wyłącznie to, co rdzeń o zleceniu
 * powiedział, bez uzupełniania brakujących pól domysłem po stronie okna.
 */
function opisZlecenia(zlecenie: AssistantAction): string {
  return [
    `Zlecenie: ${zlecenie.title ?? zlecenie.id}`,
    `Identyfikator: ${zlecenie.id}`,
    `Okno: ${zlecenie.windowId}`,
    `Stan: ${NAZWY_STANOW[zlecenie.status]}`,
    `Priorytet: ${String(zlecenie.priority ?? 0)}`,
    `Założone: ${new Date(zlecenie.createdAt).toLocaleString('pl-PL')}`,
    `Zmienione: ${new Date(zlecenie.updatedAt).toLocaleString('pl-PL')}`,
  ].join('\n');
}

/**
 * Odsłonięcie wyniku zlecenia z pola `result` odpowiedzi rdzenia; brak wyniku
 * okno nazywa wprost, zamiast podstawiać treść zastępczą.
 */
function opisWyniku(zlecenie: AssistantAction): string {
  if (zlecenie.result === undefined || zlecenie.result === '') {
    return 'Rdzeń nie oddał wyniku tego zlecenia (pole result jest puste).';
  }
  return zlecenie.result;
}
