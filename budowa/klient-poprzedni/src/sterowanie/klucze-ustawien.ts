import { ConfigScope } from '../../../shared/contract';

/**
 * Ustawienia okna zapisywane na poziomie zasięgu okna.
 *
 * Trzy sterowania kompletu nie mają odpowiednika w treści `window.update`:
 * host wykonania (nazwa hosta jest ustawieniem, nie wartością wyliczenia
 * `ExecutionEnv`), kanał zapasowy oraz nakład rozumowania. Ich drogą jest
 * `config.set` na poziomie zasięgu `window` — poziomie najwęższym,
 * wygrywającym z pozostałymi.
 *
 * Rdzeń zna wyłącznie klucze katalogu ustawień z bazy (tabela
 * `definicja_ustawienia`); klucz wymyślony po stronie klienta przechodzi zapis,
 * ale nie ma żadnego skutku. Klucze poniżej są przepisane z katalogu co do znaku.
 */
export const ZASIEG_OKNA = ConfigScope.Window;

/**
 * Klucze ustawień zapisywanych na poziomie okna komunikacji.
 * Źródło: `server/internal/store/migracja_012_katalog_ustawien.sql`
 * (kolumna `definicja_ustawienia.klucz`).
 */
export const KluczUstawieniaOkna = {
  /** Host wykonania wskazany przy zasięgu `remote`. */
  HostWykonania: 'host_wykonania',
  /** Kanał modelu użyty, gdy kanał główny nie odpowiada. */
  KanalZapasowy: 'kanal_modelu_zapasowy',
  /** Nakład rozumowania: od odpowiedzi szybkiej po najstaranniejszą. */
  NakladRozumowania: 'naklad_rozumowania',
} as const;
export type KluczUstawieniaOkna =
  (typeof KluczUstawieniaOkna)[keyof typeof KluczUstawieniaOkna];

/** Komplet ustawień okna niesionych poza treścią `window.update`. */
export interface UstawieniaOkna {
  hostWykonania: string;
  kanalZapasowy: string;
  /** Wartość wyliczenia katalogu, nie liczba — patrz STOPNIE_NAKLADU. */
  nakladRozumowania: string;
}

/**
 * Stopnie nakładu rozumowania w kolejności katalogu (`opcja_ustawienia`
 * dla klucza `naklad_rozumowania`).
 *
 * Katalog rdzenia opisuje ten klucz jako wyliczenie napisowe, więc suwak niesie
 * napis wyliczenia, nie liczbę położenia. Napis pusty jest pełnoprawnym
 * stopniem katalogu — znaczy „rozstrzyga kanał modelu", nie brak ustawienia.
 */
export const STOPNIE_NAKLADU = ['', 'low', 'medium', 'high', 'xhigh', 'max'] as const;
export type StopienNakladu = (typeof STOPNIE_NAKLADU)[number];

/** Skrajne położenia suwaka nakładu — indeksy w wykazie stopni. */
export const NAKLAD_NAJMNIEJSZY = 0;
export const NAKLAD_NAJWIEKSZY = STOPNIE_NAKLADU.length - 1;

/** Położenie suwaka dla stopnia; stopień nieznany wraca na początek wykazu. */
export function polozenieNakladu(stopien: string): number {
  const miejsce = STOPNIE_NAKLADU.indexOf(stopien as StopienNakladu);
  return miejsce < 0 ? NAKLAD_NAJMNIEJSZY : miejsce;
}

/** Stopień spod położenia suwaka; położenie spoza wykazu daje stopień pierwszy. */
export function nakladZPolozenia(polozenie: number): string {
  return STOPNIE_NAKLADU[polozenie] ?? STOPNIE_NAKLADU[NAKLAD_NAJMNIEJSZY];
}

/**
 * Wartości domyślne. Brak ustawienia znaczy wartość domyślną, nie blokadę
 * uruchomienia. Domyślna nakładu jest pusta, tak jak w kolumnie
 * `wartosc_domyslna` katalogu.
 */
