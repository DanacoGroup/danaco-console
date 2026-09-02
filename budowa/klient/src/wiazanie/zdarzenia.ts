/**
 * Uchwyty warstwy wspólnej zgłaszane do rozdzielacza zdarzeń. Sam rozdzielacz
 * stoi w warstwie połączenia (`polaczenie/rozdzielacz-zdarzen.ts`); tutaj
 * zostaje to, co ze zdarzeniem robi interfejs — sesja bieżąca i komunikaty
 * pokazywane Operatorowi.
 */
import { EventType, NotificationWeight, type EventPayloadOf } from '../../../shared/contract.ts';
import { zglosUchwyt, zwiazRozdzielacz } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { sesjaKlienta } from '../protokol/sesja.ts';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { oglos } from './ogloszenie.ts';

export { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
export type { UchwytZdarzenia } from '../polaczenie/rozdzielacz-zdarzen.ts';

/** Wiąże rozdzielacz z kanałem i zgłasza uchwyty warstwy wspólnej. Wołane raz, przy stawianiu aplikacji. */
export function zwiazZdarzenia(kanal: Kanal): void {
  zwiazRozdzielacz(kanal);
  zglosUchwyt(EventType.SessionFocusChanged, zapiszOgnisko);
  zglosUchwyt(EventType.NotificationRaised, ogloszenieZdarzenia);
  zglosUchwyt(EventType.AlertTriggered, ogloszenieAlertu);
}

/**
 * Sesja bieżąca klienta idzie za ogniskiem potwierdzonym przez rdzeń.
 * Zdarzenie niesie klienta, na którym ognisko przeniesiono — ognisko cudzego
 * urządzenia tego samego konta nie przestawia sesji nadawanej tutejszym kopertom.
 */
function zapiszOgnisko(tresc: EventPayloadOf<typeof EventType.SessionFocusChanged>): void {
  if (tresc.clientId !== tozsamoscKlienta().id) return;
  sesjaKlienta().ustaw(tresc.sessionId);
}

/** Zdarzenie rejestru centrum staje Operatorowi przed oczami w chwili wystąpienia; kontrakt czyni je podstawą komunikatu. */
function ogloszenieZdarzenia(tresc: EventPayloadOf<typeof EventType.NotificationRaised>): void {
  const wymagaDecyzji = tresc.notification.weight === NotificationWeight.WymagajacaDecyzji;
  oglos('Powiadomienie', tresc.notification.text, wymagaDecyzji ? 'ostrzezenie' : 'informacja');
}

/** Wyzwolenie reguły alertu jest ostrzeżeniem: opisuje warunek, który już zaszedł na maszynie Operatora. */
function ogloszenieAlertu(tresc: EventPayloadOf<typeof EventType.AlertTriggered>): void {
  oglos(tresc.rule.name, tresc.trigger.message, 'ostrzezenie');
}
