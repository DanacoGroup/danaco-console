/**
 * Most do powłoki wykonuje aktualizację: przekazuje żądanie założenia
 * wydania wraz z sumą kontrolną i odbiera odpowiedź powłoki natywnej.
 */
import { invoke } from '@tauri-apps/api/core';
import { listen } from '@tauri-apps/api/event';

import { czyPowlokaNatywna } from '../powloka/powloka-natywna';
import type { Wydanie } from './wykaz-wydan';

/** Nazwa polecenia powłoki wykonującego aktualizację; odpowiednik polecenia `wykonaj_aktualizacje` zaimplementowanego w powłoce natywnej. */
const POLECENIE_AKTUALIZACJI = 'wykonaj_aktualizacje';

/** Nazwa polecenia powłoki pytającego, czym jest ta kopia aplikacji i czy powłoka umie ją podmienić na nowsze wydanie. */
const POLECENIE_DROGI = 'droga_aktualizacji';

/** Nazwa zdarzenia, którym powłoka natywna rozgłasza bieżący postęp pobrania trwającego wydania aktualizacji aplikacji. */
const ZDARZENIE_POSTEPU = 'powloka:aktualizacja-postep';

/** Przebieg oddany przez powłokę po udanej aktualizacji, o kształcie zgodnym ze strukturą `aktualizacja::Przebieg`. */
export interface PrzebiegAktualizacji {
  droga: string;
  zalozone: string;
  bajtow: number;
  rdzen: string;
  restart_za_ms: number;
  /** Gotowe zdanie dla Operatora — układa je powłoka, bo tylko ona zna przebieg. */
  zdanie: string;
}

/** Odmowa powłoki wobec żądania aktualizacji: kod do rozgałęzienia obsługi i gotowe zdanie dla Operatora. */
interface OdmowaPowloki {
  powod: string;
  zdanie: string;
}

/** Wynik próby aktualizacji: powodzenie niosące przebieg albo odmowa niosąca osobno kod i zdanie dla Operatora. */
export type WynikAktualizacji =
  | { udana: true; przebieg: PrzebiegAktualizacji }
  | { udana: false; kod: string; powod: string };

/** Zbiór kodów odmów, po których ponowienie czynności ma sens, ponieważ zabrakło czegoś, co może jeszcze wrócić. */
const ODMOWY_PRZEMIJAJACE = new Set([
  'brak-lacznosci',
  'odpowiedz-serwera',
  'pobieranie-przerwane',
  'zapis-nieudany',
  'plik-ponad-pulap',
  'przebieg-przerwany',
  'aktualizacja-w-toku',
]);

/** Rozstrzyga, czy po otrzymanej odmowie warto dać Operatorowi możliwość ponownego kliknięcia przycisku. */
export function odmowaPrzemijajaca(kod: string): boolean {
  return ODMOWY_PRZEMIJAJACE.has(kod);
}

/** Droga założenia wydania na tej kopii aplikacji jest odpowiedzią powłoki, nie domysłem samego interfejsu. */
export interface DrogaAktualizacji {
  /** Czy powłoka umie podmienić TĘ kopię — odpowiednik udanego `droga::rozpoznaj()`. */
  mozliwa: boolean;
  /** Krótka nazwa drogi (`appimage`, `instalka-nsis`) — do dziennika i pomiaru. */
  nazwa: string;
  /** Gotowe zdanie dla Operatora; przy `mozliwa: false` niesie powód i radę. */
  zdanie: string;
}

/** Pyta powłokę natywną o drogę założenia wydania; oddaje `null`, gdy powłoka nie odpowiada na to pytanie. */
export async function drogaAktualizacji(): Promise<DrogaAktualizacji | null> {
  if (!czyPowlokaNatywna()) return null;
  try {
    const odpowiedz: unknown = await invoke(POLECENIE_DROGI);
    if (typeof odpowiedz !== 'object' || odpowiedz === null) return null;
    const zapis = odpowiedz as Partial<DrogaAktualizacji>;
    // `mozliwa` jest jedynym polem obowiązkowym odpowiedzi, reszta ma wartości domyślne.
    if (typeof zapis.mozliwa !== 'boolean') return null;
    return {
      mozliwa: zapis.mozliwa,
      nazwa: typeof zapis.nazwa === 'string' ? zapis.nazwa : '',
      zdanie: typeof zapis.zdanie === 'string' ? zapis.zdanie : '',
    };
  } catch {
    // Powłoka bez tego polecenia odrzuca wywołanie — to nie jest awaria, tylko brak odpowiedzi.
    return null;
  }
}

/** Etap, na którym w danej chwili stoi trwająca aktualizacja aplikacji, nazwany słownictwem samej powłoki. */
export type EtapAktualizacji = 'pobieranie' | 'sprawdzanie-sumy' | 'zakladanie';

/** Postęp aktualizacji rozgłaszany przez powłokę: etap, liczba pobranych bajtów i znana całość pobrania. */
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

/** Czyta ładunek zdarzenia rozgłoszonego przez powłokę i zwraca `null`, gdy ładunek nie ma kształtu postępu. */
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
    // Całość podana jako zero jest bezużyteczna tak samo jak niepodana.
    calosc:
      typeof calosc === 'number' && Number.isFinite(calosc) && calosc > 0 ? calosc : null,
  };
}

/** Zaczyna nasłuchiwać postępu pobrania rozgłaszanego przez powłokę natywną i oddaje odwołanie tego nasłuchu. */
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

/** Prosi powłokę o założenie pobranego wydania i ponowne uruchomienie aplikacji po jego zainstalowaniu. */
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
    // Świadoma odmowa: bez sumy nie ma czego sprawdzić, więc plik nie jest zakładany.
    return {
      udana: false,
      kod: 'wydanie-bez-sumy',
      powod:
        'wydanie ' + wydanie.wersja + ' nie niesie sumy kontrolnej, ' +
        'a pliku bez sprawdzenia sumy aplikacja nie zakłada',
    };
  }
  try {
    // Nazwy pól są nazwami z interfejsu, nie z rdzenia powłoki.
    const przebieg = await invoke<PrzebiegAktualizacji>(POLECENIE_AKTUALIZACJI, {
      adres: wydanie.plik,
      sumaSha256: wydanie.suma,
    });
    return { udana: true, przebieg };
  } catch (przyczyna) {
    return { udana: false, kod: kodPrzyczyny(przyczyna), powod: opisPrzyczyny(przyczyna) };
  }
}

/** Wyjmuje z przechwyconej odmowy zgłoszonej przez powłokę kod służący do rozgałęzienia dalszej obsługi błędu. */
function kodPrzyczyny(przyczyna: unknown): string {
  if (typeof przyczyna === 'object' && przyczyna !== null) {
    const odmowa = przyczyna as Partial<OdmowaPowloki>;
    if (typeof odmowa.powod === 'string' && odmowa.powod !== '') return odmowa.powod;
  }
  return 'odmowa-bez-kodu';
}

/** Zamienia przechwyconą odmowę zgłoszoną przez powłokę na gotowe zdanie do wyświetlenia Operatorowi na pasie. */
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
