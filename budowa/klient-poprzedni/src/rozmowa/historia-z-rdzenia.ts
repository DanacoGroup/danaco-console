import {
  Command,
  MessageRole,
  MessageStatus,
  type Message,
  type WindowRole,
} from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Kanal } from '../protokol/kanal';
import { rozpoznajNadawce } from './nadawca';
import { obiekt, tekst } from './odczyt-fragmentu';
import { czyPustaTura, nowyWpis, type StanWpisu, type WpisRozmowy } from './wpis-rozmowy';
import { odtworzZapisanyBlok } from './zlozenie-tury';

/** Funkcja wczytuje historię okna z rdzenia komendą message.list i przekłada otrzymane wiadomości na wpisy rozmowy dla wywołującego. */
export function wczytajHistorie(
  kanal: Kanal,
  idOkna: string,
  rolaOkna: WindowRole | null,
  persona: string,
  personaOperatora: string,
  przyjmij: (wpisy: WpisRozmowy[]) => void,
  zglosNiepowodzenie: (powod: string) => void,
): void {
  kanal.wyslij(Command.MessageList, { windowId: idOkna }, (wynik) => {
    if (!wynik.udany) {
      zglosNiepowodzenie(opisOdmowyBledu('Historia okna', wynik.blad));
      return;
    }
    const wiadomosci = wynik.wynik?.messages ?? [];
    przyjmij(wiadomosci.map((w) => wpisZWiadomosci(w, rolaOkna, persona, personaOperatora)));
  });
}

/** Funkcja przekłada wiadomość kontraktu na wpis rozmowy, ustalając jego treść, nadawcę i stan na podstawie pól wiadomości. */
function wpisZWiadomosci(
  wiadomosc: Message,
  rolaOkna: WindowRole | null,
  persona: string,
  personaOperatora: string,
): WpisRozmowy {
  const odOperatora = wiadomosc.role === MessageRole.User;
  const nadawca = odOperatora
    ? rozpoznajNadawce(MessageRole.User, rolaOkna)
    : rozpoznajNadawce(MessageRole.Assistant, rolaOkna);

  const wpis = nowyWpis(
    wiadomosc.id,
    nadawca,
    odOperatora ? personaOperatora : persona,
    stanZapisany(wiadomosc, odOperatora),
  );
  wpis.idWiadomosci = wiadomosc.id;
  wpis.tresc = wiadomosc.content ?? '';
  wpis.domkniety = true;
  wpis.znacznikCzasu = wiadomosc.createdAt;

  // Rdzeń oddaje nietekstowe fragmenty tury w polu metadata.blocks tą samą logiką co strumień żywy.
  for (const blok of blokiZapisane(wiadomosc.metadata)) {
    odtworzZapisanyBlok(wpis, blok.rodzaj, blok.tresc, blok.dane);
  }
  // Tura bez tekstu przestaje być pusta, jeżeli bloki przyniosły rozumowanie albo narzędzie.
  if (wpis.stan === 'pusty' && !czyPustaTura(wpis)) wpis.stan = 'zakonczony';
  return wpis;
}

/** Interfejs opisuje jeden blok wiadomości odczytany z metadanych zapisanej wiadomości kontraktu rdzenia. */
interface BlokZapisany {
  rodzaj: string;
  tresc: string;
  dane: unknown;
}

/** Funkcja odczytuje wykaz bloków z metadanych wiadomości w sposób tolerancyjny na brak albo niewłaściwy kształt danych. */
function blokiZapisane(metadata: unknown): BlokZapisany[] {
  const obszar = obiekt(metadata);
  const surowe = obszar?.['blocks'];
  if (!Array.isArray(surowe)) return [];
  const bloki: BlokZapisany[] = [];
  for (const pozycja of surowe) {
    const zrodlo = obiekt(pozycja);
    const rodzaj = tekst(zrodlo, 'kind');
    if (rodzaj.length === 0) continue;
    bloki.push({ rodzaj, tresc: tekst(zrodlo, 'text'), dane: zrodlo?.['data'] });
  }
  return bloki;
}

/** Funkcja ustala stan wpisu odtworzonego z zapisu wprost na podstawie stanu zapisanego w wiadomości kontraktu. */
function stanZapisany(wiadomosc: Message, odOperatora: boolean): StanWpisu {
  if (odOperatora) return 'zakonczony';
  if (wiadomosc.status === MessageStatus.Error) return 'bledny';
  if (wiadomosc.status === MessageStatus.Stopped) return 'przerwany';
  if ((wiadomosc.content ?? '').trim().length === 0) return 'pusty';
  return 'zakonczony';
}
