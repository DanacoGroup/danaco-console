import { invoke } from '@tauri-apps/api/core';
import { listen } from '@tauri-apps/api/event';

import { czyPowlokaNatywna } from '../powloka/powloka-natywna';
import type { Wydanie } from './wykaz-wydan';

/**
 * Most do powłoki — wykonanie aktualizacji.
 *
 * Przeglądarka nie podmieni pliku wykonywalnego i nie uruchomi aplikacji
 * ponownie; potrafi to wyłącznie powłoka natywna (`desktop/src-tauri`), bo to
 * ona stawia proces i zna swoje miejsce na dysku. Interfejs rozpoznaje, że jest
 * co zakładać, i przekazuje żądanie.
 *
 * Suma kontrolna idzie w żądaniu, bo powłoka ma odmówić założenia pliku,
 * którego suma się nie zgadza — dlatego interfejs nigdy nie woła aktualizacji
 * bez sumy z wykazu wydań. Poza powłoką natywną i przy powłoce, która polecenia
 * nie zna, wynikiem jest nazwana odmowa — nie wyjątek i nie cisza.
 */

/** Nazwa polecenia powłoki; odpowiednik `polecenia::wykonaj_aktualizacje`. */
const POLECENIE_AKTUALIZACJI = 'wykonaj_aktualizacje';

/**
 * Nazwa polecenia powłoki pytającego, czym jest ta kopia aplikacji.
 *
 * Wykaz poleceń powłoki (`desktop/src-tauri/src/main.rs`) tego polecenia nie
 * zawiera, więc `drogaAktualizacji()` oddaje `null` — „powłoka nie mówi" —
 * a baner nie obiecuje wtedy restartu. Dołożenie polecenia po stronie powłoki
 * niczego tu nie łamie: odpowiedź zaczyna przychodzić, zdanie banera robi się
 * dokładniejsze.
 */
const POLECENIE_DROGI = 'droga_aktualizacji';

/**
 * Nazwa zdarzenia, którym powłoka rozgłasza postęp pobrania.
 *
 * Powłoka tego zdarzenia jeszcze nie rozgłasza: pętla pobrania (`pobranie.rs`)
 * liczy bajty, ale ich nie wysyła. Nasłuch stoi tu założony, żeby dołożenie
 * rozgłoszenia po stronie powłoki było jedyną potrzebną zmianą; dopóki nic nie
 * przychodzi, pas mówi wprost, że powłoka postępu nie podaje, zamiast rysować
 * pasek bez danych.
 */
const ZDARZENIE_POSTEPU = 'powloka:aktualizacja-postep';

/**
 * Przebieg oddany przez powłokę po udanej aktualizacji — kształt
 * `aktualizacja::Przebieg` (`desktop/src-tauri/src/aktualizacja/mod.rs`).
 *
 * Odpowiedź dociera przed restartem: powłoka odkłada ponowne uruchomienie
 * o `restart_za_ms`, żeby baner zdążył powiedzieć „udało się". Bez tej zwłoki
 * okno ginęłoby przed odebraniem odpowiedzi i wyglądałoby to jak awaria.
 */
export interface PrzebiegAktualizacji {
  droga: string;
  zalozone: string;
  bajtow: number;
  rdzen: string;
  restart_za_ms: number;
  /** Gotowe zdanie dla Operatora — układa je powłoka, bo tylko ona zna przebieg. */
  zdanie: string;
}

/** Odmowa powłoki: kod do rozgałęzienia i gotowe zdanie dla Operatora. */
interface OdmowaPowloki {
  powod: string;
  zdanie: string;
}

/**
 * Wynik próby aktualizacji.
 *
 * Odmowa niesie osobno kod i zdanie, bo to dwie różne rzeczy dla dwóch różnych
 * odbiorców: `zdanie` czyta Operator, `kod` czyta baner. Bez kodu każda odmowa
 * wyglądałaby na trwałą, także brak łączności, który minie za minutę.
 */
export type WynikAktualizacji =
  | { udana: true; przebieg: PrzebiegAktualizacji }
  | { udana: false; kod: string; powod: string };

/**
 * Kody odmów, po których ponowienie ma sens — brakło czegoś, co może wrócić.
 *
 * Reszta kodów (`suma-niezgodna`, `droga-niedostepna`, `adres-nie-https`,
 * `suma-w-zlym-zapisie`, `plik-pusty`, …) opisuje brak trwały: powtórzone
 * kliknięcie powtórzyłoby tę samą odmowę co do słowa, więc baner zostaje przy
 * przeczytanym powodzie aż do następnego obiegu pytania.
 */
const ODMOWY_PRZEMIJAJACE = new Set([
  'brak-lacznosci',
  'odpowiedz-serwera',
  'pobieranie-przerwane',
  'zapis-nieudany',
  'plik-ponad-pulap',
  'przebieg-przerwany',
  'aktualizacja-w-toku',
]);

