import {
  Command,
  MediaOperationKind,
  type ArchivePackRequest,
  type ArchivePackResponse,
  type DesignAsset,
  type MediaInspectRequest,
  type MediaInspectResponse,
  type MediaTranscodeRequest,
  type MediaTranscodeResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLiczba, czyObiekt, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Trzy komendy arsenału nad materiałem z dysku Operatora: `media.inspect`,
 * `media.transcode` i `archive.pack`.
 *
 * Dlaczego ścieżka na dysku, a nie zasób biblioteki. Wszystkie trzy komendy
 * rozwiązują `assetId` przez repozytorium zasobów modułu Design
 * (`server/internal/core/adapter_narzedzia_media.go`,
 * `adapter_narzedzia_archiwum.go`), a plik biblioteki leży w innym rejestrze —
 * jego identyfikator wraca stamtąd odmową `not_found`. Kontrakt niesie jednak
 * drugą drogę źródła, `sourcePath`, i mówi o niej wprost: „treść jest WCIĄGANA
 * do magazynu, nie dowiązywana". Tą drogą czynność wykonuje się naprawdę,
 * więc tą drogą idzie. Przycisk wysyłający identyfikator pliku biblioteki
 * zawodziłby zawsze i byłby przyciskiem pewnej odmowy.
 *
 * Gdzie ląduje wynik. Bajty idą do magazynu zasobów Designu — tego samego,
 * którym jedzie `design.asset.upload` (`adapter_narzedzia_wynik.go`). Pola
 * `windowId` źródło NIE podaje i to jest rozstrzygnięcie, nie przeoczenie:
 * kontrakt każe podać okno modułu Design, a moduł Library zna wyłącznie okno
 * komunikacji sesji. Okno zmyślone nie zapisałoby się w ogóle (kolumna
 * `zasob_design.okno` jest NOT NULL), a okno cudze pokazałoby wynik w wykazie,
 * do którego on nie należy. Cena jest jedna i jawna, dokładnie ta z kontraktu:
 * „bajty i tak trafiają do magazynu pod sumą kontrolną, ale zasób nie pojawi
 * się w wykazie okna". Widok mówi to Operatorowi zdaniem.
 *
 * Rodzina `media.*` nie ma w kontrakcie ani jednego zdarzenia, więc nic tu nie
 * nasłuchuje — czynność kończy się swoją odpowiedzią i niczym więcej.
 */

/** Rodzaje przetworzenia materiału wraz z nazwą dla Operatora. */
export const RODZAJE_PRZETWORZENIA: readonly { kod: MediaOperationKind; nazwa: string }[] = [
  { kod: MediaOperationKind.Convert, nazwa: 'Zmiana formatu' },
  { kod: MediaOperationKind.Trim, nazwa: 'Wycięcie fragmentu' },
  { kod: MediaOperationKind.ExtractAudio, nazwa: 'Wyodrębnienie ścieżki dźwiękowej' },
  { kod: MediaOperationKind.Resize, nazwa: 'Zmiana rozdzielczości' },
  { kod: MediaOperationKind.Frame, nazwa: 'Zrzut klatki' },
];

/** Nastawy przetworzenia zbierane z formularza. */
export interface NastawaPrzetworzenia {
  sciezka: string;
  operacja: MediaOperationKind;
  format: string;
  odMs?: number;
  doMs?: number;
  szerokosc?: number;
  wysokosc?: number;
}

/** Nastawy pakowania archiwum zbierane z formularza. */
export interface NastawaPakowania {
  sciezka: string;
  format: string;
  nazwa: string;
}

export interface NarzedziaMaterialu {
  /** `media.inspect` — właściwości materiału dźwiękowego albo filmowego. */
  zbadaj(sciezka: string): Promise<Wynik<MediaInspectResponse>>;
  /** `media.transcode` — przetworzenie materiału; wynikiem jest nowy zasób. */
  przetworz(nastawa: NastawaPrzetworzenia): Promise<Wynik<MediaTranscodeResponse>>;
  /** `archive.pack` — spakowanie katalogu albo pliku do archiwum. */
  spakuj(nastawa: NastawaPakowania): Promise<Wynik<ArchivePackResponse>>;
}

export function utworzNarzedziaMaterialu(kanal: Kanal): NarzedziaMaterialu {
  return {
    async zbadaj(sciezka) {
      const zadanie: MediaInspectRequest = { sourcePath: sciezka };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MediaInspect, zadanie),
        Command.MediaInspect,
        (tresc) => czyLiczba(tresc.durationMs) && czyTekst(tresc.format),
      );
    },

    async przetworz(nastawa) {
      // Pola nieobowiązkowe idą wyłącznie wtedy, gdy Operator je wypełnił:
      // zero w `startMs` jest wartością znaczącą (początek materiału), a nie
      // brakiem, więc pusta wartość nie może jechać jako zero.
      const zadanie: MediaTranscodeRequest = {
        sourcePath: nastawa.sciezka,
        operation: nastawa.operacja,
        ...(nastawa.format === '' ? {} : { format: nastawa.format }),
        ...(nastawa.odMs === undefined ? {} : { startMs: nastawa.odMs }),
        ...(nastawa.doMs === undefined ? {} : { endMs: nastawa.doMs }),
        ...(nastawa.szerokosc === undefined ? {} : { width: nastawa.szerokosc }),
        ...(nastawa.wysokosc === undefined ? {} : { height: nastawa.wysokosc }),
      };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MediaTranscode, zadanie),
        Command.MediaTranscode,
        (tresc) => czyObiekt(tresc.asset),
      );
    },

    async spakuj(nastawa) {
      const zadanie: ArchivePackRequest = {
        sourcePath: nastawa.sciezka,
        ...(nastawa.format === '' ? {} : { format: nastawa.format }),
        ...(nastawa.nazwa === '' ? {} : { name: nastawa.nazwa }),
      };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ArchivePack, zadanie),
        Command.ArchivePack,
        (tresc) => czyObiekt(tresc.asset) && czyLiczba(tresc.entries),
      );
    },
  };
}

/**
 * Nazwa wyniku dla Operatora.
 *
 * Zasób bez nazwy własnej nie jest zasobem bez tożsamości: identyfikator
 * niesie ją zawsze, a nazwa bywa pusta, gdy rdzeń jej nie nadał. Zdanie
 * „wynik: (bez nazwy)" mówiłoby o braku, którego nie ma — mówimy więc
 * identyfikatorem, bo po nim wynik da się otworzyć.
 */
export function nazwaWyniku(zasob: DesignAsset): string {
  const nazwa = (zasob.name ?? '').trim();
  const format = (zasob.format ?? '').trim();
  const ogon = format === '' ? '' : ` (${format})`;
  return nazwa === '' ? `${zasob.id}${ogon}` : `${nazwa}${ogon}, zasób ${zasob.id}`;
}

/** Liczba wpisana w pole; pusty napis i wartość niebędąca liczbą znaczą brak. */
export function liczbaZPola(wartosc: string): number | undefined {
  const wpisana = wartosc.trim();
  if (wpisana === '') return undefined;
  const liczba = Number(wpisana);
  return Number.isFinite(liczba) ? liczba : undefined;
}
