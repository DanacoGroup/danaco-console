/**
 * Moduł suwaki-koncepcyjne opisuje wielkości ciągłe pracy nad tekstem — objętość, ton, rejestr, poziom szczegółu i stopień dopracowania — jako suwaki zamiast przycisków, wraz z ładunkiem żądania operacji kontekstowej.
 */

/** Interfejs WielkoscCiagla opisuje jedną wielkość ciągłą pracy nad tekstem: identyfikator operacji, końce skali i jej nazwę widoczną w interfejsie. */
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
 * Stała SKALA_SUWAKA niesie granice skali wspólne dla wszystkich suwaków, symetryczne wobec zera, ponieważ każdy suwak ma dwa przeciwne kierunki i środek neutralny.
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

/** Typ NastawySuwakow mapuje kod każdej wielkości ciągłej na jej bieżącą wartość liczbową na skali suwaka. */
export type NastawySuwakow = Record<string, number>;

/** Funkcja neutralneNastawy zwraca nastawy początkowe, w których każdy suwak stoi dokładnie w środku skali. */
export function neutralneNastawy(): NastawySuwakow {
  const nastawy: NastawySuwakow = {};
  for (const wielkosc of WIELKOSCI_CIAGLE) nastawy[wielkosc.kod] = wielkosc.neutralna;
  return nastawy;
}

/**
 * Funkcja akcjaWielkosci zwraca identyfikator akcji dla wielkości ustawionej w danym kierunku; objętość ma dwie operacje przeciwne, a kierunek suwaka wskazuje, która z nich pojedzie.
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

/** Funkcja nastawyCzynne zwraca nastawy różne od wartości neutralnej — tylko one jadą w żądaniu operacji kontekstowej. */
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
 * Funkcja zdanieDlaModelu składa zdanie polecenia dla modelu z nastaw suwaków, ponieważ rdzeń nie dokłada dziś tych nastaw do treści polecenia samodzielnie.
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
