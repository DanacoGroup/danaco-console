import { uruchomMotyw } from '../motyw/motyw';
import { opisPoczatkowy } from '../okno-komunikacji/opis-okna';
import { udostepnijRdzenPasomKart } from '../powloka/wpiecie-kart-sesji';
import type { KodSrodowiska } from '../strona-glowna/indeks';
import { utworzGospodarza } from './gospodarz-dokumentu';
import { zlozPolaczenieZRdzeniem, type PolaczenieZRdzeniem } from './polaczenie-z-rdzeniem';
import { utworzRouter, type Router } from './router';
import { Trasa } from './trasy';
import { utworzWidokPulpitu } from './widok-pulpitu';
import { utworzWidokSrodowiska, type WidokSrodowiska } from './widok-srodowiska';
import { utworzWidokStronyGlownej, type WidokStronyGlownej } from './widok-strony-glownej';

/** Złożona aplikacja gotowa do uruchomienia. */
export interface Aplikacja {
  /** Droga do rdzenia — punkt wejścia otwiera połączenie. */
  rdzen: PolaczenieZRdzeniem;
  /** Router tras — punkt wejścia pokazuje trasę zapisaną w adresie. */
  router: Router;
  /**
   * Gotowość Centrum dowodzenia — pierwszy odczyt strony domknięty.
   *
   * Punkt wejścia podaje ją scenie wejścia (`ladowanie/`), żeby ekran między
   * bramką a produktem schodził w chwili, gdy strona ma czym stanąć.
   * Nieodwiedzona trasa strony głównej znaczy gotowość natychmiastową: nie ma
   * na co czekać, skoro widoku nikt nie zbudował.
   */
  stronaGlownaGotowa(): Promise<void>;
}

/**
 * Składa aplikację: powołuje warstwy wspólne i wiąże trzy trasy w jeden przepływ.
 * Żadnego elementu widoku tu nie ma — każdą trasę buduje jej własny moduł.
 *
 * Motyw rusza przed montażem czegokolwiek, żeby dokument dostał żetony motywu,
 * zanim pojawi się pierwszy element. Router buduje widok leniwie, przy pierwszym
 * wejściu na trasę, i już go nie porzuca — karty opuszczonego środowiska trwają
 * w tle i wracają z pełnym stanem. Środowisko przygotowuje się od razu, bo to
 * ono podpina przepływ komunikatów. Uzgodnienia z rdzeniem tu nie rozpoczynamy:
 * robi to przepływ komunikatów, gdy transport zgłosi stan „połączony" —
 * wywołanie stąd dałoby drugie powitanie. Tożsamość klienta pochodzi z powitania
 * (`rdzen.uzgodnienie.klient`), bo `session.bind` i `session.focus` żądają
 * `clientId` przedstawionego rdzeniowi raz.
 */
export function zlozAplikacje(): Aplikacja {
  const opis = opisPoczatkowy();

  uruchomMotyw();

  const gospodarz = utworzGospodarza(opis);
  const rdzen = zlozPolaczenieZRdzeniem(opis);
  const router = utworzRouter(gospodarz);

  // Pas kart sesji powstaje w głębi powłoki, która o rdzeniu nie wie (buduje ją
  // także stanowisko podglądu). Drogę do rdzenia podaje mu złożenie aplikacji.
  udostepnijRdzenPasomKart({ kanal: rdzen.kanal, klient: rdzen.uzgodnienie.klient });

  let stronaGlowna: WidokStronyGlownej | null = null;
  let srodowisko: WidokSrodowiska | null = null;

  /**
   * Wchodzi do środowiska. Router pokazuje trasę jako pierwszy, bo dopiero
   * wtedy powstaje widok środowiska; ustawienie środowiska i pozycji nawigacji
   * idzie po nim.
   */
  function wejdzDoSrodowiska(kod: KodSrodowiska, modul?: string): void {
    router.pokaz(Trasa.Srodowisko);
    srodowisko?.ustawSrodowisko(kod, modul);
    stronaGlowna?.ustawSrodowiskoCzynne(kod);
  }

  router.zarejestruj(Trasa.StronaGlowna, () => {
    stronaGlowna = utworzWidokStronyGlownej({
      projekt: opis.projekt,
      kanal: rdzen.kanal,
      klient: rdzen.uzgodnienie.klient,
      naTrase: router.pokaz,
      naSrodowisko: wejdzDoSrodowiska,
    });
    return stronaGlowna;
  });

  router.zarejestruj(Trasa.Srodowisko, () => {
    srodowisko = utworzWidokSrodowiska({ rdzen, opis, naTrase: router.pokaz });
    return srodowisko;
  });

  router.zarejestruj(Trasa.Pulpit, () =>
    utworzWidokPulpitu({
      projekt: opis.projekt,
      kanal: rdzen.kanal,
      naTrase: router.pokaz,
      naSrodowisko: wejdzDoSrodowiska,
    }),
  );

  router.przygotuj(Trasa.Srodowisko);

  return {
    rdzen,
    router,
    stronaGlownaGotowa: () => stronaGlowna?.gotowa ?? Promise.resolve(),
  };
}
