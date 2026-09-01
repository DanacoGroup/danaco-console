/**
 * Rozdzielacz zdarzeń rdzenia: mapa typ zdarzenia → uchwyty, z rejestracją
 * przez wiązania modułów. Zdarzenie jest jedynym nośnikiem synchronizacji —
 * zmiana zrobiona w jednym oknie albo przez proces w tle dochodzi do pozostałych
 * okien wyłącznie nim. Jedna subskrypcja na typ, wspólna dla wszystkich wiązań,
 * trzyma odbiór w jednym miejscu i pozwala nazwać zdarzenie, którego nikt nie odbiera.
 */
import {
  EventType,
  NotificationWeight,
  type Envelope,
  type EventPayloadOf,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { sesjaKlienta } from '../protokol/sesja.ts';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { oglos } from './ogloszenie.ts';

/** Uchwyt zdarzenia jednego typu, wywoływany z treścią w kształcie z kontraktu. */
export type UchwytZdarzenia<K extends EventType> = (
  tresc: EventPayloadOf<K>,
  koperta: Envelope,
) => void;

/**
 * Zdarzenia warstwy wspólnej — te, które rdzeń stawia niezależnie od modułu.
 * Rozdzielacz podpina je od razu, żeby brak odbiorcy był widoczny, a nie cichy.
 */
const ZDARZENIA_WARSTWY_WSPOLNEJ: readonly EventType[] = [
  EventType.SessionChanged,
  EventType.SessionFocusChanged,
  EventType.WindowChanged,
  EventType.WindowStateChanged,
  EventType.ConfigChanged,
  EventType.QueueChanged,
  EventType.ProgressChanged,
  EventType.NotificationRaised,
  EventType.AuthChanged,
  EventType.DeviceChanged,
  EventType.AlertTriggered,
];

/* Uchwyty trzymane bez typu treści: mapa obejmuje wszystkie typy zdarzeń naraz,
   a kształt treści wiąże dopiero `zglosUchwyt` po stronie zgłaszającego. */
type UchwytDowolny = (tresc: unknown, koperta: Envelope) => void;

const UCHWYTY = new Map<EventType, Set<UchwytDowolny>>();
const SUBSKRYPCJE = new Map<EventType, Odsubskrybuj>();
/* Zdarzenie bez uchwytu nazywane jest raz na typ: rdzeń stawia je w pętli
   (postęp procesu, kolejka), a powtarzany wiersz nie niesie nic ponad pierwszy. */
const NAZWANE_BEZ_UCHWYTU = new Set<EventType>();
let kanalZdarzen: Kanal | null = null;

/**
 * Podpina rozdzielacz do kanału i zgłasza uchwyty warstwy wspólnej. Wołane raz,
 * przy stawianiu aplikacji; drugie wołanie zostawia pierwsze wiązanie nietknięte.
 */
export function zwiazZdarzenia(kanal: Kanal): void {
  if (kanalZdarzen !== null) {
    console.warn('[zdarzenia] rozdzielacz jest już związany z kanałem');
    return;
  }
  kanalZdarzen = kanal;
  zglosUchwyt(EventType.SessionFocusChanged, zapiszOgnisko);
  zglosUchwyt(EventType.NotificationRaised, ogloszenieZdarzenia);
  zglosUchwyt(EventType.AlertTriggered, ogloszenieAlertu);
  for (const zdarzenie of ZDARZENIA_WARSTWY_WSPOLNEJ) {
    podepnij(zdarzenie);
  }
}

/** Zgłasza uchwyt zdarzenia; zwraca odłączenie. Uchwyt zgłoszony przed związaniem kanału czeka na nie. */
export function zglosUchwyt<K extends EventType>(
  zdarzenie: K,
  uchwyt: UchwytZdarzenia<K>,
): Odsubskrybuj {
  const zestaw = UCHWYTY.get(zdarzenie) ?? new Set<UchwytDowolny>();
  UCHWYTY.set(zdarzenie, zestaw);
  const zgloszony = uchwyt as UchwytDowolny;
  zestaw.add(zgloszony);
  podepnij(zdarzenie);
  return () => {
    zestaw.delete(zgloszony);
  };
}

/** Zakłada w kanale subskrypcję jednego typu zdarzenia, jeżeli jeszcze nie stoi. */
function podepnij(zdarzenie: EventType): void {
  if (kanalZdarzen === null || SUBSKRYPCJE.has(zdarzenie)) return;
  SUBSKRYPCJE.set(
    zdarzenie,
    kanalZdarzen.naZdarzenie(zdarzenie, (tresc, koperta) => rozdziel(zdarzenie, tresc, koperta)),
  );
}

/** Rozdaje zdarzenie jego uchwytom; błąd jednego uchwytu nie zabiera zdarzenia pozostałym. */
function rozdziel(zdarzenie: EventType, tresc: unknown, koperta: Envelope): void {
  const zestaw = UCHWYTY.get(zdarzenie);
  if (zestaw === undefined || zestaw.size === 0) {
    if (!NAZWANE_BEZ_UCHWYTU.has(zdarzenie)) {
      NAZWANE_BEZ_UCHWYTU.add(zdarzenie);
      console.warn('[zdarzenia] zdarzenie bez uchwytu', zdarzenie);
    }
    return;
  }
  for (const uchwyt of [...zestaw]) {
    try {
      uchwyt(tresc, koperta);
    } catch (blad) {
      console.error('[zdarzenia] błąd uchwytu', zdarzenie, blad);
    }
  }
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
