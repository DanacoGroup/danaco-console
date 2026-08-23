import {
  Command,
  type SpeechAudioUploadResponse,
  type SpeechAvailabilityGetResponse,
  type SpeechTranscribeResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLogiczna, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import { opisOdmowy } from '../../komponenty/odmowa';

/**
 * Dyktowanie — wymóg osobisty Właściciela, funkcja działająca.
 *
 * ── Droga jest jedna i już stoi w kontrakcie ────────────────────────────────
 * Nie zakładamy drugiej: `speech.availability.get` sprawdza silnik,
 * `speech.audio.upload` przyjmuje bajty nagrania i oddaje `audioRef`,
 * `speech.transcribe` zamienia nagranie na tekst. Nagranie NIE OPUSZCZA maszyny
 * rdzenia — tak stanowi kontrakt tych komend i tak tu zostaje: bajty jadą do
 * magazynu nagrań rdzenia, nie do sieci.
 *
 * ── Silnik mowy jedzie z pakietem serwera ───────────────────────────────────
 * Rozstrzygnięcie Właściciela: silnik mowy jest **składnikiem pakietu serwera**,
 * razem z rozpoznaniem pisma i obsługą archiwów — nie rzeczą, którą Operator
 * sobie doinstalowuje. Dyktowanie jest więc funkcją, na którą można liczyć,
 * a nie możliwością warunkową.
 *
 * Sprawdzenie dostępności zostaje, bo zmieniło znaczenie: odmowa nie mówi już
 * „tej funkcji nie ma w produkcie", lecz **„ten serwer jest niekompletny"** —
 * i nazywa brakujący składnik. Brak silnika jest usterką wdrożenia serwera,
 * nie ograniczeniem produktu, i tak brzmi zdanie, które Operator zobaczy.
 *
 * Trzy stany odpowiedzi `speech.transcribe` są rozróżnione, bo znaczą różne
 * rzeczy: rozpoznano tekst (`processed` prawda, `transcript` niepusty),
 * przetworzono, lecz mowy nie było (`processed` prawda, `transcript` pusty — to
 * FAKT pomiaru, nie awaria) i nie przetworzono (odmowa).
 *
 * ── Dwa miejsca dyktowania, jeden rachunek ──────────────────────────────────
 * Ten sam przycisk obsadza wiersz polecenia (treść polecenia dla modelu)
 * i treść dokumentu (dyktafon wzorem wstążki pakietu biurowego). Różni je
 * wyłącznie to, komu oddają rozpoznany tekst — `naTekst` — i napis na przycisku.
 * Druga droga nagrywania rozjechałaby się z pierwszą przy pierwszej poprawce.
 */
export interface MowaZrodlo {
  /** Uczciwy stan silnika mowy. */
  mowaDostepnosc(): Promise<Wynik<SpeechAvailabilityGetResponse>>;
  /** Wnosi bajty nagrania do magazynu rdzenia i oddaje odnośnik. */
  mowaWniesNagranie(
    bajtyBase64: string,
    typTresci: string,
    idSesji: string,
    idOkna: string,
  ): Promise<Wynik<SpeechAudioUploadResponse>>;
  /** Zamienia nagranie na tekst silnikiem lokalnym. */
  mowaPrzepisz(odnosnik: string, jezyk: string): Promise<Wynik<SpeechTranscribeResponse>>;
}

export function utworzMowaZrodlo(kanal: Kanal): MowaZrodlo {
  return {
    async mowaDostepnosc() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SpeechAvailabilityGet, {}),
        Command.SpeechAvailabilityGet,
        (tresc) => czyLogiczna(tresc.available),
      );
    },

    async mowaWniesNagranie(bajtyBase64, typTresci, idSesji, idOkna) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SpeechAudioUpload, {
          audio: bajtyBase64,
          contentType: typTresci,
          ...(idSesji === '' ? {} : { sessionId: idSesji }),
          ...(idOkna === '' ? {} : { windowId: idOkna }),
        }),
        Command.SpeechAudioUpload,
        (tresc) => czyTekst(tresc.audioRef),
      );
    },

    async mowaPrzepisz(odnosnik, jezyk) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SpeechTranscribe, {
          audioRef: odnosnik,
          ...(jezyk === '' ? {} : { language: jezyk }),
        }),
        Command.SpeechTranscribe,
        (tresc) => czyLogiczna(tresc.processed) && czyTekst(tresc.transcript),
      );
    },
  };
}

/** Czym mikrofon rozporządza po stronie karty. */
export interface ZapleczeMikrofonu {
  zrodlo: MowaZrodlo;
  /** Karta sesji, której nagranie dotyczy. */
  idSesji(): string;
  /** Okno, z którego nagranie wychodzi. */
  idOkna(): string;
  /** Oddaje rozpoznany tekst: polu polecenia albo treści dokumentu. */
  naTekst(tekst: string): void;
  /** Wypisuje zdanie o stanie — powodzenie albo odmowę nazwaną. */
  naZdanie(tresc: string, udane: boolean): void;
  /** Napis na przycisku w stanie spokoju; brak bierze „Mów". */
  etykieta?: string;
}

