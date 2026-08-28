import { ExtensionKind, ExtensionOrigin, type Extension } from '../../../shared/contract';
import type { PozycjaMenu } from '../komponenty/menu-drzewo';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { ZrodloRozszerzen } from '../moduly/agents/zrodlo-rozszerzen';
import { utworzSterNastawy, type SterPaska } from './ster-nastawy';

// Rozszerzenia to ster paska zlecenia, wejście do czterech rodzin rozszerzeń zamiast jednej listy.

/** Nazwa rodzajowa nastawy rozszerzeń — idzie do etykiety dostępności aria-label, nigdy nie trafia na ekran. */
const NASTAWA = 'Rozszerzenia';

/** Przedrostek klucza liścia menu rozszerzeń; oddziela identyfikatory pozycji od kluczy gałęzi rodziny. */
const PRZEDROSTEK = 'rozszerzenie:';

/** Etykieta uchwytu przy pustym wykazie rozszerzeń — mówi wprost o braku, niczego nie udaje operatorowi. */
export const BRAK_ROZSZERZEN = 'Bez rozszerzeń';

/** Etykieta uchwytu pokazywana, dopóki katalog rozszerzeń nie wrócił jeszcze z odpowiedzią od samego rdzenia. */
export const WYKAZ_CZYTANY = 'Wczytuję…';

/** Nazwy rodzin rozszerzeń w kolejności wyświetlania — cztery gałęzie pierwszego poziomu drzewa tego menu. */
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

/** Nagłówki grup rozszerzeń pogrupowanych po źródle pochodzenia w drzewie menu paska zlecenia operatora. */
const ZRODLA: readonly { origin: ExtensionOrigin; nazwa: string }[] = [
  { origin: ExtensionOrigin.Danaco, nazwa: 'Danaco' },
  { origin: ExtensionOrigin.Personal, nazwa: 'Moje' },
];

/** Zależności steru rozszerzeń — wąskie i wstrzykiwane, obejmujące źródło komend oraz katalog rejestracji. */
export interface ZaleznosciSteruRozszerzen {
  /** Gotowe źródło komend `extension.*`; ster go nie zakłada. */
  zrodlo: ZrodloRozszerzen;
  /** Subskrypcja `extension.changed` — katalog zmienia się także spoza paska. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
  /** Droga do katalogu rozszerzeń bywa pominięta — stopka menu jej wtedy nie ma. */
  otworzKatalog?(): void;
}

/** Ster paska rozszerzeń ze zdejmowalnym nasłuchem zmiany katalogu, zgłaszanej przez zdarzenie rdzenia. */
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

  /** Stan przełączenia bierze się z katalogu, nie z menu, by nie cofnąć decyzji rdzenia w międzyczasie. */
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
    // Odpowiedź rdzenia jest stanem potwierdzonym, wpisywanym do katalogu od razu, bez pełnego odczytu.
    katalog = (katalog ?? []).map((pozycjaKatalogu) =>
      pozycjaKatalogu.id === id ? potwierdzone.extension : pozycjaKatalogu,
    );
  }

  /** Liście jednej rodziny, pogrupowane po źródle pochodzenia. */
  function galezieRodziny(rodzaj: ExtensionKind): PozycjaMenu[] {
    const wykaz = (katalog ?? []).filter((rozszerzenie) => rozszerzenie.kind === rodzaj);
    const grupy = zrodlaNiepuste(wykaz);
    // Jedno źródło znaczy jedną grupę — nagłówek nad całą zawartością gałęzi niczego nie rozróżnia.
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
      // Rodzina bez ani jednej pozycji nie staje w menu — pusta gałąź byłaby ślepym zaułkiem.
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

  /** Wartość na uchwycie to liczba pracujących rozszerzeń, nie znak działania ani nazwa jednego z nich. */
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
        // Katalog zostaje pusty jako nieznany, nie jako pusta tablica — to dwie różne prawdy dla menu.
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
