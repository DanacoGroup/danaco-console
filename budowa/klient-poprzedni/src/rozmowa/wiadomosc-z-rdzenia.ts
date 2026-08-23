import { MessageRole, MessageStatus, type Message, type WindowRole } from '../../../shared/contract';
import type { HistoriaTur } from './historia-tur';
import { kluczMiejscowy } from './historia-tur';
import { rozpoznajNadawce } from './nadawca';
import type { StanRozmowy } from './stan-rozmowy';
import { czyPustaTura, wpisOperatora, type StanWpisu, type WpisRozmowy } from './wpis-rozmowy';

/**
 * Przyjęcie zdarzenia `message.changed` do wątku okna.
 *
 * Mieszka tu całe przełożenie wiadomości kontraktu na wpis wątku: dwie
 * równorzędne gałęzie — wypowiedź roli `user` i domknięcie tury modelu — wraz
 * z obroną przed podwojeniem wpisu. Rozdzielenie od `rozmowa.ts` biegnie wzdłuż
 * odpowiedzialności, nie wzdłuż długości pliku.
 *
 * Plik nie wie, kto wypowiedź napisał, i nie zgaduje. Kontrakt nie niesie
 * sprawcy zmiany — ani `Message`, ani koperta zdarzenia nie mają pola
 * połączenia — więc jedyne, co da się udowodnić o wypowiedzi nieznanej wątkowi,
 * to że nie powstała w tym połączeniu. Tyle mówi persona i ani słowa więcej.
 */

/** Czego przyjęcie wiadomości potrzebuje od okna rozmowy. */
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
 * Tożsamość wypowiedzi, która przyszła innym połączeniem tego konta.
 *
 * Widok składa ją z etykietą rodzaju nadawcy w jeden napis — „Operator ·
 * spoza tego połączenia" (`widok-wpisu.ts`). To jest cała prawda, jaką da się
 * o tej wypowiedzi powiedzieć bez zmiany kontraktu: wiadomo, że nie powstała
 * tutaj; nie wiadomo, czy wpisał ją asystent, czy Operator z drugiego
 * urządzenia. Napis „Asystent" byłby wygodniejszy i nieprawdziwy.
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
 * Wypowiedź roli `user`.
 *
 * Autorem wypowiedzi Operatora bywa nie tylko sam Operator: prompt wpisuje też
 * asystent, innym połączeniem WebSocket, przez MCP. Wypowiedź wchodzi więc do
 * wątku zawsze — inaczej pytanie, na które model odpowiada, nie byłoby widoczne
 * ani w polu wypowiedzi, ani w wątku.
 *
 * Przed podwojeniem broni jej wiązanie po treści: echo miejscowe, założone
 * w `wyslij`, dostaje identyfikator z rdzenia zamiast drugiej pozycji obok.
 */
function przyjmijWypowiedz(o: OtoczenieWiadomosci, message: Message): void {
  if (o.historia.znajdz(message.id) !== undefined) return;
  if (o.historia.zwiazWypowiedz(message.content, message.id) !== undefined) return;

  // Wpis składamy wprost, a nie przez `historia.dlaWiadomosci`: tamta droga
  // przejmuje wpis oczekujący, czyli ramkę przygotowaną na odpowiedź modelu.
  // Wypowiedź wjechałaby wtedy w miejsce odpowiedzi i tura straciłaby swoją
  // pozycję w wątku.
  const wpis = wpisOperatora(kluczMiejscowy('zzewnatrz'), message.content, PERSONA_Z_ZEWNATRZ);
  wpis.idWiadomosci = message.id;
  wpis.znacznikCzasu = message.createdAt;
  o.oglos(o.historia.dodaj(wpis));
  o.zglosWypowiedzZZewnatrz(message.content);
}

/** Domknięcie tury modelu. */
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
    // Bez tego wiersza tura domknięta samym zdarzeniem zmiany wiadomości —
    // czyli taka, do której nie doszedł ani jeden fragment strumienia —
    // zostawiłaby wskaźnik na „Wysyłanie…" na zawsze.
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
 * Stan wpisu domkniętego zdarzeniem zmiany wiadomości.
 *
 * Rozstrzygnięcia są cztery, nie dwa: tura zamknięta błędem i tura, która nie
 * przyniosła ani jednego znaku, mają własne stany. Bez nich obie wyglądałyby na
 * ekranie tak samo jak udana — pustą ramką z napisem „zakończona".
 */
function stanDomknietej(status: MessageStatus, wpis: WpisRozmowy): StanWpisu {
  if (wpis.bledy.length > 0 || status === MessageStatus.Error) return 'bledny';
  if (status === MessageStatus.Stopped) return 'przerwany';
  return czyPustaTura(wpis) ? 'pusty' : 'zakonczony';
}
