import { Command, type ConfigEntry } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { nazwaZasiegu } from '../konfiguracja/zasiegi';
import {
  poleKlucza,
  postacWartosci,
  ustawieniaDomyslne,
  ZASIEG_OKNA,
  type UstawieniaOkna,
} from './klucze-ustawien';

// Wartość obowiązująca: to, czym naprawdę pojedzie model, obok wartości zapisanej na poziomie okna.

/** Jedna nastawa w postaci obowiązującej wraz z poziomem zasięgu i bytem poziomu, z których ona pochodzi. */
export interface NastawaObowiazujaca {
  /** Wartość rozstrzygnięta przez rdzeń; pusta znaczy wartość domyślną katalogu. */
  wartosc: string;
  // Poziom zasięgu, na którym rdzeń znalazł wartość; napis pusty znaczy brak zapisu.
  zasieg: string;
  /** Byt poziomu; pusty dla poziomu globalnego. */
  bytZasiegu: string;
}

/** Komplet nastaw okna w postaci obowiązującej, wraz ze znacznikiem tego, czy rdzeń już oddał politykę. */
export type ObowiazujaceOkna = {
  // Czy rdzeń zdążył oddać politykę; do tego czasu stery pokazują poziom okna.
  znane: boolean;
} & Record<keyof UstawieniaOkna, NastawaObowiazujaca>;

/** Nastawa nieznana: niesie samą wartość domyślną, bez poziomu zasięgu i bez bytu tego poziomu zasięgu. */
function nastawaPusta(wartosc: string): NastawaObowiazujaca {
  return { wartosc, zasieg: '', bytZasiegu: '' };
}

/** Komplet nastaw obowiązujących przed odpowiedzią rdzenia, złożony wyłącznie z pojedynczych nastaw nieznanych. */
export function obowiazujaceNieznane(): ObowiazujaceOkna {
  const domyslne = ustawieniaDomyslne();
  return {
    znane: false,
    hostWykonania: nastawaPusta(domyslne.hostWykonania),
    kanalZapasowy: nastawaPusta(domyslne.kanalZapasowy),
    nakladRozumowania: nastawaPusta(domyslne.nakladRozumowania),
  };
}

/**
 * Nanosi na komplet jeden wpis polityki efektywnej.
 *
 * Wpis klucza spoza kompletu jest pomijany bez błędu — polityka efektywna
 * oddaje cały rejestr ustawień, nie trzy klucze sterowania.
 */
export function naniesObowiazujaca(
  komplet: ObowiazujaceOkna,
  wpis: ConfigEntry,
): ObowiazujaceOkna {
  const pole = poleKlucza(wpis.key);
  if (pole === null) return komplet;
  return {
    ...komplet,
    [pole]: {
      wartosc: postacWartosci(wpis.key, wpis.value),
      zasieg: wpis.scope ?? '',
      bytZasiegu: wpis.scopeId ?? '',
    },
  };
}

/** Komplet nastaw obowiązujących całego okna złożony z wpisów polityki efektywnej oddanych przez rdzeń. */
export function obowiazujaceZWpisow(wpisy: readonly ConfigEntry[]): ObowiazujaceOkna {
  let komplet: ObowiazujaceOkna = { ...obowiazujaceNieznane(), znane: true };
  for (const wpis of wpisy) komplet = naniesObowiazujaca(komplet, wpis);
  return komplet;
}

/**
 * Czy nastawa obowiązująca pochodzi spoza poziomu okna — a więc czy widoczna
 * wartość nie została ustawiona w tym oknie.
 *
 * Poziom pusty (wartość domyślna katalogu) nie jest poziomem szerszym: nie ma
 * o czym dopisywać zdania, skoro nikt nic nie zapisał.
 */
export function zPoziomuSzerszego(
  nastawa: NastawaObowiazujaca,
  idOkna: string,
): boolean {
  if (nastawa.zasieg === '') return false;
  // Poziom okna bierzemy ze stałej kontraktu, nie z literału.
  return !(nastawa.zasieg === String(ZASIEG_OKNA) && nastawa.bytZasiegu === idOkna);
}

/**
 * Zdanie steru: nazwa wartości, którą naprawdę pojedzie model, z dopiskiem
 * o poziomie, gdy wartość przychodzi spoza okna.
 *
 * Dopóki rdzeń nie odpowiedział, zdanie niesie wartość poziomu okna — bez
 * wygaszania steru i bez pustki.
 */
export function zdanieSteru(
  nastawa: NastawaObowiazujaca,
  znane: boolean,
  wartoscOkna: string,
  idOkna: string,
  nazwij: (wartosc: string) => string,
): string {
  if (!znane) return nazwij(wartoscOkna);
  if (!zPoziomuSzerszego(nastawa, idOkna)) return nazwij(nastawa.wartosc);
  return `${nazwij(nastawa.wartosc)} · obowiązuje z poziomu ${nazwaZasiegu(nastawa.zasieg)}`;
}

/** Odczyt polityki efektywnej okna: pyta rdzeń o nastawy obowiązujące i oddaje komplet wywołaniem zwrotnym. */
export interface OdczytObowiazujacej {
  /** Pyta rdzeń o politykę efektywną i oddaje wynik wywołaniem zwrotnym. */
  wczytaj(przyKomplecie: (komplet: ObowiazujaceOkna) => void): void;
}

export function utworzOdczytObowiazujacej(
  kanal: Kanal,
  idOkna: () => string,
): OdczytObowiazujacej {
  return {
    wczytaj(przyKomplecie) {
      // Bez pola poziomu — to ono rozgałęzia adapter rdzenia na politykę efektywną.
      kanal.wyslij(Command.ConfigGet, { scopeId: idOkna() }, (wynik) => {
        // Niepowodzenie nie zmienia niczego: stery zostają przy poziomie okna.
        if (!wynik.udany) return;
        przyKomplecie(obowiazujaceZWpisow(wynik.wynik?.entries ?? []));
      });
    },
  };
}
