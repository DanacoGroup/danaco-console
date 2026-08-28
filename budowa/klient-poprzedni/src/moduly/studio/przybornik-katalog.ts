import { Command, type Action, type StudioOperation } from '../../../../shared/contract';
import { utworzMenuDrzewo, type GalazMenu, type MenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';
import { KATEGORIE_OPERACJI, nazwaOperacji } from './kategorie-operacji';
import { WIELKOSCI_CIAGLE } from './suwaki-koncepcyjne';

/** Interfejs CzynnosciKatalogu niesie czynności katalogu operacji zlecane oknu: uruchomienie, zapis operacji własnej i jej usunięcie. */
export interface CzynnosciKatalogu {
  /** Uruchamia operację o wskazanym identyfikatorze. */
  naOperacje(idAkcji: string): void;
  /** Zapisuje operację własną Operatora. */
  naZapisOperacji(nazwa: string, kategoria: string, polecenie: string): void;
  /** Usuwa operację; odmowę dla fabrycznej oddaje rdzeń. */
  naUsuniecieOperacji(idOperacji: string): void;
}

/** Interfejs KatalogOperacji niesie katalog operacji wraz z jego zasilaniem: elementem, menu, wstawianiem rejestru i operacji, odpowiedzią i zamknięciem. */
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

/** Stała STEROWANE_SUWAKIEM niesie identyfikatory czynności sterowanych suwakiem, do opisu tych czynności w menu katalogu. */
const STEROWANE_SUWAKIEM: ReadonlySet<string> = new Set(
  WIELKOSCI_CIAGLE.flatMap((wielkosc) => [wielkosc.idAkcji, `${wielkosc.kod}`]),
);

/** Stała ZDANIE_SUWAKA niesie zdanie dopisywane do opisu czynności katalogu, która jest wielkością ciągłą sterowaną suwakiem. */
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
 * Funkcja przybornikZlozDrzewo składa drzewo katalogu z trzech pochodzeń osobną funkcją, bez elementu menu, żeby reguła zaniku pozycji przejętej dała się zmierzyć.
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

/** Funkcja przybornikGaleziePoKategorii zwraca gałęzie operacji rdzenia ułożone wedle ich kategorii, gotowe do wstawienia w menu. */
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
