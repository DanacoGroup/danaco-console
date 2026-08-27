import { utworzRameOkna } from '../../komponenty/rama-okna';
import { przyciskAkcji, przyciskBezKomendy } from '../../modele/kontrolki-formularza-braki';
import { powodBezKomendy } from './braki-kontraktu';
import type { StanDiagnostyki } from './stan-diagnostyki';
import { utworzTelemetrieProcesow } from './telemetria-procesow';
import { utworzZakladkeProwenancji } from './zakladka-prowenancji';
import type { ZrodloProwenancji } from './prowenancja-zrodlo';
import { utworzZakladkeZuzycia } from './zuzycie-zakladka';
import type { ZrodloZuzycia } from './zuzycie-zrodlo';
import {
  cialoNarzedzia,
  objasnienieNarzedzia,
  pasekNarzedzia,
  utworzZakladkiNarzedzi,
} from './zakladki-narzedzi';
import type { ZrodloDiagnostics } from './zrodlo-diagnostics';
import type { ZrodloObserwowalnosci } from './zrodlo-obserwowalnosci';

/**
 * Interfejs Observability Tools łączy pięć narzędzi obserwowalności: prowenancję
 * wywołań modeli, zużycie i koszt, metryki i wydajność, kontrolę stanu oraz
 * alerty, każde jako osobna zakładka.
 */
export interface OknoObservabilityTools {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka nasłuch telemetrii oraz nasłuch zakresu wspólnego modułu. */
  zamknij(): void;
}

export function utworzOknoObservabilityTools(
  zrodlo: ZrodloDiagnostics,
  obserwowalnosc: ZrodloObserwowalnosci,
  stan: StanDiagnostyki,
  idOkna: string,
  zuzycie: ZrodloZuzycia,
  // Prowenancja wchodzi osobnym źródłem, bo rodzina komend nie należy do portu diagnostyki.
  prowenancjaZrodlo?: ZrodloProwenancji,
): OknoObservabilityTools {
  const rama = utworzRameOkna({
    tytul: 'Observability Tools',
    rola: 'pomocnicze',
    kod: 'observability-tools',
    przeznaczenie:
      'Kontener pięciu narzędzi obserwowalności: prowenancja wywołań, zużycie i koszt, metryki i wydajność, kontrola stanu, alerty.',
    modul: 'Diagnostics',
    przedrostek: 'dg',
  });

  const telemetria = utworzTelemetrieProcesow(obserwowalnosc, idOkna);
  const prowenancja = utworzZakladkeProwenancji(zrodlo, stan, prowenancjaZrodlo);
  const zakladkaZuzycia = utworzZakladkeZuzycia(zuzycie, stan);

  const zakladki = utworzZakladkiNarzedzi([
    { kod: 'provenance', nazwa: 'Provenance Explorer', element: prowenancja.element },
    zakladkaZuzycia.pozycja,
    { kod: 'metrics', nazwa: 'Metrics & Performance', element: telemetria.metryki },
    { kod: 'health', nazwa: 'Health & Uptime', element: telemetria.kondycja },
    { kod: 'alerts', nazwa: 'Alerts', element: zakladkaAlertow() },
  ]);

  rama.narzedzia.append(
    objasnienieNarzedzia(
      'Widocznością zakładek steruje wg opracowania klucz diagnostics.narzedzia.* — rdzeń nie zna ' +
        'dziś ani jednego ustawienia obszaru diagnostics, więc widocznych jest wszystkie pięć.',
    ),
  );
  rama.akcje.append(...czynnosciKontenera(telemetria.odswiez));
  rama.cialo.append(zakladki.element);

  // Prowenancja i zużycie pytają rdzeń dopiero po otwarciu zakładki, jako osobne przeliczenie.
  zakladki.naZmiane((kod) => {
    if (kod === 'provenance' && !prowenancja.czytano()) prowenancja.odswiez();
    if (kod === 'usage-cost' && !zakladkaZuzycia.czytano()) zakladkaZuzycia.odswiez();
  });

  // Wspólny zakres czasu ma odczyt tylko dla zakładki już otwartej, nigdy dla nieoglądanej.
  const odsubskrybujZakres = stan.naZmiane(() => {
    if (prowenancja.czytano()) prowenancja.odswiez();
    if (zakladkaZuzycia.czytano()) zakladkaZuzycia.odswiez();
  });

  function odswiez(): void {
    telemetria.odswiez();
    if (prowenancja.czytano()) prowenancja.odswiez();
    if (zakladkaZuzycia.czytano()) zakladkaZuzycia.odswiez();
  }

  return {
    element: rama.element,
    odswiez,
    zamknij() {
      odsubskrybujZakres();
      telemetria.zamknij();
    },
  };
}

/**
 * Panel akcji kontenera. Odświeżenie telemetrii jest jedyną czynnością
 * sprawczą całego kontenera — reszta czynności należy do zakładek i stoi przy
 * materiale, którego dotyczy.
 */
function czynnosciKontenera(odswiezTelemetrie: () => void): readonly HTMLElement[] {
  const odswiez = przyciskAkcji('Odśwież telemetrię', 'dn-btn dn-btn--atrament');
  odswiez.addEventListener('click', odswiezTelemetrie);
  return [
    odswiez,
    przyciskBezKomendy(
      'Importuj telemetrię zewnętrzną',
      powodBezKomendy('Wczytanie logów i śladów z systemu zewnętrznego wymagałoby komendy importu telemetrii.'),
    ),
    przyciskBezKomendy(
      'Redakcja danych wrażliwych w śladach',
      powodBezKomendy(
        'Redakcja jest wg opracowania ustawieniem Operatora, a rdzeń nie zna ani jednego ustawienia obszaru diagnostics.',
      ),
    ),
  ];
}

/**
 * Zakładka alertów pozostaje bez pokrycia w kontrakcie: reguła progu wymaga
 * miejsca zapisu i ewaluatora działającego poza otwartą kartą, a rdzeń dziś
 * nie zapewnia żadnego z tych dwóch elementów.
 */
function zakladkaAlertow(): HTMLElement {
  return cialoNarzedzia(
    objasnienieNarzedzia(
      'Reguł alertów nie ma gdzie zapisać ani czym ewaluować: kontrakt nie niesie rodziny komend ' +
        'alertów, a rdzeń nie zna ustawień obszaru diagnostics, w których mieszkałby próg. ' +
        powodBezKomendy('Reguła progowa, wyciszenie i eskalacja wymagałyby własnej rodziny komend.', 'alert'),
    ),
    objasnienieNarzedzia(
      'Reguła zbudowana wyłącznie w oknie żyłaby do zamknięcia karty sesji i nie zadziałałaby ' +
        'wtedy, gdy alert ma sens — gdy Operator nie patrzy. Dlatego zakładka nie stawia ' +
        'formularza reguły, którego rdzeń nie przyjmie.',
    ),
    pasekNarzedzia(
      przyciskBezKomendy(
        'Nowa reguła progowa',
        powodBezKomendy('Zapis reguły alertu wymagałby komendy zapisu reguły.', 'alert'),
      ),
      przyciskBezKomendy(
        'Historia wyzwoleń',
        powodBezKomendy('Rejestr wyzwoleń alertów wymagałby komendy oddającej ten rejestr.', 'alert'),
      ),
      przyciskBezKomendy(
        'Wypychanie do Always On Display',
        powodBezKomendy('Skierowanie alertu do nakładki wymagałoby komendy alertu; rodzina aod.* nie przyjmuje alertu.', 'alert'),
      ),
    ),
  );
}
