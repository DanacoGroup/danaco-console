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

/** Źródło zapisuje rejestr kanałów modelu: zakładanie, zmianę, usunięcie, sprawdzenie i poświadczenie. */

/** Opis kanału zakładanego w rejestrze zawiera nazwę, rodzaj, model, stan włączenia oraz ustawienia konfiguracyjne kanału. */
export interface OpisZalozeniaKanalu {
  name: string;
  kind: string;
  model?: string;
  enabled?: boolean;
  config?: unknown;
}

/** Opis zmiany wiersza rejestru pomija rodzaj kanału, ponieważ kontrakt zmiany rodzaju kanału po założeniu nie przewiduje. */
export interface OpisZmianyKanalu {
  channelId: string;
  name?: string;
  model?: string;
  enabled?: boolean;
  config?: unknown;
}

/** Interfejs udostępnia trzy czynności zapisujące rejestru kanałów: zakładanie, zmianę oraz usunięcie wiersza rejestru. */
export interface ZrodloKanalow {
  /** `channel.add` — dopisuje wiersz rejestru. */
  zaloz(opis: OpisZalozeniaKanalu): Promise<Wynik<Channel>>;
  /** `channel.update` — zmienia wiersz rejestru wybiórczo. */
  zmien(opis: OpisZmianyKanalu): Promise<Wynik<Channel>>;
  /** `channel.remove` — wykreśla wiersz rejestru. Nieodwracalne. */
  usun(idKanalu: string): Promise<Wynik<string>>;
  /** `channel.check` — sprawdzenie łączności kanału; wynik niczego w zapisie nie warunkuje. */
  sprawdz(idKanalu: string): Promise<Wynik<ChannelCheckResponse>>;
  /** `channel.credential.status` — stan poświadczenia kanału, nigdy jego treść. */
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

/** Sprawdzenie waliduje pola wiersza rejestru wymagane przez widok: identyfikator, nazwę, rodzaj kanału oraz stan włączenia. */
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
