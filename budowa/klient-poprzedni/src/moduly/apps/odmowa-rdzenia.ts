import { ErrorCode, type Command, type RequestOf, type ResponseOf } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';

/**
 * Wywołanie komendy modułu Apps, które rozstrzyga się także wtedy, gdy rdzeń
 * odmawia kopertą bez pola `status`, pomijaną przez korelację żądań. Odmowę
 * czyta z dziennika nierozpoznanych i zwraca jako zwykły `Wynik` z polem
 * `blad`.
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

      // Subskrypcja przed wysyłką, bo kolejność wpisu dziennika i odpowiedzi
      // nie jest zagwarantowana.
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
