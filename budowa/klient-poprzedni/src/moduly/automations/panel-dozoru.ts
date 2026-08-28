/**
 * Panel Execution Monitora obsługuje trzynaście czynności szuflad i paneli popover okna
 * przebiegów, wskazywanych zawsze przez Operatora, nie przez stan modułu.
 */
import { AutomationAlertTrigger, AutomationLogLevel } from '../../../../shared/contract';
import { pole, poleLiczbowe, poleTajne, wybor } from '../../modele/kontrolki-formularza';
import {
  odczytajChwile,
  odczytajLiczbe,
  odczytajWykaz,
  utworzPanelDobudowy,
  wykonajCzynnoscPanelu,
} from './panel-dobudowy';
import type { StanAutomatyki } from './stan-automatyki';
import type { ZrodloAutomations } from './zrodlo-automations';

/** Panel dopełniający Execution Monitora, osadzany wewnątrz okna przebiegów wraz z jego szufladami i panelami popover. */
export interface PanelDozoru {
  element: HTMLElement;
}

/** Poziomy logu w wykazie zawężenia wprost z wyliczenia kontraktu; pusta wartość znaczy wszystkie poziomy naraz. */
const POZIOMY: ReadonlyArray<readonly [string, string]> = [
  ['', 'wszystkie poziomy'],
  ...Object.values(AutomationLogLevel).map((poziom): readonly [string, string] => [poziom, poziom]),
];

/** Wyzwalacze alarmu wprost z wyliczenia kontraktu — mapa zupełna, złożona automatycznie z wartości wyliczenia. */
const WYZWALACZE_ALARMU: ReadonlyArray<readonly [string, string]> = Object.values(
  AutomationAlertTrigger,
).map((wyzwalacz): readonly [string, string] => [wyzwalacz, wyzwalacz]);

