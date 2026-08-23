import type { StreamChunkEvent } from '../../../shared/contract';
import { odczytajKonto } from './metadane-konta';
import { obiekt, prawda, tekst, zapisCzytelny } from './odczyt-fragmentu';
import { odczytajPodsumowanie } from './podsumowanie-tury';
import { odczytajProwenancje } from './prowenancja';
import { RodzajFragmentu } from './rodzaje-fragmentow';
import { czyPustaTura, type WpisRozmowy, type WywolanieNarzedzia } from './wpis-rozmowy';

/** Znaczniki toru strumienia niesione przez kopertę, nie przez ładunek. */
export interface ZnacznikiFragmentu {
  /** Numer fragmentu liczony od jedynki (pole `seq`). */
  numer: number;
  /** Prawda w ostatnim fragmencie tury (pole `done`). */
  ostatni: boolean;
}

/**
 * Złożenie tury — dokłada jeden fragment strumienia do wpisu rozmowy.
 *
 * Rozdzielenie rodzajów jest sednem: jedna droga niesie tekst, prowenancję,
 * metadane konta, tok rozumowania, wywołania narzędzi i błąd, a odbiorca
 * rozstrzyga po rodzaju, gdzie fragment trafi. Fragment rodzaju nieznanego nie
 * jest odrzucany — ląduje w treści wpisu z oznaczeniem rodzaju, żeby żadna
 * treść strumienia nie znikła bez śladu.
 *
 * Wpis jest zmieniany w miejscu: historia trzyma jedną tożsamość tury, a widok
 * odświeża tę samą pozycję zamiast dokładać kolejną.
 */
export function dolaczFragment(
  wpis: WpisRozmowy,
  fragment: StreamChunkEvent,
  znaczniki: ZnacznikiFragmentu,
): void {
  wpis.fragmenty += 1;
  wpis.ostatniNumer = Math.max(wpis.ostatniNumer, znaczniki.numer);
  wpis.idWiadomosci = fragment.messageId;
  if (wpis.stan === 'wysylanie') wpis.stan = 'strumien';

  rozdzielRodzaj(wpis, fragment);

  // Podsumowanie czyta się z każdego fragmentu, nie tylko z domykającego:
  // ładunek podsumowania niesie fragment kanału, a strumień domyka osobne
  // zdarzenie z wersją ostateczną (ChunkKind.final), więc odczyt tylko
  // z domykającego gubiłby podsumowanie całkiem. Parser odrzuca ładunek bez
  // własnych pól podsumowania, więc fragmenty pozostałych rodzajów przechodzą
  // bez śladu, a `??=` zatrzymuje pierwsze odczytane.
  wpis.podsumowanie ??= odczytajPodsumowanie(fragment.data);

  if (!znaczniki.ostatni) return;
  wpis.domkniety = true;
  // Stany są trzy, nie dwa. Tura, która domknęła się, nie przynosząc ani znaku
  // odpowiedzi, ani rozumowania, ani narzędzia, ani błędu, nie jest
  // „zakończona" — jest zamknięta bez odpowiedzi i tak ma się nazywać, żeby
  // pusta ramka nie kazała zgadywać, czy model milczał, czy widok czegoś nie
  // narysował.
  if (wpis.bledy.length > 0) wpis.stan = 'bledny';
  else wpis.stan = czyPustaTura(wpis) ? 'pusty' : 'zakonczony';
}

/**
 * Odtwarza we wpisie jeden blok zapisany w bazie (tabela `blok_wiadomosci`
 * z pliku `migracja_096_tresc_rozmowy.sql`), odczytany z pola
 * `metadata.blocks` wiadomości kontraktu.
 *
 * Idzie przez ten sam rozdzielacz rodzajów, co fragment żywego strumienia:
 * historia po restarcie ma składać wpis tą samą logiką, którą składała go tura
 * na żywo, a drugi czytnik ładunków byłby drugą prawdą o rozdziale.
 *
 * Blok rodzaju tekstowego przechodzi bez śladu: rdzeń tekstu w blokach nie
 * zapisuje (treść mieszka w `content` wiadomości), więc taki blok mógłby być
 * wyłącznie powtórką — doliczenie go zdublowałoby treść wpisu.
 */