/** Przycisk mikrofonu wraz z jego stanem. */
export interface PrzyciskMikrofonu {
  element: HTMLButtonElement;
  /** Przerywa nagranie, gdy trwa — okno chowające wiersz polecenia woła to. */
  przerwij(): void;
  /** Czy nagranie trwa w tej chwili. */
  nagrywa(): boolean;
}

/** Typ treści nagrania preferowany; przeglądarka bez niego dostaje wybór własny. */
const TYP_NAGRANIA = 'audio/webm';

export function utworzPrzyciskMikrofonu(zaplecze: ZapleczeMikrofonu): PrzyciskMikrofonu {
  let rejestrator: MediaRecorder | null = null;
  let porcje: Blob[] = [];

  const spokoj = zaplecze.etykieta ?? 'Mów';

  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn dn-btn--sm dn-btn--zarys ms-mikrofon';
  element.textContent = spokoj;
  element.dataset['czynnosc'] = 'mikrofon';
  element.dataset['nagrywa'] = 'nie';
  element.setAttribute('aria-label', `${spokoj} — dyktowanie głosem`);
  element.title =
    'Nagranie idzie komendą speech.audio.upload do magazynu rdzenia i wraca tekstem ' +
    'z speech.transcribe. Dźwięk nie opuszcza maszyny rdzenia. Silnik mowy jest składnikiem ' +
    'pakietu serwera, więc dyktowanie jest funkcją, na którą można liczyć.';

  element.addEventListener('click', () => {
    if (rejestrator !== null) {
      zakoncz();
      return;
    }
    void rozpocznij();
  });

  /** Sprawdza silnik, prosi o mikrofon i rusza z nagraniem. */
  async function rozpocznij(): Promise<void> {
    const dostepnosc = await zaplecze.zrodlo.mowaDostepnosc();
    if (!dostepnosc.udany || dostepnosc.wynik === undefined) {
      zaplecze.naZdanie(
        opisOdmowy('Sprawdzenie silnika mowy', dostepnosc.blad?.code, dostepnosc.blad?.message),
        false,
      );
      return;
    }
    const stanSilnika = dostepnosc.wynik;
    if (!stanSilnika.available) {
      zaplecze.naZdanie(
        'TEN SERWER JEST NIEKOMPLETNY: brakuje mu składnika pakietu — silnika rozpoznawania ' +
          'mowy. Rdzeń podaje: ' +
          (stanSilnika.reason ?? 'powodu nie nazwał.') +
          ' Silnik mowy jedzie wraz z aplikacją na serwer, więc jego brak jest usterką wdrożenia ' +
          'serwera, a nie ograniczeniem produktu — uzupełnia go administrator pakietu, nie ' +
          'Operator. Do tego czasu wpisz treść klawiaturą; pole obok działa.',
        false,
      );
      return;
    }
    if (stanSilnika.uploadAvailable === false) {
      zaplecze.naZdanie(
        'Silnik mowy stoi, ale ten rdzeń nie przyjmuje nagrania w bajtach (pole uploadAvailable ' +
          'komendy speech.availability.get). Nagranie z mikrofonu istnieje wyłącznie jako bajty ' +
          'w pamięci karty, więc bez tej drogi nie ma czym go dostarczyć — to również usterka ' +
          'wdrożenia rdzenia, nie brak funkcji.',
        false,
      );
      return;
    }

    const strumien = await pobierzStrumien();
    if (strumien === null) return;

    porcje = [];
    const wybor = typeof MediaRecorder.isTypeSupported === 'function' &&
      MediaRecorder.isTypeSupported(TYP_NAGRANIA)
      ? { mimeType: TYP_NAGRANIA }
      : {};
    rejestrator = new MediaRecorder(strumien, wybor);
    rejestrator.addEventListener('dataavailable', (zdarzenie) => {
      if (zdarzenie.data.size > 0) porcje.push(zdarzenie.data);
    });
    rejestrator.addEventListener('stop', () => {
      for (const sciezka of strumien.getTracks()) sciezka.stop();
      void dostarcz();
    });
    rejestrator.start();
    element.textContent = 'Zakończ nagranie';
    element.dataset['nagrywa'] = 'tak';
    zaplecze.naZdanie(
      `Nagranie w toku. Silnik: ${stanSilnika.engine ?? 'nienazwany'}, model ` +
        `${stanSilnika.model ?? 'z katalogu ustawień'}. Naciśnij ponownie, żeby zakończyć.`,
      true,
    );
  }

  /**
   * Prosi przeglądarkę o mikrofon.
   *
   * Odmowa Operatora i brak urządzenia to dwie różne rzeczy i obie mają tu
   * własne zdanie — „mikrofon nie działa" nie powiedziałoby, czego szukać.
   */
  async function pobierzStrumien(): Promise<MediaStream | null> {
    if (navigator.mediaDevices === undefined) {
      zaplecze.naZdanie(
        'Ta przeglądarka nie wystawia urządzeń nagrywających (navigator.mediaDevices), więc ' +
          'nagrania nie ma czym zdjąć. Droga rdzenia jest gotowa — brak jest po stronie karty.',
        false,
      );
      return null;
    }
    try {
      return await navigator.mediaDevices.getUserMedia({ audio: true });
    } catch (powod) {
      zaplecze.naZdanie(
        'Mikrofon nie został udostępniony: ' +
          (powod instanceof Error ? powod.message : String(powod)) +
          '. Zgoda na mikrofon należy do przeglądarki, nie do rdzenia — wydaj ją i spróbuj ponownie.',
        false,
      );
      return null;
    }
  }

  function zakoncz(): void {
    if (rejestrator === null) return;
    if (rejestrator.state !== 'inactive') rejestrator.stop();
    rejestrator = null;
    element.textContent = spokoj;
    element.dataset['nagrywa'] = 'nie';
  }

  /** Wnosi nagranie do rdzenia i zleca jego przepisanie. */
  async function dostarcz(): Promise<void> {
    if (porcje.length === 0) {
      zaplecze.naZdanie('Nagranie jest puste — nie zdjęto ani jednej porcji dźwięku.', false);
      return;
    }
    const calosc = new Blob(porcje, { type: porcje[0]?.type ?? TYP_NAGRANIA });
    porcje = [];
    const bajty = await naBase64(calosc);
    zaplecze.naZdanie(`Nagranie ${calosc.size} bajtów idzie do magazynu rdzenia…`, true);

    const wniesienie = await zaplecze.zrodlo.mowaWniesNagranie(
      bajty,
      calosc.type === '' ? TYP_NAGRANIA : calosc.type,
      zaplecze.idSesji(),
      zaplecze.idOkna(),
    );
    if (!wniesienie.udany || wniesienie.wynik === undefined) {
      zaplecze.naZdanie(
        opisOdmowy('Wniesienie nagrania', wniesienie.blad?.code, wniesienie.blad?.message),
        false,
      );
      return;
    }

    const przepisanie = await zaplecze.zrodlo.mowaPrzepisz(wniesienie.wynik.audioRef, 'pl');
    if (!przepisanie.udany || przepisanie.wynik === undefined) {
      zaplecze.naZdanie(
        opisOdmowy('Przepisanie nagrania', przepisanie.blad?.code, przepisanie.blad?.message),
        false,
      );
      return;
    }
    const wynik = przepisanie.wynik;
    if (!wynik.processed) {
      zaplecze.naZdanie(
        'Rdzeń nagrania nie przetworzył (pole processed fałsz) — tekstu nie ma i okno go nie ' +
          'zmyśla.',
        false,
      );
      return;
    }
    if (wynik.transcript === '') {
      zaplecze.naZdanie(
        `Rdzeń przetworzył nagranie (${wynik.durationMs} ms, model ${wynik.model}) i mowy w nim ` +
          'nie było. Pusty wynik jest tu FAKTEM pomiaru, nie usterką — cisza i szum wyglądają tak.',
        false,
      );
      return;
    }
    zaplecze.naTekst(wynik.transcript);
    zaplecze.naZdanie(
      `Podyktowano ${wynik.characters} znaków (${wynik.durationMs} ms, model ${wynik.model}` +
        `${wynik.confidence === undefined ? '' : `, pewność ${wynik.confidence}/100`}). ` +
        'Poprawisz tekst w polu, zanim zlecisz.',
      true,
    );
  }

  return {
    element,
    przerwij: () => zakoncz(),
    nagrywa: () => rejestrator !== null,
  };
}

/**
 * Bajty nagrania w postaci base64, bez przedrostka schematu danych.
 *
 * Kontrakt `speech.audio.upload` żąda samych bajtów, a `FileReader` oddaje je
 * z przedrostkiem `data:…;base64,` — przedrostek schodzi tutaj, w jednym
 * miejscu, żeby nie jechał do rdzenia jako część nagrania.
 */
async function naBase64(dane: Blob): Promise<string> {
  const bufor = await dane.arrayBuffer();
  const bajty = new Uint8Array(bufor);
  let napis = '';
  // Porcjami, bo `String.fromCharCode` z rozwinięciem całej tablicy przekracza
  // dopuszczalną liczbę argumentów wywołania przy nagraniu dłuższym niż chwila.
  const porcja = 8192;
  for (let poczatek = 0; poczatek < bajty.length; poczatek += porcja) {
    napis += String.fromCharCode(...bajty.subarray(poczatek, poczatek + porcja));
  }
  return btoa(napis);
}
