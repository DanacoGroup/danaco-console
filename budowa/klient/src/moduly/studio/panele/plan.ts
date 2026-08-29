/**
 * Okno robocze Studia — panel Plan, kształt wzięty z
 * `design/05-okna/moduly/studio.html`, `section#panel-plan`.
 *
 * Panel pokazuje rozkład zlecenia dokumentowego na zadania i stan każdego z
 * nich: `studio.plan.get` niesie rozkład, `studio.plan.create` układa nowy,
 * `studio.plan.task.update` przestawia jedno zadanie, a `studio.plan.run` i
 * `studio.plan.stop` puszczają pętlę wykonawczą i ją zatrzymują. Łańcuch
 * operacji jest tą samą pętlą uruchomioną gotową sekwencją, więc
 * `studio.chain.list` i `studio.chain.run` stoją w tej samej strefie, a
 * `studio.chain.progressed` prowadzi ich postęp.
 *
 * Dokument, którego rozkład dotyczy, panel poznaje ze zdarzenia
 * `studio.document.changed` ograniczonego do własnego okna — tak samo jak
 * `tools.ts`, bo kontrakt nie ma komendy oddającej dokument otwarty w oknie.
 *
 * `studio.chain.save` i `studio.batch.run` nie mają w źródle kształtu żadnej
 * powierzchni sterowania: pierwsza układa kroki łańcucha, druga wybiera wiele
 * dokumentów naraz. Panel ich nie woła, bo obie musiałby oprzeć na układzie,
 * którego prototyp nie niesie.
 *
 * Żadna wartość widoczna dla Operatora nie jest tu wpisana z góry: poza
 * odmową i stanem pustym każda pochodzi z odpowiedzi rdzenia.
 */

import {
  ChangeKind,
  Command,
  ErrorCode,
  EventType,
  StudioOperationScope,
  StudioPlanState,
  StudioTaskState,
  type ErrorInfo,
  type StudioChain,
  type StudioDocument,
  type StudioDocumentTask,
  type StudioTaskPlan,
} from '../../../../../shared/contract.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { ikony } from '../ikony.ts';
import { zLiczba } from '../liczebnik.ts';
import { el, zeZnacznika, type Dziecko } from '../narzedzia.ts';
import { opisOdmowy, zaloguj } from '../odmowa.ts';
import { opisPostepuLancucha, planTresci as T } from './plan-tresci.ts';
import type { MontazPanelu, ZaleznosciPanelu, ZamontowanyPanel } from './umowa.ts';

type Zaleznosc = 'okno' | 'sesja';

type Stan =
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'brakZaleznosci'; brak: Zaleznosc }
  | { rodzaj: 'odmowaRozkladu'; blad?: ErrorInfo }
  | { rodzaj: 'brakRozkladu' }
  | { rodzaj: 'rozklad'; plan: StudioTaskPlan };

type Dzialanie =
  | { rodzaj: 'spoczynek' }
  | { rodzaj: 'uruchamianie' }
  | { rodzaj: 'zatrzymywanie' }
  | { rodzaj: 'rozkladanie' }
  | { rodzaj: 'zadanie' }
  | { rodzaj: 'lancuch' }
  | { rodzaj: 'odmowa'; naglowek: string; opis: string };

const WARIANT_ZNACZNIKA: Record<StudioTaskState, string> = {
  [StudioTaskState.Pending]: 'dn-kropka--neutralna',
  [StudioTaskState.Running]: 'dn-kropka--sygnal',
  [StudioTaskState.Done]: 'dn-kropka--sukces',
  [StudioTaskState.Failed]: 'dn-kropka--blad',
  [StudioTaskState.Blocked]: 'dn-kropka--ostrzezenie',
  [StudioTaskState.Skipped]: 'dn-kropka--neutralna',
};

const WARIANT_PLAKIETKI: Partial<Record<StudioPlanState, string>> = {
  [StudioPlanState.Running]: 'dn-plakietka--sygnal',
  [StudioPlanState.Paused]: 'dn-plakietka--ostrzezenie',
  [StudioPlanState.Done]: 'dn-plakietka--sukces',
  [StudioPlanState.Stopped]: 'dn-plakietka--ostrzezenie',
};

function znak(rysunek: keyof typeof ikony): SVGElement {
  const wezelZnaku = zeZnacznika(ikony[rysunek]);
  wezelZnaku.setAttribute('aria-hidden', 'true');
  return wezelZnaku;
}

