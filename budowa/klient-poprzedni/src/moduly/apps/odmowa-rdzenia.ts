import { ErrorCode, type Command, type RequestOf, type ResponseOf } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';

/**
 * Wywołanie komendy modułu Apps, które kończy się także wtedy, gdy rdzeń
 * odmówi zdarzeniem `<obszar>.unknown`.
 *
 * Powód jest mechaniczny. Komenda bez uchwytu w rdzeniu wraca kopertą typu
 * `<obszar>.unknown` z polem `requestedType`. Ta koperta nie niesie pola
 * `status`, a korelacja żądań rozstrzyga wyłącznie koperty ze statusem
 * (`protokol/koperta.ts`). Obietnica zwykłego `wywolaj()` zostałaby więc
 * nierozstrzygnięta na zawsze, a okno stałoby w stanie ładowania bez końca.
 * Dlatego czytamy też dziennik nierozpoznanych.
 *
 * Rozwiązanie nie buduje drugiej drogi do rdzenia: żądanie idzie tym
 * samym `kanal.wyslij`, a odmowę odczytujemy z dziennika komunikatów
 * nierozpoznanych, który kanał prowadzi dla wszystkich obszarów kontraktu
 * (`polaczenie/dziennik-nieznanych.ts`). Moduł nie zna literału `apps.unknown`
 * i nie zakłada własnej subskrypcji zdarzenia odmowy.
 *
 * Wynik odmowy jest zwykłym `Wynik` z polem `blad` — okno pokazuje go tak samo
 * jak każdą inną odmowę rdzenia i nie potrzebuje osobnej gałęzi widoku.
 */
export type WywolanieApps = <K extends Command>(
  komenda: K,
  zadanie: RequestOf<K>,
) => Promise<Wynik<ResponseOf<K>>>;

export function utworzWywolanieApps(kanal: Kanal): WywolanieApps {
  return <K extends Command>(komenda: K, zadanie: RequestOf<K>) =>
    new Promise<Wynik<ResponseOf<K>>>((rozstrzygnij) => {
      let idZadania = '';
      let odsubskrybuj: Odsubskrybuj = () => undefined;
      let rozstrzygniete = false;

      function zakoncz(wynik: Wynik<ResponseOf<K>>): void {
        if (rozstrzygniete) return;
        rozstrzygniete = true;
        odsubskrybuj();
        rozstrzygnij(wynik);
      }

      // Subskrypcja przed wysyłką: wpis dziennika może paść w tej samej pętli
      // zdarzeń co odpowiedź, a kolejność ich nadejścia nie jest niczym
      // zagwarantowana.
      odsubskrybuj = kanal.dziennikNieznanych().naWpis((wpis) => {
        if (wpis.idZadania !== idZadania || idZadania === '') return;
        zakoncz({ udany: false, blad: bladOdmowy(wpis.zadanyTyp, wpis.powod) });
      });

      idZadania = kanal.wyslij(komenda, zadanie, zakoncz);
    });
}

/**
 * Błąd nazywający komendę, której rdzeń nie zna.
 *
 * Nazwa żądanego typu pochodzi z odpowiedzi rdzenia (`requestedType`), nie
 * z literału w kliencie — okno mówi Operatorowi dokładnie to, co odpowiedział
 * rdzeń, a nie to, czego się spodziewał klient.
 */
function bladOdmowy(zadanyTyp: string, powod: string) {
  const nazwa = zadanyTyp === '' ? 'komendy, której rdzeń nie nazwał' : zadanyTyp;
  const uzupelnienie = powod === '' ? '' : ` Rdzeń podał powód: ${powod}.`;
  return {
    code: ErrorCode.NotFound,
    message:
      `Rdzeń nie ma uchwytu dla ${nazwa} — komenda jest w kontrakcie, wykonania jeszcze nie ma.` +
      uzupelnienie,
    retryable: false,
  };
}
