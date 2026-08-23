import {
  Command,
  type ClipboardDeleteResponse,
  type ClipboardEntryKind,
  type ClipboardListResponse,
  type ClipboardPinResponse,
  type ClipboardPushRequest,
  type ClipboardPushResponse,
  type LauncherHotkeyGetResponse,
  type LauncherHotkeySetResponse,
  type SnippetDeleteResponse,
  type SnippetListResponse,
  type SnippetSetRequest,
  type SnippetSetResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLiczba, czyLogiczna, czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Dziewięć komend Command & Tools Hub: historia schowka, słownik skrótów
 * tekstowych i skrót globalny wywoływacza poleceń.
 *
 * ── Czyj jest schowek ───────────────────────────────────────────────────────
 * Schowek należy do maszyny Operatora i rdzeń go NIE czyta. Podział ról jest
 * taki: Operator kopiuje u siebie, okno oddaje skopiowaną treść rdzeniowi
 * (`clipboard.push`), a rdzeń daje jej trwałość — historia przestaje ginąć
 * razem z kartą i jest ta sama na każdej maszynie tego samego Operatora.
 * Wklejenie jest ruchem powrotnym i wykonuje je okno, u siebie.
 *
 * ── Skrót globalny ──────────────────────────────────────────────────────────
 * Nastawę trzyma rdzeń, przechwycenie klawiszy należy do powłoki programu
 * okiennego. Odpowiedź `launcher.hotkey.*` mówi wprost, czy rejestracji ma kto
 * dokonać — okno powtarza to zdanie zamiast obiecywać skrót, który nikogo nie
 * obudzi.
 */

/** Wpis oddawany rdzeniowi po skopiowaniu treści u Operatora. */
export interface WpisSchowkaDoZapisu {
  tresc: string;
  rodzaj?: ClipboardEntryKind;
  idOknaZrodlowego?: string;
  wrazliwy?: boolean;
}

/** Skrót tekstowy zapisywany z okna; puste `id` zakłada nowy. */
export interface SkrotDoZapisu {
  id: string;
  skrot: string;
  tresc: string;
  opis: string;
  polaSzablonu: readonly string[];
  czynny: boolean;
}

export interface ZrodloSchowka {
  /** `clipboard.list` — historia, przypięte na czele. */
  historia(fraza: string, tylkoPrzypiete: boolean): Promise<Wynik<ClipboardListResponse>>;
  /** `clipboard.push` — treść skopiowana u Operatora oddana rdzeniowi. */
  dopisz(wpis: WpisSchowkaDoZapisu): Promise<Wynik<ClipboardPushResponse>>;
  /** `clipboard.pin` — przypięcie chroniące wpis przed czyszczeniem. */
  przypnij(id: string, przypiety: boolean): Promise<Wynik<ClipboardPinResponse>>;
  /** `clipboard.delete` — jeden wpis albo cała historia nieprzypięta. */
  usun(id: string, zPrzypietymi: boolean): Promise<Wynik<ClipboardDeleteResponse>>;
  /** `snippet.list` — słownik skrótów rozwijanych w dłuższą treść. */
  skroty(fraza: string): Promise<Wynik<SnippetListResponse>>;
  /** `snippet.set` — założenie albo zmiana skrótu. */
  zapiszSkrot(skrot: SkrotDoZapisu): Promise<Wynik<SnippetSetResponse>>;
  /** `snippet.delete` — usunięcie skrótu ze słownika. */
  usunSkrot(id: string): Promise<Wynik<SnippetDeleteResponse>>;
  /** `launcher.hotkey.get` — skrót globalny wraz ze stanem wykonalności. */
  skrotGlobalny(): Promise<Wynik<LauncherHotkeyGetResponse>>;
  /** `launcher.hotkey.set` — zapis skrótu; pusty zdejmuje rejestrację. */
  zapiszSkrotGlobalny(skrot: string): Promise<Wynik<LauncherHotkeySetResponse>>;
}

/** Górna granica historii — okno pokazuje ostatnie kopie, nie całe archiwum. */
const GRANICA_HISTORII = 100;

export function utworzZrodloSchowka(kanal: Kanal): ZrodloSchowka {
  return {
    async historia(fraza, tylkoPrzypiete) {
      const zadanie: Record<string, unknown> = {
        limit: GRANICA_HISTORII,
        pinnedOnly: tylkoPrzypiete,
      };
      if (fraza.trim() !== '') zadanie['query'] = fraza.trim();
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ClipboardList, zadanie),
        Command.ClipboardList,
        (tresc) => czyTablica(tresc.entries) && czyLiczba(tresc.total),
      );
    },

    async dopisz(wpis) {
      const zadanie: ClipboardPushRequest = { content: wpis.tresc };
      if (wpis.rodzaj !== undefined) zadanie.kind = wpis.rodzaj;
      if (wpis.idOknaZrodlowego !== undefined && wpis.idOknaZrodlowego !== '') {
        zadanie.sourceWindowId = wpis.idOknaZrodlowego;
      }
      if (wpis.wrazliwy !== undefined) zadanie.sensitive = wpis.wrazliwy;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ClipboardPush, zadanie),
        Command.ClipboardPush,
        (tresc) => czyObiekt(tresc.entry) && czyLogiczna(tresc.alreadyPresent),
      );
    },

    async przypnij(id, przypiety) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ClipboardPin, { entryId: id, pinned: przypiety }),
        Command.ClipboardPin,
        (tresc) => czyObiekt(tresc.entry),
      );
    },

    async usun(id, zPrzypietymi) {
      const zadanie: Record<string, unknown> = { includePinned: zPrzypietymi };
      if (id !== '') zadanie['entryId'] = id;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ClipboardDelete, zadanie),
        Command.ClipboardDelete,
        (tresc) => czyLiczba(tresc.deleted),
      );
    },

    async skroty(fraza) {
      const zadanie: Record<string, unknown> = {};
      if (fraza.trim() !== '') zadanie['query'] = fraza.trim();
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SnippetList, zadanie),
        Command.SnippetList,
        (tresc) => czyTablica(tresc.snippets) && czyLiczba(tresc.total),
      );
    },

    async zapiszSkrot(skrot) {
      const zadanie: SnippetSetRequest = {
        shortcut: skrot.skrot,
        content: skrot.tresc,
        enabled: skrot.czynny,
      };
      if (skrot.id !== '') zadanie.snippetId = skrot.id;
      if (skrot.opis.trim() !== '') zadanie.description = skrot.opis.trim();
      if (skrot.polaSzablonu.length > 0) zadanie.variables = [...skrot.polaSzablonu];
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SnippetSet, zadanie),
        Command.SnippetSet,
        (tresc) => czyObiekt(tresc.snippet),
      );
    },

    async usunSkrot(id) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SnippetDelete, { snippetId: id }),
        Command.SnippetDelete,
        (tresc) => czyLogiczna(tresc.deleted),
      );
    },

    async skrotGlobalny() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.LauncherHotkeyGet, {}),
        Command.LauncherHotkeyGet,
        (tresc) => czyLogiczna(tresc.supported) && czyLogiczna(tresc.registered),
      );
    },

    async zapiszSkrotGlobalny(skrot) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.LauncherHotkeySet, { hotkey: skrot }),
        Command.LauncherHotkeySet,
        (tresc) => czyLogiczna(tresc.registered),
      );
    },
  };
}
