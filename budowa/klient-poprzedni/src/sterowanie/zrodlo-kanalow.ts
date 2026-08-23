import {
  Command,
  type Channel,
  type ChannelCheckResponse,
  type ChannelCredentialStatusResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyObiekt, czyTekst, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { przenies } from '../protokol/wynik-czastkowy';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Strona zapisująca rejestru kanałów modelu wraz z dwiema czynnościami
 * pytającymi: `channel.add`, `channel.update`, `channel.remove`,
 * `channel.check`, `channel.credential.status`.
 *
 * Zakładanie stoi tu razem z usuwaniem, bo świeża baza ma dokładnie jeden
 * kanał, a okno komunikacji kanału wymaga. Panel z samym usuwaniem pozwoliłby
 * nieodwracalnie skasować jedyny kanał produktu i nie zostawiłby drogi powrotu.
 *
 * Wykazu kanałów to źródło nie pobiera. Rejestr kanałów
 * (`sterowanie/rejestr-kanalow.ts`) jest czytającą pamięcią podręczną nad
 * `channel.list`, wspólną całemu klientowi, i drugiej listy kanałów nie ma.
 * Po udanym zapisie panel woła `rejestr.odswiez()`.
 *
 * Każda czynność oddaje `Wynik<T>`, nigdy samej treści: odmowa musi dojechać
 * do widoku z kodem i wiadomością rdzenia — usunięcie odrzucone przez rdzeń nie
 * może wyglądać jak usunięcie wykonane.
 */

/** Opis kanału zakładanego w rejestrze; `kind` jest wymagany tylko tutaj. */
export interface OpisZalozeniaKanalu {
  name: string;
  kind: string;
  model?: string;
  enabled?: boolean;
  config?: unknown;
}

/**
 * Opis zmiany wiersza rejestru.
 *
 * Pola `kind` tu nie ma z rozmysłu: `ChannelUpdateRequest` kontraktu
 * nie niesie rodzaju, więc rodzaju kanału po założeniu zmienić się nie da.
 * Formularz mówi to Operatorowi wprost, zamiast przyjmować wpis, który rdzeń
 * cicho zignoruje.
 */
export interface OpisZmianyKanalu {
  channelId: string;
  name?: string;
  model?: string;
  enabled?: boolean;
  config?: unknown;
}

/** Trzy czynności zapisujące rejestru kanałów. */
export interface ZrodloKanalow {
  /** `channel.add` — dopisuje wiersz rejestru. */
  zaloz(opis: OpisZalozeniaKanalu): Promise<Wynik<Channel>>;
  /** `channel.update` — zmienia wiersz rejestru wybiórczo. */
  zmien(opis: OpisZmianyKanalu): Promise<Wynik<Channel>>;
  /** `channel.remove` — wykreśla wiersz rejestru. Nieodwracalne. */
  usun(idKanalu: string): Promise<Wynik<string>>;
  /**
   * `channel.check` — sprawdzenie, czy kanał odpowiada.
   *
   * To narzędzie pomocnicze, nie bramka: wynik NICZEGO nie warunkuje — nie
   * wstrzymuje zapisu, nie wyłącza wiersza, nie blokuje wysłania tury. Kanał,
   * który nie odpowiedział minutę temu, bywa sprawny teraz.
   */
  sprawdz(idKanalu: string): Promise<Wynik<ChannelCheckResponse>>;
  /**
   * `channel.credential.status` — STAN poświadczenia, nigdy jego treść.
   *
   * Odpowiedź mówi, czy poświadczenie jest ustawione, jakiego jest rodzaju
   * i gdzie mieszka. Sekretu nie niesie i nieść nie może.
   */
  stanPoswiadczenia(idKanalu: string): Promise<Wynik<ChannelCredentialStatusResponse>>;
}

export function utworzZrodloKanalow(kanal: Kanal): ZrodloKanalow {
  return {
    async zaloz(opis) {
      const wynik = await wywolaj(kanal, Command.ChannelAdd, opis);
      return przenies(
        sprawdzKsztalt(wynik, Command.ChannelAdd, (tresc) => czyKanal(tresc.channel)),
        (tresc) => tresc.channel,
      );
    },

    async zmien(opis) {
      const wynik = await wywolaj(kanal, Command.ChannelUpdate, opis);
      return przenies(
        sprawdzKsztalt(wynik, Command.ChannelUpdate, (tresc) => czyKanal(tresc.channel)),
        (tresc) => tresc.channel,
      );
    },

    async sprawdz(idKanalu) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ChannelCheck, { channelId: idKanalu }),
        Command.ChannelCheck,
        (tresc) => typeof tresc.reachable === 'boolean',
      );
    },

    async stanPoswiadczenia(idKanalu) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ChannelCredentialStatus, { channelId: idKanalu }),
        Command.ChannelCredentialStatus,
        (tresc) => typeof tresc.status === 'object' && tresc.status !== null,
      );
    },

    async usun(idKanalu) {
      const wynik = await wywolaj(kanal, Command.ChannelRemove, { channelId: idKanalu });
      return przenies(
        sprawdzKsztalt(wynik, Command.ChannelRemove, (tresc) => czyTekst(tresc.channelId)),
        (tresc) => tresc.channelId,
      );
    },
  };
}

/**
 * Czy odpowiedź niesie wiersz rejestru o kształcie z kontraktu.
 *
 * Sprawdzamy pola, na które widok patrzy bez osłony: identyfikator (klucz
 * wiersza w wykazie), nazwę i rodzaj (`nazwaKanalu` czyta oba) oraz `enabled`.
 * Rdzeń starszej wersji, który przyśle wiersz bez `enabled`, dałby w widoku
 * kanał „nieczynny" — czyli zdanie nieprawdziwe zamiast odmowy.
 */
function czyKanal(kanal: Channel): boolean {
  if (!czyObiekt(kanal)) return false;
  return (
    czyTekst(kanal.id) &&
    kanal.id.length > 0 &&
    czyTekst(kanal.name) &&
    czyTekst(kanal.kind) &&
    typeof kanal.enabled === 'boolean'
  );
}
