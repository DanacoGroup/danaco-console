import { Command, type Action, type StudioOperation } from '../../../../shared/contract';
import { utworzMenuDrzewo, type GalazMenu, type MenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';
import { KATEGORIE_OPERACJI, nazwaOperacji } from './kategorie-operacji';
import { WIELKOSCI_CIAGLE } from './suwaki-koncepcyjne';

/**
 * Katalog operacji pod jednym uchwytem — miejsce, w którym Operator znajduje
 * czynność, której nie zna na pamięć.
 *
 * ── Trzy pochodzenia w jednym wykazie ───────────────────────────────────────
 * Rozstrzygnięcie Właściciela: **wykaz jest jeden, nie dwa** — ten sam dla
 * Operatora i dla modelu. Pozycje mają jednak trzy pochodzenia i katalog je
 * rozróżnia, bo Operator ma wiedzieć, czym rozporządza:
 *
 *   — **rejestr akcji rdzenia** (`action.list`) — pozycje wskazujące komendę
 *     `studio.contextual.op`; źródło właściwe: nowa operacja to wiersz rejestru;
 *   — **operacje zapisane w rdzeniu** (`studio.operation.list`) — fabryczne
 *     i własne Operatora, rozdzielone polem `builtin`;
 *   — **wykaz dokumentacji** (`kategorie-operacji.ts`) — 28 czynności
 *     w siedmiu grupach, dopóki dwa pierwsze źródła ich nie przejmą. Pozycja
 *     przejęta znika z wykazu dokumentacji, żeby nie stała w menu dwa razy.
 *
 * ── Czego katalog NIE odtwarza ──────────────────────────────────────────────
 * Rozwijania, wyszukiwania po nazwie, opisów przy pozycjach, wędrówki
 * strzałkami i znacznika wyboru na gałęzi. To wszystko niesie
 * `komponenty/menu-drzewo.ts` i stąd bierze się dostępność z klawiatury: cały
 * katalog jest osiągalny bez myszki, bo mechanizm menu obsługuje strzałki,
 * Enter i Escape sam. Drugi mechanizm rozwijania byłby rozjazdem.
 *
 * ── Wielkości ciągłe są tu oznaczone, nie wyjęte ────────────────────────────
 * Czynność będąca wielkością ciągłą (objętość, ton, rejestr, poziom szczegółu,
 * stopień dopracowania) zostaje w katalogu — ale jej opis mówi, że sterowanie
 * ma suwakiem, nie jednym naciśnięciem, i wskazuje pływak. Wyjęcie jej z menu
 * kazałoby Operatorowi szukać jej w dwóch miejscach.
 */

/** Czynności katalogu zlecane oknu. */
export interface CzynnosciKatalogu {
  /** Uruchamia operację o wskazanym identyfikatorze. */
  naOperacje(idAkcji: string): void;
  /** Zapisuje operację własną Operatora. */
  naZapisOperacji(nazwa: string, kategoria: string, polecenie: string): void;
  /** Usuwa operację; odmowę dla fabrycznej oddaje rdzeń. */
  naUsuniecieOperacji(idOperacji: string): void;
}

/** Katalog wraz z jego zasilaniem. */
export interface KatalogOperacji {
  element: HTMLElement;
  /** Menu — wołający może je rozwinąć albo zwinąć z zewnątrz. */
  menu: MenuDrzewo;
  /** Wstawia wiersze rejestru akcji rdzenia. */
  ustawRejestr(akcje: readonly Action[]): void;
  /** Wstawia operacje zapisane w rdzeniu — fabryczne i własne. */
  ustawOperacje(operacje: readonly StudioOperation[]): void;
  /** Wypisuje odpowiedź rdzenia na zapis albo usunięcie operacji własnej. */
  pokazOdpowiedz(tresc: string, udana: boolean): void;
  /** Zwija menu i zdejmuje jego nasłuchy dokumentu. */
  zamknij(): void;
}

/** Identyfikatory czynności sterowanych suwakiem — do opisu w menu. */
const STEROWANE_SUWAKIEM: ReadonlySet<string> = new Set(
  WIELKOSCI_CIAGLE.flatMap((wielkosc) => [wielkosc.idAkcji, `${wielkosc.kod}`]),
);

/** Zdanie dopisywane do opisu czynności, która jest wielkością ciągłą. */
const ZDANIE_SUWAKA =
  'Ta czynność jest wielkością ciągłą: na pływaku przy zaznaczeniu ma suwak, którym wskazujesz, ' +
  'o ile ma się zmienić. Uruchomienie z menu pojedzie z nastawą suwaka bieżącą.';

export function utworzKatalogOperacji(czynnosci: CzynnosciKatalogu): KatalogOperacji {
  let zRejestru: readonly Action[] = [];
  let zRdzenia: readonly StudioOperation[] = [];

  const menu = utworzMenuDrzewo({
    nastawa: 'Operacja kontekstowa',
    progSzukania: 8,
    naWybor(klucz) {
      if (klucz.startsWith('usun:')) {
        czynnosci.naUsuniecieOperacji(klucz.slice('usun:'.length));
        return;
      }
      czynnosci.naOperacje(klucz);
    },
  });

  const nazwa = document.createElement('input');
  nazwa.type = 'text';
  nazwa.className = 'dn-pole-kontrolka';
  nazwa.placeholder = 'nazwa operacji własnej';
  nazwa.setAttribute('aria-label', 'Nazwa operacji własnej Operatora');

  const kategoria = document.createElement('input');
  kategoria.type = 'text';
  kategoria.className = 'dn-pole-kontrolka';
  kategoria.placeholder = 'kategoria, w której ma stanąć';
  kategoria.setAttribute('aria-label', 'Kategoria operacji własnej');
  kategoria.setAttribute('list', 'ms-katalog-kategorie');

  const podpowiedzi = document.createElement('datalist');
  podpowiedzi.id = 'ms-katalog-kategorie';
  for (const pozycja of KATEGORIE_OPERACJI) {
    const wpis = document.createElement('option');
    wpis.value = pozycja.nazwa;
    podpowiedzi.append(wpis);
  }

  const polecenie = document.createElement('textarea');
  polecenie.className = 'dn-pole-kontrolka ms-katalog__polecenie';
  polecenie.rows = 3;
  polecenie.placeholder = 'treść polecenia dla modelu — to ona jedzie przy każdym uruchomieniu';
  polecenie.setAttribute('aria-label', 'Treść polecenia operacji własnej');

  const zapisz = document.createElement('button');
  zapisz.type = 'button';
  zapisz.className = 'dn-btn dn-btn--sm dn-btn--atrament';
  zapisz.textContent = 'Zapisz operację własną';
  zapisz.dataset['czynnosc'] = 'zapisz-operacje';
  zapisz.title =
    'Idzie komendą studio.operation.save. Operacja własna wchodzi do TEGO SAMEGO wykazu co ' +
    'fabryczne i jest w nim odróżniona pochodzeniem.';
  zapisz.addEventListener('click', () => {
    czynnosci.naZapisOperacji(nazwa.value.trim(), kategoria.value.trim(), polecenie.value.trim());
  });

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'dn-pole-opis ms-katalog__odpowiedz';

  const wlasne = document.createElement('details');
  wlasne.className = 'ms-katalog__wlasne';
  const podpis = document.createElement('summary');
  podpis.textContent = 'Operacja własna Operatora';
  wlasne.append(podpis, nazwa, kategoria, podpowiedzi, polecenie, zapisz, odpowiedz);

  const element = document.createElement('div');
  element.className = 'ms-katalog';
  element.append(menu.element, wlasne);

  function przerysuj(): void {
    menu.ustaw('cały katalog operacji ▾', przybornikZlozDrzewo(zRejestru, zRdzenia));
  }

  przerysuj();

  return {
    element,
    menu,

    ustawRejestr(akcje) {
      zRejestru = akcje;
      przerysuj();
    },

    ustawOperacje(operacje) {
      zRdzenia = operacje;
      przerysuj();
    },

    pokazOdpowiedz(tresc, udana) {
      odpowiedz.textContent = tresc;
      odpowiedz.dataset['udana'] = udana ? 'tak' : 'nie';
    },

    zamknij: () => menu.zwin(),
  };
}

/**
 * Składa drzewo katalogu z trzech pochodzeń.
 *
 * Plik oddaje je osobną funkcją, bez DOM, żeby złożenie sprawdzało się bez
 * stawiania menu: reguła „pozycja przejęta przez rdzeń znika z wykazu
 * dokumentacji" jest tu jedyną nieoczywistą rzeczą i ma dać się zmierzyć.
 */
export function przybornikZlozDrzewo(
  zRejestru: readonly Action[],
  zRdzenia: readonly StudioOperation[],
): PozycjaMenu[] {
  const kontekstowe = zRejestru.filter((akcja) => akcja.command === Command.StudioContextualOp);
  const przejete = new Set<string>([
    ...kontekstowe.map((akcja) => akcja.id),
    ...zRdzenia.map((operacja) => operacja.id),
  ]);

  const drzewo: PozycjaMenu[] = [];

  if (kontekstowe.length > 0) {
    drzewo.push({
      rodzaj: 'grupa',
      nazwa: 'Rejestr akcji rdzenia',
      dzieci: kontekstowe.map((akcja) => ({
        rodzaj: 'wybor' as const,
        klucz: akcja.id,
        nazwa: akcja.name,
        nazwaKrotka: akcja.id,
        opis: przybornikOpisPozycji(akcja.id, akcja.description),
        wybrany: false,
      })),
    });
  }

  const wlasne = zRdzenia.filter((operacja) => !operacja.builtin);
  const fabryczne = zRdzenia.filter((operacja) => operacja.builtin);
  if (fabryczne.length > 0) {
    drzewo.push({
      rodzaj: 'grupa',
      nazwa: 'Operacje fabryczne rdzenia',
      dzieci: przybornikGaleziePoKategorii(fabryczne, false),
    });
  }
  if (wlasne.length > 0) {
    drzewo.push({
      rodzaj: 'grupa',
      nazwa: 'Operacje własne Operatora',
      dzieci: przybornikGaleziePoKategorii(wlasne, true),
    });
  }

  const dokumentacja: GalazMenu[] = [];
  for (const grupa of KATEGORIE_OPERACJI) {
    const czekajace = grupa.operacje.filter((operacja) => !przejete.has(operacja.id));
    if (czekajace.length === 0) continue;
    dokumentacja.push({
      rodzaj: 'galaz',
      klucz: `kategoria:${grupa.kod}`,
      nazwa: grupa.nazwa,
      dzieci: czekajace.map((operacja) => ({
        rodzaj: 'wybor' as const,
        klucz: operacja.id,
        nazwa: operacja.nazwa,
        nazwaKrotka: operacja.id,
        opis: przybornikOpisPozycji(operacja.id, undefined),
        wybrany: false,
      })),
    });
  }
  if (dokumentacja.length > 0) {
    drzewo.push({
      rodzaj: 'grupa',
      nazwa:
        przejete.size === 0
          ? 'Siedem grup wykazu dokumentacji'
          : `Wykaz dokumentacji — bez zapisu w rdzeniu (rdzeń przejął już: ${przejete.size})`,
      dzieci: dokumentacja,
    });
  }

  return drzewo;
}

/** Gałęzie operacji rdzenia ułożone wedle ich kategorii. */
function przybornikGaleziePoKategorii(
  operacje: readonly StudioOperation[],
  zUsuwaniem: boolean,
): GalazMenu[] {
  const wedlugKategorii = new Map<string, StudioOperation[]>();
  for (const operacja of operacje) {
    const zebrane = wedlugKategorii.get(operacja.category) ?? [];
    zebrane.push(operacja);
    wedlugKategorii.set(operacja.category, zebrane);
  }
  return [...wedlugKategorii.entries()].map(([nazwaKategorii, pozycje]) => ({
    rodzaj: 'galaz' as const,
    klucz: `rdzen:${nazwaKategorii}`,
    nazwa: nazwaKategorii,
    dzieci: pozycje.flatMap((operacja) => {
      const lisc = {
        rodzaj: 'wybor' as const,
        klucz: operacja.id,
        nazwa: operacja.name,
        nazwaKrotka: operacja.id,
        opis: przybornikOpisPozycji(operacja.id, operacja.prompt),
        wybrany: false,
      };
      if (!zUsuwaniem) return [lisc];
      return [
        lisc,
        {
          rodzaj: 'wybor' as const,
          klucz: `usun:${operacja.id}`,
          nazwa: `Usuń „${operacja.name}"`,
          opis:
            'Idzie komendą studio.operation.delete. Operacji fabrycznej rdzeń nie usunie — ' +
            'odpowie odmową nazywającą powód, i tak ma zostać.',
          wybrany: false,
        },
      ];
    }),
  }));
}

/**
 * Opis pozycji katalogu.
 *
 * Bierze opis z rdzenia, gdy jest, a gdy go nie ma — nazywa czynność wykazem
 * dokumentacji. Czynność sterowana suwakiem dostaje przy tym zdanie o suwaku,
 * żeby Operator nie szukał go w drugim miejscu.
 */
export function przybornikOpisPozycji(idAkcji: string, opisRdzenia: string | undefined): string {
  const czesci: string[] = [];
  if (opisRdzenia !== undefined && opisRdzenia !== '') czesci.push(opisRdzenia);
  else {
    const nazwa = nazwaOperacji(idAkcji);
    if (nazwa !== undefined) czesci.push(nazwa);
  }
  if (STEROWANE_SUWAKIEM.has(idAkcji)) czesci.push(ZDANIE_SUWAKA);
  return czesci.join(' ');
}
