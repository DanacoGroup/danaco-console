import { Command, type Channel } from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Rejestr kanałów modelu widziany przez moduł Translate — cudzy obszar, z którego
 * moduł korzysta, a którego nie prowadzi.
 *
 * Pola `TranslateTargetAddRequest.channelId`
 * i `TranslateBacktranslationRunRequest.channelId` wskazują kanał modelu
 * wykonujący przekład; brak pola bierze kanał czynny okna. Rdzeń je czyta
 * (`adapter_modul_tlumaczenie_panele.go` → `przetlumaczModelem(…, z.ChannelId)`
 * i `kanalZadania`, który wskazany kanał odnajduje w wykazie, a niedobry
 * odrzuca odmową nazwaną, nie cichym zejściem na domyślny). Bez wykazu kanałów
 * nie ma czego wskazać, więc rejestr jest warunkiem wstępnym steru.
 *
 * `enabledOnly: true`, bo ster ma pokazywać to, czym da się przełożyć. Kanał
 * wyłączony rdzeń i tak odrzuci (`kanalZadania` czyta cały wykaz właśnie po to,
 * żeby odróżnić „nie ma takiego" od „jest, ale wyłączony"), więc stawianie go
 * na liście wyboru byłoby zaproszeniem do odmowy. Ta sama nastawa co
 * w `agents/zrodlo-zaplecza.ts`.
 *
 * Rejestr, który nie dotarł, daje wykaz pusty i nie odbiera ani dodania języka,
 * ani tłumaczenia zwrotnego: żądanie idzie wtedy bez `channelId`, a rdzeń
 * bierze kanał domyślny. Powód nieudanego odczytu nie jest połykany — mówi go
 * ster przy pozycji domyślnej.
 *
 * Faz są cztery, bo trzy pustki są różne — ten sam podział, co w
 * `magazyn-translate.ts`: `spoczynek` znaczy „jeszcze nie pytałem", `odczyt` —
 * „pytam", `gotowe` z pustym wykazem — „rdzeń nie zna ani jednego kanału
 * czynnego", `blad` — „nie udało się zapytać". Zlanie ich kazałoby zgadywać,
 * czy czekać, czy zakładać kanał.
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
        // Wykaz poprzedni zostaje: nieudane odświeżenie nie jest powodem, żeby
        // zabrać wybór sprzed chwili, a powód odmowy i tak stanie przy pozycji
        // domyślnej steru.
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
