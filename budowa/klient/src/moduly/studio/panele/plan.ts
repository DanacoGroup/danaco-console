/**
 * Okno robocze Studia — panel Plan. Rozkłada zlecenie dokumentowe Operatora na
 * zadania i pokazuje stan każdego z nich wraz ze sterowaniem pętli wykonawczej:
 * `studio.plan.get` niesie bieżący rozkład, `studio.plan.run` i `studio.plan.stop`
 * puszczają go i zatrzymują. Kształt wzięty z `design/05-okna/moduly/studio.html`,
 * `section#panel-plan`.
 *
 * Rozkład bez zadań i odmowa rdzenia są dwoma różnymi stanami: pierwszy to
 * odpowiedź rdzenia bez żadnego zadania, drugi to brak odpowiedzi w ogóle.
 * Żadna wartość widoczna dla Operatora nie jest tu wpisana z góry — poza
 * odmową źródło każdej jest odpowiedź rdzenia.
 */

import {
  Command,
  StudioPlanState,
  StudioTaskState,
  type StudioDocumentTask,
  type StudioTaskPlan,
} from '../../../../../shared/contract.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import type { MontazPanelu, ZaleznosciPanelu, ZamontowanyPanel } from './umowa.ts';
import { el, tekst } from '../narzedzia.ts';
import { planTresci } from './plan-tresci.ts';

type Zaleznosc = 'okno' | 'sesja';

type Dzialanie =
  | { rodzaj: 'spoczynek' }
  | { rodzaj: 'uruchamianie' }
  | { rodzaj: 'zatrzymywanie' }
  | { rodzaj: 'odmowa'; naglowek: string; opis: string };

type Stan =
  | { rodzaj: 'ladowanie' }
  | { rodzaj: 'brakZaleznosci'; brak: Zaleznosc }
  | { rodzaj: 'odmowaRdzenia'; naglowek: string; opis: string }
  | { rodzaj: 'plan'; plan: StudioTaskPlan; dzialanie: Dzialanie };

const WARIANT_ZNACZNIKA: Partial<Record<StudioTaskState, string>> = {
  [StudioTaskState.Done]: 'dn-kropka--sukces',
  [StudioTaskState.Failed]: 'dn-kropka--blad',
  [StudioTaskState.Blocked]: 'dn-kropka--ostrzezenie',
  [StudioTaskState.Skipped]: 'dn-kropka--neutralna',
  [StudioTaskState.Pending]: 'dn-kropka--neutralna',
};

const WARIANT_PLAKIETKI: Partial<Record<StudioPlanState, string>> = {
  [StudioPlanState.Running]: 'dn-plakietka--sygnal',
  [StudioPlanState.Paused]: 'dn-plakietka--ostrzezenie',
  [StudioPlanState.Done]: 'dn-plakietka--sukces',
};

