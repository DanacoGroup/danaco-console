import { Command, KOMENDY, PROTOCOL_VERSION } from '../../../../shared/contract';
import { przyciskBezKomendy } from '../../modele/kontrolki-formularza';
import type { Kanal } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { tozsamoscKlienta } from '../../protokol/tozsamosc-klienta';
import { wywolaj } from '../../protokol/wywolanie';

/** Powody nieczynnych kontrolek modułu Workspace składają się z wykazu komend rdzenia i kontraktu. */

/** Komendy kontraktu jako zbiór pozwalają rozpoznać, po czyjej stronie jest brak, gdy rdzeń komendy nie rejestruje. */
const W_KONTRAKCIE: ReadonlySet<string> = new Set<string>(KOMENDY);

/** Wykaz komend rdzenia wraz z kontrolkami, które z niego biorą swoje zdanie, tworzy jedno źródło powodów nieczynności. */
export interface WykazBrakow {
  /** Pyta rdzeń o wykaz komend i przerysowuje powody wszystkich kontrolek. */
  odczytaj(): Promise<void>;
  /** Kontrolka bez pokrycia, której powód dopisuje się sam po odpowiedzi rdzenia na wykaz komend. */
  przyciskBraku(
    etykieta: string,
    czynnosc: string,
    ...komendy: readonly string[]
  ): HTMLButtonElement;
}

export function utworzWykazBrakow(kanal: Kanal): WykazBrakow {
  /** Wykaz z rdzenia; `null`, dopóki powitanie nie wróciło. */
  let komendyRdzenia: ReadonlySet<string> | null = null;
  /** Treść odmowy powitania; pusta, dopóki nic się nie zepsuło. */
  let odmowa = '';
  /** Kontrolki, które trzeba przerysować po odpowiedzi rdzenia. */
  const zalezne: Array<() => void> = [];

  /** Zdanie o pokryciu komendy rozróżnia cztery stany: odmowę, pytanie w drodze, rejestrację i brak. */
  function zdanieOKomendzie(nazwa: string): string {
    if (odmowa !== '') {
      return `Rdzeń nie oddał wykazu komend (${odmowa}) — o pokryciu komendy ${nazwa} nic nie wiadomo.`;
    }
    if (komendyRdzenia === null) {
      return `Sprawdzanie u rdzenia, czy obsługuje komendę ${nazwa}…`;
    }
    if (komendyRdzenia.has(nazwa)) {
      return `Rdzeń rejestruje komendę ${nazwa}, ale to okno jej nie wywołuje — czynność niezbudowana.`;
    }
    return W_KONTRAKCIE.has(nazwa)
      ? `Rdzeń tej wersji nie rejestruje komendy ${nazwa}, choć kontrakt ją niesie — ` +
          'brak jest po stronie rdzenia (zmierzone powitaniem connection.hello).'
      : `Ani kontrakt, ani rdzeń nie znają komendy ${nazwa} — czynności nie ma czym wykonać.`;
  }

  /** Powód nieczynności kontrolki opartej na jednej komendzie albo na kilku naraz. */
  function powod(czynnosc: string, komendy: readonly string[]): string {
    if (komendy.length === 0) {
      return `${czynnosc}: okno nie zna komendy, która by tego dokonała.`;
    }
    return `${czynnosc}. ${komendy.map(zdanieOKomendzie).join(' ')}`;
  }

  return {
    async odczytaj() {
      const klient = tozsamoscKlienta();
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ConnectionHello, {
          clientId: klient.id,
          clientVersion: klient.wersja,
          protocolVersion: PROTOCOL_VERSION,
        }),
        Command.ConnectionHello,
        (tresc) => czyTablica(tresc.commands),
      );
      if (!wynik.udany || wynik.wynik === undefined) {
        odmowa = `kod ${wynik.blad?.code ?? 'brak'}`;
        komendyRdzenia = null;
      } else {
        odmowa = '';
        komendyRdzenia = new Set(wynik.wynik.commands ?? []);
      }
      for (const przerysuj of zalezne) przerysuj();
    },

    przyciskBraku(etykieta, czynnosc, ...komendy) {
      let kontrolka = przyciskBezKomendy(etykieta, powod(czynnosc, komendy));
      // Powód idzie trzema drogami naraz, więc kontrolka powstaje na nowo i wchodzi na miejsce poprzedniej.
      zalezne.push(() => {
        const nowa = przyciskBezKomendy(etykieta, powod(czynnosc, komendy));
        if (kontrolka.parentNode === null) {
          // Kontrolka poza drzewem: trzyma ją jeszcze wywołujący, więc podmiana węzła nic by nie dała.
          kontrolka.title = nowa.title;
          kontrolka.setAttribute('aria-description', nowa.title);
          return;
        }
        kontrolka.replaceWith(nowa);
        kontrolka = nowa;
      });
      return kontrolka;
    },
  };
}