/** Czy po tej odmowie warto dać Operatorowi kliknąć jeszcze raz. */
export function odmowaPrzemijajaca(kod: string): boolean {
  return ODMOWY_PRZEMIJAJACE.has(kod);
}

/**
 * Droga założenia wydania na tej kopii — odpowiedź powłoki, nie domysł strony.
 *
 * Sama obecność powłoki nie wystarcza za odpowiedź: na kopii z pakietu `.deb`
 * powłoka jest, ale podmiany nie wykona (`droga.rs`: brak zmiennej `APPIMAGE`
 * → odmowa `droga-niedostepna`). Miejsce, w którym pytamy o samą powłokę, woła
 * `czyPowlokaNatywna()` po imieniu.
 *
 * `null` znaczy „powłoka nie mówi", a nie „nie da się" — tak jest wszędzie tam,
 * gdzie powłoka polecenia nie zna. Wywołujący ma wtedy milczeć o restarcie,
 * a nie zgadywać w którąkolwiek stronę.
 */
export interface DrogaAktualizacji {
  /** Czy powłoka umie podmienić TĘ kopię — odpowiednik udanego `droga::rozpoznaj()`. */
  mozliwa: boolean;
  /** Krótka nazwa drogi (`appimage`, `instalka-nsis`) — do dziennika i pomiaru. */
  nazwa: string;
  /** Gotowe zdanie dla Operatora; przy `mozliwa: false` niesie powód i radę. */
  zdanie: string;
}

/** Pyta powłokę o drogę założenia; `null`, gdy powłoka nie odpowiada na to pytanie. */
export async function drogaAktualizacji(): Promise<DrogaAktualizacji | null> {
  if (!czyPowlokaNatywna()) return null;
  try {
    const odpowiedz: unknown = await invoke(POLECENIE_DROGI);
    if (typeof odpowiedz !== 'object' || odpowiedz === null) return null;
    const zapis = odpowiedz as Partial<DrogaAktualizacji>;
    // `mozliwa` jest jedynym polem obowiązkowym: bez niego odpowiedź nie niesie
    // tego, po co pytaliśmy, i lepiej milczeć niż zgadywać.
    if (typeof zapis.mozliwa !== 'boolean') return null;
    return {
      mozliwa: zapis.mozliwa,
      nazwa: typeof zapis.nazwa === 'string' ? zapis.nazwa : '',
      zdanie: typeof zapis.zdanie === 'string' ? zapis.zdanie : '',
    };
  } catch {
    // Powłoka bez tego polecenia odrzuca wywołanie. To nie jest awaria —
    // to jest „nie wiem", i tak właśnie ma zostać przekazane.
    return null;
  }
}

/** Etap, na którym stoi trwająca aktualizacja — słownictwo powłoki. */
export type EtapAktualizacji = 'pobieranie' | 'sprawdzanie-sumy' | 'zakladanie';

/**
 * Postęp rozgłaszany przez powłokę.
 *
 * `calosc` bywa `null` I TO NIE JEST BRAK DANYCH DO ZAŁATANIA. Serwer nie musi
 * podać nagłówka `Content-Length`; wtedy znana jest wyłącznie liczba bajtów już
 * pobranych i procentu NIE MA. Widok ma wtedy pokazać licznik megabajtów, a nie
 * wymyślony udział.
 */
export interface PostepAktualizacji {
  etap: EtapAktualizacji;
  pobrano: number;
  calosc: number | null;
}

const ETAPY: ReadonlySet<string> = new Set<EtapAktualizacji>([
  'pobieranie',
  'sprawdzanie-sumy',
  'zakladanie',
]);

/** Czyta ładunek zdarzenia; `null`, gdy nie ma kształtu postępu. */
function odczytajPostep(ladunek: unknown): PostepAktualizacji | null {
  if (typeof ladunek !== 'object' || ladunek === null) return null;
  const zapis = ladunek as Record<string, unknown>;
  const etap = zapis['etap'];
  const pobrano = zapis['pobrano'];
  if (typeof etap !== 'string' || !ETAPY.has(etap)) return null;
  if (typeof pobrano !== 'number' || !Number.isFinite(pobrano) || pobrano < 0) return null;
  const calosc = zapis['calosc'];
  return {
    etap: etap as EtapAktualizacji,
    pobrano,
    // Całość podana jako zero jest bezużyteczna tak samo jak niepodana —
    // dzielenie przez nią dałoby procent, którego nikt nie zmierzył.
    calosc:
      typeof calosc === 'number' && Number.isFinite(calosc) && calosc > 0 ? calosc : null,
  };
}

