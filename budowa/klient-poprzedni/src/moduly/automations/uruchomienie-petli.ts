import {
  QueueAction,
  type AutomationWorkflow,
  type ErrorInfo,
  type Queue,
} from '../../../../shared/contract';
import type { ZrodloAutomations } from './zrodlo-automations';

/**
 * Uruchomienie gotowej pętli jednym kliknięciem — dwa przebiegi do rdzenia
 * złożone w jedną czynność.
 *
 * Dwa, ponieważ `automation.queue.action` wymaga `queueId`, a wykaz automatyk go
 * nie niesie: `AutomationWorkflow` nie ma pola kolejki. Kolejka powstaje więc tą
 * samą drogą, którą zakłada ją Queue Manager — komendą `queue.create` z sesją
 * `automations` — a zaraz po niej idzie działanie `start` ze wskazaniem
 * automatyki. Rdzeń zasila kolejkę krokami wyłącznie wtedy, gdy nie ma ona ani
 * jednej pozycji; kolejka świeżo założona jest pusta, więc uruchomienie startuje
 * z pełnym wykazem kroków.
 *
 * Plik nie posuwa kolejki dalej i nie zna jej stanów — składa dwie komendy
 * kontraktu i oddaje to, co rdzeń powiedział. Odmowa mówi, co się nie udało
 * i czym to naprawić; powód z rdzenia wraca osobnym polem, bo zdanie okna dokłada
 * do niego kod i komunikat kontraktu.
 */

/** Skutek uruchomienia: kolejka po działaniu albo zdanie odmowy wraz z powodem. */
export type SkutekUruchomienia =
  | { udane: true; kolejka: Queue }
  | { udane: false; zdanie: string; powod?: ErrorInfo };

/**
 * Zakłada kolejkę pętli i posuwa ją działaniem `start`.
 *
 * Sesja kolejki jest ta sama co w Queue Managerze (`automations`), bo pętla
 * uruchamiana z wykazu nie ma karty sesji: moduł jest komponentem własnym strefy
 * 2 strony głównej, a nie przestrzeni roboczej sesji.
 */
export async function uruchomPetle(
  zrodlo: ZrodloAutomations,
  petla: AutomationWorkflow,
): Promise<SkutekUruchomienia> {
  const zalozona = await zrodlo.zalozKolejke({
    sessionId: 'automations',
    name: `Automatyka ${petla.id}`,
  });
  if (!zalozona.udany || zalozona.wynik === undefined) {
    return odmowa(
      `Rdzeń nie założył kolejki dla pętli „${petla.name}”, więc uruchomienie nie ruszyło. ` +
        'Pętla wykonuje się na kolejce silnika kolejek, a bez kolejki nie ma czym jej ' +
        'posunąć. Naciśnij „Uruchom” ponownie albo załóż kolejkę w Queue Managerze i wykonaj na ' +
        'niej działanie „Uruchom”.',
      zalozona.blad,
    );
  }

  const kolejka = zalozona.wynik;
  const posunieta = await zrodlo.dzialanieKolejki({
    queueId: kolejka.id,
    action: QueueAction.Start,
    workflowId: petla.id,
  });
  if (!posunieta.udany || posunieta.wynik === undefined) {
    return odmowa(
      `Rdzeń założył kolejkę ${kolejka.id}, ale odmówił jej uruchomienia dla pętli ` +
        `„${petla.name}”. Kroki pętli zasila i posuwa działanie „start” silnika kolejek, ` +
        'i to ono nie przeszło. Kolejka istnieje — wykonaj na niej działanie „Uruchom” ' +
        'w Queue Managerze albo naciśnij „Uruchom” w wykazie jeszcze raz.',
      posunieta.blad,
    );
  }

  return { udane: true, kolejka: posunieta.wynik };
}

/**
 * Zdanie potwierdzenia powstaje z oddanej kolejki, nie z żądania.
 *
 * Okno nie orzeka, że pętla „ruszyła”: odpowiedź niesie stan kolejki po
 * działaniu, licznik obiegów i liczbę zleceń oczekujących, lecz nie mówi, które
 * kroki się wykonały. Twierdzenie o wykonanych krokach byłoby potwierdzeniem
 * czynności, której okno nie zmierzyło.
 */
export function zdanieUruchomienia(petla: AutomationWorkflow, kolejka: Queue): string {
  const kroki = petla.steps?.length ?? 0;
  const oczekujace =
    kolejka.pendingCount === undefined
      ? ''
      : `, zleceń oczekujących ${kolejka.pendingCount}`;
  return (
    `Rdzeń przyjął uruchomienie pętli „${petla.name}” (${kroki} kroków definicji). ` +
    `Kolejka ${kolejka.id} po działaniu: ${kolejka.status}, ` +
    `licznik obiegów naprawczych ${kolejka.cycle ?? 0}${oczekujace}. ` +
    'Przebieg widać w Execution Monitorze.'
  );
}

/** Składa odmowę z powodem rdzenia, gdy rdzeń go podał. */
function odmowa(zdanie: string, powod?: ErrorInfo): SkutekUruchomienia {
  return powod === undefined ? { udane: false, zdanie } : { udane: false, zdanie, powod };
}
