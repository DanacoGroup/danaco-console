import { Command } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { Kanal } from '../../protokol/kanal';
import { wywolaj } from '../../protokol/wywolanie';
import type { Dostarczenie } from './dostarczenie-nagrania';

/**
 * Odpowiedź na pytanie, czy dyktowanie jest teraz dostępne: koniunkcja gotowości
 * silnika po stronie rdzenia i drogi dostarczenia nagrania po stronie klienta.
 */
export interface DostepnoscDyktowania {
  /** Koniunkcja obu warunków — dopiero ona zapala mikrofon. */
  dostepne: boolean;
  /** Wersja Pythona znaleziona przez rdzeń; pusta, gdy nie znaleziono. */
  python: string;
  /** Wersja silnika faster-whisper; pusta, gdy brak. */
  engine: string;
  /** Model ustawiony w katalogu ustawień; pusty, gdy nieustawiony. */
  model: string;
  /** Powód niedostępności — trzyczęściowy: co, dlaczego, czym Operator to zmieni. */
  powod: string;
}

export interface ZrodloDostepnosci {
  /** Pierwsze wywołanie pyta rdzeń, kolejne oddają ten sam odczyt. */
  odczytaj(): Promise<DostepnoscDyktowania>;
  /** Unieważnia odczyt — następne `odczytaj()` zapyta rdzeń ponownie. */
  odswiez(): void;
}

/**
 * Zdanie trzyczęściowe o braku drogi dostarczenia nagrania: co się nie stało,
 * dlaczego, i jaka zmiana kontraktu to naprawi.
 */
const ZDANIE_BRAKU_DOSTARCZENIA =
  'Dyktowanie jest wyłączone po stronie klienta. ' +
  'Nagranie z mikrofonu powstaje jako bajty w pamięci karty, a speech.transcribe przyjmuje audioRef, ' +
  'czyli ścieżkę pliku na maszynie silnika — komendy przenoszącej nagranie kontrakt dziś nie ma. ' +
  'Zmieni to dopisanie do kontraktu komendy przyjmującej nagranie w base64 wraz z uchwytem w rdzeniu.';

/**
 * Zdanie trzyczęściowe o braku silnika mowy po stronie rdzenia, używane, gdy rdzeń
 * sam nie podał własnego powodu odmowy.
 */
const ZDANIE_BRAKU_SILNIKA =
  'Silnik mowy nie jest gotowy po stronie rdzenia. ' +
  'Rdzeń odpowiedział na speech.availability.get, że transkrypcja nie jest tu i teraz wykonalna, ' +
  'a powodu nie nazwał — brakuje Pythona, pakietu faster-whisper albo modelu. ' +
  'Zmieni to instalacja pomocnika mowy na maszynie rdzenia i wskazanie modelu w katalogu ustawień.';

export function utworzZrodloDostepnosci(kanal: Kanal, dostarczenie: Dostarczenie): ZrodloDostepnosci {
  let odczyt: Promise<DostepnoscDyktowania> | null = null;

  async function pobierz(): Promise<DostepnoscDyktowania> {
    const odpowiedz = await wywolaj(kanal, Command.SpeechAvailabilityGet, {});

    // Drugi warunek sprawdzany zawsze — także przy odmowie rdzenia — by oba braki wyszły naraz.
    const drogaJest = dostarczenie.dostepne();
    const powody: string[] = [];

    if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
      // Odmowa komendy, w tym nieznana: pola puste, powód nazywa odmowę, nie brak silnika.
      powody.push(opisOdmowyBledu('Odczyt dostępności silnika mowy (speech.availability.get)', odpowiedz.blad));
      if (!drogaJest) powody.push(ZDANIE_BRAKU_DOSTARCZENIA);
      return { dostepne: false, python: '', engine: '', model: '', powod: powody.join(' ') };
    }

    const tresc = odpowiedz.wynik;
    if (!tresc.available) {
      const wlasny = (tresc.reason ?? '').trim();
      powody.push(wlasny === '' ? ZDANIE_BRAKU_SILNIKA : wlasny);
    }
    if (!drogaJest) powody.push(ZDANIE_BRAKU_DOSTARCZENIA);

    return {
      // Koniunkcja: jedno „nie" po którejkolwiek stronie gasi mikrofon.
      dostepne: tresc.available && drogaJest,
      python: tresc.python ?? '',
      engine: tresc.engine ?? '',
      model: tresc.model ?? '',
      powod: powody.join(' '),
    };
  }

  return {
    odczytaj() {
      odczyt ??= pobierz();
      return odczyt;
    },
    odswiez() {
      // Samo unieważnienie: nie każda zmiana ustawień kończy się otwarciem okna dyktowania.
      odczyt = null;
    },
  };
}
