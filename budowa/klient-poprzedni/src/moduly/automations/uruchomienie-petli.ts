/**
 * Uruchomienie gotowej pętli jednym kliknięciem, złożone z dwóch przebiegów do
 * rdzenia. Plik składa dwie komendy kontraktu i oddaje to, co rdzeń powiedział;
 * kolejki dalej nie posuwa i jej stanów nie zna.
 */
import {
  QueueAction,
  type AutomationWorkflow,
  type ErrorInfo,
  type Queue,
} from '../../../../shared/contract';
import type { ZrodloAutomations } from './zrodlo-automations';

/**
 * Skutek uruchomienia pętli: kolejka oddana przez rdzeń po działaniu albo
 * zdanie odmowy wraz z powodem podanym przez rdzeń, gdy rdzeń powód podał.
 */
export type SkutekUruchomienia =
  | { udane: true; kolejka: Queue }
  | { udane: false; zdanie: string; powod?: ErrorInfo };

/**
 * Zakłada kolejkę pętli i posuwa ją działaniem rozpoczęcia. Sesja kolejki jest
 * ta sama, co w oknie Queue Manager, ponieważ pętla uruchamiana z wykazu nie ma
 * własnej karty sesji.
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
 * Zdanie potwierdzenia powstaje z oddanej kolejki, a nie z żądania. Okno nie
 * orzeka, że pętla ruszyła, ponieważ odpowiedź niesie stan kolejki po działaniu,
 * lecz nie mówi, które kroki się wykonały.
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

/**
 * Składa odmowę uruchomienia wraz z powodem rdzenia, gdy rdzeń go podał. Powód
 * wraca osobnym polem, ponieważ zdanie okna dokłada do niego kod i komunikat
 * kontraktu.
 */
function odmowa(zdanie: string, powod?: ErrorInfo): SkutekUruchomienia {
  return powod === undefined ? { udane: false, zdanie } : { udane: false, zdanie, powod };
}
