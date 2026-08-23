/**
 * Suwaki koncepcyjne — wielkości ciągłe pracy nad tekstem.
 *
 * ── Dlaczego suwak, a nie przycisk ──────────────────────────────────────────
 * „Skróć" i „Rozwiń" to nie dwie czynności, tylko dwa końce jednej wielkości:
 * objętości. To samo z tonem, rejestrem, poziomem szczegółu i stopniem
 * dopracowania. Przycisk daje jeden skok o nieznanej wielkości i drugie
 * naciśnięcie skacze znowu; suwak nazywa, o ile ma się zmienić, i pokazuje to
 * przed wysłaniem.
 *
 * ── Skąd biorą się operacje ─────────────────────────────────────────────────
 * Z `kategorie-operacji.ts`, nie z nowego wykazu. Każdy suwak wskazuje
 * identyfikator akcji, który już tam stoi, i dokłada mu wielkość. Operacje,
 * które wielkością ciągłą nie są (korekta interpunkcji, spis treści,
 * tłumaczenie), zostają przyciskami w swoich miejscach.
 *
 * ── Czym jedzie nastawa ─────────────────────────────────────────────────────
 * Polem `params` żądania `studio.contextual.op` — kontrakt ma je jako `unknown`
 * i opisuje jako „parametry operacji wymagane przez pozycję rejestru". Nastawy
 * jadą więc drogą, która w kontrakcie jest, i nie wymagają jego zmiany. Rdzeń
 * dziś `params` do polecenia modelu nie dokłada — to jest potrzeba nazwana
 * w sprawozdaniu, nie brak ukryty: okno pisze przy suwakach, że wielkość jedzie
 * w żądaniu, i dokłada ją TAKŻE do treści polecenia wysyłanego wierszem
 * polecenia, gdzie model ją przeczyta.
 *
 * Plik nie zna DOM — oddaje wykaz, zdania i ładunek żądania.
 */

/** Jedna wielkość ciągła: identyfikator operacji, końce skali i jej nazwa. */
export interface WielkoscCiagla {
  /** Kod nastawy — trafia do `params` i do `data-suwak`. */
  kod: string;
  nazwa: string;
  /** Identyfikator akcji z `kategorie-operacji.ts`. */
  idAkcji: string;
  /** Nazwa końca dolnego skali. */
  koniecDolny: string;
  /** Nazwa końca górnego skali. */
  koniecGorny: string;
  /** Wartość, przy której nastawa niczego nie zmienia. */
  neutralna: number;
  /** Zdanie mówiące, co wielkość robi z tekstem. */
  opis: string;
}

/**
 * Granice skali wszystkich suwaków.
 *
 * Skala jest jedna dla wszystkich wielkości i symetryczna wobec zera, bo każdy
 * suwak ma dwa przeciwne kierunki i środek, w którym nie robi nic. Osobne skale
 * dla każdej wielkości kazałyby Operatorowi czytać liczbę inaczej przy każdym
 * suwaku.
 */
export const SKALA_SUWAKA = { dol: -3, gora: 3, krok: 1 } as const;

export const WIELKOSCI_CIAGLE: readonly WielkoscCiagla[] = [
  {
    kod: 'objetosc',
    nazwa: 'Objętość',
    idAkcji: 'studio.styl.skrocenie',
    koniecDolny: 'skróć mocno',
    koniecGorny: 'rozwiń mocno',
    neutralna: 0,
    opis:
      'Ile tekstu ma zostać. W lewo zwięźlej — do samej treści zdania; w prawo szerzej — ' +
      'z rozwinięciem i przykładami. Operacja idzie identyfikatorem studio.styl.skrocenie ' +
      'albo studio.styl.rozwiniecie, zależnie od kierunku.',
  },
  {
    kod: 'ton',
    nazwa: 'Ton',
    idAkcji: 'studio.styl.ton',
    koniecDolny: 'rzeczowo',
    koniecGorny: 'perswazyjnie',
    neutralna: 0,
    opis: 'Nastawienie wypowiedzi wobec czytającego: od suchego sprawozdania do namowy.',
  },
  {
    kod: 'rejestr',
    nazwa: 'Rejestr',
    idAkcji: 'studio.styl.rejestr',
    koniecDolny: 'potocznie',
    koniecGorny: 'urzędowo',
    neutralna: 0,
    opis: 'Wysokość stylu: od mowy codziennej do pisma urzędowego.',
  },
  {
    kod: 'szczegol',
    nazwa: 'Poziom szczegółu',
    idAkcji: 'studio.styl.uproszczenie',
    koniecDolny: 'prościej',
    koniecGorny: 'dokładniej',
    neutralna: 0,
    opis:
      'Ile wiedzy tekst zakłada u czytającego. W lewo uproszczenie języka, w prawo ' +
      'terminologia i dokładność zapisu.',
  },
  {
    kod: 'dopracowanie',
    nazwa: 'Stopień dopracowania',
    idAkcji: 'studio.korekta.czytelnosc',
    koniecDolny: 'sama korekta',
    koniecGorny: 'przepisanie',
    neutralna: 0,
    opis:
      'Jak głęboko model ma wejść w tekst: od poprawienia błędów bez ruszania zdań ' +
      'do przepisania fragmentu od nowa.',
  },
];

