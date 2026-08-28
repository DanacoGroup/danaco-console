import { MessageRole, MessageStatus, type Message, type WindowRole } from '../../../shared/contract';
import type { HistoriaTur } from './historia-tur';
import { kluczMiejscowy } from './historia-tur';
import { rozpoznajNadawce } from './nadawca';
import type { StanRozmowy } from './stan-rozmowy';
import { czyPustaTura, wpisOperatora, type StanWpisu, type WpisRozmowy } from './wpis-rozmowy';

/**
 * Przyjęcie zdarzenia message.changed do wątku okna; przełożenie wiadomości kontraktu na
 * wpis.
 */

/**
 * Zależności, których przyjęcie wiadomości do wątku okna rozmowy potrzebuje od warstwy
 * nadrzędnej rozmowy.
 */
export interface OtoczenieWiadomosci {
  historia: HistoriaTur;
  /** Rola okna w pętli; rozstrzyga, którym z dziewięciu nadawców jest model. */
  rolaOkna(): WindowRole | null;
  /** Tożsamość mówiącego po stronie modelu. */
  persona: string;
  /** Ogłasza wpis założony albo zmieniony. */
  oglos(wpis: WpisRozmowy): void;
  /** Przestawia wskaźnik wysyłania okna. */
  ustawStan(zmiana: Partial<StanRozmowy>): void;
  /** Zgłasza, że wypowiedź weszła do okna spoza tego połączenia. */
  zglosWypowiedzZZewnatrz(tresc: string): void;
}

/**
 * Tożsamość wypowiedzi, która przyszła do okna innym połączeniem tego samego konta
 * Operatora niż bieżące.
 */
export const PERSONA_Z_ZEWNATRZ = 'spoza tego połączenia';

export function przyjmijWiadomosc(o: OtoczenieWiadomosci, message: Message): void {
  if (message.role === MessageRole.User) {
    przyjmijWypowiedz(o, message);
    return;
  }
  przyjmijTureModelu(o, message);
}

/**
 * Wypowiedź roli user, wchodząca do wątku rozmowy zawsze, niezależnie od tego, kto
 * naprawdę ją wpisał.
 */
function przyjmijWypowiedz(o: OtoczenieWiadomosci, message: Message): void {
  if (o.historia.znajdz(message.id) !== undefined) return;
  if (o.historia.zwiazWypowiedz(message.content, message.id) !== undefined) return;

  // Wpis składany wprost, nie przez historię wiadomości, żeby nie przejąć wpisu czekającego
  // na model.
  const wpis = wpisOperatora(kluczMiejscowy('zzewnatrz'), message.content, PERSONA_Z_ZEWNATRZ);
  wpis.idWiadomosci = message.id;
  wpis.znacznikCzasu = message.createdAt;
  o.oglos(o.historia.dodaj(wpis));
  o.zglosWypowiedzZZewnatrz(message.content);
}

/**
 * Domknięcie tury modelu zdarzeniem zmiany wiadomości kontraktu, kończące jej
 * dotychczasowy stan wysyłania.
 */
function przyjmijTureModelu(o: OtoczenieWiadomosci, message: Message): void {
  const znany = o.historia.znajdz(message.id);
  if (znany === undefined) {
    const wpis = o.historia.dlaWiadomosci(
      message.id,
      rozpoznajNadawce(message.role, o.rolaOkna()),
      o.persona,
    );
    wpis.tresc = message.content;
    wpis.stan = stanDomknietej(message.status, wpis);
    wpis.domkniety = true;
    o.oglos(wpis);
    // Bez tego wiersza tura bez fragmentu strumienia zostałaby wskaźnikiem wysyłania na
    // zawsze.
    o.ustawStan({
      wysyla: false,
      idWiadomosci: message.id,
      fragmenty: wpis.fragmenty,
      ciszaSekundy: 0,
    });
    return;
  }
  znany.stan = stanDomknietej(message.status, znany);
  znany.domkniety = true;
  o.oglos(znany);
  o.ustawStan({ wysyla: false, ciszaSekundy: 0 });
}

/**
 * Stan wpisu domkniętego zdarzeniem zmiany wiadomości, rozróżniający cztery różne
 * rozstrzygnięcia tury.
 */
function stanDomknietej(status: MessageStatus, wpis: WpisRozmowy): StanWpisu {
  if (wpis.bledy.length > 0 || status === MessageStatus.Error) return 'bledny';
  if (status === MessageStatus.Stopped) return 'przerwany';
  return czyPustaTura(wpis) ? 'pusty' : 'zakonczony';
}
