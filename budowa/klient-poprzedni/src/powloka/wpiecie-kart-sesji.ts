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

// Wpięcie pasa kart sesji w rdzeń — jedyne miejsce powłoki, które zna nazwę komendy sesji.

/** Droga do rdzenia potrzebna pasowi kart sesji, wraz z tożsamością klienta udostępnioną tej samej powłoce. */
export interface DostepDoRdzenia {
  kanal: Kanal;
  /** Tożsamość z powitania; przeniesienie ogniska żąda właśnie jej, bo ognisko równa się klient. */
  klient: TozsamoscKlienta;
}

/** Tyle pasa, ile potrzeba do zasilenia go rdzeniem — bez jakiejkolwiek zależności zwrotnej tego samego pasa. */
export interface PasKart {
  ustawMigawke(migawka: MigawkaKart): void;
  zglosKomunikat(tekst: string, waga?: 'blad' | 'info'): void;
  wybierz(id: string): void;
  czynna(): { readonly id: string } | undefined;
}

/** Zamiary zgłaszane przez pas kart sesji; rozstrzyga je zawsze rdzeń, nigdy sam widok tego samego pasa. */
export interface ZamiaryKart {
  /** ＋ — założenie sesji, nie okna komunikacji. */
  przyNowej(tytul: string): void;
  /** Zamknięcie karty — zamknięcie sesji w rdzeniu. */
  przyZamknieciu(id: string): void;
  /** Wybór karty — przeniesienie ogniska tego klienta. */
  przyWyborze(id: string): void;
  // Kosz — trwałe usunięcie sesji po potwierdzeniu, tytuł idzie razem z identyfikatorem sesji.
  przyUsunieciu(id: string, tytul: string): void;
}

/** Scena sesji pokazuje okna sesji powiązanej z tym połączeniem, a przeniesienie wymaga jej powiązania. */
const SCENA_ZOSTAJE =
  'Ognisko przeszło na tę sesję. Scena pokazuje okna sesji bieżącego połączenia — po okna tamtej sesji wróć przyciskiem „Wróć do sesji" w Centrum dowodzenia.';

/** Droga do rdzenia udostępniona przez aplikację; wartość pusta oznacza pracę poza tym produktem klienta. */
let dostep: DostepDoRdzenia | null = null;

/** Podaje pasom kart drogę do rdzenia; wywołuje ją wyłącznie złożenie tej samej aplikacji przy jej starcie. */
export function udostepnijRdzenPasomKart(nowy: DostepDoRdzenia): void {
  dostep = nowy;
}

/** Wiąże świeżo złożony pas z rdzeniem, zwracając zamiary, którym ten pas oddaje wszystkie naciśnięcia. */
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
      // Zaznaczenie idzie od razu, bo wybór karty jest faktem widoku; odmowa rdzenia cofa je do poprzedniej.
      pas.wybierz(id);
      ogniskuj(id, poprzednia);
    },
  };
}

/** Trwałe usunięcie sesji jednej karty; pas kart dostaje samo rozliczenie zaraz po odpowiedzi tego rdzenia. */
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
