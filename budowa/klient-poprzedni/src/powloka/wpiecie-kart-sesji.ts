import { Command } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { czyObiekt, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { zadajPrzeniesienieOgniska } from '../protokol/ognisko-sesji';
import type { TozsamoscKlienta } from '../protokol/tozsamosc-klienta';
import { wywolaj } from '../protokol/wywolanie';
import { otworzUsuniecieSesji } from './potwierdzenie-usuniecia';
import { czyRozliczeniePuste, trescRozliczenia } from './rozliczenie-usuniecia';
import { zadajUsuniecieSesji } from './usuniecie-sesji';
import { utworzZrodloKartSesji, type MigawkaKart } from './zrodlo-kart-sesji';

/**
 * Wpięcie pasa kart sesji w rdzeń — jedyne miejsce powłoki, które zna nazwę
 * komendy sesji.
 *
 * Wiąże migawki `zrodlo-kart-sesji` z pasem i wykonuje cztery zamiary: założenie
 * sesji (`session.create`), zamknięcie sesji (`session.close`), trwałe usunięcie
 * (`session.delete`) i przeniesienie ogniska (`session.focus`).
 *
 * Usunięcie nie jest zamknięciem: `session.close` zmienia stan i zostawia zapis,
 * `session.delete` kasuje zapis i jest jedyną drogą utraty danych sesji.
 * Dlatego tylko usunięcie przechodzi przez potwierdzenie z wykazem tego, co
 * ginie, i tylko ono rozlicza się dwoma wykazami rdzenia.
 *
 * Droga do rdzenia udostępnia się raz, przy złożeniu aplikacji. Pas kart
 * powstaje wewnątrz powłoki, która jest widokiem i nie zna ani transportu, ani
 * kontraktu; przeciąganie kanału przez wszystkie jej warstwy wprowadziłoby
 * protokół do każdego odbiorcy powłoki. Bez udostępnienia pas zostaje samym
 * widokiem — tak pracuje stanowisko podglądu, w którym rdzenia nie ma.
 *
 * Każde naciśnięcie dostaje odpowiedź: odmowa rdzenia idzie na pas jego własną
 * treścią, a skład pasa zmienia dopiero kolejna migawka rdzenia.
 */

/** Droga do rdzenia potrzebna pasowi kart. */
export interface DostepDoRdzenia {
  kanal: Kanal;
  /** Tożsamość z powitania; `session.focus` żąda właśnie jej (ognisko = klient). */
  klient: TozsamoscKlienta;
}

/** Tyle pasa, ile potrzeba do zasilenia go rdzeniem — bez zależności zwrotnej. */
export interface PasKart {
  ustawMigawke(migawka: MigawkaKart): void;
  zglosKomunikat(tekst: string, waga?: 'blad' | 'info'): void;
  wybierz(id: string): void;
  czynna(): { readonly id: string } | undefined;
}

/** Zamiary zgłaszane przez pas; rozstrzyga je rdzeń, nie widok. */
export interface ZamiaryKart {
  /** ＋ — założenie sesji, nie okna komunikacji. */
  przyNowej(tytul: string): void;
  /** Zamknięcie karty — zamknięcie sesji w rdzeniu. */
  przyZamknieciu(id: string): void;
  /** Wybór karty — przeniesienie ogniska tego klienta. */
  przyWyborze(id: string): void;
  /**
   * Kosz — trwałe usunięcie sesji po potwierdzeniu. Tytuł idzie razem
   * z identyfikatorem, bo potwierdzenie ma nazwać stratę, a nie pokazać ciąg
   * znaków.
   */
  przyUsunieciu(id: string, tytul: string): void;
}

/**
 * Scena sesji pokazuje okna sesji powiązanej z tym połączeniem. Ognisko można
 * przenieść na inną kartę, ale przeniesienie sceny wymaga powiązania połączenia
 * z tamtą sesją (`session.bind`), a droga do tego prowadzi przez Centrum
 * dowodzenia; komunikat mówi to wprost.
 */
const SCENA_ZOSTAJE =
  'Ognisko przeszło na tę sesję. Scena pokazuje okna sesji bieżącego połączenia — po okna tamtej sesji wróć przyciskiem „Wróć do sesji" w Centrum dowodzenia.';

/** Droga do rdzenia udostępniona przez aplikację; `null` poza produktem. */
let dostep: DostepDoRdzenia | null = null;

/** Podaje pasom kart drogę do rdzenia. Wywołuje ją złożenie aplikacji. */
export function udostepnijRdzenPasomKart(nowy: DostepDoRdzenia): void {
  dostep = nowy;
}

/**
 * Wiąże świeżo złożony pas z rdzeniem. Zwraca zamiary, którym pas oddaje
 * naciśnięcia, albo `null`, gdy rdzeń nie został udostępniony.
 */
export function zwiazPasKartZRdzeniem(pas: PasKart): ZamiaryKart | null {
  if (dostep === null) return null;
  const { kanal, klient } = dostep;
  const zrodlo = utworzZrodloKartSesji(kanal, klient);

  pas.ustawMigawke(zrodlo.migawka());
  zrodlo.naZmiane((migawka) => pas.ustawMigawke(migawka));

  /** Przenosi ognisko klienta i mówi, czego przeniesienie nie robi. */
  function ogniskuj(idSesji: string, poprzednia: string | undefined): void {
    void zadajPrzeniesienieOgniska(kanal, { sessionId: idSesji, clientId: klient.id }).then(
      (wynik) => {
        if (!wynik.udany) {
          if (poprzednia !== undefined) pas.wybierz(poprzednia);
          pas.zglosKomunikat(wynik.blad?.message ?? 'Rdzeń nie przeniósł ogniska.');
          return;
        }
        if (idSesji !== kanal.sesja().id()) pas.zglosKomunikat(SCENA_ZOSTAJE, 'info');
      },
    );
  }

  return {
    przyNowej(tytul) {
      void wywolaj(kanal, Command.SessionCreate, { title: tytul })
        .then((wynik) => sprawdzKsztalt(wynik, Command.SessionCreate, (tresc) => czyObiekt(tresc.session)))
        .then((wynik) => {
          const zalozona = wynik.wynik?.session;
          if (!wynik.udany || zalozona === undefined) {
            pas.zglosKomunikat(wynik.blad?.message ?? 'Rdzeń nie założył sesji.');
            return;
          }
          ogniskuj(zalozona.id, pas.czynna()?.id);
        });
    },

    przyZamknieciu(id) {
      void wywolaj(kanal, Command.SessionClose, { sessionId: id })
        .then((wynik) => sprawdzKsztalt(wynik, Command.SessionClose, (tresc) => czyObiekt(tresc.session)))
        .then((wynik) => {
          if (!wynik.udany) pas.zglosKomunikat(wynik.blad?.message ?? 'Rdzeń nie zamknął sesji.');
        });
    },

    przyUsunieciu: (id, tytul) => usunSesjeKarty(kanal, pas, id, tytul),

    przyWyborze(id) {
      const poprzednia = pas.czynna()?.id;
      if (poprzednia === id) return;
      // Zaznaczenie idzie od razu, bo wybór karty jest faktem widoku; odmowa
      // rdzenia cofa je do karty poprzedniej.
      pas.wybierz(id);
      ogniskuj(id, poprzednia);
    },
  };
}

/**
 * Trwałe usunięcie sesji jednej karty.
 *
 * Potwierdzenie prowadzi całą czynność wraz ze stanami obowiązkowymi, a pas
 * kart dostaje samo rozliczenie, natychmiast po odpowiedzi rdzenia. Odmowę
 * nazwał już modal wraz z kodem, więc pas jej nie powtarza.
 *
 * Rozliczenie, w którym nic nie zginęło, idzie na pas z wagą błędu: „nie
 * usunąłem niczego" nie jest powodzeniem zamówionej czynności.
 */
function usunSesjeKarty(kanal: Kanal, pas: PasKart, id: string, tytul: string): void {
  const nazwa = (idSesji: string): string => (idSesji === id ? tytul : '');
  void otworzUsuniecieSesji(
    [{ id, tytul }],
    (idSesji) => zadajUsuniecieSesji(kanal, idSesji, true),
    nazwa,
  ).then((rozliczenie) => {
    if (rozliczenie === null) return;
    pas.zglosKomunikat(
      trescRozliczenia(rozliczenie, nazwa),
      czyRozliczeniePuste(rozliczenie) ? 'blad' : 'info',
    );
  });
}
