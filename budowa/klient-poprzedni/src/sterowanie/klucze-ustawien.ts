import { ConfigScope } from '../../../shared/contract';

/** Ustawienia okna zapisywane na poziomie zasięgu okna: host wykonania, kanał zapasowy i nakład rozumowania, zapisywane komendą ustawienia poziomu okna. */
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

/** Komplet ustawień okna niesionych poza treścią żądania aktualizacji okna: host wykonania, kanał zapasowy i nakład rozumowania. */
export interface UstawieniaOkna {
  hostWykonania: string;
  kanalZapasowy: string;
  /** Wartość wyliczenia katalogu STOPNIE_NAKLADU, nie liczba położenia. */
  nakladRozumowania: string;
}

/** Stopnie nakładu rozumowania w kolejności katalogu rdzenia, jako wyliczenie napisowe; napis pusty jest pełnoprawnym stopniem katalogu. */
export const STOPNIE_NAKLADU = ['', 'low', 'medium', 'high', 'xhigh', 'max'] as const;
export type StopienNakladu = (typeof STOPNIE_NAKLADU)[number];

/** Skrajne położenia suwaka nakładu rozumowania: indeksy pierwszego i ostatniego stopnia w wykazie stopni. */
export const NAKLAD_NAJMNIEJSZY = 0;
export const NAKLAD_NAJWIEKSZY = STOPNIE_NAKLADU.length - 1;

/** Położenie suwaka odpowiadające podanemu stopniowi; stopień nieznany wraca na początek wykazu stopni. */
export function polozenieNakladu(stopien: string): number {
  const miejsce = STOPNIE_NAKLADU.indexOf(stopien as StopienNakladu);
  return miejsce < 0 ? NAKLAD_NAJMNIEJSZY : miejsce;
}

/** Stopień nakładu odpowiadający położeniu suwaka; położenie spoza wykazu daje pierwszy stopień z listy stopni. */
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

/** Odwzorowanie klucza katalogu rdzenia na pole kompletu ustawień, czytane przy naniesieniu, zdjęciu i rozpoznaniu zapisu. */
const POLA_KLUCZY: Readonly<
  Record<string, { pole: keyof UstawieniaOkna; postac: (wartosc: unknown) => string }>
> = {
  [KluczUstawieniaOkna.HostWykonania]: { pole: 'hostWykonania', postac: tekst },
  [KluczUstawieniaOkna.KanalZapasowy]: { pole: 'kanalZapasowy', postac: tekst },
  [KluczUstawieniaOkna.NakladRozumowania]: { pole: 'nakladRozumowania', postac: stopien },
};

/** Klucze katalogu rdzenia, które komplet ustawień okna niesie i rozpoznaje jako własne pole ustawienia. */
export const KLUCZE_KOMPLETU: readonly string[] = Object.keys(POLA_KLUCZY);

/** Pole kompletu ustawień, do którego należy podany klucz katalogu rdzenia; wartość pusta dla klucza obcego. */
export function poleKlucza(klucz: string): keyof UstawieniaOkna | null {
  return POLA_KLUCZY[klucz]?.pole ?? null;
}

/** Wartość klucza sprowadzona do postaci przyjmowanej przez komplet ustawień; klucz obcy daje napis pusty. */
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

/** Zdejmuje zapis jednego klucza, sprowadzając jego pole do wartości domyślnej katalogu po przywróceniu ustawienia. */
export function zdejmij(ustawienia: UstawieniaOkna, klucz: string): UstawieniaOkna {
  const opis = POLA_KLUCZY[klucz];
  if (opis === undefined) return ustawienia;
  return { ...ustawienia, [opis.pole]: ustawieniaDomyslne()[opis.pole] };
}

/** Wartość tekstowa ustawienia odczytana z konfiguracji; każdy inny kształt wartości znaczy brak wartości. */
function tekst(wartosc: unknown): string {
  return typeof wartosc === 'string' ? wartosc : '';
}

/** Stopień nakładu rozumowania sprowadzony do wykazu stopni katalogu; wartość spoza wykazu daje napis pusty. */
function stopien(wartosc: unknown): string {
  const napis = tekst(wartosc);
  return STOPNIE_NAKLADU.includes(napis as StopienNakladu) ? napis : '';
}
