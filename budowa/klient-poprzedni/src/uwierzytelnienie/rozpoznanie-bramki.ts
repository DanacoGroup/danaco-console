import type { AuthSession, ErrorInfo } from '../../../shared/contract';
import type { Tryb } from './postac-bramki';
import { bezUchwytu, type MetodaWejscia, type ZrodloUwierzytelnienia } from './zrodlo-auth';

/**
 * Rozpoznanie wejścia — co ma się stać, zanim Operator cokolwiek zobaczy.
 *
 * Kolejność kroków stoi w osobnym pliku, bo to jedyna część bramki, w której da
 * się popełnić błąd cichy: przesłona stanie tam, gdzie nie powinna, albo nie
 * stanie tam, gdzie musi. Wyjęta z ekranu sprawdza się bez ani jednego elementu
 * DOM, na podstawionym źródle.
 *
 * Kolejność kroków i powód każdego:
 *
 *   1. Powitanie. Rdzeń mówi w `connection.hello` trzy rzeczy naraz: czy wymóg
 *      logowania obowiązuje, czy bramka jest ustawiona i czy to połączenie jest
 *      już związane z sesją. Powitanie i tak leci przy każdym nawiązaniu, więc
 *      krok nie kosztuje ani jednej dodatkowej koperty.
 *
 *   2. Wymóg logowania. `loginRequired === false` znaczy, że Operator wymóg
 *      wyłączył nastawą `gateway.requireLogin` — i wtedy przesłona nie staje
 *      wcale. Przycisk „Pomiń" byłby bramą z furtką, a wola Operatora jest już
 *      wyrażona.
 *
 *   3. Połączenie związane. `authenticated === true` przy braku zapisanej sesji
 *      znaczy, że rdzeń uznał to gniazdo za wejście — nie ma do czego pokazywać
 *      formularza. Rdzeń mówi wprost, że pole jest do tego: „Klient czyta
 *      `authenticated` i sam rozstrzyga, czy pokazać okno logowania"
 *      (`handlers_connection.go`).
 *
 *   4. Sesja zapisana. Przedłużenie (`auth.token.refresh`) jako jedyne
 *      rozstrzyga o życiu sesji i jako jedyne oddaje nowy czas ważności.
 *      Powitanie tego nie zastąpi, więc gdy zapis istnieje, pyta się rdzenia
 *      przedłużeniem.
 *
 *   5. Który formularz. `gatewayConfigured` z powitania; pole puste znaczy
 *      „rdzeń nie wie" i wtedy — i tylko wtedy — idzie sonda `auth.login` bez
 *      sekretu. Milczenie nie jest zamieniane na „nie".
 *
 *   6. Metody. Segment „PIN" staje wyłącznie wtedy, gdy PIN na tej maszynie
 *      naprawdę jest założony; metody niepewnej ekran nie proponuje.
 */

/** Skąd wiadomo, że przesłona nie ma stawać. */
export type PowodBezPrzeslony = 'wymog-wylaczony' | 'polaczenie-zwiazane';

/** Co ekran ma zrobić po rozpoznaniu. */
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

/** Zapis sesji bramki — wstrzykiwany, żeby rozpoznanie dało się sprawdzić bez przeglądarki. */
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
    // Sesja martwa nie jest usterką — jest wiadomością sprzed rozpoczęcia
    // pracy. Zapis znika, a zdanie o jej losie idzie nad formularz.
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

/**
 * Tryb formularza: z powitania, a przy jego milczeniu — sondą.
 *
 * Sonda idzie raz, a jej odmowa wraca w całości — powtórne pytanie tylko po to,
 * żeby wziąć z niego powód, obciążałoby dławik drugi raz tą samą wątpliwością.
 * Wartości domyślnej nie ma: domyślną byłby zgadnięty formularz.
 */
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

/**
 * Metody czynne na tej maszynie.
 *
 * Przy pierwszym uruchomieniu jest jedna z konieczności: PIN zakłada się
 * w Ustawieniach po zalogowaniu, więc przy nieustawionej bramce nie ma go skąd
 * wziąć i nie ma po co o niego pytać rdzenia.
 */
async function metodyWejscia(
  opis: OpisRozpoznania,
  tryb: 'wejscie' | 'zalozenie',
): Promise<MetodaWejscia[]> {
  if (tryb === 'zalozenie') return ['haslo'];
  const urzadzenie = opis.urzadzenie();
  if (urzadzenie === '') return ['haslo'];
  return (await opis.zrodlo.zbadajPin(urzadzenie)) ? ['haslo', 'pin'] : ['haslo'];
}
