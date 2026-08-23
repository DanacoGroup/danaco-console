import { opisOdmowy } from '../../komponenty/odmowa';
import type { TrescNotatki } from './czynnosci-notatek';
import { skutekPrzekazania, skutekZapisuNotatki, skutekZrzutu } from './skutek-zapisu';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Trzy czynności paska dolnego, które idą do rdzenia: zrzut ekranu,
 * przekazanie treści do modułu Translate i adnotacja zaznaczonego fragmentu.
 *
 * Jedna odpowiedzialność: wywołania rdzenia paska dolnego. Pasek składa
 * kontrolki, a to jest rozmowa z rdzeniem — stąd osobny plik.
 *
 * Każda z trzech odpowiada tak samo: zapowiedź, a po odpowiedzi rdzenia albo
 * wynik, albo powód odmowy. Milczenie po naciśnięciu jest zakazane.
 *
 * Zdanie końcowe buduje się z treści odpowiedzi, nie z samego `udany`
 * (`skutek-zapisu.ts`). Rdzeń pobiera dziś stronę bez uruchamiania przeglądarki
 * (`przegladarka_pobieranie.go`), więc odpowiedź udana potrafi przyjść z pustym
 * `screenshotRef` — okno nie ma wtedy prawa twierdzić o obrazie strony. Zdanie
 * o skutku czyta pole odpowiedzi, więc pozostanie prawdziwe także wtedy, gdy
 * zrzuty zacznie oddawać uchwyt `browser.screenshot.capture`.
 */
export interface CzynnosciPaskaDolnego {
  /** `browser.snapshot.get` ze zrzutem ekranu i źródłem strony. */
  zrzutEkranu(): Promise<void>;
  /** `context.transfer` treści strony do modułu Translate. */
  doTlumaczenia(): Promise<void>;
  /** `browser.note.add` z cytatem zaznaczonego fragmentu. */
  adnotujFragment(): Promise<void>;
}

/** Treść notatki zakładanej z podświetlenia — jedno brzmienie dla żądania i oceny. */
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
      // Zamówienie spisane raz i to samo idzie do rdzenia oraz do oceny jego
      // odpowiedzi: zdanie o cytacie ma się opierać na wierszu, który wrócił,
      // a nie na tym, że okno cytat wysłało.
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
