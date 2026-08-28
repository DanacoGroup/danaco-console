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
 * Źródło grupuje cztery komendy schowka rdzenia: odczyt historii wpisów wraz
 * z przypiętymi, dopisanie treści, przypięcie oraz usunięcie wpisu; nie ma stanu.
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
 * Buduje zdanie opisowe wpisu schowka z jego rodzajem, rozmiarem w bajtach,
 * stanem przypięcia oraz czasem utworzenia sformatowanym w zapisie lokalnym.
 */
export function schowekOpiszWpis(wpis: ClipboardEntry): string {
  const czas = new Date(wpis.createdAt).toLocaleString('pl-PL');
  const przypiety = wpis.pinned ? 'przypięty' : 'w historii';
  return `${wpis.kind} · ${wpis.sizeBytes} bajtów · ${przypiety} · ${czas}`;
}

/**
 * Zwraca podgląd treści wpisu do wykazu; korzysta z pola przygotowanego przez
 * rdzeń, a dla wpisu nietekstowego oddaje opis rodzaju zamiast surowych bajtów.
 */
export function schowekPodglad(wpis: ClipboardEntry): string {
  if (wpis.preview !== undefined && wpis.preview !== '') return wpis.preview;
  if (wpis.kind === ClipboardEntryKind.Text) return wpis.content;
  return `(${wpis.kind} — bajty, nie tekst; wklejenie do dokumentu wymaga wstawienia obrazu)`;
}
