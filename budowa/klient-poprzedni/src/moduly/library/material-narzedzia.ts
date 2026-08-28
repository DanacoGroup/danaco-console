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

// Trzy komendy arsenału nad materiałem idą ścieżką źródła, nie identyfikatorem zasobu biblioteki.

/** Rodzaje przetworzenia materiału wraz z nazwą widoczną dla użytkownika na liście wyboru rodzaju operacji. */
export const RODZAJE_PRZETWORZENIA: readonly { kod: MediaOperationKind; nazwa: string }[] = [
  { kod: MediaOperationKind.Convert, nazwa: 'Zmiana formatu' },
  { kod: MediaOperationKind.Trim, nazwa: 'Wycięcie fragmentu' },
  { kod: MediaOperationKind.ExtractAudio, nazwa: 'Wyodrębnienie ścieżki dźwiękowej' },
  { kod: MediaOperationKind.Resize, nazwa: 'Zmiana rozdzielczości' },
  { kod: MediaOperationKind.Frame, nazwa: 'Zrzut klatki' },
];

/** Nastawy przetworzenia zbierane z formularza, przekazywane dalej jako pojedyncze żądanie transkodowania. */
export interface NastawaPrzetworzenia {
  sciezka: string;
  operacja: MediaOperationKind;
  format: string;
  odMs?: number;
  doMs?: number;
  szerokosc?: number;
  wysokosc?: number;
}

/** Nastawy pakowania archiwum zbierane z formularza, przekazywane dalej jako pojedyncze żądanie spakowania. */
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
      // Pola nieobowiązkowe idą wyłącznie po wypełnieniu: zero w polu startu jest wartością znaczącą.
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
 * Nazwa wyniku dla przeglądającego: zasób bez nazwy własnej nie jest zasobem
 * bez tożsamości, więc puste imię zastępuje identyfikator.
 */
export function nazwaWyniku(zasob: DesignAsset): string {
  const nazwa = (zasob.name ?? '').trim();
  const format = (zasob.format ?? '').trim();
  const ogon = format === '' ? '' : ` (${format})`;
  return nazwa === '' ? `${zasob.id}${ogon}` : `${nazwa}${ogon}, zasób ${zasob.id}`;
}

/** Liczba wpisana w pole formularza; pusty napis i wartość niebędąca liczbą oznaczają wspólnie brak wartości. */
export function liczbaZPola(wartosc: string): number | undefined {
  const wpisana = wartosc.trim();
  if (wpisana === '') return undefined;
  const liczba = Number(wpisana);
  return Number.isFinite(liczba) ? liczba : undefined;
}
