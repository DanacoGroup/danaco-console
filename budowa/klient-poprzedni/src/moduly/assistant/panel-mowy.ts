import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { opisOdmowy } from '../../komponenty/odmowa';
import { pole, przyciskAkcji, wybor } from '../../modele/kontrolki-formularza';
import { KLASY_DYMKA, ODCZYTY } from './etykiety-assistant';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAssistant } from './stan-assistant';

/**
 * Rozpoznawanie mowy w Voice Console: `speech.availability.get` mówi, czy silnik
 * stoi na maszynie, a `speech.transcribe` zamienia nagranie w tekst polecenia.
 */
export interface PanelMowy {
  element: HTMLElement;
  /** Odczyt dostępności silnika mowy. */
  wczytaj(): Promise<void>;
  /** Rozpoczyna albo kończy nagranie mikrofonem — droga z paska promptu. */
  przelaczMikrofon(): void;
}

/**
 * Wykaz języków wskazywanych wprost przy rozpoznaniu nagrania; wartość pusta
 * zleca silnikowi rozpoznanie języka samodzielnie.
 */
const JEZYKI: ReadonlyArray<readonly [string, string]> = [
  ['', 'Język nagrania: rozpoznaj automatycznie'],
  ['pl', 'Język nagrania: polski'],
  ['en', 'Język nagrania: angielski'],
];