export function odtworzZapisanyBlok(
  wpis: WpisRozmowy,
  rodzaj: string,
  tresc: string,
  dane: unknown,
): void {
  if (rodzaj === RodzajFragmentu.Tekst) return;
  rozdzielRodzaj(wpis, {
    windowId: '',
    messageId: wpis.idWiadomosci,
    kind: rodzaj as RodzajFragmentu,
    text: tresc.length > 0 ? tresc : undefined,
    data: dane,
  });
}

/** Kieruje fragment do właściwej warstwy wpisu według jego rodzaju. */
function rozdzielRodzaj(wpis: WpisRozmowy, fragment: StreamChunkEvent): void {
  const tresc = fragment.text ?? '';
  switch (fragment.kind as RodzajFragmentu) {
    case RodzajFragmentu.Tekst:
      wpis.tresc += tresc;
      return;
    // Zastąpienie, nie doklejenie — to cały powód istnienia tego rodzaju.
    // Treść złożona z fragmentów jest przybliżeniem: kanał może fragment
    // powtórzyć po rotacji konta, a tura zapasowa nadaje własne. Wersja
    // ostateczna jest tą samą treścią, którą rdzeń zapisuje jako `content`
    // wiadomości — po podmianie wpis na ekranie i wiersz w bazie mówią to samo.
    case RodzajFragmentu.WersjaOstateczna:
      wpis.tresc = tresc;
      return;
    case RodzajFragmentu.Rozumowanie:
      wpis.rozumowanie += tresc;
      return;
    case RodzajFragmentu.Prowenancja:
      wpis.prowenancja = odczytajProwenancje(fragment.data) ?? wpis.prowenancja;
      return;
    case RodzajFragmentu.Konto:
      wpis.konto = odczytajKonto(fragment.data) ?? wpis.konto;
      return;
    case RodzajFragmentu.WywolanieNarzedzia:
      dopiszWywolanie(wpis, fragment.data);
      return;
    case RodzajFragmentu.WynikNarzedzia:
      dopiszWynikNarzedzia(wpis, fragment.data);
      return;
    case RodzajFragmentu.Blad:
      dopiszBlad(wpis, tresc, fragment.data);
      return;
    default:
      dopiszNierozpoznany(wpis, fragment.kind, tresc, fragment.data);
  }
}

/** Nowe wywołanie narzędzia. */
function dopiszWywolanie(wpis: WpisRozmowy, dane: unknown): void {
  const zrodlo = obiekt(dane);
  const wywolanie: WywolanieNarzedzia = {
    id: tekst(zrodlo, 'id'),
    nazwa: tekst(zrodlo, 'name') || 'narzędzie',
    wejscie: zapisCzytelny(zrodlo?.['input']),
    wynik: '',
    bledne: false,
  };
  wpis.narzedzia.push(wywolanie);
}

/** Wynik narzędzia dopisany do wywołania o tym samym identyfikatorze. */
function dopiszWynikNarzedzia(wpis: WpisRozmowy, dane: unknown): void {
  const zrodlo = obiekt(dane);
  const idWywolania = tekst(zrodlo, 'toolUseId');
  const wynik = zapisCzytelny(zrodlo?.['content']);
  const bledne = prawda(zrodlo, 'isError');
  const wywolanie = wpis.narzedzia.find((pozycja) => pozycja.id === idWywolania);
  if (wywolanie === undefined) {
    wpis.narzedzia.push({ id: idWywolania, nazwa: 'narzędzie', wejscie: '', wynik, bledne });
    return;
  }
  wywolanie.wynik = wynik;
  wywolanie.bledne = bledne;
}

/** Błąd kanału. Kończy turę i nic ponadto — okno i sesja zostają czynne. */
function dopiszBlad(wpis: WpisRozmowy, tresc: string, dane: unknown): void {
  const zrodlo = obiekt(dane);
  wpis.bledy.push({
    kod: tekst(zrodlo, 'code'),
    tresc: tresc.length > 0 ? tresc : tekst(zrodlo, 'message'),
    ponawialny: prawda(zrodlo, 'retryable'),
  });
}

/** Fragment rodzaju spoza katalogu — pokazany, nie porzucony. */
function dopiszNierozpoznany(
  wpis: WpisRozmowy,
  rodzaj: string,
  tresc: string,
  dane: unknown,
): void {
  const zapis = tresc.length > 0 ? tresc : zapisCzytelny(dane);
  wpis.tresc += `\n[${rodzaj}] ${zapis}\n`;
}
