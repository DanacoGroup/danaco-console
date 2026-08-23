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

/**
 * Odtworzenie historii okna z rdzenia.
 *
 * Komenda `message.list` daje wiadomości leżące w bazie, więc odświeżenie okna
 * albo powrót do sesji pokazuje wątek, a nie pustkę. Wpisy odtworzone są
 * domknięte: pochodzą z zapisu, nie ze strumienia, więc nie mają tury w biegu
 * ani fragmentów do doliczenia.
 *
 * Wywołujący dostaje dwie drogi — `przyjmij` (na sukces, także z pustą tablicą)
 * i `zglosNiepowodzenie` (na odmowę, z gotowym zdaniem złożonym z powodu
 * koperty) — i ma obsłużyć obie. Odmowa rdzenia i historia naprawdę pusta muszą
 * wyglądać na ekranie inaczej.
 */
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

/**
 * Przekład wiadomości kontraktu na wpis rozmowy.
 *
 * Treść wpisu bierze się z pola `content`; wiadomość przerwana albo błędna
 * wraca z treścią, którą zdążyła zebrać — pokazanie jej jest uczciwsze niż
 * ukrycie tury, która się odbyła.
 *
 * Stan wpisu bierze się z pola `status` wiadomości, nie z samego faktu zapisu:
 * tura przerwana i tura zamknięta błędem mają po odświeżeniu okna wyglądać tak
 * samo jak w chwili, w której się odbyły.
 */
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

  // Rdzeń oddaje nietekstowe fragmenty tury — tok rozumowania, narzędzia
  // z wynikami, prowenancję, konto, błędy — w `metadata.blocks`; do wpisu
  // wchodzą tą samą logiką, którą składa je strumień żywy
  // (odtworzZapisanyBlok → rozdzielacz rodzajów).
  for (const blok of blokiZapisane(wiadomosc.metadata)) {
    odtworzZapisanyBlok(wpis, blok.rodzaj, blok.tresc, blok.dane);
  }
  // Tura bez tekstu przestaje być „zamknięta bez odpowiedzi", jeżeli bloki
  // przyniosły rozumowanie albo narzędzie — pustkę mierzy się po doliczeniu
  // wszystkich warstw wpisu, tak jak przy turze żywej (czyPustaTura).
  if (wpis.stan === 'pusty' && !czyPustaTura(wpis)) wpis.stan = 'zakonczony';
  return wpis;
}

/** Jeden blok wiadomości odczytany z metadanych. */
interface BlokZapisany {
  rodzaj: string;
  tresc: string;
  dane: unknown;
}

/**
 * Odczyt wykazu `blocks` z metadanych wiadomości. Odczyt jest tolerancyjny
 * (wzorem `odczyt-fragmentu.ts`): brak obszaru, obszar cudzego kształtu albo
 * pozycja bez rodzaju dają mniej bloków — historia ma się wyświetlić, a nie
 * zniknąć od nieczytelnej pozycji.
 */
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

/**
 * Stan wpisu odtworzonego z zapisu — wprost ze stanu wiadomości kontraktu.
 *
 * Wypowiedź Operatora nie ma stanu tury: jest tym, co napisał, i nie może
 * wrócić z zapisu jako „zamknięta bez odpowiedzi", nawet gdyby zapis był pusty.
 */
function stanZapisany(wiadomosc: Message, odOperatora: boolean): StanWpisu {
  if (odOperatora) return 'zakonczony';
  if (wiadomosc.status === MessageStatus.Error) return 'bledny';
  if (wiadomosc.status === MessageStatus.Stopped) return 'przerwany';
  if ((wiadomosc.content ?? '').trim().length === 0) return 'pusty';
  return 'zakonczony';
}
