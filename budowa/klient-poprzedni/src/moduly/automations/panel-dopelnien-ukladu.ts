import {
  AutomationDependencyKind,
  OrchestrationGateRule,
} from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import type { ZrodloAutomations } from './zrodlo-automations';

/**
 * Panel dopełnień układu zależności — cztery rzeczy, których krawędź nie
 * wyraża.
 *
 * Krawędź mówi, że kroki się schodzą, i to wszystko, co mówi. Zostają cztery
 * pytania, na które trzeba odpowiedzieć osobno:
 *
 *   - KIEDY tory scalają się w kroku wspólnym — bramka dołączenia;
 *   - CZY zbiór kroków biegnie razem, czy jeden po drugim — grupa;
 *   - CO wycofuje skutki kroku, gdy przebieg pękł w pół — kompensacja;
 *   - CZYIM silnikiem jadą kolejki tej automatyki — spięcie z MultitaskingAI.
 *
 * Panel stoi obok panelu zależności, w narzędziach kontekstowych okna
 * Orchestrator: gospodaruje tym samym układem, tylko innym jego wymiarem.
 *
 * Ocena bramki nie blokuje zapisu — tak samo jak przy krawędzi. Bramka na
 * kroku, którego jeszcze nie ma, zapisuje się, a zastrzeżenie wraca
 * w odpowiedzi: Workflow Builder buduje układ krok po kroku i odmowa kazałaby
 * Operatorowi układać go w jedynej dopuszczonej kolejności.
 */
export interface PanelDopelnienUkladu {
  element: HTMLElement;
}