export const montujPanelPlanu: MontazPanelu = (wezel, zaleznosci: ZaleznosciPanelu): ZamontowanyPanel => {
  let zdjety = false;
  let stan: Stan = { rodzaj: 'ladowanie' };
  let dzialanie: Dzialanie = { rodzaj: 'spoczynek' };
  let dokument: StudioDocument | null = null;
  let rozklady: StudioTaskPlan[] = [];
  let lancuchy: StudioChain[] = [];
  let odmowaLancuchow: ErrorInfo | undefined | null = null;
  let wskazaneZadanie: string | null = null;
  let biegLancucha: string | null = null;
  let postepLancucha: string | null = null;
  let zlecenie = '';

  const znacznik = el('span', { klasa: 'sta-okno-znacznik' });
  const tresc = el('div', { klasa: 'sta-okno-tresc st-panel-lista' });
  wezel.classList.add('sta-okno');
  wezel.replaceChildren(
    el('header', { klasa: 'sta-okno-belka' }, [
      el('span', { klasa: 'sta-okno-tytul' }, [znak('plan'), el('b', { tekst: T.panel.tytul })]),
      znacznik,
    ]),
    tresc,
  );

  /* ── Wywołania ────────────────────────────────────────────────────────── */

  async function wczytajRozklad(idRozkladu?: string): Promise<void> {
    if (zaleznosci.idSesji === null) return odswiez({ rodzaj: 'brakZaleznosci', brak: 'sesja' });
    if (zaleznosci.idOkna === null) return odswiez({ rodzaj: 'brakZaleznosci', brak: 'okno' });

    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioPlanGet, {
      ...(idRozkladu === undefined ? {} : { planId: idRozkladu }),
      ...(dokument === null ? {} : { documentId: dokument.id }),
    });
    if (zdjety) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      /* Brak rozkładu nie jest usterką: dokument, którego nikt jeszcze nie
         rozłożył, wraca tym samym kodem co pozycja nieistniejąca. */
      if (wynik.blad?.code === ErrorCode.NotFound) {
        zaloguj(wynik.blad, 'plan.get');
        rozklady = [];
        return odswiez({ rodzaj: 'brakRozkladu' });
      }
      return odswiez({ rodzaj: 'odmowaRozkladu', blad: wynik.blad });
    }
    rozklady = wynik.wynik.plans ?? [];
    odswiez({ rodzaj: 'rozklad', plan: wynik.wynik.plan });
  }

  async function wczytajLancuchy(): Promise<void> {
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioChainList, {});
    if (zdjety) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odmowaLancuchow = wynik.blad ?? undefined;
      return przerysuj();
    }
    odmowaLancuchow = null;
    lancuchy = wynik.wynik.chains;
    przerysuj();
  }

  async function rozloz(): Promise<void> {
    const otwarty = dokument;
    const trescZlecenia = zlecenie.trim();
    if (otwarty === null || trescZlecenia === '') return;

    dzialanie = { rodzaj: 'rozkladanie' };
    przerysuj();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioPlanCreate, {
      documentId: otwarty.id,
      order: trescZlecenia,
    });
    if (zdjety) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      return zglosOdmowe(T.odmowa.rozlozenie, wynik.blad, 'plan.create');
    }
    zlecenie = '';
    dzialanie = { rodzaj: 'spoczynek' };
    odswiez({ rodzaj: 'rozklad', plan: wynik.wynik.plan });
  }

  async function uruchom(plan: StudioTaskPlan): Promise<void> {
    dzialanie = { rodzaj: 'uruchamianie' };
    przerysuj();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioPlanRun, { planId: plan.id });
    if (zdjety) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      return zglosOdmowe(T.odmowa.uruchomienie, wynik.blad, 'plan.run');
    }
    if (!wynik.wynik.started) {
      /* Powód spisany przez serwer idzie do dziennika: opisuje wnętrze pętli,
         nie czynność, którą Operator ma przed sobą. */
      console.warn('[studio] pętla nie ruszyła', wynik.wynik.refusalReason);
      dzialanie = { rodzaj: 'odmowa', naglowek: T.odmowa.uruchomienie, opis: opisOdmowy(undefined, 'plan.run') };
      return odswiez({ rodzaj: 'rozklad', plan: wynik.wynik.plan });
    }
    dzialanie = { rodzaj: 'spoczynek' };
    odswiez({ rodzaj: 'rozklad', plan: wynik.wynik.plan });
  }

  async function zatrzymaj(plan: StudioTaskPlan): Promise<void> {
    dzialanie = { rodzaj: 'zatrzymywanie' };
    przerysuj();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioPlanStop, { planId: plan.id });
    if (zdjety) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      return zglosOdmowe(T.odmowa.zatrzymanie, wynik.blad, 'plan.stop');
    }
    dzialanie = { rodzaj: 'spoczynek' };
    odswiez({ rodzaj: 'rozklad', plan: wynik.wynik.plan });
  }

  async function przestawZadanie(zadanie: StudioDocumentTask, docelowy: StudioTaskState): Promise<void> {
    dzialanie = { rodzaj: 'zadanie' };
    przerysuj();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioPlanTaskUpdate, {
      taskId: zadanie.id,
      state: docelowy,
    });
    if (zdjety) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      return zglosOdmowe(T.odmowa.zadanie, wynik.blad, 'plan.task.update');
    }
    dzialanie = { rodzaj: 'spoczynek' };
    odswiez({ rodzaj: 'rozklad', plan: wynik.wynik.plan });
  }

  async function uruchomLancuch(lancuch: StudioChain): Promise<void> {
    const otwarty = dokument;
    const idOkna = zaleznosci.idOkna;
    if (otwarty === null || idOkna === null) return;

    dzialanie = { rodzaj: 'lancuch' };
    postepLancucha = null;
    przerysuj();
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioChainRun, {
      windowId: idOkna,
      documentId: otwarty.id,
      chainId: lancuch.id,
      scope: StudioOperationScope.Document,
    });
    if (zdjety) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      return zglosOdmowe(T.odmowa.lancuch, wynik.blad, 'chain.run');
    }
    biegLancucha = wynik.wynik.runId;
    postepLancucha = `${T.lancuchy.przyjeto} ${zLiczba(wynik.wynik.steps, T.lancuchy.jednostkaKroki)}`;
    dzialanie = { rodzaj: 'spoczynek' };
    void wczytajRozklad();
    przerysuj();
  }

  function zglosOdmowe(naglowek: string, blad: ErrorInfo | undefined, czynnosc: string): void {
    dzialanie = { rodzaj: 'odmowa', naglowek, opis: opisOdmowy(blad, czynnosc) };
    przerysuj();
  }

  /* ── Widok ────────────────────────────────────────────────────────────── */

  function alert(naglowek: string, opis: string): HTMLElement {
    return el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
      el('span', { klasa: 'dn-alert-tresc' }, [el('b', { tekst: naglowek }), el('span', { tekst: opis })]),
    ]);
  }

  function wierszPulsu(etykieta: string): HTMLElement {
    return el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
      el('span', { klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno', 'aria-hidden': 'true' }),
      etykieta,
    ]);
  }

  function pustyStan(tytul: string, opis: string): HTMLElement {
    return el('div', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty' }, [
      el('span', { klasa: 'dn-pusty-stan-tytul', tekst: tytul }),
      el('span', { klasa: 'dn-pusty-stan-opis', tekst: opis }),
    ]);
  }

  function wyborRozkladu(biezacy: StudioTaskPlan): HTMLElement | null {
    if (rozklady.length < 2) return null;
    return el('div', { klasa: 'dn-zakladki dn-zakladki--pigulki', role: 'group', 'aria-label': T.rozklady.etykieta }, [
      ...rozklady.map((inny) => {
        const wybrany = inny.id === biezacy.id;
        const guzik = el('button', {
          klasa: wybrany ? 'dn-btn dn-btn--zarys dn-btn--sm' : 'dn-btn dn-btn--duch dn-btn--sm',
          type: 'button',
          'aria-pressed': wybrany ? 'true' : 'false',
          tekst: inny.order === '' ? T.rozklady.bezZlecenia : inny.order,
        });
        guzik.addEventListener('click', () => {
          if (wybrany) return;
          wskazaneZadanie = null;
          void wczytajRozklad(inny.id);
        });
        return guzik;
      }),
    ]);
  }

  function znacznikZadania(stanZadania: StudioTaskState): HTMLElement {
    /* Tętno zamiast kropki tylko dla zadania w biegu — jedyny stan, który sam
       się zmienia, więc jedyny, który ma się ruszać. */
    if (stanZadania === StudioTaskState.Running) return el('span', { klasa: 'pt-tetno', 'aria-hidden': 'true' });
    return el('span', { klasa: `dn-kropka ${WARIANT_ZNACZNIKA[stanZadania]}`, 'aria-hidden': 'true' });
  }

  function wierszZadania(zadanie: StudioDocumentTask): HTMLElement {
    const wskazane = wskazaneZadanie === zadanie.id;
    const opisStanu = T.stanZadania[zadanie.state];
    /* `<div>`, jak w źródle kształtu: `.st-panel-wiersz` nie niesie resetu
       wyglądu natywnego przycisku, więc rolę i klawiaturę dokłada się wprost. */
    const wiersz = el(
      'div',
      {
        klasa: 'st-panel-wiersz',
        role: 'button',
        tabindex: '0',
        'aria-pressed': wskazane ? 'true' : 'false',
      },
      [
        znacznikZadania(zadanie.state),
        zadanie.title,
        el('span', { klasa: 'dn-meta', tekst: opisStanu.glif, title: opisStanu.etykieta }),
      ],
    );
    const przelacz = (): void => {
      wskazaneZadanie = wskazane ? null : zadanie.id;
      przerysuj();
    };
    wiersz.addEventListener('click', przelacz);
    wiersz.addEventListener('keydown', (zdarzenie) => {
      if (zdarzenie.key !== 'Enter' && zdarzenie.key !== ' ') return;
      zdarzenie.preventDefault();
      przelacz();
    });
    return wiersz;
  }

  function czynnosciZadania(zadanie: StudioDocumentTask): HTMLElement[] {
    const zajete = dzialanie.rodzaj === 'zadanie';
    function guzik(etykieta: string, docelowy: StudioTaskState): HTMLElement {
      const przycisk = el('button', {
        klasa: 'dn-btn dn-btn--duch dn-btn--sm',
        type: 'button',
        disabled: zajete || zadanie.state === docelowy,
        tekst: zajete ? T.zadanie.przestawianie : etykieta,
      });
      przycisk.addEventListener('click', () => void przestawZadanie(zadanie, docelowy));
      return przycisk;
    }

    const szczegoly: Dziecko[] = [
      el('span', { klasa: 'dn-meta', tekst: T.stanZadania[zadanie.state].etykieta }),
      zadanie.assignedTo?.agentName === undefined
        ? null
        : el('span', { klasa: 'dn-meta', tekst: `${T.zadanie.wykonawca} ${zadanie.assignedTo.agentName}` }),
      zadanie.failureReason === undefined
        ? null
        : el('span', { klasa: 'dn-meta', tekst: `${T.zadanie.powod} ${zadanie.failureReason}` }),
      zadanie.result === undefined
        ? null
        : el('span', { klasa: 'dn-meta', tekst: `${T.zadanie.wynik} ${zadanie.result}` }),
    ];

    return [
      el('div', { klasa: 'dn-nota' }, szczegoly),
      el('div', { klasa: 'dn-pas-dzialan', role: 'group', 'aria-label': T.zadanie.etykietaDzialan }, [
        guzik(T.zadanie.wykonane, StudioTaskState.Done),
        guzik(T.zadanie.pomin, StudioTaskState.Skipped),
        guzik(T.zadanie.ponow, StudioTaskState.Pending),
      ]),
    ];
  }

  function plakietkaStanu(plan: StudioTaskPlan): HTMLElement {
    const wariant = WARIANT_PLAKIETKI[plan.state];
    return el('span', { klasa: wariant === undefined ? 'dn-plakietka' : `dn-plakietka ${wariant}` }, [
      plan.state === StudioPlanState.Running ? el('span', { klasa: 'pt-tetno', 'aria-hidden': 'true' }) : null,
      `${T.stanPetli} ${T.stanRozkladuSlowem[plan.state]}`,
    ]);
  }

  function przyciskPetli(plan: StudioTaskPlan): HTMLElement | null {
    const zajete = dzialanie.rodzaj === 'uruchamianie' || dzialanie.rodzaj === 'zatrzymywanie';
    if (plan.state === StudioPlanState.Running) {
      const przycisk = el('button', {
        klasa: 'dn-btn dn-btn--niebezpieczny dn-btn--sm',
        type: 'button',
        disabled: zajete,
        tekst: dzialanie.rodzaj === 'zatrzymywanie' ? T.akcje.zatrzymywanie : T.akcje.zatrzymaj,
      });
      przycisk.addEventListener('click', () => void zatrzymaj(plan));
      return przycisk;
    }
    /* Rozkład domknięty nie ma czego puszczać: przycisk uruchomienia nie stoi
       wcale, zamiast stać wyłączony bez powodu widocznego dla Operatora. */
    if (plan.state === StudioPlanState.Done) return null;
    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--atrament dn-btn--sm',
      type: 'button',
      disabled: zajete,
      tekst: dzialanie.rodzaj === 'uruchamianie' ? T.akcje.uruchamianie : T.akcje.uruchom,
    });
    przycisk.addEventListener('click', () => void uruchom(plan));
    return przycisk;
  }

  function formularzZlecenia(): HTMLElement {
    const zajete = dzialanie.rodzaj === 'rozkladanie';
    const pole = el('textarea', {
      klasa: 'dn-pole-kontrolka',
      rows: 2,
      placeholder: T.nowe.zastepcza,
      'aria-label': T.nowe.etykieta,
      disabled: zajete || dokument === null,
    }) as HTMLTextAreaElement;
    pole.value = zlecenie;

    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--atrament dn-btn--sm',
      type: 'submit',
      disabled: zajete || dokument === null || zlecenie.trim() === '',
      tekst: zajete ? T.nowe.rozkladanie : T.nowe.rozloz,
    });
    /* Wpisywanie nie przerysowuje panelu — przerysowanie zabrałoby polu
       ognisko po każdym znaku — więc guzik przestawia się wprost. */
    pole.addEventListener('input', () => {
      zlecenie = pole.value;
      przycisk.toggleAttribute('disabled', zlecenie.trim() === '');
    });

    const formularz = el('form', { klasa: 'dn-pole st-odsun-sekcja' }, [
      el('span', { klasa: 'dn-pole-etykieta', tekst: T.nowe.etykieta }),
      pole,
      dokument === null ? el('span', { klasa: 'dn-meta', tekst: T.nowe.brakDokumentu }) : null,
      el('div', { klasa: 'dn-pas-dzialan' }, [przycisk]),
    ]);
    formularz.addEventListener('submit', (zdarzenie) => {
      zdarzenie.preventDefault();
      void rozloz();
    });
    return formularz;
  }

  function sekcjaLancuchow(): Dziecko[] {
    const naglowek = el('div', { klasa: 'pt-etykieta st-odsun-sekcja', tekst: T.lancuchy.etykieta });
    if (odmowaLancuchow !== null) {
      return [naglowek, alert(T.odmowa.wykazLancuchow, opisOdmowy(odmowaLancuchow, 'chain.list'))];
    }
    if (lancuchy.length === 0) {
      return [naglowek, el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: T.lancuchy.brak })];
    }
    const zajete = dzialanie.rodzaj === 'lancuch';
    const wiersze = lancuchy.map((lancuch) => {
      const przycisk = el('button', {
        klasa: 'dn-btn dn-btn--duch dn-btn--sm',
        type: 'button',
        disabled: zajete || dokument === null,
        tekst: zajete ? T.lancuchy.uruchamianie : T.lancuchy.uruchom,
      });
      przycisk.addEventListener('click', () => void uruchomLancuch(lancuch));
      return el('div', { klasa: 'st-panel-wiersz' }, [
        el('span', { klasa: 'dn-kropka dn-kropka--neutralna', 'aria-hidden': 'true' }),
        lancuch.name,
        el('span', { klasa: 'dn-meta', tekst: zLiczba(lancuch.steps.length, T.lancuchy.jednostkaKroki) }),
        przycisk,
      ]);
    });
    const postep = postepLancucha === null ? null : el('div', { klasa: 'dn-nota', tekst: postepLancucha });
    return [naglowek, ...wiersze, postep];
  }

  function widokRozkladu(plan: StudioTaskPlan): Dziecko[] {
    const zadania = plan.tasks ?? [];
    const wskazane = zadania.find((zadanie) => zadanie.id === wskazaneZadanie);
    const wiersze: Dziecko[] = [];
    for (const zadanie of zadania) {
      wiersze.push(wierszZadania(zadanie));
      if (wskazane !== undefined && wskazane.id === zadanie.id) wiersze.push(...czynnosciZadania(zadanie));
    }
    return [
      wyborRozkladu(plan),
      el('div', { klasa: 'pt-etykieta', tekst: `${T.zlecenie} ${plan.order}` }),
      ...(zadania.length === 0 ? [pustyStan(T.pusty.tytul, T.pusty.opis)] : wiersze),
      el('div', { klasa: 'st-panel-wiersz' }, [plakietkaStanu(plan), przyciskPetli(plan)]),
    ];
  }

  function zawartosc(): Dziecko[] {
    switch (stan.rodzaj) {
      case 'ladowanie':
        return [wierszPulsu(T.ladowanie)];
      case 'brakZaleznosci':
        return [alert(T.odmowa[stan.brak], opisOdmowy(undefined, `plan.${stan.brak}`))];
      case 'odmowaRozkladu':
        return [
          alert(T.odmowa.rozklad, opisOdmowy(stan.blad, 'plan.get')),
          formularzZlecenia(),
          ...sekcjaLancuchow(),
        ];
      case 'brakRozkladu':
        return [
          pustyStan(T.brakRozkladu.tytul, T.brakRozkladu.opis),
          formularzZlecenia(),
          ...sekcjaLancuchow(),
        ];
      case 'rozklad':
        return [...widokRozkladu(stan.plan), formularzZlecenia(), ...sekcjaLancuchow()];
    }
  }

  function opisZnacznika(): string {
    if (stan.rodzaj !== 'rozklad') return '';
    const zadania = stan.plan.tasks ?? [];
    if (zadania.length === 0) return '';
    const domkniete = zadania.filter((zadanie) => zadanie.state === StudioTaskState.Done).length;
    return `${domkniete} / ${zadania.length}`;
  }

  function przerysuj(): void {
    const miara = opisZnacznika();
    znacznik.textContent = miara;
    if (miara === '') znacznik.removeAttribute('aria-label');
    else znacznik.setAttribute('aria-label', `${T.panel.etykietaZnacznika}: ${miara}`);

    const dzieci = zawartosc();
    if (dzialanie.rodzaj === 'odmowa') dzieci.push(alert(dzialanie.naglowek, dzialanie.opis));
    const widoczne = dzieci.filter(
      (dziecko): dziecko is Node | string => dziecko !== null && dziecko !== undefined && dziecko !== false,
    );
    tresc.replaceChildren(...widoczne);
  }

  function odswiez(nowy: Stan): void {
    stan = nowy;
    przerysuj();
  }

  /* ── Nasłuchy i start ─────────────────────────────────────────────────── */

  const odsubskrybujDokument = zaleznosci.kanal.naZdarzenie(EventType.StudioDocumentChanged, (zdarzenie) => {
    if (zdjety || zaleznosci.idOkna === null || zdarzenie.document.windowId !== zaleznosci.idOkna) return;
    const poprzedni = dokument?.id ?? null;
    if (zdarzenie.change === ChangeKind.Deleted) {
      if (dokument !== null && dokument.id === zdarzenie.document.id) dokument = null;
    } else {
      dokument = zdarzenie.document;
    }
    /* Rozkład należy do dokumentu, więc zmiana dokumentu okna każe wczytać go
       na nowo; sam zapis treści tego samego dokumentu rozkładu nie rusza. */
    if ((dokument?.id ?? null) !== poprzedni) {
      wskazaneZadanie = null;
      postepLancucha = null;
      void wczytajRozklad();
    } else {
      przerysuj();
    }
  });

  const odsubskrybujLancuch = zaleznosci.kanal.naZdarzenie(EventType.StudioChainProgressed, (zdarzenie) => {
    if (zdjety || biegLancucha === null || zdarzenie.runId !== biegLancucha) return;
    postepLancucha = opisPostepuLancucha(zdarzenie.stepIndex, zdarzenie.stepCount);
    /* Każdy krok łańcucha odkłada własną pracę w rozkładzie — wykaz zadań bez
       ponownego wczytania pokazywałby stan sprzed kroku. */
    void wczytajRozklad();
  });

  przerysuj();
  void wczytajRozklad();
  void wczytajLancuchy();

  return {
    zdejmij() {
      zdjety = true;
      odsubskrybujDokument();
      odsubskrybujLancuch();
      wezel.replaceChildren();
    },
  };
};
