import { SubagentStatus, type Subagent } from '../../../../shared/contract';
import { przyciskAkcji } from '../../modele/kontrolki-formularza';
import { czasOdcinka, odmienAgentow, zapiszCzasOdcinka } from './format-zadan';
import { ODZNAKA, StanPrzeplywu, type Przeplyw } from './zadania-w-tle';
import { wierszEtapu, type TelemetriaOkien } from './wskaznik-etapu';

/**
 * Karta jednego przepływu zadań w tle: nazwa z odznaką stanu, wiersz
 * podsumowania, wiersz etapu i tabela agentów.
 *
 * Kolumny „Tok." (żetony) i „Narz." (wywołania narzędzi) stoją puste: struktura
 * `Subagent` w kontrakcie nie niesie ani licznika żetonów, ani licznika
 * wywołań narzędzi. Komórki dostają znak braku wraz z powodem zamiast zera,
 * które czytałoby się jako wykonany pomiar. Z tego samego powodu wiersz
 * podsumowania nie sumuje żetonów.
 */

/** Powód pustej kolumny żetonów — jedno zdanie na cały panel. */
export const BRAK_ZETONOW =
  'Kontrakt nie niesie liczby żetonów: struktura Subagent nie ma takiego pola, więc rdzeń tej liczby nie oddaje. Zero byłoby pomiarem, którego nikt nie wykonał.';

/** Powód pustej kolumny wywołań narzędzi. */
export const BRAK_NARZEDZI =
  'Kontrakt nie niesie liczby wywołań narzędzi: struktura Subagent nie ma takiego pola, więc rdzeń tej liczby nie oddaje.';

/** Znak stawiany w komórce bez pokrycia w kontrakcie. */
const BEZ_POKRYCIA = '—';

/** Odznaki stanu wiążą się z wariantami plakietki biblioteki `dn-*`. */
const WARIANT_ODZNAKI: Readonly<Record<StanPrzeplywu, string>> = {
  [StanPrzeplywu.Oczekuje]: 'dn-plakietka',
  [StanPrzeplywu.WToku]: 'dn-plakietka dn-plakietka--sygnal',
  [StanPrzeplywu.Zakonczone]: 'dn-plakietka dn-plakietka--sukces',
  [StanPrzeplywu.Blad]: 'dn-plakietka dn-plakietka--blad',
  [StanPrzeplywu.Zatrzymane]: 'dn-plakietka dn-plakietka--ostrzezenie',
};

/** Czynności karty; zbieranie wyników woła `subagent.result.collect`. */
export interface CzynnosciKarty {
  zbierz(przeplyw: Przeplyw): void;
}

/** Buduje kartę przepływu gotową do wstawienia w wykaz panelu. */
export function kartaPrzeplywu(
  przeplyw: Przeplyw,
  telemetria: TelemetriaOkien,
  teraz: number,
  czynnosci: CzynnosciKarty,
): HTMLElement {
  const element = document.createElement('li');
  element.className = 'dm-przeplyw';
  element.dataset['stan'] = przeplyw.stan;
  element.dataset['okno'] = przeplyw.idOkna;

  element.append(
    naglowek(przeplyw, czynnosci),
    podsumowanie(przeplyw),
    wierszEtapu(przeplyw.idOkna, telemetria),
    tabelaAgentow(przeplyw.podagenci, teraz),
  );
  return element;
}

/**
 * Nazwa przepływu, odznaka stanu i zbieranie wyników.
 *
 * Przycisk zbierania pozostaje czynny także dla przepływu w toku: komenda
 * `subagent.result.collect` niesie pole `waitForAll`, więc zbieranie przed
 * końcem pracy jest czynnością przewidzianą przez kontrakt.
 */
function naglowek(przeplyw: Przeplyw, czynnosci: CzynnosciKarty): HTMLElement {
  const nazwa = document.createElement('strong');
  nazwa.className = 'dm-przeplyw__nazwa';
  nazwa.textContent = przeplyw.nazwa === '' ? 'zadanie bez treści w kontrakcie' : przeplyw.nazwa;
  nazwa.title = przeplyw.nazwa;

  const odznaka = document.createElement('span');
  odznaka.className = WARIANT_ODZNAKI[przeplyw.stan];
  odznaka.dataset['odznaka'] = przeplyw.stan;
  odznaka.textContent = ODZNAKA[przeplyw.stan];

  const zbierz = przyciskAkcji('Zbierz wyniki', 'dn-btn dn-btn--sm');
  zbierz.addEventListener('click', () => czynnosci.zbierz(przeplyw));

  const element = document.createElement('div');
  element.className = 'dm-przeplyw__naglowek';
  element.append(nazwa, odznaka, zbierz);
  return element;
}

