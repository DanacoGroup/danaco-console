/**
 * Kreator wyrażenia cyklicznego okna Scheduler — wzorce cykliczności i ich
 * przełożenie na notację cron.
 *
 * Operator ustala cykliczność wzorcem („co tydzień, poniedziałek, 07:00”),
 * a kontrakt niesie ją jednym polem `AutomationSchedule.cron`. Ten plik jest
 * przekładem w obie strony: składa zapis z nastaw i rozpoznaje wzorzec
 * w zapisie odczytanym z rdzenia, żeby okno otwarte na harmonogramie zastanym
 * pokazywało wzorzec, a nie samą składnię.
 *
 * Przekład jest zawężony do wzorców, które da się zapisać pięcioma polami
 * notacji: zapis odczytany z rdzenia i niepasujący do żadnego wzorca zostaje
 * wzorcem własnym, a jego treść idzie do pola zapisu bez zmiany. Rachunek
 * kolejnych terminów należy do `nastepne-uruchomienia.ts`.
 */

/**
 * Wzorce cykliczności udostępniane w oknie.
 *
 * Nazwy są kluczami wyboru w oknie, nie wartościami kontraktu: harmonogram
 * niesie gotowy zapis cron albo odstęp w sekundach, a nie wzorzec, z którego
 * jedno czy drugie powstało.
 */
export const WZORCE_CYKLICZNOSCI = {
  coMinute: 'co-minute',
  godzinnie: 'godzinnie',
  codziennie: 'codziennie',
  coTydzien: 'co-tydzien',
  coMiesiac: 'co-miesiac',
  kwartalnie: 'kwartalnie',
  interwalMinutowy: 'interwal-minutowy',
  interwalGodzinowy: 'interwal-godzinowy',
  wlasny: 'wlasny',
} as const;

export type WzorzecCyklicznosci =
  (typeof WZORCE_CYKLICZNOSCI)[keyof typeof WZORCE_CYKLICZNOSCI];

/** Wzorce w kolejności wykazu rozwijanego wraz z ich nazwami na ekranie. */
export const NAZWY_WZORCOW: ReadonlyArray<[WzorzecCyklicznosci, string]> = [
  [WZORCE_CYKLICZNOSCI.coMinute, 'co minutę'],
  [WZORCE_CYKLICZNOSCI.godzinnie, 'godzinnie'],
  [WZORCE_CYKLICZNOSCI.codziennie, 'codziennie'],
  [WZORCE_CYKLICZNOSCI.coTydzien, 'co tydzień'],
  [WZORCE_CYKLICZNOSCI.coMiesiac, 'co miesiąc'],
  [WZORCE_CYKLICZNOSCI.kwartalnie, 'kwartalnie'],
  [WZORCE_CYKLICZNOSCI.interwalMinutowy, 'co N minut'],
  [WZORCE_CYKLICZNOSCI.interwalGodzinowy, 'co N godzin'],
  [WZORCE_CYKLICZNOSCI.wlasny, 'niestandardowo'],
];

/** Dni tygodnia w numeracji notacji cron: niedziela ma numer zero. */
export const DNI_TYGODNIA: ReadonlyArray<[string, string]> = [
  ['1', 'poniedziałek'],
  ['2', 'wtorek'],
  ['3', 'środa'],
  ['4', 'czwartek'],
  ['5', 'piątek'],
  ['6', 'sobota'],
  ['0', 'niedziela'],
];

/**
 * Miesiące otwierające kwartały. Kwartalnie znaczy „raz na kwartał”, a notacja
 * cron nie ma pojęcia kwartału — ma wykaz miesięcy.
 */
const MIESIACE_KWARTALOW = '1,4,7,10';

/** Nastawy wzorca; pola nieużywane przez dany wzorzec zostają nietknięte. */
export interface NastawyWzorca {
  /** Minuta godziny, zakres 0–59. */
  minuta: number;
  /** Godzina doby, zakres 0–23. */
  godzina: number;
  /** Dzień tygodnia w numeracji cron, zakres 0–6. */
  dzienTygodnia: number;
  /** Dzień miesiąca, zakres 1–31. */
  dzienMiesiaca: number;
  /** Krok wzorca interwałowego — liczba minut albo godzin. */
  krok: number;
}

/** Nastawy wyjściowe okna: codziennie o siódmej, w poniedziałek, pierwszego. */
export const NASTAWY_WYJSCIOWE: NastawyWzorca = {
  minuta: 0,
  godzina: 7,
  dzienTygodnia: 1,
  dzienMiesiaca: 1,
  krok: 15,
};

/** Wzorzec rozpoznany w zapisie wraz z nastawami odczytanymi z jego pól. */
export interface RozpoznanyWzorzec {
  wzorzec: WzorzecCyklicznosci;
  nastawy: NastawyWzorca;
}

