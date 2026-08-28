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
 * Doradca eksperta prowadzi konsultację modelem silniejszym niż model
 * bazowy, przez okno komunikacji zakładane i zamykane na czas jednego
 * zapytania.
 */
const KOD_MODULU = 'agents';

/** Zlecenie konsultacji: kogo dotyczy, kogo pytamy o radę i jakie pytanie zadajemy temu doradcy modelu. */
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

/** Rada doradcy wraz z pochodzeniem; treść rady nigdy nie chodzi bez tożsamości doradcy i treści pytania. */
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
  /** Przeprowadza jedną konsultację od założenia okna do jego zamknięcia; niepowodzenie wraca odmową. */
  zapytaj(zlecenie: ZlecenieRady): Promise<Wynik<RadaDoradcy>>;
}

export function utworzZrodloDoradcy(kanal: Kanal): ZrodloDoradcy {
  /** Czeka na pierwszą domkniętą wiadomość modelu — jedyną postać, której znaczenie jest pewne. */
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
          // Sesję zna kanał, nie widok modułu; panel nie przepisuje wartości warstwy protokołu.
          sessionId: kanal.sesja().id(),
          moduleId: KOD_MODULU,
          modelChannelId: zlecenie.kanalDoradcy,
          workingDirs: [],
          // Konsultacja jest rozmową, nie robotą na plikach: tryb planistyczny znaczy tu poradę.
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

      // Nasłuch zakładamy przed wysłaniem pytania, bo rdzeń bywa szybszy od obietnicy wysyłki.
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

/** Kanały doradcze, czyli te kanały modelu widoczne w wykazie, które nie są kanałem bazowym danego eksperta. */
export function kanalyDoradcze<T extends { id: string }>(
  kanaly: readonly T[],
  kanalBazowy: string,
): T[] {
  return kanaly.filter((pozycja) => pozycja.id !== kanalBazowy);
}
