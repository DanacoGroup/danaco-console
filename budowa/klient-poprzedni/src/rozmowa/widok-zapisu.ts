import type { WpisRozmowy } from './wpis-rozmowy';

/**
 * Widok transkryptu — cztery tryby pokazywania wątku, rozstrzygające wyłącznie, które
 * warstwy wpisu się pokazują.
 */
export type WidokZapisu = 'zwykly' | 'rozumowanie' | 'pelny' | 'streszczenie';

/**
 * Pozycja wykazu trybów: klucz techniczny, nazwa widoczna dla Operatora oraz zdanie o
 * przeznaczeniu trybu.
 */
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

/**
 * Które konkretnie warstwy wpisu rozmowy rysuje dany tryb widoku transkryptu w tym oknie
 * komunikacji z rdzeniem.
 */
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
 * Rozkład warstw dla trybu; stopka i błędy stoją we wszystkich czterech trybach
 * niezależnie od ustawienia.
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
 * Czy wpis w ogóle wchodzi do zapisu w tym trybie widoku; filtruje wyłącznie tryb
 * streszczenia wątku.
 */
export function czyWpisWZapisie(widok: WidokZapisu, wpis: WpisRozmowy): boolean {
  if (widok !== 'streszczenie') return true;
  return wpis.podsumowanie !== null || wpis.narzedzia.length > 0 || wpis.bledy.length > 0;
}

/**
 * Rozpoznaje kod trybu widoku transkryptu; nieznany napis wraca trybem domyślnym, nigdy
 * błędem aplikacji.
 */
export function rozpoznajWidokZapisu(kod: string): WidokZapisu {
  const znany = WIDOKI_ZAPISU.find((pozycja) => pozycja.kod === kod);
  return znany === undefined ? WIDOK_ZAPISU_DOMYSLNY : znany.kod;
}