/**
 * Zaczyna nasłuchiwać postępu pobrania. Oddaje odwołanie nasłuchu.
 *
 * Poza powłoką i na powłoce, która tego zdarzenia nie rozgłasza, nie dzieje się
 * nic: słuchacz nigdy nie jest wołany, a odwołanie jest czynnością pustą.
 * Nasłuch nie ma prawa przeszkodzić w aktualizacji, o której opowiada.
 */
export async function nasluchujPostepu(
  sluchacz: (postep: PostepAktualizacji) => void,
): Promise<() => void> {
  if (!czyPowlokaNatywna()) return () => {};
  try {
    const odwolaj = await listen(ZDARZENIE_POSTEPU, (zdarzenie) => {
      const postep = odczytajPostep(zdarzenie.payload);
      if (postep !== null) sluchacz(postep);
    });
    return () => {
      odwolaj();
    };
  } catch {
    return () => {};
  }
}

/**
 * Prosi powłokę o założenie wydania i ponowne uruchomienie aplikacji.
 *
 * Obietnica rozstrzyga się WYŁĄCZNIE odmową: przy powodzeniu powłoka zamyka
 * proces i nikt na wynik nie czeka, bo nie ma już czego czekać. Kod po tym
 * wywołaniu musi więc zakładać, że może się nie wykonać — i tak właśnie jest
 * napisany baner.
 */
export async function wykonajAktualizacje(wydanie: Wydanie): Promise<WynikAktualizacji> {
  if (!czyPowlokaNatywna()) {
    return {
      udana: false,
      kod: 'poza-powloka',
      powod:
        'aktualizacja z poziomu aplikacji działa w oknie powłoki Danaco Console; ' +
        'w przeglądarce pobierz wydanie ze strony danaco-console.pl',
    };
  }
  if (!wydanie.plik) {
    return { udana: false, kod: 'wydanie-bez-pliku', powod: 'wydanie nie wskazuje pliku do pobrania' };
  }
  if (!wydanie.suma) {
    // Świadoma odmowa, nie przeoczenie: bez sumy nie ma czego sprawdzić, a plik
    // z sieci bez sprawdzenia jest dokładnie tym, czego ten most ma nie robić.
    return {
      udana: false,
      kod: 'wydanie-bez-sumy',
      powod:
        'wydanie ' + wydanie.wersja + ' nie niesie sumy kontrolnej, ' +
        'a pliku bez sprawdzenia sumy aplikacja nie zakłada',
    };
  }
  try {
    // Nazwy pól są NAZWAMI Z JS, nie z Rusta: powłoka deklaruje `suma_sha256`,
    // a Tauri v2 przyjmuje argumenty w camelCase. Rozbieżność jest zapisana
    // tutaj, bo to jedyne miejsce, w którym da się ją zobaczyć z tej strony.
    const przebieg = await invoke<PrzebiegAktualizacji>(POLECENIE_AKTUALIZACJI, {
      adres: wydanie.plik,
      sumaSha256: wydanie.suma,
    });
    return { udana: true, przebieg };
  } catch (przyczyna) {
    return { udana: false, kod: kodPrzyczyny(przyczyna), powod: opisPrzyczyny(przyczyna) };
  }
}

/**
 * Wyjmuje z odmowy powłoki kod do rozgałęzienia.
 *
 * Wyjątek bez kształtu odmowy to powłoka, która polecenia nie zna, albo
 * przerwany most IPC. Dostaje osobny kod zamiast zgadywanego kodu powłoki:
 * baner ma po czym poznać, że kodu nie było.
 */
function kodPrzyczyny(przyczyna: unknown): string {
  if (typeof przyczyna === 'object' && przyczyna !== null) {
    const odmowa = przyczyna as Partial<OdmowaPowloki>;
    if (typeof odmowa.powod === 'string' && odmowa.powod !== '') return odmowa.powod;
  }
  return 'odmowa-bez-kodu';
}

/**
 * Zamienia odmowę powłoki na zdanie dla Operatora.
 *
 * Zdanie układa powłoka, nie interfejs: tylko ona wie, czy zabrakło łączności,
 * prawa zapisu, czy zgodności sumy. Drugi zestaw zdań tutaj byłby drugą prawdą
 * o tej samej odmowie.
 */
function opisPrzyczyny(przyczyna: unknown): string {
  if (typeof przyczyna === 'object' && przyczyna !== null) {
    const odmowa = przyczyna as Partial<OdmowaPowloki>;
    if (typeof odmowa.zdanie === 'string' && odmowa.zdanie !== '') return odmowa.zdanie;
    if (typeof odmowa.powod === 'string' && odmowa.powod !== '') return odmowa.powod;
  }
  if (typeof przyczyna === 'string' && przyczyna !== '') return przyczyna;
  if (przyczyna instanceof Error && przyczyna.message !== '') return przyczyna.message;
  return 'powłoka nie wykonała aktualizacji i nie podała powodu';
}