export function utworzPanelMowy(
  stan: StanAssistant,
  naTranskrypcje: (tresc: string) => void,
): PanelMowy {
  const okno: StanOkna = utworzStanOkna();

  const sciezka = pole(
    'Ścieżka nagrania na maszynie silnika',
    'ścieżka pliku widziana przez maszynę silnika mowy',
  );
  const jezyk = wybor('Język nagrania', JEZYKI);

  const rozpoznaj = przyciskAkcji('Rozpoznaj nagranie', 'dn-btn dn-btn--sm dn-btn--zarys');
  rozpoznaj.addEventListener('click', () => void transkrybuj());

  // Nagranie z mikrofonu wchodzi komendą `speech.audio.upload`, oddającą odnosnik.
  const mikrofon = przyciskAkcji('Nagraj z mikrofonu', 'dn-btn dn-btn--sm dn-btn--zarys');
  mikrofon.addEventListener('click', () => void przelaczNagrywanie());

  const kontrolki = document.createElement('div');
  kontrolki.className = 'ma-mowa__kontrolki';
  kontrolki.append(sciezka, jezyk, rozpoznaj, mikrofon);

  // Panel nie stawia wskaźnika pewności rozpoznania, ponieważ kontrakt
  // miary pewności nie oddaje.
  const pewnosc = document.createElement('p');
  pewnosc.className = 'dn-pole-opis';
  pewnosc.textContent =
    'Rdzeń nie oddaje miary pewności rozpoznania — po transkrypcji przeczytaj tekst ' +
    'w polu polecenia i popraw go przed wysłaniem.';

  okno.tresc.append(kontrolki, pewnosc);

  const tytul = document.createElement('h4');
  tytul.className = 'ma-panel__tytul';
  tytul.append(
    'Rozpoznawanie mowy',
    utworzDymekObjasnienia(
      'Silnik lokalny rdzenia (speech.transcribe). Pole audioRef kontraktu to ścieżka ' +
        'pliku na maszynie silnika — dźwięk nie opuszcza tej maszyny, a nagranie z ' +
        'mikrofonu karty nie ma dziś czym do niej dojechać.',
      KLASY_DYMKA,
    ),
  );

  const element = document.createElement('div');
  element.className = 'ma-panel';
  element.dataset['panel'] = 'mowa';
  element.append(tytul, okno.element);

  /** Zdanie o silniku: co jest, czego nie ma i czym Operator to zmieni. */
  async function wczytaj(): Promise<void> {
    okno.ladowanie(ODCZYTY.silnikMowy);
    const wynik = await stan.zaplecze.silnikMowy();
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Odczyt dostępności silnika mowy', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    const stanSilnika = wynik.wynik;
    if (!stanSilnika.available) {
      // Niedostępność silnika nie jest odmową: stan pusty niesie powód rdzenia.
      okno.puste(
        stanSilnika.reason !== undefined && stanSilnika.reason !== ''
          ? `Silnik mowy nie jest gotów: ${stanSilnika.reason}`
          : 'Silnik mowy nie jest gotów, a rdzeń nie podał powodu. ' +
              'Polecenie wydasz polem transkrypcji powyżej.',
      );
      return;
    }
    okno.gotowe();
  }

  /** Nagrywanie w toku; pusty znaczy, że mikrofon stoi. */
  let nagrywanie: MediaRecorder | undefined;

  /** Nagranie z mikrofonu karty i jego droga do rdzenia oraz z powrotem. */
  async function przelaczNagrywanie(): Promise<void> {
    if (nagrywanie !== undefined) {
      nagrywanie.stop();
      return;
    }
    if (navigator.mediaDevices === undefined) {
      okno.blad(
        'Ta karta nie ma dostępu do urządzeń nagrywających (brak navigator.mediaDevices). ' +
          'Nagranie wskażesz ścieżką powyżej albo otworzysz Danaco Console w programie okiennym.',
      );
      return;
    }
    let strumien: MediaStream;
    try {
      strumien = await navigator.mediaDevices.getUserMedia({ audio: true });
    } catch (powod) {
      okno.blad(
        `Mikrofon nie został udostępniony: ${String(powod)}. ` +
          'Zgody udziela przeglądarka, nie rdzeń.',
      );
      return;
    }

    const odcinki: Blob[] = [];
    const rejestrator = new MediaRecorder(strumien);
    nagrywanie = rejestrator;
    mikrofon.textContent = 'Zakończ nagranie';

    rejestrator.addEventListener('dataavailable', (zdarzenie) => {
      if (zdarzenie.data.size > 0) odcinki.push(zdarzenie.data);
    });
    rejestrator.addEventListener('stop', () => {
      nagrywanie = undefined;
      mikrofon.textContent = 'Nagraj z mikrofonu';
      for (const sciezkaWejscia of strumien.getTracks()) sciezkaWejscia.stop();
      void przeslij(new Blob(odcinki, { type: rejestrator.mimeType }));
    });
    rejestrator.start();
    okno.ladowanie('Nagrywanie… naciśnij ponownie, żeby zakończyć.');
  }

  /** Wysyła nagranie do rdzenia i od razu prosi o jego rozpoznanie. */
  async function przeslij(nagranie: Blob): Promise<void> {
    if (nagranie.size === 0) {
      okno.puste('Nagranie jest puste — mikrofon nie oddał ani jednego bajtu.');
      return;
    }
    okno.ladowanie('Przesyłanie nagrania do rdzenia…');
    const bajty = new Uint8Array(await nagranie.arrayBuffer());
    let zapis = '';
    for (const bajt of bajty) zapis += String.fromCharCode(bajt);

    const wynik = await stan.mowa.przeslijNagranie({
      base64: btoa(zapis),
      typTresci: nagranie.type === '' ? 'audio/webm' : nagranie.type,
      idOkna: stan.idOkna(),
      idSesji: stan.idSesji(),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Przesłanie nagrania', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    sciezka.value = wynik.wynik.audioRef;
    await transkrybuj();
  }

  async function transkrybuj(): Promise<void> {
    const wskazana = sciezka.value.trim();
    if (wskazana === '') {
      okno.blad(
        'Wskaż ścieżkę nagrania na maszynie silnika. Kontrakt przyjmuje w polu audioRef ' +
          'ścieżkę pliku, nie sam dźwięk, więc bez niej rdzeń nie ma czego rozpoznać.',
      );
      return;
    }
    okno.ladowanie(ODCZYTY.transkrypcja);
    const zadanie = { audioRef: wskazana };
    const wynik = await stan.zaplecze.transkrypcja(
      jezyk.value === '' ? zadanie : { ...zadanie, language: jezyk.value },
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Rozpoznanie nagrania', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    const rozpoznane = wynik.wynik;
    if (rozpoznane.transcript === '') {
      // Cisza jest wynikiem, nie usterką, więc pole polecenia zostaje bez zmian.
      okno.puste(
        `Silnik przetworzył nagranie (${String(rozpoznane.durationMs)} ms, model ` +
          `${rozpoznane.model}) i nie rozpoznał w nim mowy. Pole polecenia zostaje bez zmian.`,
      );
      return;
    }
    naTranskrypcje(rozpoznane.transcript);
    okno.gotowe();
  }

  return {
    element,
    wczytaj,
    przelaczMikrofon: () => {
      void przelaczNagrywanie();
    },
  };
}