/** Nastawy wszystkich suwaków — kod wielkości na jej wartość. */
export type NastawySuwakow = Record<string, number>;

/** Nastawy neutralne: każdy suwak w środku skali. */
export function neutralneNastawy(): NastawySuwakow {
  const nastawy: NastawySuwakow = {};
  for (const wielkosc of WIELKOSCI_CIAGLE) nastawy[wielkosc.kod] = wielkosc.neutralna;
  return nastawy;
}

/**
 * Identyfikator akcji dla wielkości ustawionej w danym kierunku.
 *
 * Objętość ma dwie operacje przeciwne w `kategorie-operacji.ts` — skrócenie
 * i rozwinięcie — więc kierunek suwaka wskazuje, która z nich pojedzie. Reszta
 * wielkości ma jedną operację i kierunek jedzie w nastawie.
 */
export function akcjaWielkosci(wielkosc: WielkoscCiagla, wartosc: number): string {
  if (wielkosc.kod === 'objetosc' && wartosc > 0) return 'studio.styl.rozwiniecie';
  if (wielkosc.kod === 'szczegol' && wartosc > 0) return 'studio.wzbogacenie.kontekst';
  if (wielkosc.kod === 'dopracowanie' && wartosc > 0) return 'studio.styl.rejestr';
  return wielkosc.idAkcji;
}

/**
 * Zdanie o skutku nastawy — podgląd przed wysłaniem.
 *
 * Podgląd nie udaje wyniku modelu: mówi, o ile nastawa przestawia wielkość
 * i w którą stronę. Pokazywanie tu zmyślonego zdania „tak będzie wyglądał
 * tekst" byłoby wynikiem bez rachunku.
 */
export function opiszNastawe(wielkosc: WielkoscCiagla, wartosc: number): string {
  if (wartosc === wielkosc.neutralna) {
    return `${wielkosc.nazwa}: bez zmiany — suwak stoi w środku, więc ta wielkość nie pojedzie w żądaniu.`;
  }
  const kierunek = wartosc < 0 ? wielkosc.koniecDolny : wielkosc.koniecGorny;
  const stopien = Math.abs(wartosc);
  const nasilenie = stopien === 1 ? 'lekko' : stopien === 2 ? 'wyraźnie' : 'mocno';
  return `${wielkosc.nazwa}: ${nasilenie} w stronę „${kierunek}" (${wartosc > 0 ? '+' : ''}${wartosc} z ${SKALA_SUWAKA.gora}) · akcja ${akcjaWielkosci(wielkosc, wartosc)}.`;
}

/** Nastawy różne od neutralnych — tylko one jadą w żądaniu. */
export function nastawyCzynne(nastawy: NastawySuwakow): { wielkosc: WielkoscCiagla; wartosc: number }[] {
  const czynne: { wielkosc: WielkoscCiagla; wartosc: number }[] = [];
  for (const wielkosc of WIELKOSCI_CIAGLE) {
    const wartosc = nastawy[wielkosc.kod] ?? wielkosc.neutralna;
    if (wartosc !== wielkosc.neutralna) czynne.push({ wielkosc, wartosc });
  }
  return czynne;
}

/**
 * Ładunek `params` żądania operacji kontekstowej.
 *
 * Niesie nastawy czynne wraz z granicami skali, bo bez granic „ton: 2" nic nie
 * znaczy. Polecenie słowne Operatora jedzie tym samym polem, gdy je podał.
 */
export function ladunekOperacji(
  nastawy: NastawySuwakow,
  polecenie: string,
): Record<string, unknown> {
  const czynne = nastawyCzynne(nastawy);
  const ladunek: Record<string, unknown> = {
    skala: { dol: SKALA_SUWAKA.dol, gora: SKALA_SUWAKA.gora },
  };
  if (polecenie.trim() !== '') ladunek['polecenie'] = polecenie.trim();
  for (const pozycja of czynne) ladunek[pozycja.wielkosc.kod] = pozycja.wartosc;
  return ladunek;
}

/**
 * Zdanie polecenia dla modelu składane z nastaw.
 *
 * Rdzeń nie dokłada dziś `params` do treści polecenia (`trescOperacjiStudia`
 * składa czynność i treść dokumentu), więc nastawy muszą pojechać także słowem —
 * inaczej suwak przestawiałby pole, którego model nie czyta. Zdanie jest jedno
 * i powstaje tutaj, żeby wiersz polecenia i wstążka mówiły modelowi to samo.
 */
export function zdanieDlaModelu(nastawy: NastawySuwakow, polecenie: string): string {
  const czesci: string[] = [];
  if (polecenie.trim() !== '') czesci.push(polecenie.trim());
  for (const pozycja of nastawyCzynne(nastawy)) {
    const kierunek =
      pozycja.wartosc < 0 ? pozycja.wielkosc.koniecDolny : pozycja.wielkosc.koniecGorny;
    czesci.push(
      `${pozycja.wielkosc.nazwa}: ${kierunek}, natężenie ${Math.abs(pozycja.wartosc)} z ${SKALA_SUWAKA.gora}`,
    );
  }
  return czesci.join('; ');
}
