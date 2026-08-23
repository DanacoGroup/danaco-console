import type { WpisRozmowy } from './wpis-rozmowy';

/**
 * Widok transkryptu — cztery tryby pokazywania wątku.
 *
 * Klucz techniczny jest oddzielony od etykiety: kod trybu (`zwykly`,
 * `rozumowanie`, `pelny`, `streszczenie`) nie zmienia się z językiem
 * interfejsu, a polska nazwa i zdanie o przeznaczeniu stoją obok niego
 * w wykazie.
 *
 * Plik nie rysuje niczego i nie zna DOM-u — rozstrzyga wyłącznie, które
 * warstwy wpisu mają się pokazać. Rysowanie zostaje w `widok-wpisu.ts`.
 *
 * Tryb nie sięga do rdzenia: wszystkie cztery stoją na polach, które wpis już
 * niesie (`tresc`, `rozumowanie`, `narzedzia`, `prowenancja`, `podsumowanie`).
 * Przełączenie jest filtrem nad pamięcią okna — nie woła komendy, nie dociąga
 * historii i niczego nie gubi.
 */
export type WidokZapisu = 'zwykly' | 'rozumowanie' | 'pelny' | 'streszczenie';

/** Pozycja wykazu trybów: klucz techniczny, nazwa dla Operatora, przeznaczenie. */
export interface OpisWidokuZapisu {
  /** Klucz techniczny — ten sam we wszystkich językach interfejsu. */
  kod: WidokZapisu;
  /** Nazwa widoczna w menu sesji. */
  nazwa: string;
  /** Zdanie o tym, co tryb pokazuje. */
  przeznaczenie: string;
}

/**
 * Wykaz czterech trybów.
 *
 * Kolejność nie jest dowolna: idzie od najwęższego widoku do najszerszego,
 * a streszczenie stoi na końcu, bo jako jedyne nie dokłada warstw, tylko je
 * zdejmuje.
 */
export const WIDOKI_ZAPISU: readonly OpisWidokuZapisu[] = [
  {
    kod: 'zwykly',
    nazwa: 'Zwykły',
    przeznaczenie: 'Wypowiedzi; wywołania narzędzi zwinięte.',
  },
  {
    kod: 'rozumowanie',
    nazwa: 'Rozumowanie',
    przeznaczenie: 'Jak wyżej oraz bloki rozumowania.',
  },
  {
    kod: 'pelny',
    nazwa: 'Pełny',
    przeznaczenie: 'Wywołania narzędzi, argumenty, wyniki, zdarzenia.',
  },
  {
    kod: 'streszczenie',
    nazwa: 'Streszczenie',
    przeznaczenie: 'Streszczenia tur i wykaz wytworów.',
  },
];

/**
 * Tryb, w którym okno otwiera wątek.
 *
 * Wartość domyślna nie odkrywa toku rozumowania ani prowenancji — pokazanie
 * pracy modelu jest decyzją Operatora.
 */
export const WIDOK_ZAPISU_DOMYSLNY: WidokZapisu = 'zwykly';

/** Które warstwy wpisu rysuje dany tryb. */
export interface WarstwyZapisu {
  /** Blok „co poszło do modelu". */
  prowenancja: boolean;
  /** Blok toku rozumowania. */
  rozumowanie: boolean;
  /** Treść wypowiedzi. */
  tresc: boolean;
  /** Blok wywołań narzędzi. */
  narzedzia: boolean;
  /** Czy blok narzędzi ma być rozwinięty od razu — argumenty i wyniki na wierzchu. */
  narzedziaRozwiniete: boolean;
  /** Błędy tury. */
  bledy: boolean;
  /** Stopka: podsumowanie tury i konto kanału. */
  stopka: boolean;
  /** Wykaz typów zdarzeń przechwyconych w turze. */
  zdarzenia: boolean;
  /** Wykaz wytworów tury — nazwy narzędzi, które tura wywołała. */
  wytwory: boolean;
}

/**
 * Rozkład warstw dla trybu.
 *
 * Stopka stoi we wszystkich czterech trybach: podsumowanie tury i konto kanału
 * są widoczne bez przełącznika, a tryb ma zapis poszerzać albo zawężać do
 * streszczenia, nie odbierać informacji dostępnej wszędzie indziej.
 *
 * Błędy również stoją we wszystkich czterech. Błąd kanału nie jest szczegółem
 * diagnostycznym, który wolno schować za trybem: wpis, w którym tura padła,
 * ma o tym mówić niezależnie od ustawienia widoku.
 */
export function warstwyZapisu(widok: WidokZapisu): WarstwyZapisu {
  const pelny = widok === 'pelny';
  const streszczenie = widok === 'streszczenie';
  return {
    prowenancja: pelny,
    rozumowanie: widok === 'rozumowanie' || pelny,
    tresc: !streszczenie,
    narzedzia: !streszczenie,
    narzedziaRozwiniete: pelny,
    bledy: true,
    stopka: true,
    zdarzenia: pelny,
    wytwory: streszczenie,
  };
}

/**
 * Czy wpis w ogóle wchodzi do zapisu w tym trybie.
 *
 * Filtruje wyłącznie tryb `streszczenie`: wpis, z którego nie da się nic
 * streścić — bez domkniętej tury, bez wywołań narzędzi i bez błędu —
 * zostawiłby na ekranie samą ramkę z godziną. Trzy pozostałe tryby nie chowają
 * niczego.
 *
 * Funkcja nie usuwa wpisu z pamięci okna ani z listy; ukrycie jest odwracalne
 * powrotem do innego trybu.
 */
export function czyWpisWZapisie(widok: WidokZapisu, wpis: WpisRozmowy): boolean {
  if (widok !== 'streszczenie') return true;
  return wpis.podsumowanie !== null || wpis.narzedzia.length > 0 || wpis.bledy.length > 0;
}

/** Rozpoznaje kod trybu; nieznany napis wraca trybem domyślnym, nigdy błędem. */
export function rozpoznajWidokZapisu(kod: string): WidokZapisu {
  const znany = WIDOKI_ZAPISU.find((pozycja) => pozycja.kod === kod);
  return znany === undefined ? WIDOK_ZAPISU_DOMYSLNY : znany.kod;
}