export const montujPanelPlanu: MontazPanelu = (wezel, zaleznosci: ZaleznosciPanelu): ZamontowanyPanel => {
  let zdjety = false;

  const tresc = el('div', { klasa: 'sta-okno-tresc st-panel-lista' });
  wezel.replaceChildren(
    el('header', { klasa: 'sta-okno-belka' }, [
      el('span', { klasa: 'sta-okno-tytul' }, [el('b', { tekst: tekst('karty.plan') })]),
    ]),
    tresc,
  );

  odswiez({ rodzaj: 'ladowanie' });
  void wczytaj();

  async function wczytaj(): Promise<void> {
    if (zaleznosci.idSesji === null) return odswiez({ rodzaj: 'brakZaleznosci', brak: 'sesja' });
    if (zaleznosci.idOkna === null) return odswiez({ rodzaj: 'brakZaleznosci', brak: 'okno' });

    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioPlanGet, {});
    if (zdjety) return;
    if (!wynik.udany) {
      return odswiez({
        rodzaj: 'odmowaRdzenia',
        naglowek: planTresci.odmowa.rozklad,
        opis: wynik.blad?.message ?? tekst('odmowa.brakOpisu'),
      });
    }
    if (wynik.wynik === undefined) {
      return odswiez({
        rodzaj: 'odmowaRdzenia',
        naglowek: planTresci.odmowa.brakRozkladu,
        opis: tekst('odmowa.brakOpisu'),
      });
    }
    odswiez({ rodzaj: 'plan', plan: wynik.wynik.plan, dzialanie: { rodzaj: 'spoczynek' } });
  }

  async function uruchom(plan: StudioTaskPlan): Promise<void> {
    odswiez({ rodzaj: 'plan', plan, dzialanie: { rodzaj: 'uruchamianie' } });
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioPlanRun, { planId: plan.id });
    if (zdjety) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      return odswiez({
        rodzaj: 'plan',
        plan,
        dzialanie: {
          rodzaj: 'odmowa',
          naglowek: planTresci.odmowa.uruchomienie,
          opis: wynik.blad?.message ?? tekst('odmowa.brakOpisu'),
        },
      });
    }
    if (!wynik.wynik.started) {
      return odswiez({
        rodzaj: 'plan',
        plan: wynik.wynik.plan,
        dzialanie: {
          rodzaj: 'odmowa',
          naglowek: planTresci.odmowa.uruchomienie,
          opis: wynik.wynik.refusalReason ?? tekst('odmowa.brakOpisu'),
        },
      });
    }
    odswiez({ rodzaj: 'plan', plan: wynik.wynik.plan, dzialanie: { rodzaj: 'spoczynek' } });
  }

  async function zatrzymaj(plan: StudioTaskPlan): Promise<void> {
    odswiez({ rodzaj: 'plan', plan, dzialanie: { rodzaj: 'zatrzymywanie' } });
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioPlanStop, { planId: plan.id });
    if (zdjety) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      return odswiez({
        rodzaj: 'plan',
        plan,
        dzialanie: {
          rodzaj: 'odmowa',
          naglowek: planTresci.odmowa.zatrzymanie,
          opis: wynik.blad?.message ?? tekst('odmowa.brakOpisu'),
        },
      });
    }
    odswiez({ rodzaj: 'plan', plan: wynik.wynik.plan, dzialanie: { rodzaj: 'spoczynek' } });
  }

  function alertOdmowy(naglowek: string, opis: string): HTMLElement {
    return el('div', { klasa: 'dn-alert dn-alert--wstega dn-alert--blad', role: 'alert' }, [
      el('span', { klasa: 'dn-alert-tresc' }, [el('b', { tekst: naglowek }), el('span', { tekst: opis })]),
    ]);
  }

  function widokLadowania(): HTMLElement[] {
    return [
      el('div', { klasa: 'dn-wykaz-modulu-poz' }, [
        el('span', { klasa: 'dn-kropka dn-kropka--sygnal dn-kropka--tetno', 'aria-hidden': 'true' }),
        planTresci.ladowanie,
      ]),
    ];
  }

  function znacznikZadania(stanZadania: StudioTaskState): HTMLElement {
    if (stanZadania === StudioTaskState.Running) {
      return el('span', { klasa: 'pt-tetno', 'aria-hidden': 'true' });
    }
    const wariant = WARIANT_ZNACZNIKA[stanZadania] ?? 'dn-kropka--neutralna';
    return el('span', { klasa: `dn-kropka ${wariant}`, 'aria-hidden': 'true' });
  }

  function wierszZadania(zadanie: StudioDocumentTask): HTMLElement {
    return el('div', { klasa: 'st-panel-wiersz' }, [
      znacznikZadania(zadanie.state),
      zadanie.title,
      el('span', { klasa: 'dn-meta', tekst: planTresci.stanZadania[zadanie.state].glif }),
    ]);
  }

  function plakietkaStanu(stanRozkladu: StudioPlanState): HTMLElement {
    const wariant = WARIANT_PLAKIETKI[stanRozkladu];
    return el('span', { klasa: wariant ? `dn-plakietka ${wariant}` : 'dn-plakietka' }, [
      stanRozkladu === StudioPlanState.Running
        ? el('span', { klasa: 'pt-tetno', 'aria-hidden': 'true' })
        : null,
      `${planTresci.stanPetli} ${planTresci.stanRozkladuSlowem[stanRozkladu]}`,
    ]);
  }

  function przyciskPetli(plan: StudioTaskPlan, dzialanie: Dzialanie): HTMLElement | null {
    const wBiegu = dzialanie.rodzaj === 'uruchamianie' || dzialanie.rodzaj === 'zatrzymywanie';
    if (plan.state === StudioPlanState.Running) {
      const przycisk = el('button', {
        klasa: 'dn-btn dn-btn--niebezpieczny',
        type: 'button',
        tekst: dzialanie.rodzaj === 'zatrzymywanie' ? planTresci.akcje.zatrzymywanie : planTresci.akcje.zatrzymaj,
        disabled: wBiegu,
      });
      przycisk.addEventListener('click', () => void zatrzymaj(plan));
      return przycisk;
    }
    if (plan.state === StudioPlanState.Done) return null;
    const przycisk = el('button', {
      klasa: 'dn-btn dn-btn--atrament',
      type: 'button',
      tekst: dzialanie.rodzaj === 'uruchamianie' ? planTresci.akcje.uruchamianie : planTresci.akcje.uruchom,
      disabled: wBiegu,
    });
    przycisk.addEventListener('click', () => void uruchom(plan));
    return przycisk;
  }

  function widokPlanu(plan: StudioTaskPlan, dzialanie: Dzialanie): HTMLElement[] {
    const zadania = plan.tasks ?? [];
    const wiersze: HTMLElement[] =
      zadania.length === 0
        ? [
            el('div', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty' }, [
              el('span', { klasa: 'dn-pusty-stan-tytul', tekst: planTresci.pusty.tytul }),
              el('span', { klasa: 'dn-pusty-stan-opis', tekst: planTresci.pusty.opis }),
            ]),
          ]
        : zadania.map((zadanie) => wierszZadania(zadanie));

    const widok: HTMLElement[] = [
      el('div', { klasa: 'pt-etykieta', tekst: `${planTresci.zlecenie} ${plan.order}` }),
      ...wiersze,
      el('div', { klasa: 'st-panel-wiersz' }, [plakietkaStanu(plan.state), przyciskPetli(plan, dzialanie)]),
    ];
    if (dzialanie.rodzaj === 'odmowa') widok.push(alertOdmowy(dzialanie.naglowek, dzialanie.opis));
    return widok;
  }

  function zawartosc(stan: Stan): HTMLElement[] {
    switch (stan.rodzaj) {
      case 'ladowanie':
        return widokLadowania();
      case 'brakZaleznosci':
        return [alertOdmowy(planTresci.odmowa[stan.brak], tekst('odmowa.brakOpisu'))];
      case 'odmowaRdzenia':
        return [alertOdmowy(stan.naglowek, stan.opis)];
      case 'plan':
        return widokPlanu(stan.plan, stan.dzialanie);
    }
  }

  function odswiez(stan: Stan): void {
    tresc.replaceChildren(...zawartosc(stan));
  }

  return {
    zdejmij() {
      zdjety = true;
      wezel.replaceChildren();
    },
  };
};
