import type { StreamChunkEvent } from '../../../shared/contract';
import { odczytajKonto } from './metadane-konta';
import { obiekt, prawda, tekst, zapisCzytelny } from './odczyt-fragmentu';
import { odczytajPodsumowanie } from './podsumowanie-tury';
import { odczytajProwenancje } from './prowenancja';
import { RodzajFragmentu } from './rodzaje-fragmentow';
import { czyPustaTura, type WpisRozmowy, type WywolanieNarzedzia } from './wpis-rozmowy';

/**
 * Znaczniki toru strumienia niesione przez kopertę zdarzenia sieciowego, a nie przez sam
 * ładunek fragmentu strumienia.
 */
export interface ZnacznikiFragmentu {
  /** Numer fragmentu liczony od jedynki (pole `seq`). */
  numer: number;
  /** Prawda w ostatnim fragmencie tury (pole `done`). */
  ostatni: boolean;
}

/**
 * Złożenie tury dokłada jeden fragment strumienia odpowiedzi do wpisu rozmowy, według
 * rozdzielenia jego rodzaju.
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

  // Podsumowanie czyta się z każdego fragmentu, nie tylko z fragmentu domykającego strumień.
  wpis.podsumowanie ??= odczytajPodsumowanie(fragment.data);

  if (!znaczniki.ostatni) return;
  wpis.domkniety = true;
  // Tura bez znaku odpowiedzi, rozumowania czy narzędzia jest zamknięta bez odpowiedzi.
  if (wpis.bledy.length > 0) wpis.stan = 'bledny';
  else wpis.stan = czyPustaTura(wpis) ? 'pusty' : 'zakonczony';
}

/**
 * Odtwarza we wpisie jeden blok zapisany w bazie, przechodząc przez ten sam rozdzielacz
 * rodzajów co żywy fragment.
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

/**
 * Kieruje fragment strumienia odpowiedzi do właściwej warstwy wpisu rozmowy według jego
 * rodzaju wewnętrznego.
 */
function rozdzielRodzaj(wpis: WpisRozmowy, fragment: StreamChunkEvent): void {
  const tresc = fragment.text ?? '';
  switch (fragment.kind as RodzajFragmentu) {
    case RodzajFragmentu.Tekst:
      wpis.tresc += tresc;
      return;
    // Wersja ostateczna zastępuje treść złożoną z fragmentów, nie dokłada się do niej.
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

/**
 * Nowe wywołanie narzędzia, zgłoszone przez model w trakcie trwającej tury tego bieżącego
 * okna rozmowy Operatora.
 */
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

/**
 * Wynik narzędzia dopisywany do wcześniejszego wywołania o tym samym identyfikatorze tego
 * zgłoszenia narzędzia.
 */
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

/**
 * Błąd zgłoszony przez kanał kończy bieżącą turę i nic więcej ponadto — samo okno oraz
 * cała sesja rozmowy zostają nadal czynne.
 */
function dopiszBlad(wpis: WpisRozmowy, tresc: string, dane: unknown): void {
  const zrodlo = obiekt(dane);
  wpis.bledy.push({
    kod: tekst(zrodlo, 'code'),
    tresc: tresc.length > 0 ? tresc : tekst(zrodlo, 'message'),
    ponawialny: prawda(zrodlo, 'retryable'),
  });
}

/**
 * Fragment rodzaju spoza katalogu znanych rodzajów strumienia odpowiedzi, pokazany wprost,
 * a nie porzucony po cichu.
 */
function dopiszNierozpoznany(
  wpis: WpisRozmowy,
  rodzaj: string,
  tresc: string,
  dane: unknown,
): void {
  const zapis = tresc.length > 0 ? tresc : zapisCzytelny(dane);
  wpis.tresc += `\n[${rodzaj}] ${zapis}\n`;
}
