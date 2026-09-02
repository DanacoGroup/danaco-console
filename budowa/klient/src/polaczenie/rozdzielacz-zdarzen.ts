/**
 * Rozdzielacz zdarzeń rdzenia — jedna subskrypcja na typ zdarzenia, wspólna
 * wszystkim wiązaniom. Zdarzenie jest jedynym nośnikiem synchronizacji: zmiana
 * zrobiona w innym oknie albo przez proces w tle dochodzi tutaj wyłącznie nim.
 * Rozdzielacz stoi w warstwie połączenia, nad źródłem komunikatów.
 */
import { EventType, type Envelope, type EventPayloadOf } from '../../../shared/contract.ts';
import type { Odsubskrybuj } from './magistrala-zdarzen.ts';
import type { ZrodloZdarzen } from './zrodlo-zdarzen.ts';

/** Uchwyt zdarzenia jednego typu, wywoływany z treścią w kształcie z kontraktu. */
export type UchwytZdarzenia<K extends EventType> = (
  tresc: EventPayloadOf<K>,
  koperta: Envelope,
) => void;

/**
 * Zdarzenia warstwy wspólnej — te, które rdzeń stawia niezależnie od modułu.
 * Rozdzielacz subskrybuje je od chwili postawienia aplikacji: subskrypcja
 * zakładana dopiero przy montażu okna gubiłaby wszystko, co przyszło wcześniej.
 */
export const ZDARZENIA_WARSTWY_WSPOLNEJ: readonly EventType[] = [
  EventType.SessionChanged,
  EventType.SessionFocusChanged,
  EventType.WindowChanged,
  EventType.WindowStateChanged,
  EventType.MessageChanged,
  EventType.StreamChunk,
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
let zrodloZdarzen: ZrodloZdarzen | null = null;

/**
 * Podpina rozdzielacz do źródła komunikatów i zakłada subskrypcje warstwy
 * wspólnej. Wołane raz, przy stawianiu aplikacji; drugie wołanie zostawia
 * pierwsze wiązanie nietknięte, bo subskrypcje są jedne na typ.
 */
export function zwiazRozdzielacz(zrodlo: ZrodloZdarzen): void {
  if (zrodloZdarzen !== null) return;
  zrodloZdarzen = zrodlo;
  for (const zdarzenie of ZDARZENIA_WARSTWY_WSPOLNEJ) {
    podepnij(zdarzenie);
  }
}

/** Zgłasza uchwyt zdarzenia; zwraca odłączenie. Uchwyt zgłoszony przed związaniem źródła czeka na nie. */
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

/** Zakłada w źródle subskrypcję jednego typu zdarzenia, jeżeli jeszcze nie stoi. */
function podepnij(zdarzenie: EventType): void {
  if (zrodloZdarzen === null || SUBSKRYPCJE.has(zdarzenie)) return;
  SUBSKRYPCJE.set(
    zdarzenie,
    zrodloZdarzen.naZdarzenie(zdarzenie, (tresc, koperta) => rozdziel(zdarzenie, tresc, koperta)),
  );
}

/** Rozdaje zdarzenie jego uchwytom; błąd jednego uchwytu nie zabiera zdarzenia pozostałym. */
function rozdziel(zdarzenie: EventType, tresc: unknown, koperta: Envelope): void {
  const zestaw = UCHWYTY.get(zdarzenie);
  if (zestaw === undefined) return;
  for (const uchwyt of [...zestaw]) {
    try {
      uchwyt(tresc, koperta);
    } catch (blad) {
      console.error('[zdarzenia] błąd uchwytu', zdarzenie, blad);
    }
  }
}
