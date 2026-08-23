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
 * Observability Tools — kontener pięciu narzędzi obserwowalności.
 *
 * Opracowanie modułu umieszcza w warstwie funkcji eksperckich znaczną część
 * diagnostyki: prowenancję wywołań modeli, zużycie i koszt, metryki
 * i wydajność, kontrolę stanu oraz alerty. Kontrakt niesie z tego jedną
 * rodzinę: `monitor.*` wraz ze zdarzeniem `progress.changed`. Okno jest więc
 * podzielone dokładnie tak, jak opracowanie, i przy każdej zakładce mówi,
 * co pod nią stoi:
 *
 *   Provenance Explorer   — rodzina `provenance.*`: wykaz i odczyt wywołań
 *                           modeli, ocena wywołania przez Operatora i wydanie
 *                           śladu. Źródło wchodzi osobnym argumentem, bo nie
 *                           należy do portu Diagnostics; bez niego zakładka
 *                           nazywa ten brak po stronie złożenia modułu,
 *   Usage & Cost          — rodzina `usage.*`: zestawienie zużycia w jednym
 *                           wymiarze i raport rozliczeniowy wytwarzający plik,
 *   Metrics & Performance — telemetria procesów rdzenia, na żywo,
 *   Health & Uptime       — ta sama telemetria czytana jako kondycja,
 *   Alerts                — bez rodziny komend i bez ustawień progów.
 *
 * Zakładka bez pokrycia nie znika ze sceny i nie zostaje wygaszona: znika
 * dopiero funkcja, o której nikt nie wie, że jej nie ma. Powód każdej składa
 * się z wykazu komend kontraktu przy składaniu okna (`braki-kontraktu.ts`),
 * więc pojawienie się rodziny przepisze te zdania samo.
 *
 * Opracowanie wiąże widoczność zakładek z kluczami `diagnostics.narzedzia.*`.
 * Rdzeń nie zna dziś ani jednego ustawienia obszaru `diagnostics`, więc okno
 * nie udaje, że czyta nastawę — pokazuje pięć zakładek i mówi to wprost.
 * Ukrycie zakładki bez ustawienia byłoby ukryciem opartym o wartość zmyśloną.
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
  // Prowenancja wchodzi osobnym źródłem, bo rodzina `provenance.*` nie należy
  // do portu Diagnostics — zakładka bez niej nie ma czym wołać rdzenia i mówi
  // to wprost, jako brak po stronie złożenia modułu.
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

  // Prowenancja pyta rdzeń dopiero po otwarciu zakładki: to drugie
  // przeszukanie dziennika obok Logs Viewera i nie ma powodu robić go, zanim
  // Operator na tę zakładkę spojrzy.
  // Zużycie pyta rdzeń tą samą regułą i z tego samego powodu: zestawienie za
  // okres jest osobnym przeliczeniem po stronie rdzenia, więc idzie dopiero,
  // gdy Operator na tę zakładkę spojrzy.
  zakladki.naZmiane((kod) => {
    if (kod === 'provenance' && !prowenancja.czytano()) prowenancja.odswiez();
    if (kod === 'usage-cost' && !zakladkaZuzycia.czytano()) zakladkaZuzycia.odswiez();
  });

  // Zakres czasu wspólny modułowi dotyczy prowenancji i zużycia — telemetria
  // procesów nie ma w kontrakcie pól zakresu. Odczyt tylko wtedy, gdy zakładka
  // już czytała: przeliczanie zakładki nigdy nie otwartej byłoby zapytaniem
  // o materiał, którego nikt nie ogląda.
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
 * Alerts — zakładka bez pokrycia w kontrakcie.
 *
 * Reguła alertu wymaga dwóch rzeczy naraz: miejsca na jej zapis i ewaluatora,
 * który ją sprawdzi. Kontrakt nie ma ani rodziny komend alertów, ani ustawień
 * obszaru diagnostics, w których próg mógłby zamieszkać. Reguła zbudowana
 * wyłącznie w oknie żyłaby do zamknięcia karty i nie zadziałałaby ani razu,
 * gdy Operator nie patrzy — czyli dokładnie wtedy, gdy alert ma sens.
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
