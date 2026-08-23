import { opisOdmowy } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { wydajDokument } from './konwersja-dokumentu';
import type { StanOknaStudio } from './stan-okna-studio';
import type { StanStudio } from './stan-studio';
import type { ZrodloDokumentuStudio } from './zrodlo-dokumentu-studio';
import type { ZrodloPrzekazania } from './zrodlo-przekazania';

/**
 * Eksport i przekazanie do Library — dwie czynności Preview Window wyjęte
 * z wytwórni okna. Obie wychodzą do rdzenia i obie się wykonują.
 *
 * Eksport nie ma komendy w obszarze `studio` i mieć jej nie musi: zamiana
 * formatu dokumentu jest czynnością obszaru `document`, ma tam uchwyt w rdzeniu
 * i słownik ośmiu formatów. Przedmiotem zamiany jest treść zaakceptowana
 * modułu, więc żądanie niesie ją wprost, a nie ścieżkę pliku.
 *
 * Wynik zamiany zostaje zasobem magazynu rdzenia. Komendy wydającej jego bajty
 * do przeglądarki kontrakt nie niesie — okno mówi to wprost przy potwierdzeniu,
 * zamiast pozwalać czytać „wyeksportowano" jako „pobrano".
 *
 * Przekazanie idzie `context.transfer` i wykonuje się w całości: rdzeń zakłada
 * okno modułu Library i oddaje `transferred: true`.
 *
 * Odmowa zostaje w pasie stanu, powodzenie w wierszu odpowiedzi. Wskaźnik
 * odczytu zapala się wyłącznie na czas rzeczywistego wywołania — warunek
 * merytoryczny sprawdzany przed wysłaniem żadnego nie udaje.
 */
export interface ZapleczePodgladu {
  stan: StanStudio;
  /** Źródło komend obszaru dokumentów — nośnik zamiany formatu. */
  dokumenty: ZrodloDokumentuStudio;
  /** Pas stanu okna — nośnik ładowania i odmowy. */
  pas: StanOknaStudio;
  odpowiedz: WierszOdpowiedzi;
  /** Przywraca pasowi stan wynikający z danych, gdy komunikat ustępuje. */
  odswiez(): void;
}

export async function eksportujDokument(
  zaplecze: ZapleczePodgladu,
  format: string,
): Promise<void> {
  const { stan, dokumenty, pas, odpowiedz } = zaplecze;
  pas.ladowanie(`Wydanie dokumentu w formacie ${format} w toku…`);
  const udane = await wydajDokument({ stan, dokumenty, odpowiedz }, format);
  if (!udane) {
    // Powód stoi już w wierszu odpowiedzi — pas wraca do stanu wynikającego
    // z danych, żeby jeden powód nie był pokazywany dwa razy w dwóch miejscach.
    pas.gotowe();
    zaplecze.odswiez();
    return;
  }
  pas.gotowe();
  zaplecze.odswiez();
}

/**
 * Przekazanie dokumentu do Library kompletem kontekstu.
 *
 * Komplet niesie identyfikator dokumentu, bo to on jest przedmiotem
 * przekazania; polecenie wyjściowe zostaje puste, bo podgląd nie zleca pracy,
 * tylko przenosi wynik.
 */
export async function przekazDoLibrary(
  zaplecze: ZapleczePodgladu,
  przekazanie: ZrodloPrzekazania,
): Promise<void> {
  const { stan, pas, odpowiedz } = zaplecze;
  const dokument = stan.dokument();
  if (dokument === null) {
    odpowiedz.pokaz('Przekazanie dotyczy dokumentu — wczytaj go najpierw w Studio Editorze.', false);
    return;
  }
  pas.ladowanie('Przekazywanie kontekstu do Library…');
  odpowiedz.pokaz('Przekazywanie kontekstu do Library…', true);
  const wynik = await przekazanie.doLibrary(stan.idOkna(), { documentIds: [dokument.id] });
  if (!wynik.udany || wynik.wynik === undefined) {
    const opis = opisOdmowy('Przekazanie do Library', wynik.blad?.code, wynik.blad?.message);
    pas.blad(opis);
    odpowiedz.pokaz(opis, false);
    return;
  }
  pas.gotowe();
  zaplecze.odswiez();
  odpowiedz.pokaz(
    `Kontekst przeniesiony do okna ${wynik.wynik.window.id} modułu Library ` +
      `(przeniesiono: ${String(wynik.wynik.transferred)}).`,
    true,
  );
}
