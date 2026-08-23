/**
 * Podgląd kolejnych uruchomień z notacji cron — pozycja panelu akcji okna
 * Scheduler („podgląd kolejnych uruchomień”).
 *
 * Rachunek stoi po stronie klienta, bo dotyczy cykliczności jeszcze
 * niezapisanej: kontrakt nie ma ani pola z wykazem kolejnych terminów, ani
 * komendy próbnego wyliczenia. Wartością obowiązującą po zapisie pozostaje
 * `AutomationSchedule.nextRunAt` liczone przez rdzeń.
 *
 * Rachunek idzie w czasie UTC — tak samo jak w rdzeniu, żeby oba wyniki dawały
 * się porównać. Zapis nieczytelny daje wykaz pusty.
 */

/** Ile terminów pokazuje podgląd. */
const LICZBA_TERMINOW = 5;

/** Ile dni wprzód wolno szukać terminu — tyle samo, co w rdzeniu. */
const HORYZONT_DNI = 366;

export function nastepneUruchomienia(
  zapis: string,
  od: Date = new Date(),
  ile: number = LICZBA_TERMINOW,
): Date[] {
  const pola = polaCron(zapis);
  if (pola === null) return [];
  const terminy: Date[] = [];
  let kursor = new Date(Date.UTC(
    od.getUTCFullYear(), od.getUTCMonth(), od.getUTCDate(),
    od.getUTCHours(), od.getUTCMinutes(), 0, 0,
  ));
  for (let licznik = 0; licznik < ile; licznik += 1) {
    kursor = new Date(kursor.getTime() + 60_000);
    const termin = pierwszyTermin(pola, kursor);
    if (termin === null) break;
    terminy.push(termin);
    kursor = termin;
  }
  return terminy;
}

/** Zbiory wartości pięciu pól zapisu wraz z rozpoznaniem gwiazdek pól dnia. */
interface PolaCron {
  minuty: Set<number>;
  godziny: Set<number>;
  dniMiesiaca: Set<number>;
  miesiace: Set<number>;
  dniTygodnia: Set<number>;
  dzienDowolny: boolean;
  tydzienDowolny: boolean;
}

/** Rozkłada zapis pięciopolowy; zapis niezrozumiały daje `null`. */
export function polaCron(zapis: string): PolaCron | null {
  const czesci = zapis.trim().split(/\s+/);
  if (czesci.length !== 5) return null;
  const minuty = zbiorPola(czesci[0] ?? '', 0, 59);
  const godziny = zbiorPola(czesci[1] ?? '', 0, 23);
  const dniMiesiaca = zbiorPola(czesci[2] ?? '', 1, 31);
  const miesiace = zbiorPola(czesci[3] ?? '', 1, 12);
  const dniTygodnia = zbiorPola(czesci[4] ?? '', 0, 6);
  if (
    minuty === null || godziny === null || dniMiesiaca === null ||
    miesiace === null || dniTygodnia === null
  ) {
    return null;
  }
  return {
    minuty,
    godziny,
    dniMiesiaca,
    miesiace,
    dniTygodnia,
    dzienDowolny: czesci[2] === '*',
    tydzienDowolny: czesci[4] === '*',
  };
}

/** Pierwsza chwila nie wcześniejsza niż `od`, pasująca do zapisu. */
function pierwszyTermin(pola: PolaCron, od: Date): Date | null {
  for (let dzien = 0; dzien < HORYZONT_DNI; dzien += 1) {
    const data = new Date(od.getTime() + dzien * 86_400_000);
    if (!dzienPasuje(pola, data)) continue;
    const odMinuty = dzien === 0 ? od.getUTCHours() * 60 + od.getUTCMinutes() : 0;
    for (let minuta = odMinuty; minuta < 1440; minuta += 1) {
      if (!pola.godziny.has(Math.floor(minuta / 60)) || !pola.minuty.has(minuta % 60)) continue;
      return new Date(Date.UTC(
        data.getUTCFullYear(), data.getUTCMonth(), data.getUTCDate(),
        Math.floor(minuta / 60), minuta % 60, 0, 0,
      ));
    }
  }
  return null;
}

/** Reguła klasyczna: gdy oba pola dnia są wskazane, wystarczy trafienie jednego. */
function dzienPasuje(pola: PolaCron, data: Date): boolean {
  if (!pola.miesiace.has(data.getUTCMonth() + 1)) return false;
  const dzien = pola.dniMiesiaca.has(data.getUTCDate());
  const tydzien = pola.dniTygodnia.has(data.getUTCDay());
  if (pola.dzienDowolny && pola.tydzienDowolny) return true;
  if (pola.dzienDowolny) return tydzien;
  if (pola.tydzienDowolny) return dzien;
  return dzien || tydzien;
}

/** Rozkłada jedno pole zapisu na zbiór wartości; pole błędne daje `null`. */
function zbiorPola(pole: string, dolna: number, gorna: number): Set<number> | null {
  const zbior = new Set<number>();
  for (const czlon of pole.split(',')) {
    if (!dopiszCzlon(zbior, czlon.trim(), dolna, gorna)) return null;
  }
  return zbior.size === 0 ? null : zbior;
}

/** Dokłada do zbioru wartości jednego członu: gwiazdka, `a`, `a-b` oraz krok po ukośniku. */
function dopiszCzlon(zbior: Set<number>, czlon: string, dolna: number, gorna: number): boolean {
  const [zakres = '', krokTekst] = czlon.split('/');
  let krok = 1;
  if (krokTekst !== undefined) {
    krok = Number.parseInt(krokTekst, 10);
    if (!Number.isInteger(krok) || krok <= 0) return false;
  }
  let od = dolna;
  let do_ = gorna;
  if (zakres !== '*') {
    const [poczatek = '', koniec] = zakres.split('-');
    od = Number.parseInt(poczatek, 10);
    if (!Number.isInteger(od)) return false;
    do_ = od;
    if (koniec !== undefined) {
      do_ = Number.parseInt(koniec, 10);
      if (!Number.isInteger(do_)) return false;
    }
  }
  if (od < dolna || do_ > gorna || od > do_) return false;
  for (let wartosc = od; wartosc <= do_; wartosc += krok) zbior.add(wartosc);
  return true;
}
