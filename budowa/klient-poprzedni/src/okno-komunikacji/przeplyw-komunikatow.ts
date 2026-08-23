import {
  Command,
  EventType,
  MessageRole,
  type ErrorInfo,
} from '../../../shared/contract';
import type { Transport } from '../polaczenie/gniazdo';
import type { Kanal } from '../protokol/kanal';
import type { PostepUzgodnienia, Uzgodnienie } from '../protokol/uzgodnienie';
import { utworzDyktowanie } from './dyktowanie/dyktowanie';
import { nazwaEtapu } from './etykiety-okna';
import type { OknoKomunikacji } from './okno';
import type { OpisOkna } from './opis-okna';
import { nazwaNadawcy, wpisOperatora, wpisSystemowy } from './wpis';
import { personaRoli, pokazFragment } from './wpisy-strumienia';

/** Elementy, które przepływ komunikatów spina w jedną całość. */
export interface CzesciPrzeplywu {
  okno: OknoKomunikacji;
  kanal: Kanal;
  transport: Transport;
  uzgodnienie: Uzgodnienie;
  opis: OpisOkna;
}

/**
 * Przepływ komunikatów okna: transport i kontrakt po jednej stronie, widok
 * po drugiej.
 *
 * Treść wpisana przed otwarciem okna nie jest odrzucana — czeka i zostaje
 * wysłana po uzgodnieniu (brak blokad w interfejsie). Niepowodzenie
 * pojedynczej komendy zostaje pokazane operatorowi i nie kończy pracy okna.
 */