/**
 * Składa zapis cron z wzorca i nastaw.
 *
 * Wzorzec własny nie ma czego składać — treść zapisu podaje wtedy Operator
 * wprost, więc funkcja oddaje zapis zastany bez zmiany.
 */
export function zapisWzorca(
  wzorzec: WzorzecCyklicznosci,
  nastawy: NastawyWzorca,
  zapisWlasny = '',
): string {
  const minuta = ogranicz(nastawy.minuta, 0, 59);
  const godzina = ogranicz(nastawy.godzina, 0, 23);
  const dzienTygodnia = ogranicz(nastawy.dzienTygodnia, 0, 6);
  const dzienMiesiaca = ogranicz(nastawy.dzienMiesiaca, 1, 31);
  switch (wzorzec) {
    case WZORCE_CYKLICZNOSCI.coMinute:
      return '* * * * *';
    case WZORCE_CYKLICZNOSCI.godzinnie:
      return `${minuta} * * * *`;
    case WZORCE_CYKLICZNOSCI.codziennie:
      return `${minuta} ${godzina} * * *`;
    case WZORCE_CYKLICZNOSCI.coTydzien:
      return `${minuta} ${godzina} * * ${dzienTygodnia}`;
    case WZORCE_CYKLICZNOSCI.coMiesiac:
      return `${minuta} ${godzina} ${dzienMiesiaca} * *`;
    case WZORCE_CYKLICZNOSCI.kwartalnie:
      return `${minuta} ${godzina} ${dzienMiesiaca} ${MIESIACE_KWARTALOW} *`;
    case WZORCE_CYKLICZNOSCI.interwalMinutowy:
      return `*/${ogranicz(nastawy.krok, 1, 59)} * * * *`;
    case WZORCE_CYKLICZNOSCI.interwalGodzinowy:
      return `${minuta} */${ogranicz(nastawy.krok, 1, 23)} * * *`;
    default:
      return zapisWlasny;
  }
}

/**
 * Rozpoznaje wzorzec w zapisie odczytanym z rdzenia.
 *
 * Zapis, którego żaden wzorzec nie obejmuje, jest wzorcem własnym — nie
 * usterką. Nastawy nierozstrzygnięte przez pola zapisu zostają wyjściowe, żeby
 * przełączenie wzorca w oknie miało od czego zacząć.
 */
export function rozpoznajWzorzec(zapis: string): RozpoznanyWzorzec {
  const pola = zapis.trim().split(/\s+/);
  if (pola.length !== 5) return { wzorzec: WZORCE_CYKLICZNOSCI.wlasny, nastawy: { ...NASTAWY_WYJSCIOWE } };
  const [minuta = '', godzina = '', dzienMiesiaca = '', miesiac = '', dzienTygodnia = ''] = pola;
  const nastawy: NastawyWzorca = {
    minuta: liczbaPola(minuta, NASTAWY_WYJSCIOWE.minuta),
    godzina: liczbaPola(godzina, NASTAWY_WYJSCIOWE.godzina),
    dzienTygodnia: liczbaPola(dzienTygodnia, NASTAWY_WYJSCIOWE.dzienTygodnia),
    dzienMiesiaca: liczbaPola(dzienMiesiaca, NASTAWY_WYJSCIOWE.dzienMiesiaca),
    krok: NASTAWY_WYJSCIOWE.krok,
  };

  const krokMinut = krokPola(minuta);
  if (krokMinut !== null && godzina === '*' && dzienMiesiaca === '*' && miesiac === '*' && dzienTygodnia === '*') {
    return { wzorzec: WZORCE_CYKLICZNOSCI.interwalMinutowy, nastawy: { ...nastawy, krok: krokMinut } };
  }
  const krokGodzin = krokPola(godzina);
  if (krokGodzin !== null && dzienMiesiaca === '*' && miesiac === '*' && dzienTygodnia === '*') {
    return { wzorzec: WZORCE_CYKLICZNOSCI.interwalGodzinowy, nastawy: { ...nastawy, krok: krokGodzin } };
  }
  if (minuta === '*' && godzina === '*' && dzienMiesiaca === '*' && miesiac === '*' && dzienTygodnia === '*') {
    return { wzorzec: WZORCE_CYKLICZNOSCI.coMinute, nastawy };
  }
  if (liczbowe(minuta) && godzina === '*' && dzienMiesiaca === '*' && miesiac === '*' && dzienTygodnia === '*') {
    return { wzorzec: WZORCE_CYKLICZNOSCI.godzinnie, nastawy };
  }
  if (liczbowe(minuta) && liczbowe(godzina) && dzienMiesiaca === '*' && miesiac === '*' && dzienTygodnia === '*') {
    return { wzorzec: WZORCE_CYKLICZNOSCI.codziennie, nastawy };
  }
  if (liczbowe(minuta) && liczbowe(godzina) && dzienMiesiaca === '*' && miesiac === '*' && liczbowe(dzienTygodnia)) {
    return { wzorzec: WZORCE_CYKLICZNOSCI.coTydzien, nastawy };
  }
  if (liczbowe(minuta) && liczbowe(godzina) && liczbowe(dzienMiesiaca) && dzienTygodnia === '*') {
    if (miesiac === '*') return { wzorzec: WZORCE_CYKLICZNOSCI.coMiesiac, nastawy };
    if (miesiac === MIESIACE_KWARTALOW) return { wzorzec: WZORCE_CYKLICZNOSCI.kwartalnie, nastawy };
  }
  return { wzorzec: WZORCE_CYKLICZNOSCI.wlasny, nastawy };
}

