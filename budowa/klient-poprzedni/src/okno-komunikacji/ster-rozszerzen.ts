import { ExtensionKind, ExtensionOrigin, type Extension } from '../../../shared/contract';
import type { PozycjaMenu } from '../komponenty/menu-drzewo';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { ZrodloRozszerzen } from '../moduly/agents/zrodlo-rozszerzen';
import { utworzSterNastawy, type SterPaska } from './ster-nastawy';

/**
 * Rozszerzenia — ster paska zlecenia. Nie jest jednym menu, tylko wejściem do
 * rodzin: pierwszy poziom drzewa to cztery gałęzie po `ExtensionKind` —
 * serwery MCP, wtyczki, integracje API, skille — a nie jedna lista wszystkiego.
 * Podział pochodzi z wyliczenia kontraktu.
 *
 * Katalog jest czytany przez gotowe źródło komend `extension.*`
 * (`moduly/agents/zrodlo-rozszerzen.ts`); ster nie wie, jak wygląda koperta.
 * Import z `moduly/` nie zamyka pętli, bo `zrodlo-rozszerzen.ts` nie importuje
 * ani jednego pliku z `okno-komunikacji/` — zna wyłącznie kontrakt i protokół.
 *
 * Liść jest przełącznikiem, a nie wyborem: rozszerzenia nie wykluczają się
 * nawzajem, więc zgaszenie jednego nie jest wybraniem innego.
 *
 * Opis pozycji pochodzi z rdzenia — `Extension.description` idzie do menu
 * dosłownie, a rozszerzenie bez opisu zostaje w menu bez zdania.
 *
 * Wykaz jest zawężony do zainstalowanych: menu obsługuje podłączenie („czego
 * używam teraz"), a rejestrację („co w ogóle mam") prowadzi katalog rozszerzeń
 * w module Agents. Pozycja niezainstalowana nie ma czego włączać.
 */

/** Nazwa rodzajowa nastawy — idzie do `aria-label`, nie na ekran. */
const NASTAWA = 'Rozszerzenia';

/** Przedrostek klucza liścia; oddziela identyfikatory od kluczy gałęzi. */
const PRZEDROSTEK = 'rozszerzenie:';

/** Etykieta uchwytu przy pustym wykazie — mówi o braku, niczego nie udaje. */
export const BRAK_ROZSZERZEN = 'Bez rozszerzeń';

/** Etykieta uchwytu, dopóki katalog nie wrócił z rdzenia. */
export const WYKAZ_CZYTANY = 'Wczytuję…';

/** Nazwy rodzin w kolejności wyświetlania — cztery gałęzie pierwszego poziomu. */
const RODZINY: readonly { rodzaj: ExtensionKind; nazwa: string; opis: string }[] = [
  {
    rodzaj: ExtensionKind.Mcp,
    nazwa: 'Serwery MCP',
    opis: 'Mosty do narzędzi zewnętrznych mówiące protokołem MCP.',
  },
  {
    rodzaj: ExtensionKind.Plugin,
    nazwa: 'Wtyczki',
    opis: 'Rozszerzenia dokładające możliwości po stronie konsoli.',
  },
  {
    rodzaj: ExtensionKind.Api,
    nazwa: 'Integracje API',
    opis: 'Połączenia z usługami zewnętrznymi po ich własnym interfejsie.',
  },
  {
    rodzaj: ExtensionKind.Skill,
    nazwa: 'Skille',
    opis: 'Spakowane zestawy instrukcji dostępne ekspertowi w zleceniu.',
  },
];

/** Nagłówki grup po źródle pochodzenia. */
const ZRODLA: readonly { origin: ExtensionOrigin; nazwa: string }[] = [
  { origin: ExtensionOrigin.Danaco, nazwa: 'Danaco' },
  { origin: ExtensionOrigin.Personal, nazwa: 'Moje' },
];

