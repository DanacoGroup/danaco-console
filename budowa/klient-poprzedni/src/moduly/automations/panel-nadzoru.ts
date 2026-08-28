/**
 * Panel Schedulera obsługuje pięć czynności rodziny schedule, których stan musi przeżyć restart
 * rdzenia: okna wykonania, wsteczne uruchomienie, historię i nadzór.
 */
import { pole, poleLiczbowe, poleTresci } from '../../modele/kontrolki-formularza';
import {
  odczytajChwile,
  odczytajLiczbe,
  odczytajZapis,
  utworzPanelDobudowy,
  wykonajCzynnoscPanelu,
} from './panel-dobudowy';
import type { StanAutomatyki } from './stan-automatyki';
import type { ZrodloAutomations } from './zrodlo-automations';

/** Panel dopełniający Schedulera, osadzany wewnątrz okna harmonogramów wraz z jego szufladami czynności. */
export interface PanelNadzoru {
  element: HTMLElement;
}

export function utworzPanelNadzoru(
  zrodlo: ZrodloAutomations,
  stan: StanAutomatyki,
): PanelNadzoru {
  const panel = utworzPanelDobudowy(
    'Okna wykonania i nadzór uruchomień',
    'Ograniczenie przedziałów uruchomień, uruchomienie wsteczne pominiętych ' +
      'terminów, historia rzeczywistych wyzwoleń, alarm przy braku uruchomienia ' +
      'oraz adres wejściowy wyzwalacza webhook.',
  );

  const idHarmonogramu = panel.dodajPole(
    'Harmonogram',
    pole('Harmonogram', 'harm-…'),
    'Identyfikator harmonogramu z pola „id” odczytu — nie identyfikator automatyki.',
  );
  const odChwili = panel.dodajPole('Od', pole('Od', '2026-07-01T00:00:00Z'));
  const doChwili = panel.dodajPole('Do', pole('Do', '2026-08-01T00:00:00Z'));
  const granica = panel.dodajPole('Granica przebiegów', poleLiczbowe('Granica przebiegów', '10'));
  const tolerancja = panel.dodajPole(
    'Tolerancja (sekundy)',
    poleLiczbowe('Tolerancja', '900'),
    'Zero wyłącza nadzór obecności uruchomień.',
  );
  const regulaAlarmu = panel.dodajPole('Reguła alarmowania', pole('Reguła alarmowania', 'alarm-…'));
  const okna = panel.dodajPole(
    'Okna wykonania',
    poleTresci('Okna wykonania', 4, '[{ "fromMinuteOfDay": 480, "toMinuteOfDay": 1020 }]'),
    'Komplet okien po zmianie; wykaz pusty znaczy „bez ograniczenia”.',
  );

  /** Harmonogram wskazany w polu; brak wstrzymuje czynność z nazwanym powodem. */
  function harmonogram(): string | null {
    if (idHarmonogramu.value.trim() === '') {
      panel.tresc.potwierdzenie(
        'Wskaż harmonogram. Jego identyfikator oddaje odczyt harmonogramów w oknie Schedulera — ' +
          'te czynności dotyczą harmonogramu, nie automatyki.',
        false,
      );
      return null;
    }
    return idHarmonogramu.value.trim();
  }

  panel.dodajCzynnosc('Zapisz okna wykonania', () => {
    const kod = harmonogram();
    if (kod === null) return;
    const odczytane = odczytajZapis(okna.value);
    if (odczytane === null) {
      panel.tresc.potwierdzenie('Zapis okien wykonania jest nieczytelny — popraw go.', false);
      return;
    }
    const wykaz = (odczytane ?? []) as Parameters<
      ZrodloAutomations['ustawOknaWykonania']
    >[0]['windows'];
    wykonajCzynnoscPanelu(panel, 'Zapis okien wykonania…',
      zrodlo.ustawOknaWykonania({ scheduleId: kod, windows: wykaz }),
      'Rdzeń nie zapisał okien wykonania.',
      'Rdzeń zapisał okna wykonania; ograniczenie obowiązuje także po ponownym złożeniu rdzenia.');
  });

  panel.dodajCzynnosc('Uruchom wstecznie', () => {
    const kod = harmonogram();
    if (kod === null) return;
    const od = odczytajChwile(odChwili.value);
    const doTerminu = odczytajChwile(doChwili.value);
    if (od === undefined || doTerminu === undefined) {
      panel.tresc.potwierdzenie(
        'Uruchomienie wsteczne wymaga obu granic zakresu dat w polach „Od” i „Do”.',
        false,
      );
      return;
    }
    const zadanie: Parameters<ZrodloAutomations['uruchomWstecznie']>[0] = {
      scheduleId: kod, fromAt: od, toAt: doTerminu,
    };
    const limit = odczytajLiczbe(granica.value);
    if (limit !== undefined) zadanie.maxRuns = limit;
    wykonajCzynnoscPanelu(panel, 'Uruchomienie wsteczne…', zrodlo.uruchomWstecznie(zadanie),
      'Rdzeń nie wykonał przebiegów wstecznych.',
      'Rdzeń zakolejkował przebiegi dla pominiętych terminów; terminy, dla których przebieg ' +
        'już był, policzył jako pominięte.');
  });

  panel.dodajCzynnosc('Historia wyzwoleń', () => {
    const zadanie: Parameters<ZrodloAutomations['historiaWyzwolen']>[0] = {};
    if (idHarmonogramu.value.trim() !== '') zadanie.scheduleId = idHarmonogramu.value.trim();
    else if (stan.automatyka() !== '') zadanie.workflowId = stan.automatyka();
    wykonajCzynnoscPanelu(panel, 'Odczyt historii wyzwoleń…', zrodlo.historiaWyzwolen(zadanie),
      'Rdzeń nie oddał historii wyzwoleń.',
      'Rdzeń oddał rzeczywiste momenty wyzwolenia wraz z przyczyną.');
  });

  panel.dodajCzynnosc('Zapisz nadzór uruchomień', () => {
    const kod = harmonogram();
    if (kod === null) return;
    const okno = odczytajLiczbe(tolerancja.value);
    if (okno === undefined) {
      panel.tresc.potwierdzenie('Wskaż okno tolerancji w sekundach; zero wyłącza nadzór.', false);
      return;
    }
    const zadanie: Parameters<ZrodloAutomations['ustawNadzorUruchomien']>[0] = {
      scheduleId: kod, toleranceSeconds: okno,
    };
    if (regulaAlarmu.value.trim() !== '') zadanie.alertRuleId = regulaAlarmu.value.trim();
    wykonajCzynnoscPanelu(panel, 'Zapis nadzoru…', zrodlo.ustawNadzorUruchomien(zadanie),
      'Rdzeń nie zapisał nadzoru obecności uruchomień.',
      okno === 0
        ? 'Rdzeń wyłączył nadzór obecności uruchomień.'
        : 'Rdzeń zapisał okno tolerancji; alarm zadziała, gdy uruchomienie w nim nie nastąpi.');
  });

  panel.dodajCzynnosc('Adres webhooka', () => {
    if (stan.automatyka() === '') {
      panel.tresc.potwierdzenie(
        'Wskaż automatykę — adres wejściowy webhooka należy do automatyki, nie do harmonogramu.',
        false,
      );
      return;
    }
    wykonajCzynnoscPanelu(panel, 'Odczyt adresu webhooka…',
      zrodlo.adresWebhooka({ workflowId: stan.automatyka() }),
      'Rdzeń nie oddał adresu wejściowego webhooka.',
      'Rdzeń oddał adres wejściowy i REFERENCJĘ klucza podpisu; sama wartość klucza nie wraca nigdy.');
  });

  panel.dodajCzynnosc('Wymień klucz podpisu', () => {
    if (stan.automatyka() === '') {
      panel.tresc.potwierdzenie('Wskaż automatykę, której klucz podpisu ma zostać wymieniony.', false);
      return;
    }
    wykonajCzynnoscPanelu(panel, 'Wymiana klucza podpisu…',
      zrodlo.adresWebhooka({ workflowId: stan.automatyka(), rotateSecret: true }),
      'Rdzeń nie wymienił klucza podpisu.',
      'Rdzeń wymienił klucz podpisu i oddał jego nową referencję; klucza poprzedniego ' +
        'nie da się już odczytać żadną komendą.');
  });

  return { element: panel.element };
}