export function utworzPanelDozoru(
  zrodlo: ZrodloAutomations,
  stan: StanAutomatyki,
): PanelDozoru {
  const panel = utworzPanelDobudowy(
    'Log, ładunki, alarmy i skarbiec',
    'Pełny zapis zdarzeń przebiegu, stan każdego kroku, punkty wznowienia, ' +
      'podgląd i odtworzenie ładunku, reguły alarmowania, budżety czasu, ' +
      'skarbiec poświadczeń i dziennik audytu.',
  );

  const idPrzebiegu = panel.dodajPole('Przebieg', pole('Przebieg', 'przebieg-…'));
  const idKroku = panel.dodajPole('Krok', pole('Krok', 'krok-1'));
  const poziom = panel.dodajPole('Poziom logu', wybor('Poziom logu', POZIOMY));
  const granica = panel.dodajPole('Granica wykazu', poleLiczbowe('Granica wykazu', '200'));
  const idPunktu = panel.dodajPole(
    'Punkt wznowienia',
    pole('Punkt wznowienia', 'pusty bierze najnowszy'),
  );
  const wyzwalacz = panel.dodajPole('Wyzwalacz alarmu', wybor('Wyzwalacz alarmu', WYZWALACZE_ALARMU));
  const warunekAlarmu = panel.dodajPole('Warunek alarmu', pole('Warunek alarmu', ''));
  const kanaly = panel.dodajPole('Kanały powiadomień', pole('Kanały', 'mobile, poczta'));
  const idReguly = panel.dodajPole('Reguła', pole('Reguła', 'pusta zakłada nową'));
  const budzetPrzebiegu = panel.dodajPole('Budżet przebiegu (s)', poleLiczbowe('Budżet przebiegu', '0'));
  const budzetKroku = panel.dodajPole('Budżet kroku (s)', poleLiczbowe('Budżet kroku', '0'));
  const nazwaSekretu = panel.dodajPole('Nazwa poświadczenia', pole('Nazwa poświadczenia', ''));
  const wartoscSekretu = panel.dodajPole(
    'Wartość poświadczenia',
    poleTajne({ etykieta: 'Wartość poświadczenia' }).kontrolka,
    'Idzie do skarbca poza bazą. Rdzeń nie oddaje jej żadną komendą — także tą.',
  );
  const odwolanieSekretu = panel.dodajPole('Referencja poświadczenia', pole('Referencja', 'sejf:…'));
  const odChwili = panel.dodajPole('Audyt od', pole('Audyt od', '2026-08-01T00:00:00Z'));
  const doChwili = panel.dodajPole('Audyt do', pole('Audyt do', '2026-08-17T00:00:00Z'));

  /** Przebieg wskazany w polu; brak wstrzymuje czynność z nazwanym powodem. */
  function przebieg(): string | null {
    if (idPrzebiegu.value.trim() === '') {
      panel.tresc.potwierdzenie(
        'Wskaż przebieg — identyfikator stoi w wierszu historii Execution Monitora.',
        false,
      );
      return null;
    }
    return idPrzebiegu.value.trim();
  }

  /** Automatyka bieżąca modułu; pusta wstrzymuje czynność. */
  function automatyka(): string | null {
    if (stan.automatyka() === '') {
      panel.tresc.potwierdzenie('Wskaż automatykę — reguła i budżet należą do niej.', false);
      return null;
    }
    return stan.automatyka();
  }

  panel.dodajCzynnosc('Log przebiegu', () => {
    const kod = przebieg();
    if (kod === null) return;
    const zadanie: Parameters<ZrodloAutomations['dziennikPrzebiegu']>[0] = { executionId: kod };
    if (poziom.value !== '') zadanie.level = poziom.value as AutomationLogLevel;
    if (idKroku.value.trim() !== '') zadanie.stepId = idKroku.value.trim();
    const limit = odczytajLiczbe(granica.value);
    if (limit !== undefined) zadanie.limit = limit;
    wykonajCzynnoscPanelu(panel, 'Odczyt logu przebiegu…', zrodlo.dziennikPrzebiegu(zadanie),
      'Rdzeń nie oddał logu przebiegu.',
      'Rdzeń oddał zapis zdarzeń przebiegu; pole „truncated” mówi, czy wykaz przycięto granicą.');
  });

  panel.dodajCzynnosc('Kroki przebiegu', () => {
    const kod = przebieg();
    if (kod === null) return;
    wykonajCzynnoscPanelu(panel, 'Odczyt kroków przebiegu…',
      zrodlo.krokiPrzebiegu({ executionId: kod }),
      'Rdzeń nie oddał kroków przebiegu.', 'Rdzeń oddał stan każdego kroku przebiegu osobno.');
  });

  panel.dodajCzynnosc('Punkty wznowienia', () => {
    const kod = przebieg();
    if (kod === null) return;
    wykonajCzynnoscPanelu(panel, 'Odczyt punktów wznowienia…',
      zrodlo.punktyWznowienia({ executionId: kod }),
      'Rdzeń nie oddał punktów wznowienia.', 'Rdzeń oddał punkty wznowienia przebiegu.');
  });

  panel.dodajCzynnosc('Wznów przebieg', () => {
    const kod = przebieg();
    if (kod === null) return;
    const zadanie: Parameters<ZrodloAutomations['wznowPrzebieg']>[0] = { executionId: kod };
    if (idPunktu.value.trim() !== '') zadanie.checkpointId = idPunktu.value.trim();
    wykonajCzynnoscPanelu(panel, 'Wznawianie przebiegu…', zrodlo.wznowPrzebieg(zadanie),
      'Rdzeń nie wznowił przebiegu.',
      'Rdzeń wznowił przebieg od punktu wznowienia, bez powtarzania kroków ukończonych.');
  });

  panel.dodajCzynnosc('Podgląd ładunku', () => {
    const kod = przebieg();
    if (kod === null) return;
    if (idKroku.value.trim() === '') {
      panel.tresc.potwierdzenie('Wskaż krok — ładunek należy do kroku, nie do całego przebiegu.', false);
      return;
    }
    wykonajCzynnoscPanelu(panel, 'Odczyt ładunku kroku…',
      zrodlo.podgladLadunku({ executionId: kod, stepId: idKroku.value.trim() }),
      'Rdzeń nie oddał ładunku kroku.',
      'Rdzeń oddał dane wejściowe i wyjściowe kroku; pola zamaskowane regułą redakcji ' +
        'sekretów nazywa pole „redactedFields”.');
  });

  panel.dodajCzynnosc('Odtwórz przebieg', () => {
    const kod = przebieg();
    if (kod === null) return;
    const zadanie: Parameters<ZrodloAutomations['odtworzPrzebieg']>[0] = { executionId: kod };
    if (idKroku.value.trim() !== '') zadanie.fromStepId = idKroku.value.trim();
    wykonajCzynnoscPanelu(panel, 'Odtwarzanie przebiegu…', zrodlo.odtworzPrzebieg(zadanie),
      'Rdzeń nie odtworzył przebiegu.',
      'Rdzeń założył przebieg odtworzony; przebieg źródłowy został nietknięty.');
  });

  panel.dodajCzynnosc('Zapisz regułę alarmowania', () => {
    const kod = automatyka();
    if (kod === null) return;
    const zadanie: Parameters<ZrodloAutomations['ustawReguleAlarmowania']>[0] = {
      workflowId: kod,
      trigger: wyzwalacz.value as AutomationAlertTrigger,
      channels: odczytajWykaz(kanaly.value),
    };
    if (idReguly.value.trim() !== '') zadanie.ruleId = idReguly.value.trim();
    if (warunekAlarmu.value.trim() !== '') zadanie.condition = warunekAlarmu.value.trim();
    wykonajCzynnoscPanelu(panel, 'Zapis reguły alarmowania…',
      zrodlo.ustawReguleAlarmowania(zadanie),
      'Rdzeń nie zapisał reguły alarmowania.',
      'Rdzeń zapisał warunek i kanały powiadomień reguły.');
  });

  panel.dodajCzynnosc('Wykaz reguł alarmowania', () => {
    const zadanie: Parameters<ZrodloAutomations['regulyAlarmowania']>[0] = {};
    if (stan.automatyka() !== '') zadanie.workflowId = stan.automatyka();
    wykonajCzynnoscPanelu(panel, 'Odczyt reguł…', zrodlo.regulyAlarmowania(zadanie),
      'Rdzeń nie oddał reguł alarmowania.',
      'Rdzeń oddał reguły alarmowania; bez wskazania automatyki — komplet reguł Operatora.');
  });

  panel.dodajCzynnosc('Zapisz budżety czasu', () => {
    const kod = automatyka();
    if (kod === null) return;
    const zadanie: Parameters<ZrodloAutomations['ustawBudzetyPrzebiegu']>[0] = { workflowId: kod };
    const przebiegSekundy = odczytajLiczbe(budzetPrzebiegu.value);
    if (przebiegSekundy !== undefined) zadanie.runBudgetSeconds = przebiegSekundy;
    const krokSekundy = odczytajLiczbe(budzetKroku.value);
    if (krokSekundy !== undefined) zadanie.stepBudgetSeconds = krokSekundy;
    if (idReguly.value.trim() !== '') zadanie.alertRuleId = idReguly.value.trim();
    wykonajCzynnoscPanelu(panel, 'Zapis budżetów czasu…',
      zrodlo.ustawBudzetyPrzebiegu(zadanie),
      'Rdzeń nie zapisał budżetów czasu.',
      'Rdzeń zapisał budżet przebiegu i kroku; zero znaczy „bez granicy”.');
  });

  panel.dodajCzynnosc('Zapisz poświadczenie', () => {
    if (nazwaSekretu.value.trim() === '' || wartoscSekretu.value === '') {
      panel.tresc.potwierdzenie(
        'Poświadczenie wymaga nazwy, po której przywołuje je krok, oraz wartości.',
        false,
      );
      return;
    }
    const zadanie: Parameters<ZrodloAutomations['zapiszPoswiadczenie']>[0] = {
      name: nazwaSekretu.value.trim(), value: wartoscSekretu.value,
    };
    // Pole wartości czyścimy od razu po złożeniu żądania, bo trafiła już do skarbca.
    wartoscSekretu.value = '';
    wykonajCzynnoscPanelu(panel, 'Zapis poświadczenia…', zrodlo.zapiszPoswiadczenie(zadanie),
      'Rdzeń nie zapisał poświadczenia.',
      'Rdzeń zapisał poświadczenie w skarbcu i oddał jego REFERENCJĘ; wartość nie wraca nigdy.');
  });

  panel.dodajCzynnosc('Wykaz poświadczeń', () => {
    wykonajCzynnoscPanelu(panel, 'Odczyt referencji…', zrodlo.poswiadczenia({}),
      'Rdzeń nie oddał referencji poświadczeń.',
      'Rdzeń oddał referencje poświadczeń dostępnych krokom; wartości nie wracają.');
  });

  panel.dodajCzynnosc('Usuń poświadczenie', () => {
    if (odwolanieSekretu.value.trim() === '') {
      panel.tresc.potwierdzenie('Wskaż referencję poświadczenia usuwanego.', false);
      return;
    }
    wykonajCzynnoscPanelu(panel, 'Usuwanie poświadczenia…',
      zrodlo.usunPoswiadczenie({ secretRef: odwolanieSekretu.value.trim() }),
      'Rdzeń nie usunął poświadczenia.',
      'Rdzeń usunął poświadczenie; kroki, które je przywoływały, nazywa odpowiedź — ' +
        'usunięcie nie było nimi wstrzymywane.');
  });

  panel.dodajCzynnosc('Dziennik audytu', () => {
    const zadanie: Parameters<ZrodloAutomations['odczytajAudyt']>[0] = {};
    if (stan.automatyka() !== '') zadanie.workflowId = stan.automatyka();
    const od = odczytajChwile(odChwili.value);
    if (od !== undefined) zadanie.fromAt = od;
    const doChwiliWartosc = odczytajChwile(doChwili.value);
    if (doChwiliWartosc !== undefined) zadanie.toAt = doChwiliWartosc;
    const limit = odczytajLiczbe(granica.value);
    if (limit !== undefined) zadanie.limit = limit;
    wykonajCzynnoscPanelu(panel, 'Odczyt dziennika audytu…', zrodlo.odczytajAudyt(zadanie),
      'Rdzeń nie oddał dziennika audytu.',
      'Rdzeń oddał wpisy audytu od najnowszego: kto i kiedy zmienił definicję, ' +
        'uruchomił przebieg albo zmienił harmonogram.');
  });

  return { element: panel.element };
}
