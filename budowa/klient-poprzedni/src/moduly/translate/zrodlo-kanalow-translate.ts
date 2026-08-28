import { Command, type Channel } from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Rejestr kanałów modelu widziany przez moduł Translate jest cudzym obszarem, z którego moduł
 * korzysta, a którego nie prowadzi; bez wykazu kanałów nie ma czego wskazać sterowi.
 */
export type FazaKanalow = 'spoczynek' | 'odczyt' | 'gotowe' | 'blad';

export interface ZrodloKanalowTranslate {
  /** Kanały czynne znane rdzeniowi przy ostatnim odczycie. */
  kanaly(): readonly Channel[];
  /** Faza odczytu rejestru. */
  faza(): FazaKanalow;
  /** Powód nieudanego odczytu; pusty, gdy odczyt się udał. */
  powod(): string;
  /** Odczytuje rejestr z rdzenia; wolno wołać wielokrotnie. */
  wczytaj(): Promise<void>;
  /** Nasłuch zmiany wykazu — stery przerysowują się bez pytania rdzenia. */
  obserwuj(sluchacz: () => void): () => void;
  /** Zdejmuje wszystkich nasłuchujących. */
  zapomnijSluchaczy(): void;
}

export function utworzZrodloKanalowTranslate(kanal: Kanal): ZrodloKanalowTranslate {
  const sluchacze = new Set<() => void>();
  let wykaz: readonly Channel[] = [];
  let stanFazy: FazaKanalow = 'spoczynek';
  let stanPowodu = '';

  function oglos(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  function ustaw(faza: FazaKanalow, powod: string): void {
    stanFazy = faza;
    stanPowodu = powod;
    oglos();
  }

  return {
    kanaly: () => wykaz,
    faza: () => stanFazy,
    powod: () => stanPowodu,

    async wczytaj() {
      ustaw('odczyt', '');
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ChannelList, { enabledOnly: true }),
        Command.ChannelList,
        (tresc) => czyTablica(tresc.channels),
      );
      if (!wynik.udany || wynik.wynik === undefined) {
        // Wykaz poprzedni zostaje: nieudane odświeżenie nie jest powodem, żeby zabrać wybór sprzed chwili.
        ustaw('blad', wynik.blad?.message ?? 'rdzeń nie podał powodu');
        return;
      }
      wykaz = wynik.wynik.channels;
      ustaw('gotowe', '');
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    zapomnijSluchaczy: () => sluchacze.clear(),
  };
}
