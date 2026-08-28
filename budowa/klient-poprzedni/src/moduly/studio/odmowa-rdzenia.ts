import {
  ErrorCode,
  zdarzenieNieznanej,
  type Command,
  type ErrorInfo,
  type RequestOf,
  type ResponseOf,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';

/** Wywołanie komendy odporne na odmowę rdzenia bez uchwytu, łączące odpowiedź skorelowaną i zdarzenie odmowy w jeden wynik. */
export function wywolajUczciwie<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<Wynik<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    let idZadania = '';
    let zamkniete = false;

    function zakoncz(wynik: Wynik<ResponseOf<K>>): void {
      if (zamkniete) return;
      zamkniete = true;
      odsubskrybuj();
      rozstrzygnij(wynik);
    }

    // Subskrypcja przed wysyłką, inaczej natychmiastowa odmowa rdzenia wyprzedziłaby założenie nasłuchu.
    const odsubskrybuj = kanal.naZdarzenie(zdarzenieNieznanej(komenda), (tresc: unknown) => {
      const odmowa = odczytajOdmowe(tresc);
      if (odmowa === null) return;
      if (odmowa.idZadania !== '' && odmowa.idZadania !== idZadania) return;
      zakoncz({ udany: false, blad: bladOdmowy(komenda, odmowa) });
    });

    idZadania = kanal.wyslij(komenda, zadanie, zakoncz);
  });
}

/** Odmowa rdzenia odczytana z ładunku zdarzenia odmowy komendy bez uchwytu, wraz z żądanym typem i powodem. */
interface OdmowaRdzenia {
  /** Typ, którego rdzeń nie rozpoznał; pusty, gdy rdzeń go nie podał. */
  zadanyTyp: string;
  /** Identyfikator żądania, którego odmowa dotyczy. */
  idZadania: string;
  /** Powód podany przez rdzeń. */
  powod: string;
}

/** Zdanie odmowy nazywające żądany typ komendy, której rdzeń nie zna, z kodem odmowy niepowtarzalnej przez okno. */
function bladOdmowy(komenda: string, odmowa: OdmowaRdzenia): ErrorInfo {
  const typ = odmowa.zadanyTyp === '' ? komenda : odmowa.zadanyTyp;
  const powod = odmowa.powod === '' ? '' : ` (${odmowa.powod})`;
  return {
    code: ErrorCode.NotFound,
    message:
      `Rdzeń nie obsługuje jeszcze komendy ${typ}${powod}. ` +
      'Komenda jest w kontrakcie, uchwytu w rdzeniu nie ma — okno nie ma czym wykonać tej czynności.',
    retryable: false,
  };
}

/** Odczytuje treść odmowy z ładunku zdarzenia, odporny na jego brak albo na kształt inny niż oczekiwany. */
function odczytajOdmowe(ladunek: unknown): OdmowaRdzenia | null {
  if (typeof ladunek !== 'object' || ladunek === null) return null;
  const pola = ladunek as Record<string, unknown>;
  if (typeof pola['requestedType'] !== 'string') return null;
  return {
    zadanyTyp: tekst(pola['requestedType']),
    idZadania: tekst(pola['requestId']),
    powod: tekst(pola['reason']),
  };
}

/** Zwraca napis albo pusty łańcuch znaków, bo pole opcjonalne kontraktu bywa w odpowiedzi rdzenia nieobecne. */
function tekst(wartosc: unknown): string {
  return typeof wartosc === 'string' ? wartosc : '';
}