export function utworzPanelDopelnienUkladu(
  zrodlo: ZrodloAutomations,
  automatyka: () => string,
): PanelDopelnienUkladu {
  const odpowiedz = utworzWierszOdpowiedzi();

  // ── Bramka dołączenia ──────────────────────────────────────────────────────
  const krokBramki = poleTekstowe({
    etykieta: 'Krok scalający tory',
    podpowiedz: 'identyfikator kroku',
  });
  const regula = poleWyboru({ etykieta: 'Reguła scalenia' }, [
    { wartosc: OrchestrationGateRule.All, etykieta: 'wszystkie tory' },
    { wartosc: OrchestrationGateRule.Any, etykieta: 'dowolny tor' },
    { wartosc: OrchestrationGateRule.Count, etykieta: 'wskazana liczba torów' },
  ]);
  const licznikTorow = document.createElement('input');
  licznikTorow.type = 'number';
  licznikTorow.className = 'dn-pole-kontrolka';
  licznikTorow.min = '1';
  licznikTorow.value = '2';
  const zapiszBramke = przycisk('Zapisz bramkę', 'dn-btn dn-btn--sm dn-btn--zarys');

  // ── Grupa kroków ───────────────────────────────────────────────────────────
  const nazwaGrupy = poleTekstowe({ etykieta: 'Nazwa grupy kroków' });
  const krokiGrupy = poleTekstowe({
    etykieta: 'Kroki grupy',
    podpowiedz: 'identyfikatory rozdzielone przecinkiem',
    opis: 'Wykaz PUSTY usuwa grupę — grupa bez ani jednego kroku nie mówi o niczym.',
  });
  const rodzajGrupy = poleWyboru({ etykieta: 'Rodzaj grupowania' }, [
    { wartosc: AutomationDependencyKind.Parallel, etykieta: 'równolegle' },
    { wartosc: AutomationDependencyKind.Sequential, etykieta: 'w ścisłej kolejności' },
  ]);
  const zapiszGrupe = przycisk('Zapisz grupę', 'dn-btn dn-btn--sm dn-btn--zarys');

  // ── Kompensacja ────────────────────────────────────────────────────────────
  const krokGlowny = poleTekstowe({ etykieta: 'Krok główny' });
  const krokWycofu = poleTekstowe({
    etykieta: 'Krok wycofujący',
    podpowiedz: 'puste zdejmuje kompensację',
  });
  const zapiszKompensacje = przycisk('Zapisz kompensację', 'dn-btn dn-btn--sm dn-btn--zarys');

  // ── Spięcie z MultitaskingAI ───────────────────────────────────────────────
  const rolaSpiecia = poleTekstowe({
    etykieta: 'Okno roli środowiska MultitaskingAI',
    podpowiedz: 'puste spina bez wskazania roli',
  });
  const spnij = przycisk('Spnij kolejki', 'dn-btn dn-btn--sm dn-btn--zarys');
  const rozlacz = przycisk('Rozłącz kolejki', 'dn-btn dn-btn--sm dn-btn--zarys');

  const tytul = document.createElement('h4');
  tytul.className = 'da-panel__tytul';
  tytul.textContent = 'Dopełnienia układu — bramki, grupy, kompensacje, spięcie';

  const element = document.createElement('section');
  element.className = 'aa-panel aa-panel--dopelnienia';
  element.dataset['panel'] = 'dopelnienia-ukladu';
  element.append(
    tytul,
    krokBramki.element,
    regula.element,
    licznikTorow,
    zapiszBramke,
    nazwaGrupy.element,
    krokiGrupy.element,
    rodzajGrupy.element,
    zapiszGrupe,
    krokGlowny.element,
    krokWycofu.element,
    zapiszKompensacje,
    rolaSpiecia.element,
    spnij,
    rozlacz,
    odpowiedz.element,
  );

  /** Układ wskazany w oknie; pusty znaczy „nie ma czego dopełniać”. */
  function uklad(): string | null {
    const kod = automatyka();
    if (kod === '') {
      odpowiedz.pokaz('Wybierz automatykę — dopełnienia dotyczą jednego układu.', false);
      return null;
    }
    return kod;
  }

  async function bramka(): Promise<void> {
    const kod = uklad();
    if (kod === null) return;
    const krok = krokBramki.kontrolka.value.trim();
    if (krok === '') {
      odpowiedz.pokaz('Wskaż krok scalający — bramka bez kroku nie ma czego scalać.', false);
      return;
    }
    const wybrana = regula.kontrolka.value as OrchestrationGateRule;
    const zadanie: Parameters<ZrodloAutomations['ustawBramke']>[0] = {
      workflowId: kod,
      stepId: krok,
      rule: wybrana,
    };
    if (wybrana === OrchestrationGateRule.Count) {
      zadanie.count = Number.parseInt(licznikTorow.value, 10) || 0;
    }
    odpowiedz.pokaz(`Zapis bramki kroku ${krok}…`, true);
    const wynik = await zrodlo.ustawBramke(zadanie);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Zapis bramki', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    // Zastrzeżenie nie jest odmową: bramka jest zapisana, a układ ma wadę,
    // którą Operator dopiero domknie kolejnym krokiem.
    odpowiedz.pokaz(
      wynik.wynik.valid
        ? `Bramka zapisana — układ ma teraz ${wynik.wynik.gates.length} bramek.`
        : `Bramka ZAPISANA, ale układ ma zastrzeżenia: ${(wynik.wynik.issues ?? []).join('; ')}`,
      wynik.wynik.valid,
    );
  }

  async function grupa(): Promise<void> {
    const kod = uklad();
    if (kod === null) return;
    const nazwa = nazwaGrupy.kontrolka.value.trim();
    if (nazwa === '') {
      odpowiedz.pokaz('Grupa wymaga nazwy — bez niej nie da się jej wskazać w wykazie.', false);
      return;
    }
    const kroki = krokiGrupy.kontrolka.value
      .split(',')
      .map((krok) => krok.trim())
      .filter((krok) => krok !== '');
    odpowiedz.pokaz(`Zapis grupy „${nazwa}"…`, true);
    const wynik = await zrodlo.ustawGrupe({
      workflowId: kod,
      name: nazwa,
      stepIds: kroki,
      kind: rodzajGrupy.kontrolka.value as AutomationDependencyKind,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowy('Zapis grupy', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    odpowiedz.pokaz(
      kroki.length === 0
        ? `Grupa usunięta — układ ma teraz ${wynik.wynik.groups.length} grup.`
        : `Grupa „${nazwa}" obejmuje ${kroki.length} kroków; układ ma ${wynik.wynik.groups.length} grup.`,
      true,
    );
  }

  async function kompensacja(): Promise<void> {
    const kod = uklad();
    if (kod === null) return;
    const glowny = krokGlowny.kontrolka.value.trim();
    if (glowny === '') {
      odpowiedz.pokaz('Wskaż krok główny — kompensacja wycofuje skutki czegoś.', false);
      return;
    }
    const wycofujacy = krokWycofu.kontrolka.value.trim();
    odpowiedz.pokaz(`Zapis kompensacji kroku ${glowny}…`, true);
    const wynik = await zrodlo.ustawKompensacje({
      workflowId: kod,
      stepId: glowny,
      compensationStepId: wycofujacy,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Zapis kompensacji', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    odpowiedz.pokaz(
      wycofujacy === ''
        ? `Kompensacja kroku ${glowny} zdjęta — układ ma ${wynik.wynik.compensations.length} kompensacji.`
        : `Krok ${wycofujacy} wycofuje skutki kroku ${glowny}.`,
      true,
    );
  }

  async function spiecie(spiete: boolean): Promise<void> {
    const kod = uklad();
    if (kod === null) return;
    const rola = rolaSpiecia.kontrolka.value.trim();
    const zadanie: Parameters<ZrodloAutomations['spnijZMultitaskingiem']>[0] = {
      workflowId: kod,
      linked: spiete,
    };
    if (spiete && rola !== '') zadanie.roleId = rola;
    odpowiedz.pokaz(spiete ? 'Spinanie kolejek…' : 'Rozłączanie kolejek…', true);
    const wynik = await zrodlo.spnijZMultitaskingiem(zadanie);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Spięcie z MultitaskingAI', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    const kolejki = wynik.wynik.queueIds ?? [];
    odpowiedz.pokaz(
      wynik.wynik.linked
        ? `Kolejki automatyki prowadzi silnik środowiska MultitaskingAI — objętych ${kolejki.length}.`
        : `Kolejki rozłączone — wraca stan wyjściowy. Objętych ${kolejki.length}.`,
      true,
    );
  }

  zapiszBramke.addEventListener('click', () => void bramka());
  zapiszGrupe.addEventListener('click', () => void grupa());
  zapiszKompensacje.addEventListener('click', () => void kompensacja());
  spnij.addEventListener('click', () => void spiecie(true));
  rozlacz.addEventListener('click', () => void spiecie(false));

  return { element };
}
