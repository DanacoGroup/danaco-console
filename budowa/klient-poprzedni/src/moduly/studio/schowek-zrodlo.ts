import {
  ClipboardEntryKind,
  Command,
  type ClipboardDeleteResponse,
  type ClipboardEntry,
  type ClipboardListResponse,
  type ClipboardPinResponse,
  type ClipboardPushResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLiczba, czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Schowek Operatora — cztery komendy rodziny `clipboard.*`, wołane, nie pisane
 * drugi raz.
 *
 * ── Dlaczego rdzeń, a nie schowek przeglądarki ──────────────────────────────
 * Schowek przeglądarki żyje tyle, co karta, i nie ma historii: `navigator
 * .clipboard` oddaje jedną, ostatnią treść. Wymaganie Właściciela mówi
 * o wykazie wpisów sprzed kilku ruchów i o wpisach przypiętych na stałe, więc
 * historia musi leżeć w rdzeniu. Rodzina `clipboard.*` robi dokładnie to:
 * `push` dopisuje (powtórzenie identycznej treści NIE mnoży wpisów, tylko
 * podnosi zastany na czoło i oddaje `alreadyPresent`), `list` oddaje historię
 * wraz z przypiętymi, `pin` przypina (przypięty nie wygasa wraz z retencją),
 * `delete` usuwa wpis albo całą historię nieprzypiętą.
 *
 * ── Model sięga tą samą drogą ───────────────────────────────────────────────
 * Komendy schowka nie mają pola autora, więc wpis odłożony przez model i wpis
 * Operatora są w rdzeniu nierozróżnialne. Odkładając fragment za modelem,
 * okno zapisuje więc `sourceWindowId` — jedyne pole pochodzenia, które wpis
 * niesie. Rozróżnienia autora wpisu schowka w kontrakcie NIE MA i jest ono
 * wypisane w sprawozdaniu jako pozycja do dobudowy, a nie udawane tutaj
 * przedrostkiem w treści.
 *
 * Źródło nie ma stanu i nie buduje elementu.
 */
export interface SchowekZrodlo {
  /** Historia schowka; `tylkoPrzypiete` zawęża do wpisów przypiętych. */
  schowekWykaz(
    fraza: string,
    tylkoPrzypiete: boolean,
    granica: number,
  ): Promise<Wynik<ClipboardListResponse>>;
  /** Dopisuje treść do historii schowka. */
  schowekOdloz(
    tresc: string,
    rodzaj: ClipboardEntryKind,
    idOkna: string,
  ): Promise<Wynik<ClipboardPushResponse>>;
  /** Przypina wpis albo zdejmuje przypięcie. */
  schowekPrzypnij(idWpisu: string, przypiety: boolean): Promise<Wynik<ClipboardPinResponse>>;
  /** Usuwa wskazany wpis; puste wskazanie kasuje historię nieprzypiętą. */
  schowekUsun(idWpisu: string): Promise<Wynik<ClipboardDeleteResponse>>;
}

export function utworzSchowekZrodlo(kanal: Kanal): SchowekZrodlo {
  return {
    async schowekWykaz(fraza, tylkoPrzypiete, granica) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ClipboardList, {
          ...(fraza === '' ? {} : { query: fraza }),
          ...(tylkoPrzypiete ? { pinnedOnly: true } : {}),
          limit: granica,
        }),
        Command.ClipboardList,
        (tresc) => czyTablica(tresc.entries),
      );
    },

    async schowekOdloz(tresc, rodzaj, idOkna) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ClipboardPush, {
          content: tresc,
          kind: rodzaj,
          ...(idOkna === '' ? {} : { sourceWindowId: idOkna }),
        }),
        Command.ClipboardPush,
        (odpowiedz) => czyObiekt(odpowiedz.entry),
      );
    },

    async schowekPrzypnij(idWpisu, przypiety) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ClipboardPin, { entryId: idWpisu, pinned: przypiety }),
        Command.ClipboardPin,
        (odpowiedz) => czyObiekt(odpowiedz.entry),
      );
    },

    async schowekUsun(idWpisu) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ClipboardDelete, {
          ...(idWpisu === '' ? {} : { entryId: idWpisu }),
        }),
        Command.ClipboardDelete,
        (odpowiedz) => czyLiczba(odpowiedz.deleted),
      );
    },
  };
}

/**
 * Zdanie o wpisie schowka wraz z jego rozmiarem i czasem.
 *
 * Podgląd bierze się z pola `preview`, gdy rdzeń je oddał, a z treści tylko
 * wtedy, gdy wpis jest tekstem — treść wpisu obrazowego to bajty base64
 * i pokazanie ich jako podglądu byłoby pokazaniem szumu za treść.
 */
export function schowekOpiszWpis(wpis: ClipboardEntry): string {
  const czas = new Date(wpis.createdAt).toLocaleString('pl-PL');
  const przypiety = wpis.pinned ? 'przypięty' : 'w historii';
  return `${wpis.kind} · ${wpis.sizeBytes} bajtów · ${przypiety} · ${czas}`;
}

/** Podgląd treści wpisu do wykazu; wpis nietekstowy nie udaje tekstu. */
export function schowekPodglad(wpis: ClipboardEntry): string {
  if (wpis.preview !== undefined && wpis.preview !== '') return wpis.preview;
  if (wpis.kind === ClipboardEntryKind.Text) return wpis.content;
  return `(${wpis.kind} — bajty, nie tekst; wklejenie do dokumentu wymaga wstawienia obrazu)`;
}