export function polaczPrzeplyw(czesci: CzesciPrzeplywu): void {
  const { okno, kanal, transport, uzgodnienie, opis } = czesci;

  // Dyktowanie powstaje tutaj, bo tutaj jest kanał: warstwa pyta rdzeń
  // o dostępność silnika mowy i wysyła nagranie, a `okno.ts` kanału nie zna.
  // Przepływ zna oba końce, więc składa warstwę i oddaje ją oknu; pasek
  // polecenia bierze ją stamtąd i rysuje mikrofon — albo go nie rysuje, gdy
  // warstwa mówi, że nie ma czym nagrać (krótszy pasek zamiast wyszarzonego
  // przycisku).
  okno.podlaczDyktowanie(utworzDyktowanie(kanal));
  const oczekujace: string[] = [];
  const zestrumieniowane = new Set<string>();
  const persona = opis.kanalModelu;

  /** Identyfikator okna nadany przez rdzeń; pusty do chwili uzgodnienia. */
  function idOkna(): string {
    return uzgodnienie.okno()?.id ?? '';
  }

  /**
   * Czy okno prowadzi turę. Ustawia się przy nadaniu wiadomości, gaśnie ze
   * zdarzeniem domykającym strumień (koperta z `done`).
   *
   * Znacznik jest przybliżeniem: rozstrzyga wyłącznie o tym, czy poprzedzić
   * wysłanie zatrzymaniem. Prawdę o stanie okna zna rdzeń i to on odmawia —
   * znacznik nieaktualny kosztuje jedno zbędne `message.stop`, które zawsze
   * odpowiada, a nie utratę wiadomości.
   */
  let turaWBiegu = false;

  function nadajTresc(tresc: string): void {
    kanal.wyslij(
      Command.MessageSend,
      { windowId: idOkna(), content: tresc, stream: true },
      (wynik) => {
        if (wynik.udany) {
          turaWBiegu = true;
          return;
        }
        turaWBiegu = false;
        okno.dopisz(wpisSystemowy(`Wiadomość niewysłana — ${opisBledu(wynik.blad)}`));
      },
    );
  }

  /**
   * Wysyła wiadomość, a gdy okno prowadzi turę — poprzedza ją zatrzymaniem.
   *
   * Przerwanie jest jawne: `message.send` skierowany do okna, które odpowiada,
   * odmawia kodem `conflict` zamiast skasować odpowiedź w połowie zdania.
   * Klient wysyła więc parę — najpierw `message.stop`, po jego odpowiedzi
   * `message.send`.
   *
   * Kolejność jest sekwencyjna, nie równoległa: obie komendy nadane naraz
   * dotarłyby w kolejności niegwarantowanej i wysłanie mogłoby wyprzedzić
   * zatrzymanie, czyli trafić na okno nadal zajęte i odmówić.
   *
   * Nieudane zatrzymanie nie wstrzymuje wysłania. Jeżeli tura zdążyła
   * tymczasem dobiec końca sama, wysłanie przejdzie; jeżeli nie — odmowa
   * przyjdzie z rdzenia.
   */
  function wyslijTresc(tresc: string): void {
    if (!turaWBiegu) {
      nadajTresc(tresc);
      return;
    }
    kanal.wyslij(Command.MessageStop, { windowId: idOkna() }, () => {
      turaWBiegu = false;
      nadajTresc(tresc);
    });
  }

  // Stan transportu trafia do nagłówka, a nawiązanie połączenia rozpoczyna
  // uzgodnienie. Ponowne połączenie nie zakłada drugiego okna, jeżeli rdzeń
  // otworzył je już wcześniej.
  transport.naStan((stan) => {
    okno.pokazStan(stan, transport.oczekujace());
    if (stan === 'polaczony' && uzgodnienie.okno() === null) uzgodnienie.rozpocznij();
  });

  uzgodnienie.naPostep((postep) => okno.dopisz(wpisSystemowy(opisPostepu(postep))));

  uzgodnienie.naOtwarcieOkna(() => {
    for (const tresc of oczekujace.splice(0)) wyslijTresc(tresc);
  });

  okno.naWpisanie((tresc) => {
    okno.dopisz(wpisOperatora(tresc, nazwaNadawcy(MessageRole.User)));
    if (idOkna().length === 0) {
      oczekujace.push(tresc);
      return;
    }
    wyslijTresc(tresc);
  });

  kanal.naZdarzenie(EventType.StreamChunk, (fragment, koperta) => {
    if (fragment.windowId !== idOkna()) return;
    zestrumieniowane.add(fragment.messageId);
    // Zdarzenie domykające gasi znacznik tury — jedno na turę, także po błędzie
    // i po przerwaniu (server/internal/core/strumien_odpowiedzi.go).
    if (koperta.done === true) turaWBiegu = false;
    pokazFragment(okno, persona, fragment);
  });

  kanal.naZdarzenie(EventType.MessageChanged, ({ message }) => {
    if (message.windowId !== idOkna()) return;
    if (message.role === MessageRole.User || zestrumieniowane.has(message.id)) return;
    okno.dopisz({
      nadawca: message.role,
      persona: personaRoli(message.role, persona),
      tresc: message.content,
      znacznikCzasu: message.createdAt,
    });
  });

  // Zmiana okna niesie także jego moduł. `workspace.enter` przestawia moduł
  // tego okna zamiast zakładać drugie, więc okno rekonfiguruje pasek narzędzi,
  // panel akcji i kontekst, a wątek zostaje nietknięty.
  kanal.naZdarzenie(EventType.WindowChanged, ({ change, window }) => {
    if (window.id !== idOkna()) return;
    okno.ustawModul(window.moduleId);
    okno.dopisz(wpisSystemowy(`Okno komunikacji — ${change}, stan: ${window.status}`));
  });
}

/** Komunikat o etapie uzgodnienia z rdzeniem. */
function opisPostepu(postep: PostepUzgodnienia): string {
  const etap = nazwaEtapu(postep.etap);
  return postep.udany ? etap : `${etap} — ${opisBledu(postep.blad)}`;
}

/** Opis błędu kontraktu dla operatora. */
function opisBledu(blad: ErrorInfo | undefined): string {
  if (blad === undefined) return 'rdzeń nie podał przyczyny';
  return `${blad.message} (${blad.code})`;
}
