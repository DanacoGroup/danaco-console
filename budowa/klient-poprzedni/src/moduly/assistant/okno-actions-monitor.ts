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

/** Kod okna w katalogu rdzenia (`okno_operacyjne.kod`). */
export const KOD_OKNA = 'actions-monitor';

/**
 * Actions Monitor — okno monitorujące modułu Assistant.
 *
 * Dwie funkcje operatora: przegląd statusu zlecenia wieloetapowego oraz
 * wstrzymanie lub anulowanie działania. Panel akcji niesie sześć pozycji —
 * wstrzymanie, wznowienie, anulowanie, Ponów, Szczegóły, Pokaż wynik — i
 * priorytetyzację zleceń; wszystkie mieszczą się w jednej komendzie
 * `assistant.action.status` wraz z polami `control` i `priority`.
 *
 * Okno nie ma pętli odświeżania: `assistant.action.changed` wciąga zmianę
 * stanu zlecenia w chwili, w której rdzeń ją ogłasza. Subskrypcja mieszka
 * w stanie modułu, więc zdarzenie zmienia wszystkie trzy okna naraz.
 *
 * Zdanie potwierdzenia mówi, co zrobił rdzeń, a nie co wysłało okno: składa je
 * `skutek-sterowania.ts` ze zlecenia, które wróciło, bo nie każde przyjęte
 * żądanie coś zmienia. Sprawdzane są wszystkie zamówienia panelu — cztery
 * przyciski stanu tak samo jak priorytet — inaczej przycisk meldowałby sukces
 * także wtedy, gdy rdzeń oddał `failed`.
 *
 * Nagłówek kolumn zostaje w stanie pustym: pustka dotyczy ciała tabeli, nie
 * jej budowy.
 *
 * Druga tabela zbiera zlecenia spoza tego okna. Praca asystenta wydana
 * w nakładce AOD siada na oknie, które nakładka dobiera sama
 * (`adapter_modul_aod.go`, `oknoZadania`) — zwykle na oknie rozmowy, nie na
 * oknie modułu. Bez drugiej tabeli taka praca byłaby w monitorze niewidzialna,
 * bo stan modułu odrzuca zdarzenia o cudzym `windowId`. Tabela ma ten sam
 * panel akcji, bo rdzeń przyjmuje sterowanie po `actionId`, nie po oknie
 * (`adapter_modul_asystent_czynnosci.go`, `steruj`). Sekcja jest ukryta,
 * dopóki nie ma czego pokazać.
 */
export interface OknoActionsMonitor {
  element: HTMLElement;
  odswiez(): void;
}

/** Kolumny tabeli zleceń — nagłówek trwa niezależnie od zawartości. */
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
    // Ocena dostaje zamówienie w całości — sterowanie i priorytet — bo tylko
    // mając jedno i drugie może porównać je ze zleceniem, które wróciło.
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
    // Filtr zawęża oba wykazy tym samym warunkiem: zlecenie spoza okna modułu
    // przechodzi przez te same stany co własne, więc dwa różne zawężenia
    // znaczyłyby dwie prawdy o jednym filtrze.
    const zlecenia = stan.zlecenia().filter((zlecenie) => filtr.przepusc(zlecenie));
    const obce = stan.zleceniaObce().filter((zlecenie) => filtr.przepusc(zlecenie));
    wlasne.cialo.replaceChildren(...zlecenia.map(rysuj));
    spoza.cialo.replaceChildren(...obce.map(rysuj));
    // Sekcja bez wierszy nie stoi pusta: nie ma stanu własnego i nie miałaby go
    // czym nazwać. Pustkę całego wykazu nazywa okno.
    sekcjaSpoza.hidden = obce.length === 0;
    naniesFaze(okno, stan, zlecenia.length + obce.length, filtr);
  }

  odswiez();
  return { element, odswiez };
}

/** Stopka okna: brama potwierdzeń, dla której kontrakt nie ma drogi. */
function stopka(): HTMLElement {
  const brama = przycisk('Brama potwierdzeń akcji', 'dn-btn dn-btn--sm dn-btn--duch');
  brama.dataset['brak'] = 'brama-potwierdzen';
  brama.addEventListener('click', () => zglosBrak('Brama potwierdzeń', BRAKI.bramaPotwierdzen));

  const element = document.createElement('div');
  element.className = 'ma-okno__stopka';
  element.append(brama);
  return element;
}

/** Nanosi fazę odczytu; ciało tabeli znika, nagłówek kolumn zostaje. */
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
    // Pustka przed pytaniem, pustka po odpowiedzi i pustka po zawężeniu to trzy
    // różne zdania: okno nie orzeka o zleceniach rdzenia, zanim rdzeń się o nich
    // wypowie, i nie mówi „rdzeń nie prowadzi ani jednego", gdy wykaz ukrył
    // własny filtr okna. Każde zdanie nazywa przy tym samo okno — pustka ma
    // tłumaczyć, przed czym Operator stoi, a nie tylko czego nie ma.
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

/** Tabela zleceń wraz z własnym ciałem — ta sama budowa dla obu wykazów. */
function utworzTabeleZlecen(): { tabela: HTMLTableElement; cialo: HTMLTableSectionElement } {
  const cialo = document.createElement('tbody');
  const tabela = document.createElement('table');
  tabela.className = 'dn-tabela ma-zlecenia';
  tabela.append(naglowekTabeli(), cialo);
  return { tabela, cialo };
}

/** Tytuł sekcji zleceń spoza okna — stopień niżej niż nazwa okna (h3). */
function tytulSpoza(): HTMLElement {
  const element = document.createElement('h4');
  element.className = 'dn-karta-tytul ma-spoza__tytul';
  element.textContent = SPOZA_OKNA.tytul;
  return element;
}

/** Zdanie o zasięgu wykazu — mówi także, czego w nim nie ma. */
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

/** „Szczegóły" — to, co rdzeń o zleceniu powiedział, bez ani jednego domysłu. */
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

/** „Pokaż wynik" — pole `result`; jego brak nazywamy, nie zmyślamy. */
function opisWyniku(zlecenie: AssistantAction): string {
  if (zlecenie.result === undefined || zlecenie.result === '') {
    return 'Rdzeń nie oddał wyniku tego zlecenia (pole result jest puste).';
  }
  return zlecenie.result;
}