/**
 * Wiersz podsumowania: czas · liczba agentów · żetony.
 *
 * Trzeci człon nie jest liczbą, tylko zdaniem o jej braku, dopóki kontrakt
 * żetonów nie niesie.
 */
function podsumowanie(przeplyw: Przeplyw): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dm-przeplyw__podsumowanie';

  const czas = document.createElement('span');
  czas.textContent = zapiszCzasOdcinka(przeplyw.czas);

  const agenci = document.createElement('span');
  agenci.textContent = odmienAgentow(przeplyw.podagenci.length);

  const zetony = document.createElement('span');
  zetony.className = 'dm-bez-pokrycia';
  zetony.textContent = 'żetony: bez pokrycia w kontrakcie';
  zetony.title = BRAK_ZETONOW;
  zetony.setAttribute('aria-description', BRAK_ZETONOW);

  element.append(czas, rozdzielnik(), agenci, rozdzielnik(), zetony);
  return element;
}

/** Kropka rozdzielająca człony podsumowania; sam znak należy do widoku. */
function rozdzielnik(): HTMLElement {
  const element = document.createElement('span');
  element.className = 'dm-rozdzielnik';
  element.setAttribute('aria-hidden', 'true');
  element.textContent = '·';
  return element;
}

/** Tabela agentów przepływu: Agent · Tok. · Narz. · Czas. */
function tabelaAgentow(podagenci: readonly Subagent[], teraz: number): HTMLElement {
  const tabela = document.createElement('table');
  tabela.className = 'dn-tabela dm-agenci';

  const glowa = document.createElement('thead');
  const wierszGlowy = document.createElement('tr');
  wierszGlowy.append(
    komorkaGlowy('Agent', ''),
    komorkaGlowy('Tok.', BRAK_ZETONOW),
    komorkaGlowy('Narz.', BRAK_NARZEDZI),
    komorkaGlowy('Czas', ''),
  );
  glowa.append(wierszGlowy);

  const cialo = document.createElement('tbody');
  for (const podagent of podagenci) cialo.append(wierszAgenta(podagent, teraz));

  tabela.append(glowa, cialo);
  return tabela;
}

/** Komórka nagłówka; kolumna bez pokrycia niesie powód trzema drogami. */
function komorkaGlowy(napis: string, powod: string): HTMLTableCellElement {
  const komorka = document.createElement('th');
  komorka.scope = 'col';
  komorka.textContent = napis;
  if (powod !== '') {
    komorka.dataset['bezPokrycia'] = 'tak';
    komorka.title = powod;
    komorka.setAttribute('aria-description', powod);
  }
  return komorka;
}

/** Wiersz jednego podagenta. */
function wierszAgenta(podagent: Subagent, teraz: number): HTMLTableRowElement {
  const wiersz = document.createElement('tr');
  wiersz.dataset['podagent'] = podagent.id;
  wiersz.dataset['stan'] = podagent.status;

  const czas = czasOdcinka(
    podagent.startedAt,
    podagent.status === SubagentStatus.Running ? undefined : podagent.finishedAt,
    teraz,
  );

  wiersz.append(
    komorka(podagent.name !== undefined && podagent.name !== '' ? podagent.name : podagent.id, ''),
    komorka(BEZ_POKRYCIA, BRAK_ZETONOW),
    komorka(BEZ_POKRYCIA, BRAK_NARZEDZI),
    komorka(zapiszCzasOdcinka(czas), ''),
  );
  return wiersz;
}

/** Komórka danych; bez pokrycia niesie powód, a nie liczbę. */
function komorka(napis: string, powod: string): HTMLTableCellElement {
  const element = document.createElement('td');
  element.className = 'dn-dane';
  element.textContent = napis;
  if (powod !== '') {
    element.classList.add('dm-bez-pokrycia');
    element.dataset['bezPokrycia'] = 'tak';
    element.title = powod;
    element.setAttribute('aria-description', powod);
  }
  return element;
}
