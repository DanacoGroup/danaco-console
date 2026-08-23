import {
  Command,
  ExecutionEnv,
  MessageRole,
  MessageStatus,
  PermissionMode,
  WindowRole,
  EventType,
  type Message,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Doradca eksperta — konsultacja modelem silniejszym niż model bazowy.
 *
 * Kontrakt nie ma rodziny `agent.advisor.*`, więc konsultacja idzie trzema
 * komendami obszarów `window.*` i `message.*` oraz jednym zdarzeniem:
 *
 *   `window.create`  — okno konsultacji na kanale doradcy, nie na kanale
 *                      eksperta. To jedyne miejsce, w którym wybór modelu
 *                      silniejszego staje się faktem po stronie rdzenia:
 *                      `modelChannelId` jest polem wymaganym żądania.
 *   `message.send`   — pytanie do doradcy, bez strumienia; okno jest
 *                      jednorazowe i nie ma czego rysować na żywo.
 *   `message.changed`— odpowiedź doradcy. `message.send` oddaje wiadomość
 *                      przyjętą (własną, rolę `user`), więc czekanie na jej
 *                      wynik dałoby echo pytania, nie radę.
 *   `window.close`   — okno konsultacji znika po odpowiedzi. Zostawione
 *                      wisiałoby w wykazie okien sesji jako okno bez widoku.
 *
 * Rada nie jest odpowiedzią eksperta. Źródło oddaje treść rady razem
 * z tożsamością doradcy i treścią zadanego pytania — kto pyta, o co pyta i kto
 * odpowiada, są trzema osobnymi polami wyniku, a nie jednym napisem, więc widok
 * nie ma z czego złożyć rady podanej jako własne zdanie eksperta.
 *
 * Źródło nie ma stanu. Okno konsultacji żyje wyłącznie w obrębie jednego
 * wywołania `zapytaj`; nic z niego nie zostaje w polu modułu.
 */

/** Kod modułu nadawany oknu konsultacji — ten sam, którym opisuje się moduł. */
const KOD_MODULU = 'agents';

/** Zlecenie konsultacji: kogo dotyczy, kogo pytamy i o co. */
export interface ZlecenieRady {
  /** Kanał modelu doradcy — model silniejszy od bazowego modelu eksperta. */
  kanalDoradcy: string;
  /** Nazwa kanału doradcy widziana przez Operatora; wchodzi do prowenancji. */
  nazwaDoradcy: string;
  /** Ekspert, którego konsultacja dotyczy. */
  idEksperta: string;
  nazwaEksperta: string;
  /** Pytanie zadane doradcy — dokładnie w tej postaci, w jakiej pójdzie. */
  pytanie: string;
}

/** Rada doradcy wraz z pochodzeniem; treść nigdy nie chodzi bez nich. */
export interface RadaDoradcy {
  /** Treść odpowiedzi doradcy. */
  tresc: string;
  /** Pytanie w postaci wysłanej do rdzenia. */
  pytanie: string;
  nazwaDoradcy: string;
  kanalDoradcy: string;
  idEksperta: string;
  nazwaEksperta: string;
  /** Okno konsultacji, w którym rada powstała — ślad dla Mission Control. */
  idOkna: string;
  /** Czas przyjęcia rady w milisekundach epoki. */
  chwila: number;
  /** Stan wiadomości oddany przez rdzeń; `stopped` i `error` też tu wracają. */
  stan: MessageStatus;
}

export interface ZrodloDoradcy {
  /**
   * Przeprowadza jedną konsultację od założenia okna do jego zamknięcia.
   * Niepowodzenie każdego z czterech kroków wraca polem `blad`.
   */
  zapytaj(zlecenie: ZlecenieRady): Promise<Wynik<RadaDoradcy>>;
}

export function utworzZrodloDoradcy(kanal: Kanal): ZrodloDoradcy {
  /**
   * Czeka na pierwszą domkniętą wiadomość modelu w oknie konsultacji.
   *
   * `stream.chunk` nie jest tu drogą: okno konsultacji zakładamy z
   * `stream: false`, a fragment strumienia niesie treść częściową. Domknięta
   * wiadomość z `message.changed` jest jedyną postacią, której znaczenie jest
   * pewne.
   *
   * Limitu czasu nie ma z tego samego powodu, co w `protokol/wywolanie.ts`:
   * kontrakt go nie przewiduje, a rozłączenie klienta nie kończy pracy rdzenia.
   * Zwrócone `odwolaj` zdejmuje subskrypcję wtedy, gdy odpowiedź już nie
   * nadejdzie — bez niego nasłuch przeżyłby nieudane wysłanie pytania.
   */
  function poczekajNaOdpowiedz(idOkna: string): {
    odpowiedz: Promise<Message>;
    odwolaj: Odsubskrybuj;
  } {
    let odsubskrybuj: Odsubskrybuj = () => undefined;
    const odpowiedz = new Promise<Message>((rozstrzygnij) => {
      odsubskrybuj = kanal.naZdarzenie(EventType.MessageChanged, (tresc) => {
        const wiadomosc = tresc.message;
        if (wiadomosc.windowId !== idOkna) return;
        if (wiadomosc.role !== MessageRole.Assistant) return;
        if (wiadomosc.status === MessageStatus.Pending) return;
        if (wiadomosc.status === MessageStatus.Streaming) return;
        odsubskrybuj();
        rozstrzygnij(wiadomosc);
      });
    });
    return { odpowiedz, odwolaj: () => odsubskrybuj() };
  }

  return {
    async zapytaj(zlecenie) {
      const okno = sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowCreate, {
          // Sesję zna kanał, nie widok modułu — ten sam wzorzec, którym idzie
          // `window.handoff` w `multitasking/zrodlo-biegu.ts`. Przeciąganie jej
          // przez panel byłoby przepisywaniem wartości, którą warstwa protokołu
          // i tak trzyma.
          sessionId: kanal.sesja().id(),
          moduleId: KOD_MODULU,
          modelChannelId: zlecenie.kanalDoradcy,
          workingDirs: [],
          // Konsultacja jest rozmową, nie robotą na plikach: zasięg rdzenia
          // i tryb planistyczny znaczą razem „poradź, niczego nie zmieniaj".
          executionEnv: ExecutionEnv.Core,
          permissionMode: PermissionMode.Plan,
          windowRole: WindowRole.Standalone,
          title: `Konsultacja doradcy — ${zlecenie.nazwaEksperta}`,
        }),
        Command.WindowCreate,
        (tresc) => czyObiekt(tresc.window),
      );
      if (!okno.udany || okno.wynik === undefined) {
        return { udany: false, ...(okno.blad === undefined ? {} : { blad: okno.blad }) };
      }
      const idOkna = okno.wynik.window.id;

      // Nasłuch zakładamy przed wysłaniem pytania. Rdzeń bywa szybszy od
      // obietnicy `message.send`: subskrypcja założona po niej przegapiłaby
      // odpowiedź i konsultacja wisiałaby na zawsze.
      const { odpowiedz, odwolaj } = poczekajNaOdpowiedz(idOkna);

      const wyslanie = sprawdzKsztalt(
        await wywolaj(kanal, Command.MessageSend, {
          windowId: idOkna,
          content: zlecenie.pytanie,
          stream: false,
        }),
        Command.MessageSend,
        (tresc) => czyObiekt(tresc.message),
      );
      if (!wyslanie.udany) {
        odwolaj();
        await zamknij(kanal, idOkna);
        return { udany: false, ...(wyslanie.blad === undefined ? {} : { blad: wyslanie.blad }) };
      }

      const rada = await odpowiedz;
      await zamknij(kanal, idOkna);

      return {
        udany: true,
        wynik: {
          tresc: rada.content,
          pytanie: zlecenie.pytanie,
          nazwaDoradcy: zlecenie.nazwaDoradcy,
          kanalDoradcy: zlecenie.kanalDoradcy,
          idEksperta: zlecenie.idEksperta,
          nazwaEksperta: zlecenie.nazwaEksperta,
          idOkna,
          chwila: rada.createdAt,
          stan: rada.status,
        },
      };
    },
  };
}

/**
 * Zamyka okno konsultacji. Odmowa zamknięcia nie unieważnia rady — Operator
 * dostał odpowiedź i to ona jest wynikiem czynności; okno zostaje wtedy
 * w wykazie sesji i Mission Control je pokaże.
 */
async function zamknij(kanal: Kanal, idOkna: string): Promise<void> {
  await wywolaj(kanal, Command.WindowClose, { windowId: idOkna });
}

/** Kanały doradcze — te, które nie są kanałem bazowym eksperta. */
export function kanalyDoradcze<T extends { id: string }>(
  kanaly: readonly T[],
  kanalBazowy: string,
): T[] {
  return kanaly.filter((pozycja) => pozycja.id !== kanalBazowy);
}