/** Zdanie o cykliczności w mowie Operatora — pod polem zapisu w oknie. */
export function opisCyklicznosci(zapis: string): string {
  const { wzorzec, nastawy } = rozpoznajWzorzec(zapis);
  const godzina = `${dwucyfrowo(nastawy.godzina)}:${dwucyfrowo(nastawy.minuta)}`;
  switch (wzorzec) {
    case WZORCE_CYKLICZNOSCI.coMinute:
      return 'Uruchomienie co minutę, bez przerwy.';
    case WZORCE_CYKLICZNOSCI.godzinnie:
      return `Uruchomienie co godzinę, w ${dwucyfrowo(nastawy.minuta)} minucie.`;
    case WZORCE_CYKLICZNOSCI.codziennie:
      return `Uruchomienie codziennie o ${godzina}.`;
    case WZORCE_CYKLICZNOSCI.coTydzien:
      return `Uruchomienie co tydzień w ${nazwaDnia(nastawy.dzienTygodnia)} o ${godzina}.`;
    case WZORCE_CYKLICZNOSCI.coMiesiac:
      return `Uruchomienie co miesiąc, ${nastawy.dzienMiesiaca} dnia o ${godzina}.`;
    case WZORCE_CYKLICZNOSCI.kwartalnie:
      return `Uruchomienie raz na kwartał: w styczniu, kwietniu, lipcu i październiku, ${nastawy.dzienMiesiaca} dnia o ${godzina}.`;
    case WZORCE_CYKLICZNOSCI.interwalMinutowy:
      return `Uruchomienie co ${nastawy.krok} minut, licząc od pełnej godziny.`;
    case WZORCE_CYKLICZNOSCI.interwalGodzinowy:
      return `Uruchomienie co ${nastawy.krok} godzin, w ${dwucyfrowo(nastawy.minuta)} minucie.`;
    default:
      return zapis.trim() === ''
        ? 'Bez cykliczności — automatyka rusza wyzwalaczem albo ręcznie.'
        : 'Zapis własny — okno nie sprowadza go do żadnego wzorca i wysyła w całości.';
  }
}

/** Nazwa dnia tygodnia w numeracji cron; numer spoza zakresu daje sam numer. */
function nazwaDnia(numer: number): string {
  return DNI_TYGODNIA.find(([wartosc]) => wartosc === String(numer))?.[1] ?? String(numer);
}

/** Czy pole niesie pojedynczą liczbę, a nie gwiazdkę, wykaz albo zakres. */
function liczbowe(pole: string): boolean {
  return /^\d+$/.test(pole);
}

/** Krok pola zapisanego gwiazdką z ukośnikiem; pole innego kształtu daje `null`. */
function krokPola(pole: string): number | null {
  const dopasowanie = /^\*\/(\d+)$/.exec(pole);
  if (dopasowanie === null) return null;
  const krok = Number.parseInt(dopasowanie[1] ?? '', 10);
  return Number.isInteger(krok) && krok > 0 ? krok : null;
}

/** Liczba z pola zapisu; pole nieliczbowe oddaje wartość zastępczą. */
function liczbaPola(pole: string, zastepcza: number): number {
  return liczbowe(pole) ? Number.parseInt(pole, 10) : zastepcza;
}

/** Wartość wtłoczona w zakres pola zapisu — okno nie wysyła pola spoza niego. */
function ogranicz(wartosc: number, dolna: number, gorna: number): number {
  if (!Number.isInteger(wartosc)) return dolna;
  return Math.min(gorna, Math.max(dolna, wartosc));
}

/** Zapis dwucyfrowy godziny i minuty — postać przyjęta w zdaniach okna. */
function dwucyfrowo(wartosc: number): string {
  return String(wartosc).padStart(2, '0');
}