export function ustawieniaDomyslne(): UstawieniaOkna {
  return { hostWykonania: '', kanalZapasowy: '', nakladRozumowania: '' };
}

/**
 * Jedno odwzorowanie klucza katalogu na pole kompletu.
 *
 * Klucz katalogu rdzenia trzeba przełożyć na pole kompletu w czterech
 * czynnościach: naniesieniu zapisu, zdjęciu zapisu (zdarzenie `config.changed`
 * o rodzaju `deleted`), rozpoznaniu przynależności klucza do kompletu oraz
 * naniesieniu wartości obowiązującej. Wszystkie cztery czytają ten jeden wykaz,
 * więc dopisanie klucza jest zmianą w jednym miejscu.
 */
const POLA_KLUCZY: Readonly<
  Record<string, { pole: keyof UstawieniaOkna; postac: (wartosc: unknown) => string }>
> = {
  [KluczUstawieniaOkna.HostWykonania]: { pole: 'hostWykonania', postac: tekst },
  [KluczUstawieniaOkna.KanalZapasowy]: { pole: 'kanalZapasowy', postac: tekst },
  [KluczUstawieniaOkna.NakladRozumowania]: { pole: 'nakladRozumowania', postac: stopien },
};

/** Klucze katalogu, które komplet ustawień okna niesie. */
export const KLUCZE_KOMPLETU: readonly string[] = Object.keys(POLA_KLUCZY);

/** Pole kompletu, do którego należy klucz katalogu; `null` dla klucza obcego. */
export function poleKlucza(klucz: string): keyof UstawieniaOkna | null {
  return POLA_KLUCZY[klucz]?.pole ?? null;
}

/** Wartość klucza sprowadzona do postaci kompletu; klucz obcy daje napis pusty. */
export function postacWartosci(klucz: string, wartosc: unknown): string {
  return POLA_KLUCZY[klucz]?.postac(wartosc) ?? '';
}

/**
 * Nanosi na komplet ustawień wpis odczytany z konfiguracji.
 *
 * Klucz spoza kompletu jest pomijany bez błędu — konfiguracja poziomu okna
 * może nieść ustawienia innych warstw interfejsu.
 */
export function nanies(
  ustawienia: UstawieniaOkna,
  klucz: string,
  wartosc: unknown,
): UstawieniaOkna {
  const opis = POLA_KLUCZY[klucz];
  if (opis === undefined) return ustawienia;
  return { ...ustawienia, [opis.pole]: opis.postac(wartosc) };
}

/**
 * Zdejmuje zapis jednego klucza, sprowadzając pole do wartości domyślnej
 * katalogu.
 *
 * Przywrócenie wartości domyślnej (`config.reset`) rdzeń rozgłasza zdarzeniem
 * `config.changed` o rodzaju `deleted`, niosącym wpis ze starą wartością — tak,
 * żeby odbiorca wiedział, co zniknęło. Wpisu z takiego zdarzenia nie wolno
 * nanosić jak świeżego zapisu; zdjęcie zapisu jest osobną czynnością.
 */
export function zdejmij(ustawienia: UstawieniaOkna, klucz: string): UstawieniaOkna {
  const opis = POLA_KLUCZY[klucz];
  if (opis === undefined) return ustawienia;
  return { ...ustawienia, [opis.pole]: ustawieniaDomyslne()[opis.pole] };
}

/** Wartość tekstowa ustawienia; każdy inny kształt znaczy brak wartości. */
function tekst(wartosc: unknown): string {
  return typeof wartosc === 'string' ? wartosc : '';
}

/** Stopień nakładu sprowadzony do wykazu katalogu. */
function stopien(wartosc: unknown): string {
  const napis = tekst(wartosc);
  return STOPNIE_NAKLADU.includes(napis as StopienNakladu) ? napis : '';
}
