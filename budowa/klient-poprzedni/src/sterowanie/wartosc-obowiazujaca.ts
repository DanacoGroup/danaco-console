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

/**
 * Wartość obowiązująca — to, czym naprawdę pojedzie model — obok wartości
 * zapisanej na poziomie okna.
 *
 * Prawdy są dwie i zlanie ich dałoby trzeci rozjazd. Ster zapisuje na poziomie
 * okna (`config.set`, zasięg `window`), ale rdzeń buduje wywołanie z wartości
 * rozstrzygniętej po wszystkich poziomach zasięgu
 * (`core/adapter_rozmowa_wykonanie.go`, funkcja `Ustal` →
 * `injection/argumenty.go`: `--fallback-model`, `--effort`). Nakład ustawiony
 * globalnie obowiązuje model, a ster czytający sam poziom okna napisałby „Bez
 * wskazania". Etykieta niesie więc wartość obowiązującą; menu i suwak nadal
 * ustawiają poziom okna.
 *
 * `config.effective.get` jest komendą o czym innym: składa obszary konfiguracji
 * sesji zapisane pod kluczami rezolwera `sesja.konfiguracja.<obszar>`
 * (`core/sesja_konfiguracja.go`, stała `przedrostekObszaru`;
 * `core/sesja_konfiguracja_skladanie.go`, funkcja `rozstrzygnijObszary`).
 * Klucza prostego `naklad_rozumowania` — tego, który ster zapisuje i który
 * czyta budowa wywołania — ta droga nie ogląda wcale, więc pokazywałaby pustkę
 * tam, gdzie nastawa jest ustawiona.
 *
 * Drogą właściwą jest `config.get` bez poziomu. Adapter rdzenia rozgałęzia
 * odczyt: z podanym poziomem oddaje surowe wpisy tego poziomu, a bez poziomu —
 * politykę efektywną, czyli po jednym zwycięskim wpisie na klucz, z polem
 * `scope` niosącym poziom, na którym wartość znaleziono
 * (`core/adapter_ustawienia.go`, funkcja `Odczytaj`, gałąź `z.Scope == nil`;
 * `konfig/odwzorowanie_kontraktu.go`, `WpisKontraktu`). Poziom pusty znaczy
 * wartość domyślną katalogu, nie błąd.
 *
 * Czego ta droga nie obejmuje: kontekst rozstrzygania budowany przez rdzeń dla
 * tej gałęzi ma wypełnione wyłącznie okno (`core/adapter_ustawienia.go`,
 * funkcja `kontekstZasiegu`), a droga tury wypełnia dodatkowo kartę sesji i oś
 * modelu (`Ustal`: `konfig.Kontekst{Okno, KartaSesji, Model}`). Wartość
 * zapisana na karcie sesji albo na osi kanału modelu obowiązuje więc wywołanie,
 * a tą komendą się nie pokaże — klient nie ma jak tego domknąć po swojej
 * stronie i nie udaje, że ma.
 */

/** Jedna nastawa w postaci obowiązującej wraz z poziomem, z którego pochodzi. */
export interface NastawaObowiazujaca {
  /** Wartość rozstrzygnięta przez rdzeń; pusta znaczy wartość domyślną katalogu. */
  wartosc: string;
  /**
   * Poziom zasięgu, na którym rdzeń znalazł wartość. Napis pusty znaczy
   * „z żadnego zapisu" — obowiązuje warstwa definicji katalogu.
   * Typ jest napisem, nie wyliczeniem, bo rdzeń wysyła tu również poziom pusty.
   */
  zasieg: string;
  /** Byt poziomu; pusty dla poziomu globalnego. */
  bytZasiegu: string;
}

/** Komplet nastaw okna w postaci obowiązującej. */
export type ObowiazujaceOkna = {
  /**
   * Czy rdzeń zdążył oddać politykę. Do tego czasu stery pokazują poziom okna —
   * pokazywanie wartości domyślnej jako „obowiązującej" byłoby zmyśleniem.
   */
  znane: boolean;
} & Record<keyof UstawieniaOkna, NastawaObowiazujaca>;

/** Nastawa nieznana: bez wartości i bez poziomu. */
function nastawaPusta(wartosc: string): NastawaObowiazujaca {
  return { wartosc, zasieg: '', bytZasiegu: '' };
}

/** Komplet przed odpowiedzią rdzenia. */
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

/** Komplet złożony z wpisów polityki efektywnej oddanych przez rdzeń. */
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

/** Odczyt polityki efektywnej okna. */
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
      // Bez pola `scope` — to ono rozgałęzia adapter rdzenia na politykę
      // efektywną. `scopeId` niesie okno, bo kontekst rozstrzygania buduje się
      // właśnie z niego.
      kanal.wyslij(Command.ConfigGet, { scopeId: idOkna() }, (wynik) => {
        // Niepowodzenie nie zmienia niczego: stery zostają przy poziomie okna,
        // zamiast pokazać domyślne jako obowiązujące.
        if (!wynik.udany) return;
        przyKomplecie(obowiazujaceZWpisow(wynik.wynik?.entries ?? []));
      });
    },
  };
}