/** Zależności steru — wąskie i wstrzykiwane. */
export interface ZaleznosciSteruRozszerzen {
  /** Gotowe źródło komend `extension.*`; ster go nie zakłada. */
  zrodlo: ZrodloRozszerzen;
  /** Subskrypcja `extension.changed` — katalog zmienia się także spoza paska. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
  /**
   * Droga do katalogu rozszerzeń (stopka menu). Pominięta znaczy „to menu jej
   * nie ma": okno katalogu żyje wewnątrz modułu Agents i otwiera je
   * powierzchnia modułu, a stopka wołająca coś, czego wołający nie umie
   * wykonać, byłaby wierszem meldującym pracę bez pracy.
   */
  otworzKatalog?(): void;
}

/** Ster paska ze zdejmowalnym nasłuchem `extension.changed`. */
export type SterRozszerzen = SterPaska & { rozlacz(): void };

export function utworzSterRozszerzen(
  zaleznosci: ZaleznosciSteruRozszerzen,
): SterRozszerzen {
  /** Katalog potwierdzony przez rdzeń; `null` znaczy „jeszcze nie wiem". */
  let katalog: Extension[] | null = null;

  const ster = utworzSterNastawy({
    nastawa: NASTAWA,
    ikona: 'plus',
    ...(zaleznosci.otworzKatalog === undefined
      ? {}
      : {
          stopka: {
            nazwa: 'Katalog rozszerzeń',
            opis:
              'To menu włącza i wyłącza to, co już masz. Dołożenie nowej pozycji ' +
              'jest rejestracją i odbywa się w katalogu rozszerzeń.',
            ikona: 'biblioteka',
            wykonaj: () => zaleznosci.otworzKatalog?.(),
          },
        }),
    wykonaj: (klucz) => przelacz(klucz),
    odswiez: () => odswiez(),
  });

  /**
   * Przełączenie jednej pozycji.
   *
   * Stan bierze się z katalogu, nie z pozycji menu: między narysowaniem wykazu
   * a kliknięciem rdzeń mógł rozgłosić zmianę, a wysłanie `enabled` odwrotnego
   * do nieaktualnego przekonania widoku cofnęłoby cudzą decyzję.
   */
  async function przelacz(klucz: string): Promise<void> {
    const id = klucz.slice(PRZEDROSTEK.length);
    const pozycja = katalog?.find((rozszerzenie) => rozszerzenie.id === id);
    if (pozycja === undefined) {
      throw new Error(
        `Rozszerzenie ${id} zniknęło z katalogu między narysowaniem wykazu a kliknięciem. ` +
          'Wykaz odświeży się sam po zmianie zgłoszonej przez rdzeń.',
      );
    }
    const wynik = await zaleznosci.zrodlo.przelacz(id, !pozycja.enabled);
    const potwierdzone = wynik.wynik;
    if (!wynik.udany || potwierdzone === undefined) {
      throw new Error(
        opisOdmowyBledu(`Przełączenie rozszerzenia „${pozycja.name}" (extension.toggle)`, wynik.blad),
      );
    }
    // Odpowiedź rdzenia jest stanem potwierdzonym — wpisujemy ją do katalogu
    // zamiast czekać na pełny odczyt. Zdarzenie `extension.changed` i tak
    // przyjdzie; wykaz nie ma migać w międzyczasie starą wartością.
    katalog = (katalog ?? []).map((pozycjaKatalogu) =>
      pozycjaKatalogu.id === id ? potwierdzone.extension : pozycjaKatalogu,
    );
  }

  /** Liście jednej rodziny, pogrupowane po źródle pochodzenia. */
  function galezieRodziny(rodzaj: ExtensionKind): PozycjaMenu[] {
    const wykaz = (katalog ?? []).filter((rozszerzenie) => rozszerzenie.kind === rodzaj);
    const grupy = zrodlaNiepuste(wykaz);
    // Jedno źródło znaczy jedną grupę, a nagłówek nad całą zawartością gałęzi
    // niczego nie rozróżnia — grupowanie wchodzi dopiero, gdy jest co oddzielić.
    if (grupy.length < 2) return wykaz.map(lisc);
    return grupy.map((grupa) => ({
      rodzaj: 'grupa',
      nazwa: grupa.nazwa,
      dzieci: wykaz.filter((rozszerzenie) => rozszerzenie.origin === grupa.origin).map(lisc),
    }));
  }

  /** Źródła faktycznie obecne w tej rodzinie — puste grupy nie powstają. */
  function zrodlaNiepuste(wykaz: readonly Extension[]): typeof ZRODLA {
    return ZRODLA.filter((zrodlo) =>
      wykaz.some((rozszerzenie) => rozszerzenie.origin === zrodlo.origin),
    );
  }

  function lisc(rozszerzenie: Extension): PozycjaMenu {
    const opis = (rozszerzenie.description ?? '').trim();
    return {
      rodzaj: 'przelacznik',
      klucz: `${PRZEDROSTEK}${rozszerzenie.id}`,
      nazwa: rozszerzenie.name,
      ...(opis === '' ? {} : { opis }),
      wlaczony: rozszerzenie.enabled,
    };
  }

  /** Drzewo: cztery rodziny, w każdej grupy po źródle i liście-przełączniki. */
  function drzewo(): PozycjaMenu[] {
    const pozycje: PozycjaMenu[] = [];
    for (const rodzina of RODZINY) {
      const dzieci = galezieRodziny(rodzina.rodzaj);
      // Rodzina bez ani jednej pozycji nie staje w menu: gałąź otwierająca się
      // w pustkę jest ślepym zaułkiem, a nie ujawnianiem stopniowym.
      if (dzieci.length === 0) continue;
      pozycje.push({
        rodzaj: 'galaz',
        klucz: `rodzina:${rodzina.rodzaj}`,
        nazwa: rodzina.nazwa,
        opis: rodzina.opis,
        dzieci,
      });
    }
    return pozycje;
  }

  /**
   * Wartość na uchwycie: ile rozszerzeń pracuje.
   *
   * Nie „+", bo znak działania nie jest wartością nastawy, i nie nazwa jednego
   * z nich, bo pracuje ich naraz kilkanaście. Liczba włączonych mieści się
   * w jednym rzędzie i nie przekłamuje stanu.
   */
  function wartosc(): string {
    if (katalog === null) return WYKAZ_CZYTANY;
    const wlaczone = katalog.filter((rozszerzenie) => rozszerzenie.enabled).length;
    if (wlaczone === 0) return BRAK_ROZSZERZEN;
    return `${wlaczone} ${odmianaWlaczonych(wlaczone)}`;
  }

  function odswiez(): void {
    ster.ustaw(wartosc(), drzewo());
  }

  /** Odczyt katalogu; odmowa wychodzi zdaniem pod uchwytem, nie pustym wykazem. */
  function wczytaj(): void {
    void zaleznosci.zrodlo.katalog({ tylkoZainstalowane: true }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        // Katalog zostaje `null`, a nie pustą tablicą: „nie wiem" i „nie masz
        // ani jednego" to dwie różne prawdy, a menu-drzewo na pustej tablicy
        // powiedziałoby tę drugą.
        ster.zdanie(opisOdmowyBledu('Odczyt katalogu rozszerzeń (extension.list)', wynik.blad));
        odswiez();
        return;
      }
      katalog = wynik.wynik.extensions;
      ster.zdanie('');
      odswiez();
    });
  }

  const odsubskrybuj = zaleznosci.naZmiane(() => wczytaj());

  wczytaj();
  odswiez();

  return {
    element: ster.element,
    odswiez,
    rozlacz: () => odsubskrybuj(),
  };
}

/**
 * Odmiana słowa „włączone" przy liczebniku.
 *
 * Polszczyzna ma tu trzy formy: „1 włączone", „2 włączone", „5 włączonych".
 */
function odmianaWlaczonych(ile: number): string {
  const dziesiatki = ile % 100;
  const jednosci = ile % 10;
  if (dziesiatki >= 12 && dziesiatki <= 14) return 'włączonych';
  return jednosci >= 2 && jednosci <= 4 ? 'włączone' : ile === 1 ? 'włączone' : 'włączonych';
}
