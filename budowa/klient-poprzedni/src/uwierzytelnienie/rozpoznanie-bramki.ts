import type { AuthSession, ErrorInfo } from '../../../shared/contract';
import type { Tryb } from './postac-bramki';
import { bezUchwytu, type MetodaWejscia, type ZrodloUwierzytelnienia } from './zrodlo-auth';

/** Rozpoznanie wejścia w osobnym pliku ustala, co ma się stać, zanim operator cokolwiek zobaczy, sprawdzalne bez ani jednego elementu przeglądarki, na podstawionym źródle danych. */
export type PowodBezPrzeslony = 'wymog-wylaczony' | 'polaczenie-zwiazane';

/** Co ekran ma zrobić po rozpoznaniu wejścia: pominąć przesłonę, wejść sesją żywą, pokazać formularz albo zgłosić stan niepewny rdzenia. */
export type DrogaWejscia =
  /** Przesłona nie staje w ogóle — nie ma czego strzec. */
  | { rodzaj: 'bez-przeslony'; powod: PowodBezPrzeslony }
  /** Sesja żywa — wejście bez pytania o cokolwiek. */
  | { rodzaj: 'sesja'; sesja: AuthSession }
  /** Formularz wejścia albo pierwszego hasła. */
  | {
      rodzaj: 'formularz';
      tryb: Extract<Tryb, 'wejscie' | 'zalozenie'>;
      metody: MetodaWejscia[];
      /** Zdanie o losie poprzedniej sesji, gdy była i nie przeżyła. */
      notatka?: string;
    }
  /** Rdzeń nie rozstrzygnął — ekran pokazuje odmowę i „Spróbuj ponownie". */
  | { rodzaj: 'niepewny'; blad?: ErrorInfo; bezUchwytu: boolean };

/** Zapis sesji bramki jest wstrzykiwany, żeby rozpoznanie dało się sprawdzić bez przeglądarki, jedynie na podstawionym źródle danych. */
export interface ZapisSesji {
  odczytaj(): AuthSession | null;
  skasuj(): void;
}

export interface OpisRozpoznania {
  zrodlo: ZrodloUwierzytelnienia;
  zapis: ZapisSesji;
  /** Tożsamość maszyny dla sondy PIN-u; pusty napis znaczy „pamięci nie ma". */
  urzadzenie(): string;
  /** Zdanie o odmowie przedłużenia — składa je ekran, bo zna nazwy czynności. */
  opiszOdmowe(obszar: string, blad: ErrorInfo | undefined, bezUchwytu: boolean): string;
  /** Kres oczekiwania na powitanie; przekazywany źródłu. */
  limitPowitania?: number;
}

export async function rozpoznajWejscie(opis: OpisRozpoznania): Promise<DrogaWejscia> {
  const powitanie = await opis.zrodlo.powitanie(opis.limitPowitania);

  if (powitanie.loginRequired === false) {
    return { rodzaj: 'bez-przeslony', powod: 'wymog-wylaczony' };
  }

  const zapisana = opis.zapis.odczytaj();

  if (zapisana === null && powitanie.authenticated === true) {
    return { rodzaj: 'bez-przeslony', powod: 'polaczenie-zwiazane' };
  }

  let notatka: string | undefined;
  if (zapisana !== null) {
    const odpowiedz = await opis.zrodlo.przedluz(zapisana.token);
    if (odpowiedz.udana && odpowiedz.wynik !== undefined) {
      return { rodzaj: 'sesja', sesja: odpowiedz.wynik.session };
    }
    // Sesja martwa nie jest usterką; zapis znika, a zdanie o jej losie idzie nad formularz.
    opis.zapis.skasuj();
    notatka = opis.opiszOdmowe('przedluzenie', odpowiedz.blad, bezUchwytu(odpowiedz));
  }

  const tryb = await rozpoznajTryb(opis, powitanie.gatewayConfigured);
  if (typeof tryb !== 'string') return tryb;

  return {
    rodzaj: 'formularz',
    tryb,
    metody: await metodyWejscia(opis, tryb),
    ...(notatka === undefined ? {} : { notatka }),
  };
}

/** Tryb formularza pochodzi z powitania, a przy jego milczeniu ustala się sondą, bo wartości domyślnej dla zgadniętego formularza nie ma. */
async function rozpoznajTryb(
  opis: OpisRozpoznania,
  zalozona: boolean | undefined,
): Promise<'wejscie' | 'zalozenie' | DrogaWejscia> {
  if (zalozona === true) return 'wejscie';
  if (zalozona === false) return 'zalozenie';
  const rozpoznanie = await opis.zrodlo.zbadaj();
  if (rozpoznanie.stan === 'ustawiona') return 'wejscie';
  if (rozpoznanie.stan === 'nieustawiona') return 'zalozenie';
  return {
    rodzaj: 'niepewny',
    ...(rozpoznanie.blad === undefined ? {} : { blad: rozpoznanie.blad }),
    bezUchwytu: rozpoznanie.stan === 'bez-uchwytu',
  };
}

/** Metody czynne na tej maszynie: przy pierwszym uruchomieniu jest jedna z konieczności, bo PIN zakłada się dopiero po zalogowaniu w ustawieniach. */
async function metodyWejscia(
  opis: OpisRozpoznania,
  tryb: 'wejscie' | 'zalozenie',
): Promise<MetodaWejscia[]> {
  if (tryb === 'zalozenie') return ['haslo'];
  const urzadzenie = opis.urzadzenie();
  if (urzadzenie === '') return ['haslo'];
  return (await opis.zrodlo.zbadajPin(urzadzenie)) ? ['haslo', 'pin'] : ['haslo'];
}
