import { opisOdmowy } from '../../komponenty/odmowa';
import type { TrescNotatki } from './czynnosci-notatek';
import { skutekPrzekazania, skutekZapisuNotatki, skutekZrzutu } from './skutek-zapisu';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Trzy czynności paska dolnego, które idą do rdzenia: zrzut ekranu, przekazanie
 * treści do modułu Translate oraz adnotacja zaznaczonego fragmentu. Każda odpowiada
 * tak samo: zapowiedź, a po odpowiedzi rdzenia wynik albo powód odmowy.
 */
export interface CzynnosciPaskaDolnego {
  /** `browser.snapshot.get` ze zrzutem ekranu i źródłem strony. */
  zrzutEkranu(): Promise<void>;
  /** `context.transfer` treści strony do modułu Translate. */
  doTlumaczenia(): Promise<void>;
  /** `browser.note.add` z cytatem zaznaczonego fragmentu. */
  adnotujFragment(): Promise<void>;
}

/**
 * Treść notatki zakładanej z podświetlenia. Jedno brzmienie służy i żądaniu, i ocenie
 * odpowiedzi, więc zdanie o skutku porównuje się z tym samym napisem, który poszedł
 * do rdzenia, a nie z napisem przepisanym drugi raz.
 */
const TRESC_ADNOTACJI = 'Podświetlony fragment strony';

export function utworzCzynnosciPaskaDolnego(
  stan: StanPrzegladania,
  powiedz: (tresc: string, powodzenie: boolean) => void,
): CzynnosciPaskaDolnego {
  return {
    async zrzutEkranu() {
      const idOkna = stan.idOkna();
      if (idOkna === '') {
        powiedz(stan.powod(), false);
        return;
      }
      powiedz('Pobranie migawki ze zrzutem ekranu w toku…', true);
      const wynik = await stan.zrodlo.migawka({
        windowId: idOkna,
        includeScreenshot: true,
        includeHtml: true,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Zrzut ekranu', wynik.blad?.code, wynik.blad?.message), false);
        return;
      }
      stan.wchlonMigawke(wynik.wynik.snapshot);
      const skutek = skutekZrzutu(wynik.wynik.snapshot);
      powiedz(skutek.zdanie, skutek.udany);
    },

    async doTlumaczenia() {
      const idOkna = stan.idOkna();
      const migawka = stan.migawka();
      if (idOkna === '' || migawka === null) {
        powiedz('Nie ma czego tłumaczyć — najpierw pobierz migawkę strony.', false);
        return;
      }
      powiedz('Przekazanie treści strony do modułu Translate…', true);
      const zaznaczony = stan.zaznaczenie();
      const wynik = await stan.zapisy.przekaz(idOkna, 'translate', {
        prompt: zaznaczony === '' ? (migawka.text ?? migawka.url) : zaznaczony,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(
          opisOdmowy('Przekazanie do Translate', wynik.blad?.code, wynik.blad?.message),
          false,
        );
        return;
      }
      const skutek = skutekPrzekazania(wynik.wynik, 'translate', 'Treść strony przekazana');
      powiedz(skutek.zdanie, skutek.udany);
    },

    async adnotujFragment() {
      const idOkna = stan.idOkna();
      const fragment = stan.zaznaczenie();
      if (idOkna === '' || fragment === '') {
        powiedz('Podświetlenie adnotuje zaznaczony fragment — najpierw go zaznacz.', false);
        return;
      }
      powiedz('Zapis adnotacji zaznaczonego fragmentu…', true);
      // Zamówienie spisane raz idzie i do rdzenia, i do oceny jego odpowiedzi.
      const zapis: TrescNotatki = { tresc: TRESC_ADNOTACJI, idZrodla: '', cytat: fragment };
      const wynik = await stan.zrodlo.dodajNotatke({
        windowId: idOkna,
        content: zapis.tresc,
        quote: zapis.cytat,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        powiedz(opisOdmowy('Adnotacja fragmentu', wynik.blad?.code, wynik.blad?.message), false);
        return;
      }
      stan.zebrane.dopiszNotatke(wynik.wynik.note);
      const skutek = skutekZapisuNotatki(wynik.wynik.note, zapis);
      powiedz(skutek.zdanie, skutek.udany);
    },
  };
}
